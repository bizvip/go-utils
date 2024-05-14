package goutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"strings"
)

// "af","ak","am","ar","as","ay","az","be","bg","bho","bm","bn","bs","ca","ceb","ckb","zh","zh-CN","zh-TW",
// "co","cs","cy","da","de","doi","dv","ee","el","en","eo","es","et","eu","fa","fi","fr","fy","ga","gd",
// "gl","gn","gom","gu","ha","haw","he","hi","hmn","hr","ht","hu","hy","id","ig","ilo","is","it","iw","ja",
// "jv","jw","ka","kk","km","kn","ko","kri","ku","ky","la","lb","lg","ln","lo","lt","lus","lv","mai","mg",
// "mi","mk","ml","mn","mni-Mtei","mr","ms","mt","my","ne","nl","no","nso","ny","om","or","pa","pl","ps",
// "pt","qu","ro","ru","rw","sa","sd","si","sk","sl","sm","sn","so","sq","sr","st","su","sv","sw","ta",
// "te","tg","th","ti","tk","tl","tr","ts","tt","ug","uk","ur","uz","vi","xh","yi","yo","zu".

type ResultData struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

func GoogleTranslateToEn(text string, source string) (string, error) {
	text = url2.QueryEscape(text)
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
	payload := strings.NewReader(fmt.Sprintf("q=%s&target=en&source=%s", text, source))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "application/gzip")
	req.Header.Add("X-RapidAPI-Key", "2cf5688f73msh55fc8f71f9a01eap1964a2jsn12d1f9213dff")
	req.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return "", errors.New("请求google翻译返回结果出错 : " + string(body))
	}

	var result ResultData
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Data.Translations[0].TranslatedText, nil
}

func GoogleTranslateToCN(text string, source string) (string, error) {
	text = url2.QueryEscape(text)
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2"
	payload := strings.NewReader(fmt.Sprintf("q=%s&target=zh-CN&source=%s", text, source))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "application/gzip")
	req.Header.Add("X-RapidAPI-Key", "2cf5688f73msh55fc8f71f9a01eap1964a2jsn12d1f9213dff")
	req.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return "", errors.New("请求google翻译返回结果出错")
	}

	var result ResultData
	err := json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.Data.Translations[0].TranslatedText, nil
}

func GoogleDetectLang(text string) {
	url := "https://google-translate1.p.rapidapi.com/language/translate/v2/detect"
	payload := strings.NewReader("q=English%20is%20hard%2C%20but%20detectably%20so")

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "application/gzip")
	req.Header.Add("X-RapidAPI-Key", "2cf5688f73msh55fc8f71f9a01eap1964a2jsn12d1f9213dff")
	req.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
