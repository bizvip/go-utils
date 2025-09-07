package sqids

import (
	"github.com/rs/zerolog/log"

	"github.com/sqids/sqids-go"
)

var s *sqids.Sqids

// go get github.com/sqids/sqids-go
func init() {
	var err error
	s, err = sqids.New(sqids.Options{
		Alphabet:  "7QAe51ajZrfx4Bg6Sp8YzNhobJIRcOyqVTmnFCsPW9k3G2uE0liwDdHXLUMvKt",
		MinLength: 6,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create sqids instance")
	}
}

func ToAlpha(ids []uint64) string {
	id, _ := s.Encode(ids)
	return id
}

func ToInt(ids string) []uint64 {
	return s.Decode(ids)
}
