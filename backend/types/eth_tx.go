package types

type ETHTransaction struct {
	Hash string
	BlockNumber int
	From string
	To string
	Value string
}

func (ethtx *ETHTransaction) ToTransaction() *Transaction {
	return &Transaction{
		Txhash: ethtx.Hash,
		FromAddress: ethtx.From,
		TxOutputs: []TxOutput{TxOutput{ToAddress: ethtx.To, Value: ethtx.Value}},
		Timestamp: int64(ethtx.BlockNumber) * 15000 + 1356969600000,
	}
}

func ETHTransactionToTransaction (ethtxs []ETHTransaction) (txs []Transaction) {
	for _, ethtx := range ethtxs {
		txs = append(txs, *ethtx.ToTransaction())
	}
	return
}
