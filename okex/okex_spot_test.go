package okex

import (
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/primitivelab/goexchange"
)

var client = &http.Client{}
var apiKey = ""
var secretKey = ""
var passphrase = ""

func TestGetDepth(t *testing.T) {
	market := New(client, "","","","")

	response := market.GetDepth(NewSymbol("eos", "usdt"), 21, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTicker(t *testing.T) {
	market := New(client, "","","","")

	response := market.GetTicker(NewSymbol("btc", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetKline(t *testing.T) {
	market := New(client, "","","","")

	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTrade(t *testing.T) {
	market := New(client, "","","","")

	response := market.GetTrade(NewSymbol("btc", "usdt"), 21, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetSymbolList(t *testing.T) {
	market := New(client, "","","","")
	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetCoinList(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)
	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetUserBalance(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)
	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestCancelOrder(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	response := market.CancelOrder(NewSymbol("btc", "usdt"), "1111111", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetUserTrustOrders(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.GetUserTrustOrders(NewSymbol("btc", "usdt"), "7", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetUserTradeOrders(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.GetUserTradeOrders(NewSymbol("btc", "usdt"),10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetUserOrderInfo(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.GetUserOrderInfo(NewSymbol("btc", "usdt"),"6028843635265536", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetUserOpenTrustOrders(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.GetUserOpenTrustOrders(NewSymbol("btc", "usdt"),10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestOkexSpot_PlaceLimitOrderOrders(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.PlaceLimitOrder(NewSymbol("btc", "usdt"),"13", "1", BUY, "a1234444445")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestOkexSpot_BatchPlaceLimitOrder(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string

	order := LimitOrder{}
	order.Symbol = NewSymbol("link", "usdt")
	order.Price = "13"
	order.Amount = "1"
	order.Side = BUY

	order1 := LimitOrder{}
	order1.Symbol = NewSymbol("link", "usdt")
	order1.Price = "13"
	order1.Amount = "0.8"
	order1.Side = BUY

	orders := []LimitOrder{order, order1}

	response := market.BatchPlaceLimitOrder(orders)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestOkexSpot_CancelOrder(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.CancelOrder(NewSymbol("btc", "usdt"),"", "a1234444444")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestOkexSpot_BatchCancelOrder(t *testing.T) {
	market := New(client, "", apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.BatchCancelOrder(NewSymbol("link", "usdt"),"6034189181870081,6034189181870082", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHttpRequest(t *testing.T) {
	// client := &http.Client{}
	market := New(client, "","","","")
	instrumentId := NewSymbol("btc", "usdt").ToUpper().ToSymbol("-")
	params := map[string]interface{}{}
	params["granularity"] = "300"
	response := market.HttpRequest("/api/spot/v3/instruments/"+instrumentId+"/candles", "get", params, false)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}