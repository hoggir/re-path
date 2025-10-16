"""Dashboard check schemas."""

from datetime import datetime
from typing import Literal, Optional

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
