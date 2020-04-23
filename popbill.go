package popbill

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	"github.com/labstack/echo/v4"
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

func (c *Client) Request(method, service, path string, body interface{}, headers ...string) (*req.Resp, *echo.HTTPError) {

	var res *req.Resp
	var token *SessionToken
	var err error

	token, err = c.ServiceToken(service)

	if aefire.LogIfError(err) {
		println("linkhub token issue failed:" + err.Error())
		return nil, aefire.NewEchoHttpError(500, err)
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
		return res, &echo.HTTPError{
			Code:     500,
			Message:  err.Error(),
			Internal: err,
		}
	}

	defaultResponse := DefaultResponse{}

	err = res.ToJSON(&defaultResponse)
	if err != nil {
		return res, &echo.HTTPError{
			Code:     500,
			Message:  err.Error(),
			Internal: err,
		}
	}


	if res.Response().StatusCode / 100 == 4 ||  defaultResponse.Code < 0 {
		return res, echo.NewHTTPError(400, defaultResponse.Message)
	}

	return res, nil
}

func (c *Client) MultipartFormDataRequest(service, path string, params req.Param, files ...req.FileUpload) (*req.Resp, *echo.HTTPError) {

	var res *req.Resp
	var token *SessionToken
	var err error

	token, err = c.ServiceToken(service)

	if aefire.LogIfError(err) {
		println("linkhub token issue failed:" + err.Error())
		return nil, aefire.NewEchoHttpError(500, err)
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
		return res, aefire.NewEchoHttpError(500, err)
	}

	defaultResponse := DefaultResponse{}

	err = res.ToJSON(&defaultResponse)
	if err != nil {
		return res, aefire.NewEchoHttpError(500, err)
	}

	if res.Response().StatusCode / 100 == 4 ||  defaultResponse.Code < 0 {
		return res, echo.NewHTTPError(400, fmt.Sprintf("[%d]%s", defaultResponse.Code, defaultResponse.Message))
	}

	return res, nil
}

func (c *Client) MethodOverrideRequest(
	method, service, path string, body interface{}, overrideMethod string) (res *req.Resp, err *echo.HTTPError) {

	return c.Request(
		method,
		service,
		path,
		body,
		"X-HTTP-Method-Override", overrideMethod)
}
