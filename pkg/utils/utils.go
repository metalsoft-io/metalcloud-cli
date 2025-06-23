package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func CreateSlug(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}

	processedString := reg.ReplaceAllString(input, " ")

	processedString = strings.TrimSpace(processedString)

	slug := strings.ReplaceAll(processedString, " ", "-")

	slug = strings.ToLower(slug)

	return slug
}

func GetFloat32FromString(input string) (float32, error) {
	result, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, err
	}

	return float32(result), nil
}

func ReadConfigFromPipe() ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data piped to stdin")
	}

	var content = make([]byte, 100000)
	size, err := os.Stdin.Read(content)
	if err != nil {
		return nil, err
	}

	return content[:size], nil
}

func ReadConfigFromFile(configSource string) ([]byte, error) {
	content, err := os.ReadFile(configSource)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ReadConfigFromPipeOrFile(configSource string) ([]byte, error) {
	if configSource == "" {
		return nil, nil
	}

	if configSource == "pipe" {
		return ReadConfigFromPipe()
	}

	return ReadConfigFromFile(configSource)
}

func UnmarshalContent(content []byte, destination any) error {
	if len(content) == 0 {
		return fmt.Errorf("no content to unmarshal")
	}

	format := strings.ToLower(viper.GetString(formatter.ConfigFormat))

	switch format {
	case "json":
		err := json.Unmarshal(content, destination)
		if err != nil {
			return fmt.Errorf("failed to unmarshal content: %w", err)
		}
	case "yaml":
		err := yaml.Unmarshal(content, destination)
		if err != nil {
			return fmt.Errorf("failed to unmarshal content: %w", err)
		}
	default:
		err := json.Unmarshal(content, destination)
		if err != nil {
			err = yaml.Unmarshal(content, destination)
			if err != nil {
				return fmt.Errorf("failed to unmarshal content: %w", err)
			}
		}
	}

	return nil
}

func ProcessFilterStringSlice(filter []string) []string {
	parts := make([]string, len(filter))

	for i, part := range filter {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "-") {
			parts[i] = "$or:$not:$eq:" + strings.TrimPrefix(part, "-")
		} else {
			parts[i] = "$or:$eq:" + part
		}
	}

	return parts
}
