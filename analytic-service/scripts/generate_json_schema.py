#!/usr/bin/env python3
"""
Generate JSON Schema for RPC contracts.

This script generates JSON Schema files that can be used by the Go service
for validation or documentation purposes.

Usage:
    python scripts/generate_json_schema.py
"""

import json
import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.schemas.rpc_contracts import DashboardRequest, DashboardResponse


def generate_schemas():
    """Generate JSON Schema files for RPC contracts."""
    output_dir = Path("schemas")
    output_dir.mkdir(exist_ok=True)

    # Generate DashboardRequest schema
    request_schema = DashboardRequest.model_json_schema()
    request_file = output_dir / "dashboard_request.schema.json"
    with open(request_file, "w") as f:
        json.dump(request_schema, f, indent=2)
    print(f"✅ Generated: {request_file}")

    # Generate DashboardResponse schema
    response_schema = DashboardResponse.model_json_schema()
    response_file = output_dir / "dashboard_response.schema.json"
    with open(response_file, "w") as f:
        json.dump(response_schema, f, indent=2)
    print(f"✅ Generated: {response_file}")

    # Generate examples
    request_example = DashboardRequest(user_id=1).model_dump()
    request_example_file = output_dir / "dashboard_request.example.json"
    with open(request_example_file, "w") as f:
        json.dump(request_example, f, indent=2)
    print(f"✅ Generated: {request_example_file}")

    response_example = {
        "user_id": 1,
        "total_clicks": 1542,
        "total_links": 45,
        "recent_clicks": [
            {
                "short_code": "my-link",
                "clicked_at": "2025-10-20T10:00:00Z",
                "ip_address_hash": "abc123hash",
                "user_agent": "Mozilla/5.0 (X11; Linux x86_64) Chrome/120.0",
                "country_code": "ID",
                "city": "Jakarta",
                "device_type": "desktop",
                "browser_name": "Chrome",
                "is_bot": False,
            }
        ],
        "top_links": [{"short_code": "popular-link", "clicks": 350}],
        "status": "success",
    }
    response_example_file = output_dir / "dashboard_response.example.json"
    with open(response_example_file, "w") as f:
        json.dump(response_example, f, indent=2)
    print(f"✅ Generated: {response_example_file}")

    print("\n📋 Schema files generated successfully!")
    print(f"📁 Location: {output_dir.absolute()}")
    print("\nThese files can be shared with the Go team for reference.")


if __name__ == "__main__":
    generate_schemas()
