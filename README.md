# KOS API

Minimal Go HTTP API for Vercel deployment.

## Database setup

Create a Supabase project, then set `DATABASE_URL` to the Supabase pooled
connection string. Use the pooler URL for serverless deployments.

For local development, create a `.env` file or export the variable directly:

```bash
export DATABASE_URL='postgres://postgres.<project-ref>:[password]@aws-0-<region>.pooler.supabase.com:6543/postgres?sslmode=require'
```

Never commit the connection string. The `.gitignore` already excludes `.env`.

Install the PostgreSQL client on macOS if `psql` is not available:

```bash
brew install libpq
brew link --force libpq
```

Run all database migrations from the `kos-api` directory:

```bash
make migrate
```

The target loads `DATABASE_URL` from `.env` and runs each SQL file in
`migrations/` in filename order with errors enabled. The migrations are safe
to run again because they use idempotent schema statements.

## Local development

```bash
cd kos-api
make run
```

The API is exposed through the local HTTP server registered in `main.go`.

Useful commands:

- `make run` -> load `.env` and start the API
- `make migrate` -> apply the Supabase database migration
- `make test` -> run tests
- `make test-race` -> run tests with the race detector
- `make check` -> format, test, and build

## Deploy to Vercel

1. Import this folder as a new Vercel project.
2. Keep the root directory as `kos-api`.
3. Vercel will detect `go.mod` and the Go runtime automatically.
4. Deploy. The app listens on the `PORT` environment variable as required by the Vercel Go runtime.

## Endpoints

- `GET /api` -> health payload
- `GET /api/nodes` -> nodes ordered by most recently updated
- `POST /api/nodes` -> create a node

Request body for `POST /api/nodes`:

```json
{
	"title": "Subscription Lifecycle",
	"content": {"type": "doc", "content": []},
	"isStub": false
}
```
