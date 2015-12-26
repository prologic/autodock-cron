autodock-cron
=============

[![Image Layers](https://badge.imagelayers.io/prologic/autodock-cron:latest.svg)](https://imagelayers.io/?images=prologic/autodock-cron:latest)

Cron-like Plugin for autodock.

This plugin runs containers on a regular scheduled as defined by the environment variable `CRON` of the form `m h d mon dow` much like crond on many Linux/UNIX systems.

The container must first be run at least once for `autodock-cron` to pick up the new container and its configuration.

Note that `autodock-cron` makes use of the Docker API and effectively calls `docker start` on your container; a new container is **NOT** created on very run of the schedule.

> **note**
>
> See: [autodock](https://github.com/prologic/autodock)

Basic Usage
-----------

Start the daemon:

    $ docker run -d --name autodock prologic/autodock

Link and start an autodock plugin:

    $ docker run -d --link autodock prologic/autodock-cron

Run a container of your choice and set `CRON=*/1 * * * *` to run every minute:

    $ docker run --name hello -e CRON="*/1 * * * *" busybox sh -c 'echo Hello'

Now autodock-cron will schedule a timer to re-run this container every minute until the container is deleted and removed. After about ~3mins you should get:

    $ docker logs hello
    Hello
    Hello
    Hello
