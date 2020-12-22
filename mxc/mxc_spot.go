package mxc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	KLINE_PERIOD_1MINUTE:  "1m",
	KLINE_PERIOD_5MINUTE:  "5m",
	KLINE_PERIOD_15MINUTE: "15m",
	KLINE_PERIOD_30MINUTE: "30m",
	KLINE_PERIOD_60MINUTE: "60m",
	KLINE_PERIOD_1HOUR:    "1h",
	KLINE_PERIOD_1DAY:     "1d",
	KLINE_PERIOD_1MONTH:   "1M",
}

const (
	MXC_BUY  string = "BID"
	MXC_SELL string = "ASK"
)

type MxcSpot struct {
	httpClient *http.Client
	baseUrl    string
	accessKey  string
	secretKey  string
}

func New(client *http.Client, baseUrl, apiKey, secretKey string) *MxcSpot {
	instance := new(MxcSpot)
	if baseUrl == "" {
		instance.baseUrl = "https://www.mxc.ceo"
	} else {
		instance.baseUrl = baseUrl
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

func NewWithConfig(config *APIConfig) *MxcSpot {
	instance := new(MxcSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://www.mxc.ceo"
	} else {
		instance.baseUrl = config.Endpoint
	}

	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

func (spot *MxcSpot) GetExchangeName() string {
	return EXCHANGE_MCX
}

func (spot *MxcSpot) GetCoinList() interface{} {
	startTime := time.Now().UnixNano() / 1e6
	retData := map[string]interface{}{
		"code":  MethodNotExistError.Code,
		"st":    startTime,
		"et":    startTime,
		"msg":   MethodNotExistError.Msg,
		"error": MethodNotExistError.Msg,
		"data":  nil}
	return retData
}

func (spot *MxcSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/api/v2/market/symbols", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *MxcSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	params.Set("depth", strconv.Itoa(size))

	result := spot.httpGet("/open/api/v2/market/depth", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *MxcSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	result := spot.httpGet("/open/api/v2/market/ticker", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *MxcSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1m"
	}
	params.Set("interval", periodStr)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["start_time"]; ok == true {
		params.Set("start_time", startTime)
	}
	result := spot.httpGet("/open/api/v2/market/kline", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (spot *MxcSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	result := spot.httpGet("/open/api/v2/market/deals", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 获取余额
func (spot *MxcSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/api/v2/account/info", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 批量下单
func (spot *MxcSpot) PlaceOrder(order *PlaceOrder) interface{} {
	params := &url.Values{}

	result := spot.httpPost("/open/api/v2/account/info", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 下限价单
func (spot *MxcSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("trade_type", MXC_BUY)
	if side == SELL {
		params.Set("trade_type", MXC_SELL)
	}
	params.Set("order_type", "LIMIT_ORDER")
	if ClientOrderId != "" {
		params.Set("client_order_id", ClientOrderId)
	}

	result := spot.httpPost("/open/api/v2/order/place", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 下市价单
func (spot *MxcSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	startTime := time.Now().UnixNano() / 1e6
	retData := map[string]interface{}{
		"code":  MethodNotExistError.Code,
		"st":    startTime,
		"et":    startTime,
		"msg":   MethodNotExistError.Msg,
		"error": MethodNotExistError.Msg,
		"data":  nil}
	return retData
}

// 批量下限价单
func (spot *MxcSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return nil
}

// 撤单
func (spot *MxcSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("order_ids", orderId)
	params.Set("client_order_ids", clientOrderId)

	result := spot.httpDelete("/open/api/v2/order/cancel", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 批量撤单
func (spot *MxcSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	params := &url.Values{}
	params.Set("order_ids", orderIds)
	params.Set("client_order_ids", clientOrderIds)

	result := spot.httpDelete("/open/api/v2/order/cancel", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 我的当前委托单
func (spot *MxcSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["start_time"]; ok == true {
		params.Set("start_time", startTime)
	}
	if tradeType, ok := options["trade_type"]; ok == true {
		params.Set("trade_type", tradeType)
	}
	result := spot.httpGet("/open/api/v2/order/open_orders", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 委托单详情
func (spot *MxcSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("order_ids", orderId)
	result := spot.httpGet("/open/api/v2/order/query", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 我的成交单列表
func (spot *MxcSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["start_time"]; ok == true {
		params.Set("start_time", startTime)
	}
	result := spot.httpGet("/open/api/v2/order/deals", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 我的委托单列表
func (spot *MxcSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToUpper().String())
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["start_time"]; ok == true {
		params.Set("start_time", startTime)
	}
	if tradeType, ok := options["trade_type"]; ok == true {
		params.Set("trade_type", tradeType)
	}
	if state, ok := options["state"]; ok == true {
		params.Set("state", state)
	}
	result := spot.httpGet("/open/api/v2/order/list", params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (spot *MxcSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	return nil
}

func (spot *MxcSpot) httpRequest(url, method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := spot.baseUrl + url + "?" + params.Encode()
		fmt.Println(requestUrl)
		responseMap = HttpGet(spot.httpClient, requestUrl)
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

func (spot *MxcSpot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse

	params.Set("api_key", spot.accessKey)
	params.Set("req_time", GetNowTimestampStr())

	requestUrl := spot.baseUrl + url
	reqData := params.Encode()
	requestUrl = requestUrl + "?" + reqData

	if signed {
		sign := spot.sign(url, HTTP_GET, reqData)
		requestUrl = requestUrl + "&sign=" + sign
	}
	responseMap = HttpGet(spot.httpClient, requestUrl)

	fmt.Println(requestUrl)

	return spot.handlerResponse(&responseMap)
}

func (spot *MxcSpot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse

	params.Set("api_key", spot.accessKey)
	params.Set("req_time", GetNowTimestampStr())

	sign := spot.sign(url, HTTP_POST, params.Encode())

	params.Set("sign", sign)

	requestUrl := spot.baseUrl + url + "?" + params.Encode()

	bodyMap := map[string]string{}
	for key, item := range *params {
		bodyMap[key] = item[0]
	}

	jsonBody, _ := json.Marshal(bodyMap)
	reqData := string(jsonBody)
	headers := map[string]string{}
	responseMap = HttpPostWithJson(spot.httpClient, requestUrl, reqData, headers)

	fmt.Println(requestUrl)

	return spot.handlerResponse(&responseMap)
}

func (spot *MxcSpot) httpDelete(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse

	params.Set("api_key", spot.accessKey)
	params.Set("req_time", GetNowTimestampStr())
	reqData := params.Encode()
	sign := spot.sign(url, HTTP_DELETE, reqData)

	requestUrl := fmt.Sprintf("%s%s?%s&sign=%s", spot.baseUrl, url, reqData, sign)
	responseMap = HttpDelete(spot.httpClient, requestUrl)

	fmt.Println(requestUrl)

	return spot.handlerResponse(&responseMap)
}

func (spot *MxcSpot) sign(url, method, reqData string) string {
	signStr := fmt.Sprintf("%s\n%s\n%s", method, url, reqData)
	sign, _ := HmacSha256Signer(signStr, spot.secretKey)
	return sign
}

func (spot *MxcSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}

	returnData["data"] = bodyDataMap
	return returnData
}

func (spot *MxcSpot) handlerError(retData map[string]interface{}) {
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
