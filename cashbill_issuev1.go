package popbill

import (
	"errors"
	"github.com/theorders/aefire"
	"strings"
)

type CashbillIssueV1 struct {
	IdentityNum  string       `json:"identityNum" firestore:"identityNum"`
	Email        string       `json:"email" firestore:"email"`
	ItemName     string       `json:"itemName" firestore:"itemName"`
	CustomerName string       `json:"customerName" firestore:"customerName"`
	Usage        TradeUsage   `json:"usage" firestore:"usage"`
	TradeOpt     TradeOpt     `json:"tradeOpt,omitempty" firestore:"tradeOpt"` //(필수)
	Transaction  *Transaction `json:"transaction" firestore:"transaction"`
	ServiceFee   int64        `json:"serviceFee" firestore:"serviceFee"`
	Aid          *string      `json:"aid,omitempty" firestore:"aid,omitempty"`
}

func (cashbill *CashbillIssueV1) NameOrIdentityNumMasked() string {
	if cashbill.IdentityNum == SelfIssueNum {
		return SelfIssueUsage
	} else if cashbill.CustomerName != "" {
		return cashbill.CustomerName
	} else {
		return cashbill.IdentityNumMasked()
	}
}
func (cashbill *CashbillIssueV1) IdentityNumMasked() string {
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

func (customer *CashbillIssueV1) Validate() error {
	if customer.Usage == "" {
		return errors.New("발급용도가 지정되지 않았습니다")
	}

	if customer.IdentityNum == "" {
		return errors.New("고객 식별번호가 없습니다")
	}

	if customer.Transaction == nil || customer.Transaction.Supply == 0 || customer.Transaction.Sum == 0 {
		return errors.New("거래금액이 없습니다")
	}

	//소득공제용
	if customer.Usage == TradeUsageIncomeDeduction &&
		!aefire.ValidateRRN(customer.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) {
		return errors.New("소득공제용 현금영수증 발행에는 고객 주민등록번호나 휴대전화번호가 필요합니다")
	}
	//if customer.Usage == TradeUsageIncomeDeduction &&
	//	!aefire.ValidateRRN(customer.IdentityNum) &&
	//	!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) {
	//	return errors.New("소득공제용 현금영수증 발행에는 고객 주민등록번호나 휴대전화번호가 필요합니다")
	//}

	//지출증빙용
	if customer.Usage == TradeUsageProofOfExpenditure &&
		!aefire.ValidateRRN(customer.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) &&
		!aefire.ValidateCorpNum(customer.IdentityNum) {
		return errors.New("소득공제용 현금영수증 발행에는 사업자등록번호, 주민등록번호 혹은 휴대전화번호가 필요합니다")
	}
	//if customer.Usage == TradeUsageProofOfExpenditure &&
	//	!aefire.ValidateRRN(customer.IdentityNum) &&
	//	!aefire.ValidateLocalCellPhoneNumber(customer.IdentityNum) &&
	//	!aefire.ValidateCorpNum(customer.IdentityNum) {
	//	return errors.New("소득공제용 현금영수증 발행에는 사업자등록번호, 주민등록번호 혹은 휴대전화번호가 필요합니다")
	//}

	return nil
}

