APP=goweb
.PHONY: help
help: Makefile ## Show this help
	@echo
	@echo "Choose a command run in "$(APP)":"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | sed -e 's/\(\:.*\#\#\)/\:\ /' | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: build
build: ## Build an application
	@echo "Building ${APP} ..."
	go build -o build/${APP} cmd/main.go

run: ## Run an application
	@echo "Starting ${APP} ..."
	go run -race cmd/main.go

clean: ## Clean a garbage
	@echo "Cleaning"
	go clean

docker-build: ## Build docker image from Dockerfile
	@echo "Building ${APP} container image..."
	docker-compose -f ./deployments/docker-compose.yml build

docker-run:  ## Run docker containers
	@echo "Run docker containers for ${APP} ..."
	docker-compose -f ./deployments/docker-compose.yml up
