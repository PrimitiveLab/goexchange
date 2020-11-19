package gate

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/primitivelab/goexchange"
)

var client = &http.Client{}

func TestGetDepth(t *testing.T) {
	market := New(client, "","")

	response := market.GetDepth(NewSymbol("eos1", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTicker(t *testing.T) {
	market := New(client, "","")

	response := market.GetTicker(NewSymbol("btc", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetKline(t *testing.T) {
	market := New(client, "","")

	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTrade(t *testing.T) {
	market := New(client, "","")

	response := market.GetTrade(NewSymbol("btc", "usdt"), 5, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetSymbolList(t *testing.T) {
	// client := &http.Client{}
	market := New(client, "","")
	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetCoinList(t *testing.T) {
	// client := &http.Client{}
	market := New(client, "","")
	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHttpRequest(t *testing.T) {
	// client := &http.Client{}
	market := New(client, "","")
	params := map[string]string{}
	params["currency_pair"] = NewSymbol("btc", "usdt").ToUpper().ToSymbol("_")
	response := market.HttpRequest("/api/v4/spot/trades", "get", params, false)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}