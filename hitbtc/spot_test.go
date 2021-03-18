package hitbtc

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

// var proxy = "http://192.168.6.4:3128"

//
func getInstance() *Spot {
	clientConfig := &goex.HTTPClientConfig{
		HTTPTimeout:  5 * time.Second,
		MaxIdleConns: 10,
	}
	// clientConfig.SetProxyURL(proxy)
	client = goex.NewHTTPClientWithConfig(clientConfig)

	// client = &http.Client{}
	config, err := goex.LoadConfig("hitbtc1")
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

func TestHitbtcSpot_GetCoinList(t *testing.T) {
	market := getInstance()

	response := market.GetCoinList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetSymbolList(t *testing.T) {
	market := getInstance()

	response := market.GetSymbolList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetDepth(t *testing.T) {
	market := getInstance()

	response := market.GetDepth(goex.NewSymbol("xrp", "usdt"), 10, map[string]string{"type": "step0"})
	b, err := json.Marshal(response)
	t.Log(string(b))
	t.Log(err)
}

func TestHitbtcSpot_GetTicker(t *testing.T) {
	market := getInstance()

	response := market.GetTicker(goex.NewSymbol("xrp", "usdt"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetKline(t *testing.T) {
	market := getInstance()

	response := market.GetKline(goex.NewSymbol("xrp", "usdt"), goex.KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetTrade(t *testing.T) {
	market := getInstance()

	response := market.GetTrade(goex.NewSymbol("xrp", "usdt"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserBalance(t *testing.T) {
	market := getInstance()

	response := market.GetUserBalance()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserCommissionRate(t *testing.T) {
	market := getInstance()
	response := market.GetUserCommissionRate(goex.NewSymbol("eth", "btc"))
	// response := market.GetUserCommissionRate(Symbol{})
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserOpenTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserOpenTrustOrders(goex.NewSymbol("btc", "usd"), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserOrderInfo(t *testing.T) {
	market := getInstance()

	response := market.GetUserOrderInfo(goex.NewSymbol("btc", "usd"), "6d8d3d0368524ce9ab3fdb8d226caddb", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserTrustOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTrustOrders(goex.NewSymbol("btc", "usd"), "", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserTradeOrders(t *testing.T) {
	market := getInstance()

	response := market.GetUserTradeOrders(goex.NewSymbol("btc", "usd"), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserDepositAddress(t *testing.T) {
	market := getInstance()
	response := market.GetUserDepositAddress("btc", nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserDepositRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserDepositRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_GetUserWithdrawRecords(t *testing.T) {
	market := getInstance()

	response := market.GetUserWithdrawRecords("btc", 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_PlaceOrder(t *testing.T) {
	market := getInstance()

	order := goex.PlaceOrder{}
	order.Amount = "0.01"
	order.ClientOrderId = ""
	order.Price = "1000"
	order.Side = goex.BUY
	order.Symbol = goex.NewSymbol("btc", "usd")
	order.TimeInForce = goex.GTC
	order.TradeType = goex.LIMIT

	response := market.PlaceOrder(&order)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_PlaceLimitOrder(t *testing.T) {
	market := getInstance()

	response := market.PlaceLimitOrder(goex.NewSymbol("eth", "btc"), "0.046016", "0.063", goex.SELL, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_PlaceMarketOrder(t *testing.T) {
	market := getInstance()
	response := market.PlaceMarketOrder(goex.NewSymbol("xrp", "usdt"), "1", goex.BUY, "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_CancelOrder(t *testing.T) {
	market := getInstance()
	response := market.CancelOrder(goex.NewSymbol("btc", "usd"), "6d8d3d0368524ce9ab3fdb8d226caddb", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_BatchCancelOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelOrder(goex.NewSymbol("xrp", "usdt"), "", "")
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestHitbtcSpot_BatchCancelAllOrder(t *testing.T) {
	market := getInstance()

	response := market.BatchCancelAllOrder(goex.NewSymbol("btc", "usd"))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}
