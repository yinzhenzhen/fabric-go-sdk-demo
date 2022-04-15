package cli

type BlockInfo struct {
	ChannelId     string         `json:"channel_id"`
	BlockNum      uint64         `json:"block_num"`
	PrevBlockHash string         `json:"prev_block_hash"`
	BlockHash     string         `json:"block_hash"`
	DataHash      string         `json:"data_hash"`
	TxCount       int32          `json:"tx_count"`
	BlockTime     string         `json:"block_time"`
	BlockSize     uint64         `json:"block_size"`
	Trans         []*Transaction `json:"trans"`
}

type Transaction struct {
	TxId           string `json:"tx_id"`
	ChainCodeId    string `json:"chain_code_id"`
	ChainCodeName  string `json:"chain_code_name"`
	Status         int32  `json:"status"`
	ValidationCode string `json:"validation_code"`
	Type           string `json:"type"`
	CreaterMSPId   string `json:"creater_msp_id"`
	TransSize      uint64 `json:"trans_size"`
	TxTime         string `json:"tx_time"`
	// Base64编码
	Input string `json:"input"`
	// Base64编码
	Output string `json:"output"`
}
