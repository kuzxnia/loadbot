package database

import (
	"net/http"
	"time"

	"github.com/kuzxnia/loadbot/lbot/config"
)

func HTTPClient(cfg *config.Job) *http.Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          int(cfg.Connections),
		IdleConnTimeout:       cfg.Timeout,
		TLSHandshakeTimeout:   cfg.Timeout,
		ExpectContinueTimeout: cfg.Timeout,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}
}


func FastHTTPClient(cfg *config.Job) *http.Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          int(cfg.Connections),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}
}

// TODO: new fasthttp client
