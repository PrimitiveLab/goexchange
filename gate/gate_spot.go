package gate

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

const (
	GATE_BUY  string = "buy"
	GATE_SELL string = "sell"
)

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
	return spot.httpGet("/api2/1/coininfo", params, false)
}

// GetSymbolList exchange all symbol
func (spot *GateSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	return spot.httpGet(spot.getURL("currency_pairs"), params, false)
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
	result := spot.httpGet(spot.getURL("order_book"), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetTicker symbol ticker
func (spot *GateSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", symbol.ToUpper().ToSymbol("_"))
	result := spot.httpGet(spot.getURL("tickers"), params, false)
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

	result := spot.httpGet(spot.getURL("candlesticks"), params, false)
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
	if lastID, ok := options["lastId"]; ok == true {
		params.Set("last_id", lastID)
	}
	result := spot.httpGet(spot.getURL("trades"), params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// GetUserBalance user balance
func (spot *GateSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	return spot.httpGet(spot.getURL("accounts"), params, true)
}

// PlaceOrder place order
func (spot *GateSpot) PlaceOrder(order *PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(order.Symbol))
	params.Set("amount", order.Amount)
	params.Set("amount", order.Price)
	params.Set("type", "limit")
	params.Set("account", "spot")
	params.Set("side", order.Side.String())

	switch order.TimeInForce {
	case IOC:
		params.Set("time_in_force", "ioc")
	case GTC:
		params.Set("time_in_force", "gtc")
	case POC:
		params.Set("time_in_force", "poc")
	}
	if order.ClientOrderId != "" {
		params.Set("text", "t-"+order.ClientOrderId)
	}

	result := spot.httpPost(spot.getURL("orders"), params, false)
	return result
}

// PlaceLimitOrder place limit order
func (spot *GateSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("amount", amount)
	params.Set("type", "limit")
	params.Set("account", "spot")
	params.Set("side", side.String())
	if ClientOrderID != "" {
		params.Set("text", "t-"+ClientOrderID)
	}
	return spot.httpPost(spot.getURL("orders"), params, true)
}

// PlaceMarketOrder place market order
func (spot *GateSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderID string) interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// BatchPlaceLimitOrder batch place limit order
func (spot *GateSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	var params []map[string]interface{}
	for index, item := range orders {
		param := map[string]interface{}{}
		param["currency_pair"] = spot.getSymbol(item.Symbol)
		param["price"] = item.Price
		param["amount"] = item.Amount
		param["side"] = item.Side.String()
		param["type"] = LIMIT
		param["account"] = "spot"
		if item.ClientOrderId != "" {
			param["text"] = "t-" + item.ClientOrderId
		} else {
			param["text"] = "t-" + GetNowMillisecondStr() + strconv.FormatInt(int64(index), 10)
		}
		switch item.TimeInForce {
		case IOC:
			param["time_in_force"] = "ioc"
		case GTC:
			param["time_in_force"] = "gtc"
		case POC:
			param["time_in_force"] = "poc"
		}
		params = append(params, param)
	}
	return spot.httpPostBatch(spot.getURL("batch_orders"), params, true)
}

// CancelOrder cancel a order
func (spot *GateSpot) CancelOrder(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("order_id", orderID)
	params.Set("currency_pair", spot.getSymbol(symbol))

	return spot.httpDelete(spot.getURL("orders/"+orderID), params, true)
}

// BatchCancelOrder batch cancel orders
func (spot *GateSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	var params []map[string]interface{}
	orderIDList := strings.Split(orderIds, ",")
	for _, item := range orderIDList {
		param := map[string]interface{}{}
		param["currency_pair"] = spot.getSymbol(symbol)
		param["id"] = item
		params = append(params, param)
	}
	return spot.httpPostBatch(spot.getURL("cancel_batch_orders"), params, true)
}

// BatchCancelAllOrder batch cancel all orders
func (spot *GateSpot) BatchCancelAllOrder(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	params.Set("account", "spot")
	result := spot.httpDelete(spot.getURL("orders"), params, true)
	return result
}

// GetUserOpenTrustOrders get current trust order
func (spot *GateSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	params.Set("status", "open")
	params.Set("limit", strconv.FormatInt(int64(size), 10))
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	result := spot.httpGet(spot.getURL("open_orders"), params, true)
	return result
}

// GetUserOrderInfo get trust order info
func (spot *GateSpot) GetUserOrderInfo(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	params.Set("order_id", orderID)
	result := spot.httpGet(spot.getURL("orders/"+orderID), params, true)
	return result
}

// GetUserTradeOrders get trade order list
func (spot *GateSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	params.Set("limit", strconv.FormatInt(int64(size), 10))
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	if orderID, ok := options["order_id"]; ok == true {
		params.Set("order_id", orderID)
	}

	result := spot.httpGet(spot.getURL("my_trades"), params, true)
	return result
}

// GetUserTrustOrders get trust order list
func (spot *GateSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("currency_pair", spot.getSymbol(symbol))
	if status != "" {
		params.Set("status", status)
	} else {
		params.Set("status", "finished")
	}
	params.Set("limit", strconv.FormatInt(int64(size), 10))
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	result := spot.httpGet(spot.getURL("orders"), params, true)
	return result
}

func (spot *GateSpot) HttpRequest(requestURL, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}

	switch method {
	case HTTP_GET:
		return spot.httpGet(spot.getURL(requestURL), params, signed)
	case HTTP_POST:
		return spot.httpPost(spot.getURL(requestURL), params, signed)
	case HTTP_DELETE:
		return spot.httpDelete(spot.getURL(requestURL), params, signed)
	}
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
		headers["SIGN"] = spot.sign(url, HTTP_GET, timestamp, params.Encode())
		headers["Timestamp"] = timestamp
	}
	requestURL := spot.baseURL + url + "?" + params.Encode()
	responseMap = HttpGetWithHeader(spot.httpClient, requestURL, headers)
	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *GateSpot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	headers := map[string]string{}
	bodyMap := map[string]string{}
	for key, item := range *params {
		bodyMap[key] = item[0]
	}
	jsonBody, _ := json.Marshal(bodyMap)
	if signed {
		timestamp := GetNowTimestampStr()
		headers["KEY"] = spot.accessKey
		headers["SIGN"] = spot.sign(url, HTTP_POST, timestamp, string(jsonBody))
		headers["Timestamp"] = timestamp
	}
	requestURL := spot.baseURL + url
	responseMap = HttpPostWithJson(spot.httpClient, requestURL, string(jsonBody), headers)

	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *GateSpot) httpPostBatch(url string, params interface{}, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	headers := map[string]string{}
	jsonBody, _ := json.Marshal(params)
	if signed {
		timestamp := GetNowTimestampStr()
		headers["KEY"] = spot.accessKey
		headers["SIGN"] = spot.sign(url, HTTP_POST, timestamp, string(jsonBody))
		headers["Timestamp"] = timestamp
	}
	requestURL := spot.baseURL + url
	responseMap = HttpPostWithJson(spot.httpClient, requestURL, string(jsonBody), headers)

	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *GateSpot) httpDelete(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	headers := map[string]string{}

	if signed {
		timestamp := GetNowTimestampStr()
		headers["KEY"] = spot.accessKey
		headers["SIGN"] = spot.sign(url, HTTP_DELETE, timestamp, params.Encode())
		headers["Timestamp"] = timestamp
	}
	requestURL := spot.baseURL + url + "?" + params.Encode()
	responseMap = HttpDeleteWithHeader(spot.httpClient, requestURL, headers)
	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
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

	var bodyData interface{}
	err := json.Unmarshal(responseMap.Data, &bodyData)
	if err != nil {
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}
	returnData["data"] = bodyData
	return returnData
}

func (spot *GateSpot) sign(url, method, timestamp, bodyData string) string {
	queryStr := ""
	hashPayload := ""
	if method == HTTP_POST {
		hashPayload = Sha512Signer(bodyData)
	} else {
		queryStr = bodyData
		hashPayload = Sha512Signer("")
	}

	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, url, queryStr, hashPayload, timestamp)
	sign, _ := HmacSha512Signer(signStr, spot.secretKey)
	return sign
}

func (spot GateSpot) getSymbol(symbol Symbol) string {
	return symbol.ToUpper().ToSymbol("_")
}

func (spot *GateSpot) getURL(url string) string {
	return "/api/v4/spot/" + url
}
