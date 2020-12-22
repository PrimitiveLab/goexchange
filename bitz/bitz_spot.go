package bitz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	. "github.com/primitivelab/goexchange"
)

const (
	BITZ_BUY  string = "1"
	BITZ_SELL string = "2"
)

var klinePeriod = map[int]string{
	KLINE_PERIOD_1MINUTE:  "1min",
	KLINE_PERIOD_5MINUTE:  "5min",
	KLINE_PERIOD_15MINUTE: "15min",
	KLINE_PERIOD_30MINUTE: "30min",
	KLINE_PERIOD_60MINUTE: "60min",
	KLINE_PERIOD_1HOUR:    "60min",
	KLINE_PERIOD_4HOUR:    "4hour",
	KLINE_PERIOD_1DAY:     "1day",
	KLINE_PERIOD_3DAY:     "3day",
	KLINE_PERIOD_5DAY:     "5day",
	KLINE_PERIOD_1WEEK:    "1week",
	KLINE_PERIOD_1MONTH:   "1mon",
}

type BitzSpot struct {
	httpClient *http.Client
	baseUrl    string
	accessKey  string
	secretKey  string
	passphrase string
}

func New(client *http.Client, baseUrl, apiKey, secretKey, passphrase string) *BitzSpot {
	instance := new(BitzSpot)
	if baseUrl == "" {
		instance.baseUrl = "https://apiv2.bitz.com"
	} else {
		instance.baseUrl = baseUrl
	}
	instance.httpClient = client
	instance.accessKey = apiKey
	instance.secretKey = secretKey
	instance.passphrase = passphrase
	return instance
}

func NewWithConfig(config *APIConfig) *BitzSpot {
	instance := new(BitzSpot)
	if config.Endpoint == "" {
		instance.baseUrl = "https://apiv2.bitz.com"
	} else {
		instance.baseUrl = config.Endpoint
	}
	instance.httpClient = config.HttpClient
	instance.accessKey = config.ApiKey
	instance.secretKey = config.ApiSecretKey
	instance.passphrase = config.ApiPassphrase
	return instance
}

func (spot *BitzSpot) GetExchangeName() string {
	return EXCHANGE_BITZ
}

func (spot *BitzSpot) GetCoinList() interface{} {
	params := &url.Values{}
	result := spot.httpRequest("/api2/1/coininfo", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

func (spot *BitzSpot) GetSymbolList() interface{} {
	params := &url.Values{}
	result := spot.httpRequest("/Market/symbolList", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *BitzSpot) GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol("_"))
	result := spot.httpRequest("/Market/depth", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *BitzSpot) GetTicker(symbol Symbol) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol("_"))
	result := spot.httpRequest("/Market/ticker", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *BitzSpot) GetKline(symbol Symbol, period, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol("_"))
	periodStr, ok := klinePeriod[period]
	if ok != true {
		periodStr = "1min"
	}
	params.Set("resolution", periodStr)
	if size != 0 {
		params.Set("size", strconv.Itoa(size))
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("to", endTime)
	}

	result := spot.httpRequest("/Market/kline", "get", params, false)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *BitzSpot) GetTrade(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.ToSymbol("_"))
	result := spot.httpRequest("/Market/order", "get", params, false)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 获取余额
func (spot *BitzSpot) GetUserBalance() interface{} {
	params := &url.Values{}
	result := spot.httpRequest("/Assets/getUserAssets", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}

	return result
}

// 批量下单
func (spot *BitzSpot) PlaceOrder(order *PlaceOrder) interface{} {
	return nil
}

// 下限价单
func (spot *BitzSpot) PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("price", price)
	params.Set("number", amount)
	params.Set("type", BITZ_BUY)
	if side == SELL {
		params.Set("type", BITZ_SELL)
	}
	params.Set("tradePwd", spot.passphrase)
	result := spot.httpRequest("/Trade/addEntrustSheet", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 下市价单
func (spot *BitzSpot) PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("total", amount)
	params.Set("type", BITZ_BUY)
	if side == SELL {
		params.Set("type", BITZ_SELL)
	}
	params.Set("tradePwd", spot.passphrase)
	result := spot.httpRequest("/Trade/MarketTrade", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 批量下限价单
func (spot *BitzSpot) BatchPlaceLimitOrder(orders []LimitOrder) interface{} {
	var trustOrders []map[string]interface{}
	tradePwd := Md5Signer(spot.passphrase)
	for _, item := range orders {
		param := map[string]interface{}{}
		param["coins"] = item.Symbol.String()
		param["price"] = item.Price
		param["number"] = item.Amount
		param["type"] = BITZ_BUY
		if item.Side == SELL {
			param["type"] = BITZ_SELL
		}
		param["tradepwd"] = tradePwd
		trustOrders = append(trustOrders, param)
	}
	jsonBody, _ := json.Marshal(trustOrders)

	params := &url.Values{}
	params.Set("tradeData", string(jsonBody))
	return spot.httpRequest("/Trade/addEntrustSheetBatch", HTTP_POST, params, true)
}

// 撤单
func (spot *BitzSpot) CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("entrustSheetId", orderId)
	result := spot.httpRequest("/Trade/cancelEntrustSheet", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 批量撤单
func (spot *BitzSpot) BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{} {
	params := &url.Values{}
	params.Set("ids", orderIds)
	result := spot.httpRequest("/Trade/cancelAllEntrustSheet", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 我的当前委托单
func (spot *BitzSpot) GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("coinFrom", symbol.CoinFrom)
	params.Set("coinTo", symbol.CoinTo)
	if size != 0 {
		params.Set("pageSize", strconv.Itoa(size))
	}

	if side, ok := options["type"]; ok == true {
		params.Set("type", side)
	}

	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}

	result := spot.httpRequest("/Trade/getUserNowEntrustSheet", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 委托单详情
func (spot *BitzSpot) GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{} {
	params := &url.Values{}
	params.Set("entrustSheetId", orderId)
	result := spot.httpRequest("/Trade/getEntrustSheetInfo", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

// 我的成交单列表
func (spot *BitzSpot) GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{} {
	return nil
}

// 我的委托单列表
func (spot *BitzSpot) GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{} {
	params := &url.Values{}
	params.Set("coinFrom", symbol.CoinFrom)
	params.Set("coinTo", symbol.CoinTo)
	if size != 0 {
		params.Set("pageSize", strconv.Itoa(size))
	}

	if side, ok := options["type"]; ok == true {
		params.Set("type", side)
	}

	if page, ok := options["page"]; ok == true {
		params.Set("page", page)
	}
	if startTime, ok := options["startTime"]; ok == true {
		params.Set("startTime", startTime)
	}
	if endTime, ok := options["endTime"]; ok == true {
		params.Set("endTime", endTime)
	}

	result := spot.httpRequest("/Trade/getUserHistoryEntrustSheet", HTTP_POST, params, true)
	if result["code"] != 0 {
		return result
	}
	return result
}

func (spot *BitzSpot) HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{} {
	params := &url.Values{}
	mapOptions := options.(map[string]string)
	for key, val := range mapOptions {
		params.Set(key, val)
	}
	return spot.httpRequest(requestUrl, strings.ToUpper(method), params, signed)
}

func (spot *BitzSpot) httpRequest(url, method string, params *url.Values, signed bool) map[string]interface{} {
	method = strings.ToUpper(method)

	var responseMap HttpClientResponse
	requestUrl := spot.baseUrl + url
	switch method {
	case HTTP_GET:
		if params != nil {
			requestUrl = requestUrl + "?" + params.Encode()
		}
		responseMap = HttpGet(spot.httpClient, requestUrl)
	case HTTP_POST:
		params.Set("sign", spot.sign(*params))
		responseMap = HttpPost(spot.httpClient, requestUrl, params.Encode())
	}

	fmt.Println(requestUrl)

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
	resStatus := bodyDataMap["status"].(float64)
	if 200 != resStatus {
		returnData["code"] = ExchangeError.Code
		returnData["msg"] = ExchangeError.Msg
		returnData["error"] = fmt.Sprintf("%g: %s", resStatus, spot.getError(resStatus))
		return returnData
	}
	returnData["data"] = bodyDataMap["data"]
	return returnData
}

func (spot *BitzSpot) sign(params url.Values) string {

	timestamp := GetNowTimestampStr()
	params.Set("apiKey", spot.accessKey)
	params.Set("timeStamp", timestamp)
	params.Set("nonce", timestamp[3:])

	var buf strings.Builder
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := params[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(vs[0])
	}

	signStr := buf.String() + spot.secretKey
	return Md5Signer(signStr)
}

func (spot *BitzSpot) getError(code float64) string {
	errorMap := map[float64]string{
		200:      "成功",
		-102:     "参数错误",
		-103:     "校验失败",
		-104:     "网络异常-1",
		-105:     "签名不匹配",
		-106:     "网络异常-2",
		-107:     "请求路径错误",
		-109:     "scretKey错误",
		-110:     "访问请求次数超限",
		-111:     "当前IP不在可信任IP范围内",
		-112:     "服务正在维护中",
		-114:     "每日请求次数已达上限",
		-117:     "apikey有效期已到期",
		-100015:  "交易密码错误",
		-100044:  "请求数据失败",
		-100101:  "交易对信息不存在",
		-100201:  "交易对信息不存在",
		-100301:  "交易对信息不存在",
		-100401:  "交易对信息不存在",
		-100302:  "k线type值错误",
		-100303:  "k线size值超出范围",
		-200000:  "委托单已撤销",
		-200003:  "请先完善交易密码",
		-200005:  "该账号不可交易",
		-200025:  "临时停牌",
		-200027:  "价格错误",
		-200028:  "数量需大于0",
		-200029:  "交易数量需在%s到%d之间",
		-200030:  "超过价格范围",
		-200031:  "资产不足",
		-200032:  "请联系客服",
		-200033:  "下单失败",
		-200034:  "委托单已成交或已取消",
		-200035:  "撤销失败，该委托单已成交",
		-200036:  "撤销失败",
		-200037:  "交易方向错误",
		-200038:  "交易对错误",
		-200055:  "委托记录不存在",
		-200056:  "买入金额需在%s到%d之间",
		-300069:  "apiKey不合法",
		-300101:  "交易类型错误",
		-300102:  "下单金额和数量不能小于0",
		-300103:  "交易密码错误",
		-301001:  "网络异常-3",
		-200053:  "暂停充币",
		-300041:  "暂停提币",
		-200001:  "币种不存在",
		-200054:  "暂时无法转出",
		-300076:  "您的账号被限制充币",
		-200004:  "钱包创建地址错误",
		-300077:  "您的账号被限制提币",
		-300046:  "请您先绑定邮箱",
		-300007:  "请开启手机验证或谷歌验证",
		-300040:  "单笔最小提币数量为%s，最大提币数量为%d",
		-300042:  "修改安全设置后，需24小时后才能提现",
		-300043:  "请先添加提币地址",
		-100027:  "资产不足",
		-300044:  "网络费设置错误",
		-300048:  "不支持站内互转",
		-300047:  "请输入合法的转出币地址",
		-300049:  "memo最大的长度20位",
		-300050:  "memo输入有误",
		-300084:  "为使您的提币准确到账，向某些交易所地址提币必须输入正确格式的MEMO。若非此原因，请检查提币地址的格式是否错误",
		-300045:  "您的账户单日提现限额为%sBTC,您还可以提现%dBTC",
		-300075:  "需要审核且用户未绑定手机号",
		-300091:  "非常抱歉，您的API提币额度超过限额，请您登录Bitz官方APP或官方网站，进行提币操作",
		-100031:  "添加用户账号资产错误",
		-100028:  "扣除资产错误",
		-2001001: "内部错误",
		-2001003: "参数错误",
		-2001004: "签名错误",
		-2001005: "合约不存在",
		-2001006: "该市场暂停交易",
		-2001007: "uid错误",
		-2001008: "未开通合约交易",
		-2001009: "价格错误",
		-2001010: "合约账户被锁定",
		-300092:  "API提币地址没有加入白名单,且没有历史提币记录",
		-300037:  "不在提币地址列表",
	}

	msg, ok := errorMap[code]
	if ok {
		return msg
	} else {
		return "未知错误"
	}

}
