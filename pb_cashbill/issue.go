package pb_cashbill

import (
	"errors"
	"github.com/theorders/aefire"
	"github.com/theorders/popbill"
	"net/http"
	"strconv"
	"strings"
)

type Issue struct {
	Customer

	TradeType TradeType `json:"tradeType,omitempty" firestore:"tradeType,omitempty"`

	//거래관련
	MgtKey      string `json:"mgtKey,omitempty" firestore:"mgtKey,omitempty"`
	SupplyCost  string `json:"supplyCost,omitempty" firestore:"supplyCost,omitempty"`
	Tax         string `json:"tax,omitempty" firestore:"tax,omitempty"`
	ServiceFee  string `json:"serviceFee,omitempty" firestore:"serviceFee,omitempty"`
	OrderNumber string `json:"orderNumber,omitempty" firestore:"orderNumber,omitempty"`

	//발행업첻관련
	FranchiseCorpNum  string `json:"franchiseCorpNum,omitempty" firestore:"franchiseCorpNum,omitempty"`
	FranchiseCorpName string `json:"franchiseCorpName,omitempty" firestore:"franchiseCorpName,omitempty"`
	FranchiseCEOName  string `json:"franchiseCEOName,omitempty" firestore:"franchiseCEOName,omitempty"`
	FranchiseAddr     string `json:"franchiseAddr,omitempty" firestore:"franchiseAddr,omitempty"`
	FranchiseTEL      string `json:"franchiseTEL,omitempty" firestore:"franchiseTEL,omitempty"`

	//취소발행 관련
	CancelType    CancelType `json:"cancelType,omitempty" firestore:"cancelType,omitempty"`
	OrgMgtKey     string     `json:"orgMgtKey,omitempty" firestore:"orgMgtKey,omitempty"` //(필수)파트너 문서관리번호
	OrgConfirmNum string     `json:"orgConfirmNum,omitempty" firestore:"orgConfirmNum,omitempty"`
	OrgTradeDate  string     `json:"orgTradeDate,omitempty" firestore:"orgTradeDate,omitempty"`
}

func (i *Issue) Validate(cd *popbill.CloseDown) error {
	if i.TradeUsage == "" {
		return errors.New("발급용도가 지정되지 않았습니다")
	}

	i.IdentityNum = strings.ReplaceAll(i.IdentityNum, "-", "")
	i.IdentityNum = strings.ReplaceAll(i.IdentityNum, " ", "")
	i.IdentityNum = strings.ReplaceAll(i.IdentityNum, ",", "")

	if i.IdentityNum == "" {
		return errors.New("고객 식별번호가 없습니다")
	}

	i.TotalAmount = strings.TrimSpace(i.TotalAmount)
	i.TotalAmount = strings.ReplaceAll(i.TotalAmount, ",", "")

	if i.TotalAmount == "" {
		return errors.New("거래금액이 없습니다")
	}

	totalAmount, err := strconv.Atoi(i.TotalAmount)
	if err != nil {
		return errors.New("거래금액이 숫자가 아닙니다")
	}

	if cd.TaxType != popbill.TaxTypeNormal {
		i.TaxationType = TaxationTypeNoTax
	}

	if i.TaxationType == TaxationTypeWithTax && totalAmount > 10 {
		var supply, tax int
		supply = (totalAmount / 11) * 10
		tax = totalAmount - supply

		i.SupplyCost = strconv.Itoa(supply)
		i.Tax = strconv.Itoa(tax)
	} else {
		i.SupplyCost = i.TotalAmount
		i.Tax = "0"
	}

	i.ServiceFee = "0"

	//{@no.7 tradeUsage} 값이 "소득공제용" 인 경우
	//└ 주민등록/휴대폰/카드번호(현금영수증 카드)/자진발급용 번호(010-000-1234) 입력
	//{@no.7 tradeUsage} 값이 "지출증빙용" 인 경우
	//└ 사업자번호/주민등록/휴대폰/카드번호(현금영수증 카드) 입력
	//※ 주민등록번호 13자리, 휴대폰번호 10~11자리, 카드번호 13~19자리, 사업자번호 10자리 입력 가능
	//소득공제용
	if i.TradeUsage == TradeUsageIncomeDeduction &&
		!aefire.ValidateRRN(i.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(i.IdentityNum) &&
		(len(i.IdentityNum) < 13 || len(i.IdentityNum) > 19) {
		return errors.New("소득공제용 현금영수증 발급대상고객의 주민등록번호, 휴대전화번호 혹은 카드번호가 필요합니다")
	}

	//지출증빙용
	if i.TradeUsage == TradeUsageProofOfExpenditure &&
		!aefire.ValidateRRN(i.IdentityNum) &&
		!aefire.ValidateLocalCellPhoneNumber(i.IdentityNum) &&
		!aefire.ValidateCorpNum(i.IdentityNum) &&
		(len(i.IdentityNum) < 13 || len(i.IdentityNum) > 19) {
		return errors.New("지출증빙용 현금영수증 발급대상고객의 사업자등록번호, 휴대전화번호, 주민등록번호 혹은 카드번호가 필요합니다")
	}

	return nil
}

func (i *Issue) Regist(pb *popbill.Client) error {
	i.TradeType = TradeTypeApproval

	_, err := pb.MethodOverrideRequest(http.MethodPost,
		popbill.CashbillService,
		"",
		i,
		"ISSUE")

	return err
}
