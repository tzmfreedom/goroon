package goroon

import (
	"encoding/xml"
	"time"
)

type SoapEnvelope struct {
	XMLName    xml.Name    `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapHeader *SoapHeader `xml:"http://www.w3.org/2003/05/soap-envelope Header,omitempty"`
	SoapBody   *SoapBody   `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type SoapBody struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",omitempty"`
	Fault   *SoapFault  `xml:",omitempty"`
}

type SoapHeader struct {
	Action    string     `xml:"Action"`
	Security  *Security  `xml:"Security"`
	Timestamp *Timestamp `xml:"Timestamp"`
	Locale    string     `xml:"Locale"`
}

type SoapFault struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type Security struct {
	UsernameToken *UsernameToken `xml:"UsernameToken,omitempty`
}

type UsernameToken struct {
	Username string `xml:"Username"`
	Password string `xml:"Password"`
}

type Timestamp struct {
	Created *time.Time `xml:"Created"`
	Expires *time.Time `xml:"Expires"`
}

type ScheduleGetEventsByTargetRequest struct {
	XMLName    xml.Name    `xml:"ScheduleGetEventsByTarget"`
	Parameters *Parameters `xml:"parameters"`
}

type Parameters struct {
	Start *time.Time `xml:"start,attr"`
	End   *time.Time `xml:"end,attr"`
	User  *User      `xml:"user"`
}

type ScheduleGetEventsByTargetResponse struct {
	XMLName xml.Name `xml:"ScheduleGetEventsByTargetResponse"`
	Returns *Returns `xml:"returns"`
}

type Returns struct {
	ScheduleEvents []*ScheduleEvent `xml:"schedule_event"`
}

type ScheduleEvent struct {
	Members     []*Members    `xml:"members"`
	RepeatInfo  []*RepeatInfo `xml:"repeat_info"`
	When        *When         `xml:"when"`
	Detail      string        `xml:"detail,attr"`
	Description string        `xml:"description,attr"`
	Id          int           `xml:"id,attr"`
	EventType   string        `xml:"event_type,attr"`
}

type RepeatInfo struct {
	XMLName xml.Name `xml:"repeat_info"`
}

type Members struct {
	XMLName xml.Name `xml:"members"`
	Member  *Member  `xml:"member`
}

type Member struct {
	XMLName xml.Name `xml:"member"`
	User    *User    `xml:"user"`
}

type User struct {
	XMLName xml.Name `xml:"user"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type When struct {
	XMLName  xml.Name  `xml:"when"`
	Datetime *Datetime `xml:"datetime"`
	Date     *Date     `xml:"date"`
}

type Datetime struct {
	XMLName xml.Name   `xml:"datetime"`
	Start   *time.Time `xml:"start,attr"`
	End     *time.Time `xml:"end,attr"`
}

type Date struct {
	XMLName xml.Name `xml:"date"`
	Start   *xmlDate `xml:"start,attr"`
	End     *xmlDate `xml:"end,attr"`
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

func (b *SoapBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SoapFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}
