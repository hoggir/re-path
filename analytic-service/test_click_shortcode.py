"""
Test script for getting click data by shortCode.

This script demonstrates how to:
1. Insert click data with shortCode
2. Retrieve click data by shortCode

To run:
1. Start the server: uvicorn app.main:app --reload
2. Run this script: python test_click_shortcode.py
"""

import requests
import json
from datetime import datetime

BASE_URL = "http://localhost:8000/api/v1/analytics"


def insert_sample_click_data():
    """Insert sample click data into OpenSearch."""
    print("\n=== Inserting Sample Click Data ===")

    url = f"{BASE_URL}/ingest/bulk"

    # Sample click data with different shortCodes
    # event_type must be "location" for click tracking
    # shortCode must be inside metadata object
    payload = {
        "index_type": "click",
        "documents": [
            {
                "event_type": "location",
                "metadata": {
                    "shortCode": "abc123",
                    "timestamp": datetime.utcnow().isoformat(),
                    "ip": "192.168.1.1",
                    "country": "ID",
                    "city": "Jakarta",
                    "device_type": "Desktop",
                    "referer": "https://google.com/",
                    "browser": "Chrome"
                }
            },
            {
                "event_type": "location",
                "metadata": {
                    "shortCode": "abc123",
                    "timestamp": datetime.utcnow().isoformat(),
                    "ip": "192.168.1.2",
                    "country": "ID",
                    "city": "Bandung",
                    "device_type": "Mobile",
                    "referer": "https://instagram.com/",
                    "browser": "Safari"
                }
            },
            {
                "event_type": "location",
                "metadata": {
                    "shortCode": "abc123",
                    "timestamp": datetime.utcnow().isoformat(),
                    "ip": "192.168.1.3",
                    "country": "ID",
                    "city": "Surabaya",
                    "device_type": "Desktop",
                    "referer": "https://facebook.com/",
                    "browser": "Firefox"
                }
            },
            {
                "event_type": "location",
                "metadata": {
                    "shortCode": "xyz789",
                    "timestamp": datetime.utcnow().isoformat(),
                    "ip": "192.168.1.4",
                    "country": "SG",
                    "city": "Singapore",
                    "device_type": "Mobile",
                    "referer": "https://twitter.com/",
                    "browser": "Chrome"
                }
            },
        ]
    }

    print(f"Request URL: {url}")
    print(f"Inserting {len(payload['documents'])} click records...")

    try:
        response = requests.post(url, json=payload)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


def get_clicks_by_shortcode(short_code: str, size: int = 10):
    """Get click data by shortCode."""
    print(f"\n\n=== Getting Clicks for ShortCode: {short_code} ===")

    url = f"{BASE_URL}/clicks/by-shortcode"
    params = {
        "shortCode": short_code,
        "size": size,
        "from": 0
    }

    print(f"Request URL: {url}")
    print(f"Parameters: {params}")

    try:
        response = requests.get(url, params=params)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Response:\n{json.dumps(response.json(), indent=2)}")
        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


def get_clicks_with_pagination(short_code: str):
    """Get click data with pagination example."""
    print(f"\n\n=== Getting Clicks with Pagination for ShortCode: {short_code} ===")

    url = f"{BASE_URL}/clicks/by-shortcode"

    # First page
    params = {
        "shortCode": short_code,
        "size": 2,
        "from": 0
    }

    print(f"Request URL: {url}")
    print(f"Getting first page (size=2, from=0)")

    try:
        response = requests.get(url, params=params)
        print(f"\nResponse Status: {response.status_code}")
        print(f"Page 1:\n{json.dumps(response.json(), indent=2)}")

        # Second page
        params["from"] = 2
        print(f"\n\nGetting second page (size=2, from=2)")
        response = requests.get(url, params=params)
        print(f"Page 2:\n{json.dumps(response.json(), indent=2)}")

        return response.json()
    except Exception as e:
        print(f"Error: {e}")
        return None


if __name__ == "__main__":
    print("=" * 70)
    print("Click Data by ShortCode - API Test")
    print("=" * 70)

    # Step 1: Insert sample data
    insert_result = insert_sample_click_data()

    # Wait a moment for indexing
    print("\n\nWaiting 2 seconds for OpenSearch to index the data...")
    import time
    time.sleep(2)

    # Step 2: Get clicks for shortCode "abc123"
    get_clicks_by_shortcode("abc123", size=10)

    # Step 3: Get clicks for shortCode "xyz789"
    get_clicks_by_shortcode("xyz789", size=10)

    # Step 4: Get clicks for non-existent shortCode
    get_clicks_by_shortcode("nonexistent", size=10)

    # Step 5: Test pagination
    get_clicks_with_pagination("abc123")

    print("\n" + "=" * 70)
    print("Tests Completed!")
    print("=" * 70)
