package goroon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

func newLogger(outStream io.Writer, debug bool) *logrus.Logger {
	l := logrus.New()
	l.Out = outStream
	if debug {
		l.Level = logrus.DebugLevel
	}
	return l
}

type Client struct {
	username string
	password string
	endpoint string
	header   *SoapHeader
	Debugger Debugger
}

type Debugger interface {
	Debug(args ...interface{})
}

func NewClient(username string, password string, endpoint string, debug bool, w io.Writer) *Client {
	return &Client{
		endpoint: endpoint,
		Debugger: newLogger(w, debug),
		header: &SoapHeader{
			Security: &Security{
				UsernameToken: &UsernameToken{
					Username: username,
					Password: password,
				},
			},
			Locale:    "jp",
			Timestamp: &Timestamp{},
		},
	}
}

func (c *Client) SetHeader(header *SoapHeader) {
	c.header = header
}

func (c *Client) Request(action string, uri string, req interface{}, res interface{}) error {
	envelope := &SoapEnvelope{}
	envelope.SoapHeader = c.header
	envelope.SoapHeader.Action = action

	created := time.Now()
	envelope.SoapHeader.Timestamp.Created = &created
	expires := created.Add(time.Duration(1) * time.Hour)
	envelope.SoapHeader.Timestamp.Expires = &expires

	envelope.SoapBody = &SoapBody{Content: req}
	b, err := xml.MarshalIndent(envelope, "", "	")
	c.Debugger.Debug(string(b))
	msg, err := xml.Marshal(envelope)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(msg)
	resp, err := http.Post(uri, "text/xml", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	c.Debugger.Debug(bytes.NewBuffer(body).String())
	res_env := &SoapEnvelope{SoapBody: &SoapBody{Content: res}}
	err = xml.Unmarshal(body, res_env)
	return err
}

func (c *Client) ScheduleGetEventsByTarget(req *ScheduleGetEventsByTargetRequest) (*ScheduleGetEventsByTargetResponse, error) {
	uri := fmt.Sprintf("%s/cbpapi/schedule/api", c.endpoint)
	res := &ScheduleGetEventsByTargetResponse{}
	err := c.Request("ScheduleGetEventsByTarget", uri, req, res)
	return res, err
}

func (c *Client) UtilGetLoginUserId(req *UtilGetLoginUserIdRequest) (*UtilGetLoginUserIdResponse, error) {
	uri := fmt.Sprintf("%s/cbpapi/util/api", c.endpoint)
	res := &UtilGetLoginUserIdResponse{}
	err := c.Request("UtilGetLoginUserId", uri, req, res)
	return res, err
}

func (c *Client) ScheduleGetEvents(req *ScheduleGetEventsRequest) (*ScheduleGetEventsResponse, error) {
	uri := fmt.Sprintf("%s/cbpapi/schedule/api", c.endpoint)
	res := &ScheduleGetEventsResponse{}
	err := c.Request("ScheduleGetEvents", uri, req, res)
	return res, err
}

func (c *Client) BaseGetUserByLoginName(req *BaseGetUserByLoginNameRequest) (*BaseGetUserByLoginNameResponse, error) {
	uri := fmt.Sprintf("%s/cbpapi/base/api", c.endpoint)
	res := &BaseGetUserByLoginNameResponse{}
	err := c.Request("BaseGetUserByLoginName", uri, req, res)
	return res, err
}

func (c *Client) BulletinGetFollows(req *BulletinGetFollowsRequest) (*BulletinGetFollowsResponse, error) {
	uri := fmt.Sprintf("%s/cbpapi/bulletin/api", c.endpoint)
	res := &BulletinGetFollowsResponse{}
	err := c.Request("BulletinGetFollows", uri, req, res)
	return res, err
}
