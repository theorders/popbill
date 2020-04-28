package popbill

import (
	"errors"
	"github.com/theorders/aefire"
	"strings"
)

type CashbillIssueRequest struct {
	IdentityNum string `json:"identityNum" firestore:"identityNum"`
	Name        string `json:"name" firestore:"name"`
	Email       string `json:"email" firestore:"email"`
	ItemName    string `json:"itemName" firestore:"itemName"`
	TotalAmount int64  `json:"totalAmount" firestore:"totalAmount"`

	TradeUsage   `json:"usage" firestore:"usage"` //소득공제용, 지출증빙용
	TradeOpt     `json:"tradeOpt,omitempty" firestore:"tradeOpt"` //일반, 도서공연, 대중교통
	TaxationType `json:"taxationType" json:"taxationType"` //과세, 비과세
}

func (r *CashbillIssueRequest) NameOrIdentityNumMasked() string {
	if r.IdentityNum == SelfIssueNum {
		return SelfIssueUsage
	} else if r.Name != "" {
		return r.Name
	} else {
		return r.IdentityNumMasked()
	}
}
func (r *CashbillIssueRequest) IdentityNumMasked() string {
	idNum := strings.Replace(r.IdentityNum, "-", "", -1)

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

func (r *CashbillIssueRequest) Validate() error {
	if r.TradeUsage == "" {
		return errors.New("발급용도가 지정되지 않았습니다")
	}

	if r.IdentityNum == "" {
		return errors.New("고객 식별번호가 없습니다")
	}

	if r.TotalAmount == 0 {
		return errors.New("거래금액이 없습니다")
	}

	return nil
}
