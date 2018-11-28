// +build debug

package bwparse

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/bwerr/where"
)

func (p PosInfo) MarshalJSON() ([]byte, error) {

	result := map[string]interface{}{}
	result["isEOF"] = p.isEOF
	result["rune"] = string(p.rune)
	result["pos"] = p.pos
	result["line"] = p.line
	result["col"] = p.col
	result["prefix"] = p.prefix
	result["prefixStart"] = p.prefixStart
	if p.justParsed != nil {
		result["justParsed"] = p.justParsed
	}
	return json.Marshal(result)
}

func (start Start) MarshalJSON() ([]byte, error) {
	result := map[string]interface{}{}
	result["ps"] = *start.ps
	result["suffix"] = start.suffix
	if len(start.stopped) > 0 {
		result["stopped"] = start.stopped
	}
	return json.Marshal(result)
}

type Start struct {
	ps      *PosInfo
	suffix  string
	stopped where.WW
}

func (p *P) Stop(start *Start) {
	if len(start.stopped) > 0 {
		return
	}
	start.stopped = where.WWFrom(1)
	delete(p.starts, start.ps.pos)
}

func (p *proxy) Stop(start *Start) {
	if len(start.stopped) > 0 {
		return
	}
	start.stopped = where.WWFrom(1)
	delete(p.starts, start.ps.pos)
}
