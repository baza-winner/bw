#!/bin/bash

{
  # shellcheck disable=SC2154
  if [[ -n $shellcheck ]]; then 
    . bw.bash
    . core/profileSupport.bash
    . core/coreUtils.bash
    . core/funcOptionsSupport2.bash
  fi
}

# =============================================================================

_resetBash

# =============================================================================

# shellcheck disable=SC2034
{
_debugOptionsOpt=( --treatUnknownOptionAsArg --prefix ___ )
_debugBoolOptions=( 'clean' )
}
_debug() { __onlyPrepareCode='' eval "$_funcOptions2"
  # shellcheck disable=SC2154
  if [[ -n $___clean ]]; then
    echo "$@" >&2
  else
    echo -n "${_ansiBlue}${_ansiBold}${FUNCNAME[1]}:${_ansiReset} $*" >&2
    local ___joiner=', '
    _debugStack
  fi
}

# shellcheck disable=SC2034
{
_debugVarOptionsOpt=( --treatUnknownOptionAsArg --prefix ___  )
_debugVarBoolOptions=( 'clean' 'noQuotedArgs' 'asis' 'endWithJoiner')
_debugVarScalarOptions=( 'joiner' 'mark' )
}
_debugVar() { __onlyPrepareCode='' eval "$_funcOptions2"
  # shellcheck disable=SC2154
  [[ -n $___asis ]] && ___noQuotedArgs=true
  # shellcheck disable=SC2154
  [[ -n $___mark ]] && echo -n "${_ansiYellow}${_ansiBold}$___mark${_ansiReset} " >&2
  # shellcheck disable=SC2154
  [[ -z $___clean ]] && echo -n "${_ansiDebug}${_ansiBold}${FUNCNAME[1]}:${_ansiReset} " >&2
  local ___varName ___needJoiner; for ___varName in "$@"; do
    # shellcheck disable=SC2059
    [[ -n $___needJoiner ]] && printf "$___joiner" >&2
    echo -n "${_ansiOutline}$___varName${_ansiReset}" >&2
    local ___typeOfVar; ___typeOfVar=$(_getTypeOfVar "$___varName")
    if [[ $___typeOfVar == 'array' ]]; then
      dstVarName=___varValue srcVarName=$___varName codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      # shellcheck disable=SC2154
      echo -n "(${#___varValue[@]})" >&2
      if [[ ${#___varValue[@]} -gt 0 ]]; then
        if [[ -z $___noQuotedArgs ]]; then
          echo -n ": ${_ansiSecondaryLiteral}$(_quotedArgs "${___varValue[@]}")${_ansiReset}" >&2
        else
          echo -n ": ${_ansiSecondaryLiteral}${___varValue[*]})${_ansiReset}" >&2
        fi
      fi
    elif [[ $___typeOfVar == 'none' ]]; then
      echo -n "<${_ansiErr}is unset${_ansiReset}>" >&2
    elif [[ -z ${!___varName} ]]; then
      echo -n "<empty>" >&2
    elif [[ -z $___noQuotedArgs ]]; then
      echo -n ": ${_ansiPrimaryLiteral}$(_quotedArgs "${!___varName}")${_ansiReset}" >&2
    else
      echo -n ": ${_ansiPrimaryLiteral}${!___varName}${_ansiReset}" >&2
    fi
    [[ -z $___joiner ]] && ___joiner=", "
    ___needJoiner=true
  done
  # shellcheck disable=SC2154,SC2059
  [[ -n $___endWithJoiner ]] && printf "$___joiner" >&2
  if [[ -n $___clean ]]; then
    echo >&2
  else
    _debugStack
  fi
}

# shellcheck disable=SC2034
_okBoolOptions=( 'stderr' )
_ok() { eval "$_funcOptions2"
  local msg="${_ansiOK}OK: $*${_ansiReset}"
  # shellcheck disable=SC2154
  if [[ -n $stderr ]]; then
    echo "$msg" >&2
  else
    echo "$msg"
  fi
  return 0
}

# shellcheck disable=SC2034
_errScalarOptions=( 'showStack' 'prefix' )
_err() {
  local returnCode=$?; [[ $returnCode -eq 0 ]] && returnCode=1
  eval "$_funcOptions2"

  local opt 
  # shellcheck disable=SC2154
  [[ -n $showStack && -z $_noStack ]] && opt=-n

  echo $opt "$prefix${_ansiErr}ERR: $*${_ansiReset}" >&2
  [[ -n $showStack && -z $_noStack ]] && ___joiner=', ' _debugStack "$showStack"
  return $returnCode
}

# =============================================================================

