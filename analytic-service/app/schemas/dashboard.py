"""
Dashboard RPC Schemas - Dashboard service interface between services.

These schemas define the dashboard contract between redirect-service (Go) and analytic-service (Python).
Any changes to these schemas MUST be coordinated between both services.
"""

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


class TopLinkData(BaseModel):
    short_url: str = Field(..., description="Short code of the link")
    original_url: str = Field(..., description="Original URL of the link")
    clicks: int = Field(..., description="Number of clicks", ge=0)
    status: bool = Field(..., description="Link active status")

    class Config:
        json_schema_extra = {
            "example": {
                "short_url": "popular-link",
                "original_url": "https://example.com/popular",
                "clicks": 350,
                "status": True,
            }
        }


class StatsLinkData(BaseModel):
    date: str = Field(..., description="Date of the link stats")
    clicks: int = Field(..., description="Number of clicks on that date", ge=0)

    class Config:
        json_schema_extra = {
            "example": {
                "date": "21 Oct",
                "clicks": 150,
            }
        }


class DashboardResponse(BaseModel):
    user_id: int = Field(..., description="User ID")
    total_clicks: int = Field(..., description="Total number of clicks", ge=0)
    total_links: int = Field(..., description="Total number of unique links", ge=0)
    uniq_visitors: int = Field(..., description="Total number of unique visitors", ge=0)
    top_links: list[TopLinkData] = Field(
        default_factory=list, description="Top 5 links by click count"
    )
    stat_links: list[StatsLinkData] = Field(default_factory=list, description="Statistic links")
    status: Literal["success", "error", "limited"] = Field(
        default="success", description="Response status"
    )
    message: Optional[str] = Field(None, description="Optional message")

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
                "uniq_visitors": 2,
                "top_links": [
                    {
                        "short_url": "popular-link",
                        "original_url": "https://example.com/popular",
                        "clicks": 350,
                        "status": True,
                    }
                ],
                "stat_links": [{"date": "", "clicks": 123}],
                "status": "success",
            }
        }


class DashboardErrorResponse(BaseModel):
    user_id: Optional[int] = Field(None, description="User ID if available")
    status: Literal["error"] = Field(default="error", description="Error status")
    message: str = Field(..., description="Error message")
    total_clicks: int = Field(default=0, description="Default to 0 on error")
    total_links: int = Field(default=0, description="Default to 0 on error")
    top_links: list[TopLinkData] = Field(default_factory=list, description="Empty on error")

    class Config:
        json_schema_extra = {
            "example": {
                "user_id": 1,
                "status": "error",
                "message": "Database connection failed",
                "total_clicks": 0,
                "total_links": 0,
                "top_links": [],
                "stat_links": [],
            }
        }


def validate_dashboard_request(data: dict[str, Any]) -> DashboardRequest:
    return DashboardRequest.model_validate(data)


def create_dashboard_response(
    user_id: int,
    total_clicks: int,
    total_links: int,
    uniq_visitors: int,
    top_links: list[TopLinkData],
    stat_links: list[StatsLinkData],
    status: Literal["success", "error", "limited"] = "success",
    message: Optional[str] = None,
) -> dict[DashboardResponse]:
    response = DashboardResponse(
        user_id=user_id,
        total_clicks=total_clicks,
        total_links=total_links,
        uniq_visitors=uniq_visitors,
        stat_links=[StatsLinkData(**link) for link in stat_links],
        top_links=[TopLinkData(**link) for link in top_links],
        status=status,
        message=message,
    )
    return response.model_dump(mode="json", exclude_none=True)


def create_dashboard_error_response(message: str, user_id: Optional[int] = None) -> dict[str, Any]:
    response = DashboardErrorResponse(user_id=user_id, message=message)
    return response.model_dump(mode="json", exclude_none=True)
