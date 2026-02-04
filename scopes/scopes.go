package scopes

import (
	"fmt"
	"strings"
)

var scopes = []string{
	"esi-corporations.read_corporation_membership.v1",
}

func Compile() string {
	s := ""

	for _, scope := range scopes {
		s += fmt.Sprintf(" %s", scope)
	}

	s = strings.Trim(s, " ")
	return s
}
