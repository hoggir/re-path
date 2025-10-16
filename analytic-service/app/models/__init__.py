"""Database models module."""

from app.models.click_event import ClickEvent

# List of all Beanie document models
# Add new models here to register them with Beanie
ALL_MODELS: list[type] = [
    ClickEvent,
]

__all__ = ["ClickEvent", "ALL_MODELS"]
