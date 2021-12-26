package lib

import (
	peparser "github.com/saferwall/pe"
)

type GenInfo struct {
	path      string
	dllName   string
	exports   *peparser.Export
	embedUUID string
}
