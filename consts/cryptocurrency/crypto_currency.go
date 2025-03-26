package cryptocurrency

const (
	// 常见加密货币

	USDT  = "USDT"  // 泰达币（稳定币，锚定美元）
	BTC   = "BTC"   // 比特币
	ETH   = "ETH"   // 以太坊
	XRP   = "XRP"   // 瑞波币
	BCH   = "BCH"   // 比特币现金
	LTC   = "LTC"   // 莱特币
	ADA   = "ADA"   // 艾达币（卡尔达诺）
	DOT   = "DOT"   // 波卡币（波尔卡多特）
	BNB   = "BNB"   // 币安币
	SOL   = "SOL"   // Solana
	DOGE  = "DOGE"  // 狗狗币
	UNI   = "UNI"   // Uniswap
	ATOM  = "ATOM"  // 宇宙币（Cosmos）
	ICP   = "ICP"   // 互联网计算机
	XMR   = "XMR"   // 门罗币
	ETC   = "ETC"   // 以太坊经典
	XLM   = "XLM"   // 恒星币（Stellar）
	TRX   = "TRX"   // 波场币（Tron）
	NEO   = "NEO"   // 小蚁币（NEO）
	EOS   = "EOS"   // 柚子币
	FIL   = "FIL"   // Filecoin
	MATIC = "MATIC" // Polygon

	// 区块链网络名称

	BtcNetwork     = "BTC"       // 比特币主网
	EthNetwork     = "ETH"       // 以太坊主网
	BscNetwork     = "BSC"       // 币安智能链
	TronNetwork    = "TRON"      // 波场网络
	PolygonNetwork = "Polygon"   // Polygon(Matic)网络
	SolNetwork     = "Solana"    // Solana网络
	AvaxNetwork    = "Avalanche" // Avalanche网络
	FtmNetwork     = "Fantom"    // Fantom网络
	ArbNetwork     = "Arbitrum"  // Arbitrum网络
	OpNetwork      = "Optimism"  // Optimism网络
	HecoNetwork    = "HECO"      // 火币生态链
	OkcNetwork     = "OKC"       // OKX链

	// 代币协议标准

	ERC20   = "ERC20"   // 以太坊ERC20代币标准
	ERC721  = "ERC721"  // 以太坊ERC721 NFT标准
	ERC1155 = "ERC1155" // 以太坊ERC1155多重代币标准
	BEP20   = "BEP20"   // 币安智能链代币标准
	TRC20   = "TRC20"   // 波场TRC20代币标准
	TRC10   = "TRC10"   // 波场TRC10代币标准
	RC20    = "RC20"    // 通用RC20代币标准
	SPL     = "SPL"     // Solana程序库代币标准

	// 常见代币在不同链上的类型

	UsdtErc20 = "USDT-ERC20" // 以太坊USDT
	UsdtTrc20 = "USDT-TRC20" // 波场USDT
	UsdtBep20 = "USDT-BEP20" // BSC USDT
	UsdcErc20 = "USDC-ERC20" // 以太坊USDC
	UsdcBep20 = "USDC-BEP20" // BSC USDC
	// UsdcTrc20 = "USDC-TRC20"
)

// 区块链浏览器URL

var ExplorerURLs = map[string]string{
	BtcNetwork:     "https://www.blockchain.com/explorer", // BTC浏览器
	EthNetwork:     "https://etherscan.io",                // ETH浏览器
	BscNetwork:     "https://bscscan.com",                 // BSC浏览器
	TronNetwork:    "https://tronscan.org",                // 波场浏览器
	PolygonNetwork: "https://polygonscan.com",             // Polygon浏览器
	SolNetwork:     "https://explorer.solana.com",         // Solana浏览器
	AvaxNetwork:    "https://snowtrace.io",                // Avalanche浏览器
	FtmNetwork:     "https://ftmscan.com",                 // Fantom浏览器
	ArbNetwork:     "https://arbiscan.io",                 // Arbitrum浏览器
	OpNetwork:      "https://optimistic.etherscan.io",     // Optimism浏览器
	HecoNetwork:    "https://hecoinfo.com",                // HECO浏览器
	OkcNetwork:     "https://www.oklink.com/okc",          // OKC浏览器
}
