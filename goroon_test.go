package goroon

import (
	"testing"
	"time"

	"gopkg.in/h2non/gock.v1"
)

func TestGetScheduleByUserId(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/schedule/api").
		Reply(200).File("./test/fixtures/get_schedule_events/request.xml")

	client := NewClient("username", "password", "https://garoon.com")
	tm := time.Now()
	req := Parameters{
		Start: tm,
		End:   tm,
		User: User{
			Id: 1234,
		},
	}
	res, err := client.ScheduleGetEventsByTarget(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	if len(res.ScheduleEvents) != 1 {
		t.Fatalf("expect %v, get %v", 1, len(res.ScheduleEvents))
	}
	ev := res.ScheduleEvents[0]
	if ev.Id != 123 {
		t.Fatalf("expect %v, get %v", 123, ev.Id)
	}
	if ev.Detail != "fugafuga" {
		t.Fatalf("expect %v, get %v", "fugafug", ev.Detail)
	}
	if ev.Description != "hogehoge" {
		t.Fatalf("expect %v, get %v", "hogehoge", ev.Description)
	}
	if len(ev.Members.Member) != 1 {
		t.Fatalf("expect %v, get %v", 1, len(ev.Members.Member))
	}
	member := ev.Members.Member[0]
	if member.User.Id != 1 {
		t.Fatalf("expect %v, get %v", "aa", member.User.Id)
	}
	if member.User.Name != "aaa" {
		t.Fatalf("expect %v, get %v", "bb", member.User.Name)
	}
}
