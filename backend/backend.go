package backend

import (
	"fmt"
	"runtime/debug"
	"github.com/gaozhengxin/lockinmonitor/backend/request"
	"github.com/gaozhengxin/lockinmonitor/backend/types"
)

type Backend struct {}

func New () *Backend {
	return &Backend{}
}

func (b *Backend) GetTransactionsByAddress (cointype string, address string) (res request.Res) {
	defer func () {
		if e := recover(); e != nil {
			res.Err = fmt.Errorf("GetTransactionsByAddress: Runtime error:  %v\n%v", e, string(debug.Stack()))
		}
	}()
	ch := make(chan request.Res)
	f := request.Request(cointype,address)
	if f == nil {
		return request.Res{}
	}
	go f(ch)
	res, _ = <-ch
	if res.Txs != nil {
		res.Txs = LockinFilt(res.Txs, address)
	}
	return
}

// 判断是否可以lockin
func LockinFilt(txs []types.Transaction, address string) (litxs []types.Transaction) {
	for _, tx := range txs {
		canlockin := true
		for _, out := range tx.TxOutputs {
			if out.ToAddress == tx.FromAddress && out.ToAddress == address {
				canlockin = false
				break
			}
		}
		if canlockin {
			litxs = append(litxs, tx)
		}
	}
	return
}
