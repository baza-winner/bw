// Содержит релизацию множеств для базовых типов.
// Являетя демонстрацией работы инструмента bwsetter (go get github.com/baza-winner/bwcore/bwsetter).
package bwset

const (
	_StringSetTestItemA string = "a"
	_StringSetTestItemB string = "b"
)

//go:generate bwsetter -type=string -set=StringSet -test

const (
	_BoolSetTestItemA bool = false
	_BoolSetTestItemB bool = true
)

//go:generate bwsetter -type=bool -set=BoolSet -test -nosort

const (
	_IntSetTestItemA int = 0
	_IntSetTestItemB int = 1
)

//go:generate bwsetter -type=int -set=IntSet -test

const (
	_Int8SetTestItemA int8 = 0
	_Int8SetTestItemB int8 = 1
)

//go:generate bwsetter -type=int8 -set=Int8Set -test

const (
	_Int16SetTestItemA int16 = 0
	_Int16SetTestItemB int16 = 1
)

//go:generate bwsetter -type=int16 -set=Int16Set -test

const (
	_Int32SetTestItemA int32 = 0
	_Int32SetTestItemB int32 = 1
)

//go:generate bwsetter -type=int32 -set=Int32Set -test

const (
	_Int64SetTestItemA int64 = 0
	_Int64SetTestItemB int64 = 1
)

//go:generate bwsetter -type=int64 -set=Int64Set -test

const (
	_UintSetTestItemA uint = 0
	_UintSetTestItemB uint = 1
)

//go:generate bwsetter -type=uint -set=UintSet -test

const (
	_Uint8SetTestItemA uint8 = 0
	_Uint8SetTestItemB uint8 = 1
)

//go:generate bwsetter -type=uint8 -set=Uint8Set -test

const (
	_Uint16SetTestItemA uint16 = 0
	_Uint16SetTestItemB uint16 = 1
)

//go:generate bwsetter -type=uint16 -set=Uint16Set -test

const (
	_Uint32SetTestItemA uint32 = 0
	_Uint32SetTestItemB uint32 = 1
)

//go:generate bwsetter -type=uint32 -set=Uint32Set -test

const (
	_Uint64SetTestItemA uint64 = 0
	_Uint64SetTestItemB uint64 = 1
)

//go:generate bwsetter -type=uint64 -set=Uint64Set -test

const (
	_Float32SetTestItemA float32 = 0
	_Float32SetTestItemB float32 = 1
)

//go:generate bwsetter -type=float32 -set=Float32Set -test

const (
	_Float64SetTestItemA float64 = 0
	_Float64SetTestItemB float64 = 1
)

//go:generate bwsetter -type=float64 -set=Float64Set -test

const (
	_RuneSetTestItemA rune = 'a'
	_RuneSetTestItemB rune = 'b'
)

//go:generate bwsetter -type=rune -set=RuneSet -test

const (
	_InterfaceSetTestItemA bool   = true
	_InterfaceSetTestItemB string = "a"
)

//go:generate bwsetter -type=interface{} -set=InterfaceSet -test  -nosort
