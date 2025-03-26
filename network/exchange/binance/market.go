package binance

import (
	"fmt"

	"github.com/goccy/go-json"
)

type MarketService struct{}

func NewMarketService() *MarketService {
	return &MarketService{}
}

// GetPing 请求 Binance API 进行连接测试
func (m *MarketService) GetPing() interface{} {
	query := "/api/v3/ping"
	return GetApi(query)
}

// GetServerTime 请求 Binance API 获取服务器时间
func (m *MarketService) GetServerTime() interface{} {
	query := "/api/v3/time"
	return GetApi(query)
}

// GetExchangeInfo 请求 Binance API 获取交易所信息
func (m *MarketService) GetExchangeInfo() interface{} {
	query := "/api/v3/exchangeInfo"
	return GetApi(query)
}

// GetDepth 请求 Binance API 获取订单簿
func (m *MarketService) GetDepth(symbol string) interface{} {
	query := fmt.Sprintf("/api/v3/depth?symbol=%s", symbol)
	return GetApi(query)
}

// GetTrades 请求 Binance API 获取最近的交易列表
func (m *MarketService) GetTrades(symbol string, limit int) interface{} {
	query := fmt.Sprintf("/api/v3/trades?symbol=%s&limit=%d", symbol, limit)
	return GetApi(query)
}

// GetHistoricalTrades 请求 Binance API 获取旧的交易记录
func (m *MarketService) GetHistoricalTrades(symbol string, limit int) interface{} {
	query := fmt.Sprintf("/api/v3/historicalTrades?symbol=%s&limit=%d", symbol, limit)
	return GetApi(query)
}

// GetAggTrades 请求 Binance API 获取压缩/聚合交易记录
func (m *MarketService) GetAggTrades(symbol string, limit int) interface{} {
	query := fmt.Sprintf("/api/v3/aggTrades?symbol=%s&limit=%d", symbol, limit)
	return GetApi(query)
}

// GetKlines 请求 Binance API 获取K线/蜡烛图数据
func (m *MarketService) GetKlines(symbol, interval string) interface{} {
	query := fmt.Sprintf("/api/v3/klines?symbol=%s&interval=%s", symbol, interval)
	return GetApi(query)
}

// GetUIKlines 请求 Binance API 获取UIK线数据
func (m *MarketService) GetUIKlines(symbol, interval string) interface{} {
	query := fmt.Sprintf("/api/v3/uiKlines?symbol=%s&interval=%s", symbol, interval)
	return GetApi(query)
}

// GetAvgPrice 请求 Binance API 获取当前平均价格
func (m *MarketService) GetAvgPrice(symbol string) interface{} {
	query := fmt.Sprintf("/api/v3/avgPrice?symbol=%s", symbol)
	return GetApi(query)
}

// GetTicker24Hr 请求 Binance API 获取24小时价格变动情况
// symbols 是一个包含多个交易对的字符串数组，例如 ["BTCUSDT", "BNBBTC"]
// 返回结果是一个 interface{} 类型，包含了各个交易对的24小时价格变动数据
// 返回的 JSON 数据结构如下:
//
//	[{
//	  "symbol": "BNBBTC",                 // 交易对
//	  "priceChange": "-0.00010200",       // 价格变动
//	  "priceChangePercent": "-1.152",     // 价格变动百分比
//	  "weightedAvgPrice": "0.00877364",   // 加权平均价格
//	  "prevClosePrice": "0.00885700",     // 前一次收盘价
//	  "lastPrice": "0.00875500",          // 最新价格
//	  "lastQty": "3.92000000",            // 最新成交量
//	  "bidPrice": "0.00875500",           // 买入价格
//	  "bidQty": "20.61500000",            // 买入数量
//	  "askPrice": "0.00875600",           // 卖出价格
//	  "askQty": "10.44800000",            // 卖出数量
//	  "openPrice": "0.00885700",          // 开盘价
//	  "highPrice": "0.00889800",          // 最高价
//	  "lowPrice": "0.00868500",           // 最低价
//	  "volume": "48607.13000000",         // 成交量
//	  "quoteVolume": "426.46130696",      // 成交额
//	  "openTime": 1721647533545,          // 开盘时间（Unix毫秒时间戳）
//	  "closeTime": 1721733933545,         // 收盘时间（Unix毫秒时间戳）
//	  "firstId": 251266985,               // 首次成交ID
//	  "lastId": 251440381,                // 最后成交ID
//	  "count": 173397                     // 成交笔数
//	},]
func (m *MarketService) GetTicker24Hr(symbols []string, dataType string) interface{} {
	symbolsParam, _ := json.Marshal(symbols)
	query := fmt.Sprintf("/api/v3/ticker/24hr?symbols=%s&type=%s", symbolsParam, dataType)
	return GetApi(query)
}

// GetTickerTradingDay 请求 Binance API 获取交易日价格变动情况
func (m *MarketService) GetTickerTradingDay(symbols []string) interface{} {
	symbolsParam, _ := json.Marshal(symbols)
	query := fmt.Sprintf("/api/v3/ticker/tradingDay?symbols=%s", symbolsParam)
	return GetApi(query)
}

// GetTickerPrice 请求 Binance API 获取Symbol价格
func (m *MarketService) GetTickerPrice(symbols []string) interface{} {
	symbolsParam, _ := json.Marshal(symbols)
	query := fmt.Sprintf("/api/v3/ticker/price?symbols=%s", symbolsParam)
	return GetApi(query)
}

// GetTickerBookTicker 请求 Binance API 获取Symbol Order Book Ticker
func (m *MarketService) GetTickerBookTicker(symbols []string) interface{} {
	symbolsParam, _ := json.Marshal(symbols)
	query := fmt.Sprintf("/api/v3/ticker/bookTicker?symbols=%s", symbolsParam)
	return GetApi(query)
}

// GetTicker 请求 Binance API 获取滚动窗口价格变化统计数据
// symbols: 交易对列表，例如 ["BTCUSDT", "BNBUSDT"]
// windowSize: 时间窗口大小，支持的值有 "1m", "2m", ... "59m"（分钟）, "1h", "2h", ... "23h"（小时）, "1d", ... "7d"（天）
// type: 支持的值有 "FULL" 或 "MINI"。如果未提供，默认为 "FULL"
// 返回字段：
//   - symbol: 交易对
//   - priceChange: 价格变化
//   - priceChangePercent: 价格变化百分比
//   - weightedAvgPrice: 加权平均价格
//   - prevClosePrice: 前收盘价
//   - lastPrice: 最新价格
//   - lastQty: 最新成交量
//   - bidPrice: 买方出价
//   - askPrice: 卖方出价
//   - openPrice: 开盘价
//   - highPrice: 最高价
//   - lowPrice: 最低价
//   - volume: 成交量
//   - quoteVolume: 成交额
//   - openTime: 开始时间
//   - closeTime: 结束时间
//   - firstId: 首次成交ID
//   - lastId: 最后成交ID
//   - count: 成交次数
func (m *MarketService) GetTicker(symbols []string, windowSize string, dataType string) interface{} {
	symbolsParam, _ := json.Marshal(symbols)
	query := fmt.Sprintf("/api/v3/ticker?symbols=%s&windowSize=%s&type=%s", symbolsParam, windowSize, dataType)
	return GetApi(query)
}
