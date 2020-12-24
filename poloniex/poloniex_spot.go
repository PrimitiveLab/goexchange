package poloniex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	KLINE_PERIOD_5MINUTE:  "300",
	KLINE_PERIOD_15MINUTE: "900",
	KLINE_PERIOD_30MINUTE: "1800",
	KLINE_PERIOD_2HOUR:    "7200",
	KLINE_PERIOD_4HOUR:    "14400",
	KLINE_PERIOD_1DAY:     "86400",
}

const (
	POLONIEX_BUY  string = "buy"
	POLONIEX_SELL string = "sell"
)

// PoloniexSpot Poloniex exchange spot
type PoloniexSpot struct {
	httpClient *http.Client
	baseUrl    string
	accessKey  string
	secretKey  string
}

// New new instance
func New(client *http.Client, baseURL, apiKey, secretKey string) *PoloniexSpot {
	instance := new(PoloniexSpot)
	if baseURL == "" {
		instance.baseUrl = "https://poloniex.com"
	} else {
		instance.baseUrl = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

// NewWithConfig new instance by config
func NewWithConfig(config *APIConfig) *PoloniexSpot {
	instance := new(PoloniexSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://poloniex.com"
	} else {
		instance.baseUrl = config.Endpoint
	}

	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

// GetExchangeName return exchange name
func (spot *PoloniexSpot) GetExchangeName() string {
	return EXCHANGE_POLONIEX
}

// GetCoinList exchange supported coins
func (spot *PoloniexSpot) GetCoinList() interface{} {
	params := &url.Values{}
	params.Set("command", "returnCurrencies")
	return spot.httpGet("/public", params)
}

// GetSymbolList exchange all symbol
func (spot *PoloniexSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	params.Set("command", "returnTicker")
	return spot.httpGet("/public", params)
}

// GetDepth symbol depth
func (spot *PoloniexSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("command", "returnOrderBook")
	params.Set("currencyPair", spot.getSymbol(symbol))
	params.Set("depth", strconv.FormatInt(int64(size), 10))
	return spot.httpGet("/public", params)
}

// GetTicker symbol ticker
func (spot *PoloniexSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("command", "returnTicker")
	return spot.httpGet("/public", params)
}

// GetKline symbol kline
func (spot *PoloniexSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnChartData")
	params.Set("currencyPair", spot.getSymbol(symbol))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "300"
	}
	params.Set("period", periodStr)
	if start, ok := options["start"]; ok == true {
		params.Set("start", start)
	}
	if end, ok := options["end"]; ok == true {
		params.Set("end", end)
	}

	return spot.httpGet("/public", params)
}

// GetTrade symbol last trade
func (spot *PoloniexSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnTradeHistory")
	params.Set("currencyPair", spot.getSymbol(symbol))
	if start, ok := options["start"]; ok == true {
		params.Set("start", start)
	}
	if end, ok := options["end"]; ok == true {
		params.Set("end", end)
	}

	return spot.httpGet("/public", params)
}

// GetUserBalance user balance
func (spot *PoloniexSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	params.Set("command", "returnCompleteBalances")
	result := spot.httpPost("/tradingApi", params)
	return result
}

// PlaceOrder place order
func (spot *PoloniexSpot) PlaceOrder(order *PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("currencyPair", spot.getSymbol(order.Symbol))
	params.Set("rate", order.Price)
	params.Set("amount", order.Amount)
	if order.Side == BUY {
		params.Set("command", POLONIEX_BUY)
	} else {
		params.Set("command", POLONIEX_SELL)
	}
	switch order.TimeInForce {
	case IOC:
		params.Set("immediateOrCancel", "1")
	case FOK:
		params.Set("fillOrKill", "1")
	case POC:
		params.Set("postOnly", "1")
	}
	if order.ClientOrderId != "" {
		params.Set("clientOrderId", order.ClientOrderId)
	}

	result := spot.httpPost("/tradingApi", params)
	return result
}

// PlaceLimitOrder place limit order
func (spot *PoloniexSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("currencyPair", spot.getSymbol(symbol))
	params.Set("rate", price)
	params.Set("amount", amount)
	if side == BUY {
		params.Set("command", POLONIEX_BUY)
	} else {
		params.Set("command", POLONIEX_SELL)
	}
	if ClientOrderID != "" {
		params.Set("clientOrderId", ClientOrderID)
	}
	result := spot.httpPost("/tradingApi", params)
	return result
}

// PlaceMarketOrder place market order
func (spot *PoloniexSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderID string) interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// BatchPlaceLimitOrder batch place limit order
func (spot *PoloniexSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// CancelOrder cancel a order
func (spot *PoloniexSpot) CancelOrder(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("command", "cancelOrder")
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	} else {
		params.Set("orderNumber", orderID)
	}

	result := spot.httpPost("/tradingApi", params)
	return result
}

// BatchCancelOrder batch cancel orders
func (spot *PoloniexSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// GetUserOpenTrustOrders get current trust order
func (spot *PoloniexSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnOpenOrders")
	params.Set("currencyPair", spot.getSymbol(symbol))
	result := spot.httpPost("/tradingApi", params)
	return result
}

// GetUserOrderInfo get trust order info
func (spot *PoloniexSpot) GetUserOrderInfo(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnOrderStatus")
	params.Set("orderNumber", orderID)
	result := spot.httpPost("/tradingApi", params)
	return result
}

// GetUserTradeDetail get trust order trade detail
func (spot *PoloniexSpot) GetUserTradeDetail(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnOrderTrades")
	params.Set("orderNumber", orderID)
	result := spot.httpPost("/tradingApi", params)
	return result
}

// GetUserTradeOrders get trade order list
func (spot *PoloniexSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("command", "returnTradeHistory")
	params.Set("currencyPair", spot.getSymbol(symbol))
	if size != 0 {
		params.Set("limit", strconv.FormatInt(int64(size), 10))
	}
	if start, ok := options["start"]; ok == true {
		params.Set("start", start)
	}
	if end, ok := options["end"]; ok == true {
		params.Set("end", end)
	}

	result := spot.httpPost("/tradingApi", params)
	return result
}

// GetUserTrustOrders get trust order list
func (spot *PoloniexSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// HttpRequest request api
func (spot *PoloniexSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}

	switch method {
	case HTTP_GET:
		return spot.httpGet(requestUrl, params)
	case HTTP_POST:
		return spot.httpPost(requestUrl, params)
	}
	return nil
}

func (spot *PoloniexSpot) httpGet(url string, params *url.Values) map[string]interface{} {
	var responseMap HttpClientResponse

	requestUrl := spot.baseUrl + url + "?" + params.Encode()
	responseMap = HttpGet(spot.httpClient, requestUrl)
	fmt.Println(requestUrl)
	return spot.handlerResponse(&responseMap)
}

func (spot *PoloniexSpot) httpPost(url string, params *url.Values) map[string]interface{} {
	var responseMap HttpClientResponse

	requestUrl := spot.baseUrl + url
	params.Set("nonce", GetNowMillisecondStr())
	sign, _ := HmacSha512Signer(params.Encode(), spot.secretKey)
	headers := map[string]string{
		"Key":  spot.accessKey,
		"Sign": sign,
	}

	responseMap = HttpPostWithHeader(spot.httpClient, requestUrl, params.Encode(), headers)
	fmt.Println(requestUrl)
	return spot.handlerResponse(&responseMap)
}

func (spot *PoloniexSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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

	var bodyData interface{}
	err := json.Unmarshal(responseMap.Data, &bodyData)
	if err != nil {
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}

	switch bodyData.(type) {
	case map[string]interface{}:
		bodyDataMap := bodyData.(map[string]interface{})
		if _, ok := bodyDataMap["error"]; ok {
			fmt.Println(bodyDataMap)
			returnData["code"] = ExchangeError.Code
			returnData["msg"] = bodyDataMap["error"].(string)
			returnData["error"] = returnData["msg"]
			return returnData
		}
		returnData["data"] = bodyDataMap
	default:
		returnData["data"] = bodyData
	}

	return returnData
}

func (spot PoloniexSpot) getSymbol(symbol Symbol) string {
	return symbol.Reverse().ToUpper().ToSymbol("_")
}
