package request

import (
	"github.com/plzn0/go-http-1.17.1"
	"net/url"
)

func Proxy(req *http.Request) (*url.URL, error) {
	proxy, err := http.ProxyFromEnvironment(req)
	if proxy != nil || err != nil {
		return proxy, err
	}

	return nil, nil
}
