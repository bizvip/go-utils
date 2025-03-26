package ua

import (
	"log"

	"github.com/ua-parser/uap-go/uaparser"
)

var parser *uaparser.Parser
var err error

func Parser() *uaparser.Parser {
	parser, err = uaparser.New("./regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}
	return parser
}
