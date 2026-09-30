package config

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Each region serves its public API and its OAuth token endpoint from one
// host (token endpoints verified for eu, us and au on 2026-09-29).
var regionHosts = map[string]string{
	"eu": "app.aikido.dev",
	"us": "app.us.aikido.dev",
	"au": "app.au.aikido.dev",
	"me": "app.me.aikido.dev",
}

const DefaultRegion = "eu"

func regionHost(region string) (string, error) {
	if h, ok := regionHosts[region]; ok {
		return h, nil
	}
	return "", usage("unknown_region", fmt.Sprintf("unknown region %q", region),
		"use one of: "+strings.Join(slices.Sorted(maps.Keys(regionHosts)), ", "))
}
