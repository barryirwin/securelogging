package slogstring

import "strings"

// TrimWhitespace :
//
// Removes ALL the whitespaces from the string, for only trailing and heading ones, use strings.Trimspace()
func TrimWhitespace(s string) string {
	out := ""
	out = strings.TrimSpace(s)
	out = strings.Replace(out, " ", "", -1)
	out = strings.Replace(out, "	", "", -1)
	return out
}
