BIN_NAME=s3upoader
VERSION?=0.0.1

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all help test build vendor dep test clean

all: help

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

run: 
	@go run -race main.go svc

build: dep 
	@echo "Removing previous build ..."
	@rm -rf out/
	@mkdir -p out/bin
	@echo "Building binary ..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/bin/$(BIN_NAME) .
	@echo "Build Complete!"

dep: ## get dependency
	@go mod tidy
	@go mod vendor

test: ## Run test and create coverage files
	@go test -race $$(go list ./... | grep -v /vendor/) -coverprofile coverage.out 
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func coverage.out

clean: ## Remove previous build
	@rm -rf out/
