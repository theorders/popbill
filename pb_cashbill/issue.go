package pb_cashbill

import (
	"github.com/theorders/aefire"
	"github.com/theorders/popbill"
	"net/http"
)

type Issue struct {
	Customer

	CorpNm    string    `json:"corpNum" firestore:"corpNum"`
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
	OrgMgtKey     string      `json:"orgMgtKey,omitempty" firestore:"orgMgtKey,omitempty"` //(필수)파트너 문서관리번호
	OrgConfirmNum string      `json:"orgConfirmNum,omitempty" firestore:"orgConfirmNum,omitempty"`
	OrgTradeDate  string      `json:"orgTradeDate,omitempty" firestore:"orgTradeDate,omitempty"`
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
	if err := res.ToJSON(b) ; err != nil{
		return nil, aefire.NewHttpError(500, err)
	}

	return b, err
}


func (i *Issue) Detail(pb *popbill.Client) (b *Cashbill, err error) {
	res, err := pb.Request(
		http.MethodGet,
		popbill.CashbillService,
		(i.MgtKey) + "?Detail",
		nil,
	)

	if err != nil {
		return nil, err
	}

	b = &Cashbill{}
	if err := res.ToJSON(b) ; err != nil{
		return nil, aefire.NewHttpError(500, err)
	}

	return b, err
}


