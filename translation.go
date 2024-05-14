package utils

import (
	"sync"

	"github.com/longbridgeapp/opencc"
)

type Translator struct {
	T2s  *opencc.OpenCC
	S2t  *opencc.OpenCC
	once sync.Once
}

func (t *Translator) Init() {
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

func NewTranslator() *Translator {
	t := &Translator{}
	t.Init()
	return t
}
