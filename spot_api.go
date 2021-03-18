package goexchange

// SpotAPI spot api interface
type SpotAPI interface {

	// 获取余额
	GetUserBalance() interface{}

	// 批量下单
	PlaceOrder(order *PlaceOrder) interface{}

	// 下限价单
	PlaceLimitOrder(symbol Symbol, price string, amount string, side TradeSide, ClientOrderId string) interface{}

	// 下市价单
	PlaceMarketOrder(symbol Symbol, amount string, side TradeSide, ClientOrderId string) interface{}

	// 批量下限价单
	BatchPlaceLimitOrder(orders []LimitOrder) interface{}

	// 撤单
	CancelOrder(symbol Symbol, orderId, clientOrderId string) interface{}

	// 批量撤单
	BatchCancelOrder(symbol Symbol, orderIds, clientOrderIds string) interface{}

	// 我的当前委托单
	GetUserOpenTrustOrders(symbol Symbol, size int, options map[string]string) interface{}

	// 委托单详情
	GetUserOrderInfo(symbol Symbol, orderId, clientOrderId string) interface{}

	// 我的成交单列表
	GetUserTradeOrders(symbol Symbol, size int, options map[string]string) interface{}

	// 我的委托单列表
	GetUserTrustOrders(symbol Symbol, status string, size int, options map[string]string) interface{}

	GetExchangeName() string
	GetCoinList() interface{}
	GetSymbolList() interface{}
	GetDepth(symbol Symbol, size int, options map[string]string) map[string]interface{}
	GetTicker(symbol Symbol) interface{}
	GetKline(symbol Symbol, period int, size int, options map[string]string) interface{}
	GetTrade(symbol Symbol, size int, options map[string]string) interface{}
	HttpRequest(requestUrl, method string, options interface{}, signed bool) interface{}
}
