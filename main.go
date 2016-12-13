package main

import (
	_ "github.com/k0kubun/pp"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type ScheduleEventStorage struct {
	Id          int
	Start       time.Time `db:"start"`
	End         time.Time `db:"end"`
	Detail      string
	Description string
	IsNotify    bool      `db:"is_notify"`
	CreatedAt   time.Time `db:"created_at"`
}

func main() {
	cli := NewCli()
	err := cli.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
