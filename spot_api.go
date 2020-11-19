package goexchange

// api interface

type SpotAPI interface {
	// LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error)
	// LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error)
	// MarketBuy(amount, price string, currency CurrencyPair) (*Order, error)
	// MarketSell(amount, price string, currency CurrencyPair) (*Order, error)
	// CancelOrder(orderId string, currency CurrencyPair) (bool, error)
	// GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
	// GetUnfinishOrders(currency CurrencyPair) ([]Order, error)
	// GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error)
	// GetAccount() (*Account, error)

	GetExchangeName() string
	GetCoinList() interface{}
	GetSymbolList() interface{}
	GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{}
	GetTicker(symbol Symbol) interface{}
	GetKline(symbol Symbol, period int, size int, options map[string]string) interface{}
	GetTrade(symbol Symbol, size int, options map[string]string) interface{}
	HttpRequest(requestUrl, method string, options map[string]string, signed bool) interface{}
}
