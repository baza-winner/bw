package bw

import "github.com/davecgh/go-spew/spew"

var Spew spew.ConfigState

func init() {
	Spew = spew.ConfigState{SortKeys: true}
}

func Args(args ...interface{}) []interface{} {
	return args
}
