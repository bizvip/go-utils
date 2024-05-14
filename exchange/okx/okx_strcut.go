package okx

type Exchange struct {
	Currency string
	ShopName string
	Price    string
}

type DataResponse struct {
	Code         int    `json:"code"`
	Data         Data   `json:"data"`
	DetailMsg    string `json:"detailMsg"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Msg          string `json:"msg"`
	RequestID    string `json:"requestId"`
}

type Data struct {
	Buy       []interface{} `json:"buy"`
	Recommend Recommend     `json:"recommend"`
	Sell      []SellItem    `json:"sell"`
}

type Recommend struct {
	RecommendedAd interface{} `json:"recommendedAd"`
	UserGroup     int         `json:"userGroup"`
}

type SellItem struct {
	AlreadyTraded             bool        `json:"alreadyTraded"`
	AvailableAmount           string      `json:"availableAmount"`
	BaseCurrency              string      `json:"baseCurrency"`
	Black                     bool        `json:"black"`
	CancelledOrderQuantity    int         `json:"cancelledOrderQuantity"`
	CompletedOrderQuantity    int         `json:"completedOrderQuantity"`
	CompletedRate             string      `json:"completedRate"`
	CreatorType               string      `json:"creatorType"`
	GuideUpgradeKyc           bool        `json:"guideUpgradeKyc"`
	ID                        string      `json:"id"`
	Intention                 bool        `json:"intention"`
	IsInstitution             int         `json:"isInstitution"`
	MaxCompletedOrderQuantity int         `json:"maxCompletedOrderQuantity"`
	MaxUserCreatedDate        int         `json:"maxUserCreatedDate"`
	MerchantID                string      `json:"merchantId"`
	MinCompletedOrderQuantity int         `json:"minCompletedOrderQuantity"`
	MinCompletionRate         string      `json:"minCompletionRate"`
	MinKycLevel               int         `json:"minKycLevel"`
	MinSellOrders             int         `json:"minSellOrders"`
	Mine                      bool        `json:"mine"`
	NickName                  string      `json:"nickName"`
	PaymentMethods            []string    `json:"paymentMethods"`
	PaymentTimeoutMinutes     int         `json:"paymentTimeoutMinutes"`
	PosReviewPercentage       string      `json:"posReviewPercentage"`
	Price                     string      `json:"price"`
	PublicUserID              string      `json:"publicUserId"`
	QuoteCurrency             string      `json:"quoteCurrency"`
	QuoteMaxAmountPerOrder    string      `json:"quoteMaxAmountPerOrder"`
	QuoteMinAmountPerOrder    string      `json:"quoteMinAmountPerOrder"`
	QuoteScale                int         `json:"quoteScale"`
	QuoteSymbol               string      `json:"quoteSymbol"`
	ReceivingAds              bool        `json:"receivingAds"`
	SafetyLimit               bool        `json:"safetyLimit"`
	Side                      string      `json:"side"`
	UserActiveStatusVo        interface{} `json:"userActiveStatusVo"`
	UserType                  string      `json:"userType"`
	VerificationType          int         `json:"verificationType"`
	WhitelistedCountries      []string    `json:"whitelistedCountries"`
}
