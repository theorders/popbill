package popbill

import (
	"github.com/imroc/req"
	"github.com/labstack/echo/v4"
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
	Title *string       `json:"title,omitempty"`
	FCnt  int           `json:"fCnt"`
	SndDT *string       `json:"sndDT,omitempty"`
	AdsYN bool          `json:"adsYN"`
	Rcvs  []FaxReceiver `json:"rcvs"`
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
		SendResult     int      `json:"sendResult"`
		FileNames      []string `json:"fileNames"`
	} `json:"list"`
}

func (c *Client) FAXSend(sendRequest FaxSendRequest, files ...req.FileUpload) (*Receipt, error) {
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

	receipt := &Receipt{}
	if err := res.ToJSON(receipt) ; err != nil{
		return nil, echo.NewHTTPError(500, err)
	}

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
	if err := res.ToJSON(&resData) ; err != nil{
		return nil, echo.NewHTTPError(500, err)
	}

	return &resData, err
}
