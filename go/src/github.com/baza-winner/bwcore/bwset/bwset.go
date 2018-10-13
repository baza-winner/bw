// Содержит релизацию множеств для базовых типов.
// Являетя демонстрацией работы инструмента setter (go get github.com/baza-winner/bwcore/setter).
package bwset

const (
	_StringSetTestItemA string = "a"
	_StringSetTestItemB string = "b"
)

//go:generate setter -type=string -set=StringSet -test

const (
	_BoolSetTestItemA bool = false
	_BoolSetTestItemB bool = true
)

//go:generate setter -type=bool -set=BoolSet -test -nosort

const (
	_IntSetTestItemA int = 0
	_IntSetTestItemB int = 1
)

//go:generate setter -type=int -set=IntSet -test

const (
	_Int8SetTestItemA int8 = 0
	_Int8SetTestItemB int8 = 1
)

//go:generate setter -type=int8 -set=Int8Set -test

const (
	_Int16SetTestItemA int16 = 0
	_Int16SetTestItemB int16 = 1
)

//go:generate setter -type=int16 -set=Int16Set -test

const (
	_Int32SetTestItemA int32 = 0
	_Int32SetTestItemB int32 = 1
)

//go:generate setter -type=int32 -set=Int32Set -test

const (
	_Int64SetTestItemA int64 = 0
	_Int64SetTestItemB int64 = 1
)

//go:generate setter -type=int64 -set=Int64Set -test

const (
	_UintSetTestItemA uint = 0
	_UintSetTestItemB uint = 1
)

//go:generate setter -type=uint -set=UintSet -test

const (
	_Uint8SetTestItemA uint8 = 0
	_Uint8SetTestItemB uint8 = 1
)

//go:generate setter -type=uint8 -set=Uint8Set -test

const (
	_Uint16SetTestItemA uint16 = 0
	_Uint16SetTestItemB uint16 = 1
)

//go:generate setter -type=uint16 -set=Uint16Set -test

const (
	_Uint32SetTestItemA uint32 = 0
	_Uint32SetTestItemB uint32 = 1
)

//go:generate setter -type=uint32 -set=Uint32Set -test

const (
	_Uint64SetTestItemA uint64 = 0
	_Uint64SetTestItemB uint64 = 1
)

//go:generate setter -type=uint64 -set=Uint64Set -test

const (
	_Float32SetTestItemA float32 = 0
	_Float32SetTestItemB float32 = 1
)

//go:generate setter -type=float32 -set=Float32Set -test

const (
	_Float64SetTestItemA float64 = 0
	_Float64SetTestItemB float64 = 1
)

//go:generate setter -type=float64 -set=Float64Set -test
