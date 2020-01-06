package popbill

import (
	"github.com/theorders/aefire"
)

type Transaction struct {
	Supply int64 `json:"supply" firestore:"supply"`
	VAT    int64 `json:"vat" firestore:"vat"`
	Sum    int64 `json:"sum" firestore:"sum"`
}

func (trans *Transaction) Add(target *Transaction) {
	trans.Supply += target.Supply
	trans.VAT += target.VAT
	trans.Sum += target.Sum
}

func (trans *Transaction) Sub(target *Transaction) {
	trans.Supply -= target.Supply
	trans.VAT -= target.VAT
	trans.Sum -= target.Sum
}


func NewFromSupply(supply int64, taxType TaxType) *Transaction {
	t := &Transaction{
		Supply:supply,
	}

	if taxType == TaxTypeNormal {
		t.VAT = int64(float64(supply) * 0.1)
	}

	t.Sum = t.Supply + t.VAT

	return t
}

func NewFromSum(sum int64, taxType TaxType) *Transaction {
	t := &Transaction{
		Sum: sum,
	}

	if taxType == TaxTypeNormal {
		t.Supply = int64(aefire.Round(float64(sum)/1.1, 0.5))
		t.VAT = t.Sum - t.Supply
	} else {
		t.Supply = t.Sum
	}

	return t
}
