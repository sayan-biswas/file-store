BINARY_NAME=store

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

all: help

## Build:
build: build-server build-client ## Build project and place binaries in bin 

build-server: ## Build only server
	mkdir -p bin/server	
	GOARCH=amd64 GOOS=windows go build -o bin/server/$(BINARY_NAME).exe cmd/server/server.go	
	GOARCH=amd64 GOOS=linux go build -o bin/server/$(BINARY_NAME).sh cmd/server/server.go

build-client: ## Build only client
	mkdir -p bin/client	
	GOARCH=amd64 GOOS=windows go build -o bin/client/$(BINARY_NAME).exe cmd/client/client.go
	GOARCH=amd64 GOOS=linux go build -o  bin/client/$(BINARY_NAME).sh cmd/client/client.go

build-docker: ## Create docker container
	docker build -t store:latest .

## Run:
run: ## Run server
	bin/server/store.exe

run-docker: ## Run as docker container
	docker run --rm --name store -p 8080:8080 -v $(pwd)/database:/database store:latest

## Clean:
clean: ## Clean build environment and bin folder
	go clean
	rm bin/server/*.*
	rm bin/client/*.*
	rmdir bin/server
	rmdir bin/client
	docker rmi store:latest

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

