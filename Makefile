.PHONY: run test gen swagger

run:
	go run ./cmd

test:
	go test ./... -v

gen:
	go run ./tools/gen

swagger:
	swag init -g cmd/main.go -o docs