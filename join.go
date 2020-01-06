package popbill

import (
	"net/http"
)

type JoinParam struct {
	ID           string
	PWD          string
	LinkID       string
	CorpNum      string
	CEOName      string
	CorpName     string
	Addr         string
	BizType      string
	BizClass     string
	ContactName  string
	ContactEmail string
	ContactTEL   string
	ContactHP    string
	ContactFAX   string
}

func (c *Client) Join(param JoinParam) error {
	if param.BizType == "" {
		param.BizType = "-"
	}
	if param.BizClass == "" {
		param.BizClass = "-"
	}
	if param.Addr == "" {
		param.Addr = "-"
	}
	if param.ContactHP == "" {
		param.ContactHP = "-"
	}
	if param.ContactFAX == "" {
		param.ContactFAX = "-"
	}
	if param.ContactTEL == "" {
		param.ContactTEL = "-"
	}
	if param.ContactEmail == "" {
		param.ContactEmail = "-"
	}

	_, err := c.Request(
		http.MethodPost,
		JoinService,
		"",
		param)

	return err
}
