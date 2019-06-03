package types

import (
	"strings"
	"encoding/json"
)

var tokenid string = "1"

func UnmarshalEvtActions (jsonbytes []byte) ([]EvtAction) {
	var evtActions []EvtAction
	err := json.Unmarshal(jsonbytes, &evtActions)
	if err != nil {
		return nil
	}
	for i, act := range evtActions {
		if act.Name == "issuefungible" {
			evtActions[i].Data = &IssuefungibleData{Memo:act.DataI["memo"],Number:EvtAmount(act.DataI["number"]),Address:act.DataI["address"]}
		} else if act.Name == "transfertf" {
			evtActions[i].Data = &TransfertfData{From:act.DataI["from"],To:act.DataI["To"],Number:EvtAmount(act.DataI["number"])}
		} else {
			continue
		}
		if checktoken(evtActions[i], tokenid) == false {
			continue
		}
	}
	return evtActions
}

func EvtActionToTransaction (e []EvtAction) ([]Transaction) {
	txs := make([]Transaction, len(e))
	for i, act := range e {
		txs[i].Txhash = act.Trx_id
		txs[i].Timestamp = act.Timestamp
		if act.Name == "issuefungible" {
			txs[i].FromAddress = "the token issuer"
			txs[i].TxOutputs = append(txs[i].TxOutputs,TxOutput{ToAddress:act.Data.(*IssuefungibleData).Address,Value:act.Data.(*IssuefungibleData).Number.ToString()})
		}
		if act.Name == "transfertf" {
			txs[i].FromAddress = act.Data.(TransfertfData).From
			txs[i].TxOutputs = append(txs[i].TxOutputs,TxOutput{ToAddress:act.Data.(*TransfertfData).To,Value:act.Data.(*TransfertfData).Number.ToString()})
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

type TransfertfData struct {
	From string
	To string
	Number EvtAmount
}

func (evtact *EvtAction) MarshalJSON () string {
	return ""
}

func (evtact *EvtAction) UnmarshalJSON (jsonstr string) error {
	return nil
}

type EvtAmount string

func (amt EvtAmount) ToString () string {
	num := strings.Split(string(amt)," ")[0]
	str := strings.Replace(num,".","",-1)
	return str
}

func checktoken (act EvtAction, id string) bool {
	ss := strings.Split(act.Data.(*IssuefungibleData).Number.ToString(),"#")
	return len(ss) == 2 && ss[1] == id
}
