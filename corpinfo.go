package popbill

import (
	"github.com/theorders/aefire"
	"strconv"
	"strings"
)

type CorpInfo struct {
	CorpNum     string `firestore:"corpNum" json:"corpNum,omitempty"`
	Type        string `firestore:"type" json:"type"`
	TaxType     string `firestore:"taxType" json:"taxType"`
	TaxRegID    string `firestore:"taxRegID" json:"taxRegID"`
	CorpName    string `firestore:"corpName" json:"corpName"`
	CEOName     string `firestore:"ceoName" json:"ceoName"`
	Addr        string `firestore:"addr" json:"addr"`
	BizClass    string `firestore:"bizClass" json:"bizClass"`
	BizType     string `firestore:"bizType" json:"bizType"`
	ContactName string `firestore:"contactName" json:"contactName"`
	TEL         string `firestore:"tel" json:"tel"`
	Email       string `firestore:"email" json:"email"`
	Birthday    string `firestore:"birthday" json:"birthday"`
}


func (info *CorpInfo) CashbillTo(mgtKey string, customer CashbillIssueV1, trans *Transaction) *Cashbill {
	if trans == nil {
		trans = &Transaction{}
	}

	taxType := TaxationTypeWithTax

	if trans.VAT == 0 {
		taxType = TaxationTypeNoTax
	}

	return &Cashbill{
		CashbillCustomer: CashbillCustomer{
			CashbillItem: CashbillItem{
				TaxationType: taxType,
				TradeOpt:     customer.TradeOpt,
				TotalAmount:  strconv.Itoa(int(trans.Sum)),
				ItemName:     customer.ItemName,
			},
			TradeUsage:   customer.Usage,
			IdentityNum:  customer.IdentityNum,
			CustomerName: customer.CustomerName,
		},
		MgtKey:            mgtKey,
		Tax:               strconv.Itoa(int(trans.VAT)),
		SupplyCost:        strconv.Itoa(int(trans.Supply)),
		ServiceFee:        strconv.Itoa(int(customer.ServiceFee)),
		FranchiseCorpNum:  info.CorpNum,
		FranchiseAddr:     info.Addr,
		FranchiseCEOName:  info.CEOName,
		FranchiseCorpName: info.CorpOrCeoName(),
		FranchiseTEL:      aefire.LocalizePhoneNumber(info.TEL, 82),
		TradeType:         TradeTypeApproval,
		Email:             customer.Email,
		OrderNumber:       mgtKey,
		Aid:               customer.Aid,
	}

}
func (info *CorpInfo) JoinParam(linkId, frLinkId string) JoinParam {
	return JoinParam{
		CEOName:    info.CEOName,
		BizType:    info.BizType,
		BizClass:   info.BizClass,
		Addr:       info.Addr,
		LinkID:     linkId,
		CorpNum:    info.CorpNum,
		CorpName:   info.CorpOrCeoName(),
		ID:         frLinkId,
		ContactTEL: aefire.LocalizePhoneNumber(info.TEL, 82),
		//ContactEmail: info.Email,
		ContactFAX:  "",
		ContactHP:   aefire.LocalizePhoneNumber(info.TEL, 82),
		ContactName: info.CEOName,
		PWD:         info.CorpNum,
	}
}

func (info *CorpInfo) CorpOrCeoName() string {
	if info.CorpName != "" {
		return info.CorpName
	} else {
		return info.CEOName
	}
}

func (info *CorpInfo) LinkPassword() string {
	return info.CorpNum
}

func (info *CorpInfo) IsValid() bool {
	return info.CorpNum != "" && info.CEOName != ""
}

func (info *CorpInfo) IsBizTypeTransport() bool {
	return strings.Contains(info.BizType, "운수")
}
