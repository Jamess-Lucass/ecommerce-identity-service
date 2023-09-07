PROJECT_NAME=ecommerce-identity-service

.PHONY: dev
dev:
	dotenv -- go run ./cmd

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	docker run --rm --name ${PROJECT_NAME}-testing-client \
		-w /src \
		-v $(shell pwd):/src golangci/golangci-lint:latest golangci-lint run -v

.PHONY: compose
compose:
ifdef SERVICE
	docker compose up -d $(SERVICE)
else
	docker compose up -d
endif

.PHONY: compose-build
compose-build:
ifdef SERVICE
	docker compose up -d $(SERVICE) --build
else
	docker compose up -d --build
endif
