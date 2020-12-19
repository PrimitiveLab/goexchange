package goexchange

import "time"

type ApiStatusCode struct {
	Code int
	Msg  string
}

var (
	HttpClientInternalError = ApiStatusCode{Code: 1001, Msg: "http client internal error"}
	JsonUnmarshalError      = ApiStatusCode{Code: 1002, Msg: "http response data unmarshal(json) error"}
	ExchangeError           = ApiStatusCode{Code: 1003, Msg: "exchange api error"}
	HttpRequestError        = ApiStatusCode{Code: 404, Msg: "http request error"}
	MethodNotExistError     = ApiStatusCode{Code: 1004, Msg: "method is not exist"}

	// HTTP_ERR_CODE                = ApiError{Code: "HTTP_ERR_0001", Msg: "http request error"}
	// EX_ERR_API_LIMIT             = ApiError{Code: "EX_ERR_1000", Msg: "api limited"}
	// EX_ERR_SIGN                  = ApiError{Code: "EX_ERR_0001", Msg: "signature error"}
	// EX_ERR_NOT_FIND_SECRETKEY    = ApiError{Code: "EX_ERR_0002", Msg: "not find secretkey"}
	// EX_ERR_NOT_FIND_APIKEY       = ApiError{Code: "EX_ERR_0003", Msg: "not find apikey"}
	// EX_ERR_INSUFFICIENT_BALANCE  = ApiError{Code: "EX_ERR_0004", Msg: "Insufficient Balance"}
	// EX_ERR_PLACE_ORDER_FAIL      = ApiError{Code: "EX_ERR_0005", Msg: "place order failure"}
	// EX_ERR_CANCEL_ORDER_FAIL     = ApiError{Code: "EX_ERR_0006", Msg: "cancel order failure"}
	// EX_ERR_INVALID_CURRENCY_PAIR = ApiError{Code: "EX_ERR_0007", Msg: "invalid currency pair"}
	// EX_ERR_NOT_FIND_ORDER        = ApiError{Code: "EX_ERR_0008", Msg: "not find order"}
	// EX_ERR_SYMBOL_ERR            = ApiError{Code: "EX_ERR_0009", Msg: "symbol error"}
	//
	// HTTP_ERROR_401            	= ApiError{Code: "401", Msg: "http request 401 error"}
	// HTTP_ERROR_404            	= ApiError{Code: "404", Msg: "http request 404 error"}
	// HTTP_ERROR_502            	= ApiError{Code: "502", Msg: "http request 502 error"}
)

// ReturnAPIError return error
func ReturnAPIError(code ApiStatusCode) interface{} {
	startTime := time.Now().UnixNano() / 1e6
	retData := map[string]interface{}{
		"code":  code.Code,
		"st":    startTime,
		"et":    startTime,
		"msg":   code.Msg,
		"error": code.Msg,
		"data":  nil}
	return retData
}
