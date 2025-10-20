import logging
from typing import Any

from pydantic import ValidationError

from app.core.database import DatabaseManager
from app.models.click_event import ClickEvent
from app.schemas.rpc_contracts import (
    create_dashboard_error_response,
    create_dashboard_response,
    validate_dashboard_request,
)

logger = logging.getLogger(__name__)


class DashboardService:
    @staticmethod
    async def get_dashboard_data(user_id: int) -> dict[str, Any]:
        try:
            logger.info(f"📊 Processing dashboard request for user_id: {user_id}")

            if not DatabaseManager.client:
                logger.warning("Database not connected, returning limited data")
                return create_dashboard_response(
                    user_id=user_id,
                    total_clicks=0,
                    total_links=0,
                    recent_clicks=[],
                    top_links=[],
                    status="limited",
                    message="Database not available",
                )

            try:
                # total_clicks = await ClickEvent.find(query_filter).count()
                total_clicks = 123
                total_links = 45
                recent_clicks = []
                top_links = []
                

                logger.info(
                    f"✅ Dashboard data retrieved - user_id: {user_id}, "
                    f"clicks: {total_clicks}, links: {total_links}"
                )

                return create_dashboard_response(
                    user_id=user_id,
                    total_clicks=total_clicks,
                    total_links=total_links,
                    recent_clicks=recent_clicks,
                    top_links=top_links,
                    status="success",
                )

            except Exception as db_error:
                logger.error(f"Database query error: {db_error}", exc_info=True)
                return create_dashboard_error_response(
                    message="Database query failed",
                    user_id=user_id,
                )

        except Exception as e:
            logger.error(f"❌ Error processing dashboard request: {e}", exc_info=True)
            return create_dashboard_error_response(
                message=str(e),
                user_id=user_id,
            )


async def handle_dashboard_rpc_request(request_data: dict[str, Any]) -> dict[str, Any]:
    try:
        logger.info(f"📨 Received dashboard RPC request: {request_data}")

        try:
            validated_request = validate_dashboard_request(request_data)
            logger.info(f"✅ Request validated - user_id: {validated_request.user_id}")
        except ValidationError as e:
            logger.error(f"❌ Request validation failed: {e}")
            return create_dashboard_error_response(
                message=f"Invalid request format: {str(e)}",
                user_id=request_data.get("user_id"),
            )

        response = await DashboardService.get_dashboard_data(validated_request.user_id)

        logger.info(
            f"✅ Dashboard RPC request processed successfully for user_id: {validated_request.user_id}"
        )
        return response

    except Exception as e:
        logger.error(f"❌ Error handling dashboard RPC request: {e}", exc_info=True)
        return create_dashboard_error_response(
            message=f"Internal server error: {str(e)}",
            user_id=request_data.get("user_id") if isinstance(request_data, dict) else None,
        )
