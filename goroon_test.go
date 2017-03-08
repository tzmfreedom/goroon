package goroon

import (
	"testing"
	"time"

	"gopkg.in/h2non/gock.v1"
)

func TestScheduleGetEventsByTarget(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/schedule/api").
		Reply(200).File("./test/fixtures/schedule/get_events_by_target.xml")

	client := NewClient("username", "password", "https://garoon.com")
	tm := time.Now()
	req := Parameters{
		Start: XmlDateTime{tm},
		End:   XmlDateTime{tm},
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

func TestScheduleGetEvents(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/schedule/api").
		Reply(200).File("./test/fixtures/schedule/get_events.xml")

	client := NewClient("username", "password", "https://garoon.com")
	tm := time.Now()
	req := Parameters{
		Start: XmlDateTime{tm},
		End:   XmlDateTime{tm},
	}
	res, err := client.ScheduleGetEvents(&req)
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

func TestBaseGetUserByLoginName(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/base/api").
		Reply(200).File("./test/fixtures/base/get_user_by_login_name.xml")

	client := NewClient("username", "password", "https://garoon.com")
	req := Parameters{
		LoginName: []string{"hoge"},
	}
	res, err := client.BaseGetUserByLoginName(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	if len(res.User) != 2 {
		t.Fatalf("expect %v, get %v", 2, len(res.User))
	}
	adm := res.User[0]
	if adm.Key != 1 {
		t.Fatalf("expect %v, get %v", 1, adm.Phone)
	}
	if adm.Version != 1245376338 {
		t.Fatalf("expect %v, get %v", 1245376338, adm.LoginName)
	}
	if adm.Name != "Administrator" {
		t.Fatalf("expect %v, get %v", "Administrator", adm.Name)
	}
	if adm.Status != 0 {
		t.Fatalf("expect %v, get %v", 0, adm.Phone)
	}
	u1 := res.User[1]
	if u1.Key != 2 {
		t.Fatalf("expect %v, get %v", 2, u1.Phone)
	}
	if u1.Version != 1245919830 {
		t.Fatalf("expect %v, get %v", 1245919830, u1.LoginName)
	}
	if u1.Name != "u1" {
		t.Fatalf("expect %v, get %v", "u1", u1.Name)
	}
	if u1.Status != 0 {
		t.Fatalf("expect %v, get %v", 0, u1.Phone)
	}
	if u1.Phone != "9180xxxxxx" {
		t.Fatalf("expect %v, get %v", "9180xxxxxx", u1.Phone)
	}
	if u1.Description != "user1 is ..." {
		t.Fatalf("expect %v, get %v", "user1 is ...", u1.Description)
	}
	if u1.Title != "test test" {
		t.Fatalf("expect %v, get %v", "test test", u1.Title)
	}
}

func TestBulletinGetFollows(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/bulletin/api").
		Reply(200).File("./test/fixtures/bulletin/get_follows.xml")

	client := NewClient("username", "password", "https://garoon.com")
	req := Parameters{
		TopicId: 123,
		Offset:  0,
		Limit:   20,
	}
	res, err := client.BulletinGetFollows(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	if len(res.Follow) != 4 {
		t.Fatalf("expect %v, get %v", 4, len(res.Follow))
	}
}
