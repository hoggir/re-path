"""RabbitMQ RPC consumer for handling dashboard requests."""

import asyncio
import json
import logging
from concurrent.futures import ThreadPoolExecutor
from typing import Callable, Optional

import pika
from pika.channel import Channel
from pika.spec import Basic, BasicProperties

from app.core.config import settings

logger = logging.getLogger(__name__)


class RabbitMQRPCConsumer:
    """RabbitMQ RPC consumer for handling dashboard requests."""

    def __init__(self) -> None:
        """Initialize RabbitMQ RPC consumer."""
        self._connection: Optional[pika.BlockingConnection] = None
        self._channel: Optional[Channel] = None
        self._consumer_tag: Optional[str] = None
        self._request_handler: Optional[Callable] = None
        self._closing = False
        self._consuming = False
        self._thread_pool = ThreadPoolExecutor(max_workers=1)

    def set_request_handler(self, handler: Callable) -> None:
        """
        Set the request handler callback.

        Args:
            handler: Async function to handle incoming RPC requests
        """
        self._request_handler = handler

    def connect(self) -> None:
        """Connect to RabbitMQ server."""
        try:
            # Parse URL and create connection parameters
            parameters = pika.URLParameters(settings.rabbitmq_url)
            parameters.heartbeat = 600
            parameters.blocked_connection_timeout = 300

            self._connection = pika.BlockingConnection(parameters)
            logger.info("✅ RPC Consumer connected to RabbitMQ")
        except Exception as e:
            logger.error(f"❌ Failed to connect to RabbitMQ: {e}")
            raise

    def _setup_channel(self) -> None:
        """Set up the RabbitMQ channel and queue for RPC."""
        if not self._connection:
            raise RuntimeError("RabbitMQ connection is not established")

        logger.info("Creating RabbitMQ RPC channel")
        self._channel = self._connection.channel()

        # Declare RPC queue
        self._channel.queue_declare(
            queue=settings.rabbitmq_rpc_queue,
            durable=True,
        )

        # Set QoS to prefetch 1 message at a time for fair dispatch
        self._channel.basic_qos(prefetch_count=1)

        logger.info(f"✅ RPC Channel setup complete - Queue: {settings.rabbitmq_rpc_queue}")

    def start_consuming(self) -> None:
        """Start consuming RPC requests from the queue."""
        if not self._connection:
            self.connect()

        self._setup_channel()

        if not self._channel:
            raise RuntimeError("RabbitMQ channel is not set up")

        logger.info(f"Starting to consume RPC requests from queue: {settings.rabbitmq_rpc_queue}")

        def callback_wrapper(
            ch: Channel,
            method: Basic.Deliver,
            properties: BasicProperties,
            body: bytes,
        ) -> None:
            """Wrapper for RPC callback to run async handler."""
            try:
                # Use asyncio.run to create isolated event loop
                asyncio.run(self._on_rpc_request(ch, method, properties, body))
            except Exception as e:
                logger.error(f"Error in RPC callback wrapper: {e}")

        self._consumer_tag = self._channel.basic_consume(
            queue=settings.rabbitmq_rpc_queue,
            on_message_callback=callback_wrapper,
            auto_ack=False,
        )

        self._consuming = True
        logger.info(f"✅ Started consuming RPC requests - Consumer tag: {self._consumer_tag}")

        # Start consuming in a blocking manner
        try:
            self._channel.start_consuming()
        except Exception as e:
            if not self._closing:
                logger.error(f"Error while consuming: {e}")
            self.stop_consuming()

    async def _on_rpc_request(
        self,
        channel: Channel,
        method: Basic.Deliver,
        properties: BasicProperties,
        body: bytes,
    ) -> None:
        """
        Handle incoming RPC request.

        Args:
            channel: The channel object
            method: Delivery method
            properties: Message properties
            body: Message body
        """
        correlation_id = properties.correlation_id
        reply_to = properties.reply_to

        try:
            # Decode and parse request
            message = body.decode("utf-8")
            logger.info(f"📨 Received RPC request - correlation_id: {correlation_id}")
            logger.debug(f"Request body: {message[:200]}...")

            request_data = json.loads(message)

            # Process request through handler
            if self._request_handler:
                response_data = await self._request_handler(request_data)
            else:
                logger.error("No request handler set for RPC consumer")
                response_data = {"error": "Service not configured"}

            # Convert response to JSON
            response_body = json.dumps(response_data).encode("utf-8")

            # Send response back to the reply queue
            if reply_to:
                channel.basic_publish(
                    exchange="",
                    routing_key=reply_to,
                    properties=pika.BasicProperties(
                        correlation_id=correlation_id,
                        content_type="application/json",
                    ),
                    body=response_body,
                )
                logger.info(f"✅ RPC response sent - correlation_id: {correlation_id}")
            else:
                logger.warning(f"No reply_to queue specified for correlation_id: {correlation_id}")

            # Acknowledge the request
            channel.basic_ack(delivery_tag=method.delivery_tag)

        except json.JSONDecodeError as e:
            logger.error(f"❌ Failed to decode JSON request: {e}")
            # Send error response
            if reply_to and correlation_id:
                error_response = json.dumps({"error": "Invalid JSON format"}).encode("utf-8")
                channel.basic_publish(
                    exchange="",
                    routing_key=reply_to,
                    properties=pika.BasicProperties(correlation_id=correlation_id),
                    body=error_response,
                )
            # Reject and don't requeue - bad format will never work
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

        except Exception as e:
            logger.error(f"❌ Error processing RPC request: {e}", exc_info=True)
            # Send error response
            if reply_to and correlation_id:
                error_response = json.dumps({"error": str(e)}).encode("utf-8")
                channel.basic_publish(
                    exchange="",
                    routing_key=reply_to,
                    properties=pika.BasicProperties(correlation_id=correlation_id),
                    body=error_response,
                )
            # Reject message but don't requeue
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    def stop_consuming(self) -> None:
        """Stop consuming RPC requests."""
        if self._channel and self._consumer_tag and self._consuming:
            logger.info("Stopping RPC consumer")
            try:
                self._channel.basic_cancel(self._consumer_tag)
                self._consuming = False
            except Exception as e:
                logger.error(f"Error stopping RPC consumer: {e}")

    def close(self) -> None:
        """Close the RabbitMQ connection."""
        logger.info("Closing RabbitMQ RPC connection")
        self._closing = True

        self.stop_consuming()

        if self._channel:
            try:
                self._channel.close()
            except Exception as e:
                logger.error(f"Error closing channel: {e}")

        if self._connection and self._connection.is_open:
            try:
                self._connection.close()
                logger.info("✅ RabbitMQ RPC connection closed")
            except Exception as e:
                logger.error(f"Error closing connection: {e}")


class RabbitMQRPCConsumerManager:
    """Manager for RabbitMQ RPC consumer instance."""

    consumer: Optional[RabbitMQRPCConsumer] = None
    _consumer_task: Optional[asyncio.Task] = None

    @classmethod
    async def start_consumer(cls, request_handler: Callable) -> None:
        """
        Start the RabbitMQ RPC consumer.

        Args:
            request_handler: Async function to handle incoming RPC requests
        """
        if cls.consumer is not None:
            logger.warning("RPC consumer is already running")
            return

        logger.info("Starting RabbitMQ RPC consumer manager")
        cls.consumer = RabbitMQRPCConsumer()
        cls.consumer.set_request_handler(request_handler)

        # Start consuming in background thread
        loop = asyncio.get_event_loop()
        cls._consumer_task = loop.run_in_executor(None, cls.consumer.start_consuming)
        logger.info("✅ RabbitMQ RPC consumer started in background")

    @classmethod
    async def stop_consumer(cls) -> None:
        """Stop the RabbitMQ RPC consumer."""
        if cls.consumer is None:
            logger.warning("No RPC consumer to stop")
            return

        logger.info("Stopping RabbitMQ RPC consumer manager")

        # Close consumer in executor to avoid blocking
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, cls.consumer.close)

        if cls._consumer_task and not cls._consumer_task.done():
            try:
                await asyncio.wait_for(cls._consumer_task, timeout=5.0)
            except asyncio.TimeoutError:
                logger.warning("RPC consumer task timeout, forcing stop")

        cls.consumer = None
        cls._consumer_task = None
        logger.info("✅ RabbitMQ RPC consumer stopped")

    @classmethod
    def is_running(cls) -> bool:
        """Check if RPC consumer is running."""
        return cls.consumer is not None
