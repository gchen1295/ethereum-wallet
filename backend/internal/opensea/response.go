package piratesea

import (
	queries "nft-engine/internal/opensea/queries"
)

// CollectionQuery
type CollectionData struct {
	IsEditable          bool                  `json:"isEditable"`
	BannerImageURL      string                `json:"bannerImageUrl"`
	Name                string                `json:"name"`
	Description         string                `json:"description"`
	ImageURL            string                `json:"imageUrl"`
	RelayID             string                `json:"relayId"`
	RepresentativeAsset RepresentativeAddress `json:"representativeAsset"`
	Slug                string                `json:"slug"`
	IsMintable          bool                  `json:"isMintable"`
	IsSafelisted        bool                  `json:"isSafelisted"`
	IsVerified          bool                  `json:"isVerified"`
	Owner               ContractOwner         `json:"owner"`
	Stats               CollectionStats       `json:"stats"`
	FloorPrice          string                `json:"floorPrice"`
	DiscordURL          string                `json:"discordUrl"`
	ExternalURL         string                `json:"externalUrl"`
	InstagramUsername   string                `json:"instagramUsername"`
	MediumUsername      string                `json:"mediumUsername"`
	TelegramURL         string                `json:"telegramUrl"`
	TwitterUsername     string                `json:"twitterUsername"`
	ID                  string                `json:"id"`
}

type RepresentativeAddress struct {
	AssetContract AssetContract `json:"assetContract"`
	ID            string        `json:"id"`
}

type ContractOwner struct {
	Address       string `json:"address"`
	Config        string `json:"config"`
	IsCompromised bool   `json:"isCompromised"`
	User          User   `json:"user"`
	ImageURL      string `json:"imageUrl"`
	ID            string `json:"id"`
}

type CollectionStats struct {
	NumOwners   float64 `json:"numOwners"`
	TotalSupply float64 `json:"totalSupply"`
	TotalVolume float64 `json:"totalVolume"`
	ID          string  `json:"id"`
}

type User struct {
	PublicUsername string `json:"publicUsername"`
	ID             string `json:"id"`
}

type CollectionQueryResponse struct {
	Data struct {
		Collection *CollectionData        `json:"collection"`
		Assets     *CollectionQueryAssets `json:"assets"`
	} `json:"data"`
}

type CollectionQuerySearchResult struct {
	Asset struct {
		AssetContract        AssetContract                `json:"assetContract"`
		Collection           queries.OrderAssetCollection `json:"collection"`
		RelayID              string                       `json:"relayId"`
		TokenID              string                       `json:"tokenId"`
		BackgroundColor      string                       `json:"backgroundColor"`
		ImageURL             string                       `json:"imageUrl"`
		Name                 string                       `json:"name"`
		ID                   string                       `json:"id"`
		IsDelisted           bool                         `json:"isDelisted"`
		AnimationURL         string                       `json:"animationUrl"`
		DisplayImageURL      string                       `json:"displayImageUrl"`
		Decimals             int                          `json:"decimals"`
		FavoritesCount       int                          `json:"favoritesCount"`
		IsFavorite           bool                         `json:"isFavorite"`
		IsFrozen             bool                         `json:"isFrozen"`
		HasUnlockableContent bool                         `json:"hasUnlockableContent"`
		OrderData            struct {
			BestAsk struct {
				RelayID   string `json:"relayId"`
				OrderType string `json:"orderType"`
				Maker     struct {
					Address string `json:"address"`
				} `json:"maker"`
				ClosedAt               string      `json:"closedAt"`
				DutchAuctionFinalPrice interface{} `json:"dutchAuctionFinalPrice"`
				OpenedAt               string      `json:"openedAt"`
				PriceFnEndedAt         interface{} `json:"priceFnEndedAt"`
				Quantity               string      `json:"quantity"`
				Decimals               string      `json:"decimals"`
				PaymentAssetQuantity   struct {
					Quantity string `json:"quantity"`
					Asset    struct {
						Decimals      int           `json:"decimals"`
						ImageURL      string        `json:"imageUrl"`
						Symbol        string        `json:"symbol"`
						UsdSpotPrice  float64       `json:"usdSpotPrice"`
						AssetContract AssetContract `json:"assetContract"`
						ID            string        `json:"id"`
					} `json:"asset"`
					ID            string `json:"id"`
					QuantityInEth string `json:"quantityInEth"`
				} `json:"paymentAssetQuantity"`
			} `json:"bestAsk"`
			BestBid interface{} `json:"bestBid"`
		} `json:"orderData"`
		IsEditable struct {
			Value  bool   `json:"value"`
			Reason string `json:"reason"`
		} `json:"isEditable"`
		IsListable bool        `json:"isListable"`
		Ownership  interface{} `json:"ownership"`
		Creator    struct {
			Address string `json:"address"`
			ID      string `json:"id"`
		} `json:"creator"`
		AssetEventData struct {
			LastSale struct {
				UnitPriceQuantity struct {
					Asset struct {
						Decimals      int           `json:"decimals"`
						ImageURL      string        `json:"imageUrl"`
						Symbol        string        `json:"symbol"`
						UsdSpotPrice  float64       `json:"usdSpotPrice"`
						AssetContract AssetContract `json:"assetContract"`
						ID            string        `json:"id"`
					} `json:"asset"`
					Quantity string `json:"quantity"`
					ID       string `json:"id"`
				} `json:"unitPriceQuantity"`
			} `json:"lastSale"`
		} `json:"assetEventData"`
	} `json:"asset"`
	AssetBundle interface{} `json:"assetBundle"`
	Typename    string      `json:"__typename"`
}

type CollectionQueryAssets struct {
	Search CollectionQuertSearchResults `json:"search"`
}

type CollectionQuertSearchResults struct {
	Results []CollectionQueryAssetsResults `json:"edges"`
}

type CollectionQueryAssetsResults struct {
	Node CollectionQuerySearchResult `json:"node"`
}

// EventHistoryPoll
type EventHistoryPollResponse struct {
	Data Data `json:"data"`
}

type DisplayData struct {
	CardDisplayStyle string `json:"cardDisplayStyle"`
}
type Collection struct {
	Name               string      `json:"name"`
	ID                 string      `json:"id"`
	DisplayData        DisplayData `json:"displayData"`
	Slug               string      `json:"slug"`
	IsMintable         bool        `json:"isMintable"`
	IsSafelisted       bool        `json:"isSafelisted"`
	IsVerified         bool        `json:"isVerified"`
	IsAuthorizedEditor bool        `json:"isAuthorizedEditor"`
	ImageURL           interface{} `json:"imageUrl"`

	RelayID string `json:"relayId"`
}
type Asset struct {
	RelayID         string        `json:"relayId"`
	AssetContract   AssetContract `json:"assetContract"`
	Collection      Collection    `json:"collection"`
	Name            string        `json:"name"`
	AnimationURL    string        `json:"animationUrl"`
	BackgroundColor string        `json:"backgroundColor"`
	IsDelisted      bool          `json:"isDelisted"`
	DisplayImageURL string        `json:"displayImageUrl"`
	TokenID         string        `json:"tokenId"`
	ID              string        `json:"id"`
	Decimals        int           `json:"decimals"`
	ImageURL        string        `json:"imageUrl"`
	Symbol          string        `json:"symbol"`
	UsdSpotPrice    float64       `json:"usdSpotPrice"`
}

type AssetQuantity struct {
	Asset    Asset  `json:"asset"`
	Quantity string `json:"quantity"`
	ID       string `json:"id"`
}

type AssetContract struct {
	Address           string `json:"address"`
	Chain             string `json:"chain"`
	ID                string `json:"id"`
	BlockExplorerLink string `json:"blockExplorerLink"`
	OpenseaVersion    string `json:"openseaVersion"`
}

type DevFee struct {
	Asset    Asset  `json:"asset"`
	Quantity string `json:"quantity"`
	ID       string `json:"id"`
}

type FromAccount struct {
	Address       string      `json:"address"`
	Config        interface{} `json:"config"`
	IsCompromised bool        `json:"isCompromised"`
	User          User        `json:"user"`
	ImageURL      string      `json:"imageUrl"`
	ID            string      `json:"id"`
}
type Price struct {
	Quantity      string `json:"quantity"`
	QuantityInEth string `json:"quantityInEth"`
	Asset         Asset  `json:"asset"`
	ID            string `json:"id"`
}
type EndingPrice struct {
	Quantity string `json:"quantity"`
	Asset    Asset  `json:"asset"`
	ID       string `json:"id"`
}
type Seller struct {
	Address       string      `json:"address"`
	Config        interface{} `json:"config"`
	IsCompromised bool        `json:"isCompromised"`
	User          User        `json:"user"`
	ImageURL      string      `json:"imageUrl"`
	ID            string      `json:"id"`
}
type Node struct {
	AssetBundle        interface{}   `json:"assetBundle"`
	AssetQuantity      AssetQuantity `json:"assetQuantity"`
	RelayID            string        `json:"relayId"`
	EventTimestamp     string        `json:"eventTimestamp"`
	EventType          string        `json:"eventType"`
	CustomEventName    string        `json:"customEventName"`
	OfferExpired       string        `json:"offerExpired"`
	DevFee             DevFee        `json:"devFee"`
	DevFeePaymentEvent interface{}   `json:"devFeePaymentEvent"`
	FromAccount        FromAccount   `json:"fromAccount"`
	Price              Price         `json:"price"`
	EndingPrice        EndingPrice   `json:"endingPrice"`
	Seller             Seller        `json:"seller"`
	ToAccount          string        `json:"toAccount"`
	WinnerAccount      string        `json:"winnerAccount"`
	Transaction        string        `json:"transaction"`
	ID                 string        `json:"id"`
}
type Edges struct {
	Node Node `json:"node"`
}
type AssetEvents struct {
	Edges []Edges `json:"edges"`
}
type Data struct {
	AssetEvents *AssetEvents `json:"assetEvents"`
}
