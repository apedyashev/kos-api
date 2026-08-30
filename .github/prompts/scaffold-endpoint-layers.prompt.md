---
description: "Scaffold empty API/Service/Repository layer files for a new kos-api endpoint, following the Layered Architecture Go conventions, without implementing any logic"
name: "Scaffold Endpoint Layers"
argument-hint: "domain name (new or existing), e.g. tag"
agent: "Layered Architecture Go"
---
Scaffold the file structure for a new endpoint in `kos-api`, strictly following the Layered Architecture Go conventions (API -> Service -> Repository).

**Before creating anything, ask the user for:**
1. The domain name (e.g. `tag`, `space`) — reuse an existing domain under `kos-api/domains/` if it already exists, otherwise this is a new domain.
2. The HTTP method and route path for the new endpoint (e.g. `GET /api/tags/:id`).
3. The name of the use case / handler method (e.g. `GetTag`, `ListTags`, `DeleteTag`).

**Do not implement any business logic, SQL, validation, or mapping bodies.** Only create/update the files and signatures needed so the user can fill in the logic themselves. File names follow the `node` domain's existing conventions (singular domain name, no generic suffixes like `http_handlers`):

- If the domain is new, create the full skeleton:
  - `domains/<domain>/<domain>.go` — module bootstrap with `New(...)` and `RegisterRoutes(e *echo.Echo)` wiring the repository -> service -> handler, with a `TODO` where the new route registration goes.
  - `domains/<domain>/api/dto/<domain>s.go` — empty request/response DTO struct stubs for this endpoint only (no fields unless the user specifies them).
  - `domains/<domain>/api/handlers/<domain>s.go` — an `echo.HandlerFunc` method stub named after the use case, with a `TODO` body and a call signature to the service interface (no implementation).
  - `domains/<domain>/api/handlers/mappers.go` — empty `MapXToY` function stubs for this endpoint's DTO<->Domain conversion (signatures only, `TODO` bodies).
  - `domains/<domain>/service/domain/<domain>.go` — domain model / input type stub for this use case (types only, no invariants implemented).
  - `domains/<domain>/service/<domain>_service.go` — service interface method signature plus a stub implementation with a `TODO` body.
  - `domains/<domain>/service/mappers.go` — empty Domain<->Entity mapper stubs for this use case.
  - `domains/<domain>/persistence/entities/<domain>.go` — entity struct stub (types only, no query logic).
  - `domains/<domain>/persistence/repository/<domain>_repository.go` — repository interface method signature.
  - `domains/<domain>/persistence/repository/postgres/<domain>_repository.go` — a `NewNodeRepository`-style constructor that actually initializes the `pgxpool.Pool` database connection (mirror the working connection setup in `domains/node/persistence/repository/postgres/node_repository.go`, including `Close()`), plus a method stub with a `TODO` body (no SQL). The database connection itself must be real and working, not a stub — only the query logic inside the method stub is left as `TODO`.

- If the domain already exists, only **add** the new method signatures/stubs to the existing files above (do not touch unrelated existing endpoints or overwrite implemented logic), and register the new route in the domain's `RegisterRoutes`.

- Register/confirm the domain module is wired in `kos-api/main.go` if it is new, and that `<domain>.go`'s `New(...)` actually constructs the repository (e.g. via its `NewXRepository(databaseURL)` constructor) so the database connection is initialized when the domain module boots — do not leave repository construction as a `TODO`.

After scaffolding, list every file created or modified and leave all logic as `TODO` comments so the user writes the implementation.
