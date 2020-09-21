package goexchange

import "net/http"

type APIConfig struct {
	HttpClient    *http.Client
	Endpoint      string
	ApiKey        string
	ApiSecretKey  string
	ApiPassphrase string
	ClientId      string
	Proxy      	  string
}