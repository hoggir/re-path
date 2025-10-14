"""Health check schemas."""

from datetime import datetime
from typing import Literal, Optional

from pydantic import BaseModel, Field


class DatabaseHealth(BaseModel):
    """Database health information."""

    connected: bool = Field(
        description="Database connection status",
    )
    database: Optional[str] = Field(
        default=None,
        description="Database name",
    )
    ping: bool = Field(
        default=False,
        description="Database ping successful",
    )


class HealthData(BaseModel):
    """Health check data schema (used inside ApiResponse)."""

    status: Literal["healthy", "unhealthy"] = Field(
        description="Service health status",
        examples=["healthy"],
    )
    service: str = Field(
        description="Service name",
        examples=["analytic-service"],
    )
    version: str = Field(
        description="Service version",
        examples=["0.1.0"],
    )
    timestamp: datetime = Field(
        description="Current server timestamp",
        examples=["2025-10-14T10:30:00Z"],
    )
    environment: str = Field(
        description="Current environment",
        examples=["development", "production"],
    )
    database: Optional[DatabaseHealth] = Field(
        default=None,
        description="Database health information",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "status": "healthy",
                "service": "analytic-service",
                "version": "0.1.0",
                "timestamp": "2025-10-14T10:30:00Z",
                "environment": "development",
                "database": {
                    "connected": True,
                    "database": "analytic_db",
                    "ping": True,
                },
            }
        }
    }


# Backward compatibility - keep old name as alias
HealthResponse = HealthData
