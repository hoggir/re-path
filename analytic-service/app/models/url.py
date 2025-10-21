"""URL model for MongoDB using Beanie ODM."""

from datetime import datetime
from typing import Optional

from beanie import Document
from pydantic import BaseModel, Field
from pymongo import ASCENDING, DESCENDING, IndexModel


class URLMetadata(BaseModel):
    """URL metadata subdocument."""

    title: Optional[str] = None
    description: Optional[str] = None
    tags: Optional[list[str]] = None


class URL(Document):
    """
    URL document model.

    This model represents a shortened URL link.
    Matches the Go URL struct for compatibility.
    """

    short_code: str = Field(
        alias="shortCode",
        description="Short code of the link",
    )

    original_url: str = Field(
        alias="originalUrl",
        description="Original long URL",
    )

    custom_alias: Optional[str] = Field(
        default=None,
        alias="customAlias",
        description="Custom alias for the short code",
    )

    user_id: int = Field(
        alias="userId",
        description="User ID who owns this link",
    )

    click_count: int = Field(
        default=0,
        alias="clickCount",
        description="Total number of clicks",
    )

    is_active: bool = Field(
        default=True,
        alias="isActive",
        description="Whether the link is active",
    )

    expires_at: Optional[datetime] = Field(
        default=None,
        alias="expiresAt",
        description="Expiration timestamp",
    )

    metadata: Optional[URLMetadata] = Field(
        default=None,
        description="Additional metadata",
    )

    created_at: datetime = Field(
        alias="createdAt",
        description="Creation timestamp",
    )

    updated_at: datetime = Field(
        alias="updatedAt",
        description="Last update timestamp",
    )

    class Settings:
        """Beanie document settings."""

        name = "urls"
        use_state_management = True

        # Indexes for better query performance
        indexes = [
            IndexModel([("userId", ASCENDING)]),
            IndexModel([("isActive", ASCENDING)]),
            IndexModel([("createdAt", DESCENDING)]),
        ]

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True
        json_schema_extra = {
            "example": {
                "shortCode": "abc123",
                "originalUrl": "https://example.com/very/long/url",
                "customAlias": "my-link",
                "userId": 1,
                "clickCount": 42,
                "isActive": True,
                "expiresAt": "2025-12-31T23:59:59Z",
                "metadata": {
                    "title": "Example Link",
                    "description": "This is an example",
                    "tags": ["example", "test"],
                },
                "createdAt": "2025-10-16T10:00:00Z",
                "updatedAt": "2025-10-16T10:00:00Z",
            }
        }
