package defparse

type parseStackItemType uint16

const (
	parseStackItem_below_ parseStackItemType = iota
	parseStackItemKey
	parseStackItemString
	parseStackItemMap
	parseStackItemArray
	parseStackItemQw
	parseStackItemQwItem
	parseStackItemNumber
	parseStackItemWord
	parseStackItem_above_
)

//go:generate stringer -type=parseStackItemType
