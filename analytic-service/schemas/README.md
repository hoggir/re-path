# JSON Schema Examples

## Purpose

File-file JSON ini adalah **example data** dan **quick reference** untuk:
1. **Go developers** - Copy-paste untuk testing
2. **Documentation** - Menunjukkan struktur data yang expected
3. **Testing** - Manual testing via RabbitMQ
4. **Validation** - Memastikan format benar

---

## Files Overview

| File | Purpose | Usage |
|------|---------|-------|
| `dashboard_request.example.json` | Request format | Test RPC call, Go struct reference |
| `dashboard_response.example.json` | Success response | Expected Python output, Go parsing test |
| `dashboard_error.example.json` | Error response | Error handling reference |

---

## 📖 Use Case 1: Go Struct Reference

**Problem:** Go developer perlu tahu struktur JSON yang akan diterima dari Python.

**Solution:** Lihat example files untuk struktur lengkap.

### Example: Creating Go Struct

**Look at:** `dashboard_response.example.json`
```json
{
  "user_id": 1,
  "total_clicks": 1542,
  "total_links": 45,
  "recent_clicks": [...],
  "top_links": [...],
  "status": "success"
}
```

**Create Go struct:**
```go
type DashboardResponse struct {
    UserID       int           `json:"user_id"`       // From example
    TotalClicks  int           `json:"total_clicks"`  // From example
    TotalLinks   int           `json:"total_links"`   // From example
    RecentClicks []RecentClick `json:"recent_clicks"` // From example
    TopLinks     []TopLink     `json:"top_links"`     // From example
    Status       string        `json:"status"`        // From example
}
```

**Benefit:** Visual reference tanpa perlu baca Pydantic schema.

---

## 🧪 Use Case 2: Manual Testing

**Problem:** Test RPC integration tanpa run semua service.

**Solution:** Publish example JSON ke RabbitMQ queue secara manual.

### Testing Dashboard RPC

#### 1. Start analytic-service
```bash
cd analytic-service
uvicorn app.main:app --reload
```

#### 2. Publish Test Request

**Using Python:**
```python
import pika
import json

# Load example request
with open('schemas/dashboard_request.example.json') as f:
    request = json.load(f)

# Connect to RabbitMQ
connection = pika.BlockingConnection(
    pika.URLParameters('amqp://repath:repath123@localhost:5672/repath')
)
channel = connection.channel()

# Create reply queue
result = channel.queue_declare(queue='', exclusive=True)
reply_queue = result.method.queue

correlation_id = 'test-12345'

# Publish request
channel.basic_publish(
    exchange='',
    routing_key='dashboard_request',
    properties=pika.BasicProperties(
        reply_to=reply_queue,
        correlation_id=correlation_id,
    ),
    body=json.dumps(request)
)

print(f"✅ Sent request: {request}")

# Wait for response
def on_response(ch, method, props, body):
    if props.correlation_id == correlation_id:
        response = json.loads(body.decode())
        print(f"✅ Received response:")
        print(json.dumps(response, indent=2))
        connection.close()

channel.basic_consume(
    queue=reply_queue,
    on_message_callback=on_response,
    auto_ack=True
)

print("⏳ Waiting for response...")
channel.start_consuming()
```

#### 3. Verify Response Matches Example

Compare response dengan `dashboard_response.example.json`:
```bash
# Response should match structure
{
  "user_id": 1,           # ✓ Same
  "total_clicks": ...,    # ✓ Same type (int)
  "total_links": ...,     # ✓ Same type (int)
  "status": "success"     # ✓ Same
}
```

---

## 🔍 Use Case 3: Integration Testing

**Problem:** Go team perlu verify parsing logic.

**Solution:** Load example JSON dan test unmarshaling.

### Go Test Example

```go
package domain_test

import (
    "encoding/json"
    "os"
    "testing"

    "github.com/hoggir/re-path/redirect-service/internal/domain"
)

func TestDashboardResponseParsing(t *testing.T) {
    // Load example JSON
    data, err := os.ReadFile("../../schemas/dashboard_response.example.json")
    if err != nil {
        t.Fatalf("Failed to read example: %v", err)
    }

    // Parse into Go struct
    var response domain.DashboardResponse
    if err := json.Unmarshal(data, &response); err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    // Validate structure
    if err := response.Validate(); err != nil {
        t.Errorf("Validation failed: %v", err)
    }

    // Verify values
    if response.UserID != 1 {
        t.Errorf("UserID = %d, want 1", response.UserID)
    }
    if response.Status != "success" {
        t.Errorf("Status = %s, want success", response.Status)
    }

    t.Logf("✅ Successfully parsed example JSON")
}

func TestErrorResponseParsing(t *testing.T) {
    // Load error example
    data, err := os.ReadFile("../../schemas/dashboard_error.example.json")
    if err != nil {
        t.Fatalf("Failed to read example: %v", err)
    }

    var response domain.DashboardResponse
    if err := json.Unmarshal(data, &response); err != nil {
        t.Fatalf("Failed to unmarshal: %v", err)
    }

    // Verify error status
    if !response.IsError() {
        t.Errorf("Expected error status")
    }

    expectedMsg := "Database connection failed"
    if response.GetMessage() != expectedMsg {
        t.Errorf("Message = %s, want %s", response.GetMessage(), expectedMsg)
    }

    t.Logf("✅ Error response parsed correctly")
}
```

**Run test:**
```bash
cd redirect-service
go test ./internal/domain -v
```

---

## 📝 Use Case 4: Documentation

**Problem:** New developer perlu understand data format.

**Solution:** Read example files untuk quick reference.

### For Go Developers

**Question:** "What does a successful response look like?"

**Answer:** Look at `dashboard_response.example.json`
```json
{
  "user_id": 1,
  "total_clicks": 1542,
  "total_links": 45,
  "recent_clicks": [
    {
      "short_code": "my-link",
      "clicked_at": "2025-10-20T10:00:00Z",
      "country_code": "ID",
      "city": "Jakarta",
      "is_bot": false
    }
  ],
  "top_links": [
    {"short_code": "popular-link", "clicks": 350}
  ],
  "status": "success"
}
```

**Quick insights:**
- ✅ `user_id` is integer
- ✅ `recent_clicks` is array of objects
- ✅ `status` is string ("success", "error", "limited")
- ✅ `is_bot` is boolean
- ✅ Optional fields like `country_code` can be null

---

## 🔄 Use Case 5: Contract Verification

**Problem:** Verify Python output matches contract.

**Solution:** Compare actual response with example.

### Verification Script

```python
# scripts/verify_contract.py
import json
from app.schemas.rpc_contracts import DashboardResponse

# Load example
with open('schemas/dashboard_response.example.json') as f:
    example = json.load(f)

# Validate against schema
try:
    validated = DashboardResponse.model_validate(example)
    print("✅ Example matches schema")
except Exception as e:
    print(f"❌ Example doesn't match schema: {e}")

# Check all required fields present
required_fields = ['user_id', 'total_clicks', 'total_links', 'status']
for field in required_fields:
    if field not in example:
        print(f"❌ Missing required field: {field}")
    else:
        print(f"✅ {field} present")
```

**Run:**
```bash
cd analytic-service
python scripts/verify_contract.py
```

---

## 🛠️ Maintaining Examples

### When to Update

Update example files when:
1. ✅ Adding new fields to schema
2. ✅ Changing field types
3. ✅ Changing validation rules
4. ✅ Adding new response types

### How to Update

#### Manual Update
```bash
# Edit JSON files directly
vim schemas/dashboard_response.example.json
```

#### Auto-generate (Future)
```python
# scripts/generate_examples.py
from app.schemas.rpc_contracts import DashboardResponse

# Create example
example = DashboardResponse(
    user_id=1,
    total_clicks=1542,
    total_links=45,
    recent_clicks=[...],
    top_links=[...],
    status="success"
)

# Save to file
with open('schemas/dashboard_response.example.json', 'w') as f:
    json.dump(example.model_dump(mode='json'), f, indent=2)
```

---

## 📋 Quick Reference

### File Naming Convention
```
{entity}_{type}.example.json

Examples:
- dashboard_request.example.json   → Request example
- dashboard_response.example.json  → Success response example
- dashboard_error.example.json     → Error response example
- user_request.example.json        → User request (future)
```

### File Structure
```
schemas/
├── README.md                           # This file
├── dashboard_request.example.json      # Request format
├── dashboard_response.example.json     # Success response
├── dashboard_error.example.json        # Error response
└── (future) dashboard_*.schema.json    # JSON Schema definitions
```

---

## 🎯 Best Practices

### 1. Keep Examples Realistic
```json
// ✅ Good - Realistic data
{
  "user_id": 123,
  "total_clicks": 1542,
  "status": "success"
}

// ❌ Bad - Placeholder data
{
  "user_id": 1,
  "total_clicks": 999999,
  "status": "test"
}
```

### 2. Include Edge Cases
```json
// Example with minimal data
{
  "user_id": 1,
  "total_clicks": 0,
  "total_links": 0,
  "recent_clicks": [],
  "top_links": [],
  "status": "success"
}

// Example with full data
{
  "user_id": 123,
  "total_clicks": 1542,
  "recent_clicks": [...],  // Max 10 items
  "top_links": [...],       // Max 5 items
  "status": "success"
}
```

### 3. Validate Examples
```bash
# Before committing, validate examples
cd analytic-service
python -c "
import json
from app.schemas.rpc_contracts import DashboardResponse

with open('schemas/dashboard_response.example.json') as f:
    data = json.load(f)

DashboardResponse.model_validate(data)
print('✅ Example is valid')
"
```

### 4. Document Special Cases
```json
{
  "user_id": 1,
  "status": "limited",
  "message": "Database not available",  // Only present on error/limited
  "total_clicks": 0,
  "total_links": 0,
  "recent_clicks": [],
  "top_links": []
}
```

---

## 🚀 Quick Start

### For Go Developers

1. **See request format:**
   ```bash
   cat schemas/dashboard_request.example.json
   ```

2. **See response format:**
   ```bash
   cat schemas/dashboard_response.example.json
   ```

3. **Test parsing:**
   ```go
   data, _ := os.ReadFile("schemas/dashboard_response.example.json")
   var resp DashboardResponse
   json.Unmarshal(data, &resp)
   ```

### For Python Developers

1. **Validate example:**
   ```python
   from app.schemas.rpc_contracts import DashboardResponse
   import json

   with open('schemas/dashboard_response.example.json') as f:
       DashboardResponse.model_validate(json.load(f))
   ```

2. **Generate new example:**
   ```python
   from app.schemas.rpc_contracts import create_dashboard_response

   example = create_dashboard_response(
       user_id=1,
       total_clicks=100,
       total_links=10,
       recent_clicks=[],
       top_links=[],
       status="success"
   )

   with open('schemas/new_example.json', 'w') as f:
       json.dump(example, f, indent=2)
   ```

---

## Summary

**JSON example files serve as:**
1. ✅ **Quick Reference** - Visual structure tanpa baca code
2. ✅ **Testing Tool** - Manual testing via RabbitMQ
3. ✅ **Documentation** - Living examples of data format
4. ✅ **Validation** - Ensure Go/Python compatibility
5. ✅ **Onboarding** - Help new developers understand format

**Update examples whenever schema changes!**
