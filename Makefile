.PHONY: dev build install image test deps clean

CGO_ENABLED=0

all: dev

dev: build
	@./autodock-cron -d

build:
	@go build \
		-tags "netgo static_build" \
		-installsuffix netgo \
		.

install: build
	@go install

image:
	@docker build -t prologic/autodock-cron .

profile:
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench .

bench:
	@go test -v -bench .

test:
	@go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic .

clean:
	@git clean -f -d -X
