"""Message handler service for RabbitMQ consumer."""

import json
import logging
from typing import Any

import requests

from app.core.config import settings

logger = logging.getLogger(__name__)


class MessageHandler:
    """Handler for processing RabbitMQ messages."""

    @staticmethod
    def index_to_opensearch_sync(document: dict[str, Any], index_type: str) -> None:
        """
        Index document to OpenSearch using synchronous HTTP request.

        Args:
            document: Document to index
            index_type: Index type (e.g., "click_events")
        """
        # Parse OpenSearch URL
        opensearch_url = settings.opensearch_url

        # Build index name
        index_prefix = settings.opensearch_index_prefix
        index_name = f"{index_prefix}_{index_type}"

        # Build full URL for indexing
        index_url = f"{opensearch_url}/{index_name}/_doc"

        # Send POST request to index document
        response = requests.post(
            index_url,
            json=document,
            headers={"Content-Type": "application/json"},
            timeout=10,
        )

        # Check response
        response.raise_for_status()
        result = response.json()
        print(f"OpenSearch indexed: {result}")
        logger.info(f"OpenSearch indexed: {result}")

    @staticmethod
    async def handle_click_event(data: dict[str, Any]) -> None:
        """
        Handle click event message from RabbitMQ.

        Args:
            data: Message data containing click event information
        """
        try:
            # Log received message to both logger and console
            separator = "=" * 80
            print(separator)
            print("📨 Received message from RabbitMQ queue")
            print(separator)
            print(f"Payload:\n{json.dumps(data, indent=2)}")
            print(separator)

            logger.info(separator)
            logger.info("📨 Received message from RabbitMQ queue")
            logger.info(separator)
            logger.info(f"Payload:\n{json.dumps(data, indent=2)}")
            logger.info(separator)

            index_type = data.get("index_type")

            if not index_type:
                logger.warning("Missing index_type in message, skipping")
                return

            if index_type != "click_events":
                logger.warning(f"Unknown index_type: {index_type}, skipping message")
                return

            # Get the data payload - this is already formatted for OpenSearch
            event_data = data.get("data", {})

            if not event_data:
                logger.warning("Missing data in message, skipping")
                return

            # Extract short code for logging
            short_code = event_data.get("short_code", "unknown")
            print(f"📊 Processing click event for short_code: {short_code}")
            logger.info(f"📊 Processing click event for short_code: {short_code}")

            # Index directly to OpenSearch - data is already in correct format
            try:
                MessageHandler.index_to_opensearch_sync(event_data, index_type)
                print(f"✅ Click event indexed to OpenSearch for short_code: {short_code}")
                logger.info(f"✅ Click event indexed to OpenSearch for short_code: {short_code}")
            except Exception as e:
                print(f"❌ Failed to index event to OpenSearch: {e}")
                logger.error(f"❌ Failed to index event to OpenSearch: {e}")
                raise

        except Exception as e:
            print(f"❌ Error handling click event: {e}")
            logger.error(f"❌ Error handling click event: {e}", exc_info=True)
            raise


async def create_message_handler(data: dict[str, Any]) -> None:
    """
    Create and execute message handler.

    This function is passed to the RabbitMQ consumer.

    Args:
        data: Message data from RabbitMQ
    """
    await MessageHandler.handle_click_event(data)
