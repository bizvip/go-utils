/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package okx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type OKX struct {
	Client *http.Client
}
type PayMethod string

const (
	PayMethodWxPay  PayMethod = "wxPay"
	PayMethodAlipay PayMethod = "alipay"
	PayMethodBank   PayMethod = "bank"
	PayMethodAll    PayMethod = "all"
)

func NewOkxExchangeService() *OKX {
	return &OKX{Client: &http.Client{}}
}

// GetTop10Exchanges 获取前10交易商家的 C2C 汇率 可指定不同的货币和支付方式
func (o *OKX) GetTop10Exchanges(baseCurrency, quoteCurrency string, okxPayMethod PayMethod) ([]*Exchange, error) {
	url := fmt.Sprintf(
		"https://www.okx.com/v3/c2c/tradingOrders/books?quoteCurrency=%s&baseCurrency=%s&side=sell&paymentMethod=%s&userType=all&receivingAds=false&limit=10&t=%d",
		quoteCurrency, baseCurrency, string(okxPayMethod), time.Now().UnixMilli())
	result, err := o.doRequest(url)
	if err != nil {
		return nil, err
	}

	var data DataResponse
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		return nil, err
	}

	var exchanges []*Exchange
	for _, v := range data.Data.Sell {
		exchanges = append(exchanges, &Exchange{
			Currency: v.BaseCurrency,
			ShopName: v.NickName,
			Price:    v.Price,
		})
	}

	return exchanges, nil
}

// GetUsdtCnyExchangeList 获取usdt到cny的前10个实时汇率结果列表
func (o *OKX) GetUsdtCnyExchangeList(okxPayMethod PayMethod) ([]*Exchange, error) {
	return o.GetTop10Exchanges("USDT", "CNY", okxPayMethod)
}

// GetUsdtCnyRateOnly 获取实时人民币到USDT的汇率
func (o *OKX) GetUsdtCnyRateOnly(okxPayMethod PayMethod) decimal.Decimal {
	usdtRates, err := o.GetUsdtCnyExchangeList(okxPayMethod)
	if err != nil {
		return decimal.Zero
	}
	var total decimal.Decimal
	for _, usdtRate := range usdtRates {
		price, err := decimal.NewFromString(usdtRate.Price)
		if err != nil {
			continue
		}
		total = total.Add(price)
	}
	return total.DivRound(decimal.NewFromInt(int64(len(usdtRates))), 2)
}

func (o *OKX) doRequest(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("App-Type", "web")

	res, err := o.Client.Do(req)
	if err != nil {
		log.Printf("Error performing request: %v", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("API response error: %s", res.Status)
		return "", fmt.Errorf("API responded with status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	return string(body), nil
}
