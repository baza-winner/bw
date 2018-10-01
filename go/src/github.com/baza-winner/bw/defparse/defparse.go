/*
Предоставляет JSON-like парсер текста вида:

 [
	 {
		 "keyOfStringValue": "stringValue",
		 "keyOfBoolValue": false,
		 "keyOfNumberValue": 12345000.678001
	 }, {
	 	  "keyOfNull": null,
	 	  "keyOfNil": nil,
		  "keyOfArrayValue": [ "stringValue", true, 876.54321 ],
		  "keyOfMapValue": {
		 	  "key1": "value1",
		 	  "key2": true,
		 	  "key3": -3.14,
		 	  "key4": nil,
			  "key5": [ "one", "two", "three" ]
		 }
	 }
 ]


Но с менее строгими требованиями, а именно:

1. Позволяет ставить запятую перед ] или }:

 [
	 {
		 "keyOfStringValue": "stringValue",
		 "keyOfBoolValue": false,
		 "keyOfNumberValue": 12345000.678001,
	 }, {
	 	  "keyOfNull": null,
	 	  "keyOfNil": nil,
		  "keyOfArrayValue": [ "stringValue", true, 876.54321 ],
		  "keyOfMapValue": {
		 	  "key1": "value1",
		 	  "key2": true,
		 	  "key3": -3.14,
		 	  "key4": nil,
			  "key5": [ "one", "two", "three" ],
		 }
	 }
 ]

2. Позволяет не ставить запятую внутри [...] и {...}:

 [
	 {
		 "keyOfStringValue": "stringValue"
		 "keyOfBoolValue": false
		 "keyOfNumberValue": 12345000.678001
	 } {
	 	  "keyOfNull": null
	 	  "keyOfNil": nil
		  "keyOfArrayValue": [ "stringValue", true, 876.54321 ]
		  "keyOfMapValue": {
		 	  "key1": "value1"
		 	  "key2": true
		 	  "key3": -3.14
		 	  "key4": nil
			  "key5": [ "one", "two", "three" ]
		 }
	 }
 ]

3. Позволяет не заключать ключ Map'а в кавычки, если ключ начинается с буквы и не содержит пробелов:

 [
	 {
		 keyOfStringValue: "stringValue"
		 keyOfBoolValue: false
		 keyOfNumberValue: 12345000.678001
	 } {
	 	  keyOfNull: null
	 	  keyOfNil: nil
		  keyOfArrayValue: [ "stringValue", true, 876.54321 ]
		  keyOfMapValue: {
		 	  key1: "value1"
		 	  key2: true
		 	  key3: -3.14
		 	  key4: nil
			  key5: [ "one", "two", "three" ]
		 }
	 }
 ]

4. Позволяет использовать для разделения ключа и значения в Map'е rocket ('=>'), как в Perl'е, вместо ':':

 [
	 {
		 keyOfStringValue => "stringValue"
		 keyOfBoolValue => false
		 keyOfNumberValue => 12345000.678001
	 } {
	 	  keyOfNull => null
	 	  keyOfNil => nil
		  keyOfArrayValue => [ "stringValue", true, 876.54321 ]
		  keyOfMapValue => {
		 	  key1 => "value1"
		 	  key2 => true
		 	  key3 => -3.14
		 	  key4 => nil
			  key5 => [ "one", "two", "three" ]
		 }
	 }
 ]

5. Или вообще не ставить разделитель между ключом и значением, как и между элементами массива:

 [
	 {
		 keyOfStringValue "stringValue"
		 keyOfBoolValue false
		 keyOfNumberValue 12345000.678001
	 } {
	 	  keyOfNull null
	 	  keyOfNil nil
		  keyOfArrayValue [ "stringValue" true 876.54321 ]
		  keyOfMapValue {
		 	  key1 "value1"
		 	  key2 true
		 	  key3 -3.14
		 	  key4 nil
			  key5 [ "one" "two" "three" ]
		 }
	 }
 ]

6. Поддерживает qw// внутри [...], как в Perl (http://oooportal.ru/?cat=article&id=463)

 [
	 {
		 keyOfStringValue "stringValue"
		 keyOfBoolValue false
		 keyOfNumberValue 12345000.678001
	 } {
	 	  keyOfNull null
	 	  keyOfNil nil
		  keyOfArrayValue [ "stringValue" true 876.54321 ]
		  keyOfMapValue {
		 	  key1 "value1"
		 	  key2 true
		 	  key3 -3.14
		 	  key4 nil
			  key5 [ qw/one two three/ ]
		 }
	 }
 ]

7. Поддерживает underscore ('_') внутри Number как в Swift (http://omarrr.com/underscores-in-apples-swift-numbers/)

 [
	 {
		 keyOfStringValue "stringValue"
		 keyOfBoolValue false
		 keyOfNumberValue 12_345_000.678_001
	 } {
	 	  keyOfNull null
	 	  keyOfNil nil
		  keyOfArrayValue [ "stringValue" true 876.543_21 ]
		  keyOfMapValue {
		 	  key1 "value1"
		 	  key2 true
		 	  key3 -3.14
		 	  key4 nil
			  key5 [ qw/one two three/ ]
		 }
	 }
 ]
*/
package defparse

import (
	"github.com/baza-winner/bw/core"
)

/*
ParseMap - парсит строку с определением Map
*/
func ParseMap(source string) (result map[string]interface{}, err error) {
	var _result interface{}
	if _result, err = Parse(source); err == nil {
		var ok bool
		if result, ok = _result.(map[string]interface{}); !ok {
			err = core.Error("is not Map definition: <ansiSecondaryLiteral>%s", source)
		}
	}
	return result, err
}

/*
MustParseMap is like ParseMap but panics if the expression cannot be parsed.
It simplifies safe initialization of global variables holding parsed values.
*/
func MustParseMap(source string) (result map[string]interface{}) {
	var err error
	if result, err = ParseMap(source); err != nil {
		core.Panic(err.Error())
	}
	return result
}

/*
Parse - парсит строку
*/
func Parse(source string) (interface{}, error) {
	pfa := pfaStruct{stack: parseStack{}, state: parseState{primary: expectValueOrSpace}}
	var err error
	for pos, char := range source {
		if err = pfa.processCharAtPos(char, pos); err != nil {
			break
		}
	}
	if err == nil {
		err = pfa.processEOF()
	}
	if err != nil {
		err = pfa.arrangeError(err, source)
	}
	return pfa.result, err
}

/*
MustParse is like Parse but panics if the expression cannot be parsed.
It simplifies safe initialization of global variables holding parsed values.
*/
func MustParse(source string) (result interface{}) {
	var err error
	if result, err = Parse(source); err != nil {
		core.Panic(err.Error())
	}
	return result
}
