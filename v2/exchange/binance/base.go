/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.
 * Author ORCID: https://orcid.org/0009-0003-8150-367X
 ******************************************************************************/

package binance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

// 定义API的Base URL
const (
	baseURL1 = "https://api.binance.com"
	baseURL2 = "https://api-gcp.binance.com"
	baseURL3 = "https://api1.binance.com"
	baseURL4 = "https://api2.binance.com"
	baseURL5 = "https://api3.binance.com"
	baseURL6 = "https://api4.binance.com"
	baseURL7 = "https://data-api.binance.vision"
)

// API URLs 列表
var apiURLs = []string{
	baseURL1,
	baseURL2,
	baseURL3,
	baseURL4,
	baseURL5,
	baseURL6,
	baseURL7,
}

// 封装了带超时功能的HTTP客户端
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// 创建请求的函数
func createGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// 执行请求并处理失败的函数
func doRequestWithFallback(req *http.Request) (*http.Response, error) {
	var lastErr error
	for _, baseURL := range apiURLs {
		// 更新请求URL
		req.URL, _ = req.URL.Parse(baseURL + req.URL.RequestURI())
		// 执行请求
		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			fmt.Printf("Error accessing %s: %v\n", baseURL, err)
			continue
		}
		// 检查返回状态
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, baseURL)
			fmt.Printf("Error from %s: %v\n", baseURL, lastErr)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("failed to get data from all API endpoints: %v", lastErr)
}

// 解析返回数据的函数
func parseResponse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetApi Get请求流程
func GetApi(query string) interface{} {
	// 创建请求
	req, err := createGetRequest(query)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil
	}

	// 执行请求并处理失败
	resp, err := doRequestWithFallback(req)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)
		return nil
	}

	// 解析响应
	result, err := parseResponse(resp)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil
	}

	return result
}
