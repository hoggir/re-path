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
    def __init__(self, main_loop: asyncio.AbstractEventLoop) -> None:
        self._connection: Optional[pika.BlockingConnection] = None
        self._channel: Optional[Channel] = None
        self._consumer_tag: Optional[str] = None
        self._request_handler: Optional[Callable] = None
        self._closing = False
        self._consuming = False
        self._thread_pool = ThreadPoolExecutor(max_workers=1)
        self._main_loop = main_loop

    def set_request_handler(self, handler: Callable) -> None:
        self._request_handler = handler

    def connect(self) -> None:
        try:
            parameters = pika.URLParameters(settings.rabbitmq_url)
            parameters.heartbeat = 600
            parameters.blocked_connection_timeout = 300

            self._connection = pika.BlockingConnection(parameters)
            logger.info("✅ RPC Consumer connected to RabbitMQ")
        except Exception as e:
            logger.error(f"❌ Failed to connect to RabbitMQ: {e}")
            raise

    def _setup_channel(self) -> None:
        if not self._connection:
            raise RuntimeError("RabbitMQ connection is not established")

        logger.info("Creating RabbitMQ RPC channel")
        self._channel = self._connection.channel()

        self._channel.queue_declare(
            queue=settings.rabbitmq_rpc_queue,
            durable=True,
        )

        self._channel.basic_qos(prefetch_count=1)

        logger.info(f"✅ RPC Channel setup complete - Queue: {settings.rabbitmq_rpc_queue}")

    def start_consuming(self) -> None:
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
            try:
                future = asyncio.run_coroutine_threadsafe(
                    self._on_rpc_request(ch, method, properties, body), self._main_loop
                )
                future.result()
            except Exception as e:
                logger.error(f"Error in RPC callback wrapper: {e}", exc_info=True)

        self._consumer_tag = self._channel.basic_consume(
            queue=settings.rabbitmq_rpc_queue,
            on_message_callback=callback_wrapper,
            auto_ack=False,
        )

        self._consuming = True
        logger.info(f"✅ Started consuming RPC requests - Consumer tag: {self._consumer_tag}")

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
        correlation_id = properties.correlation_id
        reply_to = properties.reply_to

        try:
            message = body.decode("utf-8")
            logger.info(f"📨 Received RPC request - correlation_id: {correlation_id}")
            logger.debug(f"Request body: {message[:200]}...")

            request_data = json.loads(message)

            if self._request_handler:
                response_data = await self._request_handler(request_data)
            else:
                logger.error("No request handler set for RPC consumer")
                response_data = {"error": "Service not configured"}

            response_body = json.dumps(response_data).encode("utf-8")

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

            channel.basic_ack(delivery_tag=method.delivery_tag)

        except json.JSONDecodeError as e:
            logger.error(f"❌ Failed to decode JSON request: {e}")
            if reply_to and correlation_id:
                error_response = json.dumps({"error": "Invalid JSON format"}).encode("utf-8")
                channel.basic_publish(
                    exchange="",
                    routing_key=reply_to,
                    properties=pika.BasicProperties(correlation_id=correlation_id),
                    body=error_response,
                )
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

        except Exception as e:
            logger.error(f"❌ Error processing RPC request: {e}", exc_info=True)
            if reply_to and correlation_id:
                error_response = json.dumps({"error": str(e)}).encode("utf-8")
                channel.basic_publish(
                    exchange="",
                    routing_key=reply_to,
                    properties=pika.BasicProperties(correlation_id=correlation_id),
                    body=error_response,
                )
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    def stop_consuming(self) -> None:
        if self._channel and self._consumer_tag and self._consuming:
            logger.info("Stopping RPC consumer")
            try:
                self._channel.basic_cancel(self._consumer_tag)
                self._consuming = False
            except Exception as e:
                logger.error(f"Error stopping RPC consumer: {e}")

    def close(self) -> None:
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
    consumer: Optional[RabbitMQRPCConsumer] = None
    _consumer_task: Optional[asyncio.Task] = None

    @classmethod
    async def start_consumer(cls, request_handler: Callable) -> None:
        if cls.consumer is not None:
            logger.warning("RPC consumer is already running")
            return

        logger.info("Starting RabbitMQ RPC consumer manager")

        loop = asyncio.get_event_loop()

        cls.consumer = RabbitMQRPCConsumer(main_loop=loop)
        cls.consumer.set_request_handler(request_handler)

        cls._consumer_task = loop.run_in_executor(None, cls.consumer.start_consuming)
        logger.info("✅ RabbitMQ RPC consumer started in background")

    @classmethod
    async def stop_consumer(cls) -> None:
        if cls.consumer is None:
            logger.warning("No RPC consumer to stop")
            return

        logger.info("Stopping RabbitMQ RPC consumer manager")

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
        return cls.consumer is not None
