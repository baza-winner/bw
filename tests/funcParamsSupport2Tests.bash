
# =============================================================================

_resetBash

# =============================================================================

_prepareCodeToValidateVarValueAgainstVarTypeWrapper() {
	_prepareCodeToValidateVarValueAgainstVarType "$@" && eval "$codeToValidateVarValueAgainstVarType"
}
_prepareCodeToParseFuncParams2TestFuncBoolOptions=(
  'treatUnknownOptionAsArg'
  'isCommandWrapper'
  'canBeMixedOptionsAndArgs'
  'canBeMoreParams'
)
_prepareCodeToParseFuncParams2TestFuncListOptions=(
  'run'
)
_prepareCodeToParseFuncParams2TestFuncOptionsOpt=(--treatUnknownOptionAsArg)
_prepareCodeToParseFuncParams2TestFunc() { eval "$_funcOptions2"
  local -a _prepareCodeToParseFuncParams2TestFuncParamsOpt=(
    "${OPT_canBeMixedOptionsAndArgs[@]}"
    "${OPT_isCommandWrapper[@]}"
    "${OPT_treatUnknownOptionAsArg[@]}"
    "${OPT_canBeMoreParams[@]}"
  )
  local -a _prepareCodeToParseFuncParams2TestFuncParams=()
  while [[ $# -gt 0 ]]; do
    [[ $1 == -- ]] && shift && break
    _prepareCodeToParseFuncParams2TestFuncParams+=( "$1" )
    shift
  done
  codeHolder=_codeToCallFuncParams2 eval "$_evalCode"

  if [[ -z $isCommandWrapper ]]; then
    local item; for item in "${run[@]}"; do
      eval "$item"
    done
  fi
}

_prepareCodeToParseFuncParams2Tests=(

# ==== compile time (_prepareCodeToParseFuncParams2) errors

# проверка _codeToCheckParams
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2${_ansiErr} не ожидает параметр ${_ansiPrimaryLiteral}--canBeMoreParams${_ansiErr}, используйте ${_ansiOutline}_runBashTestHelperParamsOpt${_ansiErr} для задания опций ${_ansiCmd}_prepareCodeToParseFuncParams2${_ansiReset}"
    "_prepareCodeToParseFuncParams2 --canBeMoreParams"
  '

# проверка  _codeToCheckNoArgsInOpt
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --before "local -a __someFuncParams=() __someFuncParamsOpt=( some )"
    --before "eval \"__someFunc() { eval \\\"\\\$_funcParams2\\\"; }\""
    --after "unset -f __someFunc"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2${_ansiErr} не ожидает в ${_ansiOutline}__someFuncParamsOpt${_ansiErr} ни одного аргумента, но получен ${_ansiPrimaryLiteral}some${_ansiReset}"
    "__someFunc"
  '

# проверка реакции на *Complete
  # '
  #   --before "_substitute noStack true"
  #   --before "local -a __someCompleteParams=( some )"
  #   --before "eval \"__someComplete() { eval \\\"\\\$_funcParams2\\\"; }\""
  #   --after "unset -f __someComplete"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}__someComplete${_ansiErr} не ожидает, что будет определена переменная ${_ansiOutline}__someCompleteParams${_ansiErr}, т.к. все ${_ansiOutline}*${_ansiCmd}Complete${_ansiErr}-функции имеют предопределенный ${_ansiOutline}_completeParams${_ansiErr}: ${_ansiSecondaryLiteral}( --varName --argIdx:0.. \"compWord:?\" )${_ansiReset}"
  #   "__someComplete"
  # '
  # '
  #   --before "_substitute noStack true"
  #   --before "eval \"__someCompleteParams() { __someCompleteParams=( some ); }\""
  #   --before "eval \"__someComplete() { eval \\\"\\\$_funcParams2\\\"; }\""
  #   --after "unset -f __someComplete"
  #   --after "unset -f __someCompleteParams"
  #   --after "unset __someCompleteParams"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}__someComplete${_ansiErr} не ожидает, что будет определена функция ${_ansiCmd}__someCompleteParams${_ansiErr}, т.к. все ${_ansiOutline}*${_ansiCmd}Complete${_ansiErr}-функции имеют предопределенный ${_ansiOutline}_completeParams${_ansiErr}: ${_ansiSecondaryLiteral}( --varName --argIdx:0.. \"compWord:?\" )${_ansiReset}"
  #   "__someComplete"
  # '

# проверка на границы @args
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что в параметре ${_ansiCmd}@4..2args${_ansiErr} левая граница ${_ansiPrimaryLiteral}4${_ansiErr} не должна превосходить правую ${_ansiPrimaryLiteral}2${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      @4..2args \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что в параметре ${_ansiCmd}@..0args${_ansiErr} правая граница ${_ansiPrimaryLiteral}0${_ansiErr} не должна быть меньше ${_ansiPrimaryLiteral}1${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      @..0args \
    "
  '

# проверка следования опций и аргументов в *Params
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что опции будут следовать строго перед аргументами, т.к. ${_ansiCmd}_prepareCodeToParseFuncParams2${_ansiErr} вызвана без ${_ansiCmd}--canBeMixedOptionsAndArgs${_ansiErr}, но обнаружено ${_ansiOutline}определение опции ${_ansiCmd}--option-b${_ansiErr} после ${_ansiOutline}определения аргумента ${_ansiCmd}arg${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --option-a arg --option-b"
  '

# проверка отсутствия аргументов в *Params в режиме --isCommandWrapper
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает определения аргументов в режиме ${_ansiCmd}--isCommandWrapper${_ansiErr}, но получено ${_ansiCmd}arg${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --isCommandWrapper --option-a arg"
  '

# проверка того, что списочный аргумент, если задан, является последним аргументом
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что определение списочного аргумента ${_ansiCmd}@args${_ansiErr} будет последним в списке определений аргументов, но после него следует ещё ${_ansiCmd}argA${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --canBeMixedOptionsAndArgs @args --option-a argA argB"
  '

# проверка отсутсвия модификаторов списочного типа в определении скаляроного параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает тип ${_ansiPrimaryLiteral}sorted${_ansiErr} в определении ${_ansiOutline}скалярного${_ansiErr} параметра ${_ansiCmd}--:sorted=some${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --:sorted=some"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает тип ${_ansiPrimaryLiteral}unique${_ansiErr} в определении ${_ansiOutline}скалярного${_ansiErr} параметра ${_ansiCmd}--:unique=some${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --:unique=some"
  '

# проверка наличия закрывающейся скобки в определении перечислимого типа параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что определение перечислимого типа ${_ansiPrimaryLiteral}(1 \"a b\" 2 3${_ansiErr} для параметра ${_ansiCmd}--a:(1 \"a b\" 2 3=1${_ansiErr} должно заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}(1 \"a b\" 2 3)${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--a:'$_stOpenBraceInQ'1 '$_stqq'a b'$_stqq' 2 3=1'$_stq' \
    "
  '

# проверка наличия закрывающейся скобки в вычисляемом определении перечислимого типа параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что вычисляемое определение перечислимого типа ${_ansiPrimaryLiteral}( \$(echo 1 \"a b\" 2 3)${_ansiErr} для параметра ${_ansiCmd}--a:( \$(echo 1 \"a b\" 2 3)=1${_ansiErr} должно заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}( \$(echo 1 \"a b\" 2 3))${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--a:'$_stOpenBraceInQ' '$_stDollarInQ$_stOpenBraceInQ'echo 1 '$_stqq'a b'$_stqq' 2 3'$_stCloseBraceInQ'=1'$_stq' \
    "
  '

# проверка границ диапазона в определении типа параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что в определении диапазона для параметра ${_ansiCmd}--a:4..2${_ansiErr} левая граница ${_ansiPrimaryLiteral}4${_ansiErr} не превосходит правую ${_ansiSecondaryLiteral}2${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --a:4..2"
  '

# проверка определения типа параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "
      ${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}:!${_ansiErr} в определении параметра ${_ansiCmd}--:!=some${_ansiErr} ожидает:
      -- (необязательно) один из следующих типов:
        перечислимый: ${_ansiSecondaryLiteral}:( ${_ansiOutline}значение1 значение 2 ... ${_ansiSecondaryLiteral})${_ansiErr}
        целочисленный диапазон: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}..${_ansiOutline}max${_ansiErr}
        целочисленный, не менее: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}..${_ansiErr}
        целочисленный, не более: ${_ansiSecondaryLiteral}:..${_ansiOutline}max${_ansiErr}
        целочисленный: ${_ansiSecondaryLiteral}:..${_ansiErr}
        фиксированное целое число: ${_ansiSecondaryLiteral}:${_ansiOutline}intValue${_ansiErr}
      -- и/или (необязательно) указание пустого значения:
        значение может быть пустым: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}?${_ansiErr}
      ${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc --:!=some"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "
      ${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}:!${_ansiErr} в определении параметра ${_ansiCmd}@--:!=some${_ansiErr} ожидает:
      -- (необязательно) один из следующих типов элемента списка:
        перечислимый: ${_ansiSecondaryLiteral}:( ${_ansiOutline}значение1 значение 2 ... ${_ansiSecondaryLiteral})${_ansiErr}
        целочисленный диапазон: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}..${_ansiOutline}max${_ansiErr}
        целочисленный, не менее: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}..${_ansiErr}
        целочисленный, не более: ${_ansiSecondaryLiteral}:..${_ansiOutline}max${_ansiErr}
        целочисленный: ${_ansiSecondaryLiteral}:..${_ansiErr}
        фиксированное целое число: ${_ansiSecondaryLiteral}:${_ansiOutline}intValue${_ansiErr}
      -- и/или (необязательно) указание пустого значения элемента списка:
        значение может быть пустым: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}?${_ansiErr}
      -- и/или (необязательно) указание инструкции для формирования списка значений:
        список будет содержать только уникальные значения: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}unique${_ansiErr}
        список будет содержать отсортированные значения: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}sorted${_ansiErr}
      ${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc @--:!=some"
  '

# проверка имени переменной для параметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает непустое имя переменной для параметра ${_ansiCmd}--:?=some${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --:?=some"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что имя переменной ${_ansiPrimaryLiteral}some option${_ansiErr} для параметра ${_ansiCmd}--some option:1..${_ansiErr} $_mustBeValidVarName${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'--some option:1..'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не может использовать то же самое имя переменной ${_ansiPrimaryLiteral}a${_ansiErr} для параметра ${_ansiCmd}a${_ansiErr}, что было использовано для параметра ${_ansiCmd}--a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --a a"
  '

# проверка соответтсвия значения по умолчанию типу парметра
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что списочный параметр ${_ansiCmd}@--a=1 2 3${_ansiErr} в качестве значения по умолчанию не будет иметь скалярное значение ${_ansiPrimaryLiteral}1 2 3${_ansiErr} (не заключенное в круглые скобки)${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'@--a=1 2 3'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что списочное значение ${_ansiPrimaryLiteral}(1 2 3${_ansiErr} параметра ${_ansiCmd}@--a=(1 2 3${_ansiErr} будет заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}(1 2 3)${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'@--a='$_stOpenBraceInQ'1 2 3'$_stq'"
  '
    '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что скалярный параметр ${_ansiCmd}a=(1 2 3${_ansiErr} в качестве значения по умолчанию не будет иметь списочное значение ${_ansiPrimaryLiteral}(1 2 3${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'a='$_stOpenBraceInQ'1 2 3'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}--a:0..=-1${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'--a:0..=-1'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}--a:1..=0${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'--a:1..=0'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}d${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"d e\"${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}--a:(b \"d e\")=d${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'--a:'$_stOpenBraceInQ'b '$_stqq'd e'$_stqq''$_stCloseBraceInQ'=d'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --before "local -a __enumValues=( b \"d e\" ) __someFuncParams=()"
    --before "eval \"__enumValues() { echo \\\$(_quotedArgs \\\"\\\${__enumValues[@]}\\\"); }\""
    --before "eval \"__someFunc() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f __enumValues"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}__someFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}d${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"e f\"${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}--a:( \$(__enumValues) )=d${_ansiReset}"
    "
      true \
      && local -a __someFuncParams=('$_stq'--a:'$_stOpenBraceInQ' '$_stDollarInQ$_stOpenBraceInQ'__enumValues'$_stCloseBraceInQ' '$_stCloseBraceInQ'=d'$_stq') \
      && __onlyPrepareCode=true __someFunc \
      || __enumValues=( b \"e f\" ) && __someFunc \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}a:0..=-1${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'a:0..=-1'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве значения по умолчанию параметра ${_ansiCmd}a:1..=0${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'a:1..=0'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}d${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"d e\"${_ansiErr} в качестве элемента значения по умолчанию ${_ansiSecondaryLiteral}(d \"b e\")${_ansiErr} параметра ${_ansiCmd}@a:(b \"d e\")=(d \"b e\")${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@a:'$_stOpenBraceInQ'b '$_stqq'd e'$_stqq''$_stCloseBraceInQ'='$_stOpenBraceInQ'd '$_stqq'b e'$_stqq''$_stCloseBraceInQ''$_stq' \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве элемента значения по умолчанию ${_ansiSecondaryLiteral}(1 0 -1)${_ansiErr} параметра ${_ansiCmd}@--a:0..=(1 0 -1)${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'@--a:0..='$_stOpenBraceInQ'1 0 -1'$_stCloseBraceInQ''$_stq'"
  '

# проверка ункальности сокращения для опции
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не может использовать для опции ${_ansiCmd}--opt-b/a${_ansiErr} то же сокращение ${_ansiPrimaryLiteral}a${_ansiErr}, что и для опции ${_ansiCmd}--opt-a${_ansiReset}"
    --return "1"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --opt-a/a \
      --opt-b/a \
    "
  '

# проверки для режима --isCommandWrapper
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что будет определена по крайней мере одна функция ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc_${_ansiOutline}*${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --isCommandWrapper"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "local _prepareCodeToParseFuncParams2TestFunc_alphaShortcuts=( \"alef bravo\" )"
    --after "unset -f _prepareCodeToParseFuncParams2TestFunc_alpha"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что сокращение ${_ansiPrimaryLiteral}\"alef bravo\"${_ansiErr} для команды ${_ansiCmd}alpha $_mustBeValidCommandShortcut${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --isCommandWrapper"
  '
  '
    --before "_substitute noStack true"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "local _prepareCodeToParseFuncParams2TestFunc_alphaShortcuts=( alef bravo )"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_bravo() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeToParseFuncParams2TestFunc_alpha _prepareCodeToParseFuncParams2TestFunc_bravo"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что команда ${_ansiCmd}bravo${_ansiErr} не будет совпадать с одним из сокращений для команды ${_ansiCmd}alpha${_ansiErr}: ${_ansiSecondaryLiteral}alef bravo${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --isCommandWrapper"
  '
  '
    --before "_substitute noStack true"
    --before "local -a _prepareCodeToParseFuncParams2TestFunc_alphaParams=( @5..3arg )"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeToParseFuncParams2TestFunc_alpha"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc_alpha${_ansiErr} ожидает, что в параметре ${_ansiCmd}@5..3arg${_ansiErr} левая граница ${_ansiPrimaryLiteral}5${_ansiErr} не должна превосходить правую ${_ansiPrimaryLiteral}3${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --isCommandWrapper -- alpha "
  '

# ==== runtime (_parseFuncParams2) errors

# проверка опции
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает опцию ${_ansiCmd}--a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc -- --a"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает опцию ${_ansiCmd}--b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --a= -- --b"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--a${_ansiErr} будет снабжена ${_ansiOutline}значением${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --a= -- --a"
  '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает, что опция ${_ansiCmd}--a${_ansiErr} будет снабжена значением (${_ansiPrimaryLiteral}-b${_ansiErr}) не похожим на опцию${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --a= -- --a -b"
  # '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает ${_ansiOutline}непустое${_ansiErr} значение для опции ${_ansiCmd}--a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc --a= -- --a '$_stqq$_stqq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает, что опция ${_ansiCmd}-b${_ansiErr} будет указана повторно с другим значением ${_ansiPrimaryLiteral}thing${_ansiErr} против ${_ansiPrimaryLiteral}SOME${_ansiErr}, указанного первоначально${_ansiReset}"
    --return "1"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --opt-b/b=some \
      -- \
      --optB SOME \
      -b thing \
    "
  '

# проверка краткой опции
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает значение для опции ${_ansiCmd}--opt-b${_ansiErr}, поэтому её краткая форма ${_ansiCmd}-b${_ansiErr} не может быть использоваана в ${_ansiOutline}объединении опций ${_ansiCmd}-ab${_ansiReset}"
    --return "1"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --opt-a/a \
      --opt-b/b=some \
      -- \
      -ab \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает краткую опцию ${_ansiPrimaryLiteral}c${_ansiErr} в ${_ansiOutline}объединении опций ${_ansiCmd}-abc${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      --opt-a/a \
      --opt-b/b \
      -- \
      -abc \
    "
  '

# проверка превышения числа аргументов в цикле разбора параметров
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает ни одного аргумента, но обнаружен: ${_ansiCmd}\"f g\"${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc -- '$_stq'f g'$_stq' '$_stqq'h i'$_stqq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не более 1 аргумента, но обнаружен 2-й: ${_ansiCmd}\"h i\"${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc a:? -- '$_stq'f g'$_stq' '$_stqq'h i'$_stqq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не более 2 аргументов, но обнаружен 3-й: ${_ansiCmd}e${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc a b:? -- c d e '$_stq'f g'$_stq' '$_stqq'h i'$_stqq'"
  '

# валидация аргументов в цикле разбора параметров
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}c${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"c d\"${_ansiErr} в качестве 1-го аргумента ${_ansiOutline}a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      -- \
      c \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве 2-го аргумента ${_ansiOutline}b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'b:1..'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве 3-го аргумента ${_ansiOutline}c${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'b:1..'$_stq' \
      '$_stq'c:0..'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      1 \
      -1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает ${_ansiOutline}непустое${_ansiErr} значение в качестве 4-го аргумента ${_ansiOutline}d${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'b:1..'$_stq' \
      '$_stq'c:0..'$_stq' \
      '$_stq'd'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      1 \
      0 \
      '$_stqq''$_stqq' \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}c${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"c d\"${_ansiErr} в качестве 3-го аргумента ${_ansiOutline}a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      -- \
      b \
      '$_stqq'c d'$_stqq' \
      c \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве 4-го аргумента ${_ansiOutline}b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'@b:1..'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      2 \
      1 \
      0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве 4-го аргумента ${_ansiOutline}c${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'b:1..'$_stq' \
      '$_stq'@c:0..'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      1 \
      0 \
      -1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает ${_ansiOutline}непустое${_ansiErr} значение в качестве 6-го аргумента ${_ansiOutline}d${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq' \
      '$_stq'b:1..'$_stq' \
      '$_stq'c:0..'$_stq' \
      '$_stq'@d'$_stq' \
      -- \
      '$_stq'c d'$_stq' \
      1 \
      0 \
      a \
      b \
      '$_stqq''$_stqq' \
    "
  '

# ожидание команды
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --before "local -a _prepareCodeToParseFuncParams2TestFunc_alphaShortcuts=( alef a )"
    --before "local -a _prepareCodeToParseFuncParams2TestFunc_bravoCharlieShortcuts=( beta-gamma b )"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeToParseFuncParams2TestFunc_alpha _prepareCodeToParseFuncParams2TestFunc_bravoCharlie"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}charlie${_ansiErr} ожидает одну из следующих команд: ${_ansiSecondaryLiteral}alpha alef a bravo-charlie beta-gamma b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      --isCommandWrapper \
      -- \
      charlie \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --before "local -a _prepareCodeToParseFuncParams2TestFunc_alphaShortcuts=( alef a )"
    --before "local -a _prepareCodeToParseFuncParams2TestFunc_bravoCharlieShortcuts=( beta-gamma b )"
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "eval \"_prepareCodeToParseFuncParams2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeToParseFuncParams2TestFunc_alpha _prepareCodeToParseFuncParams2TestFunc_bravoCharlie"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} в качестве первого аргумента ожидает одну из следующих команд: ${_ansiSecondaryLiteral}alpha alef a bravo-charlie beta-gamma b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      --isCommandWrapper \
    "
  '

# ожидание аргумента
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает 1-й аргумент ${_ansiOutline}a${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc a -- "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает 2-й аргумент ${_ansiOutline}b${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc a b -- c"
  '

# _postProcessVarNames

  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не менее ${_ansiPrimaryLiteral}2${_ansiErr} элементов в списке значений аргумента ${_ansiOutline}args${_ansiErr}, но не получено ничего${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      a '$_stq'b:'$_stOpenBraceInQ'x y'$_stCloseBraceInQ''$_stq' d:? @2..args \
      -- \
      '$_stq'f g'$_stq' '$_stqq'x'$_stqq' \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не менее ${_ansiPrimaryLiteral}2${_ansiErr} элементов в списке значений аргумента ${_ansiOutline}args${_ansiErr}, но получен ${_ansiSecondaryLiteral}1${_ansiErr}: ${_ansiCmd}itemValue0${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      a '$_stq'b:'$_stOpenBraceInQ'x y'$_stCloseBraceInQ''$_stq' d:? @2..args \
      -- \
      '$_stq'f g'$_stq' '$_stqq'x'$_stqq' valueOfDArg itemValue0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не менее ${_ansiPrimaryLiteral}3${_ansiErr} элементов в списке значений аргумента ${_ansiOutline}args${_ansiErr}, но получены ${_ansiSecondaryLiteral}2${_ansiErr}: ${_ansiCmd}itemValue0 itemValue1${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      a '$_stq'b:'$_stOpenBraceInQ'x y'$_stCloseBraceInQ''$_stq' d:? @3..args \
      -- \
      '$_stq'f g'$_stq' '$_stqq'x'$_stqq' valueOfDArg itemValue0 itemValue1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} ожидает не более ${_ansiPrimaryLiteral}2${_ansiErr} элементов в списке значений аргумента ${_ansiOutline}args${_ansiErr}, но получены ${_ansiSecondaryLiteral}3${_ansiErr}: ${_ansiCmd}itemValue0 itemValue1 itemValue2${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc \
      a '$_stq'b:'$_stOpenBraceInQ'x y'$_stCloseBraceInQ''$_stq' d:? @..2args \
      -- \
      '$_stq'f g'$_stq' '$_stqq'x'$_stqq' valueOfDArg itemValue0 itemValue1 itemValue2 \
    "
  '


  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}<empty>
      ${_ansiOutline}optB${_ansiReset}<empty>
      ${_ansiOutline}arg${_ansiReset}: ${_ansiPrimaryLiteral}-abc${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --treatUnknownOptionAsArg \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        optA \
        optB \
        arg \
      '$_stq' \
      --opt-a/a \
      --opt-b/b \
      arg \
      -- \
      -abc \
    "
  '

  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "0"
    --stderr "
      extraParam
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}thing${_ansiReset}
      ${_ansiOutline}argB${_ansiReset}(4): ${_ansiSecondaryLiteral}h f d b${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debug --clean '$_stDollarInQ'@'$_stq' \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' argA argB'$_stq' \
      argA=some \
      '$_stq'@2..argB:'$_stOpenBraceInQ'i h g f e d c b a'$_stCloseBraceInQ':sorted'$_stq' \
      -- \
      thing b d f h \
      -- \
      extraParam \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "0"
    --stderr "
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}thing${_ansiReset}
      ${_ansiOutline}argB${_ansiReset}(6): ${_ansiSecondaryLiteral}1 1 2 2 3 3${_ansiReset}
      extraParam
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' argA argB'$_stq' \
      --run '$_stq'_debug --clean '$_stDollarInQ'@'$_stq' \
      argA=some @2..argB:1..:sorted \
      -- \
      thing 3 1 2 1 2 3 \
      -- \
      extraParam \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "0"
    --stderr "
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}thing${_ansiReset}
      ${_ansiOutline}argB${_ansiReset}(3): ${_ansiSecondaryLiteral}3 1 2${_ansiReset}
      extraParam
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' argA argB'$_stq' \
      --run '$_stq'_debug --clean '$_stDollarInQ'@'$_stq' \
      argA=some @2..argB:1..:unique \
      -- \
      thing 3 1 3 1 2 2 \
      -- \
      extraParam \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "0"
    --stderr "
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}thing${_ansiReset}
      ${_ansiOutline}argB${_ansiReset}(6): ${_ansiSecondaryLiteral}1 2 3 3 1 2${_ansiReset}
      extraParam
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' argA argB'$_stq' \
      --run '$_stq'_debug --clean '$_stDollarInQ'@'$_stq' \
      argA=some @2..argB:1.. \
      -- \
      thing 1 2 3 3 1 2 \
      -- \
      extraParam \
    "
  '

#!===


  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}<empty>
      ${_ansiOutline}OPT_a${_ansiReset}(2): ${_ansiSecondaryLiteral}--a '$_stqq$_stqq'${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          a \
          OPT_a \
      '$_stq' \
      --a:? \
      -- \
      --a '$_stqq$_stqq' \
    "
  '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   "--stderr=
  #     ${_ansiOutline}a${_ansiReset}<empty>
  #     ${_ansiOutline}OPT_a${_ansiReset}(2): ${_ansiSecondaryLiteral}--a '$_stqq$_stqq'${_ansiReset}
  #   "
  #   --stderrParaWithIndent "0"
  #   "_prepareCodeToParseFuncParams2TestFunc
  #     --run '$_stq'
  #       _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq'
  #         a
  #         OPT_a
  #     '$_stq'
  #     --
  #     --a:!= --'$_stqq$_stqq'
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}c${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"c d\"${_ansiErr} для опции ${_ansiCmd}--a${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --
  #     '$_stq'--a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq'
  #     --
  #     --a=c
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} для опции ${_ansiCmd}--a${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --
  #     '$_stq'--a:1..'$_stq'
  #     --
  #     --a=0
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} для опции ${_ansiCmd}--a${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --
  #     '$_stq'--a:0..'$_stq'
  #     --
  #     --a=-1
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}c${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"c d\"${_ansiErr} в качестве 1-го аргумента ${_ansiOutline}a${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --
  #     '$_stq'@a:'$_stOpenBraceInQ'b '$_stqq'c d'$_stqq''$_stCloseBraceInQ''$_stq'
  #     --
  #     c
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве 2-го аргумента ${_ansiOutline}a${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc --
  #     '$_stq'@a:0..'$_stq'
  #     --
  #     0 -1
  #   "
  # '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}boolOptA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_boolOptA${_ansiReset}(1): ${_ansiSecondaryLiteral}--boolOptA${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          boolOptA \
          OPT_boolOptA \
      '$_stq' \
      --bool-opt-a \
      -- \
      --bool-opt-a \
      --boolOptA \
    "
  '
  # ${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не ожидает, что опция ${_ansiCmd}--boolOptA${_ansiErr} будет указана повторно${_ansiReset}
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "0"
  #   "--stderr=
  #     ${_ansiOutline}a${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
  #   "
  #   --stderrParaWithIndent "0"
  #   "_prepareCodeToParseFuncParams2TestFunc
  #     --run '$_stq'
  #       _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq'
  #         a
  #     '$_stq'
  #     --
  #     --a
  #     --
  #     --a
  #     --a true
  #   "
  # '
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --return "1"
  #   "--stderr=${_ansiErr}ERR: Значение по умолчанию (задаваемое после ${_ansiSecondaryLiteral}=${_ansiErr}) в определении аргумента ${_ansiCmd}arg:1..=${_ansiErr} функции ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} не может быть пустым, в противном случае аргумент должен быть помечен суффиксом ${_ansiPrimaryLiteral}?${_ansiErr} как ${_ansiOutline}опциональный${_ansiErr}: ${_ansiCmd}arg:1..?${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc -- arg:1..="
  # '


  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(4): ${_ansiSecondaryLiteral}100 20 -3 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 -a -3 --a 20 -a -3${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
        --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:sorted \
      -- \
      -a 100 \
      -a -3 \
      --a 20 \
      -a -3 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(4): ${_ansiSecondaryLiteral}-3 20 20 100${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 --a 20${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:sorted:.. \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      --a 20 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(4): ${_ansiSecondaryLiteral}20 100 100 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 -a 100${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      '$_stq'@--a/a:sorted:'$_stOpenBraceInQ'20 100 -3'$_stCloseBraceInQ$_stq' \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      -a 100 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}100 20 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 -a -3 --a 20 -a -3${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:sorted:unique \
      -- \
      -a 100 \
      -a -3 \
      --a 20 \
      -a -3 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}-3 20 100${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 --a 20${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:sorted:..:unique \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      --a 20 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}20 100 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 -a 100${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      '$_stq'@--a/a:sorted:unique:'$_stOpenBraceInQ'20 100 -3'$_stCloseBraceInQ$_stq' \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      -a 100 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}100 -3 20${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 -a -3 --a 20 -a -3${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:unique \
      -- \
      -a 100 \
      -a -3 \
      --a 20 \
      -a -3 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}100 20 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 --a 20${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      @--a/a:..:unique \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      --a 20 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}a${_ansiReset}(3): ${_ansiSecondaryLiteral}100 20 -3${_ansiReset}
      ${_ansiOutline}OPT_a${_ansiReset}(8): ${_ansiSecondaryLiteral}-a 100 --a 20 -a -3 -a 100${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
        a \
        OPT_a \
      '$_stq' \
      '$_stq'@--a/a:unique:'$_stOpenBraceInQ'20 100 -3'$_stCloseBraceInQ$_stq' \
      -- \
      -a 100 \
      --a 20 \
      -a -3 \
      -a 100 \
    "
  '

  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}d${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"d e\"${_ansiErr} в качестве вычисляемого значения по умолчанию параметра ${_ansiCmd}a:(b \"d e\")=\$defaultValueForArgA${_ansiReset}"
    --before "local defaultValueForArgA=d"
    "_prepareCodeToParseFuncParams2TestFunc '$_stq'a:'$_stOpenBraceInQ'b '$_stqq'd e'$_stqq''$_stCloseBraceInQ'='$_stDollarInQ'defaultValueForArgA'$_stq'"
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}d${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}b \"d e\"${_ansiErr} в качестве элемента вычисляемого значения по умолчанию параметра ${_ansiCmd}@--a:(b \"d e\")=( \${defaultValueForOptA[@]} )${_ansiReset}"
    --before "local -a defaultValueForOptA=( d e )"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq'_debugVar a'$_stq' \
      '$_stq'@--a:'$_stOpenBraceInQ'b '$_stqq'd e'$_stqq''$_stCloseBraceInQ'='$_stOpenBraceInQ' '$_stDollarInQ'{defaultValueForOptA[@]} '$_stCloseBraceInQ''$_stq' \
    "
  '

  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}<empty>
      ${_ansiOutline}optA2${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optA3${_ansiReset}: ${_ansiPrimaryLiteral}defaultValueOfA${_ansiReset}
      ${_ansiOutline}optA4${_ansiReset}: ${_ansiPrimaryLiteral}customValueOfA4${_ansiReset}
      ${_ansiOutline}optA5${_ansiReset}: ${_ansiPrimaryLiteral}calculatedDefaultValueOfA${_ansiReset}
      ${_ansiOutline}optA6${_ansiReset}: ${_ansiPrimaryLiteral}customValueOfA6${_ansiReset}
      ${_ansiOutline}optA7${_ansiReset}: ${_ansiPrimaryLiteral}echoOfDefaultValueOfA${_ansiReset}
      ${_ansiOutline}optA8${_ansiReset}: ${_ansiPrimaryLiteral}customValueOfA8${_ansiReset}
      ${_ansiOutline}optAa${_ansiReset}(0)
      ${_ansiOutline}optAa2${_ansiReset}(2): ${_ansiSecondaryLiteral}firstItemOfCustomValueOfA2 secondItemOfCustomValueOfA2${_ansiReset}
      ${_ansiOutline}optAa3${_ansiReset}(2): ${_ansiSecondaryLiteral}firstItemOfDefaultValueOfA secondItemOfDefaultValueOfA${_ansiReset}
      ${_ansiOutline}optAa4${_ansiReset}(2): ${_ansiSecondaryLiteral}firstItemOfCustomValueOfA4 secondItemOfCustomValueOfA4${_ansiReset}
      ${_ansiOutline}optAa5${_ansiReset}(2): ${_ansiSecondaryLiteral}calculatedFirstItemOfDefaultValueOfA calculatedSecondItemOfDefaultValueOfA${_ansiReset}
      ${_ansiOutline}optAa6${_ansiReset}(2): ${_ansiSecondaryLiteral}firstItemOfCustomValueOfA6 secondItemOfCustomValueOfA6${_ansiReset}
      ${_ansiOutline}optAa7${_ansiReset}(2): ${_ansiSecondaryLiteral}echoOfFirstItemOfDefaultValueOfA echoOfSecondItemOfDefaultValueOfA${_ansiReset}
      ${_ansiOutline}optAa8${_ansiReset}(2): ${_ansiSecondaryLiteral}firstItemOfCustomValueOfA8 secondItemOfCustomValueOfA8${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}<empty>
      ${_ansiOutline}optB2${_ansiReset}: ${_ansiPrimaryLiteral}0${_ansiReset}
      ${_ansiOutline}optB3${_ansiReset}: ${_ansiPrimaryLiteral}0${_ansiReset}
      ${_ansiOutline}optB4${_ansiReset}: ${_ansiPrimaryLiteral}1${_ansiReset}
      ${_ansiOutline}optB5${_ansiReset}: ${_ansiPrimaryLiteral}1${_ansiReset}
      ${_ansiOutline}optB6${_ansiReset}: ${_ansiPrimaryLiteral}2${_ansiReset}
      ${_ansiOutline}optB7${_ansiReset}: ${_ansiPrimaryLiteral}2${_ansiReset}
      ${_ansiOutline}optB8${_ansiReset}: ${_ansiPrimaryLiteral}3${_ansiReset}
      ${_ansiOutline}optBb${_ansiReset}(0)
      ${_ansiOutline}optBb2${_ansiReset}(2): ${_ansiSecondaryLiteral}5 6${_ansiReset}
      ${_ansiOutline}optBb3${_ansiReset}(2): ${_ansiSecondaryLiteral}0 1${_ansiReset}
      ${_ansiOutline}optBb4${_ansiReset}(2): ${_ansiSecondaryLiteral}7 8${_ansiReset}
      ${_ansiOutline}optBb5${_ansiReset}(2): ${_ansiSecondaryLiteral}2 3${_ansiReset}
      ${_ansiOutline}optBb6${_ansiReset}(2): ${_ansiSecondaryLiteral}9 10${_ansiReset}
      ${_ansiOutline}optBb7${_ansiReset}(2): ${_ansiSecondaryLiteral}4 5${_ansiReset}
      ${_ansiOutline}optBb8${_ansiReset}(2): ${_ansiSecondaryLiteral}11 12${_ansiReset}
      ${_ansiOutline}optC${_ansiReset}<empty>
      ${_ansiOutline}optC2${_ansiReset}: ${_ansiPrimaryLiteral}1${_ansiReset}
      ${_ansiOutline}optC3${_ansiReset}: ${_ansiPrimaryLiteral}1${_ansiReset}
      ${_ansiOutline}optC4${_ansiReset}: ${_ansiPrimaryLiteral}2${_ansiReset}
      ${_ansiOutline}optC5${_ansiReset}: ${_ansiPrimaryLiteral}3${_ansiReset}
      ${_ansiOutline}optC6${_ansiReset}: ${_ansiPrimaryLiteral}3${_ansiReset}
      ${_ansiOutline}optC7${_ansiReset}: ${_ansiPrimaryLiteral}2${_ansiReset}
      ${_ansiOutline}optC8${_ansiReset}: ${_ansiPrimaryLiteral}4${_ansiReset}
      ${_ansiOutline}optCc${_ansiReset}(0)
      ${_ansiOutline}optCc2${_ansiReset}(2): ${_ansiSecondaryLiteral}5 6${_ansiReset}
      ${_ansiOutline}optCc3${_ansiReset}(2): ${_ansiSecondaryLiteral}1 2${_ansiReset}
      ${_ansiOutline}optCc4${_ansiReset}(2): ${_ansiSecondaryLiteral}7 8${_ansiReset}
      ${_ansiOutline}optCc5${_ansiReset}(2): ${_ansiSecondaryLiteral}4 5${_ansiReset}
      ${_ansiOutline}optCc6${_ansiReset}(2): ${_ansiSecondaryLiteral}9 10${_ansiReset}
      ${_ansiOutline}optCc7${_ansiReset}(2): ${_ansiSecondaryLiteral}5 6${_ansiReset}
      ${_ansiOutline}optCc8${_ansiReset}(2): ${_ansiSecondaryLiteral}11 12${_ansiReset}
      ${_ansiOutline}optD${_ansiReset}<empty>
      ${_ansiOutline}optD2${_ansiReset}: ${_ansiPrimaryLiteral}g${_ansiReset}
      ${_ansiOutline}optD3${_ansiReset}: ${_ansiPrimaryLiteral}f${_ansiReset}
      ${_ansiOutline}optD4${_ansiReset}: ${_ansiPrimaryLiteral}e${_ansiReset}
      ${_ansiOutline}optD5${_ansiReset}: ${_ansiPrimaryLiteral}g${_ansiReset}
      ${_ansiOutline}optD6${_ansiReset}: ${_ansiPrimaryLiteral}f${_ansiReset}
      ${_ansiOutline}optD7${_ansiReset}: ${_ansiPrimaryLiteral}e${_ansiReset}
      ${_ansiOutline}optD8${_ansiReset}: ${_ansiPrimaryLiteral}g${_ansiReset}
      ${_ansiOutline}optDd${_ansiReset}(0)
      ${_ansiOutline}optDd2${_ansiReset}(2): ${_ansiSecondaryLiteral}e f${_ansiReset}
      ${_ansiOutline}optDd3${_ansiReset}(2): ${_ansiSecondaryLiteral}e g${_ansiReset}
      ${_ansiOutline}optDd4${_ansiReset}(2): ${_ansiSecondaryLiteral}e f${_ansiReset}
      ${_ansiOutline}optDd5${_ansiReset}(2): ${_ansiSecondaryLiteral}g f${_ansiReset}
      ${_ansiOutline}optDd6${_ansiReset}(2): ${_ansiSecondaryLiteral}e f${_ansiReset}
      ${_ansiOutline}optDd7${_ansiReset}(2): ${_ansiSecondaryLiteral}f e${_ansiReset}
      ${_ansiOutline}optDd8${_ansiReset}(2): ${_ansiSecondaryLiteral}g f${_ansiReset}
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}valueOfArgA${_ansiReset}
      ${_ansiOutline}argB${_ansiReset}: ${_ansiPrimaryLiteral}0${_ansiReset}
      ${_ansiOutline}argC${_ansiReset}: ${_ansiPrimaryLiteral}1${_ansiReset}
      ${_ansiOutline}argD${_ansiReset}: ${_ansiPrimaryLiteral}f${_ansiReset}
    "
    --return "0"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfA=calculatedDefaultValueOfA"
    --before "local __firstItemOfDefaultValueOfA=calculatedFirstItemOfDefaultValueOfA"
    --before "local __secondItemOfDefaultValueOfA=calculatedSecondItemOfDefaultValueOfA"
    --before "local __defaultValueOfB=1"
    --before "local __firstItemOfDefaultValueOfB=2"
    --before "local __secondItemOfDefaultValueOfB=3"
    --before "local __defaultValueOfC=3"
    --before "local __firstItemOfDefaultValueOfC=4"
    --before "local __secondItemOfDefaultValueOfC=5"
    --before "local __defaultValueOfD=g"
    --before "local __firstItemOfDefaultValueOfD=g"
    --before "local __secondItemOfDefaultValueOfD=f"

    --before "local __defaultValueOfOptE=devaultValueOfOptE"
    --before "local __partOfDefaultValueOfOptD=DevaultValueOfOptD"
    --before "local __partOfDefaultValueOfOptF=DevaultValueOfOptF"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optA2 optA3 optA4 optA5 optA6 optA7 optA8 \
          optAa optAa2 optAa3 optAa4 optAa5 optAa6 optAa7 optAa8 \
          optB optB2 optB3 optB4 optB5 optB6 optB7 optB8 \
          optBb optBb2 optBb3 optBb4 optBb5 optBb6 optBb7 optBb8 \
          optC optC2 optC3 optC4 optC5 optC6 optC7 optC8 \
          optCc optCc2 optCc3 optCc4 optCc5 optCc6 optCc7 optCc8 \
          optD optD2 optD3 optD4 optD5 optD6 optD7 optD8 \
          optDd optDd2 optDd3 optDd4 optDd5 optDd6 optDd7 optDd8 \
          argA argB argC argD \
      '$_stq' \
      --canBeMixedOptionsAndArgs \
      '$_stq'--opt-a'$_stq' \
      '$_stq'--opt-a2'$_stq' \
      '$_stq'--opt-a3=defaultValueOfA'$_stq' \
      '$_stq'--opt-a4=defaultValueOfA'$_stq' \
      '$_stq'--opt-a5='$_stDollarInQ'__defaultValueOfA'$_stq' \
      '$_stq'--opt-a6='$_stDollarInQ'__defaultValueOfA'$_stq' \
      '$_stq'--opt-a7='$_stDollarInQ''$_stOpenBraceInQ'echo echoOfDefaultValueOfA'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-a8='$_stDollarInQ''$_stOpenBraceInQ'echo echoOfDefaultValueOfA'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa'$_stq' \
      '$_stq'@--opt-aa2'$_stq' \
      '$_stq'@--opt-aa3='$_stOpenBraceInQ' firstItemOfDefaultValueOfA secondItemOfDefaultValueOfA '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa4='$_stOpenBraceInQ' firstItemOfDefaultValueOfA secondItemOfDefaultValueOfA '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa5='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfA '$_stDollarInQ'__secondItemOfDefaultValueOfA '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa6='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfA '$_stDollarInQ'__secondItemOfDefaultValueOfA '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa7='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo echoOfFirstItemOfDefaultValueOfA echoOfSecondItemOfDefaultValueOfA'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-aa8='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo echoOfFirstItemOfDefaultValueOfA echoOfSecondItemOfDefaultValueOfA'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
 \
      '$_stq'--opt-b:0..'$_stq' \
      '$_stq'--opt-b2:0..'$_stq' \
      '$_stq'--opt-b3:0..=0'$_stq' \
      '$_stq'--opt-b4:0..=0'$_stq' \
      '$_stq'--opt-b5:0..='$_stDollarInQ'__defaultValueOfB'$_stq' \
      '$_stq'--opt-b6:0..='$_stDollarInQ'__defaultValueOfB'$_stq' \
      '$_stq'--opt-b7:0..='$_stDollarInQ''$_stOpenBraceInQ'echo 2'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-b8:0..='$_stDollarInQ''$_stOpenBraceInQ'echo 2'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb:0..'$_stq' \
      '$_stq'@--opt-bb2:0..'$_stq' \
      '$_stq'@--opt-bb3:0..='$_stOpenBraceInQ' 0 1 '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb4:0..='$_stOpenBraceInQ' 0 1 '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb5:0..='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfB '$_stDollarInQ'__secondItemOfDefaultValueOfB '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb6:0..='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfB '$_stDollarInQ'__secondItemOfDefaultValueOfB '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb7:0..='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo 4 5'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-bb8:0..='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo 4 5'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
 \
      '$_stq'--opt-c:1..'$_stq' \
      '$_stq'--opt-c2:1..'$_stq' \
      '$_stq'--opt-c3:1..=1'$_stq' \
      '$_stq'--opt-c4:1..=1'$_stq' \
      '$_stq'--opt-c5:1..='$_stDollarInQ'__defaultValueOfC'$_stq' \
      '$_stq'--opt-c6:1..='$_stDollarInQ'__defaultValueOfC'$_stq' \
      '$_stq'--opt-c7:1..='$_stDollarInQ''$_stOpenBraceInQ'echo 2'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-c8:1..='$_stDollarInQ''$_stOpenBraceInQ'echo 2'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc:1..'$_stq' \
      '$_stq'@--opt-cc2:1..'$_stq' \
      '$_stq'@--opt-cc3:1..='$_stOpenBraceInQ' 1 2 '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc4:1..='$_stOpenBraceInQ' 1 2 '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc5:1..='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfC '$_stDollarInQ'__secondItemOfDefaultValueOfC '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc6:1..='$_stOpenBraceInQ' '$_stDollarInQ'__firstItemOfDefaultValueOfC '$_stDollarInQ'__secondItemOfDefaultValueOfC '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc7:1..='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo 5 6'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-cc8:1..='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo 5 6'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
 \
      '$_stq'--opt-d:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-d2:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-d3:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'=f'$_stq' \
      '$_stq'--opt-d4:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'=f'$_stq' \
      '$_stq'--opt-d5:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stDollarInQ'__defaultValueOfD'$_stq' \
      '$_stq'--opt-d6:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stDollarInQ'__defaultValueOfD'$_stq' \
      '$_stq'--opt-d7:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stDollarInQ''$_stOpenBraceInQ'echo e'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-d8:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stDollarInQ''$_stOpenBraceInQ'echo e'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd2:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd3:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' e g '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd4:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' e g '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd5:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' '$_stqq$_stDollarInQ'__firstItemOfDefaultValueOfD'$_stqq' '$_stDollarInQ'__secondItemOfDefaultValueOfD '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd6:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' '$_stqq$_stDollarInQ'__firstItemOfDefaultValueOfD'$_stqq' '$_stDollarInQ'__secondItemOfDefaultValueOfD '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd7:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo f e'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
      '$_stq'@--opt-dd8:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ'='$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'echo f e'$_stCloseBraceInQ' '$_stCloseBraceInQ''$_stq' \
      '$_stq'argA'$_stq' \
      '$_stq'argB:0..'$_stq' \
      '$_stq'argC:1..'$_stq' \
      '$_stq'argD:'$_stOpenBraceInQ'e f g'$_stCloseBraceInQ''$_stq' \
      -- \
      --opt-a2 \
      --opt-a4 customValueOfA4 \
      --opt-a6 customValueOfA6 \
      --opt-a8 customValueOfA8 \
      --opt-aa2 firstItemOfCustomValueOfA2 \
      --opt-aa2 secondItemOfCustomValueOfA2 \
      --opt-aa4 firstItemOfCustomValueOfA4 \
      --opt-aa4 secondItemOfCustomValueOfA4 \
      --opt-aa6 firstItemOfCustomValueOfA6 \
      --opt-aa6 secondItemOfCustomValueOfA6 \
      --opt-aa8 firstItemOfCustomValueOfA8 \
      --opt-aa8 secondItemOfCustomValueOfA8 \
 \
      --opt-b2 0 \
      --opt-b4 1 \
      --opt-b6 2 \
      --opt-b8 3 \
      --opt-bb2 5 \
      --opt-bb2 6 \
      --opt-bb4 7 \
      --opt-bb4 8 \
      --opt-bb6 9 \
      --opt-bb6 10 \
      --opt-bb8 11 \
      --opt-bb8 12 \
 \
      --opt-c2 1 \
      --opt-c4 2 \
      --opt-c6 3 \
      --opt-c8 4 \
      --opt-cc2 5 \
      --opt-cc2 6 \
      --opt-cc4 7 \
      --opt-cc4 8 \
      --opt-cc6 9 \
      --opt-cc6 10 \
      --opt-cc8 11 \
      --opt-cc8 12 \
 \
      --opt-d2 g \
      --opt-d4 e \
      --opt-d6 f \
      --opt-d8 g \
      --opt-dd2 e \
      --opt-dd2 f \
      --opt-dd4 e \
      --opt-dd4 f \
      --opt-dd6 e \
      --opt-dd6 f \
      --opt-dd8 g \
      --opt-dd8 f \
 \
      valueOfArgA \
      0 \
      1 \
      f \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}--unknownOption${_ansiReset}
    "
    --return "0"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA argA \
      '$_stq' \
      --treatUnknownOptionAsArg \
      --opt-a \
      '$_stq'argA'$_stq' \
      -- \
      --optA \
      --unknownOption \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}--unknownOption${_ansiReset}
    "
    --return "0"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA argA \
      '$_stq' \
      --canBeMixedOptionsAndArgs \
      --treatUnknownOptionAsArg \
      --opt-a \
      '$_stq'argA'$_stq' \
      -- \
      --optA \
      --unknownOption \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}argA${_ansiReset}: ${_ansiPrimaryLiteral}--unknownOption${_ansiReset}
    "
    --return "0"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA argA \
      '$_stq' \
      --canBeMixedOptionsAndArgs \
      --treatUnknownOptionAsArg \
      --opt-a \
      '$_stq'argA'$_stq' \
      -- \
      --unknownOption \
      --optA \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}: ${_ansiPrimaryLiteral}\"some thing\"${_ansiReset}
      ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optD${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optE${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optF${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optG${_ansiReset}(3): ${_ansiSecondaryLiteral}\"g value\" g2 g3${_ansiReset}
      ${_ansiOutline}optH${_ansiReset}: ${_ansiPrimaryLiteral}defaultValueOfOptH${_ansiReset}
      ${_ansiOutline}optI${_ansiReset}(3): ${_ansiSecondaryLiteral}\"iVal 0\" iVal1 iVal2${_ansiReset}
      ${_ansiOutline}optJ${_ansiReset}: ${_ansiPrimaryLiteral}\"optJ default value\"${_ansiReset}
      ${_ansiOutline}optK${_ansiReset}(2): ${_ansiSecondaryLiteral}\"optK default value\" kVal1${_ansiReset}
      ${_ansiOutline}optL${_ansiReset}<empty>
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(2): ${_ansiSecondaryLiteral}-b \"some thing\"${_ansiReset}
      ${_ansiOutline}OPT_optC${_ansiReset}(1): ${_ansiSecondaryLiteral}-c${_ansiReset}
      ${_ansiOutline}OPT_optD${_ansiReset}(1): ${_ansiSecondaryLiteral}-d${_ansiReset}
      ${_ansiOutline}OPT_optE${_ansiReset}(1): ${_ansiSecondaryLiteral}--optE${_ansiReset}
      ${_ansiOutline}OPT_optF${_ansiReset}(1): ${_ansiSecondaryLiteral}--opt-f${_ansiReset}
      ${_ansiOutline}OPT_optG${_ansiReset}(6): ${_ansiSecondaryLiteral}--opt-g \"g value\" --optG g2 -g g3${_ansiReset}
      ${_ansiOutline}OPT_optH${_ansiReset}(2): ${_ansiSecondaryLiteral}--opt-h defaultValueOfOptH${_ansiReset}
      ${_ansiOutline}OPT_optI${_ansiReset}(6): ${_ansiSecondaryLiteral}--opt-i \"iVal 0\" --opt-i iVal1 --opt-i iVal2${_ansiReset}
      ${_ansiOutline}OPT_optJ${_ansiReset}(2): ${_ansiSecondaryLiteral}--opt-j \"optJ default value\"${_ansiReset}
      ${_ansiOutline}OPT_optK${_ansiReset}(4): ${_ansiSecondaryLiteral}--opt-k \"optK default value\" --opt-k kVal1${_ansiReset}
      ${_ansiOutline}OPT_optL${_ansiReset}(0)
    "
    --stderrParaWithIndent "0"
    --before "local defaultValueOfOptJ=\"optJ default value\""
    --before "local defaultValueOfOptK=\"optK default value\""
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC optD optE optF optG optH optI optJ optK optL \
          OPT_optA OPT_optB OPT_optC OPT_optD OPT_optE OPT_optF OPT_optG OPT_optH OPT_optI OPT_optJ OPT_optK OPT_optL \
      '$_stq' \
      --opt-a/a \
      --opt-b/b=some \
      --opt-c/c \
      --opt-d/d \
      --opt-e/e \
      --opt-f/f \
      @--opt-g/g \
      --opt-h/H=defaultValueOfOptH \
      '$_stq'@--opt-i='$_stOpenBraceInQ$_stqq'iVal 0'$_stqq' iVal1 iVal2'$_stCloseBraceInQ''$_stq' \
      '$_stq'--opt-j='$_stDollarInQ'defaultValueOfOptJ'$_stq' \
      '$_stq'@--opt-k='$_stOpenBraceInQ' '$_stqq''$_stDollarInQ'defaultValueOfOptK'$_stqq' kVal1 '$_stCloseBraceInQ$_stq' \
      --opt-l/l \
      -- \
      -b '$_stq'some thing'$_stq' \
      -ac \
      -d \
      --optE \
      --opt-f \
      --opt-g '$_stq'g value'$_stq' \
      --optG '$_stq'g2'$_stq' \
      -g '$_stq'g3'$_stq' \
    "
  '


  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(1): ${_ansiSecondaryLiteral}-b${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB \
          OPT_optA OPT_optB \
      '$_stq' \
      --opt-a/a \
      --opt-b/b \
      -- \
      --optA \
      -ab \
    "
  '

  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}valueC${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}valueA valueB${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:'$_stOpenBraceInQ'valueA valueB'$_stCloseBraceInQ$_stq' \
      -- \
      -a valueC \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}valueC${_ansiErr} ожидает одно из следующих значений: ${_ansiSecondaryLiteral}valueA valueB${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:'$_stOpenBraceInQ'valueA valueB'$_stCloseBraceInQ$_stq' \
      -- \
      --opt-a valueA -a valueC \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}valueC${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:..'$_stq' \
      -- \
      -a valueC \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}valueC${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:..'$_stq' \
      -- \
      --opt-a 0 -a valueC \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:1..'$_stq' \
      -- \
      -a 0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}положительное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:1..'$_stq' \
      -- \
      --opt-a 2 -a 0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:0..'$_stq' \
      -- \
      -a -1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}-1${_ansiErr} ожидает ${_ansiOutline}неотрицательное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:0..'$_stq' \
      -- \
      --opt-a 0 -a -1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}1${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} не менее ${_ansiPrimaryLiteral}2${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:2..'$_stq' \
      -- \
      -a 1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}1${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} не менее ${_ansiPrimaryLiteral}2${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:2..'$_stq' \
      -- \
      --opt-a 3 -a 1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}отрицательное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:..-1'$_stq' \
      -- \
      -a 0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}0${_ansiErr} ожидает ${_ansiOutline}отрицательное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:..-1'$_stq' \
      -- \
      --opt-a -2 -a 0 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}1${_ansiErr} ожидает ${_ansiOutline}неположительное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:..0'$_stq' \
      -- \
      -a 1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}1${_ansiErr} ожидает ${_ansiOutline}неположительное целое число${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:..0'$_stq' \
      -- \
      --opt-a 0 -a 1 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}3${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} не более ${_ansiPrimaryLiteral}2${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:..2'$_stq' \
      -- \
      -a 3 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}3${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} не более ${_ansiPrimaryLiteral}2${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:..2'$_stq' \
      -- \
      --opt-a 2 -a 3 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}2${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} из диапазона ${_ansiPrimaryLiteral}4..6${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'--opt-a/a:4..6'$_stq' \
      -- \
      -a 2 \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc${_ansiErr} вместо ${_ansiPrimaryLiteral}2${_ansiErr} ожидает ${_ansiOutline}целое число${_ansiErr} из диапазона ${_ansiPrimaryLiteral}4..6${_ansiErr} в качестве значения опции ${_ansiCmd}-a${_ansiReset}"
    --return "1"
    "_prepareCodeToParseFuncParams2TestFunc \
      '$_stq'@--opt-a/a:4..6'$_stq' \
      -- \
      --opt-a 4 -a 2 \
    "
  '

  # '
  #   --return "1"
  #   "--stderr:${_ansiErr}ERR: Похоже на то, что в ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc2_alpha${_ansiErr} пропущено ${_ansiCmd}eval \"\$_funcParams2\"${_ansiReset}"
  #   "_prepareCodeToParseFuncParams2TestFunc2 --isCommandWrapper"
  # '

    # --before "_substitute noStack true"
    # --after "_restore noStack"
  '
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}valueOfOptC${_ansiReset}
      ${_ansiOutline}optD${_ansiReset}(2): ${_ansiSecondaryLiteral}some thing${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(1): ${_ansiSecondaryLiteral}-b${_ansiReset}
      ${_ansiOutline}OPT_optC${_ansiReset}(2): ${_ansiSecondaryLiteral}--opt-c valueOfOptC${_ansiReset}
      ${_ansiOutline}OPT_optD${_ansiReset}(4): ${_ansiSecondaryLiteral}-d some --opt-d thing${_ansiReset}
    "
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC optD \
          OPT_optA OPT_optB OPT_optC OPT_optD \
      '$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      --opt-a/a \
      --opt-b/b \
      '$_stq'--opt-c/c:'$_stOpenBraceInQ' a b c valueOfOptC'$_stCloseBraceInQ'='$_stq' \
      @1..3--opt-d/d \
      -- \
      --optA --opt-a --opt-c valueOfOptC alpha -ab -d some --opt-d thing \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc3${_ansiErr} ожидает не менее ${_ansiPrimaryLiteral}1${_ansiErr} элемента в списке значений опции ${_ansiCmd}--opt-d${_ansiErr}, но не получено ничего${_ansiReset}"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC optD \
          OPT_optA OPT_optB OPT_optC OPT_optD \
      '$_stq' \
      --isCommandWrapper \
      --opt-a/a \
      --opt-b/b \
      '$_stq'--opt-c/c:'$_stOpenBraceInQ' a b c valueOfOptC'$_stCloseBraceInQ'='$_stq' \
      @1..3--opt-d/d \
      -- \
      --optA --opt-a --opt-c valueOfOptC alpha -ab -d some --opt-d thing \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc3 alpha${_ansiErr} не ожидает краткую опцию ${_ansiPrimaryLiteral}a${_ansiErr} в ${_ansiOutline}объединении опций ${_ansiCmd}-ab${_ansiReset}"
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC optD \
          OPT_optA OPT_optB OPT_optC OPT_optD \
      '$_stq' \
      --isCommandWrapper \
      --opt-a/a \
      --opt-b/b \
      '$_stq'--opt-c/c:'$_stOpenBraceInQ' a b c valueOfOptC'$_stCloseBraceInQ'='$_stq' \
      -- \
      --optA --opt-a --opt-c valueOfOptC alpha -ab \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}<empty>
      ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}valueOfOptC${_ansiReset}
      ${_ansiOutline}arg${_ansiReset}: ${_ansiPrimaryLiteral}-ab${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}--opt-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(0)
      ${_ansiOutline}OPT_optC${_ansiReset}(2): ${_ansiSecondaryLiteral}--opt-c valueOfOptC${_ansiReset}
    "
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC arg \
          OPT_optA OPT_optB OPT_optC \
      '$_stq' \
      --aOpt '$_stq' --treatUnknownOptionAsArg'$_stq' \
      --aParams '$_stq' arg'$_stq' \
      --isCommandWrapper \
      --opt-a/a \
      --opt-b/b \
      '$_stq'--opt-c/c:'$_stOpenBraceInQ' a b c valueOfOptC'$_stCloseBraceInQ'='$_stq' \
      -- \
      --optA --opt-a --opt-c valueOfOptC alpha -ab \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}<empty>
      ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}--opt-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(0)
      ${_ansiOutline}OPT_optC${_ansiReset}(1): ${_ansiSecondaryLiteral}-c${_ansiReset}
    "
    --stderrParaWithIndent "0"
    --before "local __defaultValueOfArgC=defaultValueOfArgC"
    --before "local __partOfDefaultValueOfArgD=DefaultValueOfArgD"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC \
          OPT_optA OPT_optB OPT_optC \
      '$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      --aParams '$_stqq' --opt-c/c'$_stqq' \
      --opt-a/a \
      --opt-b/b \
      -- \
      --opt-a alpha -c \
    "
  '

 # Влияние --canBeMixedOptionsAndArgs на наследование опций "через голову"
  # '
  #   --before "_substitute noStack true"
  #   --after "_restore noStack"
  #   --stderr "
  #     ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
  #     ${_ansiOutline}optB${_ansiReset}<empty>
  #     ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
  #     ${_ansiOutline}optD${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
  #     ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
  #     ${_ansiOutline}OPT_optB${_ansiReset}(0)
  #     ${_ansiOutline}OPT_optC${_ansiReset}(1): ${_ansiSecondaryLiteral}-c${_ansiReset}
  #     ${_ansiOutline}OPT_optD${_ansiReset}(1): ${_ansiSecondaryLiteral}-d${_ansiReset}
  #   "
  #   --stderrParaWithIndent "0"
  #   "_prepareCodeToParseFuncParams2TestFunc3 \
  #     --run '$_stq' \
  #       _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
  #         optA optB optC optD \
  #         OPT_optA OPT_optB OPT_optC OPT_optD \
  #     '$_stq' \
  #     --aOpt '$_stq' --isCommandWrapper'$_stq' \
  #     --aParams '$_stq' --opt-c/c'$_stq' \
  #     --adParams '$_stq' --opt-d/d'$_stq' \
  #     --isCommandWrapper \
  #     --canBeMixedOptionsAndArgs \
  #     --opt-a/a \
  #     --opt-b/b \
  #     -- \
  #     alpha -c delta -ad \
  #   "
  # '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeToParseFuncParams2TestFunc3 alpha delta${_ansiErr} не ожидает опцию ${_ansiCmd}-c${_ansiReset}"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --aOpt '$_stq' --isCommandWrapper'$_stq' \
      --aParams '$_stq' --opt-c/c'$_stq' \
      --adParams '$_stq' --opt-d/d'$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      --opt-a/a \
      --opt-b/b \
      -- \
      alpha delta -c -ad \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optB${_ansiReset}<empty>
      ${_ansiOutline}optC${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}optD${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
      ${_ansiOutline}OPT_optB${_ansiReset}(0)
      ${_ansiOutline}OPT_optC${_ansiReset}(1): ${_ansiSecondaryLiteral}-c${_ansiReset}
      ${_ansiOutline}OPT_optD${_ansiReset}(1): ${_ansiSecondaryLiteral}-d${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA optB optC optD \
          OPT_optA OPT_optB OPT_optC OPT_optD \
      '$_stq' \
      --aOpt '$_stq' --canBeMixedOptionsAndArgs --isCommandWrapper'$_stq' \
      --aParams '$_stq' --opt-c/c'$_stq' \
      --adParams '$_stq' --opt-d/d'$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      --opt-a/a \
      --opt-b/b \
      -- \
      alpha delta -c -ad \
    "
  '

 # important-опции
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}<empty>
      ${_ansiOutline}OPT_optA${_ansiReset}(0)
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA \
          OPT_optA \
      '$_stq' \
      --aParams '$_stq' --opt-a/a'$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      --opt-a/a \
      -- \
      -a alpha \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA \
          OPT_optA \
      '$_stq' \
      --aParams '$_stq' --opt-a/a'$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      !--opt-a/a \
      -- \
      -a alpha \
    "
  '
  '
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "
      ${_ansiOutline}optA${_ansiReset}: ${_ansiPrimaryLiteral}true${_ansiReset}
      ${_ansiOutline}OPT_optA${_ansiReset}(1): ${_ansiSecondaryLiteral}-a${_ansiReset}
    "
    --stderrParaWithIndent "0"
    "_prepareCodeToParseFuncParams2TestFunc3 \
      --run '$_stq' \
        _debugVar --clean --joiner '$_stqq''$_stSlashInQ'n'$_stqq' \
          optA \
          OPT_optA \
      '$_stq' \
      --aParams '$_stq' !--opt-a/a'$_stq' \
      --isCommandWrapper \
      --canBeMixedOptionsAndArgs \
      !--opt-a/a \
      -- \
      -a alpha \
    "
  '
)
_prepareCodeToParseFuncParams2TestFunc3BoolOptions=(
  'treatUnknownOptionAsArg'
  'isCommandWrapper'
  'canBeMixedOptionsAndArgs'
  'canBeMoreParams'
)
_prepareCodeToParseFuncParams2TestFunc3ScalarOptions=(
  'aOpt'
  'aParams'
  'adOpt'
  'adParams'
)
_prepareCodeToParseFuncParams2TestFunc3ListOptions=(
  'run'
)
_prepareCodeToParseFuncParams2TestFunc3OptionsOpt=(--treatUnknownOptionAsArg)
_prepareCodeToParseFuncParams2TestFunc3() { eval "$_funcOptions2"
  local -a __run=( "${run[@]}" )
  local __isCommandWrapper="$isCommandWrapper"
  local -a _prepareCodeToParseFuncParams2TestFunc3ParamsOpt=(
    --additionalSuffix $testNum
    "${OPT_canBeMixedOptionsAndArgs[@]}"
    "${OPT_isCommandWrapper[@]}"
    "${OPT_treatUnknownOptionAsArg[@]}"
    "${OPT_canBeMoreParams[@]}"
  )
  local -a _prepareCodeToParseFuncParams2TestFunc3Params=()
  while [[ $# -gt 0 ]]; do
    [[ $1 == -- ]] && shift && break
    _prepareCodeToParseFuncParams2TestFunc3Params+=( "$1" )
    shift
  done

  eval local -a _prepareCodeToParseFuncParams2TestFunc3_alphaParams=\( \$aParams \)
  local -a _prepareCodeToParseFuncParams2TestFunc3_alphaParamsOpt=(
    --additionalSuffix $testNum
    $aOpt
  )

  eval local -a _prepareCodeToParseFuncParams2TestFunc3_alpha_deltaParams=\( \$adParams \)
  local -a _prepareCodeToParseFuncParams2TestFunc3_alpha_deltaParamsOpt=(
    --additionalSuffix $testNum
    $adOpt
  )

  eval "$_funcParams2"

  _prepareCodeToParseFuncParams2TestFunc3Helper
}
_prepareCodeToParseFuncParams2TestFunc3Helper() {
  if [[ -z $__isCommandWrapper ]]; then
    local item; for item in "${__run[@]}"; do
      eval "$item"
    done
  fi
}
_prepareCodeToParseFuncParams2TestFunc3_alphaShortcuts=( 'alef' 'a' )
_prepareCodeToParseFuncParams2TestFunc3_alpha() { eval "$_codeToCallFuncParams2"
  _prepareCodeToParseFuncParams2TestFunc3Helper
}
_prepareCodeToParseFuncParams2TestFunc3_bravoCharlieShortcuts=( 'beta-gamma' 'b' )
_prepareCodeToParseFuncParams2TestFunc3_bravoCharlie() { eval "$_codeToCallFuncParams2"
  true
}
_prepareCodeToParseFuncParams2TestFunc3_alpha_delta() { eval "$_codeToCallFuncParams2"
  _prepareCodeToParseFuncParams2TestFunc3Helper
}
_prepareCodeToParseFuncParams2TestFunc3_alpha_echoFoxtrot() { eval "$_codeToCallFuncParams2"
  true
}

# =============================================================================

_parseFuncParams2Tests=(
# completion
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"--help \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"--help \" [1]=\"-? \" [2]=\"-h \" [3]=\"--aOpt \" [4]=\"--a-opt \" [5]=\"-a \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"--help \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"--\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"-h \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-h\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"-? \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" ) \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"--aOpt \" [1]=\"--a-opt \" [2]=\"-a \" [3]=\"--bOpt \" [4]=\"--b-opt \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"-\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a --bOpt ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc ожидает целое число\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"-a\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a:.. ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc ожидает целое число\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"--aOpt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a:.. ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc ожидает целое число из диапазона 4..5 или пустое значение\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"-a\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a:?:4..5 ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc ожидает целое число из диапазона 4..5 или пустое значение\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"--aOpt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a:?:4..5 ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc не ожидает опцию --some\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"--some\" \"--aOpt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( --aOpt/a:?:4..5 ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "eval \"_parseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _parseFuncParams2TestFunc_alpha"
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"<-- HINT: _parseFuncParams2TestFunc a не ожидает опцию --some\" [1]=\"-->\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" a \"--some\" \"--aOpt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local -a _parseFuncParams2TestFuncParams=() \
      && local -a _parseFuncParams2TestFuncParamsOpt=( --isCommandWrapper ) \
      && local -a _parseFuncParams2TestFunc_alphaParams=( --aOpt/a:?:4..5 ) \
      && local -a _parseFuncParams2TestFunc_alphaShortcuts=( a ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc_alpha.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"value1 \" [1]=\"value2 \" [2]=\"value3 \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"-a\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( '$_stq'--aOpt/a:?:(value1 value2 value3)'$_stq' ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"value1 \" [1]=\"value2 \" [2]=\"value3 \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"--aOpt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( '$_stq'--aOpt/a:?:(value1 value2 value3)'$_stq' ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"value1 \" [1]=\"value2 \" [2]=\"value3 \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"-?\" \"--a-opt\" \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( '$_stq'--aOpt/a:?:(value1 value2 value3)'$_stq' ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=()"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( '$_stq'--aOpt/a:?:(value1 value2 value3)'$_stq' ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"value1 \" [1]=\"value2 \" [2]=\"value3 \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=( '$_stq'arg:(value1 value2 value3)'$_stq' ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"alpha \" [1]=\"alef \" [2]=\"alfa \" [3]=\"bravo \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    --before "eval \"_parseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "eval \"_parseFuncParams2TestFunc_bravo() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "local _parseFuncParams2TestFunc_alphaShortcuts=(alef alfa)"
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFuncParams=() \
      && local _parseFuncParams2TestFuncParamsOpt=( --isCommandWrapper ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc_alpha.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
  '
    --before "local -a COMPREPLY=()"
    --varName "COMPREPLY"
    --varValue "declare -a COMPREPLY=([0]=\"--help \" [1]=\"-? \" [2]=\"-h \" [3]=\"--optA \" [4]=\"--opt-a \" [5]=\"-a \" [6]=\"--optB \" [7]=\"--opt-b \" [8]=\"-b \")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    --before "eval \"_parseFuncParams2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --before "local _parseFuncParams2TestFunc_alphaShortcuts=(alef alfa)"
    "true \
      && local -a COMP_WORDS=( _parseFuncParams2TestFunc \"alfa\" \"-\") \
      && local COMP_CWORD=\$(( \${#COMP_WORDS[@]} - 1 )) \
      && local _parseFuncParams2TestFunc_alphaParams=( --optB/b ) \
      && local _parseFuncParams2TestFuncParams=( --optA/a ) \
      && local _parseFuncParams2TestFuncParamsOpt=( --isCommandWrapper --canBeMixedOptionsAndArgs ) \
      && __onlyPrepareCode=true _parseFuncParams2TestFunc \
      || . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc.\$testNum.completion.code.bash\" \
      && . \"\$_bwDir/tests/generated/_parseFuncParams2TestFunc_alpha.\$testNum.completion.code.bash\" \
      && __complete__parseFuncParams2TestFunc \
    "
  '
)
_codeToCallFuncParams2='
  dstVarName=__paramsOpt srcVarName=${FUNCNAME[0]}ParamsOpt codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
  __paramsOpt+=( --additionalDependencies "'${BASH_SOURCE[0]}'" --additionalSuffix $testNum )
  dstVarName=${FUNCNAME[0]}ParamsOpt srcVarName=__paramsOpt codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
  # _debugVar ${FUNCNAME[0]}ParamsOpt
  eval "$_funcParams2"
'
_parseFuncParams2TestFunc() { eval "$_codeToCallFuncParams2"
}

# =============================================================================

_prepareCodeOfAutoHelp2Tests=(
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    "_prepareCodeOfAutoHelp2TestFunc -h"
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset} Однострочное описание
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    "true \
      && local _prepareCodeOfAutoHelp2TestFunc_description=\"Однострочное описание\" \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset}
        Длинное
        многострочное
        описание
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    "true \
      && local _prepareCodeOfAutoHelp2TestFunc_description=\"Длинное\${_nl}многострочное\${_nl}описание\" \
      && _prepareCodeOfAutoHelp2TestFunc${testNum} -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}a${_ansiReset} ${_ansiOutline}c${_ansiReset} ${_ansiOutline}d${_ansiReset} ${_ansiOutline}\"Аргумент b\"${_ansiReset} ${_ansiOutline}АргументE${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}a${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_a_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}c${_ansiReset} Однострочное \"описание\" аргумента
      ${_ansiOutline}d${_ansiReset}
        Многострочное
        \"описание\"
        аргумента
      ${_ansiOutline}\"Аргумент b\"${_ansiReset}
      ${_ansiOutline}АргументE${_ansiReset} Однострочное \"описание\" аргумента
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "local _prepareCodeOfAutoHelp2TestFunc_b_name=\"Аргумент b\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_c_description=\"Однострочное \\\"описание\\\" аргумента\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_d_description=\"Многострочное\${_nl}\\\"описание\\\"\${_nl}аргумента\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_e_name=\"АргументE\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_e_description=\"Однострочное \\\"описание\\\" аргумента\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( a c d b e ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}a${_ansiReset} ${_ansiOutline}b${_ansiReset} [${_ansiOutline}c${_ansiReset} [${_ansiOutline}d${_ansiReset}...]]
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}a${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_a_description${_ansiErr} аргумента${_ansiReset}
        ${_ansiOutline}Значение${_ansiReset} может быть пустым
      ${_ansiOutline}b${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_b_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}c${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_c_description${_ansiErr} аргумента${_ansiReset}
        ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}d${_ansiReset}
      ${_ansiOutline}d${_ansiReset}... возможно пустой список ${_ansiUnderline}уникальных${_ansiReset} ${_ansiOutline}значений${_ansiReset}
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( \
        a:? \
        b \
        c=d \
        @d:unique \
      ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}a${_ansiReset} ${_ansiOutline}b${_ansiReset} ${_ansiOutline}c${_ansiReset} ${_ansiOutline}d${_ansiReset}...
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}a${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_a_description${_ansiErr} аргумента${_ansiReset}
        ${_ansiOutline}Значение${_ansiReset} может быть пустым
      ${_ansiOutline}b${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_b_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}c${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_c_description${_ansiErr} аргумента${_ansiReset}
        ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}d${_ansiReset}
      ${_ansiOutline}d${_ansiReset}... список (не менее 2 элементов) ${_ansiOutline}значений${_ansiReset} имен \"тестов\"
        Варианты ${_ansiOutline}значения${_ansiReset}: ${_ansiSecondaryLiteral}testA testB${_ansiReset}
        ${_ansiOutline}Значение${_ansiReset} может быть пустым
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "local _prepareCodeOfAutoHelp2TestFunc_d_description=\"имен \\\"тестов\\\"\""
    --before "local -a __testNames=( testA testB )"
    --before "eval \"__testNames() { echo \\\${__testNames[@]}; }\""
    --after "unset -f __testNames"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( \
        a:? \
        b \
        c=d \
        '$_stq'@2..d:'$_stOpenBraceInQ' '$_stDollarInQ''$_stOpenBraceInQ'__testNames'$_stCloseBraceInQ' '$_stCloseBraceInQ':?'$_stq' \
      ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}testNames${_ansiReset}...
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}testNames${_ansiReset}... непустой список ${_ansiOutline}значений${_ansiReset} имен \"тестов\"
        Варианты ${_ansiOutline}значения${_ansiReset}: ${_ansiSecondaryLiteral}testA testB${_ansiReset}
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "local -a __testNames=( testA )"
    --before "eval \"__testNames() { echo \\\${__testNames[@]}; }\""
    --after "unset -f __testNames"
    "true \
      && local _prepareCodeOfAutoHelp2TestFunc_testNames_description=\"имен \\\"тестов\\\"\" \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( '$_stq'@1..testNames:( '$_stDollarInQ'(__testNames) )'$_stq' ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --needPregenHelp ) \
      && __onlyPrepareCode=true _prepareCodeOfAutoHelp2TestFunc \
      || __testNames=( testA testB ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опции${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}Опции${_ansiReset}
        ${_ansiCmd}--opt-a${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optA_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--opt-b${_ansiReset} или ${_ansiCmd}-b${_ansiReset}
          Однострочное \"описание\" опции
        ${_ansiCmd}--opt-c${_ansiReset} или ${_ansiCmd}-c${_ansiReset} или ${_ansiCmd}-o${_ansiReset}
          Многострочное
          \"описание\"
          опции
        ${_ansiCmd}--opt-d${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optD_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--opt-e${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optE_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} может быть пустым
        ${_ansiCmd}--opt-f${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optF_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - целое число
        ${_ansiCmd}--opt-g${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optG_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - неотрицательное целое число
        ${_ansiCmd}--opt-h${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optH_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - положительное целое число, может быть пустым
          ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}1${_ansiReset}
        ${_ansiCmd}--opt-l${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optL_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - целое число не менее ${_ansiPrimaryLiteral}2${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}2${_ansiReset}
        ${_ansiCmd}--opt-m${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optM_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - неположительное целое число
        ${_ansiCmd}--opt-n${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optN_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - отрицательное целое число
        ${_ansiCmd}--opt-o${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optO_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - целое число не более ${_ansiPrimaryLiteral}-2${_ansiReset}
        ${_ansiCmd}--opt-i${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optI_description${_ansiErr} опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} - целое число из диапазона ${_ansiSecondaryLiteral}2..6${_ansiReset}, может быть пустым
          ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiOutline}\$ defaultValueOfH${_ansiReset}
        ${_ansiCmd}--opt-j${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optJ_description${_ansiErr} опции${_ansiReset}
          Варианты ${_ansiOutline}значения${_ansiReset}: ${_ansiSecondaryLiteral}a \"b c\" d${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} может быть пустым
          ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}d${_ansiReset}
        ${_ansiCmd}--opt-k${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          Варианты ${_ansiOutline}значения${_ansiReset}:
            ${_ansiPrimaryLiteral}a${_ansiReset} Однострочное \"описание\" значения опции
            ${_ansiPrimaryLiteral}\"b ? c\"${_ansiReset}
              Многострочное
              \"описание\"
              значения опции
            ${_ansiPrimaryLiteral}d${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optK_d_description${_ansiErr} значения опции${_ansiReset}
          ${_ansiOutline}Значение${_ansiReset} может быть пустым
          ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiOutline}\$ defaultValueOfK${_ansiReset}
        ${_ansiCmd}--opt-p${_ansiReset} ${_ansiOutline}значение${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_optP_description${_ansiErr} опции${_ansiReset}
          Опция предназначена для того, чтобы сформировать
          возможно пустой список ${_ansiOutline}значений${_ansiReset}
          путем eё многократного использования
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --stdoutEtaPreProcess "perl -pe \"s/\\\\\$\s/\\\\\$/g\""
    --stdoutEtaPreProcess "perl -pe \"\\\$indent=\\\" \\\" x 10; s/^\s{12}/\\\$indent/g\""
    --stdoutEtaPreProcess "perl -pe \"\\\$indent=\\\" \\\" x 12; s/^\s{14}/\\\$indent/g\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_optB_description=\"Однострочное \\\"описание\\\" опции\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_optC_description=\"Многострочное\${_nl}\\\"описание\\\"\${_nl}опции\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_optK_a_description=\"Однострочное \\\"описание\\\" значения опции\""
    --before "local _prepareCodeOfAutoHelp2TestFunc_optK_bc_description=\"Многострочное\${_nl}\\\"описание\\\"\${_nl} значения опции\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( \
        --opt-a \
        --opt-b/b \
        --opt-c/co \
        --opt-d: \
        --opt-e:? \
        --opt-f:.. \
        --opt-g:0.. \
        --opt-h:1..:?=1 \
        --opt-l:2..=2 \
        --opt-m:..0 \
        --opt-n:..-1 \
        --opt-o:..-2 \
        '$_stq'--opt-i:2..6:?='$_stDollarInQ'defaultValueOfH'$_stq' \
        '$_stq'--opt-j:?:'$_stOpenBraceInQ'a '$_stqq'b c'$_stqq' d'$_stCloseBraceInQ'=d'$_stq' \
        '$_stq'--opt-k:?:'$_stOpenBraceInQ'a '$_stqq'b ? c'$_stqq' d'$_stCloseBraceInQ'='$_stDollarInQ'defaultValueOfK'$_stq' \
        @--opt-p \
      ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '

  # --isCommandWrapper
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}Команда${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}Команда${_ansiReset} - один из следующих вариантов: ${_ansiSecondaryLiteral}alpha alef bravo-charlie bravo beta charlie gamma delta${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc alpha${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc alef${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_description${_ansiErr} команды${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc bravo-charlie${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc bravo${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc beta${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc charlie${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc gamma${_ansiReset}
          Однострочное описание
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc delta${_ansiReset}
          Многострочное
          длинное
          описание
      ${_ansiHeader}Подробнее см.${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Команда ${_ansiCmd}--help${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Команда ${_ansiCmd}-?${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Команда ${_ansiCmd}-h${_ansiReset}
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_alpha"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaShortcuts=(alef)"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_bravoCharlie"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_bravoCharlieShortcuts=(bravo beta charlie gamma)"
    --before "local _prepareCodeOfAutoHelp2TestFunc_bravoCharlie_description=\"Однострочное описание\""
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_delta() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_delta"
    --before "local _prepareCodeOfAutoHelp2TestFunc_delta_description=\"Многострочное\${_nl}длинное\${_nl}описание\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --isCommandWrapper ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опция${_ansiReset}] ${_ansiOutline}Цель${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}Цель${_ansiReset} - один из следующих вариантов: ${_ansiSecondaryLiteral}alpha alef bravo-charlie bravo beta charlie gamma delta${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc alpha${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc alef${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_description${_ansiErr} команды${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc bravo-charlie${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc bravo${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc beta${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc charlie${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc gamma${_ansiReset}
          Однострочное описание
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc delta${_ansiReset}
          Многострочное
          длинное
          описание
      ${_ansiHeader}Подробнее см.${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Цель ${_ansiCmd}--help${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Цель ${_ansiCmd}-?${_ansiReset}
        ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc ${_ansiOutline}Цель ${_ansiCmd}-h${_ansiReset}
      ${_ansiOutline}Опция${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_alpha"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaShortcuts=(alef)"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_bravoCharlie"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_bravoCharlieShortcuts=(bravo beta charlie gamma)"
    --before "local _prepareCodeOfAutoHelp2TestFunc_bravoCharlie_description=\"Однострочное описание\""
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_delta() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_delta"
    --before "local _prepareCodeOfAutoHelp2TestFunc_delta_description=\"Многострочное\${_nl}длинное\${_nl}описание\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --isCommandWrapper ) \
      && local _prepareCodeOfAutoHelp2TestFunc_cmd_name=Цель \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc alpha${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}arg${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_description${_ansiReset}
      ${_ansiOutline}arg${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_arg_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}Опции${_ansiReset}
        ${_ansiCmd}--a${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_a_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_alpha"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaShortcuts=(alef)"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaParams=(--a arg)"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_bravoCharlie"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_bravoCharlieShortcuts=(bravo beta charlie gamma)"
    --before "local _prepareCodeOfAutoHelp2TestFunc_bravoCharlieDescription=\"Однострочное описание\""
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_delta() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_delta"
    --before "local _prepareCodeOfAutoHelp2TestFunc_deltaDescription=\"Многострочное\${_nl}длинное\${_nl}описание\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --isCommandWrapper ) \
      && _prepareCodeOfAutoHelp2TestFunc -h alpha \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc_alpha${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}arg${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_description${_ansiReset}
      ${_ansiOutline}arg${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_arg_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}Опции${_ansiReset}
        ${_ansiCmd}--a${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_alpha_a_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_alpha() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_alpha"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaShortcuts=(alef)"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_alphaParams=(--a arg)"
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_bravoCharlie() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_bravoCharlie"
    --before "local -a _prepareCodeOfAutoHelp2TestFunc_bravoCharlieShortcuts=(bravo beta charlie gamma)"
    --before "local _prepareCodeOfAutoHelp2TestFunc_bravoCharlieDescription=\"Однострочное описание\""
    --before "eval \"_prepareCodeOfAutoHelp2TestFunc_delta() { eval \\\"\\\$_codeToCallFuncParams2\\\"; }\""
    --after "unset -f _prepareCodeOfAutoHelp2TestFunc_delta"
    --before "local _prepareCodeOfAutoHelp2TestFunc_deltaDescription=\"Многострочное\${_nl}длинное\${_nl}описание\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --isCommandWrapper ) \
      && _prepareCodeOfAutoHelp2TestFunc_alpha -h \
    "
  '

# --canBeMoreParams
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}arg${_ansiReset} ${_ansiOutline}moreArg1${_ansiReset} ${_ansiOutline}\"Дополнительный аргумент\"${_ansiReset} ${_ansiOutline}moreArg3${_ansiReset}
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}arg${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_arg_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}moreArg1${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_moreArg1_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}\"Дополнительный аргумент\"${_ansiReset}
      ${_ansiOutline}moreArg3${_ansiReset} Однострочное описание
      ${_ansiOutline}Опции${_ansiReset}
        ${_ansiCmd}--a${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_a_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( --a arg ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncMoreArgs=( moreArg1 moreArg2 moreArg3 ) \
      && local _prepareCodeOfAutoHelp2TestFunc_moreArg2_name=\"Дополнительный аргумент\" \
      && local _prepareCodeOfAutoHelp2TestFunc_moreArg3_description=\"Однострочное описание\" \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}arg${_ansiReset} ${_ansiOutline}moreArg1${_ansiReset} [ ${_ansiOutline}\"Дополнительный аргумент\"${_ansiReset} [ ${_ansiOutline}moreArg1${_ansiReset} ] ] ...
      ${_ansiHeader}Описание:${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_description${_ansiReset}
      ${_ansiOutline}arg${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_arg_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}moreArg1${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_moreArg1_description${_ansiErr} аргумента${_ansiReset}
      ${_ansiOutline}\"Дополнительный аргумент\"${_ansiReset}
      ${_ansiOutline}moreArg3${_ansiReset} Однострочное описание
      ${_ansiOutline}Опции${_ansiReset}
        ${_ansiCmd}--a${_ansiReset}
          ${_ansiErr}Нет описания ${_ansiOutline}_prepareCodeOfAutoHelp2TestFunc_a_description${_ansiErr} опции${_ansiReset}
        ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
          $_help_description
    "
    --stdoutParaWithIndent "0"
    --before "local _prepareCodeOfAutoHelp2TestFunc_moreArg2_name=\"Дополнительный аргумент\""
    --before "local _prepareCodeOfAutoHelp2TestFuncMoreArgsUsage=\"\${_ansiOutline}moreArg1\${_ansiReset} [ \${_ansiOutline}\\\"\$_prepareCodeOfAutoHelp2TestFunc_moreArg2_name\\\"\${_ansiReset} [ \${_ansiOutline}moreArg1\${_ansiReset} ] ] ...\""
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( --a arg ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncMoreArgs=( moreArg1 moreArg2 moreArg3 ) \
      && local _prepareCodeOfAutoHelp2TestFunc_moreArg3_description=\"Однострочное описание\" \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '

# --canBeMoreParams errors

  '
    --return "1"
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiErr} ожидает, что ${_ansiOutline}_prepareCodeOfAutoHelp2TestFuncMoreArgsUsage${_ansiErr} будет скаляром, а не массивом${_ansiReset}"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncMoreArgsUsage=( some ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "1"
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiErr} ожидает, что ${_ansiOutline}_prepareCodeOfAutoHelp2TestFuncMoreArgs${_ansiErr} будет массивом, а не скаляром${_ansiReset}"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local _prepareCodeOfAutoHelp2TestFuncMoreArgs=\"moreArg1 moreArg2 moreArg3\" \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "1"
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiErr} ожидает, что имя переменной ${_ansiPrimaryLiteral}Дополнительный аргумент${_ansiErr} указанное в ${_ansiOutline}_prepareCodeOfAutoHelp2TestFuncMoreArgs${_ansiErr} $_mustBeValidVarName${_ansiReset}"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncMoreArgs=( \"Дополнительный аргумент\" ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
  '
    --return "1"
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_prepareCodeOfAutoHelp2TestFunc${_ansiErr} ${_ansiOutline}_prepareCodeOfAutoHelp2TestFuncMoreArgs${_ansiErr} содержит аргумент ${_ansiPrimaryLiteral}arg${_ansiErr} одноименный параметру из ${_ansiOutline}_prepareCodeOfAutoHelp2TestFuncParams${_ansiReset}"
    "true \
      && local -a _prepareCodeOfAutoHelp2TestFuncParamsOpt=( --canBeMoreParams ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncParams=( --a arg ) \
      && local -a _prepareCodeOfAutoHelp2TestFuncMoreArgs=( arg ) \
      && _prepareCodeOfAutoHelp2TestFunc -h \
    "
  '
)

_prepareCodeOfAutoHelp2TestFunc() {
  codeHolder=_codeToCallFuncParams2 eval "$_evalCode"
}

# =============================================================================
