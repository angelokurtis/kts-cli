.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

##@ Build

.PHONY: build
build: ## Initialize Terraform configurations
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o out/kts main.go

.PHONY: install
install: build ## Initialize Terraform configurations
	sudo install -o root -g root -m 0755 out/kts /usr/local/bin/kts

.PHONY: format
format: goimports-reviser gofumpt wsl ## Cleans up the code for easier reading and collaboration.
	$(GOIMPORTS_REVISER) -set-alias -use-cache -rm-unused -format ./...
	$(GOFUMPT) -w -extra .
	$(WSL) -fix ./...

.PHONY: lint
lint: golangci-lint ## Analyze and report style, formatting, and syntax issues in the source code.
	$(GOLANGCI_LINT) run ./...

##@ Tool Binaries

GOIMPORTS_REVISER = $(shell pwd)/bin/goimports-reviser
.PHONY: goimports-reviser
goimports-reviser: ## Checks for goimports-reviser installation and downloads it if not found.
	$(call go-get-tool,$(GOIMPORTS_REVISER),github.com/incu6us/goimports-reviser/v3@v3.6.4)

GOFUMPT = $(shell pwd)/bin/gofumpt
.PHONY: gofumpt
gofumpt: ## Checks for gofumpt installation and downloads it if not found.
	$(call go-get-tool,$(GOFUMPT),mvdan.cc/gofumpt@v0.6.0)

WSL = $(shell pwd)/bin/wsl
.PHONY: wsl
wsl: ## Checks for wsl installation and downloads it if not found.
	$(call go-get-tool,$(WSL),github.com/bombsimon/wsl/v4/cmd...@v4.2.1)

WIRE = $(shell pwd)/bin/wire
.PHONY: wire
wire: ## Checks for wire installation and downloads it if not found.
	$(call go-get-tool,$(WIRE),github.com/google/wire/cmd/wire@v0.6.0)

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Checks for golangci-lint installation and downloads it if not found.
	$(call go-get-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
