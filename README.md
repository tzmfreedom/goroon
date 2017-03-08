[![Build Status](https://travis-ci.org/tzmfreedom/goroon.svg?branch=master)](https://travis-ci.org/tzmfreedom/goroon)

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

```
NAME:
   goroon - garoon utility

USAGE:
   goroon [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     schedule, s  get today's your schedule
     bulletin, b  get bulletin
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Get schedule on target date
```
NAME:
   goroon schedule - get today's your schedule

USAGE:
   goroon schedule [command options] [arguments...]

OPTIONS:
   --username value, -u value   [$GAROON_USERNAME]
   --password value, -p value   [$GAROON_PASSWORD]
   --endpoint value, -e value   [$GAROON_ENDPOINT]
   --userid value, -i value
   --debug, -d
   --start value
   --end value
   --type value                (default: "all")
```

Get my schedule on target date
```bash
$ goroon schedule -u {USERNAME} -p {PASSWORD} -e {ENDPOINT} --start {START DATETIME} --end {END DATETIME}

ex)
$ goroon schedule -u hoge -p fuga -e https//hoge.garoon.com/grn.exe \
  --start '2017-03-01 00:00:00' --end '2017-03-02 00:00:00'
```

Get target user's schedule
```bash
$ goroon schedule -u {USERNAME} -p {PASSWORD} -e {ENDPOINT} --start {START DATETIME} --end {END DATETIME} --userid {USER ID}
```

Get bulletin
```
NAME:
   goroon bulletin - get bulletin

USAGE:
   goroon bulletin [command options] [arguments...]

OPTIONS:
   --username value, -u value   [$GAROON_USERNAME]
   --password value, -p value   [$GAROON_PASSWORD]
   --endpoint value, -e value   [$GAROON_ENDPOINT]
   --topic_id value            (default: 0)
   --debug, -d
   --offset value, -o value    (default: 0)
   --limit value, -l value     (default: 0)
```

```bash
$ goroon bulleting -u {USERNAME} -p {PASSWORD} -e {ENDPOINT}
```

### Library

Initialize
```
client := goroon.NewClient("username", "password", "https://garoon.hogehoge.com/grn.exe")
```

Get my schedule on target date

```golang
start, err := time.Parse("2006-01-02 15:04:05", "2017-03-01 00:00:00")
end, err := time.Parse("2006-01-02 15:04:05","2017-03-02 00:00:00")

res, err := client.ScheduleGetEvents(&goroon.Parameters{
  Start: goroon.XmlDateTime{start},
  End:   goroon.XmlDateTime{end},
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
  fmt.Println(fmt.Sprint(follow.Number))
  fmt.Println(follow.Creator.Name)
  fmt.Println(follow.Text)
}
```

Debug your request
```
client := goroon.NewClient("username", "password", "https://garoon.hogehoge.com/grn.exe")
client.Debugger = os.Stdout
```