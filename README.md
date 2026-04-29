# Interneers Lab Inventory API

This document covers the current inventory system API implemented in `backend/go`.

## Scope

The inventory system is the `products` module in the Go backend. It exposes a small CRUD API over HTTP and is currently wired into the application entrypoint at `backend/go/cmd/app/main.go`.

Current implementation notes:

- The API is served by Go's standard `net/http` package.
- Product routes are mounted under `/products`.
- Data is stored in an in-memory map, not in MongoDB.
- The process currently listens on port `8080` in code.

## Inventory Domain Model

Each inventory item is represented as a product with the following JSON shape:

```json
{
  "id": "sku-1001",
  "name": "Wireless Mouse",
  "description": "Ergonomic mouse with USB receiver",
  "category": "Accessories",
  "price": 24.99,
  "brand": "LogiTech",
  "quantity": 12
}
```

### Fields

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | `string` | Yes in practice | Must be supplied by the client. The API does not generate IDs. |
| `name` | `string` | Yes | Must not be empty. |
| `description` | `string` | No | Free-form text. |
| `category` | `string` | No | Free-form text. |
| `price` | `number` | Yes | Must be greater than `0`. |
| `brand` | `string` | No | Free-form text. |
| `quantity` | `number` | Yes | Must be `0` or greater. |

## Validation Rules

The controller currently enforces these business rules on create and update:

- `name` is required.
- `price` must be greater than `0`.
- `quantity` cannot be negative.

Validation failures return plain-text error messages with `400 Bad Request` for create requests and either `400 Bad Request` or `404 Not Found` depending on the handler path.

## Base URL

Use this base URL when running the Go server locally:

```text
http://localhost:8080
```

## Endpoints

### 1. List all inventory items

**Request**

```http
GET /products
```

**Success response**

- Status: `200 OK`
- Body: JSON array of products

**Example**

```bash
curl http://localhost:8080/products
```

### 2. Get one inventory item by ID

**Request**

```http
GET /products/{id}
```

**Success response**

- Status: `200 OK`
- Body: JSON product object

**Not found**

- Status: `404 Not Found`
- Body: `Product not found`

**Example**

```bash
curl http://localhost:8080/products/sku-1001
```

### 3. Create an inventory item

**Request**

```http
POST /products
Content-Type: application/json
```

**Example body**

```json
{
  "id": "sku-1001",
  "name": "Wireless Mouse",
  "description": "Ergonomic mouse with USB receiver",
  "category": "Accessories",
  "price": 24.99,
  "brand": "LogiTech",
  "quantity": 12
}
```

**Success response**

- Status: `201 Created`
- Body: created product as JSON

**Example**

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "id": "sku-1001",
    "name": "Wireless Mouse",
    "description": "Ergonomic mouse with USB receiver",
    "category": "Accessories",
    "price": 24.99,
    "brand": "LogiTech",
    "quantity": 12
  }'
```

### 4. Update an inventory item

**Request**

```http
PUT /products/{id}
Content-Type: application/json
```

**Example body**

```json
{
  "name": "Wireless Mouse Pro",
  "description": "Updated ergonomic mouse",
  "category": "Accessories",
  "price": 29.99,
  "brand": "LogiTech",
  "quantity": 8
}
```

**Success response**

- Status: `200 OK`
- Body: updated product as JSON

**Example**

```bash
curl -X PUT http://localhost:8080/products/sku-1001 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Mouse Pro",
    "description": "Updated ergonomic mouse",
    "category": "Accessories",
    "price": 29.99,
    "brand": "LogiTech",
    "quantity": 8
  }'
```

### 5. Delete an inventory item

**Request**

```http
DELETE /products/{id}
```

**Success response**

- Status: `204 No Content`

**Not found**

- Status: `404 Not Found`
- Body: `product not found`

**Example**

```bash
curl -X DELETE http://localhost:8080/products/sku-1001
```

## How To Run The Inventory API

From `backend/go`:

```bash
make setup
make build-and-run
```

If you want MongoDB available locally as well:

```bash
docker compose up -d --env-file .env.local
```

## Important Implementation Caveats

These behaviors are important if you are integrating with the current API:

- Products are stored only in memory via `MapProductRepository`. Restarting the server clears all inventory data.
- MongoDB is configured in the project but the inventory API does not currently persist product data to MongoDB.
- `POST /products` does not enforce unique IDs. Creating a product with an existing `id` silently overwrites the previous map entry.
- Creating a product with an empty `id` is currently possible because there is no explicit ID validation.
- `GET /products` returns items in map iteration order, so response ordering is not stable.
- Error responses are plain text, not structured JSON.
- The sample environment file and Docker config reference `APP_PORT=8000`, but the Go server code currently listens on `8080`.

## Relevant Source Locations

- `backend/go/cmd/app/main.go`: application entrypoint and route registration
- `backend/go/pkg/products/handler`: HTTP handlers and route dispatch
- `backend/go/pkg/products/controller`: validation and business logic
- `backend/go/pkg/products/repository`: in-memory repository implementation
- `backend/go/pkg/products/entity/product.go`: product schema

## Week 3 Progress (MongoDB Integration)

Week 3 focused on moving the inventory service from in-memory storage toward MongoDB-backed persistence while preserving the existing API flow.

- Added MongoDB integration using the official Go driver (`go.mongodb.org/mongo-driver`).
- Introduced `MongoProductRepository` so the repository layer can work with MongoDB collections.
- Added connection health checks with `mongo.Connect` and `client.Ping`.
- Introduced context propagation (`handler -> controller -> repository`) for DB calls.
- Added timeout guards in repository operations using `context.WithTimeout`.
- Improved CRUD reliability with proper Mongo result checks (`MatchedCount`, `DeletedCount`) and corresponding API responses.
- Validated endpoints through Postman and curl, and verified stored data in MongoDB Compass.

## Week 4 Progress

- Added ProductCategory domain model, repository layer (`MongoCategoryRepository`), and service layer (`ProductCategoryService`).
- Implemented full Category CRUD APIs under `/categories` endpoint matching Product handler patterns.
- Established Product ↔ Category relationship using `category_id` field in Product entity.
- Added validation in ProductService to ensure products reference only valid categories via `GetByID` check.
- Implemented category filtering for products using `GET /products?category_id=...` query parameter.
- Refactored Product entity to remove redundant `category` field; use only `category_id` for database consistency.
- Enforced `brand` as a required field in ProductService CreateProduct validation.
- Implemented lazy schema evolution: existing products without brand remain valid on reads, but new/updated products require brand field.

### TODOs

- Move ID generation to backend (use UUID or Mongo ObjectID instead of client-provided IDs).
- Implement batching + controlled concurrency (WaitGroups with batching/worker pool) for bulk product creation.

## Week 6 - Frontend Basics

- Built a basic frontend using plain HTML, CSS, and vanilla JavaScript.
- Created a product display tile and styled it with CSS.
- Integrated frontend with backend APIs using `fetch`.
- Used browser developer tools (DOM inspector, console, and network tab) for debugging.
- Dynamically updated the UI using DOM manipulation based on API responses.
- Rendered multiple products by generating HTML elements dynamically.
- Added simple CSS animations to improve user experience.