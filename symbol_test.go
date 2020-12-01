package goexchange

import (
	"fmt"
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

	t.Log(TradeSide(1))
	t.Log(TradeSide(-1))
	t.Log(TradeSide(2))
	t.Log(BUY)
	t.Log(BUY.String())
	t.Log(SELL)

	var a1 map[string]string
	a2 := map[string]string{}
	if a1 == nil {
		fmt.Println("a1 is nil")
	}
	if a2 == nil || len(a2) == 0 {
		fmt.Println("a2 is nil")
	}
	md1(a1)
	md2(a2)
	fmt.Println(a1)
	fmt.Println(a2)

}

func md1(val map[string]string)  {
	val = make(map[string]string)
	val["1"] = "1"
	val["2"] = "2"
}

func md2(val map[string]string)  {
	val["1"] = "1"
	val["2"] = "2"
}