package filtering

import (
	"fmt"
	"strings"
)

// ConvertToSearchFieldFormat converts the filtering from "id:a,v,c status:available,unavailable"
// into "id:a id:b id:c status:available status:unavailable"
// it will default to the regular multi-field filter format otherwise
func ConvertToSearchFieldFormat(filter string) string {

	conditions := []string{}

	trimmedFilter := strings.Trim(filter, " ")
	trimmedFilter = strings.Replace(trimmedFilter, "=", ":", -1)

	parts := strings.Split(trimmedFilter, " ")

	for _, part := range parts {
		condition := strings.Split(part, ":")
		if len(condition) == 2 {
			variants := strings.Split(condition[1], ",")
			for _, variant := range variants {
				conditions = append(conditions, fmt.Sprintf("+%s:%s", condition[0], variant))
			}
		}
	}

	if len(conditions) > 0 {
		return strings.Join(conditions, " ")
	}

	return filter

}
