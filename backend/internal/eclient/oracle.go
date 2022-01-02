package eclient

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"

	"github.com/plzn0/go-http-1.17.1"
)

// GasOracleResponse is the response sent back by the oracle
type GasOracleResponse struct {
	Code int        `json:"code"`
	Data OracleData `json:"data" `
}

// OracleData contains gas estimates for 
type OracleData struct {
	Rapid     int64   `json:"rapid"`
	Fast      int64   `json:"fast"`
	Standard  int64   `json:"standard"`
	Slow      int64   `json:"slow"`
	TimeStamp int64   `json:"timestamp"`
	USDPrice  float64 `json:"priceUSD"`
}

// GetGasEstimate fetches a gas estimate from etherchain
func (c *Client) GetGasEstimate() (*OracleData, error) {
	// https://etherchain.org/api/gasnow
	ep, _ := url.Parse("https://etherchain.org/api/gasnow")

	headers := http.Header{}
	headers.Add("Accept", "application/json")

	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	res, err := c.request.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &GasOracleResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, errors.New("failed to fetch gas price")
	}

	return &result.Data, nil
}
