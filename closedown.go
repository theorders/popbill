package popbill

import (
	"errors"
	"github.com/theorders/aefire"
	"net/http"
)

type TaxType string
type State string

const (
	TaxTypeNormal    TaxType = "1" //부가가치세 일반과세자
	TaxTypeFree      TaxType = "2" //부가가치세 면세과세자
	TaxTypeSimple    TaxType = "3" //부가가치세 간이과세자
	TaxTypeNonProfit TaxType = "4" //비영리법인 또는 국가기관, 고유번호가 부여된 단체

	StateUnregistered State = "0" //0 : 미등록 (등록되지 않은 사업자번호)
	StateInBusiness   State = "1" //1 : 사업중
	StateClose        State = "2" //2 : 폐업
	StateDown         State = "3" //3 : 휴업
)

type CloseDown struct {
	CorpNum   string  `json:"corpNum" firestore:"corpNum"`
	CheckDate string  `json:"checkDate" firestore:"checkDate"`
	StateDate string  `json:"stateDate" firestore:"stateDate"`
	TypeDate  string  `json:"typeDate" firestore:"typeDate"`
	TaxType   TaxType `json:"type" firestore:"type"`
	State     State   `json:"state" firestore:"state"`
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
