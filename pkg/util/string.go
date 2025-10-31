package util

import "strings"

const EMPTY_STRING = ""

func EmptyString(value string) bool {
	return strings.TrimSpace(value) == EMPTY_STRING
}
