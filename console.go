package goutils

import (
	"fmt"

	"github.com/fatih/color"
)

type C struct{}

func Console() *C { return &C{} }

func (c *C) BlackBold(txt string) {
	bold := color.New(color.FgHiWhite, color.BgBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *C) Black(txt string) {
	bold := color.New(color.FgHiBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *C) Red(txt string) {
	bold := color.New(color.FgHiRed).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func (c *C) Green(txt string) {
	bold := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s", bold(txt))
}
