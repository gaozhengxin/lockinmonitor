package main

import (
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"fmt"
)

type EvtTx struct {
	Trx_id string `json:trx_id`
}

func main () {
	evttxs := new([]EvtTx)
	parseres := func (response gorequest.Response, v interface{}, body []byte, errs []error) {
		if len(errs) > 0 {
			return
		}
		jsonerr := json.Unmarshal(body, v)
		if jsonerr != nil {
			panic(jsonerr)
		}
	}

	_, _, err := gorequest.New().
	Post("https://testnet1.everitoken.io/v1/history/get_fungible_actions").
	Set("Accept","application/json").
	Send(`{"sym_id":1,"addr":"EVT5EbwKfAUyTEpQCX2U4WGf73yUmbTVGjVgrikG3Ve5ufoQyXWYc"}`).
	EndStruct(evttxs, parseres)

	if err != nil {
		panic(err)
	}
	fmt.Printf("evttxs is %+v\n",evttxs)
}
