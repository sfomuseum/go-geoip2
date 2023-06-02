// Package data provides embedded data files.
package data

import (
	_ "embed"
)

//go:embed GeoLite2-City.mmdb
var GeoLite2City []byte
