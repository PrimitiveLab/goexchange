package huobi

import (
	// "errors"
	// "fmt"
	"net/http"
	// "strings"
	// "time"

	."github.com/mo-zhe/goexchange"
)

var _INTERNAL_KLINE_PERIOD_CONVERTER = map[int]string{
	KLINE_PERIOD_1MIN:   "1min",
	KLINE_PERIOD_5MIN:   "5min",
	KLINE_PERIOD_15MIN:  "15min",
	KLINE_PERIOD_30MIN:  "30min",
	KLINE_PERIOD_60MIN:  "60min",
	KLINE_PERIOD_1DAY:   "1day",
	KLINE_PERIOD_1WEEK:  "1week",
	KLINE_PERIOD_1MONTH: "1mon",
	KLINE_PERIOD_1YEAR:  "1year",
}

type HuobiSpot struct {
	httpClient *http.Client
	baseUrl string
	accountId string
	accessKey string
	secretKey string
}

func NewHuobi(client *http.Client, apiKey, secretKey, accountId string) *HuobiSpot {
	hb := new(HuobiSpot)
	hb.baseUrl = "https://api.huobi.pro"
	hb.httpClient = client
	hb.accessKey = apiKey
	hb.secretKey = secretKey
	hb.accountId = accountId
	return hb
}

// func (hb *HuobiSpot) GetDepth(size int, symbol Symbol) (*Depth, error) {
func (hb *HuobiSpot) GetDepth(size int, symbol Symbol) (error) {
	// url := hb.baseUrl + "/market/depth?symbol=%s&type=step0&depth=%d"
	// n := 5
	// pair := currency.AdaptUsdToUsdt()
	// if size <= 5 {
	// 	n = 5
	// } else if size <= 10 {
	// 	n = 10
	// } else if size <= 20 {
	// 	n = 20
	// } else {
	// 	url = hb.baseUrl + "/market/depth?symbol=%s&type=step0&d=%d"
	// }
	// respmap, err := HttpGet(hb.httpClient, fmt.Sprintf(url, strings.ToLower(pair.ToSymbol("")), n))
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if "ok" != respmap["status"].(string) {
	// 	return nil, errors.New(respmap["err-msg"].(string))
	// }
	//
	// tick, _ := respmap["tick"].(map[string]interface{})
	//
	// dep := hb.parseDepthData(tick, size)
	// dep.Pair = currency
	// mills := ToUint64(tick["ts"])
	// dep.UTime = time.Unix(int64(mills/1000), int64(mills%1000)*int64(time.Millisecond))

	return nil
}