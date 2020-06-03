package pb_cashbill

import (
	"github.com/theorders/aefire"
	"strings"
)

type Customer struct {
	Item

	TradeUsage   TradeUsage `json:"tradeUsage,omitempty" firestore:"tradeUsage,omitempty"`
	IdentityNum  string     `json:"identityNum,omitempty" firestore:"identityNum,omitempty"`
	CustomerName string     `json:"customerName,omitempty" firestore:"customerName,omitempty"`
	Email        string     `json:"email,omitempty" firestore:"email,omitempty"`
}


func (c *Customer) CustomerNameOrIdentityNum() string {
	if c.CustomerName != ""{
		return c.CustomerName
	} else {
		return c.IdentityNum
	}
}

func (c *Customer) NameOrIdentityNumMasked() string {
	if c.IdentityNum == SelfIssueNum {
		return "자진발급"
	} else if c.CustomerName != "" {
		return c.CustomerName
	} else {
		return c.IdentityNumMasked()
	}
}
func (c *Customer) IdentityNumMasked() string {
	idNum := strings.Replace(c.IdentityNum, "-", "", -1)

	if aefire.ValidateCorpNum(idNum) {
		return idNum[:2] + "****" + idNum[6:]
	} else if len(idNum) == 10 {
		return idNum[:3] + "***" + idNum[6:]
	} else if len(idNum) == 11 {
		return idNum[:3] + "****" + idNum[7:]
	} else if len(idNum) == 13 {
		return idNum[:6] + "*******"
	} else {
		return idNum[:8] + "****" + idNum[12:]
	}
}
