"""Dashboard endpoints."""

import logging
from typing import Optional

from fastapi import APIRouter, HTTPException, Query, status
from datetime import datetime, timedelta

from app.models.click_event import ClickEvent
from app.schemas.dashboard import (
    ClickEventsDashboard,
    ClickEventsStats,
)
from app.schemas.response import ApiResponse, create_response

router = APIRouter(tags=["Dashboard"])
logger = logging.getLogger(__name__)


@router.get(
    "/dashboard/clicks",
    response_model=ApiResponse[ClickEventsDashboard],
    status_code=status.HTTP_200_OK,
    summary="Get Click Events Dashboard",
    description="Retrieve click events from MongoDB with statistics and filtering",
)
async def get_click_events_dashboard(
    short_code: Optional[str] = Query(None, description="Filter by short code"),
) -> ApiResponse[ClickEventsDashboard]:
    try:
        query_filter = {}
        if short_code:
            query_filter["shortCode"] = short_code

        total_clicks = await ClickEvent.find(query_filter).count()
        bot_clicks = await ClickEvent.find({**query_filter, "isBot": True}).count()
        human_clicks = total_clicks - bot_clicks

        country_stats = []
        # Fallback: fetch all and count manually
        all_events_for_countries = await ClickEvent.find(query_filter).to_list()
        events_as_dict = [event.dict() for event in all_events_for_countries]
        print(events_as_dict)
        country_counts = {}
        for event in all_events_for_countries:
            if event.country_code:
                country_counts[event.country_code] = country_counts.get(event.country_code, 0) + 1
        unique_countries = len(country_counts)
        # Convert to CountryClickStats and sort by click count
        # country_stats = [
        #     CountryClickStats(country_code=code, click_count=count)
        #     for code, count in sorted(
        #         country_counts.items(), key=lambda x: x[1], reverse=True
        #     )
        # ]

        events = [
            {
                **e.dict(),
                "clicked_at": e.clicked_at.strftime("%d-%m-%Y"),
                "clicked_date": e.clicked_at.date()  # simpan juga tanggal untuk filter mudah
            }
            for e in all_events_for_countries
        ]
        # Get country-based click statistics using aggregation
        # dapatkan rentang tanggal dari tanggal 1 bulan ini sampai hari ini
        today = datetime.now().date()
        start_date = today.replace(day=1)

        current_date = start_date
        while current_date <= today:
            # filter event sesuai tanggal loop
            events_on_date = [ev for ev in events if ev["clicked_date"] == current_date]

            print(f"Tanggal: {current_date} — Jumlah klik: {len(events_on_date)}")
            # kalau mau lihat detail event:
            # for ev in events_on_date:
            #     print(ev)

            # lanjut ke tanggal berikutnya
            current_date += timedelta(days=1)

        stats = ClickEventsStats(
            total_clicks=total_clicks,
            unique_visitors=666,
            unique_countries=unique_countries,
            bot_clicks=bot_clicks,
            human_clicks=human_clicks,
        )

        # Create dashboard response
        dashboard = ClickEventsDashboard(
            short_code=short_code,
            stats=stats,
        )

        return create_response(
            data=dashboard,
            message=f"Found {total_clicks} click event(s)"
            + (f" for short code '{short_code}'" if short_code else ""),
        )

    except Exception as e:
        logger.error(f"Error retrieving click events: {e}", exc_info=True)
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to retrieve click events: {str(e)}",
        ) from e
