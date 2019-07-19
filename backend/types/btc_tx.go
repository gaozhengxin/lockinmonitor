package types

import (
	"strconv"
	"strings"
//	"time"
	"github.com/gaozhengxin/lockinmonitor/backend/utils"
)

type BtcTx struct {
	Txid string
	Vin []VinData
	Vout []VoutData
	Status TxStatus
}

type TxStatus struct {
	Block_time float64
}

type VinData struct {
	Txid string
	Vout int
	Prevout VoutData
}

type VoutData struct {
	Scriptpubkey string
	Scriptpubkey_address string
	Value float64
}

func BtcTxToTransaction(btxs []BtcTx) []Transaction {
	btxs = ConvertAddresses(btxs)
	var txs = make([]Transaction,len(btxs))
	for i, tx := range btxs {
		txs[i].Txhash = tx.Txid
		var froms = make(map[string]int)
		for _, vin := range tx.Vin {
			froms[vin.Prevout.Scriptpubkey_address]=1
		}
		fromStr := ""
		for k, _ := range froms {
			fromStr += k + "|"
		}
		txs[i].FromAddress = strings.TrimSuffix(fromStr,"|")
		for _, vout := range tx.Vout {
			if froms[vout.Scriptpubkey_address] == 1 {
				// 去掉找回
				continue
			}
			txs[i].TxOutputs = append(txs[i].TxOutputs,TxOutput{ToAddress:vout.Scriptpubkey_address,Value:strconv.FormatInt(int64(vout.Value),10)})
		}
		txs[i].Timestamp = int64(tx.Status.Block_time) * 1000
		//txs[i].Timestamp = time.Unix(0,int64(tx.Status.Block_time)).Format("2019-05-29T11:11:36+00")
	}
	return txs
}

// 合并同一地址的vout
func MergeVout (btx BtcTx) BtcTx {
	var mtx BtcTx
	mtx = btx
	var vm = make(map[string]float64)
	for _, vout := range btx.Vout {
		k := vout.Scriptpubkey + "|" + vout.Scriptpubkey_address
		vm[k] = vm[k] + vout.Value
	}
	var vouts []VoutData
	for k, v := range vm {
		vout := VoutData{
			Scriptpubkey:strings.Split(k,"|")[0],
			Scriptpubkey_address:strings.Split(k,"|")[1],
			Value:v,
		}
		vouts = append(vouts, vout)
	}
	copy(mtx.Vout, vouts)
	return mtx
}

func ConvertAddresses(btxs []BtcTx) []BtcTx {
	for i, _ := range btxs {
		for j, _ := range btxs[i].Vin {
			btxs[i].Vin[j].Prevout.Scriptpubkey_address = utils.ConvertToTestnet3Address(btxs[i].Vin[j].Prevout.Scriptpubkey_address)
		}
		for j, _ := range btxs[i].Vout {
			btxs[i].Vout[j].Scriptpubkey_address = utils.ConvertToTestnet3Address(btxs[i].Vout[j].Scriptpubkey_address)
		}
	}
	return btxs
}
