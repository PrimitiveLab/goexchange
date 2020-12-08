package bitz

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
var baseUrl = ""

func TestGetDepth(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	response := market.GetDepth(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTicker(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	response := market.GetTicker(NewSymbol("btc", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetKline(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetTrade(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	response := market.GetTrade(NewSymbol("btc", "usdt"), 5, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetSymbolList(t *testing.T) {
	// client := &http.Client{}
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestGetCoinList(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHttpRequest(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	params := map[string]string{}
	params["symbol"] = NewSymbol("btc", "usdt").ToSymbol("_")
	response := market.HttpRequest("/Market/order", "get", params, false)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_GetUserBalance(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_GetUserTrustOrders(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetUserTrustOrders(NewSymbol("eos", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetUserOpenTrustOrders(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_GetUserOrderInfo(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)
	response := market.GetUserOrderInfo(NewSymbol("eos", "usdt"), "4439453", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_PlaceLimitOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.PlaceLimitOrder(NewSymbol("eos", "usdt"), "10", "1", SELL, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_PlaceMarketOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.PlaceMarketOrder(NewSymbol("eos", "usdt"), "1", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_BatchPlaceLimitOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string

	order := LimitOrder{}
	order.Symbol = NewSymbol("eos", "usdt")
	order.Price = "1.2"
	order.Amount = "10"
	order.Side = BUY

	order1 := LimitOrder{}
	order1.Symbol = NewSymbol("eos", "usdt")
	order1.Price = "2"
	order1.Amount = "0.8"
	order1.Side = BUY

	orders := []LimitOrder{order, order1}

	response := market.BatchPlaceLimitOrder(orders)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_CancelOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.CancelOrder(NewSymbol("eos", "usdt"), "4439453", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBitzSpot_BatchCancelOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey, passphrase)

	// symbol Symbol, status string, size int, options map[string]string
	response := market.BatchCancelOrder(NewSymbol("eos", "usdt"), "4439457,4439458", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
