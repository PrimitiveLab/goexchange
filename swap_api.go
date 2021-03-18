package goexchange

// SwapAPI swap api interface
type SwapAPI interface {

	// exchange name
	GetExchangeName() string
	// Get exchange contract market list
	GetContractList() interface{}
	// Get exchange contract depth
	GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{}
	// Get exchange contract ticker
	GetTicker(symbol Symbol) interface{}
	// Get exchange contract kline
	GetKline(symbol Symbol, period int, size int, options map[string]string) interface{}
	// Get exchange contract trade
	GetTrade(symbol Symbol, size int, options map[string]string) interface{}
	// Get exchange contract trade
	GetPremiumIndex(symbol Symbol) interface{}
	// Get exchange http request
	HTTPRequest(requestURL, method string, options interface{}, signed bool) interface{}
}
