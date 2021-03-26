package hoo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	KLINE_PERIOD_1MINUTE:  "1Min",
	KLINE_PERIOD_5MINUTE:  "5Min",
	KLINE_PERIOD_15MINUTE: "15Min",
	KLINE_PERIOD_30MINUTE: "30Min",
	KLINE_PERIOD_60MINUTE: "1Hour",
	KLINE_PERIOD_1HOUR:    "1Hour",
	KLINE_PERIOD_1DAY:     "1Day",
}

const (
	MXC_BUY  string = "BID"
	MXC_SELL string = "ASK"
)

type HooSpot struct {
	httpClient *http.Client
	baseUrl    string
	accessKey  string
	secretKey  string
}

func New(client *http.Client, baseUrl, apiKey, secretKey string) *HooSpot {
	instance := new(HooSpot)
	if baseUrl == "" {
		instance.baseUrl = "https://api.hoolgd.com"
	} else {
		instance.baseUrl = baseUrl
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

func NewWithConfig(config *APIConfig) *HooSpot {
	instance := new(HooSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://api.hoolgd.com"
	} else {
		instance.baseUrl = config.Endpoint
	}

	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

func (spot *HooSpot) GetExchangeName() string {
	return EXCHANGE_HOO
}

func (spot *HooSpot) GetCoinList() interface{} {
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

func (spot *HooSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/v1/tickers/market", params, false)
	return result
}

func (spot *HooSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))

	result := spot.httpGet("/open/v1/depth/market", params, false)
	return result
}

func (spot *HooSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/v1/tickers/market", params, false)
	return result
}

func (spot *HooSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1Min"
	}
	params.Set("type", periodStr)
	result := spot.httpGet("/open/v1/kline/market", params, false)
	return result
}

func (spot *HooSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpGet("/open/v1/trade/market", params, false)
	return result
}

// 获取余额
func (spot *HooSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/v1/balance", params, true)
	return result
}

// 批量下单
func (spot *HooSpot) PlaceOrder(order *PlaceOrder) interface{} {
	params := &url.Values{}

	params.Set("symbol", spot.getSymbol(order.Symbol))
	params.Set("price", order.Price)
	params.Set("quantity", order.Amount)
	params.Set("side", "1")
	if order.Side == SELL {
		params.Set("side", "-1")
	}

	result := spot.httpPost("/open/v1/orders/place", params, true)
	return result
}

// 下限价单
func (spot *HooSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("side", "1")
	if side == SELL {
		params.Set("side", "-1")
	}

	result := spot.httpPost("/open/v1/orders/place", params, true)
	return result
}

// 下市价单
func (spot *HooSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// 批量下限价单
func (spot *HooSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return nil
}

// 撤单
func (spot *HooSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("order_id", orderId)
	params.Set("trade_no", clientOrderId)

	result := spot.httpPost("/open/v1/orders/cancel", params, true)
	return result
}

// 批量撤单
func (spot *HooSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	return nil
}

// 我的当前委托单
func (spot *HooSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpGet("/open/v1/orders/last", params, true)
	return result
}

// 委托单详情
func (spot *HooSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("order_id", orderId)

	result := spot.httpGet("/open/v1/orders/detail", params, true)
	return result
}

// 我的成交单列表
func (spot *HooSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 我的委托单列表
func (spot *HooSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if size != 0 {
		params.Set("pagesize", strconv.Itoa(size))
	}
	if start, ok := options["start"]; ok == true {
		params.Set("start", start)
	}
	if end, ok := options["end"]; ok == true {
		params.Set("end", end)
	}

	if pageNum, ok := options["pagenum"]; ok == true {
		params.Set("pagenum", pageNum)
	}

	if side, ok := options["side"]; ok == true {
		params.Set("side", side)
	}

	result := spot.httpGet("/open/v1/orders", params, true)
	return result
}

func (spot *HooSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}

	switch method {
	case HTTP_GET:
		return spot.httpGet(requestUrl, params, signed)
	case HTTP_POST:
		return spot.httpPost(requestUrl, params, signed)
	}
	return nil
}

func (spot *HooSpot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	if signed {
		spot.sign(params)
	}

	requestUrl := spot.baseUrl + url
	if params != nil {
		reqData := params.Encode()
		requestUrl = requestUrl + "?" + reqData
	}

	responseMap = HttpGet(spot.httpClient, requestUrl)
	return spot.handlerResponse(&responseMap)
}

func (spot *HooSpot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse

	// sign := spot.sign(url, HTTP_POST, params.Encode())
	spot.sign(params)
	// params.Set("sign", "")

	requestUrl := spot.baseUrl + url

	responseMap = HttpPost(spot.httpClient, requestUrl, params.Encode())
	return spot.handlerResponse(&responseMap)
}

func (spot *HooSpot) sign(params *url.Values) {
	timestamp := GetNowTimestampStr()

	signParams := url.Values{}
	signParams.Set("client_id", spot.accessKey)
	signParams.Set("nonce", timestamp[3:])
	signParams.Set("ts", timestamp)
	sign, _ := HmacSha256Signer(signParams.Encode(), spot.secretKey)
	params.Set("client_id", spot.accessKey)
	params.Set("nonce", timestamp[3:])
	params.Set("ts", timestamp)
	params.Set("sign", sign)
}

func (spot *HooSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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

	var bodyDataMap map[string]interface{}
	err := json.Unmarshal(responseMap.Data, &bodyDataMap)
	if err != nil {
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}
	if bodyDataMap["code"].(float64) != 0 {
		returnData["code"] = ExchangeError.Code
		returnData["msg"] = bodyDataMap["msg"].(string)
		returnData["error"] = returnData["msg"]
		return returnData
	}
	returnData["data"] = bodyDataMap["data"]
	return returnData
}

func (spot HooSpot) getSymbol(symbol Symbol) string {
	return symbol.ToUpper().ToSymbol("-")
}
