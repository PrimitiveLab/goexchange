package goexchange

import (
	"testing"
)

func TestSymbol_String(t *testing.T) {
	btc_usdt := NewSymbol("btc", "usdt")
	t.Log(btc_usdt.String() == "BTC_USDT")
	t.Log(btc_usdt.String())
	t.Log(btc_usdt.ToSymbol("-"))
	t.Log(btc_usdt.ToSymbol(""))
	t.Log(btc_usdt.ToLower().ToSymbol("-"))
	t.Log(btc_usdt.ToUpper().ToSymbol("-"))

	// t.Log(btc_usdt.ToLower().ToSymbol("") == "btcusd")
	// t.Log(btc_usdt.ToLower().String() == "btc_usd")
	// t.Log(btc_usdt.Reverse().String() == "USD_BTC")
	// t.Log(btc_usdt.Eq(BTC_USD))
}
