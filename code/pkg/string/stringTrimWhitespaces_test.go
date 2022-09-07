package slogstring

import (
	"testing"
)

func TestTrimWhitespaces(t *testing.T) {

	tests := map[string]string{
		"": "      			",
		"space-in-the-start":  "      space-in-the-start",
		"space-in-the-end":    "space-in-the-end      ",
		"space-in-the-middle": "space-in        -the-middle",
		"tab-in-the-start": "			tab-in-the-start",
		"tab-in-the-end": "tab-in-the-end		",
		"tab-in-the-middle": "tab-		in-the-middle",
		"mix": "  		m   	i  		  x",
	}

	for k, v := range tests {
		got := TrimWhitespace(v)
		if k != got {
			t.Error("Unexpected result, want: ", k, " got: ", v)
		}
	}
}
