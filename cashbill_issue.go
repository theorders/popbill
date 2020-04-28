package popbill

type CashbillIssueRequest struct {
	IdentityNum string `json:"identityNum" firestore:"identityNum"`
	Name        string `json:"name" firestore:"name"`
	Email       string `json:"email" firestore:"email"`
	ItemName    string `json:"itemName" firestore:"itemName"`
	TotalAmount int64  `json:"totalAmount" firestore:"totalAmount"`

	TradeUsage   `json:"tradeUsage" firestore:"tradeUsage"` //소득공제용, 지출증빙용
	TradeOpt     `json:"tradeOpt,omitempty" firestore:"tradeOpt"` //일반, 도서공연, 대중교통
	TaxationType `json:"taxationType" json:"taxationType"` //과세, 비과세
}

func (r *CashbillIssueRequest) ToCustomer() (c *CashbillCustomer) {
	c = &CashbillCustomer{
		IdentityNum: r.IdentityNum,
		Name:        r.Name,
		Email:       r.Email,
		ItemName:    r.ItemName,
		Usage:       r.TradeUsage,
		TradeOpt:    r.TradeOpt,
		Transaction: nil,
		ServiceFee:  0,
		Aid:         nil,
	}

	if r.TaxationType == TaxationTypeWithTax {
		c.Transaction = TransactionFromSum(r.TotalAmount, TaxTypeNormal)
	} else{
		c.Transaction = TransactionFromSum(r.TotalAmount, TaxTypeFree)
	}

	return
}
