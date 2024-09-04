/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package console

import (
	"fmt"

	"github.com/fatih/color"
)

func BlackBold(txt string) {
	bold := color.New(color.FgHiWhite, color.BgBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func Black(txt string) {
	bold := color.New(color.FgHiBlack).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func Red(txt string) {
	bold := color.New(color.FgHiRed).SprintFunc()
	fmt.Printf("%s", bold(txt))
}

func Green(txt string) {
	bold := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s", bold(txt))
}
