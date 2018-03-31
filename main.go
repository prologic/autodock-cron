package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prologic/autodock/plugin"
)

// CronPlugin ...
var CronPlugin = &plugin.Plugin{
	Name:    "CronPlugin",
	Version: "0.0.1",
	Description: `CronPlugin is a cron-like plugin for autodock which watches
	for contaienr and service startup events and reschedules those
	contaienrs and services according to their configured schedule.`,
	Run: func(ctx plugin.Context) error {
		ctx.On("container", func(id uint64, payload []byte, created time.Time) error {
			log.Infof("id=%v payload=%s created=%v", id, payload, created)
			return nil
		})

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Infof("caught %s, shutting down ...", sig)
		return nil
	},
}

func main() {
	log.Fatal(CronPlugin.Execute())
}
