package popbill

import (
	"github.com/theorders/aefire"
	"net/http"
)


func (c *Client) SendATS(template *ATSTemplate) (receipt *Receipt, err error) {
	res, err := c.Request(
		http.MethodPost,
		ATSService,
		"",
		aefire.MapOf(
			"templateCode", template.TemplateCode,
			"snd", SendNumber,
			"content", template.Content,
			"altContent", template.AltContent(),
			"altSendType", "A",
			"msgs", []map[string]string{
				aefire.StringMapOf(
					"rcv", aefire.LocalizePhoneNumber(template.PhoneNumber, 82),
				),
			}))

	if err != nil {
		return
	}

	receipt = &Receipt{}
	res.ToJSON(receipt)

	return receipt, err
}

type ATSTemplate struct {
	TemplateCode string
	Content      string
	SMSContent   string
	PhoneNumber  string
}

func (t *ATSTemplate) AltContent() string {
	if t.SMSContent != "" {
		return t.SMSContent
	} else
	{
		return t.Content
	}
}

func (t *ATSTemplate) AltSentType() string {
	if t.SMSContent != "" {
		return "A"
	} else
	{
		return "C"
	}
}
