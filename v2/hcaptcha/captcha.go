/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package hcaptcha

import (
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

type hCaptchaResponse struct {
	Success bool     `json:"success"`
	Error   []string `json:"error-codes"`
}

func VerifyHCaptcha(response string) (bool, error) {
	secret := "ES_fe0e15b17fc74f66906f5229a330199b"
	verifyURL := "https://api.hcaptcha.com/siteverify"

	data := url.Values{
		"secret":   {secret},
		"response": {response},
	}

	resp, err := http.PostForm(verifyURL, data)
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
