package types

import (
	"log"
	"strings"
	"time"
	"encoding/json"
)

var tokenid string

func UnmarshalEvtActions (jsonbytes []byte) ([]EvtAction) {
	var evtActions []EvtAction
	err := json.Unmarshal(jsonbytes, &evtActions)
	if err != nil {
		return nil
	}
	for i, act := range evtActions {
		if act.Name == "issuefungible" {
			evtActions[i].Data = &IssuefungibleData{Memo:act.DataI["memo"],Number:EvtAmount(act.DataI["number"]),Address:act.DataI["address"]}
		} else if act.Name == "transferft" {
			evtActions[i].Data = &TransferftData{From:act.DataI["from"],To:act.DataI["to"],Number:EvtAmount(act.DataI["number"])}
		} else {
			continue
		}
		//if checktoken(evtActions[i], tokenid) == false {
		//	continue
		//}
	}
	return evtActions
}

func EvtActionToTransaction (e []EvtAction) ([]Transaction) {
	txs := make([]Transaction, len(e))
	for i, act := range e {
		txs[i].Txhash = act.Trx_id
		timestamp, err := time.Parse("2006-01-02T15:04:05+00",act.Timestamp)
		if err == nil {
			txs[i].Timestamp = timestamp.Unix() * 1000
		} else {
			log.Printf("get Evt transaction time stamp error: %v\n",err)
		}
		if act.Name == "issuefungible" {
			txs[i].FromAddress = "the token issuer"
			txs[i].TxOutputs = append(txs[i].TxOutputs,TxOutput{ToAddress:act.Data.(*IssuefungibleData).Address,Value:act.Data.(*IssuefungibleData).Number.ToString()})
		}
		if act.Name == "transferft" {
			txs[i].FromAddress = act.Data.(*TransferftData).From
			txs[i].TxOutputs = append(txs[i].TxOutputs,TxOutput{ToAddress:act.Data.(*TransferftData).To,Value:act.Data.(*TransferftData).Number.ToString()})
		}
	}
	return txs
}

type EvtAction struct {
	Trx_id string
	Name string
	Domain string
	Key string
	DataI map[string]string `json:"data"`
	Data interface{}
	Timestamp string
}

type IssuefungibleData struct {
	Memo string
	Number EvtAmount
	Address string
}

type TransferftData struct {
	From string
	To string
	Number EvtAmount
}

type EvtAmount string

func (amt EvtAmount) ToString () string {
	num := strings.Split(string(amt)," ")[0]
	str := strings.Replace(num,".","",-1)
	return str
}

func checktoken (act EvtAction, id string) bool {
	var number string
	if act.Name == "issuefungible" {
		number = act.Data.(*IssuefungibleData).Number.ToString()
	} else if act.Name == "transferft" {
		number = act.Data.(*TransferftData).Number.ToString()
	} else {
		return false
	}
	ss := strings.Split(number,"#")
	return len(ss) == 2 && ss[1] == id
}
