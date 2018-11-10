package bwval

import (
	"fmt"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
)

// ============================================================================

func varsJSON(path bw.ValPath, optVars []map[string]interface{}) (result string) {
	if hasVar(path) {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		result = fmt.Sprintf(ansiVars, bwjson.Pretty(vars))
	}
	return
}

func hasVar(path bw.ValPath) bool {
	for _, vpi := range path {
		switch vpi.Type {
		case bw.ValPathItemVar:
			return true
		case bw.ValPathItemPath:
			if hasVar(vpi.Path) {
				return true
			}
		}
	}
	return false
}

// ============================================================================
