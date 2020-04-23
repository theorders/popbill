package popbill

import (
	"github.com/theorders/aefire"
	"net/http"
)

func (c *Client)SendLMS(rcvNum, messaage string) (receipt *Receipt, err error) {
	rcvNum = aefire.LocalizePhoneNumber(rcvNum, 82)

	res, err := c.Request(
		http.MethodPost,
		LMSService,
		"",
		aefire.MapOf(
		"msgs", []map[string]string{
			aefire.StringMapOf(
				"snd", SendNumber,
				"rcv", rcvNum,
				"msg", messaage,
			),
		}))

	if err != nil {
		return nil, err
	}

	receipt = &Receipt{}
	if err := res.ToJSON(receipt) ; err != nil{
		return nil, aefire.NewHttpError(err)
	}

	return receipt, err
}
