package goexchange


// 现货交易方向
type TradeSide int
const (
	LIMIT_BUY TradeSide = iota + 1
	LIMIT_SELL
	MARKET_BUY
	MARKET_SELL
)

func (ts TradeSide) String() string {
	switch ts {
	case 1:
		return "limit-buy"
	case 2:
		return "limit-sell"
	case 3:
		return "market-buy"
	case 4:
		return "market-sell"
	default:
		return "unknown"
	}
}

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
	ECHANGE_OKEX       	= "okex"
	ECHANGE_HUOBI   	= "huobi"
	ECHANGE_BINANCE    	= "binance"
	ECHANGE_GATE      	= "gate"
	ECHANGE_KUCOIN  	= "kucoin"
	ECHANGE_BITZ  		= "bitz"
)