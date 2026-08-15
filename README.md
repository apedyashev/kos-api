# KOS API

Minimal Go HTTP API for Vercel deployment.

## Local development

```bash
cd kos-api
go run ./api
```

The API is exposed through Vercel-compatible Go serverless handlers in `api/`.

## Deploy to Vercel

1. Import this folder as a new Vercel project.
2. Keep the root directory as `kos-api`.
3. Vercel will detect `go.mod` and the Go runtime automatically.
4. Deploy. The app listens on the `PORT` environment variable as required by the Vercel Go runtime.

## Endpoints

- `GET /api` -> health payload
- `GET /api/nodes` -> demo node list

For a production app, replace the demo payloads with your real persistence layer.
