package popbill

import (
	"strings"
)

func endpoint(isTest bool, service string, pathComponents ...string) string {
	if !isTest {
		return "https://popbill.linkhub.co.kr/" + service + "/" + strings.TrimLeft(strings.Join(pathComponents, "/"), "/")
	} else
	{
		return "https://popbill_test.linkhub.co.kr/" + service + "/" + strings.TrimLeft(strings.Join(pathComponents, "/"), "/")
	}
}

type DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (response DefaultResponse) Error() string {
	return response.Message
}

