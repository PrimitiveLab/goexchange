package goexchange

import (
	"strings"
)

type Symbol struct {
	CoinFrom string
	CoinTo   string
}

func NewSymbol(coinFrom string, coinTo string) Symbol {
	return Symbol{coinFrom, coinTo}
}

func (symbol Symbol) String() string {
	return symbol.ToSymbol("_")
}

func (symbol Symbol) ToSymbol(sep string) string {
	return strings.Join([]string{symbol.CoinFrom, symbol.CoinTo}, sep)
}

func (symbol Symbol) Reverse() Symbol {
	return Symbol{symbol.CoinTo, symbol.CoinFrom}
}

func (symbol Symbol) ToLower() Symbol {
	return Symbol{strings.ToLower(symbol.CoinFrom), strings.ToLower(symbol.CoinTo)}
}

func (symbol Symbol) ToUpper() Symbol {
	return Symbol{strings.ToUpper(symbol.CoinFrom), strings.ToUpper(symbol.CoinTo)}
}