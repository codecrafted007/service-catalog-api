# Makefile

BINARY_NAME=service-catalog-api
CMD_PATH=./cmd/api
TEST_TARGET=./...

.PHONY: all build test clean run

all: clean build test run

build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH) --db-driver=sqlite3 --db-dsn=services.db --port=8080

test:
	go test $(TEST_TARGET) -v

clean:
	@echo "Removing: $(BINARY_NAME), services.db"
	rm -f $(BINARY_NAME)
	rm -f ./services.db
