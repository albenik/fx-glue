package fxglue

import (
	"fmt"
	"regexp"
)

var tagRegex = regexp.MustCompile(`^supply(?:,(name|group)=(\w+))?$`)

func parseTag(tag string) (string, string, error) {
	m := tagRegex.FindStringSubmatch(tag)
	if m == nil {
		return "", "", fmt.Errorf("tag \"fx\" invalid value %q", tag)
	}
	return m[1], m[2], nil
}
