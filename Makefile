.PHONY: dep
dep: ## Get the dependencies
	@go get -v -d ./...

.PHONY: upgrade-dep
upgrade-dep: ## Upgrade dependencies
	@go get -u -t -v ./...

.PHONY: build
build: dep ## Build the binary file
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o handler ./cmd/handler
	rm -f handler.zip
	zip handler.zip  ./handler

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
