/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package goutils

import (
	"fmt"

	"github.com/fatih/color"
)

type ConsoleUtils struct{}

func NewConsoleUtils() *ConsoleUtils { return &ConsoleUtils{} }

func (c *ConsoleUtils) BlackBold(txt string) {
	bold := color.New(color.FgHiWhite, color.BgBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *ConsoleUtils) Black(txt string) {
	bold := color.New(color.FgHiBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *ConsoleUtils) Red(txt string) {
	bold := color.New(color.FgHiRed).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *ConsoleUtils) Green(txt string) {
	bold := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s", bold(txt))
}
