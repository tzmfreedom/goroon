package goroon

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	endpoint string
	header   *SoapHeader
	Debugger io.Writer
}

type NopWriter struct{}

func (d *NopWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func NewClient(username string, password string, endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		Debugger: &NopWriter{},
		header: &SoapHeader{
			Security: Security{
				UsernameToken: UsernameToken{
					Username: username,
					Password: password,
				},
			},
			Locale:    "jp",
			Timestamp: Timestamp{},
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
	envelope.SoapHeader.Timestamp.Created = created
	expires := created.Add(time.Duration(1) * time.Hour)
	envelope.SoapHeader.Timestamp.Expires = expires

	envelope.SoapBody = &SoapBody{Content: req}
	b, err := xml.MarshalIndent(envelope, "", "	")
	c.Debugger.Write(b)
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

	c.Debugger.Write(body)
	res_env := &SoapEnvelope{SoapBody: &SoapBody{Content: res}}
	err = xml.Unmarshal(body, res_env)
	if err != nil {
		return err
	}
	if res_env.SoapBody.Fault != nil {
		msg := fmt.Sprintf("Soap Fault is occured: %s: %s", res_env.SoapBody.Fault.Code, res_env.SoapBody.Fault.Detail)
		return errors.New(msg)
	}
	return nil
}

func (c *Client) ScheduleGetEventsByTarget(params *Parameters) (*Returns, error) {
	uri := fmt.Sprintf("%s/cbpapi/schedule/api", c.endpoint)
	req := &ScheduleGetEventsByTargetRequest{
		Parameters: params,
	}
	res := &ScheduleGetEventsByTargetResponse{}
	err := c.Request("ScheduleGetEventsByTarget", uri, req, res)
	if err != nil {
		return nil, err
	}
	return res.Returns, nil
}

func (c *Client) UtilGetLoginUserId(params *Parameters) (*Returns, error) {
	uri := fmt.Sprintf("%s/cbpapi/util/api", c.endpoint)
	req := &UtilGetLoginUserIdRequest{
		Parameters: params,
	}
	res := &UtilGetLoginUserIdResponse{}
	err := c.Request("UtilGetLoginUserId", uri, req, res)
	if err != nil {
		return nil, err
	}
	return res.Returns, nil
}

func (c *Client) ScheduleGetEvents(params *Parameters) (*Returns, error) {
	uri := fmt.Sprintf("%s/cbpapi/schedule/api", c.endpoint)
	req := &ScheduleGetEventsRequest{
		Parameters: params,
	}
	res := &ScheduleGetEventsResponse{}
	err := c.Request("ScheduleGetEvents", uri, req, res)
	if err != nil {
		return nil, err
	}
	return res.Returns, nil
}

func (c *Client) BaseGetUserByLoginName(params *Parameters) (*Returns, error) {
	uri := fmt.Sprintf("%s/cbpapi/base/api", c.endpoint)
	req := &BaseGetUserByLoginNameRequest{
		Parameters: params,
	}
	res := &BaseGetUserByLoginNameResponse{}
	err := c.Request("BaseGetUserByLoginName", uri, req, res)
	if err != nil {
		return nil, err
	}
	return res.Returns, nil
}

func (c *Client) BulletinGetFollows(params *Parameters) (*Returns, error) {
	uri := fmt.Sprintf("%s/cbpapi/bulletin/api", c.endpoint)
	req := &BulletinGetFollowsRequest{
		Parameters: params,
	}
	res := &BulletinGetFollowsResponse{}
	err := c.Request("BulletinGetFollows", uri, req, res)
	if err != nil {
		return nil, err
	}
	return res.Returns, nil
}
