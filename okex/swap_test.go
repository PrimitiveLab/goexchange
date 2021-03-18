package okex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	goex "github.com/primitivelab/goexchange"
)

// var apiKey = ""
// var secretKey = ""
// var baseURL = ""
var CoinFrom = "dot"
var CoinTo = "usdt"

func getSwapInstance() Swap {

	client = &http.Client{}
	config, err := goex.LoadConfig("okex")
	if err != nil {
		fmt.Println(err)
	}
	if config != nil {
		apiKey = config["key"].(string)
		secretKey = config["secret"].(string)
		// baseURL = config["url"].(string)
	}

	// market := NewSwapUsdt(client, baseURL, apiKey, secretKey)
	// market := NewSwapCoin(client, baseURL, apiKey, secretKey)
	conf := goex.APIConfig{}
	conf.ApiKey = apiKey
	conf.ApiSecretKey = secretKey
	conf.HttpClient = client
	market := NewSwapUsdtWithConfig(&conf)
	return market
}

func TestSwap_GetContractList(t *testing.T) {
	market := getSwapInstance()

	response := market.GetContractList()
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestSwap_GetDepth(t *testing.T) {
	market := getSwapInstance()

	response := market.GetDepth(goex.NewSymbol(CoinFrom, CoinTo), 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestSwap_GetTicker(t *testing.T) {
	market := getSwapInstance()

	response := market.GetTicker(goex.NewSymbol(CoinFrom, CoinTo))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestSwap_GetKline(t *testing.T) {
	market := getSwapInstance()

	response := market.GetKline(goex.NewSymbol(CoinFrom, CoinTo), goex.KLINE_PERIOD_5MINUTE, 10, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestSwap_GetTrade(t *testing.T) {
	market := getSwapInstance()

	response := market.GetTrade(goex.NewSymbol(CoinFrom, CoinTo), 2, nil)
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

func TestSwap_GetPremiumIndex(t *testing.T) {
	market := getSwapInstance()

	response := market.GetPremiumIndex(goex.NewSymbol(CoinFrom, CoinTo))
	b, _ := json.Marshal(response)
	t.Log(string(b))
}

// func TestSwap_GetUserBalance(t *testing.T) {
// 	market := getSwapInstance()
// 	response := market.GetUserBalance()
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserAssets(t *testing.T) {
// 	market := getSwapInstance()
// 	response := market.GetUserAssets()
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserPositions(t *testing.T) {
// 	market := getSwapInstance()
// 	response := market.GetUserPositions(goex.Symbol{})
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserAssetsIncomes(t *testing.T) {
// 	market := getSwapInstance()
// 	response := market.GetUserAssetsIncomes(goex.Symbol{}, 5, nil)
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserCommissionRate(t *testing.T) {
// 	market := getSwapInstance()
// 	response := market.GetUserCommissionRate(goex.NewSymbol(CoinFrom, CoinTo))
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserOpenTrustOrders(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.GetUserOpenTrustOrders(goex.NewSymbol(CoinFrom, CoinTo), 2, nil)
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserOrderInfo(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.GetUserOrderInfo(goex.NewSymbol(CoinFrom, CoinTo), "2785058797", "")
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserTrustOrders(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.GetUserTrustOrders(goex.NewSymbol(CoinFrom, CoinTo), "", 10, nil)
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_GetUserTradeOrders(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.GetUserTradeOrders(goex.NewSymbol(CoinFrom, CoinTo), 10, nil)
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_PlaceLimitOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.PlaceLimitOrder(goex.NewSymbol(CoinFrom, CoinTo), "1", "10", goex.BUY, "")
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_BatchPlaceLimitOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	// symbol Symbol, status string, size int, options map[string]string

// 	order := goex.LimitOrder{}
// 	order.Symbol = goex.NewSymbol(CoinFrom, CoinTo)
// 	order.Price = "1"
// 	order.Amount = "10"
// 	order.Side = goex.BUY

// 	order1 := goex.LimitOrder{}
// 	order1.Symbol = goex.NewSymbol(CoinFrom, CoinTo)
// 	order1.Price = "1"
// 	order1.Amount = "20"
// 	order1.Side = goex.BUY

// 	orders := []goex.LimitOrder{order, order1}

// 	response := market.BatchPlaceLimitOrder(orders)
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_PlaceMarketOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.PlaceMarketOrder(goex.NewSymbol(CoinFrom, CoinTo), "1", BUY, "")
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_CancelOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.CancelOrder(goex.NewSymbol(CoinFrom, CoinTo), "2786207147", "")
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_BatchCancelOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.BatchCancelOrder(goex.NewSymbol(CoinFrom, CoinTo), "2786678083,2786678832", "")
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }

// func TestSwap_BatchCancelAllOrder(t *testing.T) {
// 	market := getSwapInstance()

// 	response := market.BatchCancelAllOrder(goex.NewSymbol(CoinFrom, CoinTo))
// 	b, _ := json.Marshal(response)
// 	t.Log(string(b))
// }
