package osqueries

var (
	OrderQueryID = "OrdersQuery"
	orderQuery   = "query OrdersQuery(\n  $cursor: String\n  $count: Int = 10\n  $excludeMaker: IdentityInputType\n  $isExpired: Boolean\n  $isValid: Boolean\n  $maker: IdentityInputType\n  $makerArchetype: ArchetypeInputType\n  $makerAssetIsPayment: Boolean\n  $takerArchetype: ArchetypeInputType\n  $takerAssetCategories: [CollectionSlug!]\n  $takerAssetCollections: [CollectionSlug!]\n  $takerAssetIsOwnedBy: IdentityInputType\n  $takerAssetIsPayment: Boolean\n  $sortAscending: Boolean\n  $sortBy: OrderSortOption\n  $makerAssetBundle: BundleSlug\n  $takerAssetBundle: BundleSlug\n  $expandedMode: Boolean = false\n  $isBid: Boolean = false\n  $filterByOrderRules: Boolean = false\n) {\n  ...Orders_data_hzrfn\n}\n\nfragment AccountLink_data on AccountType {\n  address\n  config\n  isCompromised\n  user {\n    publicUsername\n    id\n  }\n  ...ProfileImage_data\n  ...wallet_accountKey\n  ...accounts_url\n}\n\nfragment AskPrice_data on OrderV2Type {\n  dutchAuctionFinalPrice\n  openedAt\n  priceFnEndedAt\n  makerAssetBundle {\n    assetQuantities(first: 30) {\n      edges {\n        node {\n          ...quantity_data\n          id\n        }\n      }\n    }\n    id\n  }\n  takerAssetBundle {\n    assetQuantities(first: 1) {\n      edges {\n        node {\n          ...AssetQuantity_data\n          id\n        }\n      }\n    }\n    id\n  }\n}\n\nfragment AssetCell_assetBundle on AssetBundleType {\n  assetQuantities(first: 2) {\n    edges {\n      node {\n        asset {\n          collection {\n            name\n            id\n          }\n          name\n          ...AssetMedia_asset\n          ...asset_url\n          id\n        }\n        relayId\n        id\n      }\n    }\n  }\n  name\n  slug\n}\n\nfragment AssetMedia_asset on AssetType {\n  animationUrl\n  backgroundColor\n  collection {\n    displayData {\n      cardDisplayStyle\n    }\n    id\n  }\n  isDelisted\n  displayImageUrl\n}\n\nfragment AssetQuantity_data on AssetQuantityType {\n  asset {\n    ...Price_data\n    id\n  }\n  quantity\n}\n\nfragment Orders_data_hzrfn on Query {\n  orders(after: $cursor, excludeMaker: $excludeMaker, first: $count, isExpired: $isExpired, isValid: $isValid, maker: $maker, makerArchetype: $makerArchetype, makerAssetIsPayment: $makerAssetIsPayment, takerArchetype: $takerArchetype, takerAssetCategories: $takerAssetCategories, takerAssetCollections: $takerAssetCollections, takerAssetIsOwnedBy: $takerAssetIsOwnedBy, takerAssetIsPayment: $takerAssetIsPayment, sortAscending: $sortAscending, sortBy: $sortBy, makerAssetBundle: $makerAssetBundle, takerAssetBundle: $takerAssetBundle, filterByOrderRules: $filterByOrderRules) {\n    edges {\n      node {\n        closedAt\n        isFulfillable\n        isValid\n        oldOrder\n        openedAt\n        orderType\n        maker {\n          address\n          ...AccountLink_data\n          ...wallet_accountKey\n          id\n        }\n        makerAsset: makerAssetBundle {\n          assetQuantities(first: 1) {\n            edges {\n              node {\n                asset {\n                  assetContract {\n                    address\n                    chain\n                    id\n                  }\n                  id\n                }\n                id\n              }\n            }\n          }\n          id\n        }\n        makerAssetBundle {\n          assetQuantities(first: 30) {\n            edges {\n              node {\n                ...AssetQuantity_data\n                ...quantity_data\n                id\n              }\n            }\n          }\n          id\n        }\n        relayId\n        side\n        taker {\n          address\n          ...AccountLink_data\n          ...wallet_accountKey\n          id\n        }\n        perUnitPrice {\n          eth\n        }\n        price {\n          usd\n        }\n        item @include(if: $isBid) {\n          __typename\n          ... on AssetType {\n            collection {\n              floorPrice\n              id\n            }\n          }\n          ... on AssetBundleType {\n            collection {\n              floorPrice\n              id\n            }\n          }\n          ... on Node {\n            __isNode: __typename\n            id\n          }\n        }\n        takerAssetBundle {\n          slug\n          ...bundle_url\n          assetQuantities(first: 1) {\n            edges {\n              node {\n                asset {\n                  ownedQuantity(identity: {})\n                  decimals\n                  symbol\n                  relayId\n                  assetContract {\n                    address\n                    id\n                  }\n                  ...asset_url\n                  id\n                }\n                quantity\n                ...AssetQuantity_data\n                ...quantity_data\n                id\n              }\n            }\n          }\n          id\n        }\n        ...AskPrice_data\n        ...orderLink_data\n        makerAssetBundleDisplay: makerAssetBundle @include(if: $expandedMode) {\n          ...AssetCell_assetBundle\n          id\n        }\n        takerAssetBundleDisplay: takerAssetBundle @include(if: $expandedMode) {\n          ...AssetCell_assetBundle\n          id\n        }\n        ...quantity_remaining\n        id\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment Price_data on AssetType {\n  decimals\n  imageUrl\n  symbol\n  usdSpotPrice\n  assetContract {\n    blockExplorerLink\n    chain\n    id\n  }\n}\n\nfragment ProfileImage_data on AccountType {\n  imageUrl\n  address\n}\n\nfragment accounts_url on AccountType {\n  address\n  user {\n    publicUsername\n    id\n  }\n}\n\nfragment asset_url on AssetType {\n  assetContract {\n    address\n    chain\n    id\n  }\n  tokenId\n}\n\nfragment bundle_url on AssetBundleType {\n  slug\n}\n\nfragment orderLink_data on OrderV2Type {\n  makerAssetBundle {\n    assetQuantities(first: 30) {\n      edges {\n        node {\n          asset {\n            externalLink\n            collection {\n              externalUrl\n              id\n            }\n            id\n          }\n          id\n        }\n      }\n    }\n    id\n  }\n}\n\nfragment quantity_data on AssetQuantityType {\n  asset {\n    decimals\n    id\n  }\n  quantity\n}\n\nfragment quantity_remaining on OrderV2Type {\n  makerAsset: makerAssetBundle {\n    assetQuantities(first: 1) {\n      edges {\n        node {\n          asset {\n            decimals\n            id\n          }\n          quantity\n          id\n        }\n      }\n    }\n    id\n  }\n  takerAsset: takerAssetBundle {\n    assetQuantities(first: 1) {\n      edges {\n        node {\n          asset {\n            decimals\n            id\n          }\n          quantity\n          id\n        }\n      }\n    }\n    id\n  }\n  remainingQuantity\n  side\n}\n\nfragment wallet_accountKey on AccountType {\n  address\n}\n"
)

type OrderAssetContract struct {
	Address                     string `json:"address"`
	AssetContractType           string `json:"asset_contract_type"`
	CreatedDate                 string `json:"created_date"`
	Name                        string `json:"name"`
	NftVersion                  string `json:"nft_version"`
	OpenseaVersion              string `json:"opensea_version"`
	Owner                       int    `json:"owner"`
	SchemaName                  string `json:"schema_name"`
	Symbol                      string `json:"symbol"`
	TotalSupply                 string `json:"total_supply"`
	Description                 string `json:"description"`
	ExternalLink                string `json:"external_link"`
	ImageURL                    string `json:"image_url"`
	DefaultToFiat               bool   `json:"default_to_fiat"`
	DevBuyerFeeBasisPoints      int    `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     int    `json:"dev_seller_fee_basis_points"`
	OnlyProxiedTransfers        bool   `json:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  int    `json:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints int    `json:"opensea_seller_fee_basis_points"`
	BuyerFeeBasisPoints         int    `json:"buyer_fee_basis_points"`
	SellerFeeBasisPoints        int    `json:"seller_fee_basis_points"`
	PayoutAddress               string `json:"payout_address"`
}

type OrderAssetCollection struct {
	BannerImageURL          string `json:"banner_image_url"`
	ChatURL                 string `json:"chat_url"`
	CreatedDate             string `json:"created_date"`
	DefaultToFiat           bool   `json:"default_to_fiat"`
	Description             string `json:"description"`
	DevBuyerFeeBasisPoints  string `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints string `json:"dev_seller_fee_basis_points"`
	DiscordURL              string `json:"discord_url"`
	DisplayData             struct {
		CardDisplayStyle string `json:"card_display_style"`
	} `json:"display_data"`
	ExternalURL                 string `json:"external_url"`
	Featured                    bool   `json:"featured"`
	FeaturedImageURL            string `json:"featured_image_url"`
	Hidden                      bool   `json:"hidden"`
	SafelistRequestStatus       string `json:"safelist_request_status"`
	ImageURL                    string `json:"image_url"`
	IsSubjectToWhitelist        bool   `json:"is_subject_to_whitelist"`
	LargeImageURL               string `json:"large_image_url"`
	MediumUsername              string `json:"medium_username"`
	Name                        string `json:"name"`
	OnlyProxiedTransfers        bool   `json:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  string `json:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints string `json:"opensea_seller_fee_basis_points"`
	PayoutAddress               string `json:"payout_address"`
	RequireEmail                bool   `json:"require_email"`
	ShortDescription            string `json:"short_description"`
	Slug                        string `json:"slug"`
	TelegramURL                 string `json:"telegram_url"`
	TwitterUsername             string `json:"twitter_username"`
	InstagramUsername           string `json:"instagram_username"`
	WikiURL                     string `json:"wiki_url"`
}

type OrderOwner struct {
	User struct {
		Username string `json:"username"`
	} `json:"user"`
	ProfileImgURL string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}

type OrderAsset struct {
	ID                   int                  `json:"id"`
	TokenID              string               `json:"token_id"`
	NumSales             int                  `json:"num_sales"`
	BackgroundColor      string               `json:"background_color"`
	ImageURL             string               `json:"image_url"`
	ImagePreviewURL      string               `json:"image_preview_url"`
	ImageThumbnailURL    string               `json:"image_thumbnail_url"`
	ImageOriginalURL     string               `json:"image_original_url"`
	AnimationURL         string               `json:"animation_url"`
	AnimationOriginalURL string               `json:"animation_original_url"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	ExternalLink         string               `json:"external_link"`
	AssetContract        OrderAssetContract   `json:"asset_contract"`
	Permalink            string               `json:"permalink"`
	Collection           OrderAssetCollection `json:"collection"`
	Decimals             int                  `json:"decimals"`
	TokenMetadata        string               `json:"token_metadata"`
	Owner                OrderUser            `json:"owner"`
}

type OrderUser struct {
	User struct {
		Username string `json:"username"`
	} `json:"user"`
	ProfileImgURL string `json:"profile_img_url"`
	Address       string `json:"address"`
	Config        string `json:"config"`
}

type Order struct {
	ID                int         `json:"id"`
	Asset             OrderAsset  `json:"asset"`
	AssetBundle       interface{} `json:"asset_bundle"`
	CreatedDate       string      `json:"created_date"`
	ClosingDate       string      `json:"closing_date"`
	ClosingExtendable bool        `json:"closing_extendable"`
	ExpirationTime    int         `json:"expiration_time"`
	ListingTime       int         `json:"listing_time"`
	OrderHash         string      `json:"order_hash"`
	Metadata          struct {
		Asset struct {
			ID      string `json:"id"`
			Address string `json:"address"`
		} `json:"asset"`
		Schema          string `json:"schema"`
		ReferrerAddress string `json:"referrerAddress"`
	} `json:"metadata"`
	Exchange             string    `json:"exchange"`
	Maker                OrderUser `json:"maker"`
	Taker                OrderUser `json:"taker"`
	CurrentPrice         string    `json:"current_price"`
	CurrentBounty        string    `json:"current_bounty"`
	BountyMultiple       string    `json:"bounty_multiple"`
	MakerRelayerFee      string    `json:"maker_relayer_fee"`
	TakerRelayerFee      string    `json:"taker_relayer_fee"`
	MakerProtocolFee     string    `json:"maker_protocol_fee"`
	TakerProtocolFee     string    `json:"taker_protocol_fee"`
	MakerReferrerFee     string    `json:"maker_referrer_fee"`
	FeeRecipient         OrderUser `json:"fee_recipient"`
	FeeMethod            int       `json:"fee_method"`
	Side                 int       `json:"side"`
	SaleKind             int       `json:"sale_kind"`
	Target               string    `json:"target"`
	HowToCall            int       `json:"how_to_call"`
	Calldata             string    `json:"calldata"`
	ReplacementPattern   string    `json:"replacement_pattern"`
	StaticTarget         string    `json:"static_target"`
	StaticExtradata      string    `json:"static_extradata"`
	PaymentToken         string    `json:"payment_token"`
	PaymentTokenContract struct {
		ID       int    `json:"id"`
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		ImageURL string `json:"image_url"`
		Name     string `json:"name"`
		Decimals int    `json:"decimals"`
		EthPrice string `json:"eth_price"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token_contract"`
	BasePrice       string      `json:"base_price"`
	Extra           string      `json:"extra"`
	Quantity        string      `json:"quantity"`
	Salt            string      `json:"salt"`
	V               int         `json:"v"`
	R               string      `json:"r"`
	S               string      `json:"s"`
	ApprovedOnChain bool        `json:"approved_on_chain"`
	Cancelled       bool        `json:"cancelled"`
	Finalized       bool        `json:"finalized"`
	MarkedInvalid   bool        `json:"marked_invalid"`
	PrefixedHash    string      `json:"prefixed_hash"`
	WaitingForBO    interface{} `json:"waitingForBestCounterOrder"`
}

type Archetype struct {
	AssetContractAddress string `json:"assetContractAddress"`
	TokenID              string `json:"tokenId"`
	// ETHEREUM
	Chain string `json:"chain"`
}

type OrderQueryVariables struct {
	Cursor       *string `json:"cursor"`
	Count        int     `json:"count"`
	ExcludeMaker *bool   `json:"excludeMaker"`
	IsExpired    bool    `json:"isExpired"`
	IsValid      bool    `json:"isValid"`
	// TODO create maker struct
	Maker                 *interface{}   `json:"maker"`
	MakerArchetype        *Archetype     `json:"makerArchetype"`
	MakerAssetIsPayment   *bool          `json:"makerAssetIsPayment"`
	TakerArchetype        *Archetype     `json:"takerArchetype"`
	TakerAssetCategories  *[]interface{} `json:"takerAssetCategories"`
	TakerAssetCollections *[]interface{} `json:"takerAssetCollections"`
	TakerAssetIsOwnedBy   *string        `json:"takerAssetIsOwnedBy"`
	TakerAssetIsPayment   *bool          `json:"takerAssetIsPayment"`
	SortAscending         bool           `json:"sortAscending"`
	// Buy => TAKER_ASSETS_USD_PRICE
	SortBy           string       `json:"sortBy"`
	MakerAssetBundle *interface{} `json:"makerAssetBundle"`
	TakerAssetBundle *interface{} `json:"takerAssetBundle"`
}

type OrdersQueryPayload struct {
	ID        string              `json:"id"`
	Query     string              `json:"query"`
	Variables OrderQueryVariables `json:"variables"`
}

// NewOrdersQueryPayload returns a query for sell orders for the given token.
// Returns a payload struct and the x-signed-query header.
func NewOrdersQueryPayload(contractAddress, tokenId string) (*OrdersQueryPayload, string) {
	return &OrdersQueryPayload{
		ID:    "EventHistoryPollQuery",
		Query: orderQuery,
		Variables: OrderQueryVariables{
			Count:     25,
			IsExpired: false,
			IsValid:   true,
			MakerArchetype: &Archetype{
				AssetContractAddress: contractAddress,
				TokenID:              tokenId,
				Chain:                "ETHEREUM",
			},
			SortAscending: true,
			SortBy:        "TAKER_ASSETS_USD_PRICE",
		},
	}, "1a98f73596b31bba66c339f9585a096bdfa32b85e53e853b9b5d3602628237c7"
}

type OrderQueryResponse struct {
	Data Data `json:"data"`
}
type User struct {
	PublicUsername string `json:"publicUsername"`
	ID             string `json:"id"`
}
type Maker struct {
	Address       string      `json:"address"`
	Config        interface{} `json:"config"`
	IsCompromised bool        `json:"isCompromised"`
	User          User        `json:"user"`
	ImageURL      string      `json:"imageUrl"`
	ID            string      `json:"id"`
}
type AssetContract struct {
	Address           string `json:"address"`
	Chain             string `json:"chain"`
	ID                string `json:"id"`
	BlockExplorerLink string `json:"blockExplorerLink"`
}

type Asset struct {
	UsdSpotPrice  float64       `json:"usdSpotPrice"`
	AssetContract AssetContract `json:"assetContract"`
	Collection    struct {
		Name                        string `json:"name"`
		Slug                        string `json:"slug"`
		Hidden                      bool   `json:"hidden"`
		DevSellerFeeBasisPoints     int    `json:"devSellerFeeBasisPoints"`
		OpenseaSellerFeeBasisPoints int    `json:"openseaSellerFeeBasisPoints"`
		IsMintable                  bool   `json:"isMintable"`
		IsSafelisted                bool   `json:"isSafelisted"`
		IsVerified                  bool   `json:"isVerified"`
		ID                          string `json:"id"`
		DisplayData                 struct {
			CardDisplayStyle string `json:"cardDisplayStyle"`
		} `json:"displayData"`
		ExternalURL string `json:"externalUrl"`
	} `json:"collection"`
	Decimals        int         `json:"decimals"`
	ImageURL        string      `json:"imageUrl"`
	Name            string      `json:"name"`
	Symbol          string      `json:"symbol"`
	RelayID         string      `json:"relayId"`
	AnimationURL    string      `json:"animationUrl"`
	BackgroundColor interface{} `json:"backgroundColor"`
	IsDelisted      bool        `json:"isDelisted"`
	DisplayImageURL string      `json:"displayImageUrl"`
	TokenID         string      `json:"tokenId"`
	ID              string      `json:"id"`
	ExternalLink    string      `json:"externalLink"`
	Quantity        string      `json:"quantity"`
	QuantityInEth   string      `json:"quantityInEth"`
}

type AssetQuantity struct {
	Asset    Asset  `json:"asset"`
	ID       string `json:"id"`
	Quantity string `json:"quantity"`
}
type AssetQuantitiesResult struct {
	Node AssetQuantity `json:"node"`
}
type AssetQuantities struct {
	Edges []AssetQuantitiesResult `json:"edges"`
}
type MakerAsset struct {
	AssetQuantities AssetQuantities `json:"assetQuantities"`
	ID              string          `json:"id"`
}
type MakerAssetBundle struct {
	AssetQuantities AssetQuantities `json:"assetQuantities"`
	ID              string          `json:"id"`
}
type PerUnitPrice struct {
	Eth string `json:"eth"`
}
type Price struct {
	Usd string `json:"usd"`
}
type TakerAssetBundle struct {
	Slug            string          `json:"slug"`
	AssetQuantities AssetQuantities `json:"assetQuantities"`
	ID              string          `json:"id"`
}
type TakerAsset struct {
	AssetQuantities AssetQuantities `json:"assetQuantities"`
	ID              string          `json:"id"`
}
type OrderQuery struct {
	ClosedAt               string           `json:"closedAt"`
	IsFulfillable          bool             `json:"isFulfillable"`
	IsValid                bool             `json:"isValid"`
	OldOrder               string           `json:"oldOrder"`
	OpenedAt               string           `json:"openedAt"`
	OrderType              string           `json:"orderType"`
	Maker                  Maker            `json:"maker"`
	MakerAsset             MakerAsset       `json:"makerAsset"`
	MakerAssetBundle       MakerAssetBundle `json:"makerAssetBundle"`
	RelayID                string           `json:"relayId"`
	Side                   string           `json:"side"`
	Taker                  interface{}      `json:"taker"`
	PerUnitPrice           PerUnitPrice     `json:"perUnitPrice"`
	Price                  Price            `json:"price"`
	TakerAssetBundle       TakerAssetBundle `json:"takerAssetBundle"`
	DutchAuctionFinalPrice interface{}      `json:"dutchAuctionFinalPrice"`
	PriceFnEndedAt         interface{}      `json:"priceFnEndedAt"`
	TakerAsset             TakerAsset       `json:"takerAsset"`
	RemainingQuantity      string           `json:"remainingQuantity"`
	ID                     string           `json:"id"`
	Typename               string           `json:"__typename"`
}
type OrderQueryResult struct {
	Node   OrderQuery `json:"node"`
	Cursor string     `json:"cursor"`
}
type PageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}
type Orders struct {
	Edges    []OrderQueryResult `json:"edges"`
	PageInfo PageInfo           `json:"pageInfo"`
}
type Data struct {
	Orders Orders `json:"orders"`
}
