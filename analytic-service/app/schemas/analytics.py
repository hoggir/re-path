"""Analytics schemas for request and response models."""

from datetime import datetime
from typing import Any, Optional

from pydantic import BaseModel, Field


class AnalyticsDataRequest(BaseModel):
    """Request schema for analytics data submission."""

    index_type: str = Field(
        default="analytics",
        description="Type of index to store data (e.g., 'analytics', 'events', 'logs')",
        examples=["analytics", "events", "logs"],
    )
    data: dict[str, Any] = Field(
        description="JSON data to be indexed in OpenSearch",
        examples=[
            {
                "user_id": "user123",
                "action": "page_view",
                "page": "/products",
                "metadata": {"category": "electronics"},
            }
        ],
    )
    doc_id: Optional[str] = Field(
        default=None,
        description="Optional document ID. If not provided, OpenSearch will generate one",
        examples=["doc_123"],
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "index_type": "analytics",
                "data": {
                    "user_id": "user123",
                    "event_type": "page_view",
                    "page_url": "https://example.com/products",
                    "timestamp": "2025-10-15T10:30:00Z",
                    "metadata": {
                        "category": "electronics",
                        "product_id": "prod_456",
                    },
                },
                "doc_id": None,
            }
        }
    }


class BulkAnalyticsDataRequest(BaseModel):
    """Request schema for bulk analytics data submission."""

    index_type: str = Field(
        default="analytics",
        description="Type of index to store data",
        examples=["analytics", "events"],
    )
    documents: list[dict[str, Any]] = Field(
        description="List of JSON documents to be indexed",
        min_length=1,
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "index_type": "analytics",
                "documents": [
                    {
                        "user_id": "user123",
                        "event_type": "page_view",
                        "page_url": "https://example.com/products",
                    },
                    {
                        "user_id": "user456",
                        "event_type": "click",
                        "page_url": "https://example.com/cart",
                    },
                ],
            }
        }
    }


class AnalyticsDataResponse(BaseModel):
    """Response schema for analytics data submission."""

    id: str = Field(
        description="Document ID in OpenSearch",
        examples=["doc_123"],
    )
    index: str = Field(
        description="Index name where document was stored",
        examples=["repath_analytics"],
    )
    result: str = Field(
        description="Result of the operation",
        examples=["created", "updated"],
    )
    timestamp: datetime = Field(
        default_factory=datetime.utcnow,
        description="Timestamp when data was indexed",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "id": "doc_123",
                "index": "repath_analytics",
                "result": "created",
                "timestamp": "2025-10-15T10:30:00Z",
            }
        }
    }


class BulkAnalyticsDataResponse(BaseModel):
    """Response schema for bulk analytics data submission."""

    items_count: int = Field(
        description="Number of documents processed",
        examples=[10],
    )
    success_count: int = Field(
        description="Number of successfully indexed documents",
        examples=[9],
    )
    error_count: int = Field(
        description="Number of failed documents",
        examples=[1],
    )
    errors: Optional[list[dict[str, Any]]] = Field(
        default=None,
        description="List of errors if any occurred",
    )
    timestamp: datetime = Field(
        default_factory=datetime.utcnow,
        description="Timestamp when bulk operation completed",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "items_count": 10,
                "success_count": 10,
                "error_count": 0,
                "errors": None,
                "timestamp": "2025-10-15T10:30:00Z",
            }
        }
    }


class SearchRequest(BaseModel):
    """Request schema for searching analytics data."""

    index_type: str = Field(
        default="analytics",
        description="Type of index to search",
        examples=["analytics", "events"],
    )
    query: dict[str, Any] = Field(
        description="OpenSearch query DSL",
        examples=[{"match": {"event_type": "page_view"}}],
    )
    size: int = Field(
        default=10,
        ge=1,
        le=100,
        description="Number of results to return",
    )
    from_: int = Field(
        default=0,
        ge=0,
        description="Offset for pagination",
        alias="from",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "index_type": "analytics",
                "query": {
                    "match": {"event_type": "page_view"},
                },
                "size": 10,
                "from": 0,
            }
        }
    }
