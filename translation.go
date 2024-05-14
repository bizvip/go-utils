/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package goutils

import (
	"sync"

	"github.com/longbridgeapp/opencc"
)

type LangTranslatorUtils struct {
	T2s  *opencc.OpenCC
	S2t  *opencc.OpenCC
	once sync.Once
}

func NewLangTranslatorUtils() *LangTranslatorUtils {
	t := &LangTranslatorUtils{}
	t.Init()
	return t
}

type Translator struct {
	T2s  *opencc.OpenCC
	S2t  *opencc.OpenCC
	once sync.Once
}

func (t *LangTranslatorUtils) Init() {
	t.once.Do(
		func() {
			var err error
			t.T2s, err = opencc.New("t2s")
			t.S2t, err = opencc.New("s2t")
			if err != nil {
				panic(err)
			}
		},
	)
}
