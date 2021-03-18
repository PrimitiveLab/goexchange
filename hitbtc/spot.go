package hitbtc

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	goex "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	goex.KLINE_PERIOD_1MINUTE:  "M1",
	goex.KLINE_PERIOD_3MINUTE:  "M3",
	goex.KLINE_PERIOD_5MINUTE:  "M5",
	goex.KLINE_PERIOD_15MINUTE: "M15",
	goex.KLINE_PERIOD_30MINUTE: "M30",
	goex.KLINE_PERIOD_60MINUTE: "H1",
	goex.KLINE_PERIOD_1HOUR:    "H1",
	goex.KLINE_PERIOD_4HOUR:    "H4",
	goex.KLINE_PERIOD_1DAY:     "D1",
	goex.KLINE_PERIOD_1WEEK:    "D7",
	goex.KLINE_PERIOD_1MONTH:   "1M",
}

// Spot hitbtc struct
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
		instance.baseURL = "https://api.hitbtc.com"
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
		instance.baseURL = "https://api.hitbtc.com"
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
	return goex.EXCHANGE_HITBTC
}

// GetCoinList exchange coin list
func (spot *Spot) GetCoinList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/2/public/currency", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetSymbolList exchange symbol list
func (spot *Spot) GetSymbolList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/2/public/symbol", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetDepth exchange depth data
func (spot *Spot) GetDepth(symbol goex.Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	fmtSymbol := spot.getSymbol(symbol)
	params.Set("symbols", fmtSymbol)
	params.Set("limit", strconv.Itoa(size))
	result := spot.httpGet("/api/2/public/orderbook", params, false)
	if result["code"] != 0 {
		return result
	}
	if data, ok := result["data"].(map[string]interface{})[fmtSymbol]; ok {
		result["data"] = data
	}
	return result
}

// GetTicker exchange ticker data
func (spot *Spot) GetTicker(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbols", spot.getSymbol(symbol))
	result := spot.httpGet("/api/2/public/ticker", params, false)
	if result["code"] != 0 {
		return result
	}

	if len(result["data"].([]interface{})) > 0 {
		result["data"] = result["data"].([]interface{})[0]
	}
	return result
}

// GetKline exchange kline data
func (spot *Spot) GetKline(symbol goex.Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	fmtSymbol := spot.getSymbol(symbol)
	params.Set("symbols", fmtSymbol)
	periodStr, isOk := klinePeriod[period]
	if !isOk {
		periodStr = "M1"
	}
	params.Set("period", periodStr)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	if sort, ok := options["sort"]; ok {
		params.Set("sort", sort)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	result := spot.httpGet("/api/2/public/candles", params, false)
	if result["code"] != 0 {
		return result
	}
	if data, ok := result["data"].(map[string]interface{})[fmtSymbol]; ok {
		result["data"] = data
	}
	return result
}

// GetTrade exchange trade order data
func (spot *Spot) GetTrade(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	fmtSymbol := spot.getSymbol(symbol)
	params.Set("symbols", fmtSymbol)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	if sort, ok := options["sort"]; ok {
		params.Set("sort", sort)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	result := spot.httpGet("/api/2/public/trades", params, false)
	if result["code"] != 0 {
		return result
	}
	if data, ok := result["data"].(map[string]interface{})[fmtSymbol]; ok {
		result["data"] = data
	}
	return result
}

// GetUserBalance user account balance
func (spot *Spot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/2/trading/balance", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserCommissionRate user current commission rate
func (spot *Spot) GetUserCommissionRate(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	result := spot.httpGet("/api/2/trading/fee/"+spot.getSymbol(symbol), params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// PlaceOrder place order
func (spot *Spot) PlaceOrder(order *goex.PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(order.Symbol))
	params.Set("price", order.Price)
	params.Set("quantity", order.Amount)
	params.Set("side", order.Side.String())
	params.Set("type", order.TradeType)
	if order.ClientOrderId != "" {
		params.Set("clientOrderId", order.ClientOrderId)
	}
	result := spot.httpPost("/api/2/order", params, true)
	return result
}

// PlaceLimitOrder place limit order
func (spot *Spot) PlaceLimitOrder(symbol goex.Symbol, price string, amount string, side goex.TradeSide, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("side", side.String())
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	result := spot.httpPost("/api/2/order", params, true)
	return result
}

// PlaceMarketOrder place market order
func (spot *Spot) PlaceMarketOrder(symbol goex.Symbol, amount string, side goex.TradeSide, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("quantity", amount)
	params.Set("timeInForce", "GTC")
	params.Set("type", goex.MARKET)
	params.Set("side", side.String())
	if clientOrderID != "" {
		params.Set("clientOrderId", clientOrderID)
	}
	result := spot.httpPost("/api/2/order", params, true)
	return result
}

// BatchPlaceLimitOrder batch place limit order
func (spot *Spot) BatchPlaceLimitOrder(orders []goex.LimitOrder) interface{} {
	return goex.ReturnAPIError(goex.MethodNotExistError)
}

// CancelOrder cancel user trust order
func (spot *Spot) CancelOrder(symbol goex.Symbol, orderId, clientOrderId string) interface{} {
	result := spot.httpDelete("/api/2/order/"+orderId, nil, true)
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
	result := spot.httpDelete("/api/2/order", params, true)
	return result
}

// GetUserOpenTrustOrders user open trust order list
func (spot *Spot) GetUserOpenTrustOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpGet("/api/2/order", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserOrderInfo user trust order info
func (spot *Spot) GetUserOrderInfo(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("clientOrderId", orderID)
	result := spot.httpGet("/api/2/history/order", params, true)
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
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	result := spot.httpGet("/api/2/history/trades", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserTradeOrders user trust order list
func (spot *Spot) GetUserTrustOrders(symbol goex.Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	result := spot.httpGet("/api/2/history/order", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserDepositAddress user deposit address
func (spot *Spot) GetUserDepositAddress(coin string, options map[string]string) interface{} {
	result := spot.httpGet("/api/2/account/crypto/address/"+coin, nil, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// Withdraw user withdraw
func (spot *Spot) Withdraw(coin, address, tag, amount, chain string, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("address", address)
	params.Set("amount", amount)
	params.Set("currency", coin)
	if includeFee, ok := options["includeFee"]; ok {
		params.Set("includeFee", includeFee)
	}
	if autoCommit, ok := options["autoCommit"]; ok {
		params.Set("autoCommit", autoCommit)
	}
	if publicComment, ok := options["publicComment"]; ok {
		params.Set("publicComment", publicComment)
	}
	if paymentId, ok := options["paymentId"]; ok {
		params.Set("paymentId", paymentId)
	}
	result := spot.httpPost("/api/2/account/crypto/withdraw", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserDepositRecords user deposit record list
func (spot *Spot) GetUserDepositRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency", coin)
	params.Set("showSenders", "true")

	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	if by, ok := options["by"]; ok {
		params.Set("by", by)
	}
	if sort, ok := options["sort"]; ok {
		params.Set("sort", sort)
	}

	result := spot.httpGet("/api/2/account/transactions", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserWithdrawRecords user withdraw record list
func (spot *Spot) GetUserWithdrawRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency", coin)
	params.Set("showSenders", "true")

	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if till, ok := options["till"]; ok {
		params.Set("till", till)
	}
	if offset, ok := options["offset"]; ok {
		params.Set("offset", offset)
	}
	if by, ok := options["by"]; ok {
		params.Set("by", by)
	}
	if sort, ok := options["sort"]; ok {
		params.Set("sort", sort)
	}

	result := spot.httpGet("/api/2/account/transactions", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

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
		return spot.httpDelete(requestURL, params, signed)
	}
	return nil
}

// httpGet Get request method
func (spot *Spot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	headers := map[string]string{}
	if signed {
		headers["Authorization"] = spot.sign()
	}
	requestURL := spot.baseURL + url
	if params != nil {
		requestURL = requestURL + "?" + params.Encode()
	}
	responseMap = goex.HttpGetWithHeader(spot.httpClient, requestURL, headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Post request method
func (spot *Spot) httpPost(path string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse

	headers := map[string]string{}
	headers["Authorization"] = spot.sign()
	requestURL := spot.baseURL + path
	responseMap = goex.HttpPostWithHeader(spot.httpClient, requestURL, params.Encode(), headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Delete request method
func (spot *Spot) httpDelete(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse

	headers := map[string]string{}
	headers["Authorization"] = spot.sign()
	requestURL := spot.baseURL + url
	if params != nil {
		requestURL = requestURL + "?" + params.Encode()
	}

	responseMap = goex.HttpDeleteWithHeader(spot.httpClient, requestURL, headers)
	return spot.handlerResponse(&responseMap)
}

// httpGet Handler response data format
func (spot *Spot) handlerResponse(responseMap *goex.HttpClientResponse) map[string]interface{} {
	// var retData map[string]interface{}
	retData := make(map[string]interface{})

	retData["code"] = responseMap.Code
	retData["st"] = responseMap.St
	retData["et"] = responseMap.Et
	if responseMap.Code != 0 {
		retData["msg"] = responseMap.Msg
		retData["error"] = responseMap.Error
		return retData
	}

	var bodyDataMap interface{}
	err := json.Unmarshal(responseMap.Data, &bodyDataMap)
	if err != nil {
		retData["code"] = goex.JsonUnmarshalError.Code
		retData["msg"] = goex.JsonUnmarshalError.Msg
		retData["error"] = err.Error()
		return retData
	}

	retData["data"] = bodyDataMap
	return retData
}

// sign signature method
func (spot *Spot) sign() string {
	return "Basic " + goex.Base64Signer(spot.accessKey+":"+spot.secretKey)
}

// httpGet format symbol method
func (spot Spot) getSymbol(symbol goex.Symbol) string {
	if symbol.CoinTo == "usdt" {
		symbol.CoinTo = "usd"
	}
	return symbol.ToUpper().ToSymbol("")
}
