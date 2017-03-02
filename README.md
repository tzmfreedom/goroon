# Goroon

Command Line Interface for Garoon with Go language

## Install

CLI

```bash
$ go get github.com/tzmfreedom/goroon/goroon
```

Library

```bash
$ go get github.com/tzmfreedom/goroon
```

## Usage

Get my schedule on target date

```bash
$ goroon schedule -u {USERNAME} -p {PASSWORD} -e {ENDPOINT} --start {START DATETIME} --end {END DATETIME}

ex)
$ goroon schedule -u hoge -p fuga -e https//hoge.garoon.com/grn.exe \
  --start '2017-03-01 00:00:00' --end '2017-03-02 00:00:00'
```

### Library

Initialize
```
client := goroon.NewClient("username", "password", "https://garoon.hogehoge.com/grn.exe", false, os.Stdout)
```

Get my schedule on target date

```golang
start, err := time.Parse("2006-01-02 15:04:05", "2017-03-01 00:00:00")
if err != nil {
  return err
}
end, err := time.Parse("2006-01-02 15:04:05","2017-03-02 00:00:00")
if err != nil {
  return err
}

req := &goroon.ScheduleGetEventsRequest{
  Parameters: &goroon.Parameters{
    Start: &start,
    End:   &end,
  },
}
res, err := client.ScheduleGetEvents(req)
if err != nil {
  return err
}
for _, sch := range res.Returns.Schedule.ScheduleEvents {
  fmt.Println(event.Id)
  fmt.Println(event.Members)
  fmt.Println(event.EventType)
  fmt.Println(event.Detail)
  fmt.Println(event.Description)
  if event.EventType == "banner" {
    event.When.Date.Start.Format("2006-01-02")
    event.When.Date.End.Format("2006-01-02")
  } else {
    event.When.Datetime.Start.Format("2006-01-02T15:04:05")
    event.When.Datetime.End.Format("2006-01-02T15:04:05")
  }
}
```

Get UserId from login name

```golang
req := &goroon.BaseGetUserByLoginNameRequest{
  Parameters: &goroon.Parameters{
    LoginName: []*string{"hogehoge"},
  },
}
res, err := client.BaseGetUserByLoginName(req)
if err != nil {
  return err
}
fmt.Println(res.Returns.LoginId)
```

Get Bulletin folows

```golang
req := &goroon.BulletinGetFollowsRequest{
  Parameters: &goroon.Parameters{
    TopicId: 1234,
    Offset:  0,
    Limit:   20,
  },
}

res, err := client.BulletinGetFollows(req)
if err != nil {
  return err
}

for _, follow := range res.Returns.Follow {
  fmt.Println(fmt.Sprint(follow.Number)
  fmt.Println(follow.Creator.Name)
  fmt.Println(follow.Text)
}
```
