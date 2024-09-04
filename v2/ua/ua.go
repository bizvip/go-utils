/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

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
