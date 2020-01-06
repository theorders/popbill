package popbill

import (
	"fmt"
	"net/http"
)

type Status struct {
	Code             int    `json:"code"`
	Message          string `json:"message"`
	InterOPYN        bool   `json:"interOPYN"`
	InvoiceeCorpNum  string `json:"invoiceeCorpNum"`
	ItemKey          string `json:"itemKey"`
	SupplyCostTotal  int    `json:"supplyCostTotal,string"`
	TaxTotal         int    `json:"taxTotal,string"`
	RegDT            string `json:"regDT"`
	WriteDate        string `json:"writeDate"`
	IssueType        string `json:"issueType"`
	StateDT          string `json:"stateDT"`
	OpenYN           bool   `json:"openYN"`
	TaxType          string `json:"taxType"`
	InvoicerCorpNum  string `json:"invoicerCorpNum"`
	InvoicerMgtKey   string `json:"invoicerMgtKey"`
	InvoicerCorpName string `json:"invoicerCorpName"`
	InvoicerPrintYN  bool   `json:"invoicerPrintYN"`
	InvoiceeCorpName string `json:"invoiceeCorpName"`
	InvoiceePrintYN  bool   `json:"invoiceePrintYN"`
	TrusteePrintYN   bool   `json:"trusteePrintYN"`
	PurposeType      string `json:"purposeType"`
	LateIssueYN      bool   `json:"lateIssueYN"`
	StateCode        int    `json:"stateCode"`
}

func (res *Status) GetCode() int {
	return res.Code
}

func (res *Status) GetMessage() string {
	return res.Message
}

func (res *Status) Error() string {
	return res.Message
}


func (c *Client) CheckMgtKeyInUse(mgtKeyType, mgtKey string) (invoiceStatus *Status, err error) {
	res, err := c.Request(
		http.MethodGet,
		TaxinvoiceService,
		fmt.Sprintf("%s/%s", mgtKeyType, mgtKey),
		"")

	if err != nil {
		return nil, err
	}

	invoiceStatus = &Status{}

	err = res.ToJSON(invoiceStatus)

	return invoiceStatus, err
}

