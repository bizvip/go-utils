package hcaptcha

import (
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

var hSecret = "ES_fe0e15b17fc74f66906f5229a330199b"
var hVerifyURL = "https://api.hcaptcha.com/siteverify"

// Verifier 验证的配置
type Verifier struct {
	secret    string
	verifyURL string
}

// NewHCaptchaVerifier 初始化一个 HCaptchaVerifier 实例
func NewHCaptchaVerifier(secret, verifyURL string) *Verifier {
	return &Verifier{
		secret:    secret,
		verifyURL: verifyURL,
	}
}

type hCaptchaResponse struct {
	Success bool     `json:"success"`
	Error   []string `json:"error-codes"`
}

// Verify 验证 hCaptcha 的响应
func (v *Verifier) Verify(response string) (bool, error) {
	data := url.Values{
		"secret":   {v.secret},
		"response": {response},
	}

	resp, err := http.PostForm(v.verifyURL, data)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var hResp hCaptchaResponse
	if err = json.NewDecoder(resp.Body).Decode(&hResp); err != nil {
		return false, err
	}

	return hResp.Success, nil
}
