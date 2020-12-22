package biki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	. "github.com/primitivelab/goexchange"
)

var klinePeriod = map[int]string{
	KLINE_PERIOD_1MINUTE:  "1",
	KLINE_PERIOD_5MINUTE:  "5",
	KLINE_PERIOD_15MINUTE: "15",
	KLINE_PERIOD_30MINUTE: "30",
	KLINE_PERIOD_60MINUTE: "60",
	KLINE_PERIOD_1HOUR:    "60",
	KLINE_PERIOD_1DAY:     "1440",
	KLINE_PERIOD_1WEEK:    "10080",
	KLINE_PERIOD_1MONTH:   "43200",
}

const (
	BIKI_BUY  string = "BUY"
	BIKI_SELL string = "SELL"
)

// BikiSpot biki exchange spot
type BikiSpot struct {
	httpClient *http.Client
	baseURL    string
	accessKey  string
	secretKey  string
}

// New new instance
func New(client *http.Client, baseURL, apiKey, secretKey string) *BikiSpot {
	instance := new(BikiSpot)
	if baseURL == "" {
		instance.baseURL = "https://openapi.biki.com"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	return instance
}

// NewWithConfig new instance by config
func NewWithConfig(config *APIConfig) *BikiSpot {
	instance := new(BikiSpot)
	if config.Endpoint == "" {
		instance.baseURL = "https://openapi.biki.com"
	} else {
		instance.baseURL = config.Endpoint
	}

	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	return instance
}

// GetExchangeName return exchange name
func (spot *BikiSpot) GetExchangeName() string {
	return EXCHANGE_BIKI
}

// GetCoinList exchange supported coins
func (spot *BikiSpot) GetCoinList() interface{} {
	return ReturnAPIError(MethodNotExistError)
}

// GetSymbolList exchange all symbol
func (spot *BikiSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	return spot.httpGet("/open/api/common/symbols", params, false)
}

// GetDepth symbol depth
func (spot *BikiSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	if step, ok := options["type"]; ok {
		params.Set("type", step)
	} else {
		params.Set("type", "step0")
	}
	return spot.httpGet("/open/api/market_dept", params, false)
}

// GetTicker symbol ticker
func (spot *BikiSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	return spot.httpGet("/open/api/get_ticker", params, false)
}

// GetKline symbol kline
func (spot *BikiSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1"
	}
	params.Set("period", periodStr)
	return spot.httpGet("/open/api/get_records", params, false)
}

// GetTrade symbol last trade
func (spot *BikiSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	return spot.httpGet("/open/api/get_trades", params, false)
}

// GetUserBalance user balance
func (spot *BikiSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpGet("/open/api/user/account", params, true)
	return result
}

// PlaceOrder place order
func (spot *BikiSpot) PlaceOrder(order *PlaceOrder) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(order.Symbol))
	params.Set("volume", order.Amount)
	if order.Side == BUY {
		params.Set("side", BIKI_BUY)
	} else {
		params.Set("side", BIKI_SELL)
	}
	if order.TradeType == LIMIT {
		params.Set("type", "1")
		params.Set("price", order.Price)
	} else {
		params.Set("type", "2")
	}
	result := spot.httpPost("/open/api/create_order", params, false)
	return result
}

// PlaceLimitOrder place limit order
func (spot *BikiSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("price", price)
	params.Set("volume", amount)
	params.Set("type", "1")
	if side == BUY {
		params.Set("side", BIKI_BUY)
	} else {
		params.Set("side", BIKI_SELL)
	}
	result := spot.httpPost("/open/api/create_order", params, true)
	return result
}

// PlaceMarketOrder place market order
func (spot *BikiSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("volume", amount)
	params.Set("type", "2")
	if side == BUY {
		params.Set("side", BIKI_BUY)
	} else {
		params.Set("side", BIKI_SELL)
	}
	result := spot.httpPost("/open/api/create_order", params, true)
	return result
}

// BatchPlaceLimitOrder batch place limit order
func (spot *BikiSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	params := &url.Values{}
	var trustOrders []map[string]interface{}
	var symbol Symbol
	for _, item := range orders {
		param := map[string]interface{}{}
		param["side"] = BIKI_BUY
		if item.Side == SELL {
			param["side"] = BIKI_SELL
		}
		param["volume"] = item.Amount
		param["price"] = item.Price
		param["type"] = "1"

		symbol = item.Symbol
		trustOrders = append(trustOrders, param)
	}
	jsonBody, _ := json.Marshal(trustOrders)
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("mass_place", string(jsonBody))

	result := spot.httpPost("/open/api/mass_replaceV2", params, true)
	return result
}

// CancelOrder cancel a order
func (spot *BikiSpot) CancelOrder(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("order_id", orderID)
	result := spot.httpPost("/open/api/cancel_order", params, true)
	return result
}

// BatchCancelOrder batch cancel orders
func (spot *BikiSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("mass_cancel", fmt.Sprintf("[%s]", orderIds))
	result := spot.httpPost("/open/api/mass_replaceV2", params, true)
	return result
}

// BatchCancelAllOrder batch cancel all orders
func (spot *BikiSpot) BatchCancelAllOrder(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	result := spot.httpPost("/open/api/cancel_order_all", params, true)
	return result
}

// GetUserOpenTrustOrders get current trust order
func (spot *BikiSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("pageSize", strconv.FormatInt(int64(size), 10))
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	result := spot.httpGet("/open/api/v2/new_order", params, true)
	return result
}

// GetUserOrderInfo get trust order info
func (spot *BikiSpot) GetUserOrderInfo(symbol Symbol, orderID, clientOrderID string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("order_id", orderID)
	result := spot.httpGet("/open/api/order_info", params, true)
	return result
}

// GetUserTradeOrders get trade order list
func (spot *BikiSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("pageSize", strconv.FormatInt(int64(size), 10))

	if start, ok := options["startDate"]; ok == true {
		params.Set("startDate", start)
	}
	if end, ok := options["endDate"]; ok == true {
		params.Set("endDate", end)
	}
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	if sort, ok := options["sort"]; ok == true {
		params.Set("sort", sort)
	}

	result := spot.httpGet("/open/api/all_trade", params, true)
	return result
}

// GetUserTrustOrders get trust order list
func (spot *BikiSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", spot.getSymbol(symbol))
	params.Set("pageSize", strconv.FormatInt(int64(size), 10))

	if start, ok := options["startDate"]; ok == true {
		params.Set("startDate", start)
	}
	if end, ok := options["endDate"]; ok == true {
		params.Set("endDate", end)
	}
	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}

	result := spot.httpGet("/open/api/v2/all_order", params, true)
	return result
}

// HttpRequest request api
func (spot *BikiSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
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

func (spot *BikiSpot) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse

	if signed {
		spot.sign(params)
	}
	requestURL := spot.baseURL + url + "?" + params.Encode()
	responseMap = HttpGet(spot.httpClient, requestURL)
	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *BikiSpot) httpPost(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap HttpClientResponse
	requestURL := spot.baseURL + url
	if signed {
		spot.sign(params)
	}

	responseMap = HttpPost(spot.httpClient, requestURL, params.Encode())

	fmt.Println(requestURL)
	return spot.handlerResponse(&responseMap)
}

func (spot *BikiSpot) handlerResponse(responseMap *HttpClientResponse) map[string]interface{} {
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

func (spot *BikiSpot) sign(params *url.Values) {
	params.Set("api_key", spot.accessKey)
	params.Set("time", GetNowTimestampStr())

	var buf strings.Builder
	keys := make([]string, 0, len(*params))
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := (*params)[k]
		buf.WriteString(k)
		buf.WriteString(vs[0])
	}
	signStr := buf.String() + spot.secretKey
	sign := Md5Signer(signStr)
	params.Set("sign", sign)
}

func (spot BikiSpot) getSymbol(symbol Symbol) string {
	return symbol.ToSymbol("")
}
