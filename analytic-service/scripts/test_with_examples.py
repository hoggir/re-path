#!/usr/bin/env python3
"""
Test RPC contract using example JSON files.

This script demonstrates how to use the example JSON files for testing.
"""

import json
import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))

from app.schemas.rpc_contracts import (
    DashboardRequest,
    DashboardResponse,
    validate_dashboard_request,
)


def test_request_example():
    """Test that request example is valid."""
    print("=" * 60)
    print("Testing Request Example")
    print("=" * 60)

    example_file = Path(__file__).parent.parent / "schemas" / "dashboard_request.example.json"

    with open(example_file) as f:
        example_data = json.load(f)

    print(f"📄 Loaded: {example_file.name}")
    print(f"Content: {json.dumps(example_data, indent=2)}")

    try:
        validated_request = validate_dashboard_request(example_data)
        print(f"\n✅ Request example is VALID")
        print(f"   user_id: {validated_request.user_id}")
        return True
    except Exception as e:
        print(f"\n❌ Request example is INVALID: {e}")
        return False


def test_response_example():
    """Test that response example is valid."""
    print("\n" + "=" * 60)
    print("Testing Success Response Example")
    print("=" * 60)

    example_file = Path(__file__).parent.parent / "schemas" / "dashboard_response.example.json"

    with open(example_file) as f:
        example_data = json.load(f)

    print(f"📄 Loaded: {example_file.name}")
    print(f"Preview:")
    print(f"  user_id: {example_data.get('user_id')}")
    print(f"  total_clicks: {example_data.get('total_clicks')}")
    print(f"  total_links: {example_data.get('total_links')}")
    print(f"  status: {example_data.get('status')}")
    print(f"  recent_clicks count: {len(example_data.get('recent_clicks', []))}")
    print(f"  top_links count: {len(example_data.get('top_links', []))}")

    try:
        validated_response = DashboardResponse.model_validate(example_data)
        print(f"\n✅ Response example is VALID")
        print(f"   Validation passed all constraints:")
        print(f"   - user_id > 0: ✓")
        print(f"   - total_clicks >= 0: ✓")
        print(f"   - total_links >= 0: ✓")
        print(f"   - recent_clicks <= 10: ✓")
        print(f"   - top_links <= 5: ✓")
        print(f"   - status in enum: ✓")
        return True
    except Exception as e:
        print(f"\n❌ Response example is INVALID: {e}")
        return False


def test_error_example():
    """Test that error example is valid."""
    print("\n" + "=" * 60)
    print("Testing Error Response Example")
    print("=" * 60)

    example_file = Path(__file__).parent.parent / "schemas" / "dashboard_error.example.json"

    with open(example_file) as f:
        example_data = json.load(f)

    print(f"📄 Loaded: {example_file.name}")
    print(f"Content: {json.dumps(example_data, indent=2)}")

    try:
        validated_response = DashboardResponse.model_validate(example_data)
        print(f"\n✅ Error example is VALID")
        print(f"   status: {validated_response.status}")
        print(f"   message: {validated_response.message}")
        return True
    except Exception as e:
        print(f"\n❌ Error example is INVALID: {e}")
        return False


def demonstrate_go_usage():
    """Show how Go developers would use these examples."""
    print("\n" + "=" * 60)
    print("Go Usage Example")
    print("=" * 60)

    print("""
Go developers can use these JSON files for:

1. Testing JSON unmarshaling:
   ----------------------------
   data, _ := os.ReadFile("schemas/dashboard_response.example.json")
   var response domain.DashboardResponse
   if err := json.Unmarshal(data, &response); err != nil {
       t.Fatalf("Failed to parse: %v", err)
   }

2. Creating test cases:
   ---------------------
   func TestResponseParsing(t *testing.T) {
       // Load example
       data := loadExample("dashboard_response.example.json")

       // Parse
       var resp DashboardResponse
       json.Unmarshal(data, &resp)

       // Validate
       if err := resp.Validate(); err != nil {
           t.Error(err)
       }
   }

3. Manual RabbitMQ testing:
   -------------------------
   # Publish example request to RabbitMQ
   cat schemas/dashboard_request.example.json | \\
       rabbitmqadmin publish routing_key=dashboard_request

4. Reference for struct creation:
   ------------------------------
   Look at JSON structure → Create matching Go struct
    """)


def main():
    """Run all tests."""
    print("🧪 Testing RPC Contract Examples\n")

    results = []

    # Test all examples
    results.append(("Request Example", test_request_example()))
    results.append(("Response Example", test_response_example()))
    results.append(("Error Example", test_error_example()))

    # Show Go usage
    demonstrate_go_usage()

    # Summary
    print("\n" + "=" * 60)
    print("Summary")
    print("=" * 60)

    all_passed = all(result for _, result in results)

    for name, passed in results:
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"{status} - {name}")

    if all_passed:
        print("\n🎉 All examples are valid!")
        print("✅ JSON files can be used for:")
        print("   - Go struct reference")
        print("   - Manual testing")
        print("   - Integration tests")
        print("   - Documentation")
        return 0
    else:
        print("\n❌ Some examples are invalid!")
        print("⚠️  Please fix the examples before using them.")
        return 1


if __name__ == "__main__":
    exit(main())
