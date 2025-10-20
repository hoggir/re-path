"""
RPC Contract Schemas - Shared interface between services.

These schemas define the contract between redirect-service (Go) and analytic-service (Python).
Any changes to these schemas MUST be coordinated between both services.
"""

from datetime import datetime
from typing import Any, Literal, Optional

from pydantic import BaseModel, Field, field_validator

class DashboardRequest(BaseModel):
    user_id: int = Field(
        ...,
        description="User ID to get dashboard data for",
        examples=[1, 123, 456],
        gt=0,  # Greater than 0
    )

    class Config:
        json_schema_extra = {
            "example": {"user_id": 1},
            "description": "Request to get dashboard data for a specific user",
        }


class RecentClickData(BaseModel):
    short_code: str = Field(..., description="Short code of the clicked link")
    clicked_at: str = Field(..., description="ISO 8601 timestamp of the click")
    ip_address_hash: Optional[str] = Field(None, description="Hashed IP address")
    user_agent: Optional[str] = Field(None, description="User agent string")
    country_code: Optional[str] = Field(None, description="ISO country code")
    city: Optional[str] = Field(None, description="City name")
    device_type: Optional[str] = Field(None, description="Device type (desktop, mobile, tablet)")
    browser_name: Optional[str] = Field(None, description="Browser name")
    is_bot: bool = Field(default=False, description="Whether the click is from a bot")

    class Config:
        json_schema_extra = {
            "example": {
                "short_code": "my-link",
                "clicked_at": "2025-10-20T10:00:00Z",
                "ip_address_hash": "abc123hash",
                "user_agent": "Mozilla/5.0...",
                "country_code": "ID",
                "city": "Jakarta",
                "device_type": "desktop",
                "browser_name": "Chrome",
                "is_bot": False,
            }
        }


class TopLinkData(BaseModel):
    short_code: str = Field(..., description="Short code of the link")
    clicks: int = Field(..., description="Number of clicks", ge=0)

    class Config:
        json_schema_extra = {"example": {"short_code": "popular-link", "clicks": 350}}


class DashboardResponse(BaseModel):
    user_id: int = Field(..., description="User ID")
    total_clicks: int = Field(..., description="Total number of clicks", ge=0)
    total_links: int = Field(..., description="Total number of unique links", ge=0)
    recent_clicks: list[RecentClickData] = Field(
        default_factory=list, description="List of recent clicks (max 10)"
    )
    top_links: list[TopLinkData] = Field(
        default_factory=list, description="Top 5 links by click count"
    )
    status: Literal["success", "error", "limited"] = Field(
        ..., description="Response status"
    )
    message: Optional[str] = Field(None, description="Optional message or error description")

    @field_validator("recent_clicks")
    @classmethod
    def validate_recent_clicks_limit(cls, v: list[RecentClickData]) -> list[RecentClickData]:
        if len(v) > 10:
            return v[:10]
        return v

    @field_validator("top_links")
    @classmethod
    def validate_top_links_limit(cls, v: list[TopLinkData]) -> list[TopLinkData]:
        if len(v) > 5:
            return v[:5]
        return v

    class Config:
        json_schema_extra = {
            "example": {
                "user_id": 1,
                "total_clicks": 1542,
                "total_links": 45,
                "recent_clicks": [
                    {
                        "short_code": "my-link",
                        "clicked_at": "2025-10-20T10:00:00Z",
                        "ip_address_hash": "abc123",
                        "user_agent": "Mozilla/5.0...",
                        "country_code": "ID",
                        "city": "Jakarta",
                        "device_type": "desktop",
                        "browser_name": "Chrome",
                        "is_bot": False,
                    }
                ],
                "top_links": [{"short_code": "popular-link", "clicks": 350}],
                "status": "success",
            }
        }


class DashboardErrorResponse(BaseModel):
    user_id: Optional[int] = Field(None, description="User ID if available")
    status: Literal["error"] = Field(default="error", description="Error status")
    message: str = Field(..., description="Error message")
    total_clicks: int = Field(default=0, description="Default to 0 on error")
    total_links: int = Field(default=0, description="Default to 0 on error")
    recent_clicks: list[Any] = Field(default_factory=list, description="Empty on error")
    top_links: list[Any] = Field(default_factory=list, description="Empty on error")

    class Config:
        json_schema_extra = {
            "example": {
                "user_id": 1,
                "status": "error",
                "message": "Database connection failed",
                "total_clicks": 0,
                "total_links": 0,
                "recent_clicks": [],
                "top_links": [],
            }
        }

class ClickEventMessage(BaseModel):
    index_type: Literal["click_events"] = Field(
        ..., description="Message type identifier"
    )
    data: dict[str, Any] = Field(..., description="Click event data")

    class Config:
        json_schema_extra = {
            "example": {
                "index_type": "click_events",
                "data": {
                    "short_code": "my-link",
                    "clicked_at": "2025-10-20T10:00:00Z",
                    "ip_address_hash": "abc123hash",
                    "user_agent": "Mozilla/5.0...",
                    "is_bot": False,
                },
            }
        }

class SchemaVersion(BaseModel):
    version: str = Field(default="1.0.0", description="Schema version (semver)")
    compatible_versions: list[str] = Field(
        default_factory=lambda: ["1.0.0"],
        description="List of compatible versions",
    )
    last_updated: datetime = Field(
        default_factory=datetime.utcnow,
        description="Last update timestamp",
    )


CURRENT_SCHEMA_VERSION = SchemaVersion(
    version="1.0.0",
    compatible_versions=["1.0.0"],
)

def validate_dashboard_request(data: dict[str, Any]) -> DashboardRequest:
    return DashboardRequest.model_validate(data)


def create_dashboard_response(
    user_id: int,
    total_clicks: int,
    total_links: int,
    recent_clicks: list[dict[str, Any]],
    top_links: list[dict[str, Any]],
    status: Literal["success", "error", "limited"] = "success",
    message: Optional[str] = None,
) -> dict[str, Any]:
    response = DashboardResponse(
        user_id=user_id,
        total_clicks=total_clicks,
        total_links=total_links,
        recent_clicks=[RecentClickData(**click) for click in recent_clicks],
        top_links=[TopLinkData(**link) for link in top_links],
        status=status,
        message=message,
    )
    return response.model_dump(mode="json", exclude_none=True)


def create_dashboard_error_response(
    message: str, user_id: Optional[int] = None
) -> dict[str, Any]:
    response = DashboardErrorResponse(user_id=user_id, message=message)
    return response.model_dump(mode="json", exclude_none=True)
