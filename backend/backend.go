package backend

import (
	"github.com/gaozhengxin/lockinmonitor/backend/request"
)

type Backend struct {}

func New () *Backend {
	return &Backend{}
}

func (b *Backend) GetTransactionsByAddress (cointype string, address string) (res request.Res) {
	ch := make(chan request.Res)
	f := request.Request(cointype,address)
	if f == nil {
		return request.Res{}
	}
	go f(ch)
	res, _ = <-ch
	return
}
