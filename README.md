# Interneers Lab — Inventory Management

A full-stack inventory management SPA built as a weekly learning project. Go REST API backend with MongoDB persistence, React + TypeScript frontend.

---

## Quick Start

### Backend (Go)

```bash
cd backend/go
make setup          # go mod tidy + copy .env.sample → .env.local
make build-and-run  # compile → bin/app, start server on :8080
```

MongoDB is required for persistence. To spin up a local instance:

```bash
docker compose up -d --env-file .env.local
```

### Frontend (React + TypeScript)

```bash
cd frontend
yarn install
yarn start          # dev server on :3000
```

The frontend expects the Go backend running at `http://localhost:8080`.

---

## Week-by-Week Progress

### Week 1 — Project Setup

- Initialised Go module (`go.mod`), committed to Git
- Adopted hexagonal architecture: `pkg/<module>/entity`, `handler`, `repository` layout
- Built HTTP server using `net/http`, listening on `:8080`
- Added structured logging with [zerolog](https://github.com/rs/zerolog)
- Wired application in `cmd/app/main.go` with explicit dependency injection
- Tested endpoints using Postman and `curl`

---

### Week 2 — In-Memory CRUD APIs

**Product model** introduced (`pkg/products/entity/product.go`):

| Field         | Type     | Required | Notes                    |
|---------------|----------|----------|--------------------------|
| `id`          | `string` | Yes      | Supplied by client       |
| `name`        | `string` | Yes      | Must not be empty        |
| `description` | `string` | No       | Free-form text           |
| `category`    | `string` | No       | Free-form text           |
| `price`       | `float64`| Yes      | Must be > 0              |
| `brand`       | `string` | No       | Free-form text           |
| `quantity`    | `int`    | Yes      | Must be ≥ 0              |

**Implemented:**
- In-memory map storage (`MapProductRepository`)
- Full CRUD: `POST /products`, `GET /products`, `GET /products/{id}`, `PUT /products/{id}`, `DELETE /products/{id}`
- Validation: `name` required, `price > 0`, `quantity ≥ 0`
- Request logging middleware
- Correct HTTP status codes: `200`, `201`, `204`, `400`, `404`, `405`

---

### Week 3 — Layered Architecture + MongoDB

**Refactored into three explicit layers:**

```
Handler  →  Service  →  Repository  →  MongoDB
(HTTP)      (business)  (interface)
```

- Repository layer uses Go interfaces — in-memory and MongoDB implementations are interchangeable
- `context.Context` propagated from HTTP handler through service to DB driver for cancellation
- Per-operation 5-second timeout via `context.WithTimeout`
- MongoDB integration using official Go driver (`go.mongodb.org/mongo-driver`)
- Added `MongoProductRepository` — data now persists across server restarts
- Added CORS middleware (allow-all origins, standard methods/headers)
- Connection health-checked on startup via `client.Ping`

---

### Week 4 — Categories & Relations

**New model** — `ProductCategory` (`pkg/products/entity/category.go`):

| Field         | Type     | Required | Notes             |
|---------------|----------|----------|-------------------|
| `id`          | `string` | Yes      | Generated backend |
| `title`       | `string` | Yes      | Must not be empty |
| `description` | `string` | No       | Free-form text    |

**Implemented:**
- `MongoCategoryRepository` + `ProductCategoryService`
- Full Category CRUD under `/categories`
- `Product.CategoryID` field links products to categories; `category` object embedded on fetch
- Filter: `GET /products?category_id=<id>`
- `brand` added as required field in product validation
- Special endpoint: `DELETE /categories/empty` — remove categories with no products
- Lazy schema evolution: existing products without `brand` remain valid on reads

---

### Week 5 — Testing, Workflow & Bulk Create

**Tests:**
- Unit tests — `pkg/products/service/product_test.go`: ProductService and ProductCategoryService tested with mocked repositories
- Integration tests — `pkg/products/handler/product_integration_test.go`: end-to-end handler tests against a real MongoDB instance

**Workflow:**
- Adopted feature-branch development; all changes merged via PRs, no direct commits to `main`

**Bulk product creation — `POST /products/bulk`:**
- Accepts `multipart/form-data` with a CSV file in the `file` field
- CSV format: `name,price,quantity,brand` (header required, no ID column)
- Up to 10 rows processed concurrently via worker-pool semaphore
- Returns `207 Multi-Status` with `created[]` and `errors[]` arrays
- Malformed rows (< 4 columns, non-numeric price/quantity) silently skipped; business-rule failures reported in `errors[]`

```bash
curl -X POST http://localhost:8080/products/bulk \
  -F "file=@products.csv"
```

---

### Week 6 — Basic Frontend (Vanilla JS)

Location: `frontend/basic/`

- Plain HTML page (`index.html`) with a product-list container
- `script.js` fetches `GET /products` and generates product tile HTML dynamically
- `styles.css` with staggered `fadeInUp` CSS animation on tiles
- Used browser DevTools (DOM inspector, console, network tab) for debugging

---

### Week 7 — React + TypeScript Basics

- Set up React 19 + TypeScript 5.7 project (Create React App, Yarn)
- Created `<ProductCard />` component for product summary display
- Created basic `<ProductList />` to render multiple cards
- Added `<Navbar />` header component
- Wired React Router v7 with an initial `/products` route
- Progressed from dummy data to real `fetch()` calls against the backend
- Added product detail view with basic edit functionality

---

### Week 8 — Full Frontend Integration

**Core:**
- Typed API layer (`src/api/client.ts`, `products.ts`, `categories.ts`) with generic fetch wrapper
- `ProductPage` (`/products/:id`) — view product details and edit inline
- `ProductListPage` (`/products`) — products grouped by category with add-product button
- React Router routes: `/` → redirects to `/products`, `/products`, `/products/:id`
- Loading states (skeleton/spinner) and error states on all data-fetching pages

**Advanced (category pages):**
- `CategoryListPage` (`/categories`) — list all categories
- `CategoryPage` (`/categories/:id`) — category details + grid of products in that category
- Navigation: product card links to its parent category; navbar links both sections
- `LoadingSpinner` and `ErrorMessage` components reused across pages

---

### Week 9 — Clean Architecture & Design System

**Frontend restructure:**

```
src/
  api/          # typed fetch wrappers (products, categories, client)
  types/        # Product, Category, ProductFormData interfaces
  components/
    layout/     # Navbar, PageShell
    product/    # ProductCard, ProductEditForm
    category/   # CategoryCard
    ui/         # Button, Card, ErrorMessage, LoadingSpinner
  pages/        # ProductListPage, ProductPage, CategoryListPage, CategoryPage
  styles/       # variables.css (design tokens), reset.css
```

**OKLCH design system** (`src/styles/variables.css`):
- Full token set: palette, semantic, spacing, radius, typography, shadows — no hard-coded hex/rgb values anywhere
- See [DESIGN.md](./DESIGN.md) for full system documentation

**Backend fix — ID generation:**
- `MongoProductRepository.Create` now generates a MongoDB ObjectID hex string (`primitive.NewObjectID().Hex()`) when no ID is supplied, matching the existing behaviour in `MongoCategoryRepository`
- Clients must **not** include an `id` field on create requests — the backend always generates one

---

## API Reference

Base URL: `http://localhost:8080`

### Products

| Method   | Path                  | Description                              | Body / Notes                              |
|----------|-----------------------|------------------------------------------|-------------------------------------------|
| `GET`    | `/products`           | List all products                        | Optional `?category_id=<id>` filter       |
| `POST`   | `/products`           | Create a product                         | JSON body; do **not** include `id`        |
| `GET`    | `/products/{id}`      | Get product by ID                        |                                           |
| `PUT`    | `/products/{id}`      | Update product by ID                     | JSON body                                 |
| `DELETE` | `/products/{id}`      | Delete product by ID                     | Returns `204 No Content`                  |
| `POST`   | `/products/bulk`      | Bulk-create from CSV                     | `multipart/form-data`, field `file`; returns `207` |

### Categories

| Method   | Path                  | Description                              | Body / Notes                              |
|----------|-----------------------|------------------------------------------|-------------------------------------------|
| `GET`    | `/categories`         | List all categories                      |                                           |
| `POST`   | `/categories`         | Create a category                        | JSON body; do **not** include `id`        |
| `GET`    | `/categories/{id}`    | Get category by ID                       |                                           |
| `PUT`    | `/categories/{id}`    | Update category by ID                    | JSON body                                 |
| `DELETE` | `/categories/{id}`    | Delete category by ID                    | Returns `204 No Content`                  |
| `DELETE` | `/categories/empty`   | Delete all categories with no products   |                                           |

### Status Codes

| Code  | Meaning                         |
|-------|---------------------------------|
| `200` | OK                              |
| `201` | Created                         |
| `204` | No Content (DELETE success)     |
| `207` | Multi-Status (bulk operations)  |
| `400` | Bad Request (validation failed) |
| `404` | Not Found                       |
| `405` | Method Not Allowed              |

---

## Data Models

### Product

```json
{
  "id":          "68ab12cd...",
  "name":        "Wireless Mouse",
  "description": "Ergonomic mouse with USB receiver",
  "category_id": "68ab12ef...",
  "category":    { "id": "...", "title": "Accessories", "description": "" },
  "price":       24.99,
  "brand":       "Logitech",
  "quantity":    12
}
```

**Validation:** `name` required · `price > 0` · `quantity ≥ 0` · `brand` required

### ProductCategory

```json
{
  "id":          "68ab12ef...",
  "title":       "Accessories",
  "description": "Peripherals and accessories"
}
```

**Validation:** `title` required

### Bulk Create CSV

```
name,price,quantity,brand
Milk,50,10,Amul
Phone,5000,2,Samsung
```

### Bulk Create Response (207)

```json
{
  "created": [
    { "id": "68ab12...", "name": "Milk", "price": 50, "quantity": 10, "brand": "Amul" }
  ],
  "errors": [
    { "index": 1, "reason": "Brand is required" }
  ]
}
```

`index` is the 0-based position in the submitted product list (after malformed-row filtering).

---

## Design System

The frontend uses a custom OKLCH-based token system defined in `src/styles/variables.css`. All color, spacing, radius, typography, and shadow values are CSS custom properties. See [DESIGN.md](./DESIGN.md) for full documentation.

### Palette tokens

| Token                | Description                   |
|----------------------|-------------------------------|
| `--color-sage`       | Primary green                 |
| `--color-sage-deep`  | Darker sage (hover states)    |
| `--color-pale-sage`  | Light sage (section headers)  |
| `--color-parchment`  | Surface background            |
| `--color-canvas`     | Page background               |
| `--color-chalk`      | Borders                       |
| `--color-stone`      | Subtle accents                |
| `--color-ink`        | Primary text / Navbar bg      |
| `--color-ash`        | Muted text                    |
| `--color-clay`       | Error states                  |

### Semantic tokens

| Token                 | Maps to            |
|-----------------------|--------------------|
| `--color-primary`     | `--color-sage`     |
| `--color-surface`     | `--color-parchment`|
| `--color-border`      | `--color-chalk`    |
| `--color-text`        | `--color-ink`      |
| `--color-text-muted`  | `--color-ash`      |
| `--color-error`       | `--color-clay`     |
