package popbill

import (
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/theorders/aefire"
	"net/http"
)

type Client struct {
	context.Context
	Test    bool
	LinkId  string
	CorpNum string
	Secret  string
}

func NewClient(c context.Context, test bool, linkId, corpNum, secret string) *Client {
	return &Client{
		Context: c,
		Test:    test,
		LinkId:  linkId,
		CorpNum: corpNum,
		Secret:  secret,
	}
}

func (c *Client) Request(method, service, path string, body interface{}, headers ...string) (res *req.Resp, err error) {

	var token *SessionToken

	token, err = c.ServiceToken(service)

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
			endpoint(c.Test, service, path),
			req.BodyJSON(body),
			req.Header(header),
			c.Context)
	} else {
		m := aefire.ToMap(body)

		res, err = req.Get(
			endpoint(c.Test, service, path),
			req.QueryParam(m),
			req.Header(header),
			c.Context)
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

func (c *Client) MultipartFormDataRequest(service, path string, params req.Param, files ...req.FileUpload) (res *req.Resp, err error) {

	var token *SessionToken

	token, err = c.ServiceToken(service)

	if aefire.LogIfError(err) {
		println("linkhub token issue failed:" + err.Error())
		return
	}

	header := aefire.StringMapOf()

	//ods.DebugLog("sessionToken : %s", token.SessionToken)

	header["Authorization"] = "Bearer " + token.SessionToken
	header["Content-Type"] = "application/json; charset=utf8"
	header["Accept-Encoding"] = "application/json; charset=utf8"
	header["x-pb-version"] = apiVersion
	header["x-lh-forwarded"] = "*"

	res, err = req.Post(
		endpoint(c.Test, service, path),
		c.Context,
		req.Header(header),
		params,
		files)

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
