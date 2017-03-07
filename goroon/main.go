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
	UserId   string
	Debug    bool
	Start    string
	End      string
	TopicId  int
	Offset   int
	Limit    int
	Type     string
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
					Destination: &c.UserId,
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
				cli.StringFlag{
					Name:        "type",
					Destination: &c.Type,
					Value:       "all",
				},
			},
			Action: func(ctx *cli.Context) error {
				client := goroon.NewClient(c.Username, c.Password, c.Endpoint)
				if c.Debug {
					client.Debugger = os.Stdout
				}
				loc, _ := time.LoadLocation("Asia/Tokyo")
				start, err := time.ParseInLocation("2006-01-02 15:04:05", c.Start, loc)
				if err != nil {
					return err
				}
				end, err := time.ParseInLocation("2006-01-02 15:04:05", c.End, loc)
				if err != nil {
					return err
				}

				var returns *goroon.Returns
				if c.UserId != "" {
					res, err := client.BaseGetUserByLoginName(&goroon.Parameters{
						LoginName: []string{c.UserId},
					})
					if err != nil {
						return err
					}
					returns, err = client.ScheduleGetEventsByTarget(&goroon.Parameters{
						Start: start,
						End:   end,
						User: goroon.User{
							Id: res.UserId,
						},
					})
					if err != nil {
						return err
					}
				} else {
					returns, err = client.ScheduleGetEvents(&goroon.Parameters{
						Start: start,
						End:   end,
					})
					if err != nil {
						return err
					}
				}

				for _, event := range returns.ScheduleEvents {
					if c.Type != "all" && event.EventType != c.Type {
						continue
					}
					fmt.Println(strings.Join([]string{
						fmt.Sprint(event.Id),
						members2str(&event.Members),
						event.EventType,
						strings.Replace(event.Detail, "\n", "", -1),
						strings.Replace(event.Description, "\n", "", -1),
						startStr(&event),
						endStr(&event),
					}, "\t"))
				}
				return nil
			},
		},
		{
			Name:    "bulletin",
			Aliases: []string{"b"},
			Usage:   "get bulletin",
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
				cli.IntFlag{
					Name:        "topic_id",
					Destination: &c.TopicId,
				},
				cli.BoolFlag{
					Name:        "debug, d",
					Destination: &c.Debug,
				},
				cli.IntFlag{
					Name:        "offset, o",
					Destination: &c.Offset,
				},
				cli.IntFlag{
					Name:        "limit, l",
					Destination: &c.Limit,
				},
			},
			Action: func(ctx *cli.Context) error {
				client := goroon.NewClient(c.Username, c.Password, c.Endpoint)
				res, err := client.BulletinGetFollows(&goroon.Parameters{
					TopicId: c.TopicId,
					Offset:  c.Offset,
					Limit:   c.Limit,
				})
				if err != nil {
					return err
				}

				for _, follow := range res.Follow {
					fmt.Println(strings.Join([]string{
						fmt.Sprint(follow.Number),
						follow.Creator.Name,
						follow.Text,
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

func members2str(members *goroon.Members) string {
	ret := []string{}
	for _, m := range members.Member {
		ret = append(ret, m.User.Name)
	}
	return strings.Join(ret, ":")
}
