.DEFAULT_GOAL := compile

compile: lint
	@go build -o betting cmd/main.go

lint:
	 golangci-lint run ./...

execute: compile
	./betting serve