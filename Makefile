build:
	@go build -o bin/tasklist_backendGo

run: build
	@./bin/tasklist_backendGo

test:
	@go test -v ./...

