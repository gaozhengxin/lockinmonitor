package backend

import (
	"fmt"
	"runtime/debug"
	"github.com/gaozhengxin/lockinmonitor/backend/request"
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
	return
}
