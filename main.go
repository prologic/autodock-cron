package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/robfig/cron"

	"github.com/prologic/autodock/plugin"
)

// Key ...
const Key = "autodock.cron"

// CronPlugin ...
var CronPlugin = &plugin.Plugin{
	Name:    "CronPlugin",
	Version: "0.0.4",
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
			scid := cid[:10]

			log.Infof("container %s creaed: %#v", scid, evt)

			schedule, ok := evt.Actor.Attributes[Key]
			if !ok || schedule == "" {
				log.Warnf("ignoring container %s with no valid label", scid)
				return nil
			}

			err = c.AddFunc(schedule, func() {
				err := ctx.Docker().ContainerStart(
					context.Background(),
					cid,
					types.ContainerStartOptions{},
				)
				if err != nil {
					log.Errorf("error starting container %s: %s", scid, err)
				}
			})
			if err != nil {
				log.Errorf(
					"error adding schedule %s for container %s: %s",
					schedule, scid, err,
				)
				return err
			}

			return nil
		})

		args := filters.NewArgs(
			filters.Arg("label", Key),
		)
		containers, err := ctx.Docker().ContainerList(
			context.Background(),
			types.ContainerListOptions{
				All:     true,
				Filters: args,
			},
		)
		if err != nil {
			log.Errorf("error listing containers: %s", err)
		} else {
			for _, container := range containers {
				cid := container.ID
				scid := cid[:10]
				schedule := container.Labels[Key]

				log.Infof(
					"found existing container %s running %s with schedule %s",
					scid, container.Image, schedule,
				)
				err = c.AddFunc(schedule, func() {
					err := ctx.Docker().ContainerStart(
						context.Background(),
						cid,
						types.ContainerStartOptions{},
					)
					if err != nil {
						log.Errorf("error starting container %s: %s", scid, err)
					}
				})
				if err != nil {
					log.Errorf(
						"error adding schedule %s for container %s: %s",
						schedule, scid, err,
					)
				}
			}
		}

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
