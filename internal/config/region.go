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

// Regions lists the regions there are, sorted.
func Regions() []string { return slices.Sorted(maps.Keys(regionHosts)) }

// RegionHost is the host serving region's API and token endpoint, or a usage
// error listing the regions there are.
func RegionHost(region string) (string, error) {
	if h, ok := regionHosts[region]; ok {
		return h, nil
	}
	return "", usage("unknown_region", fmt.Sprintf("unknown region %q", region),
		"use one of: "+strings.Join(Regions(), ", "))
}
