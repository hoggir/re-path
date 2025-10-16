"""Dashboard check schemas."""

from typing import Literal

from pydantic import BaseModel, Field


class DashboardData(BaseModel):
    """Dashboard check data schema (used inside ApiResponse)."""

    status: Literal["healthy", "unhealthy"] = Field(
        description="Service health status",
        examples=["healthy"],
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "status": "healthy",
            }
        }
    }


# Backward compatibility - keep old name as alias
DashboardResponse = DashboardData
