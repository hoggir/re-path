"""Analytics event model for MongoDB using Beanie ODM."""

from datetime import datetime
from typing import Dict, Optional

from beanie import Document, Indexed
from pydantic import Field


class AnalyticsEvent(Document):
    """
    Analytics event document model.

    This model represents an analytics event stored in MongoDB.
    Uses Beanie ODM for async MongoDB operations.
    """

    event_type: Indexed(str) = Field(  # type: ignore
        description="Type of analytics event (e.g., page_view, click, conversion)",
        examples=["page_view", "click", "conversion"],
    )

    user_id: Optional[str] = Field(
        default=None,
        description="User identifier (if authenticated)",
        examples=["user_123"],
    )

    session_id: Optional[str] = Field(
        default=None,
        description="Session identifier",
        examples=["session_abc123"],
    )

    page_url: Optional[str] = Field(
        default=None,
        description="URL of the page where event occurred",
        examples=["https://example.com/products"],
    )

    referrer: Optional[str] = Field(
        default=None,
        description="HTTP referrer URL",
        examples=["https://google.com"],
    )

    user_agent: Optional[str] = Field(
        default=None,
        description="User agent string",
    )

    ip_address: Optional[str] = Field(
        default=None,
        description="Client IP address",
        examples=["192.168.1.1"],
    )

    metadata: Dict[str, str] = Field(
        default_factory=dict,
        description="Additional event metadata",
        examples=[{"product_id": "123", "category": "electronics"}],
    )

    created_at: Indexed(datetime) = Field(  # type: ignore
        default_factory=datetime.utcnow,
        description="Event creation timestamp",
    )

    class Settings:
        """Beanie document settings."""

        name = "analytics_events"  # MongoDB collection name
        use_state_management = True
        use_revision = True

        # Indexes for better query performance
        indexes = [
            [("event_type", 1), ("created_at", -1)],  # Compound index
            "user_id",
            "session_id",
        ]

    class Config:
        """Pydantic model configuration."""

        json_schema_extra = {
            "example": {
                "event_type": "page_view",
                "user_id": "user_123",
                "session_id": "session_abc123",
                "page_url": "https://example.com/products",
                "referrer": "https://google.com",
                "user_agent": "Mozilla/5.0...",
                "ip_address": "192.168.1.1",
                "metadata": {"product_id": "123", "category": "electronics"},
                "created_at": "2025-10-14T10:30:00Z",
            }
        }
