package goexchange

import "net/http"

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string
	AccountId     string
	Proxy      	  string
}

type HttpClientResponse struct {
	Code 		int			`json:"code"`
	Msg 		string		`json:"msg"`
	Error 		string		`json:"error"`
	St 			int64		`json:"st"`
	Et 			int64		`json:"et"`
	Data 		[]byte		`json:"data"`
}

type PlaceOrder struct {
	Symbol 		    Symbol
	ClientOrderId   string
	Price 		    string
	Amount 	        string
	Side            TradeSide
	TradeType       string
	TimeInForce     TimeInForce
	options         map[string]string
}

type LimitOrder struct {
	Symbol 		    Symbol
	ClientOrderId   string
	Price 		    string
	Amount 	        string
	Side            TradeSide
	TimeInForce     TimeInForce
	options         map[string]string
}
