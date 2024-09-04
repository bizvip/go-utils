/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.
 * Author ORCID: https://orcid.org/0009-0003-8150-367X
 ******************************************************************************/

package opencc

import (
	"log"
	"sync"

	"github.com/longbridgeapp/opencc"
)

var (
	t2s  *opencc.OpenCC
	s2t  *opencc.OpenCC
	once sync.Once
)

func init() {
	once.Do(
		func() {
			var err error
			t2s, err = opencc.New("t2s")
			if err != nil {
				log.Fatalf("Failed to initialize OpenCC for t2s: %v", err)
			}
			s2t, err = opencc.New("s2t")
			if err != nil {
				log.Fatalf("Failed to initialize OpenCC for s2t: %v", err)
			}
		},
	)
}

// TwToChs 将繁体中文转换为简体中文
func TwToChs(text string) (string, error) {
	return t2s.Convert(text)
}

// ChsToTW 将简体中文转换为繁体中文
func ChsToTW(text string) (string, error) {
	return s2t.Convert(text)
}
