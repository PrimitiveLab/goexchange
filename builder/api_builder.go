package builder

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	. "github.com/primitivelab/goexchange"
	"github.com/primitivelab/goexchange/biki"
	"github.com/primitivelab/goexchange/binance"
	"github.com/primitivelab/goexchange/bitz"
	"github.com/primitivelab/goexchange/gate"
	"github.com/primitivelab/goexchange/hoo"
	"github.com/primitivelab/goexchange/huobi"
	"github.com/primitivelab/goexchange/mxc"
	"github.com/primitivelab/goexchange/okex"
	"github.com/primitivelab/goexchange/poloniex"
)

type APIBuilder struct {
	HttpClientConfig *HttpClientConfig
	client           *http.Client
	httpTimeout      time.Duration
	apiKey           string
	secretKey        string
	accountId        string
	passphrase       string
	endPoint         string
}

type HttpClientConfig struct {
	HttpTimeout  time.Duration
	Proxy        *url.URL
	MaxIdleConns int
}

func (c HttpClientConfig) String() string {
	return fmt.Sprintf("{ProxyUrl:\"%s\",HttpTimeout:%s,MaxIdleConns:%d}", c.Proxy, c.HttpTimeout.String(), c.MaxIdleConns)
}

func (c *HttpClientConfig) SetHttpTimeout(timeout time.Duration) *HttpClientConfig {
	c.HttpTimeout = timeout
	return c
}

func (c *HttpClientConfig) SetProxyUrl(proxyUrl string) *HttpClientConfig {
	if proxyUrl == "" {
		return c
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return c
	}
	c.Proxy = proxy
	return c
}

func (c *HttpClientConfig) SetMaxIdleConns(max int) *HttpClientConfig {
	c.MaxIdleConns = max
	return c
}

var (
	DefaultHttpClientConfig = &HttpClientConfig{
		Proxy:        nil,
		HttpTimeout:  5 * time.Second,
		MaxIdleConns: 10}
	DefaultAPIBuilder = NewAPIBuilder()
)

func NewAPIBuilder() (builder *APIBuilder) {
	return NewAPIBuilderWithHttpClientConfig(DefaultHttpClientConfig)
}

func NewAPIBuilderWithHttpClientConfig(config *HttpClientConfig) *APIBuilder {
	if config == nil {
		config = DefaultHttpClientConfig
	}

	return &APIBuilder{
		HttpClientConfig: config,
		client: &http.Client{
			Timeout: config.HttpTimeout,
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return config.Proxy, nil
				},
				MaxIdleConns:          config.MaxIdleConns,
				IdleConnTimeout:       5 * config.HttpTimeout,
				MaxConnsPerHost:       2,
				MaxIdleConnsPerHost:   2,
				TLSHandshakeTimeout:   config.HttpTimeout,
				ResponseHeaderTimeout: config.HttpTimeout,
				ExpectContinueTimeout: config.HttpTimeout,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, config.HttpTimeout)
				}},
		}}
}

func NewAPIBuilderWithCustomHttpClient(client *http.Client) (builder *APIBuilder) {
	return &APIBuilder{client: client}
}

func (builder *APIBuilder) GetHttpClientConfig() *HttpClientConfig {
	return builder.HttpClientConfig
}

func (builder *APIBuilder) GetHttpClient() *http.Client {
	return builder.client
}

func (builder *APIBuilder) HttpProxy(proxyUrl string) (_builder *APIBuilder) {
	if proxyUrl == "" {
		return builder
	}
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return builder
	}
	builder.HttpClientConfig.Proxy = proxy
	transport := builder.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxy)
	return builder
}

func (builder *APIBuilder) HttpTimeout(timeout time.Duration) (_builder *APIBuilder) {
	builder.HttpClientConfig.HttpTimeout = timeout
	builder.httpTimeout = timeout
	builder.client.Timeout = timeout
	transport := builder.client.Transport.(*http.Transport)
	if transport != nil {
		// transport.ResponseHeaderTimeout = timeout
		// transport.TLSHandshakeTimeout = timeout
		transport.IdleConnTimeout = timeout
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}
	return builder
}

func (builder *APIBuilder) APIKey(key string) (_builder *APIBuilder) {
	builder.apiKey = key
	return builder
}

func (builder *APIBuilder) APISecretKey(key string) (_builder *APIBuilder) {
	builder.secretKey = key
	return builder
}

func (builder *APIBuilder) AccountId(id string) (_builder *APIBuilder) {
	builder.accountId = id
	return builder
}

func (builder *APIBuilder) Passphrase(passphrase string) (_builder *APIBuilder) {
	builder.passphrase = passphrase
	return builder
}

func (builder *APIBuilder) Endpoint(endpoint string) (_builer *APIBuilder) {
	builder.endPoint = endpoint
	return builder
}

func (builder *APIBuilder) Build(exName string) (api SpotAPI) {
	config := APIConfig{}
	config.HttpClient = builder.client
	config.ApiKey = builder.apiKey
	config.ApiSecretKey = builder.secretKey
	config.ApiPassphrase = builder.passphrase
	config.Endpoint = builder.endPoint
	config.AccountId = builder.accountId

	switch exName {
	case ECHANGE_BINANCE:
		api = binance.NewWithConfig(&config)
	case ECHANGE_HUOBI:
		api = huobi.NewWithConfig(&config)
	case ECHANGE_OKEX:
		api = okex.NewWithConfig(&config)
	case ECHANGE_GATE:
		api = gate.NewWithConfig(&config)
	case ECHANGE_BITZ:
		api = bitz.NewWithConfig(&config)
	case ECHANGE_MCX:
		api = mxc.NewWithConfig(&config)
	case ECHANGE_HOO:
		api = hoo.NewWithConfig(&config)
	case ECHANGE_BIKI:
		api = biki.NewWithConfig(&config)
	case ECHANGE_POLONIEX:
		api = poloniex.NewWithConfig(&config)
	default:
		println("exchange name error [" + exName + "].")
	}
	return api
}
