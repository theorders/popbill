package pb_cashbill

import (
	"github.com/labstack/echo/v4"
	"github.com/theorders/aefire"
	"github.com/theorders/popbill"
	"net/http"
	"strconv"
	"strings"
)

type TradeUsage string
type TradeType string
type TradeOpt string
type TaxationType string
type WebHookEventType string
type CancelType int

const (
	TradeUsageIncomeDeduction    TradeUsage = "소득공제용"
	TradeUsageProofOfExpenditure TradeUsage = "지출증빙용"

	TradeTypeApproval TradeType = "승인거래"
	TradeTypeCancel   TradeType = "취소거래"

	TaxationTypeWithTax TaxationType = "과세"
	TaxationTypeNoTax   TaxationType = "비과세"

	TradeOptN TradeOpt = "일반"
	TradeOptB TradeOpt = "도서공연"
	TradeOptT TradeOpt = "대중교통"

	WebHookEventTypeIssue  WebHookEventType = "Issue"
	WebHookEventTypeCancel WebHookEventType = "Cancel"
	WebHookEventTypeNTS    WebHookEventType = "NTS"

	CancelTypeTrade CancelType = 1
	CancelTypeError CancelType = 2
	CancelTypeEtc   CancelType = 3

	SelfIssueNum   = "0100001234"
	SelfIssueUsage = "자진발급"
)

type Cashbill struct {
	Issue

	ConfirmNum string `json:"confirmNum,omitempty" firestore:"confirmNum"`
	TradeDate  string `json:"tradeDate,omitempty" firestore:"tradeDate"`

	//팝빌 값들
	ItemKey   string `json:"itemKey,omitempty" firestore:"itemKey"`
	StateMemo string `json:"stateMemo,omitempty" firestore:"stateMemo"`
	StateCode int64  `json:"stateCode,omitempty" firestore:"stateCode"`
	StateDT   string `json:"stateDT,omitempty" firestore:"stateDT"`

	//국세청 값들
	NtssendDT        string `json:"ntssendDT,omitempty" firestore:"ntssendDT"`
	NtsresultDT      string `json:"ntsresultDT,omitempty" firestore:"ntsresultDT"`
	NtsresultCode    string `json:"ntsresultCode,omitempty" firestore:"ntsresultCode"`
	NtsresultMessage string `json:"ntsresultMessage,omitempty" firestore:"ntsresultMessage"`
}

type EventMessage struct {
	CorpNum          string           `json:"corpNum"`
	EventType        WebHookEventType `json:"eventType"`
	EventDT          string           `json:"eventDT"`
	MgtKey           string           `json:"mgtKey"`
	StateMemo        string           `json:"stateMemo"`
	StateCode        int64            `json:"stateCode"`
	StateDT          string           `json:"stateDT"`
	ConfirmNum       string           `json:"confirmNum"`
	NtsresultCode    string           `json:"ntsresultCode"`
	NtssendDT        string           `json:"ntssendDT"`
	NtsresultDT      string           `json:"ntsresultDT"`
	NtsresultMessage string           `json:"ntsresultMessage"`
	ItemKey          string           `json:"itemKey"`
}

func (b *Cashbill) IdentityNumMasked() string {
	idNum := strings.Replace(b.IdentityNum, "-", "", -1)

	if idNum == SelfIssueNum {
		return SelfIssueUsage
	} else if aefire.ValidateCorpNum(idNum) {
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

func (b *Cashbill) TaxValue() int64 {
	if b.Tax == "" {
		return 0
	} else {
		i, _ := strconv.ParseInt(b.Tax, 10, 64)
		return i
	}
}

func (b *Cashbill) SupplyCostValue() int64 {
	if b.SupplyCost == "" {
		return 0
	} else {
		i, _ := strconv.ParseInt(b.SupplyCost, 10, 64)
		return i
	}
}

func (b *Cashbill) TotalAmountValue() int64 {
	if b.TotalAmount == "" {
		return 0
	} else {
		i, _ := strconv.ParseInt(b.TotalAmount, 10, 64)
		return i
	}
}

func (b *Cashbill) Revoke(pb *popbill.Client, mgtKey string) (*Cashbill, error) {
	if b.ConfirmNum == "" || b.TradeDate == "" {
		return nil, echo.NewHTTPError(http.StatusPreconditionFailed, "국세청 승인번호와 거래일자가 확정된 영수증만 취소할 수 있습니다")
	}

	var revoked Cashbill

	revoked = *b
	revoked.TradeType = TradeTypeCancel


	revoked.OrgMgtKey = b.MgtKey
	revoked.OrgConfirmNum = b.ConfirmNum
	revoked.OrgTradeDate = b.TradeDate

	revoked.MgtKey = mgtKey
	revoked.ConfirmNum = ""
	revoked.TradeDate = ""
	revoked.StateCode = 100

	_, err := pb.MethodOverrideRequest(http.MethodPost,
		popbill.CashbillService,
		"",
		&revoked,
		"REVOKEISSUE")

	//now := time.Now()
	//revoke.CreatedAt = &now

	return &revoked, err
}





func Cancel(pb *popbill.Client, mgtKey string) error {
	_, err := pb.MethodOverrideRequest(
		http.MethodPost,
		popbill.CashbillService,
		mgtKey,
		nil,
		"CANCELISSUE")

	return err
}

func Delete(pb *popbill.Client, mgtKey string) error {
	_, err := pb.MethodOverrideRequest(
		http.MethodPost,
		popbill.CashbillService,
		mgtKey,
		nil,
		"DELETE")

	return err
}

func  CancelAndDelete(pb *popbill.Client, mgtKey string) error {
	if err := Cancel(pb, mgtKey); err != nil {
		return err
	}

	if err := Delete(pb, mgtKey); err != nil {
		return err
	}

	return nil
}

func  Info(pb *popbill.Client, mgtKey string) (b *Cashbill, err error) {
	res, err := pb.Request(
		http.MethodGet,
		popbill.CashbillService,
		mgtKey,
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

func  Detail(pb *popbill.Client, mgtKey string) (b *Cashbill, err error) {
	res, err := pb.Request(
		http.MethodGet,
		popbill.CashbillService,
		(mgtKey)+"?Detail",
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
