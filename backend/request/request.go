package request

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"github.com/BurntSushi/toml"
	"github.com/gaozhengxin/lockinmonitor/backend/types"
	"github.com/parnurzeal/gorequest"
)

type Res struct {
	Err error `json:"error,omitempty"`
	Txs []types.Transaction `json:"txs,omitempty"`
}

func Request (cointype string, address string) func (chan Res) () {
	if strings.HasPrefix(cointype,"EVT") {
		return func (ch chan Res) {
			defer func () {
				if e := recover(); e != nil {
					err := fmt.Errorf("GetTransactionsByAddress: Runtime error:  %v\n%v", e, string(debug.Stack()))
					ch <- Res{Err: err}
				}
			}()
			id := strings.TrimPrefix(cointype,"EVT")
			var txs []types.Transaction
			_, b, errs := gorequest.New().
			Post(TheApiConfigs.EVT).
			Set("Accept","application/json").
			Send(`{"sym_id":`+id+`,"addr":"` + address + `"}`).
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
			Get(TheApiConfigs.BTC + "/"+"address/"+address+"/txs").
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
	if cointype == "ETH" {
		return func (ch chan Res) {
			var txs []types.Transaction
			_, b, errs := gorequest.New().
			Get(TheApiConfigs.ETH + "/"+"txs/"+address).
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

var TheApiConfigs *ApiConfigs

type ApiConfigs struct {
	BTC string
	ETH string
	EVT string
}

func LoadConfig (configfile string) error {
	if TheApiConfigs == nil {
		TheApiConfigs = new(ApiConfigs)
	}

	if exists, _ := PathExists(configfile); exists {
		fmt.Printf("use config file: %s\n", configfile)
		_, err := toml.DecodeFile(configfile, TheApiConfigs)
		return err
	} else {
		fmt.Printf("use default config: %s\n", defaultConfig)
		_, err := toml.Decode(defaultConfig, TheApiConfigs)
		return err
	}

	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
