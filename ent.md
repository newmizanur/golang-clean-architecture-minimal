# Ent ORM — Developer Guide

This project uses [entgo.io/ent](https://entgo.io) as the ORM and [Goose](https://github.com/pressly/goose) for production database migrations. They are **kept in sync manually** — ent is the Go query layer, goose is the SQL migration layer.

## How It Works

| Tool | Purpose |
|------|---------|
| **Goose** | Version-controlled SQL migrations (source of truth for DB schema) |
| **Ent** | Type-safe Go query builder (generated from hand-written schemas) |

Ent does **not** generate migrations and does **not** read from the database. You write both the goose SQL and the ent schema to match each other.

---

## Adding a New Feature (e.g. `Product`)

### 1. Write the Goose migration

```bash
task goose:create -- create_table_products
```

Edit the generated file in `db/migrations/`:

```sql
-- +goose Up
CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    price       INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE products;
```

### 2. Create the ent schema

```bash
task ent:new -- Product
```

This creates `ent/schema/product.go`. Define fields and edges to **match your migration**:

```go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
)

type Product struct {
    ent.Schema
}

func (Product) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").MaxLen(255).NotEmpty(),
        field.Int("price"),
        field.Time("created_at").Default(time.Now),
        field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
    }
}

func (Product) Edges() []ent.Edge {
    return nil
}
```

### 3. Generate ent code

```bash
task ent:generate
```

This produces `ent/product.go`, `ent/product_create.go`, `ent/product_query.go`, etc.

### 4. Write the application code

Create the following files (use existing entities as reference):

| File | Purpose |
|------|---------|
| `internal/dto/product_model.go` | Request/response structs + validation tags |
| `internal/dto/converter/product_converter.go` | `*ent.Product` → `*dto.ProductResponse` |
| `internal/repository/product_repository.go` | ent query builders (CRUD) |
| `internal/usecase/product_usecase.go` | Business logic |
| `internal/delivery/http/product_controller.go` | HTTP handler |

Add `ProductRepositoryPort` to `internal/usecase/interfaces.go`.

Wire up in `internal/config/app.go` (repository → usecase → controller).

Register routes in `internal/delivery/http/route/route.go`.

### 5. Apply migration and run

```bash
task goose:up
task run
```

---

## Useful Commands

| Command | Description |
|---------|-------------|
| `task ent:new -- Name` | Create a new ent schema |
| `task ent:generate` | Regenerate ent code from all schemas |
| `task goose:create -- name` | Create a new SQL migration file |
| `task goose:up` | Apply pending migrations |
| `task goose:down` | Roll back last migration |
| `task tools:install` | Install goose and ent CLI tools |

---

## Ent Schema Tips

- **String IDs (UUID)**: use `field.String("id").MaxLen(36).NotEmpty()` and set it explicitly on create
- **Auto-increment ID**: omit the `id` field — ent generates it as `int` (maps to `BIGSERIAL`)
- **Nullable fields**: use `.Optional().Nillable()` → generates `*string` in Go
- **Foreign key edge**: declare the FK field (e.g. `user_id`) then reference it in the edge with `.Field("user_id")`
- **Relationships**: use `edge.To` on the owner side and `edge.From` + `.Ref()` on the back-reference side

## File Locations

```
ent/
  schema/         ← hand-written schemas (edit these)
    user.go
    contact.go
    address.go
    item.go
  generate.go     ← go:generate directive
  *.go            ← generated code (do not edit)

db/migrations/    ← goose SQL files (hand-written)
```
