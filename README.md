# autodock-cron

[![Image Layers](https://badge.imagelayers.io/prologic/autodock-cron:latest.svg)](https://imagelayers.io/?images=prologic/autodock-cron:latest)

Cron plugin for autodock. autodock-cron is a cron-like plugin for autodock
which watches for `contaienr` and `service` startup events and reschedules
those contaienrs and services according to their configured schedule. THe
schedule is configured by container or service labels of the form:

```#!yaml
    deploy:
      labels:
        - "autodock.cron.schedule=*/5"
```

autodock-cron is MIT licensed.

> **note**
>
> Please see [autodock](https://github.com/prologic/autodock) for the main project and file issues there.

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
