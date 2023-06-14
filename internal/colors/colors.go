package colors

import (
	"fmt"

	"github.com/jwalton/gchalk"
)

func SetColoringEnabled(enabled bool) {
	if enabled {
		gchalk.SetLevel(gchalk.LevelBasic)
	} else {
		gchalk.SetLevel(gchalk.LevelNone)
	}
}

func Red(i interface{}) string {
	return gchalk.Red(fmt.Sprintf("%v", i))
}

func Blue(i interface{}) string {
	return gchalk.Blue(fmt.Sprintf("%v", i))
}

func Yellow(i interface{}) string {
	return gchalk.Yellow(fmt.Sprintf("%v", i))
}

func Green(i interface{}) string {
	return gchalk.Green(fmt.Sprintf("%v", i))
}

func Magenta(i interface{}) string {
	return gchalk.Magenta(fmt.Sprintf("%v", i))
}

func WhiteOnRed(i interface{}) string {
	return gchalk.WithBgRed().White(fmt.Sprintf("%v", i))
}

func Bold(i interface{}) string {
	return gchalk.Bold(fmt.Sprintf("%v", i))
}
