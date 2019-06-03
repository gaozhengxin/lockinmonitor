package request

import (
	"fmt"
	"time"
	"github.com/gaozhengxin/lockinmonitor/backend/types"
	"github.com/parnurzeal/gorequest"
)

type Res struct {
	Err error `json:"error,omitempty"`
	Txs []types.Transaction `json:"txs,omitempty"`
}

func Request (cointype string, address string) func (chan Res) () {
	if cointype == "EVT-1" {
		return func (ch chan Res) {
			var txs []types.Transaction
			_, b, _ := gorequest.New().
			Post("https://testnet1.everitoken.io/v1/history/get_fungible_actions").
			Set("Accept","application/json").
			Send(`{"sym_id":1,"addr":"` + address + `"}`).
			EndStruct(&txs)
			txs = types.ParseTransactions(cointype)(b)
			ch <- Res{Txs:txs}
			time.Sleep(time.Duration(5) * time.Second)
			close(ch)
		}
	}
	return func (ch chan Res) {
		ch <- Res{Err: fmt.Errorf("cointype %v not supported",cointype)}
	}
}
