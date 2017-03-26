[![Build Status](https://travis-ci.org/tzmfreedom/goroon.svg?branch=master)](https://travis-ci.org/tzmfreedom/goroon)

# Goroon

Command Line Interface and Library for Garoon with Go language

## Install

### CLI

MacOS user can use homebrew to install goroon
```bash
$ brew tap tzmfreedom/goroon
$ brew install goroon
```

If you want to use most recently version of goroon, execute following command.
```bash
$ go get -u github.com/tzmfreedom/goroon/goroon
```

### Library

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
   0.1.0

COMMANDS:
     login, l     login to garoon
     schedule, s  get your schedule
     bulletin, b  get bulletin
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Get schedule on target date
```
NAME:
   goroon schedule - get your schedule

USAGE:
   goroon schedule [command options] [arguments...]

OPTIONS:
   --username value, -u value   [$GAROON_USERNAME]
   --password value, -p value   [$GAROON_PASSWORD]
   --endpoint value, -e value   [$GAROON_ENDPOINT]
   --userid value, -i value
   --start value
   --end value
   --type value, -t value      (default: "all")
   --columns value, -c value   (default: "detail,start,end")
   --debug, -D
   --date value, -d value
```

Get my schedule on target date
```bash
$ goroon schedule -u {USERNAME} -p {PASSWORD} -e {ENDPOINT} --start {START DATETIME} --end {END DATETIME}
$ goroon schedule -u {USERNAME} -p {PASSWORD} -e {ENDPOINT} -d {DATE:[today|yesterday]}

ex)
$ goroon schedule -u hoge -p fuga -e https//hoge.garoon.com/grn.exe \
  --start '2017-03-01 00:00:00' --end '2017-03-01 23:59:59'
$ goroon schedule -u hoge -p fuga -e https//hoge.garoon.com/grn.exe -d today # get today's your schedule
$ goroon schedule -u hoge -p fuga -e https//hoge.garoon.com/grn.exe # get today's your schedule
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
   --offset value, -o value    (default: 0)
   --limit value, -l value     (default: 20)
   --debug, -D
   --columns value, -c value   (default: "creator,text")
```

```bash
$ goroon bulleting -u {USERNAME} -p {PASSWORD} -e {ENDPOINT}
```

Login to garoon
```bash
NAME:
   goroon login - login to garoon

USAGE:
   goroon login [command options] [arguments...]

OPTIONS:
   --username value, -u value   [$GAROON_USERNAME]
   --password value, -p value   [$GAROON_PASSWORD]
   --endpoint value, -e value   [$GAROON_ENDPOINT]
   --debug, -D
ex)
$ goroon login -u hoge -p fuga -e "https://hoge.garoon.com/grn.exe"
```
The credentials, that is username, password, endpoint, is required on any goroon command.
For this reason, it's hard to use these commands. Login subcommand solve the problem.
If you call garoon login subcommand, goroon create `.goroon` file on your home path and write session id and endpoint on it.
When `.goroon` file exists, every command read `.goroon` file and use information on file to login to garoon.

For example, if you called login command, you can use the following command to get your schedule on today.
```bash
$ goroon schedule
```

### Library

Initialize by credentials

```
client := goroon.NewClient("input garoon endpoint")
client.Username = "input your username"
client.Password = "input your password"
```

Initialize by sessionId

```
client := goroon.NewClient("input garoon endpoint")
client.SessionId = "input your sessionId"
```

Get my schedule on target date

```golang
start := time.Date(2017, 1, 1, 0, 0, 0, 0, time.Local)
end := time.Date(2017, 1, 1, 23, 59, 59, 999999, time.Local)

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
  if event.When.Datetime.Start == new(time.Time) {
    fmt.Println(event.When.Date.Start.Format("2006-01-02"))
    fmt.Println(event.When.Date.End.Format("2006-01-02")
  } else {
    fmt.Println(event.When.Datetime.Start.Format("2006-01-02T15:04:05"))
    fmt.Println(event.When.Datetime.End.Format("2006-01-02T15:04:05"))
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
client := goroon.NewClient("https://garoon.hogehoge.com/grn.exe")
client.Debugger = os.Stdout
```
Debugger property is io.Writer interface.
You can debug out to anywhere with io.Writer.