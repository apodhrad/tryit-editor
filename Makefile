TARGET = $(CURDIR)/target

clean:
	@rm -rf $(TARGET)

build: clean
	@mkdir -p $(TARGET)
	@go build -o $(TARGET)

test: clean
	@go clean -testcache
	@go test ./...

test-coverage: clean
	@mkdir -p $(TARGET)
	@go test ./... -coverprofile=$(TARGET)/coverage.out
	go tool cover -html=$(TARGET)/coverage.out
