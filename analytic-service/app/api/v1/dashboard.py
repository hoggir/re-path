"""Health check endpoints."""

import logging

from fastapi import APIRouter, status

from app.schemas.dashboard import DashboardData
from app.schemas.response import ApiResponse, create_response
from app.services.opensearch_service import OpenSearchService

router = APIRouter(tags=["Dashboard"])
logger = logging.getLogger(__name__)


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
async def first_dashboard(
    # request: SearchRequest,
) -> ApiResponse[DashboardData]:
    response = await OpenSearchService.search_all(
        index_type="click_events",
        query={"match": {"short_code": "my-custom-link"}},
        from_=0,
    )

    # first_short_code = response["hits"]
    hits = response.get("hits", {})
    # total_data = hits.get{"value", 0}
    total_data = hits.get("total", {}).get("value", 0)

    # short_codes = [hit["_source"]]
    print(hits)
    print(total_data)

    # logger.info("Data", response)
    # logger.info(f"✅ Data : {response}, ")
    dashboard_data = DashboardData(
        status="healthy",
    )

    return create_response(
        data=dashboard_data,
        message="Data retrieved successfully",
    )
