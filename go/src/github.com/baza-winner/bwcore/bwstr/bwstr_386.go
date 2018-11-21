package bwstr

import (
	"strconv"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

func ParseInt(s string) (result int, err error) {
	if _int64, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else if int64(bw.MinInt) <= _int64 && _int64 <= int64(bw.MaxInt) {
		return int(_int64), nil
	} else {
		return 0, bwerr.From(ansiOutOfRange, _int64, bw.MinInt, bw.MaxInt)
	}
}

var ansiOutOfRange string

func init() {
	ansiOutOfRange = ansi.String("<ansiVal>%d<ansi> is out of range <ansiVal>%d..%d")
}
