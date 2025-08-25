gen:
	oapi-codegen -package=gen -generate "types,fiber-server" api/scheme.yaml > internal/gen/gen.go

run.seed:
	go run ./cmd/app -seed

run:
	go run ./cmd/app