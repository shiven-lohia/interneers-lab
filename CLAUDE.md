# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a full-stack inventory management SPA built as a learning/portfolio project ("interneers-lab"). It has:
- A **Go REST API** backend (`backend/go/`) — the primary active backend
- A **React + TypeScript** frontend (`frontend/`)
- A legacy **Python/Django** backend (`backend/python/`) — not actively developed

---

## Commands

### Backend (Go)

All commands run from `backend/go/`.

```bash
make setup          # go mod tidy + copy .env.sample → .env.local
make build-and-run  # compile to bin/app and start server on :8080
make run            # run previously compiled binary

# Tests
go test ./...                                            # all tests
go test ./pkg/products/service -run TestCreateProduct -v # single unit test
go test ./pkg/products/handler -run TestIntegration_CreateAndGetProduct -v  # single integration test
```

Integration tests (handler layer) require a running MongoDB instance. Set connection string in `.env.local`.

### Frontend

All commands run from `frontend/`.

```bash
yarn install              # install dependencies
yarn start                # dev server on :3000
yarn build                # production build

# Unit tests (Jest)
yarn test                                        # watch mode
yarn test src/__tests__/App.test.tsx             # single file
yarn test -- --coverage                          # with coverage

# E2E tests (Playwright)
yarn playwright install                          # one-time browser install
yarn playwright test                             # run all E2E tests
yarn playwright test integration-tests/example.spec.ts  # single spec
yarn playwright test -- --headed                 # headed mode
yarn playwright show-report                      # open last report
```

---

## Architecture

### Backend (Go) — Clean Architecture

```
cmd/app/main.go          # entry point: wires dependencies, starts HTTP server
pkg/
  products/
    entity/              # domain models: Product, ProductCategory
    repository/          # MongoDB data access; also has in-memory impl
    service/             # business logic; depends on repository interfaces
    handler/             # HTTP handlers, route registration
  middleware/            # CORS, request logging
```

Dependency flow: `main.go` → constructs repository → injects into service → injects into handler. The repository layer uses interfaces, so unit tests swap in in-memory implementations.

Context propagates from HTTP handler through service to DB for cancellation/timeout.

### Frontend (React + TypeScript)

```
src/
  App.tsx                        # router setup (React Router v6), Navbar + Routes
  index.tsx                      # root mount
  index.css                      # imports variables.css and reset.css
  styles/
    variables.css                # OKLCH design tokens (palette + semantic)
    reset.css                    # base reset
  types/
    index.ts                     # Product, Category, ProductFormData interfaces
  api/
    client.ts                    # base fetch wrapper
    products.ts                  # product API calls
    categories.ts                # category API calls
  pages/
    ProductListPage.tsx/css      # grouped-by-category product grid
    ProductPage.tsx/css          # product detail + inline edit form
    CategoryListPage.tsx/css     # category list (ledger-row layout)
    CategoryPage.tsx/css         # category detail + products grid
  components/
    layout/
      Navbar.tsx/css             # top nav (ink bg, sage active indicator)
      PageShell.tsx/css          # max-width content wrapper
    product/
      ProductCard.tsx/css        # card: name/brand + price/stock footer
      ProductEditForm.tsx/css    # inline edit form with category creation
    category/
      CategoryCard.tsx/css       # list-row: title left, description right
    ui/
      Button.tsx/css             # primary (sage) / secondary / danger variants
      Card.tsx/css               # base clickable card with hover lift
      ErrorMessage.tsx/css       # error banner (clay palette)
      LoadingSpinner.tsx/css     # centered spinner
  __tests__/                     # Jest unit tests
integration-tests/               # Playwright E2E specs
```

Routes: `/` → redirects to `/products`, `/products`, `/products/:id`, `/categories`, `/categories/:id`.

The frontend fetches directly from the Go backend at `http://localhost:8080`.

---

## Design System

The frontend uses a custom OKLCH-based design system defined in `src/styles/variables.css`. All colors, spacing, radius, typography sizes, and shadows are CSS custom properties. Do not use hard-coded hex/rgb values — use tokens.

**Palette tokens:** `--color-sage`, `--color-sage-deep`, `--color-pale-sage`, `--color-parchment`, `--color-canvas`, `--color-chalk`, `--color-stone`, `--color-ink`, `--color-ash`, `--color-clay`

**Semantic tokens:** `--color-primary` (sage), `--color-surface` (parchment), `--color-border` (chalk), `--color-text` (ink), `--color-text-muted` (ash), `--color-error` (clay)

**Key design rules:**
- Cards are flat at rest; hover lift via `--shadow-lift` only on interactive cards
- `--color-primary` is sage (green), not dark. Navbar background uses `--color-ink` directly
- Category section headers use pale sage background + uppercase tracked label (13px, 600, 0.04em)
- No `border-left` accent stripes, no gradient text, no glassmorphism
- See `DESIGN.md` for full system documentation

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/products` | List products (supports `?category=` filter) |
| POST | `/products` | Create product |
| GET | `/products/{id}` | Get product |
| PUT | `/products/{id}` | Update product |
| DELETE | `/products/{id}` | Delete product |
| GET | `/categories` | List categories |
| POST | `/categories` | Create category |
| GET | `/categories/{id}` | Get category |
| PUT | `/categories/{id}` | Update category |
| DELETE | `/categories/{id}` | Delete category |
| POST | `/products/bulk` | Bulk-create products from CSV (multipart/form-data, field `file`) |

---

## ID Generation

Both products and categories have their IDs generated by the backend using `primitive.NewObjectID().Hex()` (MongoDB ObjectID as hex string). Clients must **not** supply an `id` field on create — the backend ignores any client-provided ID for products and always generates a fresh one in `MongoProductRepository.Create`. Categories behave the same via `MongoCategoryRepository.Create`.

## Bulk Create Products

`POST /products/bulk` accepts a `multipart/form-data` request with a single field named `file` containing a CSV.

**CSV format** (header row required, 6 columns):

```
name,description,category_id,price,quantity,brand
Milk,Fresh milk,,50,10,Amul
Phone,Android phone,<category_id>,5000,2,Samsung
```

- IDs are generated by the backend; do not include an ID column.
- `description` and `category_id` may be empty strings.
- If `category_id` is provided it must match an existing category; otherwise the row is rejected.
- Rows with fewer than 6 columns or non-numeric price/quantity are silently skipped.
- Up to 10 rows are processed concurrently (worker pool).
- Response is always **207 Multi-Status** with a JSON body:

```json
{
  "created": [ /* successfully inserted Product objects */ ],
  "errors":  [ { "index": 1, "reason": "Brand is required" } ]
}
```

`index` is the 0-based position of the row in the submitted product list (after header and malformed-row filtering).

---

## Known TODOs

- Add create-product and create-category routes/pages (Add product/category buttons exist but navigate to unimplemented routes)
