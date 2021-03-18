package goexchange

// 交易方向
type TradeSide int

const (
	BUY  TradeSide = 1
	SELL TradeSide = -1
)

func (ts TradeSide) String() string {
	switch ts {
	case BUY:
		return "buy"
	case SELL:
		return "sell"
	default:
		return "unknown"
	}
}

// Time in force 策略
type TimeInForce int

const (
	// 成交为止
	GTC TimeInForce = 0
	// 只吃单不挂单(只做maker单).
	POC TimeInForce = 1
	// 立即成交并取消剩余,只吃单不挂单(只做taker单).
	IOC TimeInForce = 2
	// 全部成交或立即取消,如果无法全部成交，订单会失效
	FOK TimeInForce = 3
	// 全部成交或立即取消,如果无法全部成交，订单会失效
	GTX TimeInForce = 4
)

// 交易类型
const (
	LIMIT  string = "limit"
	MARKET string = "market"
)

// k线周期
const (
	KLINE_PERIOD_1MINUTE = iota + 1
	KLINE_PERIOD_3MINUTE
	KLINE_PERIOD_5MINUTE
	KLINE_PERIOD_15MINUTE
	KLINE_PERIOD_30MINUTE
	KLINE_PERIOD_60MINUTE
	KLINE_PERIOD_1HOUR
	KLINE_PERIOD_2HOUR
	KLINE_PERIOD_3HOUR
	KLINE_PERIOD_4HOUR
	KLINE_PERIOD_6HOUR
	KLINE_PERIOD_8HOUR
	KLINE_PERIOD_12HOUR
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_3DAY
	KLINE_PERIOD_5DAY
	KLINE_PERIOD_7DAY
	KLINE_PERIOD_1WEEK
	KLINE_PERIOD_1MONTH
	KLINE_PERIOD_1YEAR
)

// exchange name const
const (
	EXCHANGE_OKEX     = "okex"
	EXCHANGE_HUOBI    = "huobi"
	EXCHANGE_BINANCE  = "binance"
	EXCHANGE_GATE     = "gate"
	EXCHANGE_KUCOIN   = "kucoin"
	EXCHANGE_BITZ     = "bitz"
	EXCHANGE_MCX      = "mxc"
	EXCHANGE_HOO      = "hoo"
	EXCHANGE_POLONIEX = "poloniex"
	EXCHANGE_BIKI     = "biki"
	EXCHANGE_HITBTC   = "hitbtc"
)
