package huobi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	goex "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	goex.KLINE_PERIOD_1MINUTE:  "1min",
	goex.KLINE_PERIOD_5MINUTE:  "5min",
	goex.KLINE_PERIOD_15MINUTE: "15min",
	goex.KLINE_PERIOD_30MINUTE: "30min",
	goex.KLINE_PERIOD_60MINUTE: "60min",
	goex.KLINE_PERIOD_1DAY:     "1day",
	goex.KLINE_PERIOD_1WEEK:    "1week",
	goex.KLINE_PERIOD_1MONTH:   "1mon",
	goex.KLINE_PERIOD_1YEAR:    "1year",
}

const (
	HUOBI_SPOT_ACCOUNT = "spot"
)

// Spot huobi struct
type Spot struct {
	httpClient *http.Client
	baseURL    string
	accountId  string
	accessKey  string
	secretKey  string
}

// New new instance
func New(client *http.Client, baseURL, apiKey, secretKey, accountID string) *Spot {
	instance := new(Spot)
	if baseURL == "" {
		instance.baseURL = "https://api.huobi.pro"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	instance.accountId = accountID
	return instance
}

// NewWithConfig new instance with config struct
func NewWithConfig(config *goex.APIConfig) *Spot {
	instance := new(Spot)
	if config.Endpoint == "" {
		instance.baseURL = "https://api.huobi.pro"
	} else {
		instance.baseURL = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	instance.accountId = config.AccountId
	return instance
}

// GetExchangeName get exchange name
func (spot *Spot) GetExchangeName() string {
	return goex.EXCHANGE_HUOBI
}

// GetCoinList exchange coin list
func (spot *Spot) GetCoinList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/v2/reference/currencies", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetSymbolList exchange symbol list
func (spot *Spot) GetSymbolList() interface{} {

	params := &url.Values{}
	result := spot.httpGet("/v1/common/symbols", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetDepth exchange depth data
func (spot *Spot) GetDepth(symbol goex.Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))

	// depth ranges [5,10,20]
	if size < 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else {
		size = 20
	}
	params.Set("depth", strconv.Itoa(size))

	if depthType, ok := options["type"]; ok {
		params.Set("type", depthType)
	}

	result := spot.httpGet("/market/depth", params, false)
	fmt.Println(result)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["tick"]
	return result
}

// GetTicker exchange ticker data
func (spot *Spot) GetTicker(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpGet("/market/detail/merged", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["tick"]
	return result
}

// GetKline exchange kline data
func (spot *Spot) GetKline(symbol goex.Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	periodStr, isOk := klinePeriod[period]
	if isOk != true {
		periodStr = "1min"
	}
	params.Set("period", periodStr)
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	result := spot.httpGet("/market/history/kline", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetTrade exchange trade order data
func (spot *Spot) GetTrade(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	result := spot.httpGet("/market/history/trade", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetUserBalance user account balance
func (spot *Spot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet(fmt.Sprintf("/v1/account/accounts/%s/balance", spot.accountId), params, true)

	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserCommissionRate user current commission rate
func (spot *Spot) GetUserCommissionRate(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbols", spot.getSymbol(symbol))
	result := spot.httpGet("/v2/reference/transact-fee-rate", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// PlaceOrder place order
func (spot *Spot) PlaceOrder(order *goex.PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("account-id", spot.accountId)
	params.Set("symbol", spot.getSymbol(order.Symbol))
	params.Set("price", order.Price)
	params.Set("amount", order.Amount)
	params.Set("source", "spot-api")
	if order.ClientOrderId != "" {
		params.Set("client-order-id", order.ClientOrderId)
	}
	side := order.Side.String()
	tradeType := ""
	switch order.TimeInForce {
	case goex.IOC:
		tradeType = "ioc"
	case goex.FOK:
		tradeType = "limit-fok"
	case goex.POC:
		tradeType = "limit-maker"
	default:
		tradeType = order.TradeType
	}
	params.Set("type", fmt.Sprintf("%s-%s", side, tradeType))
	result := spot.httpPost("/v1/order/orders/place", params, true)
	return result
}

// PlaceLimitOrder place limit order
func (spot *Spot) PlaceLimitOrder(symbol goex.Symbol, price string, amount string, side goex.TradeSide, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("account-id", spot.accountId)
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("amount", amount)
	if side == goex.BUY {
		params.Set("type", "buy-limit")
	} else {
		params.Set("type", "sell-limit")
	}
	params.Set("source", "spot-api")
	if clientOrderId != "" {
		params.Set("client-order-id", clientOrderId)
	}
	result := spot.httpPost("/v1/order/orders/place", params, true)
	return result
}

// PlaceMarketOrder place market order
func (spot *Spot) PlaceMarketOrder(symbol goex.Symbol, amount string, side goex.TradeSide, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("account-id", spot.accountId)
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("amount", amount)
	if side == goex.BUY {
		params.Set("type", "buy-market")
	} else {
		params.Set("type", "sell-market")
	}
	params.Set("source", "spot-api")
	if clientOrderId != "" {
		params.Set("client-order-id", clientOrderId)
	}
	result := spot.httpPost("/v1/order/orders/place", params, true)
	return result
}

// BatchPlaceLimitOrder batch place limit order
func (spot *Spot) BatchPlaceLimitOrder(orders []goex.LimitOrder) interface{} {
	var trustOrders []map[string]interface{}
	for _, item := range orders {
		param := map[string]interface{}{}

		param["account-id"] = spot.accountId
		param["symbol"] = spot.getSymbol(item.Symbol)
		param["price"] = item.Price
		param["amount"] = item.Amount
		if item.Side == goex.BUY {
			param["type"] = "buy-limit"
		} else {
			param["type"] = "sell-limit"
		}
		param["source"] = "spot-api"
		param["client-order-id"] = item.ClientOrderId
		trustOrders = append(trustOrders, param)
	}
	result := spot.httpPostBatch("/v1/order/batch-orders", trustOrders, true)
	return result
}

// CancelOrder cancel user trust order
func (spot *Spot) CancelOrder(symbol goex.Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	var result map[string]interface{}
	if clientOrderId != "" {
		params.Set("client-order-id", clientOrderId)
		result = spot.httpPost("/v1/order/orders/submitCancelClientOrder", params, true)
	} else {
		params.Set("order-id", orderId)
		result = spot.httpPost(fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId), params, true)
	}
	return result
}

// BatchCancelOrder batch cancel trust order
func (spot *Spot) BatchCancelOrder(symbol goex.Symbol, orderIds, clientOrderIds string) interface{} {
	params := map[string]interface{}{}
	if clientOrderIds != "" {
		params["client-order-ids"] = strings.Split(clientOrderIds, ",")
	} else {
		params["order-ids"] = strings.Split(orderIds, ",")
	}
	return spot.httpPostBatch("/v1/order/orders/batchcancel", params, true)
}

// BatchCancelAllOrder batch cancel all orders
func (spot *Spot) BatchCancelAllOrder(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	params.Set("account-id", spot.accountId)
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpPost("/v1/order/orders/batchCancelOpenOrders", params, true)
	return result
}

// GetUserOpenTrustOrders user open trust order list
func (spot *Spot) GetUserOpenTrustOrders(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("account-id", spot.accountId)
	params.Set("symbol", spot.getSymbol(symbol))
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	if side, ok := options["side"]; ok {
		params.Set("side", side)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if direct, ok := options["direct"]; ok {
		params.Set("direct", direct)
	}
	result := spot.httpGet("/v1/order/openOrders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserOrderInfo user trust order info
func (spot *Spot) GetUserOrderInfo(symbol goex.Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	var result map[string]interface{}
	if clientOrderID == "" {
		result = spot.httpGet(fmt.Sprintf("/v1/order/orders/%s", orderID), params, true)
	} else {
		params.Set("clientOrderId", clientOrderID)
		result = spot.httpGet("/v1/order/orders/getClientOrder", params, true)
	}
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
		params.Set("size", strconv.Itoa(size))
	}
	if types, ok := options["types"]; ok {
		params.Set("types", types)
	}
	if startTime, ok := options["start-time"]; ok {
		params.Set("start-time", startTime)
	}
	if endTime, ok := options["end-time"]; ok {
		params.Set("end-time", endTime)
	}
	if startDate, ok := options["start-date"]; ok {
		params.Set("start-date", startDate)
	}
	if endDate, ok := options["end-date"]; ok {
		params.Set("end-date", endDate)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if direct, ok := options["direct"]; ok {
		params.Set("direct", direct)
	}
	result := spot.httpGet("/v1/order/matchresults", params, true)
	if result["code"] != 0 {
		return result
	}
	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetUserTradeOrders user trust order list
func (spot *Spot) GetUserTrustOrders(symbol goex.Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("states", status)
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	if status != "" {
		params.Set("states", status)
	} else {
		params.Set("states", "created,submitted,partial-filled,filled,partial-canceled,canceling,canceled")
	}
	if types, ok := options["types"]; ok {
		params.Set("types", types)
	}
	if startTime, ok := options["start-time"]; ok {
		params.Set("start-time", startTime)
	}
	if endTime, ok := options["end-time"]; ok {
		params.Set("end-time", endTime)
	}
	if startDate, ok := options["start-date"]; ok {
		params.Set("start-date", startDate)
	}
	if endDate, ok := options["end-date"]; ok {
		params.Set("end-date", endDate)
	}
	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if direct, ok := options["direct"]; ok {
		params.Set("direct", direct)
	}
	result := spot.httpGet("/v1/order/orders", params, true)
	if result["code"] != 0 {
		return result
	}
	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

// GetUserDepositAddress user deposit address
func (spot *Spot) GetUserDepositAddress(coin string, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency", coin)

	result := spot.httpGet("/v2/account/deposit/address", params, true)
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
	if fee, ok := options["fee"]; ok {
		params.Set("fee", fee)
	}
	if tag != "" {
		params.Set("addr-tag", tag)
	}
	if chain != "" {
		params.Set("chain", chain)
	}
	result := spot.httpPost("/v1/dw/withdraw/api/create", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserDepositRecords user deposit record list
func (spot *Spot) GetUserDepositRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("type", "deposit")
	if coin != "" {
		params.Set("currency", coin)
	}
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}

	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if direct, ok := options["direct"]; ok {
		params.Set("direct", direct)
	}
	result := spot.httpGet("/v1/query/deposit-withdraw", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetUserWithdrawRecords user withdraw record list
func (spot *Spot) GetUserWithdrawRecords(coin string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("type", "withdraw")
	if coin != "" {
		params.Set("currency", coin)
	}
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}

	if from, ok := options["from"]; ok {
		params.Set("from", from)
	}
	if direct, ok := options["direct"]; ok {
		params.Set("direct", direct)
	}
	result := spot.httpGet("/v1/query/deposit-withdraw", params, true)
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
	}
	return nil
}

// httpGet Get request method
func (spot *Spot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	sign := ""
	if signed {
		sign = spot.sign(goex.HTTP_GET, url, params)
	}

	requestURL := spot.baseURL + url
	if params != nil {
		requestURL = requestURL + "?" + params.Encode()
		if sign != "" {
			requestURL = requestURL + "&Signature=" + sign
		}
	}
	responseMap = goex.HttpGet(spot.httpClient, requestURL)
	return spot.handlerResponse(&responseMap)
}

// httpPost Post request method
func (spot *Spot) httpPost(path string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse

	signParams := &url.Values{}
	sign := spot.sign(goex.HTTP_POST, path, signParams)
	requestURL := spot.baseURL + path + "?" + signParams.Encode() + "&Signature=" + sign

	bodyMap := map[string]string{}
	for key, item := range *params {
		bodyMap[key] = item[0]
	}
	jsonBody, _ := json.Marshal(bodyMap)
	responseMap = goex.HttpPostWithJson(spot.httpClient, requestURL, string(jsonBody), map[string]string{})
	return spot.handlerResponse(&responseMap)
}

func (spot *Spot) httpPostBatch(path string, params interface{}, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	jsonBody, _ := json.Marshal(params)

	signParams := &url.Values{}
	sign := spot.sign(goex.HTTP_POST, path, signParams)
	requestURL := spot.baseURL + path + "?" + signParams.Encode() + "&Signature=" + sign
	responseMap = goex.HttpPostWithJson(spot.httpClient, requestURL, string(jsonBody), map[string]string{})
	return spot.handlerResponse(&responseMap)
}

// handlerResponse Handler response data format
func (spot *Spot) handlerResponse(responseMap *goex.HttpClientResponse) map[string]interface{} {
	retData := make(map[string]interface{})

	retData["code"] = responseMap.Code
	retData["st"] = responseMap.St
	retData["et"] = responseMap.Et
	if responseMap.Code != 0 {
		retData["msg"] = responseMap.Msg
		retData["error"] = responseMap.Error
		return retData
	}

	var bodyDataMap map[string]interface{}
	err := json.Unmarshal(responseMap.Data, &bodyDataMap)
	if err != nil {
		retData["code"] = goex.JsonUnmarshalError.Code
		retData["msg"] = goex.JsonUnmarshalError.Msg
		retData["error"] = err.Error()
		return retData
	}

	if status, ok := bodyDataMap["status"]; ok && status.(string) != "ok" {
		retData["code"] = goex.ExchangeError.Code
		retData["msg"] = goex.ExchangeError.Msg
		retData["error"] = bodyDataMap["err-msg"].(string)
		return retData
	}
	if code, ok := bodyDataMap["code"]; ok && code.(float64) != 200 {
		retData["code"] = goex.ExchangeError.Code
		retData["msg"] = goex.ExchangeError.Msg
		retData["error"] = bodyDataMap["message"].(string)
		return retData
	}

	retData["data"] = bodyDataMap
	return retData
}

// sign signature method
func (spot *Spot) sign(method, path string, params *url.Values) string {
	host, _ := url.Parse(spot.baseURL)
	params.Set("AccessKeyId", spot.accessKey)
	params.Set("SignatureMethod", "HmacSHA256")
	params.Set("SignatureVersion", "2")
	params.Set("Timestamp", goex.GetNowUtcTime())
	parameters := params.Encode()

	var sb strings.Builder
	sb.WriteString(method)
	sb.WriteString("\n")
	sb.WriteString(host.Host)
	sb.WriteString("\n")
	sb.WriteString(path)
	sb.WriteString("\n")
	sb.WriteString(parameters)

	sign, _ := goex.HmacSha256Base64Signer(sb.String(), spot.secretKey)
	return sign
}

// getSymbol format symbol method
func (spot Spot) getSymbol(symbol goex.Symbol) string {
	return symbol.ToSymbol("")
}
