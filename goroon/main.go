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
	Start    string
	End      string
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
				cli.StringFlag{
					Name:        "start",
					Destination: &c.Start,
				},
				cli.StringFlag{
					Name:        "end",
					Destination: &c.End,
				},
			},
			Action: func(ctx *cli.Context) error {
				now := time.Now()
				client := goroon.NewClient(c.Username, c.Password, c.Endpoint, c.Debug, os.Stdout)

				start, err := time.Parse("2006-01-02 15:04:05", c.Start)
				if err != nil {
					return err
				}
				end, err := time.Parse("2006-01-02 15:04:05", c.End)
				if err != nil {
					return err
				}
				res := &goroon.ScheduleGetEventsByTargetResponse{}
				req := &goroon.ScheduleGetEventsByTargetRequest{
					Parameters: &goroon.Parameters{
						Start: start,
						End:   end,
						User: &goroon.User{
							Id: c.Userid,
						},
					},
				}

				err := client.ScheduleGetEventsByTarget(req, res)
				if err != nil {
					return err
				}

				for _, event := range res.Returns.ScheduleEvents {
					fmt.Println(strings.Join([]string{
						fmt.Sprint(event.Id),
						fmt.Sprint(event.Members),
						event.EventType,
						strings.Replace(event.Detail, "\n", "", -1),
						strings.Replace(event.Description, "\n", "", -1),
						startStr(event),
						endStr(event),
					}, "\t"))
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func startStr(event *goroon.ScheduleEvent) string {
	if event.EventType == "banner" {
		return fmt.Sprintf("%s00:00:00", event.When.Date.Start.Format("2006-01-02T"))
	}
	return event.When.Datetime.Start.Format("2006-01-02T15:04:05")
}

func endStr(event *goroon.ScheduleEvent) string {
	if event.EventType == "banner" {
		return fmt.Sprintf("%s00:00:00", event.When.Date.Start.Format("2006-01-02T"))
	}
	return event.When.Datetime.End.Format("2006-01-02T15:04:05")
}
