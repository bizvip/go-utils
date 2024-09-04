package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"strings"
)

var (
	raidApiKey  string = "2cf5688f73msh55fc8f71f9a01eap1964a2jsn12d1f9213dff"
	raidApiHost string = "google-translate1.p.rapidapi.com"
)

// TranslationService 包含用于Google翻译API的API密钥和主机信息
type TranslationService struct {
	apiKey  string
	apiHost string
}

// NewTranslationService 初始化一个 TranslationService 实例，并设置 API 密钥和主机
func NewTranslationService(apiKey, apiHost string) *TranslationService {
	return &TranslationService{
		apiKey:  apiKey,
		apiHost: apiHost,
	}
}

// ResultData 用于解析翻译API的JSON响应
type ResultData struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

// sendRequest 发送HTTP请求并返回响应体
func (g *TranslationService) sendRequest(url, payload string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "application/gzip")
	req.Header.Add("X-RapidAPI-Key", g.apiKey)
	req.Header.Add("X-RapidAPI-Host", g.apiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("请求google翻译返回结果出错: " + string(body))
	}

	return body, nil
}

// GoogleTranslateToEn 将文本从指定语言翻译为英文
func (g *TranslationService) GoogleTranslateToEn(text string, source string) (string, error) {
	text = url2.QueryEscape(text)
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
	payload := fmt.Sprintf("q=%s&target=en&source=%s", text, source)

	body, err := g.sendRequest(url, payload)
	if err != nil {
		return "", err
	}

	var result ResultData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Data.Translations[0].TranslatedText, nil
}

// GoogleTranslateToCN 将文本从指定语言翻译为中文
func (g *TranslationService) GoogleTranslateToCN(text string, source string) (string, error) {
	text = url2.QueryEscape(text)
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
	payload := fmt.Sprintf("q=%s&target=zh-CN&source=%s", text, source)

	body, err := g.sendRequest(url, payload)
	if err != nil {
		return "", err
	}

	var result ResultData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Data.Translations[0].TranslatedText, nil
}

// GoogleDetectLang 检测文本的语言
func (g *TranslationService) GoogleDetectLang(text string) (string, error) {
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2/detect"
	payload := fmt.Sprintf("q=%s", url2.QueryEscape(text))

	body, err := g.sendRequest(url, payload)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
