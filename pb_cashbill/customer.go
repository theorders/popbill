package pb_cashbill

type Customer struct {
	Item

	TradeUsage   TradeUsage `json:"tradeUsage" firestore:"tradeUsage"`
	IdentityNum  string     `json:"identityNum" firestore:"identityNum"`
	CustomerName string     `json:"customerName" firestore:"customerName"`
	Email        string     `json:"email,omitempty" firestore:"email,omitempty"`
}


func (c *Customer) CustomerNameOrIdentityNum() string {
	if c.CustomerName != ""{
		return c.CustomerName
	} else {
		return c.IdentityNum
	}
}