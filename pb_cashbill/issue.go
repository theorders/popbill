package pb_cashbill

import (
	"errors"
	"github.com/theorders/aefire"
	"github.com/theorders/popbill"
	"net/http"
)

type Issue struct {
	Customer

	CorpNum   string    `json:"corpNum" firestore:"corpNum"`
	TradeType TradeType `json:"tradeType" firestore:"tradeType"`

	//거래관련
	MgtKey     string `json:"mgtKey" firestore:"mgtKey"`
	SupplyCost string `json:"supplyCost" firestore:"supplyCost"`
	Tax        string `json:"tax" firestore:"tax"`
	ServiceFee string `json:"serviceFee" firestore:"serviceFee"`

	//발행업첻관련
	FranchiseCorpNum  string `json:"franchiseCorpNum" firestore:"franchiseCorpNum"`
	FranchiseCorpName string `json:"franchiseCorpName" firestore:"franchiseCorpName"`
	FranchiseCEOName  string `json:"franchiseCEOName" firestore:"franchiseCEOName"`
	FranchiseAddr     string `json:"franchiseAddr" firestore:"franchiseAddr"`
	FranchiseTEL      string `json:"franchiseTEL" firestore:"franchiseTEL"`

	//취소발행 관련
	OrgMgtKey     string `json:"orgMgtKey,omitempty" firestore:"orgMgtKey,omitempty"` //(필수)파트너 문서관리번호
	OrgConfirmNum string `json:"orgConfirmNum,omitempty" firestore:"orgConfirmNum,omitempty"`
	OrgTradeDate  string `json:"orgTradeDate,omitempty" firestore:"orgTradeDate,omitempty"`
}

func (i *Issue) Validate() error {
	if i.TradeUsage == "" {
		return errors.New("발급용도가 지정되지 않았습니다")
	}

	if i.IdentityNum == "" {
		return errors.New("고객 식별번호가 없습니다")
	}

	if i.TotalAmount == "" {
		return errors.New("거래금액이 없습니다")
	}

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
	_, err := pb.MethodOverrideRequest(http.MethodPost,
		popbill.CashbillService,
		"",
		i,
		"ISSUE")

	return err
}

func (i *Issue) Cancel(pb *popbill.Client) error {
	_, err := pb.MethodOverrideRequest(
		http.MethodPost,
		popbill.CashbillService,
		i.MgtKey,
		nil,
		"CANCELISSUE")

	return err
}

func (i *Issue) Delete(pb *popbill.Client) error {
	_, err := pb.MethodOverrideRequest(
		http.MethodPost,
		popbill.CashbillService,
		i.MgtKey,
		nil,
		"DELETE")

	return err
}

func (i *Issue) CancelAndDelete(pb *popbill.Client) error {
	if err := i.Cancel(pb); err != nil {
		return err
	}

	if err := i.Delete(pb); err != nil {
		return err
	}

	return nil
}

func (i *Issue) Info(pb *popbill.Client) (b *Cashbill, err error) {
	res, err := pb.Request(
		http.MethodGet,
		popbill.CashbillService,
		i.MgtKey,
		nil,
	)

	if err != nil {
		return nil, err
	}

	b = &Cashbill{}
	if err := res.ToJSON(b); err != nil {
		return nil, aefire.NewHttpError(500, err)
	}

	return b, err
}

func (i *Issue) Detail(pb *popbill.Client) (b *Cashbill, err error) {
	res, err := pb.Request(
		http.MethodGet,
		popbill.CashbillService,
		(i.MgtKey)+"?Detail",
		nil,
	)

	if err != nil {
		return nil, err
	}

	b = &Cashbill{}
	if err := res.ToJSON(b); err != nil {
		return nil, aefire.NewHttpError(500, err)
	}

	return b, err
}
