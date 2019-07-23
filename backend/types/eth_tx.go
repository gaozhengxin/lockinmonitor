package types

import (
	"strings"
)

var Tokens map[string]string = map[string]string{
	"ERC20GUSD":"0x28a79f9b0fe54a39a0ff4c10feeefa832eeceb78",
	"ERC20BNB":"0x7f30b414a814a6326d38535ca8eb7b9a62bceae2",
	"ERC20MKR":"0x2c111ede2538400f39368f3a3f22a9ac90a496c7",
	"ERC20HT":"0x3c3d51f6be72b265fe5a5c6326648c4e204c8b9a",
	"ERC20BNT":"0x14d5913c8396d43ab979d4b29f2102c1c65e18db",
}

type ETHTransaction struct {
	Hash string
	BlockNumber int
	From string
	To string
	Value string
	Contractaddress string
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
		if ethtx.Contractaddress != "" {
			continue
		}
		txs = append(txs, *ethtx.ToTransaction())
	}
	return
}

func ERC20TransactionToTransaction (cointype string, erc20txs []ETHTransaction) (txs []Transaction) {
	for _, erc20tx := range erc20txs {
		if strings.EqualFold(erc20tx.Contractaddress, Tokens[cointype]) == false {
			continue
		}
		txs = append(txs, *erc20tx.ToTransaction())
	}
	return
}
