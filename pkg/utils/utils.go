package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/internal/colors"
	"github.com/metalsoft-io/metalcloud-cli/internal/command"
)

func FormattedStatus(status string) string {
	switch status {
	case "active":
		return colors.Blue(status)
	case "maintenance":
		return colors.Green(status)
	case "":
		return colors.Green(status)
	default:
		return colors.Yellow(status)
	}
}

func FormattedCapacity(usedPercentage float64, capacity string) string {
	if usedPercentage >= 0.8 {
		capacity = colors.Red(capacity)
	} else if usedPercentage >= 0.5 {
		capacity = colors.Red(capacity)
	} else {
		capacity = colors.Green(capacity)
	}

	return capacity
}

func GetConfirmation(autoconfirm bool, message string) (bool, error) {
	if autoconfirm {
		return true, nil
	}

	confirmationMessage := fmt.Sprintf("%s.  Are you sure? Type \"yes\" to continue:", message)

	// this is simply so that we don't output a text on the command line under go test
	if strings.HasSuffix(os.Args[0], ".test") {
		confirmationMessage = ""
	}

	return command.RequestConfirmation(confirmationMessage)
}
