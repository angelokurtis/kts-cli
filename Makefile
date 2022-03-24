.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: build
build: ## Initialize Terraform configurations
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o out/kts main.go

.PHONY: install
install: build ## Initialize Terraform configurations
	sudo install -o root -g root -m 0755 out/kts /usr/local/bin/kts
