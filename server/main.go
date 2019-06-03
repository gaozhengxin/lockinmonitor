package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/gaozhengxin/lockinmonitor/backend"
	"github.com/gaozhengxin/lockinmonitor/backend/types"
)

func main() {
	p := flag.String("port","8080","port")
	flag.Parse()
	port := ":" + *p

	b := backend.New()
	router := gin.Default()
	router.GET("txs/:cointype/:address",func(c *gin.Context) {
		HandleGetTxs(b, c)
	})
	router.Run(port)
}

func HandleGetTxs(b *backend.Backend, c *gin.Context) {
		cointype := c.Param("cointype")
		address := c.Param("address")
		if cointype == "" {
			c.String(400, "param cointype not set")
			return
		}
		if address == "" {
			c.String(400, "param address not set")
			return
		}
		res := b.GetTransactionsByAddress(cointype,address)
		if res.Err == nil && res.Txs != nil {
			data := types.MarshalTransactions(res.Txs)
			c.Data(200, "application/json", data)
			return
		} else if res.Err != nil {
			c.String(500, res.Err.Error())
			return
		} else {
			c.String(500, "unknown error")
			return
		}
}
