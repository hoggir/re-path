"""Health check endpoints."""

from datetime import datetime

from fastapi import APIRouter, status

from app.core.config import settings
from app.core.database import get_database_health
from app.schemas.health import HealthData
from app.schemas.response import ApiResponse, create_response

router = APIRouter(tags=["Health"])


@router.get(
    "/health",
    response_model=ApiResponse[HealthData],
    status_code=status.HTTP_200_OK,
    summary="Health Check",
    description="Check service health status and get service information",
    responses={
        200: {
            "description": "Service is healthy",
            "content": {
                "application/json": {
                    "example": {
                        "success": True,
                        "message": "Service is healthy",
                        "data": {
                            "status": "healthy",
                            "service": "analytic-service",
                            "version": "0.1.0",
                            "timestamp": "2025-10-14T10:30:00Z",
                            "environment": "development",
                        },
                        "meta": {
                            "timestamp": "2025-10-14T10:30:00Z",
                            "version": "v1",
                            "request_id": None,
                        },
                    }
                }
            },
        }
    },
)
async def health_check() -> ApiResponse[HealthData]:
    """
    Health check endpoint.

    Returns service status, version, and environment information.
    This endpoint can be used for:
    - Kubernetes/Docker health probes
    - Load balancer health checks
    - Service monitoring and alerting

    Returns a standardized API response with health data including database status.
    """
    # Get database health
    db_health_check = await get_database_health()

    # Convert to DatabaseHealth schema
    from app.schemas.health import DatabaseHealth
    db_health = DatabaseHealth(
        connected=db_health_check.connected,
        database=db_health_check.database,
        ping=db_health_check.ping,
    )

    health_data = HealthData(
        status="healthy",
        service=settings.app_name,
        version=settings.app_version,
        timestamp=datetime.utcnow(),
        environment=settings.app_env,
        database=db_health,
    )

    return create_response(
        data=health_data,
        message="Service is healthy",
    )
