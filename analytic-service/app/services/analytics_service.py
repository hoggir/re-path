"""Analytics service for business logic."""

from datetime import datetime
from typing import List, Optional

from beanie import PydanticObjectId
from beanie.operators import In

from app.models.analytics_event import AnalyticsEvent


class AnalyticsService:
    """Service class for analytics operations."""

    @staticmethod
    async def create_event(
        event_type: str,
        user_id: Optional[str] = None,
        session_id: Optional[str] = None,
        page_url: Optional[str] = None,
        referrer: Optional[str] = None,
        user_agent: Optional[str] = None,
        ip_address: Optional[str] = None,
        metadata: Optional[dict] = None,
    ) -> AnalyticsEvent:
        """
        Create a new analytics event.

        Args:
            event_type: Type of event (e.g., page_view, click)
            user_id: User identifier
            session_id: Session identifier
            page_url: Page URL where event occurred
            referrer: HTTP referrer
            user_agent: User agent string
            ip_address: Client IP address
            metadata: Additional event metadata

        Returns:
            AnalyticsEvent: Created analytics event
        """
        event = AnalyticsEvent(
            event_type=event_type,
            user_id=user_id,
            session_id=session_id,
            page_url=page_url,
            referrer=referrer,
            user_agent=user_agent,
            ip_address=ip_address,
            metadata=metadata or {},
        )

        await event.insert()
        return event

    @staticmethod
    async def get_event_by_id(event_id: str) -> Optional[AnalyticsEvent]:
        """
        Get an analytics event by ID.

        Args:
            event_id: Event ID

        Returns:
            Optional[AnalyticsEvent]: Analytics event if found, None otherwise
        """
        try:
            return await AnalyticsEvent.get(PydanticObjectId(event_id))
        except Exception:
            return None

    @staticmethod
    async def get_events_by_type(
        event_type: str,
        limit: int = 100,
        skip: int = 0,
    ) -> List[AnalyticsEvent]:
        """
        Get analytics events by type.

        Args:
            event_type: Event type to filter
            limit: Maximum number of events to return
            skip: Number of events to skip (for pagination)

        Returns:
            List[AnalyticsEvent]: List of analytics events
        """
        return (
            await AnalyticsEvent.find(AnalyticsEvent.event_type == event_type)
            .sort(-AnalyticsEvent.created_at)
            .skip(skip)
            .limit(limit)
            .to_list()
        )

    @staticmethod
    async def get_events_by_user(
        user_id: str,
        limit: int = 100,
        skip: int = 0,
    ) -> List[AnalyticsEvent]:
        """
        Get analytics events by user ID.

        Args:
            user_id: User ID to filter
            limit: Maximum number of events to return
            skip: Number of events to skip (for pagination)

        Returns:
            List[AnalyticsEvent]: List of analytics events
        """
        return (
            await AnalyticsEvent.find(AnalyticsEvent.user_id == user_id)
            .sort(-AnalyticsEvent.created_at)
            .skip(skip)
            .limit(limit)
            .to_list()
        )

    @staticmethod
    async def get_events_by_session(
        session_id: str,
        limit: int = 100,
    ) -> List[AnalyticsEvent]:
        """
        Get analytics events by session ID.

        Args:
            session_id: Session ID to filter
            limit: Maximum number of events to return

        Returns:
            List[AnalyticsEvent]: List of analytics events
        """
        return (
            await AnalyticsEvent.find(AnalyticsEvent.session_id == session_id)
            .sort(AnalyticsEvent.created_at)
            .limit(limit)
            .to_list()
        )

    @staticmethod
    async def get_events_by_date_range(
        start_date: datetime,
        end_date: datetime,
        event_types: Optional[List[str]] = None,
        limit: int = 100,
        skip: int = 0,
    ) -> List[AnalyticsEvent]:
        """
        Get analytics events within a date range.

        Args:
            start_date: Start date
            end_date: End date
            event_types: Optional list of event types to filter
            limit: Maximum number of events to return
            skip: Number of events to skip (for pagination)

        Returns:
            List[AnalyticsEvent]: List of analytics events
        """
        query = AnalyticsEvent.find(
            AnalyticsEvent.created_at >= start_date,
            AnalyticsEvent.created_at <= end_date,
        )

        if event_types:
            query = query.find(In(AnalyticsEvent.event_type, event_types))

        return (
            await query.sort(-AnalyticsEvent.created_at)
            .skip(skip)
            .limit(limit)
            .to_list()
        )

    @staticmethod
    async def count_events(
        event_type: Optional[str] = None,
        user_id: Optional[str] = None,
        start_date: Optional[datetime] = None,
        end_date: Optional[datetime] = None,
    ) -> int:
        """
        Count analytics events with optional filters.

        Args:
            event_type: Optional event type filter
            user_id: Optional user ID filter
            start_date: Optional start date filter
            end_date: Optional end date filter

        Returns:
            int: Count of events
        """
        query = AnalyticsEvent.find()

        if event_type:
            query = query.find(AnalyticsEvent.event_type == event_type)

        if user_id:
            query = query.find(AnalyticsEvent.user_id == user_id)

        if start_date:
            query = query.find(AnalyticsEvent.created_at >= start_date)

        if end_date:
            query = query.find(AnalyticsEvent.created_at <= end_date)

        return await query.count()

    @staticmethod
    async def delete_event(event_id: str) -> bool:
        """
        Delete an analytics event by ID.

        Args:
            event_id: Event ID to delete

        Returns:
            bool: True if deleted, False if not found
        """
        try:
            event = await AnalyticsEvent.get(PydanticObjectId(event_id))
            if event:
                await event.delete()
                return True
            return False
        except Exception:
            return False
