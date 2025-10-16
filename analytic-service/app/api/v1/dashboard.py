"""Health check endpoints."""

from datetime import datetime

from fastapi import APIRouter, status

from app.core.config import settings
from app.schemas.dashboard import DashboardData
from app.schemas.response import ApiResponse, create_response

router = APIRouter(tags=["Dashboard"])


@router.get(
    "/dashboard",
    response_model=ApiResponse[DashboardData],
    status_code=status.HTTP_200_OK,
    summary="Dashboard",
    description="Dashboard endpoint",
    responses={
        200: {
            "description": "Data retrieved successfully",
            "content": {
                "application/json": {
                    "example": {
                        "success": True,
                        "message": "Data retrieved successfully",
                        "data": {
                            "status": "healthy",
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
async def first_dashboard() -> ApiResponse[DashboardData]:
    dashboard_data = DashboardData(
        status="healthy",
    )

    return create_response(
        data=dashboard_data,
        message="Data retrieved successfully",
    )
