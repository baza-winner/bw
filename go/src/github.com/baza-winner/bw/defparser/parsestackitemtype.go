package defparser

type parseStackItemType uint16

const (
	parseStackItemArray parseStackItemType = iota
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemMap
	parseStackItemNumber
	parseStackItemString
	parseStackItemWord
	parseStackItemKey
)

//go:generate stringer -type=parseStackItemType
