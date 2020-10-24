GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := ../api-service/cmd/api-service/*.go

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: build
build:
		go build ../api-service/cmd/api-service

.PHONY: clean
clean:
		go clean
		rm ./api-service