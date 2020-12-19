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

var klinePeriod = map[int]string{
	KLINE_PERIOD_1MINUTE:  "1m",
	KLINE_PERIOD_5MINUTE:  "5m",
	KLINE_PERIOD_15MINUTE: "15m",
	KLINE_PERIOD_30MINUTE: "30m",
	KLINE_PERIOD_60MINUTE: "1h",
	KLINE_PERIOD_1HOUR:    "1h",
	KLINE_PERIOD_4HOUR:    "4h",
	KLINE_PERIOD_8HOUR:    "8h",
	KLINE_PERIOD_1DAY:     "1d",
	KLINE_PERIOD_7DAY:     "7d",
}

// GateSpot gate exchange spot
type GateSpot struct {
	httpClient *http.Client
	baseURL    string
	accessKey  string
	secretKey  string
}

// New new instance
func New(client *http.Client, baseURL, apiKey, secretKey string) *GateSpot {
	instance := new(GateSpot)
	if baseURL == "" {
		instance.baseURL = "https://api.gateio.ws"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

// NewWithConfig new instance by config
func NewWithConfig(config *APIConfig) *GateSpot {
	instance := new(GateSpot)
	if config.Endpoint == "" {
		instance.baseURL = "https://api.gateio.ws"
	} else {
		instance.baseURL = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

// GetExchangeName return exchange name
func (spot *GateSpot) GetExchangeName() string {
	return ECHANGE_GATE
}

// GetCoinList exchange supported coins
func (spot *GateSpot) GetCoinList() interface{} {
	params := &url.Values{}
	return spot.httpRequest("/api2/1/coininfo", "get", params, false)
}

// GetSymbolList exchange all symbol
func (spot *GateSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	return spot.httpRequest(spot.getUrl("currency_pairs"), "get", params, false)
}

// GetDepth symbol depth
func (spot *GateSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if depthType, ok := options["depth"]; ok == true {
		params.Set("interval", depthType)
	}
	result := spot.httpRequest(spot.getUrl("order_book"), "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetTicker symbol ticker
func (spot *GateSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	result := spot.httpRequest(spot.getUrl("tickers"), "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetKline symbol kline
func (spot *GateSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
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

	result := spot.httpRequest(spot.getUrl("candlesticks"), "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetTrade symbol last trade
func (spot *GateSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if lastId, ok := options["lastId"]; ok == true {
		params.Set("last_id", lastId)
	}
	result := spot.httpRequest(spot.getUrl("trades"), "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetUserBalance user balance
func (spot *GateSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	return spot.httpRequest(spot.getUrl("accounts"), "get", params, true)
}

// PlaceOrder place order
func (spot *GateSpot) PlaceOrder(order *PlaceOrder) interface{} {
	return nil
}

// PlaceLimitOrder place limit order
func (spot *GateSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// PlaceMarketOrder place market order
func (spot *GateSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	return nil
}

// BatchPlaceLimitOrder batch place limit order
func (spot *GateSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	return nil
}

// CancelOrder cancel a order
func (spot *GateSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// BatchCancelOrder batch cancel orders
func (spot *GateSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	return nil
}

// GetUserOpenTrustOrders get current trust order
func (spot *GateSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// GetUserOrderInfo get trust order info
func (spot *GateSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	return nil
}

// GetUserTradeOrders get trade order list
func (spot *GateSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// GetUserTrustOrders get trust order list
func (spot *GateSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	return nil
}

func (spot *GateSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	return nil
}

func (spot *GateSpot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if signed {
		timestamp := GetNowTimestampStr()
		headers["KEY"] = spot.accessKey
		headers["SIGN"] = spot.sign("", "", "", "")
		headers["Timestamp"] = timestamp
	}

	requestURL := spot.baseURL + url + "?" + params.Encode()
	responseMap = HttpGetWithHeader(spot.httpClient, requestURL, headers)
	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *GateSpot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	requestURL := spot.baseURL + url
	if signed {
		spot.sign("", "", "", "")
	}

	responseMap = HttpPost(spot.httpClient, requestURL, params.Encode())

	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *GateSpot) httpRequest(url, method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := spot.baseURL + url + "?" + params.Encode()
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

func (spot *GateSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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

	var bodyData map[string]interface{}

	err := json.Unmarshal(responseMap.Data, &bodyData)
	fmt.Println(bodyData)
	if err != nil {
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}
	if bodyData["code"].(string) != "0" {
		returnData["code"] = ExchangeError.Code
		returnData["msg"] = ExchangeError.Msg
		returnData["error"] = fmt.Sprintf("code: %v, msg: %v", bodyData["code"], bodyData["msg"])
		return returnData
	}
	returnData["data"] = bodyData["data"]
	return returnData
}

func (spot *GateSpot) sign(url, method, timestamp, reqData string) string {
	signStr := timestamp + method + url + reqData
	sign, _ := HmacSha256Base64Signer(signStr, spot.secretKey)
	return sign
}

func (spot GateSpot) getSymbol(symbol Symbol) string {
	return symbol.ToUpper().ToSymbol("_")
}

func (spot *GateSpot) getUrl(url string) string {
	return "/api/v4/spot/" + url
}
