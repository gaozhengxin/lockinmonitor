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
			_, b, errs := gorequest.New().
			Post("https://testnet1.everitoken.io/v1/history/get_fungible_actions").
			Set("Accept","application/json").
			Send(`{"sym_id":1,"addr":"` + address + `"}`).
			EndBytes()
			if len(errs) > 0 {
				err := fmt.Errorf("Request errors: %+v",errs)
				ch <- Res{Txs:txs,Err:err}
			}
			txs = types.ParseTransactions(cointype)(b)
			ch <- Res{Txs:txs}
			time.Sleep(time.Duration(5) * time.Second)
			close(ch)
		}
	}
	if cointype == "BTC" {
		return func (ch chan Res) {
			var txs []types.Transaction
			_, b, errs := gorequest.New().
			Get("http://5.189.139.168:4000/"+"address/"+address+"/txs").
			Set("Accept","application/json").
			EndBytes()
			if len(errs) > 0 {
				err := fmt.Errorf("Request errors: %+v",errs)
				ch <- Res{Txs:txs,Err:err}
			}
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
