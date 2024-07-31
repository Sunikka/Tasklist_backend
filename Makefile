build:
	@go build -o cmd/tasklist_backendGo

run: build
	@./cmd/tasklist_backendGo

