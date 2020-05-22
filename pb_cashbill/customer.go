package pb_cashbill

type Customer struct {
	Item

	TradeUsage   TradeUsage `json:"tradeUsage,omitempty" firestore:"tradeUsage,omitempty"`
	IdentityNum  string     `json:"identityNum,omitempty" firestore:"identityNum,omitempty"`
	CustomerName string     `json:"customerName,omitempty" firestore:"customerName,omitempty"`
	Email        string     `json:"email,omitempty" firestore:"email,omitempty"`
}


func (c *Customer) CustomerNameOrIdentityNum() string {
	if c.CustomerName != ""{
		return c.CustomerName
	} else {
		return c.IdentityNum
	}
}