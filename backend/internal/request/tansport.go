package request

import (
	"crypto/tls"
	"net/url"
	"strings"

	"github.com/plzn0/go-http-1.17.1"
	utls "github.com/plzn0/go-utls"
)

type Transport struct {
	*http.Transport
	*Options
}

type TLSConfig interface {
	GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

func NewTransport(opts *Options, proxy *url.URL) *Transport {
	var proxyTransport func(*http.Request) (*url.URL, error)

	if proxy == nil {
		proxyTransport = Proxy
	} else {
		proxyTransport = http.ProxyURL(proxy)
	}

	upstream := &http.Transport{
		Proxy: proxyTransport,
		// Disables Go's built in decompression
		DisableCompression: true,
		ForceAttemptHTTP2:  true,
		TLSClientConfig: &utls.Config{
			TLSServerNameSpoofs: opts.TLSServerNameSpoofs,
			InsecureSkipVerify:  opts.InsecureSkipVerify,
			// UserAgent: "Go-http-client/1.1",
			
		},
		ClientHelloID:      utls.HelloIOS_12_1,
		HTTP2FrameSettings: *http.FrameSettingsChrome,
		// chromePseudoHeaders  = []string{":method", ":authority", ":scheme", ":path"}
		// fireFoxPsuedoHeaders = []string{":method", ":path", ":authority", ":scheme"}
		PseudoHeaderOrder: []string{":method", ":authority", ":scheme", ":path"},
		UserAgent: opts.UserAgent,
	}

	transport := &Transport{upstream, opts}

	// err := http2.ConfigureTransport(upstream)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return transport
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	hasUserAgent := false
	for k, _ := range req.Header {
		if strings.ToLower(k) == "user-agent" {
			hasUserAgent = true
			break
		}
	}
	if !hasUserAgent {
		if t.Options.UserAgent != "" {
			req.Header.Set("User-Agent", t.Options.UserAgent)
		}
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if !t.DisableDecompression {
		err := Decompress(resp)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
