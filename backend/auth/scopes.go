package auth

import (
	"fmt"
	"strings"
)

type Scopes []string

var ScopeList = Scopes{
	"esi-corporations.read_corporation_membership.v1",
	"esi-corporations.read_structures.v1",
	"esi-corporations.track_members.v1",
	"esi-corporations.read_divisions.v1",
	"esi-corporations.read_contacts.v1",
	"esi-corporations.read_titles.v1",
	"esi-corporations.read_blueprints.v1",
	"esi-corporations.read_standings.v1",
	"esi-corporations.read_starbases.v1",
	"esi-corporations.read_container_logs.v1",
	"esi-corporations.read_facilities.v1",
	"esi-corporations.read_medals.v1",
	"esi-alliances.read_contacts.v1",
	"esi-corporations.read_fw_stats.v1",
	"esi-corporations.read_projects.v1",
	"esi-corporations.read_freelance_jobs.v1",
}

func (s Scopes) Compile() string {
	compiled := ""

	for _, scope := range s {
		compiled += fmt.Sprintf(" %s", scope)
	}

	compiled = strings.Trim(compiled, " ")
	return compiled
}
