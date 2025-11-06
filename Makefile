BINARY_NAME := sheldon
BUILD_DIR := bin
ENTRYPOINT := ./cmd/sheldon

.PHONY: build test tidy install shell-alias clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(ENTRYPOINT)

test:
	GOCACHE=$$(mktemp -d) go test ./...

tidy:
	go mod tidy

install: build
	@GOBIN=$${GOBIN:-$$(go env GOBIN)}; \
	if [ -z "$$GOBIN" ]; then \
		GOBIN=$$(go env GOPATH)/bin; \
	fi; \
	mkdir -p "$$GOBIN"; \
	install -m 755 $(BUILD_DIR)/$(BINARY_NAME) "$$GOBIN/$(BINARY_NAME)"; \
	printf "Installed $(BINARY_NAME) to %s\n" "$$GOBIN"

shell-alias: build
	@scripts/add-shell-alias.sh $(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
