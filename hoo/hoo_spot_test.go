package hoo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "github.com/primitivelab/goexchange"
)

var client = &http.Client{}
var apiKey = ""
var secretKey = ""
var baseURL = ""

func getInstance() *HooSpot {

	client = &http.Client{}
	config, err := LoadConfig("hoo")
	if err != nil {
		fmt.Println(err)
	}
	if config != nil {
		apiKey = config["key"].(string)
		secretKey = config["secret"].(string)
		baseURL = config["url"].(string)
	}

	market := New(client, baseURL, apiKey, secretKey)
	return market
}

func TestHooSpot_GetSymbolList(t *testing.T) {
	market := getInstance()

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetDepth(t *testing.T) {
	market := getInstance()

	response := market.GetDepth(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetTicker(t *testing.T) {
	market := getInstance()

	response := market.GetTicker(NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetKline(t *testing.T) {
	market := getInstance()

	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetTrade(t *testing.T) {
	market := getInstance()

	response := market.GetTrade(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetUserBalance(t *testing.T) {
	market := getInstance()

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserOpenTrustOrders(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetUserOrderInfo(t *testing.T) {
	market := getInstance()

	response := market.GetUserOrderInfo(NewSymbol("eos", "usdt"), "1111111", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetUserTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTrustOrders(NewSymbol("eos", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_GetUserTradeOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeOrders(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_PlaceLimitOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceLimitOrder(NewSymbol("eos", "usdt"), "1", "10", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_PlaceMarketOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceMarketOrder(NewSymbol("eos", "usdt"), "1", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_CancelOrder(t *testing.T) {
	market := getInstance()

	response := market.CancelOrder(NewSymbol("eos", "usdt"), "4439453", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHooSpot_BatchCancelOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelOrder(NewSymbol("eos", "usdt"), "4439453,4439454", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
