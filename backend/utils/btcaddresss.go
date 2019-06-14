package utils

import (
	"reflect"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func ConvertToTestnet3Address (addr string) string {
	cfg := &chaincfg.MainNetParams
	testnet := &chaincfg.TestNet3Params
	address, err := btcutil.DecodeAddress(addr, cfg)
	if err != nil {
		return addr
	}
	addrtype := reflect.TypeOf(address).String()

	if addrtype == "*btcutil.AddressPubKeyHash" {
		hash160 := address.(*btcutil.AddressPubKeyHash).Hash160()
		taddress, err := btcutil.NewAddressPubKeyHash(hash160[:],testnet)
		if err != nil {
			return addr
		}
		taddr := taddress.EncodeAddress()
		return taddr
	} else if addrtype == "*btcutil.AddressScriptHash" {
		hash160 := address.(*btcutil.AddressScriptHash).Hash160()
		taddress, err := btcutil.NewAddressScriptHash(hash160[:],testnet)
		if err != nil {
			return addr
		}
		taddr := taddress.EncodeAddress()
		return taddr
	}
	return addr
}
