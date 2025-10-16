"""Click event model for MongoDB using Beanie ODM."""

from datetime import datetime
from typing import Optional

from beanie import Document
from pydantic import Field
from pymongo import ASCENDING, DESCENDING, IndexModel


class ClickEvent(Document):
    """
    Click event document model.

    This model represents a click event for link tracking.
    Matches the Go ClickEvent struct for compatibility.
    """

    short_code: str = Field(
        alias="shortCode",
        description="Short code of the link",
    )

    clicked_at: datetime = Field(
        alias="clickedAt",
        description="Timestamp when clicked",
    )

    ip_address_hash: str = Field(
        alias="ipAddressHash",
        description="Hashed IP address",
    )

    user_agent: str = Field(
        alias="userAgent",
        description="User agent string",
    )

    referrer_url: Optional[str] = Field(
        default=None,
        alias="referrerUrl",
        description="Referrer URL",
    )

    referrer_domain: Optional[str] = Field(
        default=None,
        alias="referrerDomain",
        description="Referrer domain",
    )

    country_code: Optional[str] = Field(
        default=None,
        alias="countryCode",
        description="Country code",
    )

    city: Optional[str] = Field(
        default=None,
        alias="city",
        description="City name",
    )

    region: Optional[str] = Field(
        default=None,
        alias="region",
        description="Region/state name",
    )

    device_type: Optional[str] = Field(
        default=None,
        alias="deviceType",
        description="Device type (desktop, mobile, tablet)",
    )

    browser_name: Optional[str] = Field(
        default=None,
        alias="browserName",
        description="Browser name",
    )

    browser_version: Optional[str] = Field(
        default=None,
        alias="browserVersion",
        description="Browser version",
    )

    os_name: Optional[str] = Field(
        default=None,
        alias="osName",
        description="Operating system name",
    )

    os_version: Optional[str] = Field(
        default=None,
        alias="osVersion",
        description="Operating system version",
    )

    lat: Optional[float] = Field(
        default=None,
        alias="lat",
        description="Latitude",
    )

    lon: Optional[float] = Field(
        default=None,
        alias="lon",
        description="Longitude",
    )

    is_bot: bool = Field(
        alias="isBot",
        description="Whether the click is from a bot",
    )

    class Settings:
        """Beanie document settings."""

        name = "click_events"
        use_state_management = True

        # Indexes for better query performance
        indexes = [
            IndexModel([("shortCode", ASCENDING), ("clickedAt", DESCENDING)]),
            IndexModel([("clickedAt", DESCENDING)]),
            IndexModel([("countryCode", ASCENDING)]),
            IndexModel([("isBot", ASCENDING)]),
        ]

    class Config:
        """Pydantic model configuration."""

        populate_by_name = True  # Allow both alias and field name
        json_schema_extra = {
            "example": {
                "shortCode": "my-link",
                "clickedAt": "2025-10-16T10:00:00Z",
                "ipAddressHash": "abc123hash",
                "userAgent": "Mozilla/5.0...",
                "referrerUrl": "https://google.com",
                "referrerDomain": "google.com",
                "countryCode": "ID",
                "city": "Jakarta",
                "region": "Jakarta",
                "deviceType": "desktop",
                "browserName": "Chrome",
                "browserVersion": "120.0",
                "osName": "Linux",
                "osVersion": "x86_64",
                "lat": -6.2088,
                "lon": 106.8456,
                "isBot": False,
            }
        }
