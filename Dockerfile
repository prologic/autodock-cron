# Build
FROM golang:alpine AS build

RUN apk add --no-cache -U git make build-base

WORKDIR /src/autodock-cron
COPY . /src/autodock-cron
RUN make build install

# Runtime
FROM alpine:latest

COPY --from=build /go/bin/autodock-cron /autodock-cron

ENTRYPOINT ["/autodock-cron"]
CMD []
