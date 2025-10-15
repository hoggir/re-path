"""Analytics endpoints for data ingestion and search."""

from typing import Optional

from fastapi import APIRouter, HTTPException, Query, status

from app.schemas.analytics import (
    AnalyticsDataRequest,
    AnalyticsDataResponse,
    BulkAnalyticsDataRequest,
    BulkAnalyticsDataResponse,
    SearchRequest,
)
from app.schemas.response import ApiResponse, create_response
from app.services.opensearch_service import OpenSearchService

router = APIRouter(prefix="/elastic", tags=["Elastic"])


@router.post(
    "/ingest",
    response_model=ApiResponse[AnalyticsDataResponse],
    status_code=status.HTTP_201_CREATED,
    summary="Ingest Analytics Data",
    description="Submit a single JSON document to OpenSearch for analytics tracking",
)
async def ingest_analytics_data(
    request: AnalyticsDataRequest,
) -> ApiResponse[AnalyticsDataResponse]:
    """
    Ingest a single analytics data document into OpenSearch.

    This endpoint accepts any JSON data structure and stores it in the specified
    OpenSearch index type. A timestamp will be automatically added if not present.

    Args:
        request: Analytics data request containing index type, data, and optional doc ID

    Returns:
        ApiResponse containing the indexed document information

    Raises:
        HTTPException: If indexing fails
    """
    try:
        # Index the document
        response = await OpenSearchService.index_document(
            index_type=request.index_type,
            document=request.data,
            doc_id=request.doc_id,
        )

        # Create response data
        response_data = AnalyticsDataResponse(
            id=response.get("_id"),
            index=response.get("_index"),
            result=response.get("result"),
        )

        return create_response(
            data=response_data,
            message="Analytics data ingested successfully",
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to ingest analytics data: {str(e)}",
        ) from e


@router.post(
    "/ingest/bulk",
    response_model=ApiResponse[BulkAnalyticsDataResponse],
    status_code=status.HTTP_201_CREATED,
    summary="Bulk Ingest Analytics Data",
    description="Submit multiple JSON documents to OpenSearch in a single request",
)
async def bulk_ingest_analytics_data(
    request: BulkAnalyticsDataRequest,
) -> ApiResponse[BulkAnalyticsDataResponse]:
    """
    Ingest multiple analytics data documents into OpenSearch in bulk.

    This endpoint is optimized for high-throughput data ingestion. It accepts
    a list of JSON documents and indexes them in a single bulk operation.

    Args:
        request: Bulk analytics data request containing index type and documents

    Returns:
        ApiResponse containing bulk operation statistics

    Raises:
        HTTPException: If bulk indexing fails
    """
    try:
        # Bulk index documents
        response = await OpenSearchService.bulk_index(
            index_type=request.index_type,
            documents=request.documents,
        )

        # Parse bulk response
        items = response.get("items", [])
        items_count = len(items)
        errors = []
        success_count = 0

        for item in items:
            if "index" in item:
                if item["index"].get("status") in [200, 201]:
                    success_count += 1
                else:
                    errors.append(item["index"])

        error_count = items_count - success_count

        # Create response data
        response_data = BulkAnalyticsDataResponse(
            items_count=items_count,
            success_count=success_count,
            error_count=error_count,
            errors=errors if errors else None,
        )

        return create_response(
            data=response_data,
            message=f"Bulk ingestion completed: {success_count}/{items_count} documents indexed successfully",
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to bulk ingest analytics data: {str(e)}",
        ) from e


@router.post(
    "/search",
    response_model=ApiResponse[dict],
    status_code=status.HTTP_200_OK,
    summary="Search Analytics Data",
    description="Search analytics data in OpenSearch using query DSL",
)
async def search_analytics_data(
    request: SearchRequest,
) -> ApiResponse[dict]:
    """
    Search analytics data in OpenSearch.

    This endpoint allows you to query indexed analytics data using
    OpenSearch Query DSL for advanced filtering and aggregations.

    Args:
        request: Search request containing index type, query, and pagination

    Returns:
        ApiResponse containing search results

    Raises:
        HTTPException: If search fails
    """
    try:
        # Search documents
        response = await OpenSearchService.search(
            index_type=request.index_type,
            query=request.query,
            size=request.size,
            from_=request.from_,
        )

        return create_response(
            data=response,
            message="Search completed successfully",
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to search analytics data: {str(e)}",
        ) from e


@router.get(
    "/{index_type}/{doc_id}",
    response_model=ApiResponse[dict],
    status_code=status.HTTP_200_OK,
    summary="Get Analytics Document",
    description="Retrieve a specific analytics document by ID",
)
async def get_analytics_document(
    index_type: str,
    doc_id: str,
) -> ApiResponse[dict]:
    """
    Retrieve a specific analytics document by ID.

    Args:
        index_type: Type of index to search in
        doc_id: Document ID to retrieve

    Returns:
        ApiResponse containing the document

    Raises:
        HTTPException: If document not found or retrieval fails
    """
    try:
        document = await OpenSearchService.get_document(
            index_type=index_type,
            doc_id=doc_id,
        )

        if document is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Document with ID '{doc_id}' not found in index type '{index_type}'",
            )

        return create_response(
            data=document,
            message="Document retrieved successfully",
        )

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to retrieve document: {str(e)}",
        ) from e


@router.delete(
    "/{index_type}/{doc_id}",
    response_model=ApiResponse[dict],
    status_code=status.HTTP_200_OK,
    summary="Delete Analytics Document",
    description="Delete a specific analytics document by ID",
)
async def delete_analytics_document(
    index_type: str,
    doc_id: str,
) -> ApiResponse[dict]:
    """
    Delete a specific analytics document by ID.

    Args:
        index_type: Type of index to delete from
        doc_id: Document ID to delete

    Returns:
        ApiResponse containing deletion result

    Raises:
        HTTPException: If deletion fails
    """
    try:
        response = await OpenSearchService.delete_document(
            index_type=index_type,
            doc_id=doc_id,
        )

        return create_response(
            data=response,
            message="Document deleted successfully",
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to delete document: {str(e)}",
        ) from e


@router.get(
    "/clicks",
    response_model=ApiResponse[dict],
    status_code=status.HTTP_200_OK,
    summary="Get Click Data by ShortCode",
    description="Retrieve click analytics data filtered by shortCode from OpenSearch",
)
async def get_clicks_by_shortcode(
    short_code: str = Query(..., description="Short code to filter click data", alias="shortCode"),
    size: int = Query(10, ge=1, le=100, description="Number of results to return"),
    from_: int = Query(0, ge=0, description="Offset for pagination", alias="from"),
    sort_by: Optional[str] = Query("timestamp", description="Field to sort by"),
    sort_order: Optional[str] = Query("desc", description="Sort order (asc or desc)"),
) -> ApiResponse[dict]:
    """
    Get click analytics data filtered by shortCode.

    This endpoint retrieves all click events for a specific shortCode from the
    'click' index type in OpenSearch.

    Args:
        short_code: The short code to filter by (required)
        size: Number of results to return (default: 10, max: 100)
        from_: Offset for pagination (default: 0)
        sort_by: Field to sort by (default: timestamp)
        sort_order: Sort order - 'asc' or 'desc' (default: desc)

    Returns:
        ApiResponse containing click data and statistics

    Raises:
        HTTPException: If search fails
    """
    try:
        # Build OpenSearch query to match metadata.shortCode and event_type location
        # shortCode is nested inside metadata object
        query = {
            "bool": {
                "must": [
                    {"term": {"metadata.shortCode.keyword": short_code}},
                    {"term": {"event_type.keyword": "location"}}
                ]
            }
        }

        # Search in click index
        response = await OpenSearchService.search(
            index_type="click",
            query=query,
            size=size,
            from_=from_,
        )

        # Extract hits and total
        hits = response.get("hits", {})
        total = hits.get("total", {}).get("value", 0)
        documents = [hit.get("_source", {}) for hit in hits.get("hits", [])]

        # Prepare response data
        result = {
            "shortCode": short_code,
            "total_clicks": total,
            "clicks": documents,
            "pagination": {
                "size": size,
                "from": from_,
                "returned": len(documents)
            }
        }

        return create_response(
            data=result,
            message=f"Found {total} click(s) for shortCode '{short_code}'",
        )

    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to retrieve click data: {str(e)}",
        ) from e