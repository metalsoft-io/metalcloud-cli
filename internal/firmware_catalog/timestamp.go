package firmware_catalog

import "time"

// parseVendorTimestamp converts a vendor supplied date string into a *time.Time,
// as required by the SDK. It returns nil when the value is empty or can't be parsed.
func parseVendorTimestamp(value string) *time.Time {
	if value == "" {
		return nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return &t
		}
	}

	return nil
}
