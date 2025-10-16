"""Global response schemas for standardized API responses."""

from datetime import datetime
from typing import Any, Generic, Optional, TypeVar

from pydantic import BaseModel, Field

# Generic type for data payload
DataT = TypeVar("DataT")


class Meta(BaseModel):
    """Metadata information for response."""

    timestamp: datetime = Field(
        default_factory=datetime.utcnow,
        description="Response timestamp in UTC",
    )
    version: str = Field(
        default="v1",
        description="API version",
    )
    request_id: Optional[str] = Field(
        default=None,
        description="Unique request identifier for tracing",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "timestamp": "2025-10-14T10:30:00Z",
                "version": "v1",
                "request_id": "req_abc123xyz",
            }
        }
    }


class ErrorDetail(BaseModel):
    """Detailed error information."""

    code: str = Field(
        description="Error code for client handling",
        examples=["VALIDATION_ERROR", "NOT_FOUND", "INTERNAL_ERROR"],
    )
    message: str = Field(
        description="Human-readable error message",
        examples=["Invalid input data", "Resource not found"],
    )
    field: Optional[str] = Field(
        default=None,
        description="Field name if validation error",
        examples=["email", "age"],
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "code": "VALIDATION_ERROR",
                "message": "Invalid email format",
                "field": "email",
            }
        }
    }


class ApiResponse(BaseModel, Generic[DataT]):
    """
    Standard API response wrapper.

    This provides a consistent response structure across all endpoints.
    """

    success: bool = Field(
        description="Indicates if the request was successful",
        examples=[True, False],
    )
    message: str = Field(
        description="Response message describing the result",
        examples=["Operation successful", "Data retrieved successfully"],
    )
    data: Optional[DataT] = Field(
        default=None,
        description="Response payload data",
    )
    meta: Meta = Field(
        default_factory=Meta,
        description="Response metadata",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "success": True,
                "message": "Operation successful",
                "data": {"key": "value"},
                "meta": {
                    "timestamp": "2025-10-14T10:30:00Z",
                    "version": "v1",
                    "request_id": "req_abc123xyz",
                },
            }
        }
    }


class ErrorResponse(BaseModel):
    """
    Standard error response.

    Used for error cases with detailed error information.
    """

    success: bool = Field(
        default=False,
        description="Always false for error responses",
    )
    message: str = Field(
        description="General error message",
        examples=["Request failed", "Validation error occurred"],
    )
    errors: list[ErrorDetail] = Field(
        default_factory=list,
        description="List of detailed errors",
    )
    meta: Meta = Field(
        default_factory=Meta,
        description="Response metadata",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "success": False,
                "message": "Validation error occurred",
                "errors": [
                    {
                        "code": "VALIDATION_ERROR",
                        "message": "Invalid email format",
                        "field": "email",
                    }
                ],
                "meta": {
                    "timestamp": "2025-10-14T10:30:00Z",
                    "version": "v1",
                    "request_id": "req_abc123xyz",
                },
            }
        }
    }


class PaginationMeta(BaseModel):
    """Pagination metadata for list responses."""

    page: int = Field(
        ge=1,
        description="Current page number",
        examples=[1],
    )
    page_size: int = Field(
        ge=1,
        le=100,
        description="Number of items per page",
        examples=[10, 20, 50],
    )
    total_items: int = Field(
        ge=0,
        description="Total number of items",
        examples=[100],
    )
    total_pages: int = Field(
        ge=0,
        description="Total number of pages",
        examples=[10],
    )
    has_next: bool = Field(
        description="Indicates if there is a next page",
    )
    has_previous: bool = Field(
        description="Indicates if there is a previous page",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "page": 1,
                "page_size": 10,
                "total_items": 100,
                "total_pages": 10,
                "has_next": True,
                "has_previous": False,
            }
        }
    }


class PaginatedResponse(BaseModel, Generic[DataT]):
    """
    Paginated response for list endpoints.

    Includes pagination metadata along with data.
    """

    success: bool = Field(
        default=True,
        description="Indicates if the request was successful",
    )
    message: str = Field(
        description="Response message",
        examples=["Data retrieved successfully"],
    )
    data: list[DataT] = Field(
        default_factory=list,
        description="List of items",
    )
    pagination: PaginationMeta = Field(
        description="Pagination metadata",
    )
    meta: Meta = Field(
        default_factory=Meta,
        description="Response metadata",
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "success": True,
                "message": "Data retrieved successfully",
                "data": [{"id": 1, "name": "Item 1"}, {"id": 2, "name": "Item 2"}],
                "pagination": {
                    "page": 1,
                    "page_size": 10,
                    "total_items": 100,
                    "total_pages": 10,
                    "has_next": True,
                    "has_previous": False,
                },
                "meta": {
                    "timestamp": "2025-10-14T10:30:00Z",
                    "version": "v1",
                },
            }
        }
    }


# Helper function to create success response
def create_response(
    data: Any = None,
    message: str = "Operation successful",
    request_id: Optional[str] = None,
) -> ApiResponse[Any]:
    """
    Helper function to create a successful API response.

    Args:
        data: Response payload data
        message: Success message
        request_id: Optional request ID for tracing

    Returns:
        ApiResponse: Standardized success response
    """
    return ApiResponse(
        success=True,
        message=message,
        data=data,
        meta=Meta(request_id=request_id),
    )


# Helper function to create error response
def create_error_response(
    message: str,
    errors: Optional[list[ErrorDetail]] = None,
    request_id: Optional[str] = None,
) -> ErrorResponse:
    """
    Helper function to create an error response.

    Args:
        message: Error message
        errors: List of detailed errors
        request_id: Optional request ID for tracing

    Returns:
        ErrorResponse: Standardized error response
    """
    return ErrorResponse(
        message=message,
        errors=errors or [],
        meta=Meta(request_id=request_id),
    )
