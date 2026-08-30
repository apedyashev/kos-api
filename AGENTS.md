# AGENTS.md

This document defines engineering standards for contributors and AI agents working in this repository.

## Project Context

- Product: Knowledge Operating System (KOS) — API backend
- Stack: Go 1.25.0, [Echo v4](https://echo.labstack.com) HTTP framework, `pgx/v5` (PostgreSQL, no ORM)
- Deployment: single Go binary on Vercel (listens on `PORT`)
- This is its own git repository, nested inside the parent workspace `knowledge-operating-system-ui` alongside the separate `kos-ui` Angular frontend repo. Treat `kos-api/` as the repo root for all commands.

Primary source of product behavior: `../docs/KOS_MVP_Blueprint_v1.0.md`.

## Mission

Serve as the persistence and business-logic layer for a personal knowledge graph (nodes + links), exposed as a small, predictable REST API.

Priorities:

1. Strict layering (API → Service → Repository)
2. Predictable, explicit mapping at every boundary
3. Small, reviewable, scoped changes
4. Correctness over cleverness in SQL and validation

## Non-Negotiable Rules

1. Keep changes scoped and minimal. Do not refactor unrelated areas.
2. Never let API DTOs leak into `service/` or `persistence/`.
3. Never let persistence entities leak into `api/handlers` or `api/dto`.
4. Business invariants (required fields, validation) live only in `service/domain/*.go` `Validate()` methods — never in handlers or repositories.
5. Repository methods take/return entities only; SQL lives only under `persistence/repository/postgres/`.
6. Always map explicitly at boundaries: `MapXToY` functions in `api/handlers/mappers.go` (DTO↔Domain) and `service/mappers.go` (Domain↔Entity).
7. Escape user input used in SQL `ILIKE`/pattern matching (see `escapeLikePattern` in the node Postgres repository) — never interpolate raw strings into SQL.
8. Add or update tests for all meaningful behavior changes in the `service/` layer using the `fakeRepository` pattern (no real DB in unit tests).

## Architecture

Domain-centric layered architecture. Mirror this exact layout for every domain:

```
main.go                                    # Echo app bootstrap, route registration entrypoint
internal/config/config.go                 # Env var loading (DATABASE_URL, DEBUG_ERRORS)
domains/<domain>/
  <domain>.go                              # Module bootstrap: wires repository -> service -> handler, RegisterRoutes(e *echo.Echo)
  api/dto/<domain>s.go                     # Request/response DTOs
  api/handlers/<domain>s.go                # echo.HandlerFunc methods
  api/handlers/mappers.go                  # DTO <-> Domain mappers
  service/<domain>_service.go              # Service interface + implementation, Repository interface, orchestration
  service/mappers.go                       # Domain <-> Entity mappers
  service/service_test.go                  # Unit tests using a fakeRepository
  service/domain/<domain>.go               # Domain model, *Input types + Validate(), sentinel errors
  persistence/entities/<domain>.go         # Entity (DB row shape)
  persistence/repository/<domain>_repository.go              # Repository interface (entities only)
  persistence/repository/postgres/<domain>_repository.go     # pgx-based Postgres implementation
migrations/                                # Postgres schema, applied manually via `make migrate`
```

Naming: `SomethingDTO` (API), `SomethingDomain` (business model), `SomethingEntity` (persistence).

Currently implemented domain: `node` (`GET/POST /api/nodes`, `GET/PATCH /api/nodes/:id`, `GET /api/nodes?search=<title prefix>`).

## Development Commands

Run all commands inside `kos-api/`. There is no CI workflow in this repo — these Makefile targets are the full validation surface.

- Format: `make fmt` (`gofmt -w main.go domains internal`)
- Test: `make test` (`go test ./...`) — no DB or env vars required
- Test with race detector: `make test-race`
- Build: `make build` (`go build ./...`)
- Full pre-submit check: `make check` (fmt + test + build) — **always run this before finishing a change**
- Additional static check: `go vet ./...` (not in the Makefile, safe and fast to run)

`make run` and `make migrate` require a real `DATABASE_URL` (Supabase Postgres) via a git-ignored `kos-api/.env`, and `make migrate` also requires `psql`. Do not rely on these for validation in a sandboxed/agent environment — use `make check` instead.

## Commenting Standards

1. Comments should state what the code cannot show on its own — one short line, not a restatement of the next line.
2. Avoid multi-paragraph doc comments where one line suffices.
3. Keep comments up to date; outdated comments are treated as defects.

## Error Handling

1. Domain invariants return errors (`ErrInvalidNode`, `ErrNodeNotFound`-style sentinels) via `errors.Join`/`errors.Is`.
2. Services propagate and wrap errors with `fmt.Errorf("...: %w", err)`.
3. Handlers translate sentinel errors to HTTP status codes (400 for invalid input, 404 for not found, 500 otherwise) — never expose internal error details unless `DEBUG_ERRORS` is enabled.
4. Before returning a 500-mapped error from a handler, log it via `internal/logger`'s `logger.Error(message string, err error)` — never a raw `log.Printf` — so failures are always visible in server logs regardless of `DEBUG_ERRORS`.
