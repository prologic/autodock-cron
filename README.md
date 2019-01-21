# autodock-cron

[![Build Status](https://cloud.drone.io/api/badges/prologic/autodock/status.svg)](https://cloud.drone.io/prologic/autodock)

Cron plugin for autodock. autodock-cron is a cron-like plugin for autodock
which watches for `container` and `service` startup events and reschedules
those containers and services according to their configured schedule. The
schedule is configured by container or service labels of the form:

:bulb: See [autodock](https://github.com/prologic/autodock) for more info.

## Building

From source:
```#!bash
$ go build .
```

Using Docker:
```#!bash
$ docker build -t autodock-cron .
```

## Usage

From source:
```#!bash
$ ./autodock-cron -h <autodock_host>
```

Using Docker:
```#!bash
$ docker run -d autodock-cron -H <autodock_host>
```

`autodock-cron` then looks for containers started with a label of
`autodock.cron=<schedule>` where schedule is a valid Cron-like
expression of the form:

- `<seconds> <minutes> <hour> <dom> <month> <dow>`
- `@yearly` (or `@annually`)
- `@monthly`
- `@weekly`
- `@daily` (or `@midnight`)
- `@hourly`
- `@every <duration>`

where `<duration>` is a string accepted by
[time.ParseDuration](http://golang.org/pkg/time/#ParseDuration) for example
`@every 5m` or `@every 20m30s`.

The following is a sample `docker-compose.yml` snippet:
```#!yaml
    deploy:
      labels:
        - "autodock.cron=@every 5m"
```

## License

autodock-cron is MIT licensed.
