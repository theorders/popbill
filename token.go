package popbill

import (
	"encoding/json"
	"fmt"
	"github.com/theorders/aefire"
	"io/ioutil"
	"os"
	"time"
)

type SessionToken struct {
	Code         int      `json:"code"`
	Message      string   `json:"message"`
	SessionToken string   `json:"session_token"`
	ServiceID    string   `json:"serviceID"`
	LinkID       string   `json:"linkID"`
	Usercode     string   `json:"usercode"`
	Ipaddress    string   `json:"ipaddress"`
	Expiration   string   `json:"expiration"`
	Scope        []string `json:"scope"`
}

func (token SessionToken) ExpiresAt() time.Time {
	t, e := time.Parse("2006-01-02T15:04:05.999Z07:00", token.Expiration)

	if e != nil {
		return time.Time{}
	}

	return t
}

func (c *Client) ServiceToken(service string) (token *SessionToken, err error) {
	if service == JoinService {
		service = TaxinvoiceService
	}

	tokenKey := fmt.Sprintf("LINKHUB_%t_%s_%s_TOKEN", c.Test, c.CorpNum, service)
	keyPath := "/tmp/" + tokenKey

	if _, err := os.Stat(keyPath); err == nil {
		//ods.DebugLog("key exists in tmp dir")

		b, err := ioutil.ReadFile(keyPath)
		token = &SessionToken{}

		if !aefire.LogIfError(err) &&
			!aefire.LogIfError(json.Unmarshal(b, token)) {

			//ods.DebugLog("token.ExpiresAt()=" + token.ExpiresAt().String())
			if token.ExpiresAt().After(time.Now().Add(time.Minute)){
				return token, nil
			} else {
				aefire.LogIfError(os.Remove(keyPath))
			}
		} else {
			aefire.LogIfError(os.Remove(keyPath))
		}
	}


	builder := c.Instance()

	switch service {
	case SMSService, LMSService:
		builder.AddScopes("150", "151", "152")
	case ATSService:
		builder.AddScopes("153", "154", "155")
	case FAXService:
		builder.AddScopes("160")
	case CloseDownService:
		builder.AddScopes("170")
	case TaxinvoiceService:
		builder.AddScopes("110")
	case CashbillService:
		builder.AddScopes("140")
	}

	token, err = builder.Build("*")

	aefire.PanicIfError(err)

	aefire.LogIfError(ioutil.WriteFile(keyPath, []byte(aefire.ToJson(*token)), os.ModePerm))

	return
}
