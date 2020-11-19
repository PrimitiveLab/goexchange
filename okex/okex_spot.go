package okex

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
}

func New(client *http.Client, apiKey, secretKey string) *OkexSpot {
	instance := new(OkexSpot)
	instance.baseUrl = "https://www.okex.com"
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
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
	return instance
}

func (okexSpot *OkexSpot) GetExchangeName() string {
	return ECHANGE_OKEX
}

func (okexSpot *OkexSpot) GetCoinList() interface{} {
	params := &url.Values{}
	result := okexSpot.httpRequest("/api/account/v3/currencies", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (okexSpot *OkexSpot) GetSymbolList() interface{} {

	params := &url.Values{}
	result := okexSpot.httpRequest("/api/spot/v3/instruments", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (okexSpot *OkexSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	params.Set("size", strconv.Itoa(size))
	if depthType, ok := options["depth"]; ok == true {
		params.Set("depth", depthType)
	}

	result := okexSpot.httpRequest("/api/spot/v3/instruments/" + instrumentId + "/book", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (okexSpot *OkexSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	result := okexSpot.httpRequest("/api/spot/v3/instruments/" + instrumentId + "/ticker", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (okexSpot *OkexSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "60"
	}
	params.Set("granularity", periodStr)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("start", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("end", endTime)
	}

	result := okexSpot.httpRequest("/api/spot/v3/instruments/" + instrumentId + "/history/candles", "get", params,false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (okexSpot *OkexSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {

	params := &url.Values{}
	instrumentId := symbol.ToUpper().ToSymbol("-")
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	result := okexSpot.httpRequest("/api/spot/v3/instruments/" + instrumentId + "/trades", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (okexSpot *OkexSpot) HttpRequest(requestUrl, method string, options map[string]string, signed bool) interface{} {

	params := &url.Values{}
	for key, value := range options {
		params.Set(key, value)
	}

	return okexSpot.httpRequest(requestUrl, method, params, signed)
}

func (okexSpot *OkexSpot) httpRequest(url , method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	switch method {
	case "GET":
		requestUrl := okexSpot.baseUrl + url + "?" + params.Encode()
		fmt.Println(requestUrl)
		responseMap = HttpGet(okexSpot.httpClient, requestUrl)
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
