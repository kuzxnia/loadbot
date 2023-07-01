package database

import (
	"net/http"
	"time"

	"github.com/kuzxnia/mongoload/pkg/config"
)

func HTTPClient(config *config.Config) *http.Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          int(config.ConcurrentConnections),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}
}

// TODO: new fasthttp client
