package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
)

// DecodeRawObject decodes a JSON object body into a generic map. Use it for
// responses whose typed SDK model rejects valid payloads (schema drift,
// polymorphic oneOf unions). The formatter renders such maps with the same
// PrintConfig field names as structs, matching camelCase JSON keys.
func DecodeRawObject(body []byte) (map[string]any, error) {
	var object map[string]any
	if err := json.Unmarshal(body, &object); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return object, nil
}

// DecodeRawItems decodes either a bare JSON array or a paginated
// {"data": [...]} envelope into raw items, so callers need not know which
// shape an endpoint returns.
func DecodeRawItems(body []byte) ([]json.RawMessage, error) {
	var items []json.RawMessage
	if err := json.Unmarshal(body, &items); err == nil {
		return items, nil
	}

	var envelope struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("failed to decode response list: %w", err)
	}
	return envelope.Data, nil
}

// PrintRawObject decodes a JSON object body and prints it with printConfig.
func PrintRawObject(body []byte, printConfig *formatter.PrintConfig) error {
	object, err := DecodeRawObject(body)
	if err != nil {
		return err
	}
	return formatter.PrintResult(object, printConfig)
}

// RevisionFromRaw extracts the optimistic-concurrency revision from a raw
// object, tolerating both numeric and string encodings. It returns "" when the
// object carries no revision.
func RevisionFromRaw(object map[string]any) string {
	switch revision := object["revision"].(type) {
	case string:
		return revision
	case float64:
		return strconv.FormatInt(int64(revision), 10)
	case json.Number:
		return revision.String()
	default:
		return ""
	}
}
