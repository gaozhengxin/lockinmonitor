package types

import "encoding/json"

type Transaction struct {
	Txhash string
	FromAddress string
	TxOutputs []TxOutput
	Timestamp string
}

type TxOutput struct {
	ToAddress string
	Value string
}

func ParseTransactions (cointype string) func ([]byte) []Transaction {
	if cointype == "EVT-1" {
		return func (b []byte) []Transaction {
			return EvtActionToTransaction(UnmarshalEvtActions(b))
		}
	}
	return nil
}

func MarshalTransactions (txs []Transaction) []byte {
	res := "["
	for i, tx := range txs {
		b, err := json.Marshal(tx)
		if err != nil {
			continue
		}
		res = res + string(b)
		if i + 1 < len(txs) {
			res += ","
		}
	}
	res += "]"
	return []byte(res)
}