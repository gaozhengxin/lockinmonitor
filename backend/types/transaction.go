package types

import (
	"encoding/json"
	//"strings"
)

type Transaction struct {
	Txhash string
	FromAddress string
	TxOutputs []TxOutput
	Timestamp int64
	//Timestamp string
}

type TxOutput struct {
	ToAddress string
	Value string
}

func ParseEVTTransactions (b []byte) []Transaction {
	return EvtActionToTransaction(UnmarshalEvtActions(b))
}

func ParseBTCTransactions (b []byte) []Transaction {
	var btctxs []BtcTx
	json.Unmarshal(b,&btctxs)
	if btctxs == nil {
		return nil
	}
	return BtcTxToTransaction(btctxs)
}

func ParseETHTransactions (b []byte) []Transaction {
	var ethtxs []ETHTransaction
	json.Unmarshal(b,&ethtxs)
	if ethtxs == nil {
		return nil
	}
	return ETHTransactionToTransaction(ethtxs)
}

func ParseERC20Transactions (cointype string, b []byte) []Transaction {
	var erc20txs []ETHTransaction
	json.Unmarshal(b,&erc20txs)
	if erc20txs == nil {
		return nil
	}
	return ERC20TransactionToTransaction(cointype, erc20txs)
}

/*
func ParseTransactions (cointype string) func ([]byte) []Transaction {
	if strings.HasPrefix(cointype, "EVT") {
		return func (b []byte) []Transaction {
			return EvtActionToTransaction(UnmarshalEvtActions(b))
		}
	}
	if strings.EqualFold(cointype, "BTC") {
		return func (b []byte) []Transaction {
			var btctxs []BtcTx
			json.Unmarshal(b,&btctxs)
			if btctxs == nil {
				return nil
			}
			return BtcTxToTransaction(btctxs)
		}
	}
	if strings.EqualFold(cointype, "ETH") {
		return func (b []byte) []Transaction {
			var ethtxs []ETHTransaction
			json.Unmarshal(b,&ethtxs)
			if ethtxs == nil {
				return nil
			}
			return ETHTransactionToTransaction(ethtxs)
		}
	}
	return nil
}
*/

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
