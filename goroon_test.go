package goroon

import (
	"strings"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func TestScheduleGetEventsByTarget(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/schedule/api").
		Reply(200).File("./test/fixtures/schedule/get_events_by_target.xml")

	client := NewClient("https://garoon.com")
	client.Username = "username"
	client.Password = "password"

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
	assert(t, len(res.ScheduleEvents), 1)
	ev := res.ScheduleEvents[0]
	assert(t, ev.Id, 123)
	assert(t, ev.Detail, "fugafuga")
	assert(t, ev.Description, "hogehoge")
	assert(t, len(ev.Members.Member), 1)
	member := ev.Members.Member[0]
	assert(t, member.User.Id, 1)
	assert(t, member.User.Name, "aaa")
}

func TestScheduleGetEvents(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/schedule/api").
		Reply(200).File("./test/fixtures/schedule/get_events.xml")

	client := NewClient("https://garoon.com")
	client.Username = "username"
	client.Password = "password"

	tm := time.Now()
	req := Parameters{
		Start: XmlDateTime{tm},
		End:   XmlDateTime{tm},
	}
	res, err := client.ScheduleGetEvents(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	assert(t, len(res.ScheduleEvents), 1)
	ev := res.ScheduleEvents[0]
	assert(t, ev.Id, 123)
	assert(t, ev.Detail, "fugafuga")
	assert(t, ev.Description, "hogehoge")
	assert(t, len(ev.Members.Member), 1)
	member := ev.Members.Member[0]
	assert(t, member.User.Id, 1)
	assert(t, member.User.Name, "aaa")
	assert(t, ev.RepeatInfo.Condition.StartDate, XmlDate{time.Date(2016, 11, 22, 0, 0, 0, 0, time.UTC)})
	assert(t, ev.RepeatInfo.Condition.EndDate, XmlDate{time.Date(2017, 4, 1, 0, 0, 0, 0, time.UTC)})
	assert(t, ev.RepeatInfo.Condition.StartTime, "14:00:00")
	assert(t, ev.RepeatInfo.Condition.EndTime, "14:30:00")
	assert(t, ev.RepeatInfo.Condition.Day, 20)
	assert(t, ev.RepeatInfo.Condition.Week, 2)
	assert(t, ev.RepeatInfo.Condition.Type, "week")
}

func TestBaseGetUserByLoginName(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/base/api").
		Reply(200).File("./test/fixtures/base/get_user_by_login_name.xml")

	client := NewClient("https://garoon.com")
	client.Username = "username"
	client.Password = "password"

	req := Parameters{
		LoginName: []string{"hoge"},
	}
	res, err := client.BaseGetUserByLoginName(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	assert(t, len(res.User), 2)
	adm := res.User[0]
	assert(t, adm.Key, 1)
	assert(t, adm.Version, 1245376338)
	assert(t, adm.Name, "Administrator")
	assert(t, adm.Status, 0)
	u1 := res.User[1]
	assert(t, u1.Key, 2)
	assert(t, u1.Version, 1245919830)
	assert(t, u1.Name, "u1")
	assert(t, u1.Status, 0)
	assert(t, u1.Phone, "9180xxxxxx")
	assert(t, u1.Description, "user1 is ...")
	assert(t, u1.Title, "test test")
}

func TestBulletinGetFollows(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/cbpapi/bulletin/api").
		Reply(200).File("./test/fixtures/bulletin/get_follows.xml")

	client := NewClient("https://garoon.com")
	client.Username = "username"
	client.Password = "password"

	req := Parameters{
		TopicId: 123,
		Offset:  0,
		Limit:   20,
	}
	res, err := client.BulletinGetFollows(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	assert(t, len(res.Follow), 4)
	assert(t, res.Follow[0].Creator.Name, "huy")
}

func TestUtilLogin(t *testing.T) {
	defer gock.Off()

	gock.New("https://garoon.com").
		Post("/util_api/util/api").
		Reply(200).File("./test/fixtures/util_api/util_login.xml")

	client := NewClient("https://garoon.com")
	client.Username = "username"
	client.Password = "password"

	req := Parameters{
		LoginName: []string{"username"},
		Password:  "password",
	}
	res, err := client.UtilLogin(&req)
	if err != nil {
		t.Fatalf("error is occured. %s", err.Error())
	}
	if !strings.Contains(res.Cookie, "CBSESSID=C735B4069ccf104Ce0f2bf12a7cc62f115db9e676f6e72f2;") {
		t.Fatalf("expect %v, get %v", "", res.Cookie)
	}
	assert(t, res.LoginName, "Administrator")
	assert(t, res.Status, "Login")
}

func assert(t *testing.T, expect interface{}, actual interface{}) {
	if expect != actual {
		t.Fatalf("expect %v, get %v", expect, actual)
	}
}
