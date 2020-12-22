package okex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string {
	KLINE_PERIOD_1MINUTE:   "60",
	KLINE_PERIOD_3MINUTE:   "180",
	KLINE_PERIOD_5MINUTE:   "300",
	KLINE_PERIOD_15MINUTE:  "900",
	KLINE_PERIOD_30MINUTE:  "1800",
	KLINE_PERIOD_60MINUTE:  "3600",
	KLINE_PERIOD_1HOUR:  	"3600",
	KLINE_PERIOD_2HOUR:  	"7200",
	KLINE_PERIOD_4HOUR:  	"14400",
	KLINE_PERIOD_6HOUR:  	"21600",
	KLINE_PERIOD_12HOUR:  	"43200",
	KLINE_PERIOD_1DAY:   	"86400",
	KLINE_PERIOD_1WEEK:  	"604800",
}

type OkexSpot struct {
	httpClient *http.Client
	baseUrl string
	accountId string
	accessKey string
	secretKey string
	passphrase string
}

func New(client *http.Client, baseUrl string, apiKey, secretKey string, passphrase string) *OkexSpot {
	instance := new(OkexSpot)
	if baseUrl == "" {
		instance.baseUrl = "https://www.okex.com"
	} else {
		instance.baseUrl = baseUrl
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	instance.passphrase = passphrase
	return instance
}

func NewWithConfig(config *APIConfig) *OkexSpot {
	instance := new(OkexSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://www.okex.com"
	} else {
		instance.baseUrl = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	instance.passphrase = config.ApiPassphrase
	return instance
}

// 交易所名称
func (spot *OkexSpot) GetExchangeName() string {
	return EXCHANGE_OKEX
}

// 币种列表
func (spot *OkexSpot) GetCoinList() interface{} {
	result := spot.httpGet("/api/account/v3/currencies",nil,true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 交易对列表
func (spot *OkexSpot) GetSymbolList() interface{} {
	result := spot.httpGet("/api/spot/v3/instruments", nil, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 深度
func (spot *OkexSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := map[string]string{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	params["size"] = strconv.Itoa(size)
	if depthType, ok := options["depth"]; ok == true {
		params["depth"] = depthType
	}

	result := spot.httpGet("/api/spot/v3/instruments/" + instrumentId + "/book", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 牌价
func (spot *OkexSpot) GetTicker(symbol Symbol) interface{} {
	instrumentId := symbol.ToUpper().ToSymbol("-")
	result := spot.httpGet("/api/spot/v3/instruments/" + instrumentId + "/ticker", nil, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// Kline
func (spot *OkexSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := map[string]string{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "60"
	}
	params["granularity"] = periodStr
	if size != 0 {
		params["limit"] = strconv.Itoa(size)
	}
	if startTime, ok := options["startTime"]; ok == true {
		params["start"] = startTime
	}
	if endTime, ok := options["endTime"]; ok == true {
		params["end"] = endTime
	}

	result := spot.httpGet("/api/spot/v3/instruments/" + instrumentId + "/history/candles", params,false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 最新成交
func (spot *OkexSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := map[string]string{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	if size != 0 {
		params["limit"] = strconv.Itoa(size)
	}
	result := spot.httpGet("/api/spot/v3/instruments/" + instrumentId + "/trades", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 获取余额
func (spot *OkexSpot) GetUserBalance() interface{} {
	result := spot.httpGet("/api/spot/v3/accounts", nil, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 批量下单
func (spot *OkexSpot) PlaceOrder(order *PlaceOrder) interface{} {

	params := map[string]interface{}{}
	params["instrument_id"] = order.Symbol.ToUpper().ToSymbol("-")
	if order.ClientOrderId != "" {
		params["client_oid"] = order.ClientOrderId
	}
	params["side"] = order.Side.String()
	if order.TradeType == LIMIT {
		params["price"] = order.Price
		params["amount"] = order.Amount
		switch order.TimeInForce {
		case IOC :
			params["order_type"] = 3
		case FOK:
			params["order_type"] = 2
		case POC:
			params["order_type"] = 1
		}
	} else {
		if order.Side == BUY {
			params["notional"] = order.Amount
		} else {
			params["size"] = order.Amount
		}
	}
	retData := spot.httpPost("/api/spot/v3/orders", params, true)
	spot.handlerError(retData)
	return retData
}

// 下限价单
func (spot *OkexSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := map[string]interface{}{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	params["price"] = price
	params["size"] = amount
	params["side"] = side.String()
	params["type"] = "limit"
	if ClientOrderId != "" {
		params["client_oid"] = ClientOrderId
	}

	retData := spot.httpPost("/api/spot/v3/orders", params, true)
	spot.handlerError(retData)
	return retData
}

// 下市价单
func (spot *OkexSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := map[string]interface{}{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	if side == BUY {
		params["notional"] = amount
	} else {
		params["size"] = amount
	}
	params["side"] = side.String()
	params["type"] = "market"
	if ClientOrderId != "" {
		params["client_oid"] = ClientOrderId
	}

	retData := spot.httpPost("/api/spot/v3/orders", params, true)
	spot.handlerError(retData)
	return retData
}

// 批量下限价单
func (spot *OkexSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {

	var params []map[string]interface{}
	for _, item := range orders {
		param := map[string]interface{}{}
		param["instrument_id"] = item.Symbol.ToUpper().ToSymbol("-")
		param["price"] = item.Price
		param["size"] = item.Amount
		param["side"] = item.Side.String()
		param["type"] = LIMIT
		if item.ClientOrderId != "" {
			param["client_oid"] = item.ClientOrderId
		}
		switch item.TimeInForce {
		case IOC :
			param["order_type"] = 3
		case FOK:
			param["order_type"] = 2
		case POC:
			param["order_type"] = 1
		}
		params = append(params, param)
	}

	return spot.httpPost("/api/spot/v3/batch_orders", params, true)
}

// 撤单
func (spot *OkexSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := map[string]string{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	id := orderId

	params["instrument_id"] = instrumentId
	if clientOrderId != "" {
		params["client_oid"] = clientOrderId
		id = clientOrderId
	} else {
		params["order_id"] = orderId
	}

	retData := spot.httpPost("/api/spot/v3/cancel_orders/" + id, params, true)
	spot.handlerError(retData)
	return retData
}

// 批量撤单
func (spot *OkexSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	param := map[string]interface{}{}
	param["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	if clientOrderIds != "" {
		param["client_oids"] = strings.Split(clientOrderIds, ",")
	} else {
		param["order_ids"] = strings.Split(orderIds, ",")
	}
	params := [1]map[string]interface{}{param}
	return spot.httpPost("/api/spot/v3/cancel_batch_orders", params, true)
}

// 我的当前委托单
func (spot *OkexSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := map[string]string{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	if size != 0 {
		params["limit"] = strconv.Itoa(size)
	}

	if after, ok := options["after"]; ok == true {
		params["after"] = after
	}

	if before, ok := options["before"]; ok == true {
		params["before"] = before
	}

	result := spot.httpGet("/api/spot/v3/orders_pending", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 委托单详情
func (spot *OkexSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := map[string]string{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	id := orderId
	if clientOrderId != "" {
		params["client_oid"] = clientOrderId
		id = clientOrderId
	} else {
		params["order_id"] = orderId
	}

	result := spot.httpGet("/api/spot/v3/orders/" + id, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 我的成交单列表
func (spot *OkexSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := map[string]string{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	if size != 0 {
		params["limit"] = strconv.Itoa(size)
	}

	if orderId, ok := options["order_id"]; ok == true {
		params["order_id"] = orderId
	}

	if after, ok := options["after"]; ok == true {
		params["after"] = after
	}

	if before, ok := options["before"]; ok == true {
		params["before"] = before
	}

	result := spot.httpGet("/api/spot/v3/fills", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 我的委托单列表
func (spot *OkexSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := map[string]string{}
	params["instrument_id"] = symbol.ToUpper().ToSymbol("-")
	params["state"] = status

	if size != 0 {
		params["limit"] = strconv.Itoa(size)
	}

	if after, ok := options["after"]; ok == true {
		params["after"] = after
	}

	if before, ok := options["before"]; ok == true {
		params["before"] = before
	}

	result := spot.httpGet("/api/spot/v3/orders", params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *OkexSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	switch method {
	case HTTP_GET:
		return spot.httpGet(requestUrl, options.(map[string]string), signed)
	case HTTP_POST:
		return spot.httpPost(requestUrl, options, signed)
	}
	return nil
}

func (spot *OkexSpot) httpGet(url string, params map[string]string, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	var headers map[string]string
	requestUrl := spot.baseUrl + url
	reqData := ""
	if params != nil && len(params) > 0 {
		reqData = "?" + BuildParams(params)
		requestUrl = requestUrl + reqData
	}

	if signed {
		timestamp := IsoTime()
		sign := spot.sign(url, HTTP_GET, timestamp, reqData)
		headers = map[string]string {
			"OK-ACCESS-KEY": spot.accessKey,
			"OK-ACCESS-SIGN": sign,
			"OK-ACCESS-PASSPHRASE": spot.passphrase,
			"OK-ACCESS-TIMESTAMP": timestamp,
		}
	}

	responseMap = HttpGetWithHeader(spot.httpClient, requestUrl, headers)

	fmt.Println(requestUrl)

	return spot.handlerResponse(&responseMap)
}

func (spot *OkexSpot) httpPost(url string, params interface{}, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	var headers map[string]string
	requestUrl := spot.baseUrl + url
	reqData := ""
	if params != nil {
		jsonBody, _ := json.Marshal(params)
		reqData = string(jsonBody)
	}

	timestamp := IsoTime()
	sign := spot.sign(url, HTTP_POST, timestamp, reqData)
	headers = map[string]string {
		"OK-ACCESS-KEY": spot.accessKey,
		"OK-ACCESS-SIGN": sign,
		"OK-ACCESS-PASSPHRASE": spot.passphrase,
		"OK-ACCESS-TIMESTAMP": timestamp,
	}

	responseMap = HttpPostWithJson(spot.httpClient, requestUrl, reqData, headers)

	fmt.Println(requestUrl)

	return spot.handlerResponse(&responseMap)
}

func (spot *OkexSpot) sign(url, method, timestamp, reqData string) string {
	signStr := timestamp + method + url + reqData
	sign, _ := HmacSha256Base64Signer(signStr, spot.secretKey)
	return sign
}

func (spot *OkexSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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

func (spot *OkexSpot) handlerError(retData map[string]interface{}) {
	if retData["code"] == 0 {
		data := retData["data"].(map[string]interface{})
		if data["error_code"].(string) != "0" {
			retData["code"] = data["error_code"].(string)
			retData["msg"] = data["error_message"].(string)
			retData["error"] = retData["msg"]
			retData["data"] = nil
		}
	}
}