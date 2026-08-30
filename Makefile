.PHONY: run migrate fmt test test-race build check

run:
	@test -f .env || (echo "Missing kos-api/.env" && exit 1)
	@set -a; . ./.env; set +a; go run .

migrate:
	@test -f .env || (echo "Missing kos-api/.env" && exit 1)
	@command -v psql >/dev/null 2>&1 || (echo "psql is required; install PostgreSQL with: brew install libpq && brew link --force libpq" && exit 1)
	@set -a; . ./.env; set +a; for migration in migrations/*.sql; do psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$$migration"; done

fmt:
	gofmt -w main.go domains internal

test:
	go test ./...

test-race:
	go test -race ./...

build:
	go build ./...

check: fmt test build
