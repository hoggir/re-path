"""Tests for RPC contract schema validation."""

import json

import pytest
from pydantic import ValidationError

from app.schemas.rpc_contracts import (
    DashboardErrorResponse,
    DashboardRequest,
    DashboardResponse,
    RecentClickData,
    TopLinkData,
    create_dashboard_error_response,
    create_dashboard_response,
    validate_dashboard_request,
)


class TestDashboardRequest:
    """Test DashboardRequest schema validation."""

    def test_valid_request(self):
        """Test valid dashboard request."""
        data = {"user_id": 1}
        request = validate_dashboard_request(data)
        assert request.user_id == 1

    def test_valid_request_with_large_user_id(self):
        """Test request with large user ID."""
        data = {"user_id": 999999}
        request = validate_dashboard_request(data)
        assert request.user_id == 999999

    def test_missing_user_id(self):
        """Test request without user_id."""
        data = {}
        with pytest.raises(ValidationError) as exc_info:
            validate_dashboard_request(data)
        assert "user_id" in str(exc_info.value)

    def test_invalid_user_id_type(self):
        """Test request with invalid user_id type."""
        data = {"user_id": "not_an_int"}
        with pytest.raises(ValidationError):
            validate_dashboard_request(data)

    def test_zero_user_id(self):
        """Test request with user_id = 0 (invalid)."""
        data = {"user_id": 0}
        with pytest.raises(ValidationError) as exc_info:
            validate_dashboard_request(data)
        assert "greater than 0" in str(exc_info.value).lower()

    def test_negative_user_id(self):
        """Test request with negative user_id."""
        data = {"user_id": -1}
        with pytest.raises(ValidationError):
            validate_dashboard_request(data)

    def test_json_serialization(self):
        """Test that request can be JSON serialized."""
        request = DashboardRequest(user_id=123)
        json_str = request.model_dump_json()
        data = json.loads(json_str)
        assert data["user_id"] == 123


class TestDashboardResponse:
    """Test DashboardResponse schema validation."""

    def test_valid_success_response(self):
        """Test valid success response."""
        response = create_dashboard_response(
            user_id=1,
            total_clicks=100,
            total_links=10,
            recent_clicks=[
                {
                    "short_code": "test",
                    "clicked_at": "2025-10-20T10:00:00Z",
                    "is_bot": False,
                }
            ],
            top_links=[{"short_code": "test", "clicks": 50}],
            status="success",
        )

        assert response["user_id"] == 1
        assert response["total_clicks"] == 100
        assert response["total_links"] == 10
        assert response["status"] == "success"
        assert len(response["recent_clicks"]) == 1
        assert len(response["top_links"]) == 1

    def test_empty_lists(self):
        """Test response with empty lists."""
        response = create_dashboard_response(
            user_id=1,
            total_clicks=0,
            total_links=0,
            recent_clicks=[],
            top_links=[],
            status="success",
        )

        assert response["recent_clicks"] == []
        assert response["top_links"] == []

    def test_recent_clicks_limit(self):
        """Test that recent_clicks is limited to 10 items."""
        clicks = [
            {
                "short_code": f"link-{i}",
                "clicked_at": "2025-10-20T10:00:00Z",
                "is_bot": False,
            }
            for i in range(15)  # Try to add 15 items
        ]

        response = create_dashboard_response(
            user_id=1,
            total_clicks=15,
            total_links=15,
            recent_clicks=clicks,
            top_links=[],
            status="success",
        )

        # Should be limited to 10
        assert len(response["recent_clicks"]) == 10

    def test_top_links_limit(self):
        """Test that top_links is limited to 5 items."""
        links = [{"short_code": f"link-{i}", "clicks": i * 10} for i in range(10)]

        response = create_dashboard_response(
            user_id=1,
            total_clicks=100,
            total_links=10,
            recent_clicks=[],
            top_links=links,
            status="success",
        )

        # Should be limited to 5
        assert len(response["top_links"]) == 5

    def test_limited_status_response(self):
        """Test response with limited status."""
        response = create_dashboard_response(
            user_id=1,
            total_clicks=0,
            total_links=0,
            recent_clicks=[],
            top_links=[],
            status="limited",
            message="Database not available",
        )

        assert response["status"] == "limited"
        assert response["message"] == "Database not available"

    def test_invalid_status(self):
        """Test that invalid status is rejected."""
        with pytest.raises(ValidationError):
            DashboardResponse(
                user_id=1,
                total_clicks=0,
                total_links=0,
                recent_clicks=[],
                top_links=[],
                status="invalid_status",  # type: ignore
            )

    def test_negative_counts_rejected(self):
        """Test that negative counts are rejected."""
        with pytest.raises(ValidationError):
            DashboardResponse(
                user_id=1,
                total_clicks=-1,  # Invalid
                total_links=0,
                recent_clicks=[],
                top_links=[],
                status="success",
            )

    def test_json_serialization(self):
        """Test that response can be JSON serialized and deserialized."""
        response = create_dashboard_response(
            user_id=1,
            total_clicks=100,
            total_links=10,
            recent_clicks=[],
            top_links=[],
            status="success",
        )

        json_str = json.dumps(response)
        parsed = json.loads(json_str)

        assert parsed["user_id"] == 1
        assert parsed["total_clicks"] == 100
        assert parsed["status"] == "success"


class TestDashboardErrorResponse:
    """Test DashboardErrorResponse schema validation."""

    def test_valid_error_response(self):
        """Test valid error response."""
        response = create_dashboard_error_response(
            message="Something went wrong", user_id=1
        )

        assert response["status"] == "error"
        assert response["message"] == "Something went wrong"
        assert response["user_id"] == 1
        assert response["total_clicks"] == 0
        assert response["total_links"] == 0

    def test_error_response_without_user_id(self):
        """Test error response without user_id."""
        response = create_dashboard_error_response(message="Invalid request")

        assert response["status"] == "error"
        assert response["message"] == "Invalid request"
        # user_id should not be present if None
        assert "user_id" not in response or response["user_id"] is None


class TestRecentClickData:
    """Test RecentClickData schema validation."""

    def test_valid_recent_click(self):
        """Test valid recent click data."""
        click = RecentClickData(
            short_code="test",
            clicked_at="2025-10-20T10:00:00Z",
            ip_address_hash="abc123",
            user_agent="Mozilla/5.0",
            country_code="ID",
            city="Jakarta",
            device_type="desktop",
            browser_name="Chrome",
            is_bot=False,
        )

        assert click.short_code == "test"
        assert click.country_code == "ID"
        assert click.is_bot is False

    def test_minimal_recent_click(self):
        """Test recent click with only required fields."""
        click = RecentClickData(
            short_code="test",
            clicked_at="2025-10-20T10:00:00Z",
            is_bot=False,
        )

        assert click.short_code == "test"
        assert click.ip_address_hash is None
        assert click.user_agent is None


class TestTopLinkData:
    """Test TopLinkData schema validation."""

    def test_valid_top_link(self):
        """Test valid top link data."""
        link = TopLinkData(short_code="popular", clicks=100)

        assert link.short_code == "popular"
        assert link.clicks == 100

    def test_zero_clicks(self):
        """Test top link with zero clicks."""
        link = TopLinkData(short_code="unused", clicks=0)
        assert link.clicks == 0

    def test_negative_clicks_rejected(self):
        """Test that negative clicks are rejected."""
        with pytest.raises(ValidationError):
            TopLinkData(short_code="test", clicks=-1)


class TestGoCompatibility:
    """Test compatibility with Go service expectations."""

    def test_go_request_format(self):
        """Test that request matches Go DashboardRequest struct."""
        # Simulate what Go sends
        go_request = {"user_id": 123}

        # Python should be able to parse it
        request = validate_dashboard_request(go_request)
        assert request.user_id == 123

    def test_go_response_format(self):
        """Test that response matches what Go expects."""
        response = create_dashboard_response(
            user_id=123,
            total_clicks=1000,
            total_links=50,
            recent_clicks=[
                {
                    "short_code": "abc",
                    "clicked_at": "2025-10-20T10:00:00Z",
                    "country_code": "ID",
                    "is_bot": False,
                }
            ],
            top_links=[{"short_code": "popular", "clicks": 500}],
            status="success",
        )

        # Verify Go can unmarshal this (simulated)
        # Go unmarshals to map[string]interface{}
        assert isinstance(response, dict)
        assert "user_id" in response
        assert "total_clicks" in response
        assert "status" in response

        # All values should be JSON-serializable
        json_str = json.dumps(response)
        parsed = json.loads(json_str)

        # Verify structure
        assert parsed["user_id"] == 123
        assert parsed["total_clicks"] == 1000
        assert isinstance(parsed["recent_clicks"], list)
        assert isinstance(parsed["top_links"], list)

    def test_go_error_handling(self):
        """Test error response format for Go."""
        error_response = create_dashboard_error_response(
            message="Database error", user_id=123
        )

        # Go checks status field
        assert error_response["status"] == "error"
        assert "message" in error_response

        # Ensure it's JSON-serializable
        json_str = json.dumps(error_response)
        parsed = json.loads(json_str)
        assert parsed["status"] == "error"


if __name__ == "__main__":
    # Run with: python -m pytest tests/test_rpc_contract.py -v
    pytest.main([__file__, "-v"])
