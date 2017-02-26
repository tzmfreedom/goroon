package goroon

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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
	logger   *logrus.Logger
}

var Debug = func(args ...interface{}) {

}

func NewClient(username string, password string, endpoint string, debug bool, w io.Writer) *Client {
	logger := newLogger(w, debug)
	if debug {
		Debug = logger.Debug
	}
	return &Client{
		endpoint: endpoint,
		logger:   newLogger(w, debug),
		header: &SoapHeader{
			Security: &Security{
				UsernameToken: &UsernameToken{
					Username: username,
					Password: password,
				},
			},
			Locale: "ja",
		},
	}
}

func (c *Client) SetHeader(header *SoapHeader) {
	c.header = header
}

func (c *Client) Request(uri string, req interface{}, res interface{}) error {
	envelope := &SoapEnvelope{}
	envelope.SoapBody = &SoapBody{Content: req}
	c.logger.Debug(xml.MarshalIndent(envelope, "", "	"))
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

	c.logger.Debug(bytes.NewBuffer(body).String())

	err = xml.Unmarshal(body, res)
	return err
}

func (c *Client) ScheduleGetEventsByTarget(req *ScheduleGetEventsByTargetRequest, res *ScheduleGetEventsByTargetResponse) error {
	uri := fmt.Sprintf("%s/cbpapi/schedule/api", c.endpoint)
	return c.Request(uri, req, res)
}
