build:
`go build -o bin/time-capsule-backend cmd/main.go`

run: build
`./bin/time-capsule-backend`

run:
`go run cmd/main.go`

test:
`go test -v ./...`

migration:
`migration create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))`

migrate-up:
`go run cmd/migrate/main.go up`

migrate-down:
`go run cmd/migrate/main.go down`
