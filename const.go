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
type KlinePeriod int

const (
	KLINE_PERIOD_1MIN = 1 + iota
	KLINE_PERIOD_3MIN
	KLINE_PERIOD_5MIN
	KLINE_PERIOD_15MIN
	KLINE_PERIOD_30MIN
	KLINE_PERIOD_60MIN
	KLINE_PERIOD_1H
	KLINE_PERIOD_2H
	KLINE_PERIOD_3H
	KLINE_PERIOD_4H
	KLINE_PERIOD_6H
	KLINE_PERIOD_8H
	KLINE_PERIOD_12H
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_3DAY
	KLINE_PERIOD_1WEEK
	KLINE_PERIOD_1MONTH
	KLINE_PERIOD_1YEAR
)