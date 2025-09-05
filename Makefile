build:
	GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o graphql.wasm ./cmd/graphql