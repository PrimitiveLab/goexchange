package gate

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
	KLINE_PERIOD_5MINUTE:   "5m",
	KLINE_PERIOD_15MINUTE:  "15m",
	KLINE_PERIOD_30MINUTE:  "30m",
	KLINE_PERIOD_60MINUTE:  "1h",
	KLINE_PERIOD_1HOUR:  	"1h",
	KLINE_PERIOD_4HOUR:  	"4h",
	KLINE_PERIOD_8HOUR:  	"8h",
	KLINE_PERIOD_1DAY:   	"1d",
	KLINE_PERIOD_7DAY:   	"7d",
}

type GateSpot struct {
	httpClient *http.Client
	baseUrl string
	accessKey string
	secretKey string
}

func New(client *http.Client, apiKey, secretKey string) *GateSpot {
	instance := new(GateSpot)
	instance.baseUrl = "https://api.gateio.ws"
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

func NewWithConfig(config *APIConfig) *GateSpot {
	instance := new(GateSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://api.gateio.ws"
	} else {
		instance.baseUrl = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

func (gateSpot *GateSpot) GetExchangeName() string {
	return ECHANGE_GATE
}

func (gateSpot *GateSpot) GetCoinList() interface{} {
	params := &url.Values{}
	result := gateSpot.httpRequest("/api2/1/coininfo", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (gateSpot *GateSpot) GetSymbolList() interface{} {

	params := &url.Values{}
	result := gateSpot.httpRequest(gateSpot.getUrl("currency_pairs"),"get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (gateSpot *GateSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if depthType, ok := options["depth"]; ok == true {
		params.Set("interval", depthType)
	}
	result := gateSpot.httpRequest(gateSpot.getUrl("order_book"), "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (gateSpot *GateSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	result := gateSpot.httpRequest(gateSpot.getUrl("tickers"), "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (gateSpot *GateSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1m"
	}
	params.Set("interval", periodStr)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("from", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("to", endTime)
	}

	result := gateSpot.httpRequest(gateSpot.getUrl("candlesticks"), "get", params,false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (gateSpot *GateSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if lastId, ok := options["lastId"]; ok == true {
		params.Set("last_id", lastId)
	}
	result := gateSpot.httpRequest(gateSpot.getUrl("trades"), "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 获取余额
func (spot *GateSpot) GetUserBalance() interface{} {
	return nil
}

// 批量下单
func (spot *GateSpot) PlaceOrder(order *PlaceOrder) interface{} {
	return nil
}

// 下限价单
func (spot *GateSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// 下市价单
func (spot *GateSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// 批量下限价单
func (spot *GateSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return nil
}

// 撤单
func (spot *GateSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// 批量撤单
func (spot *GateSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	return nil
}

// 我的当前委托单
func (spot *GateSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 委托单详情
func (spot *GateSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// 我的成交单列表
func (spot *GateSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 我的委托单列表
func (spot *GateSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	return nil
}

func (gateSpot *GateSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	return nil
}

func (gateSpot *GateSpot) httpRequest(url , method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := gateSpot.baseUrl + url + "?" + params.Encode()
		fmt.Println(requestUrl)
		responseMap = HttpGet(gateSpot.httpClient, requestUrl)
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

func (gateSpot *GateSpot) getUrl(url string) string  {
	return "/api/v4/spot/" + url
}