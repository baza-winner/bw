// Package bwset содержит релизации множеств для базовых типов: StringSet, BoolSet, Int*Set, Uint*Set, Float*Set, RuneSet, InterfaceSet.
// Являетя демонстрацией работы инструмента bwsetter (go get github.com/baza-winner/bwcore/bwsetter).
package bwset

const (
	_StringTestItemA string = "a"
	_StringTestItemB string = "b"
)

//go:generate bwsetter -type=string -set=String -test

const (
	_BoolTestItemA bool = false
	_BoolTestItemB bool = true
)

//go:generate bwsetter -type=bool -set=Bool -test -nosort

const (
	_IntTestItemA int = 0
	_IntTestItemB int = 1
)

//go:generate bwsetter -type=int -set=Int -test

const (
	_Int8TestItemA int8 = 0
	_Int8TestItemB int8 = 1
)

//go:generate bwsetter -type=int8 -set=Int8 -test

const (
	_Int16TestItemA int16 = 0
	_Int16TestItemB int16 = 1
)

//go:generate bwsetter -type=int16 -set=Int16 -test

const (
	_Int32TestItemA int32 = 0
	_Int32TestItemB int32 = 1
)

//go:generate bwsetter -type=int32 -set=Int32 -test

const (
	_Int64TestItemA int64 = 0
	_Int64TestItemB int64 = 1
)

//go:generate bwsetter -type=int64 -set=Int64 -test

const (
	_UintTestItemA uint = 0
	_UintTestItemB uint = 1
)

//go:generate bwsetter -type=uint -set=Uint -test

const (
	_Uint8TestItemA uint8 = 0
	_Uint8TestItemB uint8 = 1
)

//go:generate bwsetter -type=uint8 -set=Uint8 -test

const (
	_Uint16TestItemA uint16 = 0
	_Uint16TestItemB uint16 = 1
)

//go:generate bwsetter -type=uint16 -set=Uint16 -test

const (
	_Uint32TestItemA uint32 = 0
	_Uint32TestItemB uint32 = 1
)

//go:generate bwsetter -type=uint32 -set=Uint32 -test

const (
	_Uint64TestItemA uint64 = 0
	_Uint64TestItemB uint64 = 1
)

//go:generate bwsetter -type=uint64 -set=Uint64 -test

const (
	_Float32TestItemA float32 = 0
	_Float32TestItemB float32 = 1
)

//go:generate bwsetter -type=float32 -set=Float32 -test

const (
	_Float64TestItemA float64 = 0
	_Float64TestItemB float64 = 1
)

//go:generate bwsetter -type=float64 -set=Float64 -test

const (
	_RuneTestItemA rune = 'a'
	_RuneTestItemB rune = 'b'
)

//go:generate bwsetter -type=rune -set=Rune -test

const (
	_InterfaceTestItemA bool   = true
	_InterfaceTestItemB string = "a"
)

//go:generate bwsetter -type=interface{} -set=Interface -test  -nosort
