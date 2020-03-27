package popbill

import (
	"errors"
	"github.com/theorders/aefire"
	"net/http"
)

type TaxType string
type CloseDownState string

const (
	TaxTypeNormal    TaxType = "1" //부가가치세 일반과세자
	TaxTypeFree      TaxType = "2" //부가가치세 면세과세자
	TaxTypeSimple    TaxType = "3" //부가가치세 간이과세자
	TaxTypeNonProfit TaxType = "4" //비영리법인 또는 국가기관, 고유번호가 부여된 단체

	CloseDownStateUnregistered CloseDownState = "0" //0 : 미등록 (등록되지 않은 사업자번호)
	CloseDownStateInBusiness   CloseDownState = "1" //1 : 사업중
	CloseDownStateClose        CloseDownState = "2" //2 : 폐업
	CloseDownStateDown         CloseDownState = "3" //3 : 휴업
)

func (t *TaxType) ToWord() string {
	switch *t {
	case TaxTypeNormal:
		return "normal"
	case TaxTypeFree:
		return "free"
	case TaxTypeSimple:
		return "simple"
	case TaxTypeNonProfit:
		return "nonProfit"
	default:
		panic(*t + "is not popbill.TaxType")
	}
}

func WordToTaxType(w string) TaxType {
	switch w {
	case "normal":
		return TaxTypeNormal
	case "free":
		return TaxTypeFree
	case "simple":
		return TaxTypeSimple
	case "nonProfit":
		return TaxTypeNonProfit
	default:
		panic(w + "is not popbill.TaxType word")
	}
}

func (s *CloseDownState) ToWord() string {
	switch *s {
	case CloseDownStateUnregistered:
		return "unregistered"
	case CloseDownStateInBusiness:
		return "normal"
	case CloseDownStateClose:
		return "close"
	case CloseDownStateDown:
		return "down"
	default:
		panic(*s + "is not popbill.State")
	}
}

func WordToCloseDownState(w string) CloseDownState {
	switch w {
	case "unregistered":
		return CloseDownStateUnregistered
	case "normal":
		return CloseDownStateInBusiness
	case "close":
		return CloseDownStateClose
	case "down":
		return CloseDownStateDown
	default:
		panic(w + "is not popbill.CloseDownState word")
	}
}

type CloseDown struct {
	CorpNum   string  `json:"corpNum" firestore:"corpNum"`
	CheckDate string  `json:"checkDate" firestore:"checkDate"`
	StateDate string  `json:"stateDate" firestore:"stateDate"`
	TypeDate  string  `json:"typeDate" firestore:"typeDate"`
	TaxType   TaxType `json:"type" firestore:"type"`
	State     CloseDownState   `json:"state" firestore:"state"`
}

type WordCloseDown struct {
	CorpNum   string  `json:"corpNum" firestore:"corpNum"`
	CheckDate string  `json:"checkDate" firestore:"checkDate"`
	StateDate string  `json:"stateDate" firestore:"stateDate"`
	TypeDate  string  `json:"typeDate" firestore:"typeDate"`
	TaxType   string `json:"type" firestore:"type"`
	State     string   `json:"state" firestore:"state"`
}

func (cd *CloseDown) ToWordCloseDown() (wcd *WordCloseDown) {
	wcd = &WordCloseDown{}
	wcd.CorpNum = cd.CorpNum
	wcd.CheckDate = cd.CheckDate
	wcd.StateDate = cd.StateDate
	wcd.TypeDate = cd.TypeDate
	wcd.TaxType = cd.TaxType.ToWord()
	wcd.State = cd.State.ToWord()

	return
}

func (wcd *WordCloseDown) ToCloseDown() (cd *CloseDown) {
	cd = &CloseDown{}
	cd.CorpNum = wcd.CorpNum
	cd.CheckDate = wcd.CheckDate
	cd.StateDate = wcd.StateDate
	cd.TypeDate = wcd.TypeDate
	cd.TaxType = WordToTaxType(wcd.TaxType)
	cd.State = WordToCloseDownState(wcd.State)

	return
}

func (c *Client) GetCloseDown(cn string) (closeDown *CloseDown, err error) {
	res, err := c.Request(
		http.MethodGet,
		CloseDownService,
		"",
		aefire.UrlValuesOf("CN", cn))

	if aefire.LogIfError(err) {
		return nil, errors.New("홈택스 중계서버에 연결하지 못했습니다")
	}

	closeDown = &CloseDown{}

	if aefire.LogIfError(res.ToJSON(closeDown)) {
		return nil, errors.New("조회 결과를 해석하지 못했습니다")
	}

	return closeDown, err
}
