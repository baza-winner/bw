package bwfmt

type Struct struct {
	FmtString string
	FmtArgs   []interface{}
}

func StructFrom(fmtString string, fmtArgs ...interface{}) Struct {
	return Struct{fmtString, fmtArgs}
}
