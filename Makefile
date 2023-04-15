APP=sre-solution-cup-2023

swag: ## generate swagger docss
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init --parseDependency

build: swag ## build
	go build \
    -o ${APP}

run: build ## run
	./${APP}

test: ## test
	go test  -coverprofile=coverage.out models/

help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)