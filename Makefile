build:
	@go build -o bin/api
run: build
	@./bin/api
test:
	@go test -v ./...
seed:
	@go run script/seed.go