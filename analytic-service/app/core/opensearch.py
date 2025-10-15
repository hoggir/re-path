"""OpenSearch connection manager."""

from typing import Optional
from urllib.parse import urlparse

from opensearchpy import AsyncOpenSearch

from app.core.config import settings


class OpenSearchManager:
    """Manages OpenSearch connections."""

    client: Optional[AsyncOpenSearch] = None

    @classmethod
    async def connect_to_opensearch(cls) -> None:
        """
        Establish connection to OpenSearch.

        Raises:
            Exception: If connection fails.
        """
        try:
            # Parse URL to extract username and password
            parsed_url = urlparse(settings.opensearch_url)

            # Create client configuration
            host = parsed_url.hostname
            port = parsed_url.port or 443
            username = parsed_url.username
            password = parsed_url.password

            cls.client = AsyncOpenSearch(
                hosts=[{"host": host, "port": port}],
                http_auth=(username, password) if username and password else None,
                use_ssl=True,
                verify_certs=True,
                ssl_show_warn=False,
                timeout=30,
                max_retries=3,
                retry_on_timeout=True,
            )

            # Test connection
            info = await cls.client.info()

            print("✅ Connected to OpenSearch successfully")
            print(f"   Host: {host}")
            print(f"   Version: {info['version']['number']}")
            print(f"   Cluster: {info['cluster_name']}")

        except Exception as e:
            print(f"❌ Failed to connect to OpenSearch: {e}")
            if cls.client:
                await cls.client.close()
                cls.client = None
            raise

    @classmethod
    async def close_opensearch_connection(cls) -> None:
        """Close OpenSearch connection."""
        if cls.client:
            await cls.client.close()
            print("👋 OpenSearch connection closed")

    @classmethod
    def get_client(cls) -> AsyncOpenSearch:
        """
        Get OpenSearch client instance.

        Returns:
            AsyncOpenSearch: The OpenSearch client.

        Raises:
            RuntimeError: If client is not initialized.
        """
        if cls.client is None:
            raise RuntimeError(
                "OpenSearch client is not initialized. Call connect_to_opensearch() first."
            )
        return cls.client

    @classmethod
    def get_index_name(cls, index_type: str) -> str:
        """
        Get full index name with prefix.

        Args:
            index_type: Type of index (e.g., 'analytics', 'events').

        Returns:
            str: Full index name with prefix.
        """
        return f"{settings.opensearch_index_prefix}_{index_type}"
