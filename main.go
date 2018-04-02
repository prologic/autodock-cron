package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types/events"
	"github.com/robfig/cron"

	"github.com/prologic/autodock/plugin"
)

// Key ...
const Key = "autodock.cron.schedule"

// CronPlugin ...
var CronPlugin = &plugin.Plugin{
	Name:    "CronPlugin",
	Version: "0.0.2",
	Description: `CronPlugin is a cron-like plugin for autodock which watches
	for contaienr and service startup events and reschedules those
	contaienrs and services according to their configured schedule.`,
	Run: func(ctx plugin.Context) error {
		c := cron.New()
		c.Start()

		ctx.On("container", func(id uint64, payload []byte, created time.Time) error {
			var evt *events.Message
			err := json.Unmarshal(payload, &evt)
			if err != nil {
				log.Errorf("error decoding payload: %s", err)
				return err
			}

			if evt.Action != "create" {
				return nil
			}

			cid := evt.ID
			scid := cid[len(cid)-10:]

			log.Infof("container %s creaed: %#v", scid, evt)

			schedule, ok := evt.Actor.Attributes[Key]
			if !ok || schedule == "" {
				log.Warnf("ignoring container %s with no valid label", scid)
				return nil
			}

			c.AddFunc(schedule, func() {
				err := ctx.StartContainer(cid)
				if err != nil {
					log.Errorf("error starting container %s: %s", scid, err)
				}
			})

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
