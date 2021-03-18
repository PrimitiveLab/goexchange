package binance

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	goex "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	goex.KLINE_PERIOD_1MINUTE:  "1m",
	goex.KLINE_PERIOD_3MINUTE:  "3m",
	goex.KLINE_PERIOD_5MINUTE:  "5m",
	goex.KLINE_PERIOD_15MINUTE: "15m",
	goex.KLINE_PERIOD_30MINUTE: "30m",
	goex.KLINE_PERIOD_60MINUTE: "1h",
	goex.KLINE_PERIOD_1HOUR:    "1h",
	goex.KLINE_PERIOD_4HOUR:    "4h",
	goex.KLINE_PERIOD_6HOUR:    "6h",
	goex.KLINE_PERIOD_8HOUR:    "8h",
	goex.KLINE_PERIOD_12HOUR:   "12h",
	goex.KLINE_PERIOD_1DAY:     "1d",
	goex.KLINE_PERIOD_3DAY:     "3d",
	goex.KLINE_PERIOD_1WEEK:    "1w",
	goex.KLINE_PERIOD_1MONTH:   "1M",
}

// Spot binance struct
type Spot struct {
	httpClient *http.Client
	baseURL    string
	accessKey  string
	secretKey  string
}

// New new instance
func New(client *http.Client, baseURL, apiKey, secretKey string) *Spot {
	instance := new(Spot)
	if baseURL == "" {
		instance.baseURL = "https://api.binance.com"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

// NewWithConfig new instance with config struct
func NewWithConfig(config *goex.APIConfig) *Spot {
	instance := new(Spot)
	if config.Endpoint == "" {
		instance.baseURL = "https://api.binance.com"
	} else {
		instance.baseURL = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

// GetExchangeName get exchange name
func (spot *Spot) GetExchangeName() string {
	return goex.EXCHANGE_BINANCE
}

// GetCoinList exchange coin list
func (spot *Spot) GetCoinList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/sapi/v1/capital/config/getall", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetSymbolList exchange symbol list
func (spot *Spot) GetSymbolList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/v3/exchangeInfo", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetDepth exchange depth data
func (spot *Spot) GetDepth(symbol goex.Symbol, size int, options map[string]string) map[string]interface{} {

	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
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

	result := spot.httpGet("/api/v3/depth", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetTicker exchange ticker data
func (spot *Spot) GetTicker(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	result := spot.httpGet("/api/v3/ticker/24hr", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetKline exchange kline data
func (spot *Spot) GetKline(symbol goex.Symbol, period, size int, options map[string]string) interface{} {
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

	result := spot.httpGet("/api/v3/klines", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetTrade exchange trade order data
func (spot *Spot) GetTrade(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().ToSymbol(""))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	result := spot.httpGet("/api/v3/trades", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetUserBalance user account balance
func (spot *Spot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/v3/account", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserCommissionRate user current commission rate
func (spot *Spot) GetUserCommissionRate(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	if symbol.CoinFrom != "" {
		params.Set("symbol", spot.getSymbol(symbol))
	}

	result := spot.httpGet("/wapi/v3/tradeFee.html", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// PlaceOrder place order
func (spot *Spot) PlaceOrder(order *goex.PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(order.Symbol))
	if order.ClientOrderId != "" {
		params.Set("newClientOrderId", order.ClientOrderId)
	}
	params.Set("side", strings.ToUpper(order.Side.String()))
	params.Set("quantity", order.Amount)
	if order.TradeType == goex.LIMIT {
		params.Set("price", order.Price)
		params.Set("type", strings.ToUpper(goex.LIMIT))
		switch order.TimeInForce {
		case goex.IOC:
			params.Set("timeInForce", "IOC")
		case goex.FOK:
			params.Set("timeInForce", "FOK")
		default:
			params.Set("timeInForce", "GTC")
		}
	} else {
		params.Set("type", strings.ToUpper(goex.MARKET))
	}

	result := spot.httpPost("/api/v3/order", params, true)
	return result
}

// PlaceLimitOrder place limit order
func (spot *Spot) PlaceLimitOrder(symbol goex.Symbol, price string, amount string, side goex.TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("timeInForce", "GTC")
	params.Set("type", strings.ToUpper(goex.LIMIT))
	params.Set("side", strings.ToUpper(side.String()))
	if ClientOrderID != "" {
		params.Set("newClientOrderId", ClientOrderID)
	}
	result := spot.httpPost("/api/v3/order", params, true)
	return result
}

// PlaceMarketOrder place market order
func (spot *Spot) PlaceMarketOrder(symbol goex.Symbol, amount string, side goex.TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("quantity", amount)
	params.Set("type", strings.ToUpper(goex.MARKET))
	params.Set("side", strings.ToUpper(side.String()))
	if ClientOrderID != "" {
		params.Set("newClientOrderId", ClientOrderID)
	}
	result := spot.httpPost("/dapi/v1/order", params, true)
	return result
}

// BatchPlaceLimitOrder batch place limit order
func (spot *Spot) BatchPlaceLimitOrder(orders []goex.LimitOrder) interface{} {
	return goex.ReturnAPIError(goex.MethodNotExistError)
}

// CancelOrder cancel user trust order
func (spot *Spot) CancelOrder(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if clientOrderID != "" {
		params.Set("origClientOrderId", clientOrderID)
	} else {
		params.Set("orderId", orderID)
	}
	result := spot.httpDelete("/api/v3/order", params, true)
	return result
}

// BatchCancelOrder batch cancel trust order
func (spot *Spot) BatchCancelOrder(symbol goex.Symbol, orderIds, clientOrderIds string) interface{} {
	return goex.ReturnAPIError(goex.MethodNotExistError)
}

// BatchCancelAllOrder batch cancel all orders
func (spot *Spot) BatchCancelAllOrder(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpDelete("/api/v3/openOrders", params, true)
	return result
}

// GetUserOpenTrustOrders user open trust order list
func (spot *Spot) GetUserOpenTrustOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpGet("/api/v3/openOrders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserOrderInfo user trust order info
func (spot *Spot) GetUserOrderInfo(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if clientOrderID != "" {
		params.Set("origClientOrderId", clientOrderID)
	} else {
		params.Set("orderId", orderID)
	}

	result := spot.httpGet("/api/v3/order", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserTradeOrders user trade order list
func (spot *Spot) GetUserTradeOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))

	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}
	if fromID, ok := options["fromId"]; ok == true {
		params.Set("fromId", fromID)
	}

	result := spot.httpGet("/api/v3/myTrades", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserTrustOrders user trust order list
func (spot *Spot) GetUserTrustOrders(symbol goex.Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))

	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}
	if orderID, ok := options["orderId"]; ok == true {
		params.Set("orderId", orderID)
	}

	result := spot.httpGet("/api/v3/allOrders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserDepositAddress user deposit address
func (spot *Spot) GetUserDepositAddress(coin string, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("coin", coin)

	if network, ok := options["network"]; ok == true {
		params.Set("network", network)
	}

	result := spot.httpGet("/sapi/v1/capital/deposit/address", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// Withdraw user withdraw
func (spot *Spot) Withdraw(coin, address, tag, amount, chain string, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("coin", coin)
	params.Set("address", address)
	params.Set("amount", amount)
	if tag != "" {
		params.Set("addressTag", tag)
	}

	if chain != "" {
		params.Set("network", chain)
	}

	result := spot.httpPost("/sapi/v1/capital/withdraw/apply", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserDepositRecords user deposit record list
func (spot *Spot) GetUserDepositRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	if coin != "" {
		params.Set("coin", coin)
	}

	if status, ok := options["status"]; ok == true {
		params.Set("status", status)
	}
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}
	if offset, ok := options["offset"]; ok == true {
		params.Set("offset", offset)
	}

	result := spot.httpGet("/sapi/v1/capital/deposit/hisrec", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserWithdrawRecords user withdraw record list
func (spot *Spot) GetUserWithdrawRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	if coin != "" {
		params.Set("coin", coin)
	}

	if status, ok := options["status"]; ok == true {
		params.Set("status", status)
	}
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}
	if offset, ok := options["offset"]; ok == true {
		params.Set("offset", offset)
	}

	result := spot.httpGet("/sapi/v1/capital/withdraw/history", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// HttpRequest request url
func (spot *Spot) HttpRequest(requestURL, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}
	switch method {
	case goex.HTTP_GET:
		return spot.httpGet(requestURL, params, signed)
	case goex.HTTP_POST:
		return spot.httpPost(requestURL, params, signed)
	case goex.HTTP_DELETE:
		return spot.httpPost(requestURL, params, signed)
	}
	return nil
}

// httpGet Get request method
func (spot *Spot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	sign := ""
	if signed {
		headers["X-MBX-APIKEY"] = spot.accessKey
		sign = spot.sign(params)
	}

	requestURL := spot.baseURL + url
	if params != nil {
		requestURL = requestURL + "?" + params.Encode()
		if sign != "" {
			requestURL = requestURL + "&signature=" + sign
		}
	}

	responseMap = goex.HttpGetWithHeader(spot.httpClient, requestURL, headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Post request method
func (spot *Spot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	headers["X-MBX-APIKEY"] = spot.accessKey
	sign := spot.sign(params)
	requestURL := spot.baseURL + url + "?" + params.Encode() + "&signature=" + sign
	responseMap = goex.HttpPostWithHeader(spot.httpClient, requestURL, "", headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Delete request method
func (spot *Spot) httpDelete(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	headers["X-MBX-APIKEY"] = spot.accessKey
	sign := spot.sign(params)
	requestURL := spot.baseURL + url + "?" + params.Encode() + "&signature=" + sign
	responseMap = goex.HttpDeleteWithHeader(spot.httpClient, requestURL, headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Handler response data format
func (spot *Spot) handlerResponse(responseMap *goex.HttpClientResponse) map[string]interface{} {
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
		returnData["code"] = goex.JsonUnmarshalError.Code
		returnData["msg"] = goex.JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}
	returnData["data"] = bodyDataMap
	return returnData
}

// httpGet signature method
func (spot *Spot) sign(params *url.Values) string {
	timestamp := goex.GetNowMillisecondStr()
	params.Set("recvWindow", "5000")
	params.Set("timestamp", timestamp)
	sign, _ := goex.HmacSha256Signer(params.Encode(), spot.secretKey)
	// params.Set("signature", sign)
	return sign
}

// httpGet format symbol method
func (spot Spot) getSymbol(symbol goex.Symbol) string {
	return symbol.ToUpper().ToSymbol("")
}
