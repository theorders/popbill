package popbill

import (
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/theorders/aefire"
	"net/http"
)

type Client struct {
	test    bool
	linkId  string
	corpNum string
	secret string
}

func NewClient(test bool, linkId, corpNum, secret string) *Client {
	return &Client{
		test:    test,
		linkId:  linkId,
		corpNum: corpNum,
		secret: secret,
	}
}

func (c *Client) Request(method, service, path string, body interface{}, headers ...string) (res *req.Resp, err error) {

	var token *SessionToken

	token, err = ServiceToken(c.test, c.linkId, c.corpNum, c.secret, service)

	if aefire.LogIfError(err) {
		println("linkhub token issue failed:" + err.Error())
		return
	}

	header := aefire.StringMapOf(headers...)

	//ods.DebugLog("sessionToken : %s", token.SessionToken)

	header["Authorization"] = "Bearer " + token.SessionToken
	header["Content-Type"] = "application/json; charset=utf8"
	header["Accept-Encoding"] = "application/json; charset=utf8"
	header["x-pb-version"] = apiVersion
	header["x-lh-forwarded"] = "*"

	if method == http.MethodPost {
		res, err = req.Post(
			endpoint(c.test, service, path),
			req.BodyJSON(body),
			req.Header(header))
	} else {
		m := aefire.ToMap(body)

		res, err = req.Get(
			endpoint(c.test, service, path),
			req.QueryParam(m),
			req.Header(header))
	}

	if err != nil {
		return res, err
	}

	defaultResponse := DefaultResponse{}

	err = res.ToJSON(&defaultResponse)
	if err != nil {
		return res, err
	}

	if defaultResponse.Code < 0 {
		return res, errors.New(fmt.Sprintf("[%d]%s", defaultResponse.Code, defaultResponse.Message))
	}

	return res, nil
}

func (c *Client) MethodOverrideRequest(
	method, service, path string, body interface{}, overrideMethod string) (res *req.Resp, err error) {

	return c.Request(
		method,
		service,
		path,
		body,
		"X-HTTP-Method-Override", overrideMethod)
}
