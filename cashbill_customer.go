package popbill

type CashbillCustomer struct {
	CashbillItem
	
	TradeUsage TradeUsage   `json:"tradeUsage" firestore:"tradeUsage"`
	IdentityNum string  `json:"identityNum" firestore:"identityNum"`
	CustomerName string `json:"customerName" firestore:"customerName"`
}
