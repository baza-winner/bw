#!/bin/bash

{
  # shellcheck disable=SC2154
  if [[ -n $shellcheck ]]; then 
    . bw.bash
    . core/profileSupport.bash
    . core/coreUtils.bash
  fi
}

# =============================================================================

_resetBash

# =============================================================================

# shellcheck disable=SC2016,SC2034
{
_funcOptions2Helper='_prepareCodeToParseFuncOptions2 && . "$codeFileSpec" || return $?'
_funcOptions2='local codeFileSpec; '$_funcOptions2Helper
}

# shellcheck disable=SC2016,SC2034
_codeToInitLocalCopyOfScalar='
  local $dstVarName="${!srcVarName}"
'

# shellcheck disable=SC2016
_codeToInitLocalCopyOfArray='
  eval local -a $dstVarName=\( \"\${$srcVarName[@]}\" \)
'

_prepareCodeToParseFuncOptions2BoolOptions=( 'treatUnknownOptionAsArg' )
_prepareCodeToParseFuncOptions2ScalarOptions=( 'prefix' 'additionalSuffix' )
_prepareCodeToParseFuncOptions2() {
  local funcName="${FUNCNAME[1]}" __thisFuncCommand='' onlyPrepareCode="$__onlyPrepareCode" __onlyPrepareCode=''
  kind=Options codeHolder=_codeToCheckParams eval "$_evalCode"

  # shellcheck disable=SC2034
  local additionalSuffix prefix
  local optHolder="${funcName}OptionsOpt"
  local typeOfVar; typeOfVar=$(_getTypeOfVar "$optHolder")
  if [[ $typeOfVar != 'none' ]]; then
    # shellcheck disable=SC2046
    [[ $typeOfVar == 'array' ]] \
      || return $(_throw "ожидает, что ${_ansiOutline}$optHolder${_ansiErr} будет массивом, а не скаляром")
    eval set -- '"${'"$optHolder"'[@]}"'
    eval "$_funcOptions2Helper"
    eval "$_codeToCheckNoArgsInOpt"
  else
    local varName 
    # shellcheck disable=SC2154
    for varName in "${_prepareCodeToParseFuncOptions2BoolOptions[@]}" "${_prepareCodeToParseFuncOptions2ScalarOptions[@]}"; do
      eval local "$varName="
    done
    # shellcheck disable=SC2154
    for varName in "${_prepareCodeToParseFuncOptions2ListOptions[@]}"; do
      eval local -a "$varName="'()'
    done
  fi

  dstVarName='' codeType=funcOptions fileSpec='' originalCodeDeep='' eval "$_codeToPrepareCodeFileSpec"

  # shellcheck disable=SC2154
  if [[ -f $codeFileSpec ]]; then
    if [[ -z $onlyPrepareCode && -z $_isBwDevelop && -z $_isBwDevelopInherited ]]; then
      return 0
    elif _everyFileNotNewerThan "$codeFileSpec" "${BASH_SOURCE[@]::2}"; then
      if [[ -n $onlyPrepareCode ]]; then
        return 2
      else
        return 0
      fi
    fi
  fi

  local verbose=
  # verbose=true
  if [[ -n $verbose ]]; then
    _warn "Создаем ${_ansiFileSpec}$codeFileSpec"
  fi

  dstVarName=boolOptions srcVarName=${funcName}BoolOptions eval "$_codeToInitLocalCopyOfArray"
  dstVarName=scalarOptions srcVarName=${funcName}ScalarOptions eval "$_codeToInitLocalCopyOfArray"
  dstVarName=listOptions srcVarName=${funcName}ListOptions eval "$_codeToInitLocalCopyOfArray"
  local err=

  local -a boolOptionVarNames
    # shellcheck disable=SC2154,SC2046
    _initVarNamesWithOptions2 boolOptionVarNames "${boolOptions[@]}" || return $(errOrigin=1 _throw "$err")
  local -a scalarOptionVarNames
    # shellcheck disable=SC2154,SC2046
    _initVarNamesWithOptions2 scalarOptionVarNames "${scalarOptions[@]}" || return $(errOrigin=1 _throw "$err")
  local -a listOptionVarNames
    # shellcheck disable=SC2154,SC2046
    _initVarNamesWithOptions2 listOptionVarNames "${listOptions[@]}" || return $(errOrigin=1 _throw "$err")
  local -a OPTVarNames=();
    local varName; for varName in "${boolOptionVarNames[@]}" "${scalarOptionVarNames[@]}" "${listOptionVarNames[@]}"; do
      OPTVarNames+=( "OPT_$varName" )
    done

  local __ownPrefix=__

  local initIsBoolOption='local '$__ownPrefix'isBoolOption='
  if [[ ${#boolOptionVarNames[@]} -gt 0 ]]; then
    local conditionCode; _prepareHasItemConditionCode "\$${__ownPrefix}varName" "${boolOptionVarNames[@]}"
    initIsBoolOption+='; [[ '$conditionCode' ]] && '$__ownPrefix'isBoolOption=true'
  fi
  local initIsListOption='local '$__ownPrefix'isListOption='
  if [[ ${#listOptionVarNames[@]} -gt 0 ]]; then
    local conditionCode; _prepareHasItemConditionCode "\$${__ownPrefix}varName" "${listOptionVarNames[@]}"
    initIsListOption+='; [[ -z $'$__ownPrefix'isBoolOption ]] && [[ '$conditionCode' ]] && '$__ownPrefix'isListOption=true'
  fi
  local initIsScalarOption='local '$__ownPrefix'isScalarOption='
  if [[ ${#scalarOptionVarNames[@]} -gt 0 ]]; then
    local conditionCode; _prepareHasItemConditionCode "\$${__ownPrefix}varName" "${scalarOptionVarNames[@]}"
    initIsScalarOption+='; [[ -z $'$__ownPrefix'isBoolOption && -z $'$__ownPrefix'isScalarOption ]] && [[ '$conditionCode' ]] && '$__ownPrefix'isScalarOption=true'
  fi

  local unexpectedOption;
  # shellcheck disable=SC2154
  if [[ -z $treatUnknownOptionAsArg ]]; then
    # shellcheck disable=SC2016
    unexpectedOption='return $(errOrigin=1 _throw "не ожидает опцию ${_ansiCmd}$1")'
  else
    unexpectedOption+='break'
  fi

  local code
  code+=$(_getDeclarationOfVars "${boolOptionVarNames[@]}" "${scalarOptionVarNames[@]}")
  code+=$(_getDeclarationOfVars -a "${listOptionVarNames[@]}" "${OPTVarNames[@]}")

  local terminateProcessing
  if [[ -n $treatUnknownOptionAsArg ]]; then
    terminateProcessing='break'
  else
    terminateProcessing='shift && break'
  fi

  # shellcheck disable=SC2016
  code+='
  local __funcName='$funcName' __thisFuncCommand=
  [[ -z $__funcCommand ]] || local __funcCommand=
  while [[ $1 =~ ^-- ]]; do
    [[ $1 == -- ]] && '$terminateProcessing'
    local '$__ownPrefix'varName="${1:2}"
    # [[ ! $'$__ownPrefix'varName =~ - ]] || '$__ownPrefix'varName=$(_kebabCaseToCamelCase "$'$__ownPrefix'varName")
    dstVarName='$__ownPrefix'varName _kebabCaseToCamelCase "$'$__ownPrefix'varName"
    '$initIsBoolOption'
    if [[ -n $'$__ownPrefix'isBoolOption ]]; then
      eval '$prefix'$'$__ownPrefix'varName=true
      eval OPT_$'$__ownPrefix'varName+=\( \"\$1\" \)
    else
      '$initIsListOption'
      '$initIsScalarOption'
      if [[ -z $'$__ownPrefix'isListOption && -z $'$__ownPrefix'isScalarOption ]]; then
        '$unexpectedOption'
      else
        if [[ $# -le 1 ]]; then
          return $(errOrigin=1 _throw "ожидает, что опция ${_ansiCmd}$1${_ansiErr} будет снабжена значением")
        fi
        # if [[ $2 =~ ^-.+  && ! $2 =~ ^-[[:digit:]]+$  ]]; then
        #   return $(errOrigin=1 _throw "ожидает, что опция ${_ansiCmd}$1${_ansiErr} будет снабжена значением не похожим на опцию ${_ansiPrimaryLiteral}$2")
        # fi
        if [[ -n $'$__ownPrefix'isListOption ]]; then
          eval '$prefix'$'$__ownPrefix'varName+=\( \"\$2\" \)
          eval OPT_$'$__ownPrefix'varName+=\( \"\$1\" \"\$2\" \)
        else
          local '$__ownPrefix'varValueHolder="'$prefix'$'$__ownPrefix'varName"
          if [[ -z ${!'$__ownPrefix'varValueHolder} ]]; then
            eval '"$prefix"'$'"$__ownPrefix"'varName=\"\$2\"
            eval OPT_$'"$__ownPrefix"'varName=\( \"\$1\" \"\$2\" \)
          elif [[ ${!'"$__ownPrefix"'varValueHolder} != $2 ]]; then
            local __optionName="$1"
            return $(errOrigin=1 _throw "не ожидает, что опция ${_ansiCmd}$__optionName${_ansiErr} будет указана повторно")
          fi
        fi
        shift
      fi
    fi
    shift
  done'

  _assureDir "$(dirname "$codeFileSpec")" || return $?
  echo "$code" > "$codeFileSpec"
  [[ -z $onlyPrepareCode ]] || return 2
}

#shellcheck disable=SC2034
{
_mustBeValidVarName="${_ansiErr}должно быть ${_ansiOutline}валидным идентификатором${_ansiErr}, т.е. начинаться с подчерка (${_ansiSecondaryLiteral}_${_ansiErr}) или буквы (${_ansiSecondaryLiteral}a-zA-Z${_ansiErr}), после чего возможно следуют подчерки (${_ansiSecondaryLiteral}_${_ansiErr}), буквы (${_ansiSecondaryLiteral}a-zA-Z${_ansiErr}) или цифры (${_ansiSecondaryLiteral}0-9${_ansiErr}) без пробелов и других символов"
_mustBeValidCommandShortcut="${_ansiErr}должно начинаться с подчерка (${_ansiSecondaryLiteral}_${_ansiErr}), буквы в нижнем регистре (${_ansiSecondaryLiteral}a-z${_ansiErr}) или цифры (${_ansiSecondaryLiteral}0-9${_ansiErr}), после чего возможно следуют подчерки (${_ansiSecondaryLiteral}_${_ansiErr}), буквы в нижнем регистре (${_ansiSecondaryLiteral}a-z${_ansiErr}), цифры (${_ansiSecondaryLiteral}0-9${_ansiErr}) или дефисы (${_ansiSecondaryLiteral}0-9${_ansiErr}) без пробелов и других символов"
}

_prepareHasItemConditionCode() {
  local src="$1"; shift
  local joiner
  conditionCode=
  while [[ $# -gt 0 ]]; do
    conditionCode+="$joiner$src == $(_quotedArgs "$1")"
    joiner=" || "
    shift
  done
}

_throw() {
  local __errOrigin=$(( ${errOrigin:-0} + 1 ))
  local stackOfs=$((__errOrigin + 2))
  #shellcheck disable=SC2046
  return $(_err --showStack $stackOfs "${_ansiCmd}${errOriginName:-${__thisFuncCommand:-${FUNCNAME[$__errOrigin]}}}${_ansiErr} $*")
}

#shellcheck disable=SC2034,SC2016
_codeToCheckParams='
  if [[ $# -gt 0 ]]; then
    local nonExpected
    if [[ $# -eq 1 ]]; then
      nonExpected="параметр ${_ansiPrimaryLiteral}"
    else
      nonExpected="параметров ${_ansiSecondaryLiteral}"
    fi
    nonExpected+="$(_quotedArgs "$@")${_ansiErr}"
    return $(_throw "не ожидает $nonExpected, используйте ${_ansiOutline}${funcName}${kind}Opt${_ansiErr} для задания опций ${_ansiCmd}${FUNCNAME[0]}")
  fi
'

#shellcheck disable=SC2016
_codeToCheckNoArgsInOpt='
  if [[ $# -gt 0 ]]; then
    local gotArg
    if [[ $# -eq 1 ]]; then
      gotArg="получен ${_ansiPrimaryLiteral}"
    else
      gotArg="получены: ${_ansiSecondaryLiteral}"
    fi
    gotArg+=$(_quotedArgs "$@")
    return $(_throw "не ожидает в ${_ansiOutline}${funcName}ParamsOpt${_ansiErr} ни одного аргумента, но $gotArg")
  fi
'

#shellcheck disable=SC2016
_codeToPrepareCodeFileSpec='
  [[ -n $additionalSuffix && ! $additionalSuffix =~ ^\. ]] && additionalSuffix=".$additionalSuffix"
  local __generatedDirSpec=
  fileSpec="${fileSpec:-${BASH_SOURCE[${originalCodeDeep:-1}]}}" codeHolder=_codeToPrepareGeneratedDirSpec eval "$_evalCode"
  eval ${dstVarName:-codeFileSpec}="$__generatedDirSpec/${funcName:-${FUNCNAME[${originalCodeDeep:-1}]}}$additionalSuffix.$codeType$_codeBashExt"
'

_initVarNamesWithOptions2() {
  local varNamesVarName=$1; shift
  local -a varNames=()
  while [[ $# -gt 0 ]]; do
    local varName="$1"
    [[ $varName =~ ^-- ]] && varName=${varName:2}
    dstVarName=varName _kebabCaseToCamelCase "$varName"
    if [[ $varName =~ ^[[:alpha:]_][[:alnum:]_]*$ ]]; then
      varNames+=( "$varName" )
      shift
    else
      err="ожидает, что ${_ansiOutline}имя переменной ${_ansiPrimaryLiteral}$varName${_ansiErr} для опции ${_ansiCmd}$1${_ansiErr} $_mustBeValidVarName"
      return 1
    fi
  done
  eval "$varNamesVarName"'=( "${varNames[@]}" )'
}

_getDeclarationOfVars() {
  local declarationPrefix=local
  local declarationSuffix='""'
  if [[ $1 == '-a' ]]; then
    shift
    local declarationPrefix+=' -a'
    local declarationSuffix='()'
  fi
  local declaration=
  while [[ $# -gt 0 ]]; do
    declaration+=" $prefix$1=$declarationSuffix"
    shift
  done
  if [[ -n $declaration ]]; then
    echo "
  $declarationPrefix$declaration;"
  fi
}

# =============================================================================
