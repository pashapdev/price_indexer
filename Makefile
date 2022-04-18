run-example:
	ADDRESS=":8080" \
	go run ./examples/...

lint:
	golangci-lint run ./...

test:
	go test ./...