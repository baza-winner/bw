package where

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwjson"
)

type W struct {
	funcSpec string
	fileSpec string
	line     int
}

type WW []W

func WWFrom(depth uint) (result WW) {
	result = WW{}
	for i := depth + 1; true; i++ {
		w, ok := From(i)
		if !ok {
			break
		}
		result = append(result, w)
	}
	return
}

func From(depth uint) (result W, ok bool) {
	function, file, line, ok := runtime.Caller(int(depth) + 1)
	if !ok {
		return
	}
	result = W{
		runtime.FuncForPC(function).Name(),
		file,
		line,
	}
	return
}

func MustFrom(depth uint) W {
	if result, ok := From(depth + 1); ok {
		return result
	}
	panic(fmt.Sprintf("WhereFrom(depth: %d)", depth))
}

func (v W) FuncSpec() string {
	return v.funcSpec
}

func (v W) FuncName() (result string) {
	result = v.funcSpec
	pointPos := strings.LastIndex(result, ".") + 1
	if pointPos > 0 {
		result = result[pointPos:]
	}
	return
}

func (v W) FileSpec() string {
	return v.fileSpec
}

func (v W) FileName() string {
	return chopPath(v.fileSpec)
}

func (v W) Line() int {
	return v.line
}

var ansiWhere string

func init() {
	ansi.AddTag("ansiWhereFuncName",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansi.AddTag("ansiWhereAt",
		ansi.SGRCodeOfColor256(ansi.Color256{Code: 243}),
	)
	ansi.AddTag("ansiWhereFile",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorWhite, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansiWhere = ansi.String("<ansiWhereFuncName>%s<ansiWhereAt>@<ansiWhereFile>%s:%d")
}

func (v WW) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, bw.Spew.Sprintf("%s@%s:%d", i.FuncSpec(), i.FileName(), i.Line()))
	}
	return json.Marshal(result)
}

func (v WW) String() string {
	return bwjson.Pretty(v)
}

func (v W) String() string {
	return bw.Spew.Sprintf(ansiWhere, v.FuncSpec(), v.FileName(), v.Line())
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
