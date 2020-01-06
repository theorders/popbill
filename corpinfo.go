package popbill

import (
	"github.com/labstack/echo"
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

type CashbillCustomer struct {
	IdentityNum string       `json:"identityNum" firestore:"identityNum"`
	Name        string       `json:"name" firestore:"name"`
	Email       string       `json:"email" firestore:"email"`
	ItemName    string       `json:"itemName" firestore:"itemName"`
	Usage       TradeUsage   `json:"usage" firestore:"usage"`
	TradeOpt    TradeOpt     `json:"tradeOpt,omitempty" firestore:"tradeOpt"` //(필수)
	Transaction *Transaction `json:"transaction" firestore:"transaction"`
	ServiceFee  int64        `json:"serviceFee" firestore:"serviceFee"`
	Aid         *string      `json:"aid,omitempty" firestore:"aid,omitempty"`
}

func (cashbill *CashbillCustomer) NameOrIdentityNumMasked() string {
	if cashbill.IdentityNum == SelfIssueNum {
		return SelfIssueUsage
	} else if cashbill.Name != "" {
		return cashbill.Name
	} else {
		return cashbill.IdentityNumMasked()
	}
}
func (cashbill *CashbillCustomer) IdentityNumMasked() string {
	idNum := strings.Replace(cashbill.IdentityNum, "-", "", -1)

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

func (customer *CashbillCustomer) Validate() error {
	if customer.Usage == "" {
		return echo.NewHTTPError(400, "거래유형이 지정되지 않았습니다.")
	}

	if customer.IdentityNum == "" {
		return echo.NewHTTPError(400, "고객 식별번호가 없습니다.")
	}

	if customer.Transaction == nil || customer.Transaction.Supply == 0 || customer.Transaction.Sum == 0 {
		return echo.NewHTTPError(400, "거래금액이 없습니다.")
	}

	//소득공제용
	if customer.Usage == TradeUsageIncomeDeduction &&
		!aefire.ValidateRRN(customer.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) {
		return echo.NewHTTPError(400, "소득공제용 현금영수증 발행에는 고객 주민등록번호나 휴대전화번호가 필요합니다.")
	}

	//지출증빙용
	if customer.Usage == TradeUsageIncomeDeduction &&
		!aefire.ValidateRRN(customer.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) &&
		!aefire.ValidateCorpNum(customer.IdentityNum) {
		return echo.NewHTTPError(400, "소득공제용 현금영수증 발행에는 사업자등록번호, 주민등록번호 혹은 휴대전화번호가 필요합니다.")
	}

	return nil
}

func (info *CorpInfo) CashbillTo(mgtKey string, customer CashbillCustomer, trans *Transaction) *Cashbill {
	if trans == nil {
		trans = &Transaction{}
	}

	taxType := TaxationTypeWithTax

	if trans.VAT == 0 {
		taxType = TaxationTypeNoTax
	}

	return &Cashbill{
		MgtKey:            mgtKey,
		Tax:               strconv.Itoa(int(trans.VAT)),
		SupplyCost:        strconv.Itoa(int(trans.Supply)),
		TotalAmount:       strconv.Itoa(int(trans.Sum)),
		ServiceFee:        strconv.Itoa(int(customer.ServiceFee)),
		FranchiseCorpNum:  info.CorpNum,
		FranchiseAddr:     info.Addr,
		FranchiseCEOName:  info.CEOName,
		FranchiseCorpName: info.CorpOrCeoName(),
		FranchiseTEL:      aefire.LocalizePhoneNumber(info.TEL, 82),
		TaxationType:      taxType,
		TradeType:         TradeTypeApproval,
		TradeUsage:        customer.Usage,
		TradeOpt:          customer.TradeOpt,
		IdentityNum:       customer.IdentityNum,
		CustomerName:      customer.Name,
		Email:             customer.Email,
		OrderNumber:       mgtKey,
		ItemName:          customer.ItemName,
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

