package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func ProcessFilterStringSlice(filter []string) []string {
	parts := make([]string, len(filter))

	for i, part := range filter {
		parts[i] = strings.TrimSpace(part)

		if strings.HasPrefix(part, "-") {
			parts[i] = "$not:$eq:" + strings.TrimPrefix(part, "-")
		} else {
			parts[i] = "$eq:" + part
		}
	}

	return parts
}
