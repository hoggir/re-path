"""Database models module."""

from typing import List

from app.models.analytics_event import AnalyticsEvent

# List of all Beanie document models
# Add new models here to register them with Beanie
ALL_MODELS: List[type] = [
    AnalyticsEvent,
]

__all__ = ["AnalyticsEvent", "ALL_MODELS"]
