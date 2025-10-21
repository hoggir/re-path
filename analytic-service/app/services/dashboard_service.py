import logging
from collections import Counter
from typing import Any

from pydantic import ValidationError

from app.core.database import DatabaseManager
from app.models.click_event import ClickEvent
from app.models.url import URL
from app.schemas.dashboard import (
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
                    top_links=[],
                    status="limited",
                    message="Database not available",
                )

            try:
                user_links = await URL.find({"userId": user_id}).to_list()
                total_links = len(user_links)

                short_codes = [link.short_code for link in user_links]
                if short_codes:
                    click_events = await ClickEvent.find(
                        {"shortCode": {"$in": short_codes}}
                    ).to_list()

                    total_clicks = len(click_events)
                    uniq_visitors = len({click.ip_address_hash for click in click_events})

                    counter = Counter(
                        event.clicked_at.strftime("%-d %b")
                        for event in click_events
                    )

                    data_stat = [
                        {"date": date, "clicks": count} for date, count in sorted(counter.items())
                    ]

                    sorted_links = sorted(user_links, key=lambda x: x.click_count, reverse=True)
                    top_links = [
                        {
                            "short_url": link.short_code,
                            "original_url": link.original_url,
                            "clicks": link.click_count,
                            "status": link.is_active,
                        }
                        for link in sorted_links[:5]  # Top 5 links
                    ]
                else:
                    data_stat = []
                    top_links = []
                    uniq_visitors = 0
                    total_clicks = 0

                return create_dashboard_response(
                    user_id=user_id,
                    total_clicks=total_clicks,
                    total_links=total_links,
                    top_links=top_links,
                    stat_links=data_stat,
                    uniq_visitors=uniq_visitors,
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
