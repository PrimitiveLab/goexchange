package huobi

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
	KLINE_PERIOD_1MINUTE:  "1min",
	KLINE_PERIOD_5MINUTE:  "5min",
	KLINE_PERIOD_15MINUTE: "15min",
	KLINE_PERIOD_30MINUTE: "30min",
	KLINE_PERIOD_60MINUTE: "60min",
	KLINE_PERIOD_1DAY:     "1day",
	KLINE_PERIOD_1WEEK:    "1week",
	KLINE_PERIOD_1MONTH:   "1mon",
	KLINE_PERIOD_1YEAR:    "1year",
}

const (
	HUOBI_SPOT_ACCOUNT = "spot"
)

type HuobiSpot struct {
	httpClient *http.Client
	baseUrl    string
	accountId  string
	accessKey  string
	secretKey  string
}

func New(client *http.Client, apiKey, secretKey, accountId string) *HuobiSpot {
	hb := new(HuobiSpot)
	hb.baseUrl = "https://api.huobi.pro"
	hb.httpClient = client
	hb.accessKey = apiKey
	hb.secretKey = secretKey
	hb.accountId = accountId
	return hb
}

func NewWithConfig(config *APIConfig) *HuobiSpot {
	hb := new(HuobiSpot)
	if config.Endpoint == "" {
		hb.baseUrl = "https://api.huobi.pro"
	} else {
		hb.baseUrl = config.Endpoint
	}
	hb.httpClient = config.HttpClient
	hb.accessKey = config.ApiKey
	hb.secretKey = config.ApiSecretKey
	hb.accountId = config.AccountId
	return hb
}

func (hb *HuobiSpot) GetExchangeName() string {
	return ECHANGE_HUOBI
}

func (hb *HuobiSpot) GetCoinList() interface{} {

	params := &url.Values{}
	result := hb.httpRequest("/v2/reference/currencies", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

func (hb *HuobiSpot) GetSymbolList() interface{} {

	params := &url.Values{}
	result := hb.httpRequest("/v1/common/symbols", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

func (hb *HuobiSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {

	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol(""))

	// depth ranges [5,10,20]
	if size < 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else {
		size = 20
	}
	params.Set("depth", strconv.Itoa(size))

	if depthType, ok := options["type"]; ok == true {
		params.Set("type", depthType)
	}

	result := hb.httpRequest("/market/depth", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["tick"]
	return result
}

func (hb *HuobiSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol(""))
	result := hb.httpRequest("/market/detail/merged", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["tick"]
	return result
}

func (hb *HuobiSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol(""))
	periodStr, isOk := klinePeriod[period]
	if isOk != true {
		periodStr = "1min"
	}
	params.Set("period", periodStr)
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	result := hb.httpRequest("/market/history/kline", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

func (hb *HuobiSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {

	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol(""))
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	result := hb.httpRequest("/market/history/trade", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	result["data"] = result["data"].(map[string]interface{})["data"]
	return result
}

func (hb *HuobiSpot) HttpRequest(requestUrl, method string, options map[string]string, signed bool) interface{} {

	params := &url.Values{}
	for key, value := range options {
		params.Set(key, value)
	}

	return hb.httpRequest(requestUrl, method, params, signed)
}

func (hb *HuobiSpot) httpRequest(url, method string, params *url.Values, signed bool) map[string]interface{} {

	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := hb.baseUrl + url + "?" + params.Encode()
		fmt.Println(requestUrl)
		responseMap = HttpGet(hb.httpClient, requestUrl)
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

	var bodyDataMap map[string]interface{}
	err := json.Unmarshal(responseMap.Data, &bodyDataMap)
	if err != nil {
		log.Println(string(responseMap.Data))
		returnData["code"] = JsonUnmarshalError.Code
		returnData["msg"] = JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}

	// fmt.Println(returnData.Data)
	if "ok" != bodyDataMap["status"].(string) {
		returnData["code"] = ExchangeError.Code
		returnData["msg"] = ExchangeError.Msg
		returnData["error"] = bodyDataMap["err-msg"].(string)
		return returnData
	}
	returnData["data"] = bodyDataMap
	return returnData
}
