package request

import (
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"

	"github.com/plzn0/go-http-1.17.1"

	"github.com/andybalholm/brotli"
)

var ErrUnknownEncoding = errors.New("response has unknown encoding")

type initReader func() (io.Reader, error)

type Decompressor struct {
	io.ReadCloser

	init   initReader
	reader io.Reader
	body   io.ReadCloser
	err    error
}

func Decompress(resp *http.Response) error {
	if !resp.Uncompressed {
		var init initReader

		body := resp.Body

		switch resp.Header.Get("content-encoding") {
		case "":
			return nil
		case "gzip":
			init = func() (io.Reader, error) {
				return gzip.NewReader(body)
			}
		case "deflate":
			init = func() (io.Reader, error) {
				return flate.NewReader(body), nil
			}
		case "br":
			init = func() (io.Reader, error) {
				return brotli.NewReader(body), nil
			}

		default:
			return ErrUnknownEncoding
		}

		resp.Body = &Decompressor{
			body: body,
			init: init,
		}
	}

	return nil
}

func (d *Decompressor) Read(p []byte) (int, error) {
	if d.reader == nil {
		if d.err == nil {
			d.reader, d.err = d.init()
			d.init = nil
		}

		if d.err != nil {
			return 0, d.err
		}
	}

	return d.reader.Read(p)
}

func (d *Decompressor) Close() error {
	return d.body.Close()
}
