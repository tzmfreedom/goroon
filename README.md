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

res, err := client.ScheduleGetEvents(&goroon.Parameters{
  Start: start,
  End:   end,
})
if err != nil {
  return err
}
for _, sch := range res.ScheduleEvents {
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
res, err := client.BaseGetUserByLoginName(&goroon.Parameters{
  LoginName: []string{"hogehoge"},
})
if err != nil {
  return err
}
fmt.Println(res.LoginId)
```

Get Bulletin folows

```golang
res, err := client.BulletinGetFollows(&goroon.Parameters{
  TopicId: 1234,
  Offset:  0,
  Limit:   20,
})
if err != nil {
  return err
}

for _, follow := range res.Follow {
  fmt.Println(fmt.Sprint(follow.Number)
  fmt.Println(follow.Creator.Name)
  fmt.Println(follow.Text)
}
```
