package main

import (
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"github.com/tzmfreedom/goroon"
)

type DBClient struct {
	Db *sqlx.DB
}

func NewDBClient(parameter string) (client *DBClient, err error) {
	sqlxdb, err := sqlx.Open("sqlite3", parameter)
	if err != nil {
		return
	}
	client = &DBClient{Db: sqlxdb}
	return
}

func (c *DBClient) Exec(q string) error {
	var _, err = c.Db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func CreateDatabase() {
	os.Create("./data.db")
}

func (c *DBClient) CreateRecord(event goroon.ScheduleEvent) {
	sql := `INSERT INTO schedule_events (id, detail, description, start, end) VALUES (%d, '%s', '%s', '%s', '%s')`
	c.Exec(fmt.Sprintf(sql, event.Id, event.Detail, event.Description, event.When.Datetime.Start.Format("2006-01-02T15:04:05"), event.When.Datetime.End.Format("2006-01-02T15:04:05")))
}

func (c *DBClient) UpdateRecord(event goroon.ScheduleEvent, isNotify bool) {
	sql := `UPDATE schedule_events SET detail='%s', description='%s', start='%s', end='%s', is_notify='%t' WHERE id=%d`
	bind_sql := fmt.Sprintf(sql, event.Detail, event.Description, event.When.Datetime.Start.Format("2006-01-02T15:04:05"), event.When.Datetime.End.Format("2006-01-02T15:04:05"), isNotify, event.Id)
	c.Exec(bind_sql)
}

func (c *DBClient) CreateTable() {
	q := `CREATE TABLE schedule_events (
	id INTEGER PRIMARY KEY
	, description TEXT NOT NULL
	, detail TEXT NOT NULL
	, start TIMESTAMP NOT NULL
	, end TIMESTAMP NOT NULL
	, is_notify BOOL NOT NULL DEFAULT false
	, created_at TIMESTAMP DEFAULT (DATETIME('now','localtime'))
	)`
	c.Exec(q)
}
