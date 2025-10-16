"""RabbitMQ consumer manager for analytics events."""

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


class RabbitMQConsumer:
    """RabbitMQ consumer for analytics events using threading."""

    def __init__(self) -> None:
        """Initialize RabbitMQ consumer."""
        self._connection: Optional[pika.BlockingConnection] = None
        self._channel: Optional[Channel] = None
        self._consumer_tag: Optional[str] = None
        self._message_handler: Optional[Callable] = None
        self._closing = False
        self._consuming = False
        self._thread_pool = ThreadPoolExecutor(max_workers=1)

    def set_message_handler(self, handler: Callable) -> None:
        """
        Set the message handler callback.

        Args:
            handler: Async function to handle incoming messages
        """
        self._message_handler = handler

    def connect(self) -> None:
        """Connect to RabbitMQ server."""

        try:
            # Parse URL and create connection parameters
            parameters = pika.URLParameters(settings.rabbitmq_url)
            parameters.heartbeat = 600
            parameters.blocked_connection_timeout = 300

            self._connection = pika.BlockingConnection(parameters)
            logger.info("✅ Connected to RabbitMQ")
        except Exception as e:
            logger.error(f"❌ Failed to connect to RabbitMQ: {e}")
            raise

    def _setup_channel(self) -> None:
        """Set up the RabbitMQ channel and queue."""
        if not self._connection:
            raise RuntimeError("RabbitMQ connection is not established")

        logger.info("Creating RabbitMQ channel")
        self._channel = self._connection.channel()

        # Declare exchange
        self._channel.exchange_declare(
            exchange=settings.rabbitmq_exchange,
            exchange_type="topic",
            durable=True,
        )

        # Declare queue
        self._channel.queue_declare(
            queue=settings.rabbitmq_queue,
            durable=True,
        )

        # Bind queue to exchange
        self._channel.queue_bind(
            exchange=settings.rabbitmq_exchange,
            queue=settings.rabbitmq_queue,
            routing_key=settings.rabbitmq_routing_key,
        )

        # Set QoS to prefetch 1 message at a time
        self._channel.basic_qos(prefetch_count=1)

        logger.info(
            f"✅ Channel setup complete - Queue: {settings.rabbitmq_queue}, "
        )

    def start_consuming(self) -> None:
        """Start consuming messages from the queue."""
        if not self._connection:
            self.connect()

        self._setup_channel()

        if not self._channel:
            raise RuntimeError("RabbitMQ channel is not set up")

        logger.info(f"Starting to consume messages from queue: {settings.rabbitmq_queue}")

        def callback_wrapper(
            ch: Channel,
            method: Basic.Deliver,
            properties: BasicProperties,
            body: bytes,
        ) -> None:
            """Wrapper for message callback to run async handler."""
            try:
                # Use asyncio.run to create isolated event loop
                asyncio.run(self._on_message(ch, method, properties, body))
            except Exception as e:
                logger.error(f"Error in callback wrapper: {e}")

        self._consumer_tag = self._channel.basic_consume(
            queue=settings.rabbitmq_queue,
            on_message_callback=callback_wrapper,
            auto_ack=False,
        )

        self._consuming = True
        logger.info(f"✅ Started consuming - Consumer tag: {self._consumer_tag}")

        # Start consuming in a blocking manner
        try:
            self._channel.start_consuming()
        except Exception as e:
            if not self._closing:
                logger.error(f"Error while consuming: {e}")
            self.stop_consuming()

    async def _on_message(
        self,
        channel: Channel,
        method: Basic.Deliver,
        properties: BasicProperties,
        body: bytes,
    ) -> None:
        """
        Handle incoming message.

        Args:
            channel: The channel object
            method: Delivery method
            properties: Message properties
            body: Message body
        """
        try:
            # Decode message
            message = body.decode("utf-8")
            logger.info(f"📨 Received message: {message[:100]}...")

            # Parse JSON
            data = json.loads(message)

            # Call the message handler if set
            if self._message_handler:
                await self._message_handler(data)
            else:
                logger.warning("No message handler set, message will be acknowledged but not processed")

            # Acknowledge the message
            channel.basic_ack(delivery_tag=method.delivery_tag)
            logger.info("✅ Message acknowledged")

        except json.JSONDecodeError as e:
            logger.error(f"❌ Failed to decode JSON message: {e}")
            # Reject message and don't requeue - bad format will never work
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

        except Exception as e:
            logger.error(f"❌ Error processing message: {e}", exc_info=True)
            # Reject message but don't requeue to prevent infinite loop
            # In production, you might want to send to a dead letter queue instead
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    def stop_consuming(self) -> None:
        """Stop consuming messages."""
        if self._channel and self._consumer_tag and self._consuming:
            logger.info("Stopping consumer")
            try:
                self._channel.basic_cancel(self._consumer_tag)
                self._consuming = False
            except Exception as e:
                logger.error(f"Error stopping consumer: {e}")

    def close(self) -> None:
        """Close the RabbitMQ connection."""
        logger.info("Closing RabbitMQ connection")
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
                logger.info("✅ RabbitMQ connection closed")
            except Exception as e:
                logger.error(f"Error closing connection: {e}")


class RabbitMQConsumerManager:
    """Manager for RabbitMQ consumer instance."""

    consumer: Optional[RabbitMQConsumer] = None
    _consumer_task: Optional[asyncio.Task] = None

    @classmethod
    async def start_consumer(cls, message_handler: Callable) -> None:
        """
        Start the RabbitMQ consumer.

        Args:
            message_handler: Async function to handle incoming messages
        """
        if cls.consumer is not None:
            logger.warning("Consumer is already running")
            return

        logger.info("Starting RabbitMQ consumer manager")
        cls.consumer = RabbitMQConsumer()
        cls.consumer.set_message_handler(message_handler)

        # Start consuming in background thread
        loop = asyncio.get_event_loop()
        cls._consumer_task = loop.run_in_executor(None, cls.consumer.start_consuming)
        logger.info("✅ RabbitMQ consumer started in background")

    @classmethod
    async def stop_consumer(cls) -> None:
        """Stop the RabbitMQ consumer."""
        if cls.consumer is None:
            logger.warning("No consumer to stop")
            return

        logger.info("Stopping RabbitMQ consumer manager")

        # Close consumer in executor to avoid blocking
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, cls.consumer.close)

        if cls._consumer_task and not cls._consumer_task.done():
            try:
                await asyncio.wait_for(cls._consumer_task, timeout=5.0)
            except asyncio.TimeoutError:
                logger.warning("Consumer task timeout, forcing stop")

        cls.consumer = None
        cls._consumer_task = None
        logger.info("✅ RabbitMQ consumer stopped")

    @classmethod
    def is_running(cls) -> bool:
        """Check if consumer is running."""
        return cls.consumer is not None
