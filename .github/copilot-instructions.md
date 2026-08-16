# kos-api — Copilot Cloud Agent Instructions

Trust these instructions first. Only search the codebase if information here is missing or proven wrong by an actual command failure.

## Repository Summary

`kos-api` is a minimal Go HTTP API backing the "Knowledge Operating System" personal knowledge-graph app (see `../docs/KOS_MVP_Blueprint_v1.0.md`). It exposes CRUD endpoints for "nodes" (knowledge entries), backed by PostgreSQL (Supabase), and is deployed as a single Go binary on Vercel. This is its own git repository (nested inside the parent workspace `knowledge-operating-system-ui`, alongside a separate `kos-ui` Angular frontend repo) — treat `kos-api/` as the repo root for all git and build operations.

- **Language/runtime:** Go 1.25.0 (see `go.mod`; `go.mod` also pins `go 1.25.0` as the minimum toolchain — do not downgrade it).
- **Size:** Very small (~10 Go files, one domain implemented: `node`).
- **Web framework:** [Echo v4](https://echo.labstack.com) (`github.com/labstack/echo/v4`).
- **Database driver:** `github.com/jackc/pgx/v5` (raw SQL, no ORM).
- **No other languages/build systems** in this repo (the Angular UI lives in the separate `kos-ui` repo, not here).

## Architecture (strict layering — always preserve this)

This repo follows a domain-centric layered architecture: **API → Service → Repository**. Never let API DTOs leak into the service/repository, and never let persistence entities leak into API handlers. Always map explicitly at boundaries.

```
main.go                                    # Echo app bootstrap, route registration entrypoint
internal/config/config.go                 # Env var loading (DATABASE_URL, DEBUG_ERRORS)
domains/node/
  node.go                                  # Module bootstrap: wires repository -> service -> handler, RegisterRoutes(e *echo.Echo)
  api/dto/nodes.go                         # Request/response DTOs (CreateNodeRequestDTO, UpdateNodeRequestDTO, NodeResponseDTO)
  api/handlers/nodes.go                    # echo.HandlerFunc methods: ListNodes, CreateNode, GetNode, UpdateNode
  api/handlers/mappers.go                  # DTO <-> Domain mappers (MapCreateRequestToDomain, MapDomainToResponse, ...)
  service/node_service.go                  # Service interface + implementation, business orchestration, Repository interface
  service/mappers.go                       # Domain <-> Entity mappers (MapEntityToDomain, MapDomainToEntity)
  service/service_test.go                  # Unit tests using a fakeRepository (no real DB needed)
  service/domain/node.go                   # Domain model (NodeDomain), CreateNodeInput/UpdateNodeInput + Validate(), ErrInvalidNode, ErrNodeNotFound
  persistence/entities/node.go             # NodeEntity (DB row shape)
  persistence/repository/node_repository.go        # Repository interface (entities only, no DTOs)
  persistence/repository/postgres/node_repository.go  # pgx-based Postgres implementation (SQL lives only here)
migrations/001_create_nodes.sql            # Postgres schema, run manually via `make migrate` (not automatic)
```

**Naming conventions:** `SomethingDTO` (API), `SomethingDomain` (business model), `SomethingEntity` (persistence). `MapXToY` for explicit boundary mappers. When adding a new domain, mirror the exact `domains/node/` file layout above.

**Routing:** Routes are registered via `RegisterRoutes(e *echo.Echo)` on each domain module, called from `main.go`. Current routes: `GET/POST /api/nodes`, `GET/PATCH /api/nodes/:id`, plus `GET /api/nodes?search=<prefix>` for a case-insensitive title-prefix search (implemented with `ILIKE` + a manual `escapeLikePattern` helper in the Postgres repository — always escape `%`, `_`, `\` when building `ILIKE` patterns from user input).

## Build, Test, and Validation — always run from the `kos-api/` directory

There is **no CI workflow** in this repo (no `.github/workflows/`) — the `Makefile` targets below are the full validation surface. Always validate changes with `make check` (or the equivalent individual commands) before considering a task done.

```bash
cd kos-api            # always run all commands from here, not the parent workspace root
make check             # runs: gofmt -w, go test ./..., go build ./...  — the standard pre-submit check
```

Individual targets (all verified working with no preconditions beyond having Go 1.25+ installed):

- `make fmt` → `gofmt -w main.go domains internal` (formats in place; verified clean/no-op on current code)
- `make test` → `go test ./...` (only `domains/node/service` has tests today; they use an in-memory `fakeRepository`, **no database or env vars required**)
- `make test-race` → `go test -race ./...`
- `make build` → `go build ./...`
- `go vet ./...` is not wired into the Makefile but is fast and safe to run in addition — verified clean.

`make run` and `make migrate` require a real database and **will fail without one; do not attempt them as part of a normal code change/validation loop**:

- `make run` requires a `kos-api/.env` file with `DATABASE_URL` set (Supabase pooled Postgres connection string) — it will hard-fail with `Missing kos-api/.env` if absent. `.env` is git-ignored and not present in a fresh checkout.
- `make migrate` additionally requires the `psql` client (`brew install libpq && brew link --force libpq` on macOS if missing) and a reachable database.
- Since cloud agent sandboxes have no database credentials, **do not run `make run` or `make migrate`**; rely on `make check` (build + unit tests) for validation instead.

There is no linter config (no `.golangci.yml` or similar) — `gofmt` and `go vet` are the only static checks in use.

## Making Changes

- Always keep new/edited code `gofmt`-formatted (`make fmt` or `gofmt -w <files>`) before finishing — formatting is otherwise not auto-enforced by CI.
- When adding a new use case to an existing domain or a new domain, follow the layered file layout above exactly; add unit tests under `service/` following the `fakeRepository` pattern in `domains/node/service/service_test.go` (no real DB dependency needed for service-layer tests).
- Business validation (e.g., required fields) belongs in `service/domain/*.go` `Validate()` methods, not in handlers or repositories.
- Repository methods must only take/return entities — never DTOs or domain types directly in SQL query signatures beyond the entity struct.
- After any change, run `make check` from `kos-api/` and confirm it exits 0 before finishing.
