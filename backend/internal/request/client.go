package request

import (
	"net/url"
	"time"

	"github.com/plzn0/go-http-1.17.1"
	"github.com/plzn0/go-http-1.17.1/cookiejar"
)

// func init() {
// 	err := rootcerts.UpdateDefaultTransport()
// 	if err != nil {
// 		panic(err)
// 	}
// }

type Options struct {
	DisableDecompression bool
	UserAgent            string
	TLSServerNameSpoofs  map[string]string
	InsecureSkipVerify   bool
	HTTP2FrameSettings   *http.FrameSettings
}

type Client struct {
	http.Client
}

func NewClient(opts *Options, proxy string) (*Client, error) {
	var proxyUrl *url.URL
	var err error
	if proxy != "" {
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			return nil, err
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: http.Client{
			Jar:       jar,
			Transport: NewTransport(opts, proxyUrl),
			Timeout:   10 * time.Second,
		},
	}, nil
}
