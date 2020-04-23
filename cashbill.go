package popbill

import (
	"github.com/labstack/echo/v4"
	"github.com/theorders/aefire"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	MgtKey            string           `json:"mgtKey,omitempty" firestore:"mgtKey"` //(필수)파트너 문서관리번호
	ConfirmNum        string           `json:"confirmNum,omitempty" firestore:"confirmNum"`
	TradeDate         string           `json:"tradeDate,omitempty" firestore:"tradeDate"`
	TradeType         TradeType        `json:"tradeType,omitempty" firestore:"tradeType"`               //(필수)
	TradeUsage        TradeUsage       `json:"tradeUsage,omitempty" firestore:"tradeUsage"`             //(필수)
	TaxationType      TaxationType     `json:"taxationType,omitempty" firestore:"taxationType"`         //(필수)
	TradeOpt          TradeOpt         `json:"tradeOpt,omitempty" firestore:"tradeOpt"`                 //(필수)
	SupplyCost        string           `json:"supplyCost,omitempty" firestore:"supplyCost"`             //(필수)공급가액
	Tax               string           `json:"tax,omitempty" firestore:"tax"`                           //(필수)세액
	ServiceFee        string           `json:"serviceFee,omitempty" firestore:"serviceFee"`             //(필수)봉사료
	TotalAmount       string           `json:"totalAmount,omitempty" firestore:"totalAmount"`           //(필수)거래금액
	FranchiseCorpNum  string           `json:"franchiseCorpNum,omitempty" firestore:"franchiseCorpNum"` //(필수)발행자 사업자번호
	FranchiseCorpName string           `json:"franchiseCorpName,omitempty" firestore:"franchiseCorpName"`
	FranchiseCEOName  string           `json:"franchiseCEOName,omitempty" firestore:"franchiseCEOName"`
	FranchiseAddr     string           `json:"franchiseAddr,omitempty" firestore:"franchiseAddr"`
	FranchiseTEL      string           `json:"franchiseTEL,omitempty" firestore:"franchiseTEL"`
	IdentityNum       string           `json:"identityNum,omitempty" firestore:"identityNum"` //(필수)거래처 식별번호
	CustomerName      string           `json:"customerName,omitempty" firestore:"customerName"`
	ItemName          string           `json:"itemName,omitempty" firestore:"itemName"`
	ItemKey           string           `json:"itemKey,omitempty" firestore:"itemKey"`
	OrderNumber       string           `json:"orderNumber,omitempty" firestore:"orderNumber"`
	Email             string           `json:"email,omitempty" firestore:"email"`
	Hp                string           `json:"hp,omitempty" firestore:"hp"`
	StateMemo         string           `json:"stateMemo,omitempty" firestore:"stateMemo"`
	StateCode         int64            `json:"stateCode,omitempty" firestore:"stateCode"`
	StateDT           string           `json:"stateDT,omitempty" firestore:"stateDT"`
	OrgTradeDate      string           `json:"orgTradeDate,omitempty" firestore:"orgTradeDate"`
	NtssendDT         string           `json:"ntssendDT,omitempty" firestore:"ntssendDT"`
	NtsresultDT       string           `json:"ntsresultDT,omitempty" firestore:"ntsresultDT"`
	NtsresultCode     string           `json:"ntsresultCode,omitempty" firestore:"ntsresultCode"`
	NtsresultMessage  string           `json:"ntsresultMessage,omitempty" firestore:"ntsresultMessage"`
	EventType         WebHookEventType `json:"eventType,omitempty" firestore:"eventType"`
	EventDT           string           `json:"eventDT,omitempty" firestore:"eventDT"`
	CorpNum           string           `json:"corpNum,omitempty" firestore:"corpNum"`
	CreatedAt         *time.Time       `json:"createdAt,omitempty" firestore:"createdAt"`
	UID               string           `json:"uid,omitempty" firestore:"uid"`
	Memo              string           `json:"memo,omitempty" firestore:"memo"`
	NTSConfirmNotif   bool             `json:"ntsConfirmNotif,omitempty" firestore:"ntsConfirmNotif"`
	RevokeIssue       *RevokeIssue     `json:"revokeIssue,omitempty" firestore:"revokeIssue,omitempty"`
	Aid               *string          `json:"aid,omitempty" firestore:"aid,omitempty"`
}

type RevokeIssue struct {
	Cashbill
	OrgMgtKey     string     `json:"orgMgtKey,omitempty" firestore:"orgMgtKey"` //(필수)파트너 문서관리번호
	OrgConfirmNum string     `json:"orgConfirmNum,omitempty" firestore:"orgConfirmNum"`
	OrgTradeDate  string     `json:"orgTradeDate,omitempty" firestore:"orgTradeDate"`
	IsPartCancel  bool       `json:"isPartCancel,omitempty" firestore:"isPartCancel"`
	CancelType    CancelType `json:"cancelType" firestore:"cancelType"`
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

func (b *Cashbill) HasTaxationError(taxType TaxType) bool {
	if taxType == TaxTypeSimple {
		if b.TaxValue() == 0 || float64(b.TaxValue())/float64(b.TotalAmountValue()) >= 0.04 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (c *Client) CashbillIssue(cashbill *Cashbill) *echo.HTTPError {
	_, err := c.MethodOverrideRequest(http.MethodPost,
		CashbillService,
		"",
		cashbill,
		"ISSUE")

	now := time.Now()
	cashbill.CreatedAt = &now

	return err
}

func (c *Client) CashbillRevokeIssue(revoke *RevokeIssue) *echo.HTTPError {
	_, err := c.MethodOverrideRequest(http.MethodPost,
		CashbillService,
		"",
		revoke,
		"REVOKEISSUE")

	now := time.Now()
	revoke.CreatedAt = &now

	return err
}

func (c *Client) CashbillCancel(mgtKey string) *echo.HTTPError {
	_, err := c.MethodOverrideRequest(
		http.MethodPost,
		CashbillService,
		mgtKey,
		nil,
		"CANCELISSUE")

	return err
}

func (c *Client) GetCashbillInfo(mgtKey string) (cashbill *Cashbill, err *echo.HTTPError) {
	res, err := c.Request(
		http.MethodGet,
		CashbillService,
		mgtKey,
		nil,
	)

	if err != nil {
		return nil, err
	}

	cashbill = &Cashbill{}
	if err := res.ToJSON(cashbill) ; err != nil{
		return nil, aefire.NewEchoHttpError(500, err)
	}

	return cashbill, err
}

func (c *Client) GetCashbillDetail(mgtKey string) (m map[string]interface{}, err *echo.HTTPError) {
	res, err := c.Request(
		http.MethodGet,
		CashbillService,
		mgtKey+"?Detail",
		nil,
	)

	if err != nil {
		return nil, err
	}

	m = aefire.MapOf()
	if err := res.ToJSON(&m) ; err != nil{
		return nil, aefire.NewEchoHttpError(500, err)
	}

	delete(m, "smssendYN")

	return m, nil
}
