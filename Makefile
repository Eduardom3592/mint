.PHONY: all build clean test test-race lint vet fmt run cover install

APP_NAME := mint
ifeq ($(OS),Windows_NT)
VERSION := $(shell powershell -NoProfile -Command "if (Test-Path VERSION) { (Get-Content VERSION -Raw).Trim() } else { 'dev' }")
COMMIT := $(shell powershell -NoProfile -Command "$$commit = git rev-parse --short HEAD 2>$$null; if ($$LASTEXITCODE -eq 0 -and $$commit) { $$commit } else { 'unknown' }")
BUILDDATE := $(shell powershell -NoProfile -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')")
SET_CGO := set CGO_ENABLED=1&&
CLEAN_CMD := powershell -NoProfile -Command "$$paths = @('$(APP_NAME)', '$(APP_NAME).exe', 'coverage.out', 'coverage.html', 'dist'); foreach ($$path in $$paths) { if (Test-Path $$path) { Remove-Item -LiteralPath $$path -Recurse -Force } }"
else
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDDATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
SET_CGO := CGO_ENABLED=1
CLEAN_CMD := rm -f $(APP_NAME) $(APP_NAME).exe coverage.out coverage.html && rm -rf dist/
endif
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILDDATE)"

all: lint vet test build

build:
	go build $(LDFLAGS) -o $(APP_NAME) ./cmd/$(APP_NAME)

install:
	go install $(LDFLAGS) ./cmd/$(APP_NAME)

test:
	go test ./... -count=1

test-race:
	$(SET_CGO) go test ./... -count=1 -race

cover:
	go test ./... -count=1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

vet:
	go vet ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

run:
	go run $(LDFLAGS) ./cmd/$(APP_NAME)

clean:
	$(CLEAN_CMD)
