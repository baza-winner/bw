package defparser

type parseStackItemType uint16

const (
	_parseStackItemBelow parseStackItemType = iota
	parseStackItemKey
	parseStackItemString
	parseStackItemMap
	parseStackItemArray
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemNumber
	parseStackItemWord
	_parseStackItemAbove
)

//go:generate stringer -type=parseStackItemType
