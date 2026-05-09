# Go Backend - Week 3 Progress (MongoDB Integration)

This document summarizes the Week 3 backend progress for the inventory API, focused on the transition from in-memory storage to MongoDB persistence.

## 1. Week 3 Overview

- The inventory module moved from a map-based in-memory repository to MongoDB-backed persistent storage.
- The refactor followed clean architecture principles: handlers and controllers remained largely stable, while the repository implementation was replaced.
- This preserved API behavior while improving storage durability and production readiness.

## 2. MongoDB Integration

- MongoDB is now used as the backing data store for product inventory.
- The backend uses the official Go driver: `go.mongodb.org/mongo-driver`.
- Product data is stored in the `inventory.products` collection.
- Repository operations now execute MongoDB queries instead of map reads/writes.
- Connection startup includes both connect and health-check steps:

```go
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
if err != nil {
    return nil, err
}

if err := client.Ping(ctx, nil); err != nil {
    return nil, err
}
```

## 3. Repository Refactor

- The repository interface contract was preserved, so controller and handler call sites continue to use the same abstraction.
- A new `MongoProductRepository` now implements this interface.
- This validates the architecture's decoupling: infrastructure can change without rewriting business logic.

## 4. Use of BSON

- BSON (Binary JSON) is MongoDB's native document format.
- Queries and updates use `bson.M`, for example:
  - `bson.M{"id": id}` for filters
  - `bson.M{"$set": product}` for updates
- The `Product` entity should include `bson` tags alongside JSON tags to ensure explicit field mapping between Go structs and MongoDB documents.

## 5. Context Propagation

- `context.Context` was introduced across the inventory flow.
- Context now flows end-to-end: handler -> controller -> repository -> MongoDB driver.
- This enables cancellation and deadline propagation from HTTP requests down to database operations.

## 6. Timeout Handling

- Repository methods use `context.WithTimeout` for database operations.
- This guards against hanging or slow MongoDB calls and helps keep request latency bounded.
- Timeouts also improve system resilience under network or database contention.

## 7. Improved Error Handling

- MongoDB-specific outcomes are now checked explicitly:
  - `MatchedCount` in `UpdateOne`
  - `DeletedCount` in `DeleteOne`
- If no record is matched/deleted, the repository returns a not-found error.
- This enables correct HTTP-level behavior for missing resources.

## 8. API Behavior Improvements

- `GET /products` now returns an empty array (`[]`) when no records exist, instead of `null`.
- `DELETE /products/{id}` returns `204 No Content` on success.
- Missing resources correctly return `404 Not Found`.

## 9. Testing

- APIs were validated using Postman and `curl` for create/read/update/delete flows.
- Stored data was verified in MongoDB Compass.
- Manual tests confirmed behavior for both success and not-found scenarios.

## 10. Summary

- Week 3 established persistent inventory storage with MongoDB.
- The codebase now demonstrates backend patterns commonly used in production services:
  - repository abstraction
  - context-aware database operations
  - timeout-protected I/O
  - behavior-driven HTTP error mapping

The backend is now significantly closer to a production-ready service than the initial in-memory implementation.

---

## Categories and Bulk Creation Fix

### ID Generation

Previously, products required the caller to supply an `id` field — the backend had no fallback. This meant:
- A `POST /products` body missing `id` would insert a document with an empty-string ID, and a second such request would collide.
- The bulk CSV import required an ID column that callers had to populate with unique values themselves.

**Fix:** `MongoProductRepository.Create` now generates a MongoDB ObjectID hex string (`primitive.NewObjectID().Hex()`) when `product.ID` is empty, matching the existing behavior in `MongoCategoryRepository.Create`. The service and handler are unchanged — ID generation stays in the repository layer where infrastructure concerns belong.

Categories already generated IDs on the backend and required no changes.

### Bulk Create Products — `POST /products/bulk`

**Endpoint:** `POST /products/bulk`  
**Content-Type:** `multipart/form-data`  
**Field:** `file` — a CSV file

**CSV format** (header row required):

```
name,price,quantity,brand
Milk,50,10,Amul
Phone,5000,2,Samsung
Headphones,1200,15,Sony
```

- No ID column — the backend generates all IDs.
- Rows with fewer than 4 columns or unparseable price/quantity are silently skipped before processing.
- Up to **10 rows are processed concurrently** via a worker-pool semaphore, preventing unbounded goroutine creation on large files.
- Errors per row are collected and returned rather than dropped silently.

**Response: 207 Multi-Status**

```json
{
  "created": [
    { "id": "68...", "name": "Milk", "price": 50, "quantity": 10, "brand": "Amul", ... }
  ],
  "errors": [
    { "index": 1, "reason": "Brand is required" }
  ]
}
```

- `created` — products that were successfully inserted.
- `errors` — per-item failures; `index` is the 0-based position in the submitted list (after filtering malformed rows). The response is always 207 so callers can always inspect both fields.

**Problems fixed vs. the previous implementation:**

| Problem | Before | After |
|---------|--------|-------|
| Goroutine count | One per row, unbounded | Capped at 10 via semaphore |
| Per-row errors | Silently dropped | Returned in `errors` array |
| Overall error return | Always `nil` | System errors surfaced; item errors in body |
| Response status | 201 Created | 207 Multi-Status |
| ID responsibility | Caller must supply | Backend generates |
| CSV ID column | Required (`row[0]`) | Removed |
