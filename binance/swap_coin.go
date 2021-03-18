package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	goex "github.com/primitivelab/goexchange"
)

// SwapCoin binance coin margined contract
type SwapCoin struct {
	httpClient *http.Client
	baseURL    string
	accessKey  string
	secretKey  string
}

// NewSwapCoin new instance
func NewSwapCoin(client *http.Client, baseURL, apiKey, secretKey string) *SwapCoin {
	instance := new(SwapCoin)
	if baseURL == "" {
		instance.baseURL = "https://dapi.binance.com"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

// NewSwapCoinWithConfig new instance with config struct
func NewSwapCoinWithConfig(config *goex.APIConfig) *SwapCoin {
	instance := new(SwapCoin)
	if config.Endpoint == "" {
		instance.baseURL = "https://dapi.binance.com"
	} else {
		instance.baseURL = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

// GetExchangeName get exchange name
func (swap *SwapCoin) GetExchangeName() string {
	return goex.EXCHANGE_BINANCE
}

// GetContractList exchange contract list
func (swap *SwapCoin) GetContractList() interface{} {
	params := &url.Values{}
	result := swap.httpGet("/dapi/v1/exchangeInfo", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetDepth exchange depth data
func (swap *SwapCoin) GetDepth(symbol goex.Symbol, size int, options map[string]string) map[string]interface{} {
	// goex.Symbol
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
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
	} else {
		size = 1000
	}

	params.Set("limit", strconv.Itoa(size))

	result := swap.httpGet("/dapi/v1/depth", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetTicker exchange ticker data
func (swap *SwapCoin) GetTicker(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	result := swap.httpGet("/dapi/v1/ticker/24hr", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetTickerBook exchange ticker data
func (swap *SwapCoin) GetTickerBook(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	result := swap.httpGet("/fapi/v1/ticker/bookTicker", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetKline exchange kline data
func (swap *SwapCoin) GetKline(symbol goex.Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
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

	result := swap.httpGet("/dapi/v1/klines", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetTrade exchange trade order data
func (swap *SwapCoin) GetTrade(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	result := swap.httpGet("/dapi/v1/trades", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetPremiumIndex exchange index price& market price & funding rate
func (swap *SwapCoin) GetPremiumIndex(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	if symbol.CoinFrom != "" {
		params.Set("symbol", swap.getSymbol(symbol))
	}
	result := swap.httpGet("/dapi/v1/premiumIndex", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserBalance user account balance
func (swap *SwapCoin) GetUserBalance() interface{} {
	params := &url.Values{}
	result := swap.httpGet("/dapi/v1/balance", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetUserAssets user account assets
func (swap *SwapCoin) GetUserAssets() interface{} {
	params := &url.Values{}
	result := swap.httpGet("/dapi/v1/account", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetUserPositions user open position
func (swap *SwapCoin) GetUserPositions(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	result := swap.httpGet("/dapi/v1/account", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// PlaceOrder place order
func (swap *SwapCoin) PlaceOrder(order *goex.PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(order.Symbol))
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
		case goex.GTX:
			params.Set("timeInForce", "GTX")
		default:
			params.Set("timeInForce", "GTC")
		}
	} else {
		params.Set("type", strings.ToUpper(goex.MARKET))
	}

	result := swap.httpPost("/dapi/v1/order", params, true)
	return result
}

// PlaceLimitOrder place limit order
func (swap *SwapCoin) PlaceLimitOrder(symbol goex.Symbol, price string, amount string, side goex.TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("timeInForce", "GTC")
	params.Set("type", strings.ToUpper(goex.LIMIT))
	params.Set("side", strings.ToUpper(side.String()))
	if ClientOrderID != "" {
		params.Set("newClientOrderId", ClientOrderID)
	}
	result := swap.httpPost("/dapi/v1/order", params, true)
	return result
}

// PlaceMarketOrder place market order
func (swap *SwapCoin) PlaceMarketOrder(symbol goex.Symbol, amount string, side goex.TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	params.Set("quantity", amount)
	params.Set("type", strings.ToUpper(goex.MARKET))
	params.Set("side", strings.ToUpper(side.String()))
	if ClientOrderID != "" {
		params.Set("newClientOrderId", ClientOrderID)
	}
	result := swap.httpPost("/dapi/v1/order", params, true)
	return result
}

// BatchPlaceLimitOrder batch place limit order
func (swap *SwapCoin) BatchPlaceLimitOrder(orders []goex.LimitOrder) interface{} {
	params := &url.Values{}

	var trustOrders []map[string]interface{}
	for index, item := range orders {
		if index > 4 {
			break
		}
		param := map[string]interface{}{}
		param["symbol"] = swap.getSymbol(item.Symbol)
		param["price"] = item.Price
		param["quantity"] = item.Amount
		param["timeInForce"] = item.TimeInForce
		param["type"] = strings.ToUpper(goex.LIMIT)
		param["side"] = strings.ToUpper(item.Side.String())
		if item.ClientOrderId != "" {
			param["newClientOrderId"] = item.ClientOrderId
		}
		switch item.TimeInForce {
		case goex.IOC:
			param["timeInForce"] = "IOC"
		case goex.FOK:
			param["timeInForce"] = "FOK"
		case goex.GTX:
			param["timeInForce"] = "GTX"
		default:
			param["timeInForce"] = "GTC"
		}
		trustOrders = append(trustOrders, param)
	}
	jsonBody, _ := json.Marshal(trustOrders)
	params.Set("batchOrders", string(jsonBody))

	result := swap.httpPost("/dapi/v1/batchOrders", params, true)
	return result
}

// CancelOrder cancel user trust order
func (swap *SwapCoin) CancelOrder(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	if clientOrderID != "" {
		params.Set("origClientOrderId", clientOrderID)
	} else {
		params.Set("orderId", orderID)
	}
	result := swap.httpDelete("/dapi/v1/order", params, true)
	return result
}

// BatchCancelOrder batch cancel trust order
func (swap *SwapCoin) BatchCancelOrder(symbol goex.Symbol, orderIds, clientOrderIds string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	if clientOrderIds != "" {
		params.Set("origClientOrderIdList", fmt.Sprintf("[%s]", clientOrderIds))
	} else {
		params.Set("orderIdList", fmt.Sprintf("[%s]", orderIds))
	}
	result := swap.httpDelete("/dapi/v1/batchOrders", params, true)
	return result
}

// BatchCancelAllOrder batch cancel all orders
func (swap *SwapCoin) BatchCancelAllOrder(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	result := swap.httpDelete("/dapi/v1/allOpenOrders", params, true)
	return result
}

// GetUserOpenTrustOrders user open trust order list
func (swap *SwapCoin) GetUserOpenTrustOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	result := swap.httpGet("/dapi/v1/openOrders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserOrderInfo user trust order info
func (swap *SwapCoin) GetUserOrderInfo(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	if clientOrderID != "" {
		params.Set("origClientOrderId", clientOrderID)
	} else {
		params.Set("orderId", orderID)
	}

	result := swap.httpGet("/dapi/v1/order", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserTradeOrders user trade order list
func (swap *SwapCoin) GetUserTradeOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))

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

	result := swap.httpGet("/dapi/v1/income", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserTrustOrders user trust order list
func (swap *SwapCoin) GetUserTrustOrders(symbol goex.Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))

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

	result := swap.httpGet("/dapi/v1/allOrders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserAssetsIncomes user assets changes records
func (swap *SwapCoin) GetUserAssetsIncomes(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	if symbol.CoinFrom != "" {
		params.Set("symbol", swap.getSymbol(symbol))
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
	if incomeType, ok := options["incomeType"]; ok == true {
		params.Set("incomeType", incomeType)
	}

	result := swap.httpGet("/dapi/v1/income", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserCommissionRate user current commission rate
func (swap *SwapCoin) GetUserCommissionRate(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", swap.getSymbol(symbol))
	result := swap.httpGet("/dapi/v1/commissionRate", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// HTTPRequest request url
func (swap *SwapCoin) HTTPRequest(requestURL, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}
	switch method {
	case goex.HTTP_GET:
		return swap.httpGet(requestURL, params, signed)
	case goex.HTTP_POST:
		return swap.httpPost(requestURL, params, signed)
	case goex.HTTP_DELETE:
		return swap.httpPost(requestURL, params, signed)
	}
	return nil
}

// httpGet Get request method
func (swap *SwapCoin) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	if signed {
		headers["X-MBX-APIKEY"] = swap.accessKey
		swap.sign(params)
	}

	requestURL := swap.baseURL + url
	if params != nil {
		reqData := params.Encode()
		requestURL = requestURL + "?" + reqData
	}
	responseMap = goex.HttpGetWithHeader(swap.httpClient, requestURL, headers)
	return swap.handlerResponse(&responseMap)
}

// httpGet Post request method
func (swap *SwapCoin) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	headers["X-MBX-APIKEY"] = swap.accessKey
	swap.sign(params)
	requestURL := swap.baseURL + url
	responseMap = goex.HttpPostWithHeader(swap.httpClient, requestURL, params.Encode(), headers)
	return swap.handlerResponse(&responseMap)
}

// httpGet Delete request method
func (swap *SwapCoin) httpDelete(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	headers["X-MBX-APIKEY"] = swap.accessKey
	swap.sign(params)
	requestURL := swap.baseURL + url + "?" + params.Encode()
	responseMap = goex.HttpDeleteWithHeader(swap.httpClient, requestURL, headers)
	return swap.handlerResponse(&responseMap)
}

// httpGet Handler response data format
func (swap *SwapCoin) handlerResponse(responseMap *goex.HttpClientResponse) map[string]interface{} {
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
func (swap *SwapCoin) sign(params *url.Values) {
	timestamp := goex.GetNowMillisecondStr()
	params.Set("recvWindow", "5000")
	params.Set("timestamp", timestamp)
	sign, _ := goex.HmacSha256Signer(params.Encode(), swap.secretKey)
	params.Set("signature", sign)
}

// httpGet format symbol method
func (swap SwapCoin) getSymbol(symbol goex.Symbol) string {
	return symbol.ToUpper().ToSymbol("") + "_PERP"
}
