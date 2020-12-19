package poloniex

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

func getInstance() *PoloniexSpot {

	client = &http.Client{}
	config, err := LoadConfig("poloniex")
	if err != nil {
		fmt.Println(err)
	}
	if config != nil {
		if config["key"] != nil {
			apiKey = config["key"].(string)
		}
		if config["secret"] != nil {
			secretKey = config["secret"].(string)
		}
		if config["url"] != nil {
			baseURL = config["url"].(string)
		}
	}

	market := New(client, baseURL, apiKey, secretKey)
	return market
}

func TestPoloniexSpot_GetCoinList(t *testing.T) {
	market := getInstance()

	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetSymbolList(t *testing.T) {
	market := getInstance()

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetDepth(t *testing.T) {
	market := getInstance()

	response := market.GetDepth(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetTicker(t *testing.T) {
	market := getInstance()

	response := market.GetTicker(NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetKline(t *testing.T) {
	market := getInstance()

	options := map[string]string{"start": "1608284813", "end": "1608287813"}
	response := market.GetKline(NewSymbol("btc", "usdt"), KLINE_PERIOD_5MINUTE, 10, options)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetTrade(t *testing.T) {
	market := getInstance()

	response := market.GetTrade(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserBalance(t *testing.T) {
	market := getInstance()

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserOpenTrustOrders(NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserOrderInfo(t *testing.T) {
	market := getInstance()

	response := market.GetUserOrderInfo(NewSymbol("eos", "usdt"), "1111111", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserOrderDetail(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeDetail(NewSymbol("eos", "usdt"), "1111111", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTrustOrders(NewSymbol("eos", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_GetUserTradeOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeOrders(NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_PlaceLimitOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceLimitOrder(NewSymbol("eos", "usdt"), "1", "10", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_PlaceMarketOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceMarketOrder(NewSymbol("eos", "usdt"), "1", BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_CancelOrder(t *testing.T) {
	market := getInstance()

	response := market.CancelOrder(NewSymbol("eos", "usdt"), "4439453", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestPoloniexSpot_BatchCancelOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelOrder(NewSymbol("eos", "usdt"), "4439453,4439454", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
