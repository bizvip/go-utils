package okx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bizvip/go-utils/logs"
)

type OKX struct {
	Client *http.Client
}

func NewOkxService() *OKX {
	return &OKX{Client: &http.Client{}}
}

func (o *OKX) GetTop10Exchanges(baseCurrency string, quoteCurrency string) ([]Exchange, error) {
	url := fmt.Sprintf("https://www.okx.com/v3/c2c/tradingOrders/books?quoteCurrency=%s&baseCurrency=%s&side=sell&paymentMethod=all&userType=all&receivingAds=false&limit=10&t=%d", quoteCurrency, baseCurrency, time.Now().UnixMilli())
	result, err := o.doRequest(url)
	if err != nil {
		logs.Logger().Error(err)
		return nil, err
	}

	var data DataResponse
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		logs.Logger().Error(err)
		return nil, err
	}

	var exchanges []Exchange
	for _, v := range data.Data.Sell {
		exchange := Exchange{
			Currency: v.BaseCurrency,
			ShopName: v.NickName,
			Price:    v.Price,
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

func (o *OKX) doRequest(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	// 设置User-Agent头部
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
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
