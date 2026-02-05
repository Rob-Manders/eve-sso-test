package scopes

import (
	"fmt"
	"strings"
)

type Scopes []string

var ScopeList = Scopes{
	"esi-corporations.read_corporation_membership.v1",
}

func (s Scopes) Compile() string {
	compiled := ""

	for _, scope := range s {
		compiled += fmt.Sprintf(" %s", scope)
	}

	compiled = strings.Trim(compiled, " ")
	return compiled
}
