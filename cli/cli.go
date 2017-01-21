package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"github.com/tzmfreedom/goroon"
)

type Cli struct {
	Config       *Config
	Notifier     *Notifier
	Dbclient     *DBClient
	Garoonclient *goroon.GaroonClient
}

type Config struct {
	NotificationTypes *NotificationTypes `yaml:"notification_types"`
	Username          string
	Password          string
	Endpoint          string
	Userid            string
	Configfile        string
}

type NotificationTypes struct {
	Desktop bool `yaml:"desktop"`
}

func NewCli() *Cli {
	return &Cli{
		Notifier: NewNotifier(),
		Config: &Config{
			NotificationTypes: &NotificationTypes{
				Desktop: true,
			},
		},
	}
}

func (c *Cli) Run(args []string) error {
	var err error
	app := cli.NewApp()
	app.Name = "goroon"
	app.Usage = "garoon schedule notifier"
	app.Commands = []cli.Command{
		{
			Name:    "db:create",
			Aliases: []string{"c"},
			Usage:   "create database",
			Action: func(ctx *cli.Context) error {
				dbClient, _ := NewDBClient("./data.db")
				dbClient.CreateTable()
				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start check schedule and notifiy you",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "username, u",
					Destination: &c.Config.Username,
					EnvVar:      "GAROON_USERNAME",
				},
				cli.StringFlag{
					Name:        "password, p",
					Destination: &c.Config.Password,
					EnvVar:      "GAROON_PASSWORD",
				},
				cli.StringFlag{
					Name:        "endpoint, e",
					Destination: &c.Config.Endpoint,
					EnvVar:      "GAROON_ENDPOINT",
				},
				cli.StringFlag{
					Name:        "userid",
					Destination: &c.Config.Userid,
				},
				cli.StringFlag{
					Name:        "config, c",
					Destination: &c.Config.Configfile,
				},
			},
			Action: func(ctx *cli.Context) error {
				if c.Config.Configfile != "" {
					err = c.loadYaml(c.Config.Configfile)
					if err != nil {
						return err
					}
				}

				err = c.loop()
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get today's your schedule",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "username, u",
					Destination: &c.Config.Username,
					EnvVar:      "GAROON_USERNAME",
				},
				cli.StringFlag{
					Name:        "password, p",
					Destination: &c.Config.Password,
					EnvVar:      "GAROON_PASSWORD",
				},
				cli.StringFlag{
					Name:        "endpoint, e",
					Destination: &c.Config.Endpoint,
					EnvVar:      "GAROON_ENDPOINT",
				},
				cli.StringFlag{
					Name:        "userid",
					Destination: &c.Config.Userid,
				},
				cli.StringFlag{
					Name:        "config, c",
					Destination: &c.Config.Configfile,
				},
			},
			Action: func(ctx *cli.Context) error {
				now := time.Now()
				c.Garoonclient = goroon.NewGaroonClient(c.Config.Username, c.Config.Password, c.Config.Endpoint)
				response, err := c.Garoonclient.Request(
					c.Config.Userid,
					time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
					time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local),
				)
				if err != nil {
					return err
				}

				for _, event := range response.ResponseBody.ScheduleEvents {
					fmt.Println("%s: %s - %s", event.Detail,
						event.When.Datetime.Start.Format("2006-01-02T15:04:05"),
						event.When.Datetime.End.Format("2006-01-02T15:04:05"))
				}
				return nil
			},
		},
	}
	app.Run(args)
	return err
}

func (c *Cli) loadYaml(filename string) (err error) {
	c.Config = &Config{}
	readBody, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(readBody), &c.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cli) loop() error {
	dbClient, err := NewDBClient("./data.db")
	if err != nil {
		return err
	}
	now := time.Now()
	c.Garoonclient = goroon.NewGaroonClient(c.Config.Username, c.Config.Password, c.Config.Endpoint)

	for {

		response, err := c.Garoonclient.Request(
			c.Config.Userid,
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local),
		)
		if err != nil {
			return err
		}

		events_from_db := []ScheduleEventStorage{}
		dbClient.Db.Select(&events_from_db, fmt.Sprintf("SELECT * FROM schedule_events WHERE start > '%s'", now.Format("2006-01-02")))

		for _, event_from_response := range response.ResponseBody.ScheduleEvents {
			isCreate := true
			isNotify := false
			for _, event_from_db := range events_from_db {
				if event_from_db.Id == event_from_response.Id {
					isCreate = false
					isNotify = event_from_db.IsNotify
					break
				}
			}
			if isNotify == true {
				continue
			}
			if isCreate {
				dbClient.CreateRecord(event_from_response)
			} else {
				dbClient.UpdateRecord(event_from_response, false)
			}

			dt := event_from_response.When.Datetime
			if dt.Start.Add(-10 * time.Minute).Before(time.Now()) {
				if c.Config.NotificationTypes.Desktop {
					c.Notifier.Notify(
						event_from_response.Detail,
						"",
						fmt.Sprintf("%s/schedule/view?event=%d", c.Config.Endpoint, event_from_response.Id),
					)
					dbClient.UpdateRecord(event_from_response, true)
				}
			}
		}

		//処理
		time.Sleep(60000 * time.Millisecond)
		fmt.Println(time.Now().Format("2006-01-02T15:04:05"))
	}
	return nil
}
