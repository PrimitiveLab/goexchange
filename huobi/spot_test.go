package huobi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	goex "github.com/primitivelab/goexchange"
)

var client = &http.Client{}
var apiKey = ""
var secretKey = ""
var baseURL = ""
var proxy = "http://192.168.6.4:3128"
var accountID = ""

//
func getInstance() *Spot {
	clientConfig := &goex.HTTPClientConfig{
		HTTPTimeout:  5 * time.Second,
		MaxIdleConns: 10,
	}
	clientConfig.SetProxyURL(proxy)
	client = goex.NewHTTPClientWithConfig(clientConfig)

	// client = &http.Client{}
	config, err := goex.LoadConfig("huobi")
	if err != nil {
		fmt.Println(err)
	}
	if config != nil {
		apiKey = config["key"].(string)
		secretKey = config["secret"].(string)
		baseURL = config["url"].(string)
		accountID = config["account_id"].(string)
	}

	market := New(client, baseURL, apiKey, secretKey, accountID)
	return market
}

func TestHuobiSpot_GetCoinList(t *testing.T) {
	market := getInstance()

	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetSymbolList(t *testing.T) {
	market := getInstance()

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetDepth(t *testing.T) {
	market := getInstance()

	response := market.GetDepth(goex.NewSymbol("eos", "usdt"), 10, map[string]string{"type": "step0"})
	b, err := json.Marshal(response)
	t.Log(string(b))
	t.Log(err)
}

func TestHuobiSpot_GetTicker(t *testing.T) {
	market := getInstance()

	response := market.GetTicker(goex.NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetKline(t *testing.T) {
	market := getInstance()

	response := market.GetKline(goex.NewSymbol("btc", "usdt"), goex.KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetTrade(t *testing.T) {
	market := getInstance()

	response := market.GetTrade(goex.NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserBalance(t *testing.T) {
	market := getInstance()

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserCommissionRate(t *testing.T) {
	market := getInstance()
	response := market.GetUserCommissionRate(goex.NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserOpenTrustOrders(goex.NewSymbol("eos", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserOrderInfo(t *testing.T) {
	market := getInstance()

	response := market.GetUserOrderInfo(goex.NewSymbol("iost", "usdt"), "235190449677525", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTrustOrders(goex.NewSymbol("iost", "usdt"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserTradeOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeOrders(goex.NewSymbol("iost", "usdt"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserDepositAddress(t *testing.T) {
	market := getInstance()

	response := market.GetUserDepositAddress("btc", nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserDepositRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserDepositRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_GetUserWithdrawRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserWithdrawRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_PlaceOrder(t *testing.T) {
	market := getInstance()

	order := goex.PlaceOrder{}
	order.Amount = "5"
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

func TestHuobiSpot_PlaceLimitOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceLimitOrder(goex.NewSymbol("eos", "usdt"), "1", "10", goex.BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_PlaceMarketOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceMarketOrder(goex.NewSymbol("eos", "usdt"), "2", goex.BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_CancelOrder(t *testing.T) {
	market := getInstance()

	response := market.CancelOrder(goex.NewSymbol("eos", "usdt"), "235191859432411", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_BatchCancelOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelOrder(goex.NewSymbol("eos", "usdt"), "235191918533757,235191808561994", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_BatchPlaceLimitOrder(t *testing.T) {
	market := getInstance()

	// symbol Symbol, status string, size int, options map[string]string

	order := goex.LimitOrder{}
	order.Symbol = goex.NewSymbol("eos", "usdt")
	order.Price = "1"
	order.Amount = "5"
	order.Side = goex.BUY

	order1 := goex.LimitOrder{}
	order1.Symbol = goex.NewSymbol("eos", "usdt")
	order1.Price = "1.1"
	order1.Amount = "8"
	order1.Side = goex.BUY

	orders := []goex.LimitOrder{order, order1}

	response := market.BatchPlaceLimitOrder(orders)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHuobiSpot_BatchCancelAllOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelAllOrder(goex.NewSymbol("eos", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
