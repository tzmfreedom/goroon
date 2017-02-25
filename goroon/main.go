package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tzmfreedom/goroon"
	"github.com/urfave/cli"
)

type config struct {
	Username string
	Password string
	Endpoint string
	Userid   string
	Debug    bool
}

func main() {
	c := &config{}
	app := cli.NewApp()
	app.Name = "goroon"
	app.Usage = "garoon utility"
	app.Commands = []cli.Command{
		{
			Name:    "schedule",
			Aliases: []string{"s"},
			Usage:   "get today's your schedule",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "username, u",
					Destination: &c.Username,
					EnvVar:      "GAROON_USERNAME",
				},
				cli.StringFlag{
					Name:        "password, p",
					Destination: &c.Password,
					EnvVar:      "GAROON_PASSWORD",
				},
				cli.StringFlag{
					Name:        "endpoint, e",
					Destination: &c.Endpoint,
					EnvVar:      "GAROON_ENDPOINT",
				},
				cli.StringFlag{
					Name:        "userid, i",
					Destination: &c.Userid,
				},
				cli.BoolFlag{
					Name:        "debug, d",
					Destination: &c.Debug,
				},
			},
			Action: func(ctx *cli.Context) error {
				now := time.Now()
				client := goroon.NewGaroonClient(c.Username, c.Password, c.Endpoint, c.Debug, os.Stdout)
				res := &goroon.Envelope{}
				err := client.ScheduleGetEventsByTarget(
					c.Userid,
					time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
					time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local),
					res,
				)
				if err != nil {
					return err
				}

				for _, event := range res.ResponseBody.ScheduleEvents {
					fmt.Println(strings.Join([]string{
						fmt.Sprint(event.Id),
						fmt.Sprint(event.Members),
						event.EventType,
						strings.Replace(event.Detail, "\n", "", -1),
						strings.Replace(event.Description, "\n", "", -1),
						event.GetStartStr(),
						event.GetEndStr(),
					}, "\t"))
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
