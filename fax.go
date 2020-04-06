package popbill

import (
	"github.com/imroc/req"
	"github.com/theorders/aefire"
	"net/http"
)

type FaxReceiver struct {
	Rcv   string `json:"rcv"`
	Rcvnm string `json:"rcvnm"`
}

type FaxSendRequest struct {
	Snd   string        `json:"snd"`
	SndNm *string       `json:"snnm,omitempty"`
	//SndDT *string       `json:"sndDT,omitempty"`
	Rcvs  []FaxReceiver `json:"rcvs"`
	FCnt  int           `json:"fCnt"`
}

type FaxListResponse struct {
	Code      int    `json:"code"`
	Total     int    `json:"total"`
	PerPage   int    `json:"perPage"`
	PageNum   int    `json:"pageNum"`
	PageCount int    `json:"pageCount"`
	Message   string `json:"message"`
	List      []struct {
		SendState      int      `json:"sendState"`
		ConvState      int      `json:"convState"`
		SendNum        string   `json:"sendNum"`
		ReceiveNum     string   `json:"receiveNum"`
		ReceiveName    string   `json:"receiveName"`
		SendPageCnt    int      `json:"sendPageCnt"`
		SuccessPageCnt int      `json:"successPageCnt"`
		FailPageCnt    int      `json:"failPageCnt"`
		RefundPageCnt  int      `json:"refundPageCnt"`
		CancelPageCnt  int      `json:"cancelPageCnt"`
		ReserveDT      string   `json:"reserveDT"`
		ReceiptDT      string   `json:"receiptDT"`
		SendDT         string   `json:"sendDT"`
		ResultDT       string   `json:"resultDT"`
		SendResult     int   `json:"sendResult"`
		FileNames      []string `json:"fileNames"`
	} `json:"list"`
}

func (c *Client) FAXSend(sendRequest FaxSendRequest, files ...req.FileUpload) (receipt *Receipt, err error) {
	for i, r := range sendRequest.Rcvs {
		sendRequest.Rcvs[i].Rcvnm = aefire.LocalizePhoneNumber(r.Rcvnm, 82)
	}

	sendRequest.FCnt = len(files)

	res, err := c.MultipartFormDataRequest(
		FAXService,
		"",
		req.Param{"form": aefire.ToJson(sendRequest)},
		files...)

	if err != nil {
		return nil, err
	}

	receipt = &Receipt{}
	err = res.ToJSON(receipt)

	return receipt, err
}

func (c *Client) FaxList(SDate, EDate string, Page, PerPage int) (*FaxListResponse, error) {
	res, err := c.Request(http.MethodGet,
		FAXService,
		"/Search",
		aefire.MapOf(
			"SDate", SDate,
			"EDate", EDate,
			"Page", Page,
			"PerPage", PerPage), )

	if err != nil {
		return nil, err
	}

	resData := FaxListResponse{}
	err = res.ToJSON(&resData)

	return &resData, err
}
