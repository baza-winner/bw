
# =============================================================================

_resetBash

# =============================================================================

_prepareCodeToParseFuncOptions2TestFuncOptionsOpt=()
_prepareCodeToParseFuncOptions2TestFuncBoolOptions=( 'treatUnknownOptionAsArg' )
_prepareCodeToParseFuncOptions2TestFuncScalarOptions=( 'prefix' )
_prepareCodeToParseFuncOptions2TestFuncListOptions=( 'run' 'bool' 'scalar' 'list' )
_prepareCodeToParseFuncOptions2TestFunc() { eval "$_funcOptions2"
  local -a _prepareCodeToParseFuncOptions2TestFuncBoolOptions=( "${bool[@]}" )
  local -a _prepareCodeToParseFuncOptions2TestFuncScalarOptions=( "${scalar[@]}" )
  local -a _prepareCodeToParseFuncOptions2TestFuncListOptions=( "${list[@]}" )

  local -a __run=( "${run[@]}" )

  local -a _prepareCodeToParseFuncOptions2TestFuncOptionsOpt=(
    --additionalSuffix $testNum
    "${OPT_prefix[@]}"
    "${OPT_treatUnknownOptionAsArg[@]}"
  )
  eval "$_funcOptions2"

  local item; for item in "${__run[@]}"; do
    eval "$item"
  done
}
_prepareCodeToParseFuncOptions2Tests=(
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что ${_ansiOutline}имя переменной ${_ansiPrimaryLiteral}not valid id${_ansiErr} для опции ${_ansiCmd}not valid id${_ansiErr} $_mustBeValidVarName${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --scalar '$_stq'not valid id'$_stq' \
    "
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} не ожидает опцию ${_ansiCmd}--opt-a${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      -- \
      --opt-a \
    "
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--scalar-opt${_ansiErr} будет снабжена значением${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --scalar scalar-opt \
      -- \
      --scalar-opt \
    "
  '
  '
    --return "0"
    --stderr "
      ${_ansiOutline}boolOptA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}boolOptB${_ansiReset}<empty>
      ${_ansiOutline}boolOptC${_ansiReset}<empty>
      ${_ansiOutline}scalarOptA${_ansiReset}: ${_ansiPrimaryLiteral}thing${_ansiReset}
      ${_ansiOutline}scalarOptB${_ansiReset}<empty>
      ${_ansiOutline}listOptA${_ansiReset}(4): ${_ansiSecondaryLiteral}itemA \"item B\" itemC itemD${_ansiReset}
      ${_ansiOutline}listOptB${_ansiReset}(0)
      ${_ansiOutline}OPT_boolOptA${_ansiReset}(1): ${_ansiSecondaryLiteral}--bool-opt-a${_ansiReset}
      ${_ansiOutline}OPT_boolOptB${_ansiReset}(0)
      ${_ansiOutline}OPT_boolOptC${_ansiReset}(0)
      ${_ansiOutline}OPT_scalarOptA${_ansiReset}(2): ${_ansiSecondaryLiteral}--scalar-opt-a thing${_ansiReset}
      ${_ansiOutline}OPT_scalarOptB${_ansiReset}(0)
      ${_ansiOutline}OPT_listOptA${_ansiReset}(8): ${_ansiSecondaryLiteral}--list-opt-a itemA --listOptA \"item B\" --listOpt-a itemC --list-optA itemD${_ansiReset}
      ${_ansiOutline}OPT_listOptB${_ansiReset}(0)
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        boolOptA boolOptB boolOptC scalarOptA scalarOptB listOptA listOptB \
        OPT_boolOptA OPT_boolOptB OPT_boolOptC OPT_scalarOptA OPT_scalarOptB OPT_listOptA OPT_listOptB \
      '$_stq' \
      --bool boolOptA \
      --bool boolOptB \
      --bool bool-opt-c \
      --scalar scalar-opt-a \
      --scalar scalarOptB \
      --list list-opt-a \
      --list listOptB \
      -- \
      --bool-opt-a \
      --scalar-opt-a thing \
      --scalarOptA thing \
      --list-opt-a itemA \
      --listOptA '$_stq'item B'$_stq' \
      --listOpt-a itemC \
      --list-optA itemD \
      -- \
      --other-opt \
    "
  '
  '
    --return "0"
    --stderr "
      ${_ansiOutline}treatUnknownOptionAsArg${_ansiReset}<empty>
      ${_ansiOutline}run${_ansiReset}(0)
      ${_ansiOutline}OPT_treatUnknownOptionAsArg${_ansiReset}(0)
      ${_ansiOutline}OPT_run${_ansiReset}(0)
      --unknown-opt --another-unknown-opt arg
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean \
        --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        treatUnknownOptionAsArg run \
        OPT_treatUnknownOptionAsArg OPT_run \
      '$_stq' \
      --run '$_stq'_debug --clean '$_stDollarInQ'@'$_stq' \
      --treatUnknownOptionAsArg \
      --bool treatUnknownOptionAsArg \
      --list run \
      -- \
      --unknown-opt \
      --another-unknown-opt \
      arg \
    "
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} не ожидает, что опция ${_ansiCmd}--scalar-opt-a${_ansiErr} будет указана повторно${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --scalar scalar-opt-a \
      -- \
      --scalar-opt-a valueA \
      --scalar-opt-a valueB \
    "
  '
  # '
  #   --noErrorStack
  #   --return "1"
  #   --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--scalar-opt-a${_ansiErr} будет снабжена значением не похожим на опцию ${_ansiPrimaryLiteral}-value${_ansiReset}"
  #   "_prepareCodeToParseFuncOptions2TestFunc \
  #     --scalar scalar-opt-a \
  #     -- \
  #     --scalar-opt-a -value \
  #   "
  # '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--scalar-opt-a${_ansiErr} будет снабжена значением${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --scalar scalar-opt-a \
      -- \
      --scalar-opt-a \
    "
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--list-opt-a${_ansiErr} будет снабжена значением${_ansiReset}"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --list list-opt-a \
      -- \
      --list-opt-a \
    "
  '
  # '
  #   --noErrorStack
  #   --return "1"
  #   --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncOptions2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--list-opt-a${_ansiErr} будет снабжена значением не похожим на опцию ${_ansiPrimaryLiteral}-value${_ansiReset}"
  #   "_prepareCodeToParseFuncOptions2TestFunc \
  #     --list list-opt-a \
  #     -- \
  #     --list-opt-a -value \
  #   "
  # '

# Влияние treatUnknownOptionAsArg на поглощение --
  '
    --noErrorStack
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}<${_ansiErr}is unset${_ansiReset}>
      ${_ansiOutline}optC${_ansiReset}<${_ansiErr}is unset${_ansiReset}>
      --opt-b -- --opt-c
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        optA \
        optB \
        optC \
      '$_stq' \
      --run '$_stq'echo '$_stDollarInQ'@ >&2'$_stq' \
      --treatUnknownOptionAsArg \
      --bool optA \
      -- \
      --opt-a --opt-b -- --opt-c \
    "
  '
  '
    --noErrorStack
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optC${_ansiReset}<${_ansiErr}is unset${_ansiReset}>
      -- --opt-c
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        optA \
        optC \
      '$_stq' \
      --run '$_stq'echo '$_stDollarInQ'@ >&2'$_stq' \
      --treatUnknownOptionAsArg \
      --bool optA \
      -- \
      --opt-a -- --opt-c \
    "
  '
  '
    --noErrorStack
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optC${_ansiReset}<${_ansiErr}is unset${_ansiReset}>
      --opt-c
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        optA \
        optC \
      '$_stq' \
      --run '$_stq'echo '$_stDollarInQ'@ >&2'$_stq' \
      --bool optA \
      -- \
      --opt-a -- --opt-c \
    "
  '

# Поправил bug, который не позволял значению скалярной опции начинаться с пробела
  '
    --noErrorStack
    --stderr "
      ${_ansiOutline}opt${_ansiReset}: ${_ansiPrimaryLiteral}\" some\"${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncOptions2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        opt \
      '$_stq' \
      --scalar opt \
      -- \
      --opt '$_stq' some'$_stq' \
    "
  '
)

# =============================================================================
