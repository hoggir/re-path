"""OpenSearch service for analytics operations."""

from datetime import datetime, timezone
from typing import Any, Optional

from app.core.opensearch import OpenSearchManager


class OpenSearchService:
    """Service for OpenSearch operations."""

    @staticmethod
    async def index_document(
        index_type: str,
        document: dict[str, Any],
        doc_id: Optional[str] = None,
    ) -> dict[str, Any]:
        """
        Index a document in OpenSearch.

        Args:
            index_type: Type of index (e.g., 'analytics', 'events').
            document: Document to index.
            doc_id: Optional document ID.

        Returns:
            Dict containing the index response.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        # Add timestamp if not present
        if "timestamp" not in document:
            document["timestamp"] = datetime.now(timezone.utc).isoformat()

        response = await client.index(
            index=index_name,
            id=doc_id,
            body=document,
        )

        return response

    @staticmethod
    async def search(
        index_type: str,
        query: dict[str, Any],
        size: int = 10,
        from_: int = 0,
    ) -> dict[str, Any]:
        """
        Search documents in OpenSearch.

        Args:
            index_type: Type of index to search.
            query: OpenSearch query DSL.
            size: Number of results to return.
            from_: Offset for pagination.

        Returns:
            Dict containing search results.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        response = await client.search(
            index=index_name,
            body={"query": query},
            size=size,
            from_=from_,
        )

        return response

    @staticmethod
    async def get_document(index_type: str, doc_id: str) -> Optional[dict[str, Any]]:
        """
        Get a document by ID.

        Args:
            index_type: Type of index.
            doc_id: Document ID.

        Returns:
            Document if found, None otherwise.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        try:
            response = await client.get(index=index_name, id=doc_id)
            return response
        except Exception:
            return None

    @staticmethod
    async def delete_document(index_type: str, doc_id: str) -> dict[str, Any]:
        """
        Delete a document by ID.

        Args:
            index_type: Type of index.
            doc_id: Document ID.

        Returns:
            Dict containing the delete response.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        response = await client.delete(index=index_name, id=doc_id)
        return response

    @staticmethod
    async def bulk_index(
        index_type: str,
        documents: list[dict[str, Any]],
    ) -> dict[str, Any]:
        """
        Bulk index multiple documents.

        Args:
            index_type: Type of index.
            documents: List of documents to index.

        Returns:
            Dict containing the bulk response.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        body = []
        for doc in documents:
            # Add timestamp if not present
            if "timestamp" not in doc:
                doc["timestamp"] = datetime.now(timezone.utc).isoformat()

            body.append({"index": {"_index": index_name}})
            body.append(doc)

        response = await client.bulk(body=body)
        return response

    @staticmethod
    async def create_index(
        index_type: str,
        mappings: Optional[dict[str, Any]] = None,
        settings: Optional[dict[str, Any]] = None,
    ) -> dict[str, Any]:
        """
        Create an index with optional mappings and settings.

        Args:
            index_type: Type of index to create.
            mappings: Optional index mappings.
            settings: Optional index settings.

        Returns:
            Dict containing the create response.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        body = {}
        if mappings:
            body["mappings"] = mappings
        if settings:
            body["settings"] = settings

        response = await client.indices.create(index=index_name, body=body)
        return response

    @staticmethod
    async def index_exists(index_type: str) -> bool:
        """
        Check if an index exists.

        Args:
            index_type: Type of index to check.

        Returns:
            True if index exists, False otherwise.
        """
        client = OpenSearchManager.get_client()
        index_name = OpenSearchManager.get_index_name(index_type)

        return await client.indices.exists(index=index_name)