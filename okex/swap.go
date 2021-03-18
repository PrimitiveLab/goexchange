package okex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	goex "github.com/primitivelab/goexchange"
)

// Swap okex contract
type Swap struct {
	httpClient *http.Client
	baseURL    string
	accessKey  string
	secretKey  string
	passphrase string
}

// NewSwap new instance
func NewSwap(client *http.Client, baseURL, apiKey, secretKey, passphrase string) *Swap {
	instance := new(Swap)
	if baseURL == "" {
		instance.baseURL = "https://www.okex.com"
	} else {
		instance.baseURL = baseURL
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	instance.passphrase = passphrase
	return instance
}

// NewSwapWithConfig new instance with config struct
func NewSwapWithConfig(config *goex.APIConfig) *Swap {
	instance := new(Swap)
	if config.Endpoint == "" {
		instance.baseURL = "https://www.okex.com"
	} else {
		instance.baseURL = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	instance.passphrase = config.ApiPassphrase
	return instance
}

// GetExchangeName get exchange name
func (swap *Swap) GetExchangeName() string {
	return goex.EXCHANGE_OKEX
}

// GetContractList exchange contract list
func (swap *Swap) GetContractList() interface{} {
	params := &url.Values{}
	return swap.httpGet("/api/swap/v3/instruments", params, false)
}

// GetDepth exchange depth data
func (swap *Swap) GetDepth(symbol goex.Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	instrumentId := swap.getSymbol(symbol)
	params.Set("instrument_id", instrumentId)
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}

	if depth, ok := options["depth"]; ok {
		params.Set("depth", depth)
	}
	result := swap.httpGet(fmt.Sprintf("/api/swap/v3/instruments/%s/depth", instrumentId), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetTicker exchange ticker data
func (swap *Swap) GetTicker(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	instrumentId := swap.getSymbol(symbol)
	params.Set("instrument_id", instrumentId)
	result := swap.httpGet(fmt.Sprintf("/api/swap/v3/instruments/%s/ticker", instrumentId), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetKline exchange kline data
func (swap *Swap) GetKline(symbol goex.Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	instrumentId := swap.getSymbol(symbol)
	params.Set("instrument_id", instrumentId)
	periodStr, ok := klinePeriod[period]
	if !ok {
		periodStr = "60"
	}
	params.Set("granularity", periodStr)

	if start, ok := options["start"]; ok {
		params.Set("start", start)
	}
	if end, ok := options["end"]; ok {
		params.Set("end", end)
	}
	result := swap.httpGet(fmt.Sprintf("/api/swap/v3/instruments/%s/candles", instrumentId), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetTrade exchange trade order data
func (swap *Swap) GetTrade(symbol goex.Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	instrumentId := swap.getSymbol(symbol)
	params.Set("instrument_id", instrumentId)
	if size != 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if after, ok := options["after"]; ok {
		params.Set("after", after)
	}
	if before, ok := options["before"]; ok {
		params.Set("before", before)
	}
	result := swap.httpGet(fmt.Sprintf("/api/swap/v3/instruments/%s/trades", instrumentId), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// GetPremiumIndex exchange index price& market price & funding rate
func (swap *Swap) GetPremiumIndex(symbol goex.Symbol) interface{} {
	params := &url.Values{}
	instrumentId := swap.getSymbol(symbol)
	params.Set("instrument_id", instrumentId)
	result := swap.httpGet(fmt.Sprintf("/api/swap/v3/instruments/%s/index", instrumentId), params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

// HTTPRequest request url
func (swap *Swap) HTTPRequest(requestURL, method string, options interface{}, signed bool) interface{} {
	method = strings.ToUpper(method)
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}
	switch method {
	case goex.HTTP_GET:
		return swap.httpGet(requestURL, params, signed)
	case goex.HTTP_POST:
		return swap.httpPost(requestURL, params, signed)
	}
	return nil
}

// httpGet Get request method
func (swap *Swap) httpGet(url string, params *url.Values, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	var headers map[string]string

	requestUrl := swap.baseURL + url
	reqData := ""
	if params != nil {
		reqData = "?" + params.Encode()
		requestUrl = requestUrl + reqData
	}

	if signed {
		timestamp := goex.IsoTime()
		sign := swap.sign(url, goex.HTTP_GET, timestamp, reqData)
		headers = map[string]string{
			"OK-ACCESS-KEY":        swap.accessKey,
			"OK-ACCESS-SIGN":       sign,
			"OK-ACCESS-PASSPHRASE": swap.passphrase,
			"OK-ACCESS-TIMESTAMP":  timestamp,
		}
	}
	responseMap = goex.HttpGetWithHeader(swap.httpClient, requestUrl, headers)
	return swap.handlerResponse(&responseMap)
}

// httpPost Post request method
func (swap *Swap) httpPost(url string, params interface{}, signed bool) map[string]interface{} {
	var responseMap goex.HttpClientResponse
	var headers map[string]string
	requestUrl := swap.baseURL + url
	reqData := ""
	if params != nil {
		jsonBody, _ := json.Marshal(params)
		reqData = string(jsonBody)
	}

	timestamp := goex.IsoTime()
	sign := swap.sign(url, goex.HTTP_POST, timestamp, reqData)
	headers = map[string]string{
		"OK-ACCESS-KEY":        swap.accessKey,
		"OK-ACCESS-SIGN":       sign,
		"OK-ACCESS-PASSPHRASE": swap.passphrase,
		"OK-ACCESS-TIMESTAMP":  timestamp,
	}

	responseMap = goex.HttpPostWithJson(swap.httpClient, requestUrl, reqData, headers)
	return swap.handlerResponse(&responseMap)
}

// sign signature method
func (swap *Swap) sign(url, method, timestamp, reqData string) string {
	signStr := timestamp + method + url + reqData
	sign, _ := goex.HmacSha256Base64Signer(signStr, swap.secretKey)
	return sign
}

// handlerResponse Handler response data format
func (swap *Swap) handlerResponse(responseMap *goex.HttpClientResponse) map[string]interface{} {
	returnData := make(map[string]interface{})

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
		returnData["code"] = goex.JsonUnmarshalError.Code
		returnData["msg"] = goex.JsonUnmarshalError.Msg
		returnData["error"] = err.Error()
		return returnData
	}

	returnData["data"] = bodyDataMap
	return returnData
}

// getSymbol format symbol method
func (swap Swap) getSymbol(symbol goex.Symbol) string {
	return symbol.ToUpper().ToSymbol("-") + "-SWAP"
}
