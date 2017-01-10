package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

type GaroonClient struct {
	Username string
	Password string
	Endpoint string
}

func NewGaroonClient(username string, password string, endpoint string) *GaroonClient {
	return &GaroonClient{
		Username: username,
		Password: password,
		Endpoint: endpoint,
	}
}

func emulateResponse() string {
	// 本日分のスケジュールを取得
	soapResponse := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope
xmlns:soap="http://www.w3.org/2003/05/soap-envelope"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:xsd="http://www.w3.org/2001/XMLSchema"
xmlns:schedule="http://wsdl.cybozu.co.jp/schedule/2008">
<soap:Header>
    <vendor>Cybozu</vendor>
    <product>Garoon</product>
    <product_type>1</product_type>
    <version>3.7.5</version>
    <apiversion>1.3.1</apiversion>
</soap:Header>
<soap:Body>
    <schedule:ScheduleGetEventsByTargetResponse>
        <returns>
            <schedule_event id="123"
                detail="fugafuga"
                description="hogehoge"
                >
                <members xmlns="http://schemas.cybozu.co.jp/schedule/2008">
                    <member>
                        <user id="aa" name="bb" order="0" />
                    </member>
                </members>
                <repeat_info xmlns="http://schemas.cybozu.co.jp/schedule/2008">
                <condition type="week" day="20"
                    week="2" start_date="2016-11-22" end_date="2017-04-01"
                    start_time="14:00:00" end_time="14:30:00"/>
                    <exclusive_datetimes>
                        <exclusive_datetime start="2016-12-13T00:00:00+09:00" end="2016-12-14T00:00:00+09:00" />
                        <exclusive_datetime start="2016-12-20T00:00:00+09:00" end="2016-12-21T00:00:00+09:00" />
                    </exclusive_datetimes>
                </repeat_info>
                <when>
                    <!--<datetime start="2016-12-15T13:07:00Z" end="2016-12-15T16:30:00Z" />-->
                    <date start="2016-12-15" end="2016-12-16" />
                </when>
            </schedule_event>
        </returns>
    </schedule:ScheduleGetEventsByTargetResponse>
</soap:Body>
</soap:Envelope>`
	return soapResponse
}

//func (g *GaroonClient)Request(userId string, start time.Time, end time.Time) (res *Envelope, err error) {
//	soapResponse := emulateResponse()
//	res = &Envelope{}
//	err = xml.Unmarshal([]byte(soapResponse), res)
//	pp.Print(res)
//	return
//}

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
	resp, err := http.Post(g.Endpoint, "text/xml", strings.NewReader(soapMessage))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	soapResponse := bytes.NewBuffer(body).String()
	res = &Envelope{}
	err = xml.Unmarshal([]byte(soapResponse), res)
	return
}
