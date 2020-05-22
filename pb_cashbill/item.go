package pb_cashbill

type Item struct {
	TaxationType TaxationType `json:"taxationType,omitempty" firestore:"taxationType,omitempty"` //(필수)
	TradeOpt     TradeOpt     `json:"tradeOpt,omitempty" firestore:"tradeOpt,omitempty"`         //(필수)
	TotalAmount  string       `json:"totalAmount,omitempty" firestore:"totalAmount,omitempty"`   //(필수)거래금액
	ItemName     string       `json:"itemName,omitempty" firestore:"itemName,omitempty"`
}

