build:
	@go mod tidy
	@go build
run: build
	./theHeirophant
