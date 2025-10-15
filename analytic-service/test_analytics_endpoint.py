"""
Test script for analytics endpoints.

This script demonstrates how to use the analytics API endpoints.
To run this test, make sure the server is running first:
    uvicorn app.main:app --reload

Then run this script:
    python test_analytics_endpoint.py
"""

import requests
import json
from datetime import datetime

# Base URL for the API
BASE_URL = "http://localhost:8000/api/v1/analytics"


def test_single_ingest():
    """Test ingesting a single analytics document."""
    print("\n=== Testing Single Document Ingestion ===")

    url = f"{BASE_URL}/ingest"

    # Sample analytics data
    payload = {
        "index_type": "analytics",
        "data": {
            "user_id": "user123",
            "event_type": "page_view",
            "page_url": "https://example.com/products",
            "timestamp": datetime.utcnow().isoformat(),
            "metadata": {
                "category": "electronics",
                "product_id": "prod_456",
                "session_duration": 120
            }
        },
        "doc_id": None  # Let OpenSearch generate the ID
    }

    print(f"Request URL: {url}")
    print(f"Request Payload:\n{json.dumps(payload, indent=2)}")

    try:
        response = requests.post(url, json=payload)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response Body:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


def test_bulk_ingest():
    """Test ingesting multiple analytics documents in bulk."""
    print("\n\n=== Testing Bulk Document Ingestion ===")

    url = f"{BASE_URL}/ingest/bulk"

    # Sample bulk analytics data
    payload = {
        "index_type": "analytics",
        "documents": [
            {
                "user_id": "user123",
                "event_type": "page_view",
                "page_url": "https://example.com/products",
                "timestamp": datetime.utcnow().isoformat(),
            },
            {
                "user_id": "user456",
                "event_type": "click",
                "page_url": "https://example.com/cart",
                "button": "add_to_cart",
                "timestamp": datetime.utcnow().isoformat(),
            },
            {
                "user_id": "user789",
                "event_type": "conversion",
                "page_url": "https://example.com/checkout",
                "order_value": 99.99,
                "timestamp": datetime.utcnow().isoformat(),
            }
        ]
    }

    print(f"Request URL: {url}")
    print(f"Request Payload:\n{json.dumps(payload, indent=2)}")

    try:
        response = requests.post(url, json=payload)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response Body:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


def test_search():
    """Test searching analytics documents."""
    print("\n\n=== Testing Search ===")

    url = f"{BASE_URL}/search"

    # Sample search query
    payload = {
        "index_type": "analytics",
        "query": {
            "match": {
                "event_type": "page_view"
            }
        },
        "size": 10,
        "from": 0
    }

    print(f"Request URL: {url}")
    print(f"Request Payload:\n{json.dumps(payload, indent=2)}")

    try:
        response = requests.post(url, json=payload)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response Body:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


def test_get_document(doc_id: str):
    """Test retrieving a specific document by ID."""
    print(f"\n\n=== Testing Get Document (ID: {doc_id}) ===")

    url = f"{BASE_URL}/analytics/{doc_id}"

    print(f"Request URL: {url}")

    try:
        response = requests.get(url)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response Body:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


if __name__ == "__main__":
    print("=" * 60)
    print("Analytics API Endpoint Tests")
    print("=" * 60)

    # Test 1: Single document ingestion
    single_result = test_single_ingest()

    # Test 2: Bulk document ingestion
    bulk_result = test_bulk_ingest()

    # Test 3: Search documents
    search_result = test_search()

    # Test 4: Get document (if we have an ID from test 1)
    if single_result and single_result.get("data") and single_result["data"].get("id"):
        doc_id = single_result["data"]["id"]
        get_result = test_get_document(doc_id)

    print("\n" + "=" * 60)
    print("Tests Completed!")
    print("=" * 60)
