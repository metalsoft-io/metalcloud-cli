package main

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

func setColoringEnabled(enabled bool) {
	if enabled {
		gchalk.SetLevel(gchalk.LevelBasic)
	} else {
		gchalk.SetLevel(gchalk.LevelNone)
	}
}

func red(i interface{}) string {
	return gchalk.Red(fmt.Sprintf("%v", i))
}

func blue(i interface{}) string {
	return gchalk.Blue(fmt.Sprintf("%v", i))
}

func yellow(i interface{}) string {
	return gchalk.Yellow(fmt.Sprintf("%v", i))
}

func green(i interface{}) string {
	return gchalk.Green(fmt.Sprintf("%v", i))
}

func magenta(i interface{}) string {
	return gchalk.Magenta(fmt.Sprintf("%v", i))
}

func whiteOnRed(i interface{}) string {
	return gchalk.WithBgRed().White(fmt.Sprintf("%v", i))
}

/*
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	whiteOnRed := color.New(color.FgHiWhite, color.BgRed).SprintFunc()
*/
