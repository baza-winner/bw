package defparse

//go:generate bwsetter -type=parseStackItemType
//go:generate bwsetter -type=parsePrimaryState
//go:generate bwsetter -type=parseSecondaryState

//go:generate stringer -type=unicodeCategory,parseStackItemType,parsePrimaryState,parseSecondaryState,pfaErrorType,ruleKind
