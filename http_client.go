package goexchange

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HTTPClientConfig http client config
type HTTPClientConfig struct {
	HTTPTimeout  time.Duration
	Proxy        *url.URL
	MaxIdleConns int
}

// String return string's config
func (c HTTPClientConfig) String() string {
	return fmt.Sprintf("{ProxyUrl:\"%s\",HttpTimeout:%s,MaxIdleConns:%d}", c.Proxy, c.HTTPTimeout.String(), c.MaxIdleConns)
}

// SetHTTPTimeout set client timeout
func (c *HTTPClientConfig) SetHTTPTimeout(timeout time.Duration) *HTTPClientConfig {
	c.HTTPTimeout = timeout
	return c
}

// SetProxyURL set client proxy
func (c *HTTPClientConfig) SetProxyURL(proxyURL string) *HTTPClientConfig {
	if proxyURL == "" {
		return c
	}
	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return c
	}
	c.Proxy = proxy
	return c
}

// SetMaxIdleConns set client max idle connect number
func (c *HTTPClientConfig) SetMaxIdleConns(max int) *HTTPClientConfig {
	c.MaxIdleConns = max
	return c
}

var (
	// DefaultHTTPClientConfig default config
	DefaultHTTPClientConfig = &HTTPClientConfig{
		Proxy:        nil,
		HTTPTimeout:  5 * time.Second,
		MaxIdleConns: 10,
	}
)

// NewHTTPClient new http client instance
func NewHTTPClient() (client *http.Client) {
	client = &http.Client{
		Timeout: DefaultHTTPClientConfig.HTTPTimeout,
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return DefaultHTTPClientConfig.Proxy, nil
			},
			MaxIdleConns:          DefaultHTTPClientConfig.MaxIdleConns,
			IdleConnTimeout:       5 * DefaultHTTPClientConfig.HTTPTimeout,
			MaxConnsPerHost:       2,
			MaxIdleConnsPerHost:   2,
			TLSHandshakeTimeout:   DefaultHTTPClientConfig.HTTPTimeout,
			ResponseHeaderTimeout: DefaultHTTPClientConfig.HTTPTimeout,
			ExpectContinueTimeout: DefaultHTTPClientConfig.HTTPTimeout,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, DefaultHTTPClientConfig.HTTPTimeout)
			}},
	}
	return
}

// NewHTTPClientWithConfig new http client instance with config
func NewHTTPClientWithConfig(config *HTTPClientConfig) (client *http.Client) {
	if config == nil {
		config = DefaultHTTPClientConfig
	}

	client = &http.Client{
		Timeout: config.HTTPTimeout,
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (*url.URL, error) {
				return config.Proxy, nil
			},
			MaxIdleConns:          config.MaxIdleConns,
			IdleConnTimeout:       5 * config.HTTPTimeout,
			MaxConnsPerHost:       2,
			MaxIdleConnsPerHost:   2,
			TLSHandshakeTimeout:   config.HTTPTimeout,
			ResponseHeaderTimeout: config.HTTPTimeout,
			ExpectContinueTimeout: config.HTTPTimeout,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, config.HTTPTimeout)
			}},
	}
	return
}
