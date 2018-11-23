package bwstr

import (
	"strconv"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

func ParseInt(s string) (result int, err error) {
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else if int64(bw.MinInt) <= i && i <= int64(bw.MaxInt) {
		return int(i), nil
	} else {
		return 0, bwerr.From(ansiOutOfRange, i, bw.MinInt, bw.MaxInt)
	}
}

func ParseUint(s string) (result uint, err error) {
	if u, err := strconv.ParseUint(s, 10, 64); err != nil {
		return 0, err
	} else if u <= uint64(bw.MaxUint) {
		return uint(u), nil
	} else {
		return 0, bwerr.From(ansiOutOfRange, u, 0, bw.MaxUint)
	}
}

var ansiOutOfRange string

func init() {
	ansiOutOfRange = ansi.String("<ansiVal>%d<ansi> is out of range <ansiVal>%d..%d")
}
