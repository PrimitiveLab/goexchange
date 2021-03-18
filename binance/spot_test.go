package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	goex "github.com/primitivelab/goexchange"
)

var client = &http.Client{}
var apiKey = ""
var secretKey = ""
var baseURL = ""
var proxy = ""

func getInstance() *Spot {

	client = &http.Client{}

	config, err := goex.LoadConfig("binance")
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

func TestBinanceSpot_GetCoinList(t *testing.T) {
	market := getInstance()

	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetSymbolList(t *testing.T) {
	market := getInstance()

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetDepth(t *testing.T) {
	market := getInstance()

	response := market.GetDepth(goex.NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetTicker(t *testing.T) {
	market := getInstance()

	response := market.GetTicker(goex.NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetKline(t *testing.T) {
	market := getInstance()

	response := market.GetKline(goex.NewSymbol("btc", "usdt"), goex.KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetTrade(t *testing.T) {
	market := getInstance()

	response := market.GetTrade(goex.NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserBalance(t *testing.T) {
	market := getInstance()

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserCommissionRate(t *testing.T) {
	market := getInstance()
	response := market.GetUserCommissionRate(goex.NewSymbol("eos", "usdt"))
	// response := market.GetUserCommissionRate(Symbol{})
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserOpenTrustOrders(goex.NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserOrderInfo(t *testing.T) {
	market := getInstance()

	response := market.GetUserOrderInfo(goex.NewSymbol("eos", "usdt"), "1399414810", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTrustOrders(goex.NewSymbol("eos", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserTradeOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeOrders(goex.NewSymbol("eos", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserDepositAddress(t *testing.T) {
	market := getInstance()

	response := market.GetUserDepositAddress("btc", nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserDepositRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserDepositRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_GetUserWithdrawRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserWithdrawRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_PlaceOrder(t *testing.T) {
	market := getInstance()

	order := goex.PlaceOrder{}
	order.Amount = "10"
	order.ClientOrderId = ""
	order.Price = "1"
	order.Side = goex.BUY
	order.Symbol = goex.NewSymbol("eos", "usdt")
	order.TimeInForce = goex.GTC
	order.TradeType = goex.LIMIT

	response := market.PlaceOrder(&order)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_PlaceLimitOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceLimitOrder(goex.NewSymbol("eos", "usdt"), "1", "10", goex.BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_PlaceMarketOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceMarketOrder(goex.NewSymbol("eos", "usdt"), "1", goex.BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_CancelOrder(t *testing.T) {
	market := getInstance()

	response := market.CancelOrder(goex.NewSymbol("eos", "usdt"), "1402657574", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_BatchCancelOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelOrder(goex.NewSymbol("eos", "usdt"), "", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestBinanceSpot_BatchCancelAllOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelAllOrder(goex.NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
