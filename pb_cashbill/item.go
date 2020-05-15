package pb_cashbill

type Item struct {
	TaxationType TaxationType `json:"taxationType" firestore:"taxationType"` //(필수)
	TradeOpt     TradeOpt     `json:"tradeOpt" firestore:"tradeOpt"`         //(필수)
	TotalAmount  string       `json:"totalAmount" firestore:"totalAmount"`   //(필수)거래금액
	ItemName     string       `json:"itemName" firestore:"itemName"`
}

