"""Dashboard schemas."""

from typing import Optional

from pydantic import BaseModel, Field


class ClickEventsStats(BaseModel):
    """Click events statistics schema."""

    total_clicks: int = Field(description="Total number of clicks")
    unique_visitors: int = Field(description="Number of unique visitors")
    unique_countries: int = Field(description="Number of unique countries")
    bot_clicks: int = Field(description="Number of bot clicks")
    human_clicks: int = Field(description="Number of human clicks")


class ClickEventsDashboard(BaseModel):
    """Dashboard response with click events."""

    short_code: Optional[str] = Field(default=None, description="Filter by short code")
    stats: ClickEventsStats = Field(description="Click events statistics")
