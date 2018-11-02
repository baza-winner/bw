package bwval

import (
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/formatted"
)

type VarPath interface {
	FormattedString() (result formatted.String)
}

func VarPathFrom(s string) (result VarPath, err error) {

}

func MustVarPath(varPath VarPath, err error) VarPath {

}

type QualVarPath interface {
	// TypeIdxKey(i int) (itemType VarPathItemType, idx int, key string, err error)
	FormattedString(optVal ...Val) (result formatted.String)
}

type Opts struct {
	Vars   map[string]Val
	Consts map[string]Val
}

func QualVarPathFrom(varPath VarPath, opts ...Opts) (result QualVarPath, err error) {
	bwerror.TODO()
	return
}

type Val interface {
	PathVal(path QualVarPath) (result Val, err error)
	SetPathVal(path QualVarPath, val []interface{}) (err error)
	Val
}
