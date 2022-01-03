package cmdutil

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	// TODO: write test cases
	tc := []struct {
		version *Version
		vString string
	}{}

	for _, c := range tc {
		name := fmt.Sprintf("test version formatting of %v", c.vString)
		t.Run(name, func(t *testing.T) {
			vString := c.version.String()
			if vString != c.vString {
				t.Errorf("(*Version).String(): want %v, got %v", c.vString, vString)
			}
		})
	}
}
