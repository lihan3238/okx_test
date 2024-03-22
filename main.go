package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type IndexTicker struct {
	InstId  string `json:"instId"`
	IdxPx   string `json:"idxPx"`
	High24h string `json:"high24h"`
	Low24h  string `json:"low24h"`
	Open24h string `json:"open24h"`
	SodUtc0 string `json:"sodUtc0"`
	SodUtc8 string `json:"sodUtc8"`
	Ts      string `json:"ts"`
}

type ExchangeRate struct {
	UsdCny string `json:"usdCny"`
}

type IndexCandle struct {
	Ts      string `json:"ts"`
	Open    string `json:"o"`
	High    string `json:"h"`
	Low     string `json:"l"`
	Close   string `json:"c"`
	Confirm string `json:"confirm"`
}

func main() {
	r := gin.Default()

	// 设置静态资源和HTML模板文件夹的相对路径
	r.Static("/static", "static")
	r.Static("/assets", "assets")
	r.Static("/imgs", "imgs")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/task", func(c *gin.Context) {
		// 获取前端发送的数据
		quoteCcy := c.PostForm("quoteCcy")
		instId := c.PostForm("instId")
		after := c.PostForm("after")
		before := c.PostForm("before")
		bar := c.PostForm("bar")
		limit := c.PostForm("limit")

		queryType := c.PostForm("queryType")
		// 执行查询，获取结果
		results, err := executeQuery(quoteCcy, instId, after, before, bar, limit, queryType)
		if err != nil {
			log.Println("Error executing query:", err)
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		// 将查询结果以合适的格式返回给前端
		c.JSON(http.StatusOK, results)
	})

	r.Run(":8080")
}

func executeQuery(quoteCcy string, instId string, after string, before string, bar string, limit string, queryType string) (interface{}, error) {
	switch queryType {
	case "index":
		return executeIndexQuery(quoteCcy, instId)
	case "exchangeRate":
		return executeExchangeRateQuery()
	case "indexCandles":
		return executeIndexCandlesQuery(instId, after, before, bar, limit)
	default:
		return nil, errors.New("unsupported query type")
	}
}

func executeIndexQuery(quoteCcy string, instId string) ([]IndexTicker, error) {
	// 构建API请求URL
	apiURL := "https://okx.com/api/v5/market/index-tickers"
	params := url.Values{}
	params.Set("quoteCcy", quoteCcy)
	params.Set("instId", instId)

	fullURL := apiURL + "?" + params.Encode()
	//log.Printf("API Response: %s\n", string(fullURL))

	// 创建带有代理的HTTP客户端
	client := &http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(&url.URL{Host: "localhost:7890", Scheme: "http"}), // 设置Clash代理的地址和端口
		//},
	}

	// 发送GET请求获取数据
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//log.Printf("API Response: %s\n", string(body))

	// 解析JSON响应
	var response struct {
		Code string        `json:"code"`
		Msg  string        `json:"msg"`
		Data []IndexTicker `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	// 返回数据部分
	return response.Data, nil
}

func executeExchangeRateQuery() ([]ExchangeRate, error) {
	// 构建API请求URL
	apiURL := "https://okx.com/api/v5/market/exchange-rate"

	// 创建带有代理的HTTP客户端
	client := &http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(&url.URL{Host: "localhost:7890", Scheme: "http"}), // 设置Clash代理的地址和端口
		//},
	}

	// 发送GET请求获取数据
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var response struct {
		Code string         `json:"code"`
		Msg  string         `json:"msg"`
		Data []ExchangeRate `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// 返回数据部分
	return response.Data, nil
}

func executeIndexCandlesQuery(instId string, after string, before string, bar string, limit string) ([]IndexCandle, error) {
	// 构建API请求URL
	apiURL := "https://okx.com/api/v5/market/index-candles"
	params := url.Values{}
	params.Set("instId", instId)

	// 可选参数
	if after != "" {
		params.Set("after", after)
	}
	if before != "" {
		params.Set("before", before)
	}
	if bar != "" {
		params.Set("bar", bar)
	}
	if limit != "" {
		params.Set("limit", limit)
	}

	fullURL := apiURL + "?" + params.Encode()
	//log.Printf("API Response: %s\n", string(fullURL))

	// 创建带有代理的HTTP客户端
	client := &http.Client{
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(&url.URL{Host: "localhost:7890", Scheme: "http"}), // 设置Clash代理的地址和端口
		//},
	}

	// 发送GET请求获取数据
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("API Response: %s\n", string(body))

	// 解析JSON响应
	var response struct {
		Code string     `json:"code"`
		Msg  string     `json:"msg"`
		Data [][]string `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var candles []IndexCandle
	for _, candle := range response.Data {
		candles = append(candles, IndexCandle{
			Ts:      candle[0],
			Open:    candle[1],
			High:    candle[2],
			Low:     candle[3],
			Close:   candle[4],
			Confirm: candle[5],
		})
	}
	return candles, nil
}
