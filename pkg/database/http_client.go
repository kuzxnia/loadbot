package database

import (
	"net/http"
	"time"
)

func HTTPClient(maxIdleConns uint64) *http.Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          int(maxIdleConns),
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Transport: transport,
		// timeout
	}
}

// TODO: new fasthttp client
