package huobi

import "github.com/primitivelab/goexchange"

// Swap swap api interface
type Swap interface {

	// exchange name
	GetExchangeName() string
	// Get exchange contract market list
	GetContractList() interface{}
	// Get exchange contract depth
	GetDepth(symbol goexchange.Symbol, size int, options map[string]string) map[string]interface{}
	// Get exchange contract ticker
	GetTicker(symbol goexchange.Symbol) interface{}
	// Get exchange contract kline
	GetKline(symbol goexchange.Symbol, period int, size int, options map[string]string) interface{}
	// Get exchange contract trade
	GetTrade(symbol goexchange.Symbol, size int, options map[string]string) interface{}
	// GetPremiumIndex exchange index price& market price & funding rate
	GetPremiumIndex(symbol goexchange.Symbol) interface{}
	// Get exchange http request
	HTTPRequest(requestURL, method string, options interface{}, signed bool) interface{}
}
