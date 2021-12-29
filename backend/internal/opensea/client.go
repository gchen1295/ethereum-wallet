package piratesea

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	minter "nft-engine/internal/deth"
	queries "nft-engine/internal/opensea/queries"
	"nft-engine/internal/request"
	"nft-engine/pkg/opensea"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/plzn0/go-http-1.17.1"
)

// GenerateRandomString generates a secure random string of length n.
func GenerateRandomString(n int) (string, error) {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

type Bot struct {
	request *request.Client

	hostURL    *url.URL
	OpenSea    *opensea.Opensea
	Erc721Abi  *abi.ABI
	Erc1155Abi *abi.ABI
	wyvernAbi  *abi.ABI
	ethClient  *minter.Bot
}

type Network int

var (
	Main    = Network(0)
	Rinkeby = Network(1)
)

func NewBot(network Network, client *minter.Bot) (*Bot, error) {
	bot := &Bot{
		request: request.NewClient(&request.Options{
			DisableDecompression: false,
			UserAgent:            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
		}, nil),
	}
	contractAddress := OPENSEA_ADDRESS
	if network == Rinkeby {
		contractAddress = RINKEBY_OPENSEA_CONTRACT_ADDRESS
	}

	os, err := opensea.NewOpensea(contractAddress, client.Client)
	if err != nil {
		return nil, err
	}
	bot.OpenSea = os

	erc721, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return nil, errors.New("failed to read abi for erc 721 contracts")
	}
	bot.Erc721Abi = &erc721

	erc1155, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"internalType":"address","name":"_owner","type":"address"},{"internalType":"uint256","name":"_id","type":"uint256"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address[]","name":"_owners","type":"address[]"},{"internalType":"uint256[]","name":"_ids","type":"uint256[]"}],"name":"balanceOfBatch","outputs":[{"internalType":"uint256[]","name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"_operator","type":"address"},{"internalType":"bool","name":"_approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"_owner","type":"address"},{"internalType":"address","name":"_operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"_from","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_id","type":"uint256"},{"internalType":"uint256","name":"_value","type":"uint256"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"address","name":"_from","type":"address"},{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256[]","name":"_ids","type":"uint256[]"},{"internalType":"uint256[]","name":"_values","type":"uint256[]"},{"internalType":"bytes","name":"_data","type":"bytes"}],"name":"safeBatchTransferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return nil, errors.New("failed to load erc1155 abi")
	}
	bot.Erc1155Abi = &erc1155
	bot.ethClient = client

	data, err := ioutil.ReadFile("internal\\opensea\\wyvernABI.json")
	if err != nil {
		return nil, err
	}

	wyvern, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	bot.wyvernAbi = &wyvern

	switch network {
	case Main:
		bot.hostURL = API_BASE_URL_MAIN
	case Rinkeby:
		bot.hostURL = API_BASE_URL_RINKEBY
	default:
		bot.hostURL = API_BASE_URL_MAIN
	}
	return bot, nil
}

// &url.URL{
// 	Host: "192.168.0.20:8888",
// }

func (c *Bot) PullCollection(collectionSlug string) error {
	endpoint := &url.URL{}
	endpoint.Host = c.hostURL.Host
	endpoint.Scheme = c.hostURL.Scheme
	endpoint.Path = "api/v1/events"
	query := endpoint.Query()
	query.Add("collection_slug", "nope")
	endpoint.RawQuery = query.Encode()

	ep, _ := url.Parse(fmt.Sprintf("https://api.opensea.io/api/v1/events?collection_slug=%s", collectionSlug))

	origin, err := GenerateRandomString(32)
	if err != nil {
		origin = ""
	}
	headers := http.Header{}
	headers.Add("accept-encoding", "gzip, deflate, br")
	headers.Add("origin", origin)

	headers.Add("sec-ch-ua", `"Chromium";v="94", "Google Chrome";v="94", ";Not A Brand";v="99"`)
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("accept", "application/json")
	headers.Add("accept-language", "en-US,en;q=0.9")
	headers.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")

	req := &http.Request{
		Method: http.MethodGet,
		URL:    ep,
		Header: headers,
	}
	res, err := c.request.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	formatted := &bytes.Buffer{}
	err = json.Indent(formatted, body, "", "    ")
	if err != nil {
		return err
	}

	log.Print(c.request.Jar)
	// log.Printf("%+v", res.Header)

	return nil
}

func (c *Bot) PullCollectionGql(collectionSlug string) (*CollectionData, []CollectionQueryAssetsResults, error) {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "graphql/",
	})

	sortby := "SEVEN_DAY_VOLUME"
	sortByAscending := true
	isSingleCollection := true
	isCategory := false
	sortBy := "PRICE"
	payload := &CollectionQuery{
		ID:    "collectionQuery",
		Query: collectionQuery,
		Variables: CollectionQueryVariables{
			CollectionSlug:     collectionSlug,
			Collections:        []string{collectionSlug},
			CollectionSortBy:   &sortby,
			EventTypes:         &[]string{"AUCTION_SUCCESSFUL"},
			IsSingleCollection: &isSingleCollection,
			IsCategory:         &isCategory,
			SortBy:             &sortBy,
			SortAscending:      &sortByAscending,
			ShowContextMenu:    &isSingleCollection,
		},
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(encodedPayload))
	if err != nil {
		return nil, nil, err
	}

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "*/*")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("x-signed-query", "5abaaf9e000e594d1874e473da12dd3cb54766a88c6010631dd9d41254e84ab3")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")

	req.Header = headers
	res, err := c.request.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	result := &CollectionQueryResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, nil, err
	}

	return result.Data.Collection, result.Data.Assets.Search.Results, nil
}

func (c *Bot) WatchCollectionGql(collectionSlug string) error {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "graphql/",
	})
	timeStamp := getCurrentTS()

	payload := &EventHistoryQuery{
		ID:    "EventHistoryPollQuery",
		Query: eventHistoryPollQuery,
		Variables: EventHistoryQueryVariables{
			Collections:    []string{collectionSlug},
			Count:          100,
			EventTypes:     []string{"AUCTION_CREATED"},
			EventTimestamp: timeStamp,
			ShowAll:        true,
		},
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "*/*")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("x-signed-query", "f903d213d9dd47cb4a723d19122e984c5dbf37cbc441379ff87ab2784876f31a")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")

	timer := time.NewTicker(time.Second * 4)
	resetMapTimer := time.Now().Add(time.Minute * 15)

	// eventMap is a mapping of the recent event ID's to prevent duplicates
	eventMap := map[string]bool{}
	eventsMux := sync.RWMutex{}
	for range timer.C {
		if time.Now().After(resetMapTimer) {
			eventsMux.Lock()
			eventMap = map[string]bool{}
			resetMapTimer = time.Now().Add(time.Minute * 15)
			eventsMux.Unlock()
		}

		req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(encodedPayload))
		if err != nil {
			return err
		}

		payload.Variables.EventTimestamp = timeStamp
		encodedPayload, err = json.Marshal(payload)
		if err != nil {
			return err
		}

		headers["Content-Length"] = []string{strconv.Itoa(len(encodedPayload))}
		req.Header = headers

		timeStamp = getCurrentTS()
		res, err := c.request.Do(req)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		res.Body.Close()
		result := &EventHistoryPollResponse{}
		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}

		edges := result.Data.AssetEvents.Edges
		fmt.Println("Monitoring...")
		if len(edges) > 0 {
			fmt.Printf("%d events found\n", len(edges))
			embeds := []Embed{}
			for _, e := range edges {
				relayId := e.Node.RelayID
				if relayId == "" {
					continue
				}
				eventsMux.Lock()
				if _, ok := eventMap[relayId]; ok {
					continue
				}
				eventMap[relayId] = true
				eventsMux.Unlock()
				osLink := fmt.Sprintf("https://opensea.io/assets/%s/%s", e.Node.AssetQuantity.Asset.AssetContract.Address, e.Node.AssetQuantity.Asset.TokenID)
				price, err := strconv.ParseFloat(e.Node.Price.QuantityInEth, 64)
				if err != nil {
					price = 0
				}
				var convertedPrice float64

				if price > 0 {
					decimals := float64(math.Pow10(e.Node.Price.Asset.Decimals))
					convertedPrice = price / decimals
				}
				priceString := "Not found"
				if convertedPrice >= 0 {
					priceString = fmt.Sprintf("Ξ%.4f", convertedPrice)
				}

				embeds = append(embeds, Embed{
					Author: Author{
						Name:    fmt.Sprintf("New Listing by %s", e.Node.Seller.User.PublicUsername),
						URL:     osLink,
						IconURL: e.Node.Seller.ImageURL,
					},
					Color: 0x00ff00,
					Image: Image{
						URL: e.Node.AssetQuantity.Asset.DisplayImageURL,
					},
					Footer: Footer{
						Text: "Made with ❤️ by @undefined#6969",
					},
					Fields: []Fields{
						{
							Name:   "Item Name",
							Value:  e.Node.AssetQuantity.Asset.Name,
							Inline: false,
						},
						{
							Name:   "Token ID",
							Value:  e.Node.AssetQuantity.Asset.TokenID,
							Inline: true,
						},
						{
							Name:   "Price",
							Value:  priceString,
							Inline: true,
						},
						{
							Name:   "Link",
							Value:  fmt.Sprintf("[Here](%s)", osLink),
							Inline: false,
						},
						{
							Name:   "Time",
							Value:  fmt.Sprintf("<t:%d>", time.Now().Unix()),
							Inline: true,
						},
					},
				})
			}

			if len(embeds) > 0 {
				fmt.Println(len(embeds))
				fmt.Println(embeds)
				if err = c.sendWebhook(embeds); err != nil {
					fmt.Print(err)
				}
			}
		}
	}

	return nil
}

func (c *Bot) EventHistoryPoll(collectionSlug string, timestamp string) ([]Edges, string, error) {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "graphql/",
	})

	payload := &EventHistoryQuery{
		ID:    "EventHistoryPollQuery",
		Query: eventHistoryPollQuery,
		Variables: EventHistoryQueryVariables{
			Collections:    []string{collectionSlug},
			Count:          100,
			EventTypes:     []string{"AUCTION_CREATED"},
			EventTimestamp: timestamp,
			ShowAll:        true,
		},
	}

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "*/*")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("x-signed-query", "f903d213d9dd47cb4a723d19122e984c5dbf37cbc441379ff87ab2784876f31a")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(encodedPayload))
	if err != nil {
		return nil, timestamp, err
	}

	payload.Variables.EventTimestamp = timestamp
	encodedPayload, err = json.Marshal(payload)
	if err != nil {
		return nil, timestamp, err
	}

	headers["Content-Length"] = []string{strconv.Itoa(len(encodedPayload))}
	req.Header = headers

	res, err := c.request.Do(req)
	if err != nil {
		return nil, timestamp, err
	}
	defer res.Body.Close()
	newTS := getCurrentTS()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, timestamp, err
	}

	result := &EventHistoryPollResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, timestamp, err
	}

	return result.Data.AssetEvents.Edges, newTS, nil
}

// OrdersQuery queries the latest orders by ascending price
func (c *Bot) OrdersQuery(contractAddress string, tokenId string) ([]*queries.OrderQuery, error) {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "graphql/",
	})

	payload, signed := queries.NewOrdersQueryPayload(contractAddress, tokenId)

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "*/*")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("x-signed-query", signed)
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(encodedPayload))
	if err != nil {
		return nil, err
	}

	encodedPayload, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers["Content-Length"] = []string{strconv.Itoa(len(encodedPayload))}
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

	parsedResponse := &queries.OrderQueryResponse{}
	err = json.Unmarshal(body, parsedResponse)
	if err != nil {
		return nil, err
	}

	results := []*queries.OrderQuery{}
	for _, order := range parsedResponse.Data.Orders.Edges {
		results = append(results, &order.Node)
	}
	log.Println(string(body))
	return results, nil
}

// AssetSearch queries the latest orders by ascending price
func (c *Bot) AssetSearch(collectionSlug string) ([]*queries.AssetSearchResult, error) {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "graphql/",
	})

	payload, signed := queries.NewAssetSearchPayload(collectionSlug)

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "*/*")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("x-signed-query", signed)
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")

	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(encodedPayload))
	if err != nil {
		return nil, err
	}

	encodedPayload, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers["Content-Length"] = []string{strconv.Itoa(len(encodedPayload))}
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

	parsedResponse := &queries.AssetSearchResponse{}
	err = json.Unmarshal(body, parsedResponse)
	if err != nil {
		return nil, err
	}

	results := []*queries.AssetSearchResult{}
	for _, order := range parsedResponse.Data.Query.Search.Edges {
		results = append(results, &order.Node)
	}

	return results, nil
}

func (c *Bot) RetrieveOrder(contractAddr, tokenID string) ([]*queries.Order, error) {
	endpoint := c.hostURL.ResolveReference(&url.URL{
		Path: "wyvern/v1/orders",
		RawQuery: fmt.Sprintf("asset_contract_address=%s&token_id=%s", contractAddr, tokenID),
	})

	headers := http.Header{}
	headers.Add("Connection", "keep-alive")
	// headers.Add("Content-Length", strconv.Itoa(len(encodedPayload)))
	headers.Add("Accept", "application/json")
	// headers.Add("X-BUILD-ID", "B7__n40nMoFogTBgPIhJD")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36")
	headers.Add("X-API-KEY", "2f6f419a083c46de9d83ce3dbe7db601")
	// headers.Add("Content-Type", "application/json")
	headers.Add("Sec-GPC", "1")
	headers.Add("Origin", "https://opensea.io")
	headers.Add("Sec-Fetch-Site", "same-site")
	headers.Add("Sec-Fetch-Mode", "cors")
	headers.Add("Sec-Fetch-Dest", "empty")
	headers.Add("Referer", "https://opensea.io/")
	// headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")
	log.Println(endpoint.String())

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.request.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	orderResponse := struct {
		Count  int              `json:"count"`
		Orders []*queries.Order `json:"orders"`
	}{}
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		return nil, err
	}

	log.Println(len(orderResponse.Orders))
	log.Println(orderResponse.Orders[0])

	return orderResponse.Orders, nil
}
