package defvalid

import (
	// "bufio"
	// "encoding/json"
	// "errors"
	// "fmt"
	"encoding/json"
	"github.com/baza-winner/bw/ansi"
	"github.com/baza-winner/bw/defparse"
	"reflect"
	"testing"
	// "github.com/iancoleman/strcase"
	// "github.com/jimlawless/whereami"
	// "log"
	// "os"
	// "os/exec"
	// "strings"
	// "syscall"
)

func TestGetValidVal(t *testing.T) {
	cases := []struct {
		where    string
		val      interface{}
		def      map[string]interface{}
		whereDef string
		result   interface{}
		err      error
	}{
		{
			where: "somewhere",
			// val: defparse.MustParseMap(`{ exitCode: nil, s: 1, v: "ALL", some: 0 }`),
      val: defparse.MustParseMap(`{ }`),
			def: defparse.MustParseMap(`{
        type: 'map',
        keys: {
          v: {
            type: 'enum'
            enum: [ qw/all err ok none/ ]
            default: 'none'
          }
          s: {
            type: 'enum'
            enum: [ qw/none stderr stdout all/ ]
            default: 'all'
          }
          exitOnError: {
            type: 'bool'
            default: false
          }
        }
      }`),
			whereDef: "somewhere::def",
			result:   defparse.MustParseMap(`{ v: 'enum', s: 'all', exitOnError: false }`),
			err:      nil,
		},
	}
	for _, c := range cases {
		got, err := GetValidVal(c.where, c.val, c.def, c.whereDef)
		if err != nil {
			if err != c.err {
				t.Errorf(ansi.Ansi("", "GetValidVal(%s, %+v, %+v, %s)\n    => err:<ansiErr> %v<ansi>\n, want err:<ansiOK>%v"), c.where, c.val, c.def, c.whereDef, err, c.err)
			}
		} else if !reflect.DeepEqual(got, c.result) { // https://stackoverflow.com/questions/18208394/testing-equivalence-of-maps-golang
			gotJson, _ := json.MarshalIndent(got, "", "  ")
			resultJson, _ := json.MarshalIndent(c.result, "", "  ")
      valJson, _ := json.MarshalIndent(c.val, "", "  ")
      defJson, _ := json.MarshalIndent(c.def, "", "  ")
			t.Errorf(ansi.Ansi("", "GetValidVal(%s, %+s, %+s, %s)\n    => <ansiErr>%s<ansi>\n, want <ansiOK>%s"), c.where, valJson, defJson, c.whereDef, gotJson, resultJson)
		}
	}
}
