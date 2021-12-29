package deth


// FlashbotRelayParams
type FlashbotRelayParams struct {
	// hex encoded block number for which this bundle is valid on
	BlockNumber string `json:"blockNumber,omitempty"`
}

// FlashbotEthParams defines parameters that all flashbot relay rpc calls to the ethereum chain use.
type FlashbotEthParams struct {
	FlashbotRelayParams
	// list of signed transactions to execute in an atomic bundle
	Txns []string `json:"txs,omitempty"`
}

// EthSendBundleParams defines a txn bundle send over flashbots rpc
type EthSendBundleParams struct {
	FlashbotEthParams
	// (Optional) minimum timestamp for which this bundle is valid, in seconds since the unix epoch
	MinTimestamp int `json:"minTimestamp,omitempty"`
	// (Optional) maximum timestamp for which this bundle is valid, in seconds since the unix epoch
	MaxTimestamp int `json:"maxTimestamp,omitempty"`
	// (Optional) list of tx hashes that are allowed to revert
	RevertingTxns []string `json:"revertingTxHashes,omitempty"`
}

// EthSendBundleParams defines a txn bundle call over flashbots rpc
type EthCallBundleParams struct {
	FlashbotEthParams
	// either a hex encoded number or a block tag for which state to base this simulation on. Can use "latest"
	StateBlockNumber string `json:"stateBlockNumber,omitempty"`
	// the timestamp to use for this bundle simulation, in seconds since the unix epoch
	TimeStamp int `json:"timestamp,omitempty"`
}

// GetUserStatsParams requests a quick summary of how a searcher is performing in the relay
type GetUserStatsParams struct {
	FlashbotRelayParams
}

// GetBundleStatsParams requests stats for a single bundle.
// Signing address must match the one who submitted the bundle.
type GetBundleStatsParams struct {
	FlashbotRelayParams
	// returned by the flashbots api when calling eth_sendBundle
	BundleHash string `json:"bundleHash,omitempty"`
}

// FlashbotRelayRequest is used to send request to flashbot relay.
type FlashbotRelayRequest struct {
	RPC    string         `json:"jsonrpc,omitempty"`
	ID     int            `json:"id,omitempty"`
	Method FlashbotMethod `json:"method,omitempty"`
	Params interface{}    `json:"params,omitempty"`
}

type FlashbotSendBundleResponse struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Result struct {
		BundleHash       string              `json:"bundleHash,omitempty"`
		BundleGasPrice   string              `json:"bundleGasPrice,omitempty"`
		CoinbaseDiff     string              `json:"coinbaseDiff,omitempty"`
		CoinbaseEthSent  string              `json:"ethSentToCoinbase,omitempty"`
		GasFee           string              `json:"gasFees,omitempty"`
		Results          []FlashbotTxnResult `json:"results,omitempty"`
		StateBlockNumber int                 `json:"stateBlockNumber,omitempty"`
		TotalGasUsed     int                 `json:"totalGasUsed,omitempty"`
	} `json:"result"`
}

type FlashbotTxnResult struct {
	CoinbaseDiff    string `json:"coinbaseDiff,omitempty"`
	CoinbaseEthSent string `json:"ethSentToCoinbase,omitempty"`
	From            string `json:"fromAddress,omitempty"`
	To              string `json:"toAddress,omitempty"`
	GasFee          string `json:"gasFees,omitempty"`
	GasPrice        string `json:"gasPrice,omitempty"`
	TxHash          string `json:"txHash,omitempty"`
	Value           string `json:"value,omitempty"`
}

type FlashbotBundleStats struct {
	RPC    string `json:"jsonrpc,omitempty"`
	ID     int    `json:"id,omitempty"`
	Result struct {
		IsSimulated    bool   `json:"isSimulated,omitempty"`
		IsSentToMiners bool   `json:"isSentToMiners,omitempty"`
		IsHighPriority bool   `json:"isHighPriority,omitempty"`
		SimulatedAt    string `json:"simulatedAt,omitempty"`
		SubmittedAt    string `json:"submittedAt,omitempty"`
		SentToMinersAt string `json:"sentToMinersAt,omitempty"`
	} `json:"method,omitempty"`
}

type FlashbotErrorResponse struct {
	Error struct {
		Message string `json:"message,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error,omitempty"`
}