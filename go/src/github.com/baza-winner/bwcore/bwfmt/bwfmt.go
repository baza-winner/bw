package bwfmt

type I interface {
	FmtString() string
	FmtArgs() []interface{}
}

type A struct {
	Fmt  string
	Args []interface{}
}

func (v A) FmtString() string {
	return v.Fmt
}

func (v A) FmtArgs() []interface{} {
	return v.Args
}

// func Arg(fmtString string, fmtArgs ...interface{}) A {
// 	return A{fmtString, fmtArgs}
// }
