VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT_NAME := $(shell basename "$(PWD)")

PKG := "github.com/suryakencana007/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

STIME := $(shell date +%s)

.PHONY: lint test gomodgen dep build clean kill frontend serve race coverage coverhtml

lint: ## Lint the files
	@echo " >_ Linter Checking..."
	@golangci-lint run ${${PKG}/...}
	@echo " >_ Done Linter Checked"
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

test: ## Run unittests
	@echo " >_ Test Running..."
	@go test -short ${PKG_LIST}
	@echo " >_ Done Tested"
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

gomodgen:
	./gomod.sh;

dep:
	@go mod download

race: ## Run data race detector
	@go test -race ${PKG_LIST}

coverage: ## Generate global code coverage report
	@echo " >_ Coverage Test Running..."
	./coverage.sh;
	@echo " >_ Done Coverage Tested"
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"

coverhtml: ## Generate global code coverage report in HTML
	@echo " >_ Coverage Test Running..."
	./coverage.sh html;
	@echo " >_ Done Coverage Tested"
	@echo "Process took $$(($$(date +%s)-$(STIME))) seconds"
