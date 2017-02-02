package goroon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

type Envelope struct {
	XMLName      xml.Name     `xml:"Envelope"`
	ResponseBody ResponseBody `xml:"Body"`
}

type ResponseBody struct {
	XMLName        xml.Name        `xml:"Body"`
	ScheduleEvents []ScheduleEvent `xml:"ScheduleGetEventsByTargetResponse>returns>schedule_event"`
}

type ScheduleGetEventsByTargetResponse struct {
	XMLName xml.Name `xml:"ScheduleGetEventsByTargetResponse"`
	Returns Returns  `xml:"returns"`
}

type Returns struct {
	XMLName       xml.Name      `xml:"returns"`
	ScheduleEvent ScheduleEvent `xml:"schedule_event"`
}

type ScheduleEvent struct {
	XMLName     xml.Name     `xml:"schedule_event"`
	Members     []Members    `xml:"members"`
	RepeatInfo  []RepeatInfo `xml:"repeat_info"`
	When        When         `xml:"when"`
	Detail      string       `xml:"detail,attr"`
	Description string       `xml:"description,attr"`
	Id          int          `xml:"id,attr"`
	EventType   string       `xml:"event_type,attr"`
}

type RepeatInfo struct {
	XMLName xml.Name `xml:"repeat_info"`
}

type Members struct {
	XMLName xml.Name `xml:"members"`
	Member  Member   `xml:"member`
}

type Member struct {
	XMLName xml.Name `xml:"member"`
	User    User     `xml:"user"`
}

type User struct {
	XMLName xml.Name `xml:"user"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type When struct {
	XMLName  xml.Name `xml:"when"`
	Datetime Datetime `xml:"datetime"`
	Date     Date     `xml:"date"`
}

type Datetime struct {
	XMLName xml.Name  `xml:"datetime"`
	Start   time.Time `xml:"start,attr"`
	End     time.Time `xml:"end,attr"`
}

type Date struct {
	XMLName xml.Name `xml:"date"`
	Start   xmlDate  `xml:"start,attr"`
	End     xmlDate  `xml:"end,attr"`
}

type xmlDate struct {
	time.Time
}

func (c *xmlDate) UnmarshalXMLAttr(attr xml.Attr) error {
	const shortForm = "2006-01-02"
	parse, err := time.Parse(shortForm, attr.Value)
	if err != nil {
		return err
	}
	*c = xmlDate{parse}
	return nil
}


func NewLogger(outStream io.Writer) *logrus.Logger {
	outLogger := logrus.New()
	outLogger.Out = outStream
	return outLogger
}

type GaroonClient struct {
	Username string
	Password string
	Endpoint string
	IsDebug  bool
	Logger   *logrus.Logger
}

func NewGaroonClient(username string, password string, endpoint string) *GaroonClient {
	return &GaroonClient{
		Username: username,
		Password: password,
		Endpoint: endpoint,
		Logger:   NewLogger(os.Stdout),
	}
}

func (client *GaroonClient) SetDebug(isDebug bool) {
	client.IsDebug = isDebug
}

func (g *GaroonClient) Request(userId string, start time.Time, end time.Time) (res *Envelope, err error) {
	soapTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope
  xmlns:soap="http://www.w3.org/2003/05/soap-envelope">
  <soap:Header>
    <Action>
      ScheduleGetEventsByTarget
    </Action>
    <Security>
      <UsernameToken>
        <Username>%s</Username>
        <Password>%s</Password>
      </UsernameToken>
    </Security>
    <Timestamp>
      <Created>2016-12-05T14:45:00Z</Created>
      <Expires>2037-08-12T14:45:00Z</Expires>
    </Timestamp>
    <Locale>jp</Locale>
  </soap:Header>
  <soap:Body>
    <ScheduleGetEventsByTarget>
      <parameters start="%s" end="%s">
        <user id="%s"></user>
      </parameters>
    </ScheduleGetEventsByTarget >
  </soap:Body>
</soap:Envelope>`

	soapMessage := fmt.Sprintf(soapTemplate, g.Username, g.Password, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), userId)
	g.Debug(soapMessage)
	resp, err := http.Post(fmt.Sprintf("%s/cbpapi/schedule/api", g.Endpoint), "text/xml", strings.NewReader(soapMessage))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	soapResponse := bytes.NewBuffer(body).String()
	g.Debug(soapResponse)
	res = &Envelope{}
	err = xml.Unmarshal([]byte(soapResponse), res)
	return
}

func (event *ScheduleEvent) IsBanner() bool {
	return event.EventType == "banner"
}

func (event *ScheduleEvent) GetStartStr() string {
	if event.IsBanner() {
		return fmt.Sprintf("%s00:00:00", event.When.Date.Start.Format("2006-01-02T"))
	}
	return event.When.Datetime.Start.Format("2006-01-02T15:04:05")
}

func (event *ScheduleEvent) GetEndStr() string {
	if event.IsBanner() {
		return fmt.Sprintf("%s00:00:00", event.When.Date.Start.Format("2006-01-02T"))
	}
	return event.When.Datetime.End.Format("2006-01-02T15:04:05")
}

func (event *ScheduleEvent) GetId() string {
	var tm int64
	if event.IsBanner() {
		tm = event.When.Date.Start.Unix()
	} else {
		tm = event.When.Datetime.Start.Unix()
	}
	return fmt.Sprintf("%s-%s", event.GetId(), tm)
}

func (g *GaroonClient) Debug(args ...interface{}) {
	if g.IsDebug {
		g.Logger.Debug(args...)
	}
}

func (g *GaroonClient) Debugf(format string, args ...interface{}) {
	if g.IsDebug {
		g.Logger.Debugf(format, args...)
	}
}
