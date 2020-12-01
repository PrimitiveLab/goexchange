package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string {
	KLINE_PERIOD_1MINUTE:   "1m",
	KLINE_PERIOD_3MINUTE:   "3m",
	KLINE_PERIOD_5MINUTE:   "5m",
	KLINE_PERIOD_15MINUTE:  "15m",
	KLINE_PERIOD_30MINUTE:  "30m",
	KLINE_PERIOD_60MINUTE:  "1h",
	KLINE_PERIOD_1HOUR:  	"1h",
	KLINE_PERIOD_4HOUR:  	"4h",
	KLINE_PERIOD_6HOUR:  	"6h",
	KLINE_PERIOD_8HOUR:  	"8h",
	KLINE_PERIOD_12HOUR:  	"12h",
	KLINE_PERIOD_1DAY:   	"1d",
	KLINE_PERIOD_3DAY:   	"3d",
	KLINE_PERIOD_1WEEK:  	"1w",
	KLINE_PERIOD_1MONTH: 	"1M",
}

type BinanceSpot struct {
	httpClient *http.Client
	baseUrl string
	accessKey string
	secretKey string
}

func New(client *http.Client, apiKey, secretKey string) *BinanceSpot {
	instance := new(BinanceSpot)
	instance.baseUrl = "https://api.binance.com"
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

func NewWithConfig(config *APIConfig) *BinanceSpot {
	instance := new(BinanceSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://api.binance.com"
	} else {
		instance.baseUrl = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

func (binanceSpot *BinanceSpot) GetExchangeName() string {
	return ECHANGE_BINANCE
}

func (binanceSpot *BinanceSpot) GetCoinList() interface{} {

	params := &url.Values{}
	result := binanceSpot.httpRequest("/sapi/v1/capital/config/getall", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (binanceSpot *BinanceSpot) GetSymbolList() interface{} {

	params := &url.Values{}
	result := binanceSpot.httpRequest("/api/v3/exchangeInfo", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (binanceSpot *BinanceSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {

	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	// depth ranges [5, 10, 20, 50, 100, 500, 1000, 5000]
	if size < 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else if size <= 1000 {
		size = 1000
	} else {
		size = 5000
	}
	params.Set("limit", strconv.Itoa(size))

	result := binanceSpot.httpRequest("/api/v3/depth", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (binanceSpot *BinanceSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	result := binanceSpot.httpRequest("/api/v3/ticker/24hr", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (binanceSpot *BinanceSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1m"
	}
	params.Set("interval", periodStr)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}

	result := binanceSpot.httpRequest("/api/v3/klines", "get", params,false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (binanceSpot *BinanceSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {

	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	result := binanceSpot.httpRequest("/api/v3/trades", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 获取余额
func (spot *BinanceSpot) GetUserBalance() interface{} {
	return nil
}

// 批量下单
func (spot *BinanceSpot) PlaceOrder(order *PlaceOrder) interface{} {
	return nil
}

// 下限价单
func (spot *BinanceSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// 下市价单
func (spot *BinanceSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// 批量下限价单
func (spot *BinanceSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return nil
}

// 撤单
func (spot *BinanceSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// 批量撤单
func (spot *BinanceSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	return nil
}

// 我的当前委托单
func (spot *BinanceSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 委托单详情
func (spot *BinanceSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// 我的成交单列表
func (spot *BinanceSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 我的委托单列表
func (spot *BinanceSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	return nil
}

func (binanceSpot *BinanceSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	return nil
}

func (binanceSpot *BinanceSpot) httpRequest(url , method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := binanceSpot.baseUrl + url + "?" + params.Encode()
		fmt.Println(requestUrl)
		responseMap = HttpGet(binanceSpot.httpClient, requestUrl)
	// case "POST":
	// 	return nil
	}

	var returnData map[string]interface{}
	returnData = make(map[string]interface{})

	returnData["code"] = responseMap.Code
	returnData["st"] = responseMap.St
	returnData["et"] = responseMap.Et
	if responseMap.Code != 0 {
		returnData["msg"] = responseMap.Msg
		returnData["error"] = responseMap.Error
		return returnData
	}

	var bodyDataMap interface{}
	err := json.Unmarshal(responseMap.Data, &bodyDataMap)
	if err != nil {
		log.Println(string(responseMap.Data))
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}

	returnData["data"] = bodyDataMap
	return returnData
}
