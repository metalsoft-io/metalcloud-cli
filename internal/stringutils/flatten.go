package stringutils

import (
	"fmt"
	"strings"
)

func FlattenAndJoinStrings(arr [][]string) string {
	var strs []string

	for _, v1 := range arr {
		s := strings.Join(v1, ", ")
		strs = append(strs, s)

	}
	return strings.Join(strs, ", ")
}

func flattenAndJoinInts(arr [][]int) string {
	var strs []string

	for _, v1 := range arr {
		var sarr []string
		for _, v2 := range v1 {
			sarr = append(sarr, fmt.Sprintf("%d", v2))
		}

		s := strings.Join(sarr, ", ")
		strs = append(strs, s)

	}
	return strings.Join(strs, ", ")
}
