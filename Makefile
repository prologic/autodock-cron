.PHONY: dev build image test deps clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
APP=autodock-cron
REPO?=prologic/$(APP)
TAG?=latest
BUILD?=-dev

all: dev

dev: build
	@./$(APP) -debug

deps:
	@go get ./...

build: clean deps
	@go build -tags "netgo static_build" -installsuffix netgo .

image:
	@docker build -t $(REPO):$(TAG) .
	@echo "Image created: $(REPO):$(TAG)"

test:
	@go test -v -cover -race $(TEST_ARGS)

clean:
	@rm -rf $(APP)
