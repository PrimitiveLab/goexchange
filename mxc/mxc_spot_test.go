package mxc

import (
	"encoding/json"
	. "github.com/primitivelab/goexchange"
	"net/http"
	"testing"
)

var client = &http.Client{}
var apiKey = ""
var secretKey = ""
var baseUrl = ""

func TestMxcSpot_GetSymbolList(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetDepth(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetDepth(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetTicker(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetTicker(NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetKline(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetTrade(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetTrade(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetUserBalance(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetUserOpenTrustOrders(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetUserOrderInfo(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetUserOrderInfo(NewSymbol("eos", "usdt"), "1111111", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetUserTrustOrders(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetUserTrustOrders(NewSymbol("eos", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_GetUserTradeOrders(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.GetUserTradeOrders(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_PlaceLimitOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.PlaceLimitOrder(NewSymbol("eos", "usdt"), "1", "10", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_PlaceMarketOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.PlaceMarketOrder(NewSymbol("eos", "usdt"), "1", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_CancelOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.CancelOrder(NewSymbol("eos", "usdt"), "4439453", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestMxcSpot_BatchCancelOrder(t *testing.T) {
	market := New(client, baseUrl, apiKey, secretKey)

	response := market.BatchCancelOrder(NewSymbol("eos", "usdt"), "4439453,4439454", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
