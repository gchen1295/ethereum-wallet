package osqueries

var (
	AssetSearchQueryID = "AssetSearchQuery"
	AssetSearchQuery   = "query AssetSearchQuery(\n  $categories: [CollectionSlug!]\n  $chains: [ChainScalar!]\n  $collection: CollectionSlug\n  $collectionQuery: String\n  $collectionSortBy: CollectionSort\n  $collections: [CollectionSlug!]\n  $count: Int\n  $cursor: String\n  $identity: IdentityInputType\n  $includeHiddenCollections: Boolean\n  $numericTraits: [TraitRangeType!]\n  $paymentAssets: [PaymentAssetSymbol!]\n  $priceFilter: PriceFilterType\n  $query: String\n  $resultModel: SearchResultModel\n  $showContextMenu: Boolean = false\n  $shouldShowQuantity: Boolean = false\n  $sortAscending: Boolean\n  $sortBy: SearchSortBy\n  $stringTraits: [TraitInputType!]\n  $toggles: [SearchToggle!]\n  $creator: IdentityInputType\n  $assetOwner: IdentityInputType\n  $isPrivate: Boolean\n  $safelistRequestStatuses: [SafelistRequestStatus!]\n) {\n  query {\n    ...AssetSearch_data_2hBjZ1\n  }\n}\n\nfragment AssetCardAnnotations_assetBundle on AssetBundleType {\n  assetCount\n}\n\nfragment AssetCardAnnotations_asset_3Aax2O on AssetType {\n  assetContract {\n    chain\n    id\n  }\n  decimals\n  ownedQuantity(identity: $identity) @include(if: $shouldShowQuantity)\n  relayId\n  favoritesCount\n  isDelisted\n  isFavorite\n  isFrozen\n  hasUnlockableContent\n  ...AssetCardBuyNow_data\n  orderData {\n    bestAsk {\n      orderType\n      relayId\n      maker {\n        address\n      }\n    }\n  }\n  ...AssetContextMenu_data_3z4lq0 @include(if: $showContextMenu)\n}\n\nfragment AssetCardBuyNow_data on AssetType {\n  tokenId\n  relayId\n  assetContract {\n    address\n    chain\n    id\n  }\n  collection {\n    slug\n    id\n  }\n  orderData {\n    bestAsk {\n      relayId\n    }\n  }\n}\n\nfragment AssetCardContent_asset on AssetType {\n  relayId\n  name\n  ...AssetMedia_asset\n  assetContract {\n    address\n    chain\n    openseaVersion\n    id\n  }\n  tokenId\n  collection {\n    slug\n    id\n  }\n  isDelisted\n}\n\nfragment AssetCardContent_assetBundle on AssetBundleType {\n  assetQuantities(first: 18) {\n    edges {\n      node {\n        asset {\n          relayId\n          ...AssetMedia_asset\n          id\n        }\n        id\n      }\n    }\n  }\n}\n\nfragment AssetCardFooter_assetBundle on AssetBundleType {\n  ...AssetCardAnnotations_assetBundle\n  name\n  assetCount\n  assetQuantities(first: 18) {\n    edges {\n      node {\n        asset {\n          collection {\n            name\n            relayId\n            isVerified\n            ...collection_url\n            id\n          }\n          id\n        }\n        id\n      }\n    }\n  }\n  assetEventData {\n    lastSale {\n      unitPriceQuantity {\n        ...AssetQuantity_data\n        id\n      }\n    }\n  }\n  orderData {\n    bestBid {\n      orderType\n      paymentAssetQuantity {\n        ...AssetQuantity_data\n        id\n      }\n    }\n    bestAsk {\n      closedAt\n      orderType\n      dutchAuctionFinalPrice\n      openedAt\n      priceFnEndedAt\n      quantity\n      decimals\n      paymentAssetQuantity {\n        quantity\n        ...AssetQuantity_data\n        id\n      }\n    }\n  }\n}\n\nfragment AssetCardFooter_asset_3Aax2O on AssetType {\n  ...AssetCardAnnotations_asset_3Aax2O\n  name\n  tokenId\n  collection {\n    name\n    isVerified\n    ...collection_url\n    id\n  }\n  isDelisted\n  assetContract {\n    address\n    chain\n    openseaVersion\n    id\n  }\n  assetEventData {\n    lastSale {\n      unitPriceQuantity {\n        ...AssetQuantity_data\n        id\n      }\n    }\n  }\n  orderData {\n    bestBid {\n      orderType\n      paymentAssetQuantity {\n        ...AssetQuantity_data\n        id\n      }\n    }\n    bestAsk {\n      closedAt\n      orderType\n      dutchAuctionFinalPrice\n      openedAt\n      priceFnEndedAt\n      quantity\n      decimals\n      paymentAssetQuantity {\n        quantity\n        ...AssetQuantity_data\n        id\n      }\n    }\n  }\n}\n\nfragment AssetContextMenu_data_3z4lq0 on AssetType {\n  ...asset_edit_url\n  ...asset_url\n  ...itemEvents_data\n  isDelisted\n  isEditable {\n    value\n    reason\n  }\n  isListable\n  ownership(identity: {}) {\n    isPrivate\n    quantity\n  }\n  creator {\n    address\n    id\n  }\n  collection {\n    isAuthorizedEditor\n    id\n  }\n}\n\nfragment AssetMedia_asset on AssetType {\n  animationUrl\n  backgroundColor\n  collection {\n    displayData {\n      cardDisplayStyle\n    }\n    id\n  }\n  isDelisted\n  displayImageUrl\n}\n\nfragment AssetQuantity_data on AssetQuantityType {\n  asset {\n    ...Price_data\n    id\n  }\n  quantity\n}\n\nfragment AssetSearchFilter_data_3KTzFc on Query {\n  ...CollectionFilter_data_2qccfC\n  collection(collection: $collection) {\n    numericTraits {\n      key\n      value {\n        max\n        min\n      }\n      ...NumericTraitFilter_data\n    }\n    stringTraits {\n      key\n      ...StringTraitFilter_data\n    }\n    id\n  }\n  ...PaymentFilter_data_2YoIWt\n}\n\nfragment AssetSearchList_data_3Aax2O on SearchResultType {\n  asset {\n    assetContract {\n      address\n      chain\n      id\n    }\n    collection {\n      isVerified\n      relayId\n      id\n    }\n    relayId\n    tokenId\n    ...AssetSelectionItem_data\n    ...asset_url\n    id\n  }\n  assetBundle {\n    relayId\n    id\n  }\n  ...Asset_data_3Aax2O\n}\n\nfragment AssetSearch_data_2hBjZ1 on Query {\n  ...AssetSearchFilter_data_3KTzFc\n  ...SearchPills_data_2Kg4Sq\n  search(after: $cursor, chains: $chains, categories: $categories, collections: $collections, first: $count, identity: $identity, numericTraits: $numericTraits, paymentAssets: $paymentAssets, priceFilter: $priceFilter, querystring: $query, resultType: $resultModel, sortAscending: $sortAscending, sortBy: $sortBy, stringTraits: $stringTraits, toggles: $toggles, creator: $creator, isPrivate: $isPrivate, safelistRequestStatuses: $safelistRequestStatuses) {\n    edges {\n      node {\n        ...AssetSearchList_data_3Aax2O\n        __typename\n      }\n      cursor\n    }\n    totalCount\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment AssetSelectionItem_data on AssetType {\n  backgroundColor\n  collection {\n    displayData {\n      cardDisplayStyle\n    }\n    imageUrl\n    id\n  }\n  imageUrl\n  name\n  relayId\n}\n\nfragment Asset_data_3Aax2O on SearchResultType {\n  asset {\n    relayId\n    isDelisted\n    ...AssetCardContent_asset\n    ...AssetCardFooter_asset_3Aax2O\n    ...AssetMedia_asset\n    ...asset_url\n    ...itemEvents_data\n    orderData {\n      bestAsk {\n        paymentAssetQuantity {\n          quantityInEth\n          id\n        }\n      }\n    }\n    id\n  }\n  assetBundle {\n    relayId\n    ...bundle_url\n    ...AssetCardContent_assetBundle\n    ...AssetCardFooter_assetBundle\n    orderData {\n      bestAsk {\n        paymentAssetQuantity {\n          quantityInEth\n          id\n        }\n      }\n    }\n    id\n  }\n}\n\nfragment CollectionFilter_data_2qccfC on Query {\n  selectedCollections: collections(first: 25, collections: $collections, includeHidden: true) {\n    edges {\n      node {\n        assetCount\n        imageUrl\n        name\n        slug\n        isVerified\n        id\n      }\n    }\n  }\n  collections(assetOwner: $assetOwner, assetCreator: $creator, onlyPrivateAssets: $isPrivate, chains: $chains, first: 100, includeHidden: $includeHiddenCollections, parents: $categories, query: $collectionQuery, sortBy: $collectionSortBy) {\n    edges {\n      node {\n        assetCount\n        imageUrl\n        name\n        slug\n        isVerified\n        id\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment CollectionModalContent_data on CollectionType {\n  description\n  imageUrl\n  name\n  slug\n}\n\nfragment NumericTraitFilter_data on NumericTraitTypePair {\n  key\n  value {\n    max\n    min\n  }\n}\n\nfragment PaymentFilter_data_2YoIWt on Query {\n  paymentAssets(first: 10) {\n    edges {\n      node {\n        symbol\n        relayId\n        id\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n  PaymentFilter_collection: collection(collection: $collection) {\n    paymentAssets {\n      symbol\n      relayId\n      id\n    }\n    id\n  }\n}\n\nfragment Price_data on AssetType {\n  decimals\n  imageUrl\n  symbol\n  usdSpotPrice\n  assetContract {\n    blockExplorerLink\n    chain\n    id\n  }\n}\n\nfragment SearchPills_data_2Kg4Sq on Query {\n  selectedCollections: collections(first: 25, collections: $collections, includeHidden: true) {\n    edges {\n      node {\n        imageUrl\n        name\n        slug\n        ...CollectionModalContent_data\n        id\n      }\n    }\n  }\n}\n\nfragment StringTraitFilter_data on StringTraitType {\n  counts {\n    count\n    value\n  }\n  key\n}\n\nfragment asset_edit_url on AssetType {\n  assetContract {\n    address\n    chain\n    id\n  }\n  tokenId\n  collection {\n    slug\n    id\n  }\n}\n\nfragment asset_url on AssetType {\n  assetContract {\n    address\n    chain\n    id\n  }\n  tokenId\n}\n\nfragment bundle_url on AssetBundleType {\n  slug\n}\n\nfragment collection_url on CollectionType {\n  slug\n}\n\nfragment itemEvents_data on AssetType {\n  assetContract {\n    address\n    chain\n    id\n  }\n  tokenId\n}\n"
)

type AssetSearchVariables struct {
	Categories               interface{} `json:"categories"`
	Chains                   interface{} `json:"chains"`
	Collection               string      `json:"collection"`
	CollectionQuery          interface{} `json:"collectionQuery"`
	CollectionSortBy         string      `json:"collectionSortBy"`
	Collections              []string    `json:"collections"`
	Count                    int         `json:"count"`
	Cursor                   interface{} `json:"cursor"`
	Identity                 interface{} `json:"identity"`
	IncludeHiddenCollections interface{} `json:"includeHiddenCollections"`
	NumericTraits            interface{} `json:"numericTraits"`
	PaymentAssets            interface{} `json:"paymentAssets"`
	PriceFilter              interface{} `json:"priceFilter"`
	Query                    string      `json:"query"`
	ResultModel              string      `json:"resultModel"`
	ShowContextMenu          bool        `json:"showContextMenu"`
	ShouldShowQuantity       bool        `json:"shouldShowQuantity"`
	SortAscending            bool        `json:"sortAscending"`
	SortBy                   string      `json:"sortBy"`
	StringTraits             interface{} `json:"stringTraits"`
	Toggles                  []string    `json:"toggles"`
	Creator                  interface{} `json:"creator"`
	AssetOwner               interface{} `json:"assetOwner"`
	IsPrivate                interface{} `json:"isPrivate"`
	SafelistRequestStatuses  interface{} `json:"safelistRequestStatuses"`
}

type AssetSearchPayload struct {
	ID        string               `json:"id"`
	Query     string               `json:"query"`
	Variables AssetSearchVariables `json:"variables"`
}

func NewAssetSearchPayload(collection string) (*AssetSearchPayload, string) {
	return &AssetSearchPayload{
		ID:    AssetSearchQueryID,
		Query: AssetSearchQuery,
		Variables: AssetSearchVariables{
			Collection:         collection,
			CollectionSortBy:   "SEVEN_DAY_VOLUME",
			Collections:        []string{collection},
			Count:              50,
			ResultModel:        "ASSETS",
			ShowContextMenu:    true,
			ShouldShowQuantity: false,
			SortAscending:      true,
			SortBy:             "PRICE",
			Toggles:            []string{"BUY_NOW"}, //gives us only buy listings
		},
	}, "1c03875ecca199e4f680456647a41a6ba6d56dd0fc8a09811938f8077a2987d1"
}

type AssetSearchResponse struct {
	Data struct {
		Query struct {
			SelectedCollections struct {
				Edges []struct {
					Node struct {
						AssetCount  interface{} `json:"assetCount"`
						ImageURL    interface{} `json:"imageUrl"`
						Name        string      `json:"name"`
						Slug        string      `json:"slug"`
						IsVerified  bool        `json:"isVerified"`
						ID          string      `json:"id"`
						Description interface{} `json:"description"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"selectedCollections"`
			Collections struct {
				Edges []struct {
					Node struct {
						AssetCount interface{} `json:"assetCount"`
						ImageURL   interface{} `json:"imageUrl"`
						Name       string      `json:"name"`
						Slug       string      `json:"slug"`
						IsVerified bool        `json:"isVerified"`
						ID         string      `json:"id"`
						Typename   string      `json:"__typename"`
					} `json:"node"`
					Cursor string `json:"cursor"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"collections"`
			Collection struct {
				NumericTraits []interface{} `json:"numericTraits"`
				StringTraits  []interface{} `json:"stringTraits"`
				ID            string        `json:"id"`
			} `json:"collection"`
			PaymentAssets struct {
				Edges []struct {
					Node struct {
						Symbol   string `json:"symbol"`
						RelayID  string `json:"relayId"`
						ID       string `json:"id"`
						Typename string `json:"__typename"`
					} `json:"node"`
					Cursor string `json:"cursor"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"paymentAssets"`
			PaymentFilterCollection struct {
				PaymentAssets []struct {
					Symbol  string `json:"symbol"`
					RelayID string `json:"relayId"`
					ID      string `json:"id"`
				} `json:"paymentAssets"`
				ID string `json:"id"`
			} `json:"PaymentFilter_collection"`
			Search struct {
				Edges []struct {
					Node   AssetSearchResult `json:"node"`
					Cursor string            `json:"cursor"`
				} `json:"edges"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					EndCursor   string `json:"endCursor"`
					HasNextPage bool   `json:"hasNextPage"`
				} `json:"pageInfo"`
			} `json:"search"`
		} `json:"query"`
	} `json:"data"`
}

type AssetSearchResult struct {
	Asset struct {
		AssetContract struct {
			Address        string      `json:"address"`
			Chain          string      `json:"chain"`
			ID             string      `json:"id"`
			OpenseaVersion interface{} `json:"openseaVersion"`
		} `json:"assetContract"`
		Collection struct {
			IsVerified  bool   `json:"isVerified"`
			RelayID     string `json:"relayId"`
			ID          string `json:"id"`
			DisplayData struct {
				CardDisplayStyle string `json:"cardDisplayStyle"`
			} `json:"displayData"`
			ImageURL           interface{} `json:"imageUrl"`
			Slug               string      `json:"slug"`
			IsAuthorizedEditor bool        `json:"isAuthorizedEditor"`
			Name               string      `json:"name"`
		} `json:"collection"`
		RelayID              string      `json:"relayId"`
		TokenID              string      `json:"tokenId"`
		BackgroundColor      interface{} `json:"backgroundColor"`
		ImageURL             string      `json:"imageUrl"`
		Name                 string      `json:"name"`
		ID                   string      `json:"id"`
		IsDelisted           bool        `json:"isDelisted"`
		AnimationURL         interface{} `json:"animationUrl"`
		DisplayImageURL      string      `json:"displayImageUrl"`
		Decimals             int         `json:"decimals"`
		FavoritesCount       int         `json:"favoritesCount"`
		IsFavorite           bool        `json:"isFavorite"`
		IsFrozen             bool        `json:"isFrozen"`
		HasUnlockableContent bool        `json:"hasUnlockableContent"`
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
						Decimals      int     `json:"decimals"`
						ImageURL      string  `json:"imageUrl"`
						Symbol        string  `json:"symbol"`
						UsdSpotPrice  float64 `json:"usdSpotPrice"`
						AssetContract struct {
							BlockExplorerLink string `json:"blockExplorerLink"`
							Chain             string `json:"chain"`
							ID                string `json:"id"`
						} `json:"assetContract"`
						ID string `json:"id"`
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
			LastSale interface{} `json:"lastSale"`
		} `json:"assetEventData"`
	} `json:"asset"`
	AssetBundle interface{} `json:"assetBundle"`
	Typename    string      `json:"__typename"`
}
