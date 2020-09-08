/*  transport.go
*
* @Author:             Nanang Suryadi
* @Date:               January 19, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 19/01/20 03:27
 */

package mimir

import (
	"net"
	"net/http"
	"time"
)

func DefaultTransport(duration int) *http.Transport {
	transport := DefaultPooledTransport(duration)
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

func DefaultPooledTransport(duration int) *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(duration) * time.Second,
			KeepAlive: time.Duration(duration) * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}

func DefaultClient(duration int) *http.Client {
	return &http.Client{
		Transport: DefaultTransport(duration),
	}
}
