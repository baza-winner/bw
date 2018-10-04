package defvalid

import (
	"fmt"
	"github.com/baza-winner/bwcore/bwerror"
	"github.com/baza-winner/bwcore/bwjson"
	"reflect"
)

type value struct {
	where string
	value interface{}
	error *valueError
}

func (v value) String() string {
	return fmt.Sprintf(v.where+`<ansi> (<ansiSecondaryLiteral>%s<ansi>)`, bwjson.PrettyJson(v.value))
}

func (v *value) asMap() (result map[string]interface{}, err error) {
	var ok bool
	if result, ok = v.value.(map[string]interface{}); !ok {
		err = v.err(valueErrorIsNotOfTypes, "map")
	}
	return
}

func (v *value) mustBeMap() (result map[string]interface{}) {
	var err error
	if result, err = v.asMap(); err != nil {
		bwerror.Panic(err.Error())
	}
	return
}

func (v *value) asString() (result string, err error) {
	var ok bool
	if result, ok = v.value.(string); !ok {
		err = v.err(valueErrorIsNotOfTypes, "string")
	}
	return
}

func (v *value) mustBeString() (result string) {
	var err error
	if result, err = v.asString(); err != nil {
		bwerror.Panic(err.Error())
	}
	return
}

func (v *value) asBool() (result bool, err error) {
	var ok bool
	if result, ok = v.value.(bool); !ok {
		err = v.err(valueErrorIsNotOfTypes, "bool")
	}
	return
}

func (v *value) getKey(keyName string, opts ...interface{}) (result value, err error) {
	var ofTypes *[]string
	var defaultValue *interface{}
	var ofTypeStrings []string
	if opts != nil {
		if _isOfType(opts[0], "string") {
			ofTypeStrings = []string{_mustBeString(opts[0])}
		} else if _isOfType(opts[0], "[]string") {
			ofTypeStrings = _mustBeSliceOfStrings(opts[0])
		} else {
			_ = _mustBeOfTypes(opts[0], "string", "[]string")
		}
		ofTypes = &ofTypeStrings
		if len(opts) > 1 {
			defaultValueIntf := opts[1]
			defaultValue = &defaultValueIntf
		}
		if len(opts) > 2 {
			bwerror.Panic("expects max 2 opts (ofTypes, defaultValue), but found <ansiSecondaryLiteral>%v", opts)
		}
	}
	var m map[string]interface{}
	if m, err = v.asMap(); err == nil {
		var ok bool
		result.where = v.where + "." + keyName
		if result.value, ok = m[keyName]; !ok {
			if defaultValue == nil {
				err = v.err(valueErrorHasNoKey, keyName)
			} else {
				result.value = *defaultValue
			}
		} else if ofTypes != nil && !_isOfType(result.value, (*ofTypes)...) {
			var ofTypeIntfs = []interface{}{}
			for _, i := range ofTypeStrings {
				ofTypeIntfs = append(ofTypeIntfs, i)
			}
			err = result.err(valueErrorIsNotOfTypes, ofTypeIntfs...)
		}
	}
	return
}

func _isOfType(v interface{}, ofTypes ...string) (ok bool) {
	vType := reflect.TypeOf(v)
	for _, ofType := range ofTypes {
		switch ofType {
		case "string":
			ok = vType.Kind() == reflect.String
		case "[]string":
			ok = vType.Kind() == reflect.Slice && vType.Elem().Kind() == reflect.String
		case "map":
			ok = vType.Kind() == reflect.Map && vType.Key().Kind() == reflect.String && vType.Elem().Kind() == reflect.Interface
		case "bool":
			ok = vType.Kind() == reflect.Bool
		default:
			bwerror.Panic("unsupported type <ansiPrimaryLiteral>%s", ofType)
		}
		if ok {
			break
		}
	}
	return
}

func _mustBeOfTypes(v interface{}, ofTypes ...string) (result interface{}) {
	if !_isOfType(v, ofTypes...) {
		bwerror.Panic("<ansiSecondaryLiteral>%+v<ansi> is not of types <ansiSecondaryLiteral>%v", v, ofTypes)
	}
	return v
}

func _mustBeString(v interface{}) (result string) {
	result, _ = _mustBeOfTypes(v, "string").(string)
	return
}

func _mustBeSliceOfStrings(v interface{}) (result []string) {
	result, _ = _mustBeOfTypes(v, "[]string").([]string)
	return
}
