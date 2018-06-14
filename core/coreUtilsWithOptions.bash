
# =============================================================================

_resetBash

# =============================================================================


_debugOptionsOpt=( --treatUnknownOptionAsArg --prefix ___ )
_debugBoolOptions=( 'clean' )
_debug() { __onlyPrepareCode= eval "$_funcOptions2"
  if [[ -n $___clean ]]; then
    echo "$@" >&2
  else
    echo -n "${_ansiBlue}${_ansiBold}${FUNCNAME[1]}:${_ansiReset} $@" >&2
    local ___joiner=', '
    _debugStack
  fi
}

_debugVarOptionsOpt=( --treatUnknownOptionAsArg --prefix ___  )
_debugVarBoolOptions=( 'clean' 'noQuotedArgs' 'asis' 'endWithJoiner')
_debugVarScalarOptions=( 'joiner' 'mark' )
_debugVar() { __onlyPrepareCode= eval "$_funcOptions2"
  [[ -n $___asis ]] && ___noQuotedArgs=true
  [[ -n $___mark ]] && echo -n "${_ansiYellow}${_ansiBold}$___mark${_ansiReset} " >&2
  [[ -z $___clean ]] && echo -n "${_ansiDebug}${_ansiBold}${FUNCNAME[1]}:${_ansiReset} " >&2
  local ___varName ___needJoiner; for ___varName in $@; do
    [[ -n $___needJoiner ]] && printf "$___joiner" >&2
    echo -n "${_ansiOutline}$___varName${_ansiReset}" >&2
    local ___typeSignature=$(declare -p "$___varName" 2>/dev/null)
    if [[ "$___typeSignature" =~ ^declare[[:space:]]-a ]]; then
      dstVarName=___varValue srcVarName=$___varName codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      echo -n "(${#___varValue[@]})" >&2
      if [[ ${#___varValue[@]} -gt 0 ]]; then
        if [[ -z $___noQuotedArgs ]]; then
          echo -n ": ${_ansiSecondaryLiteral}$(_quotedArgs "${___varValue[@]}")${_ansiReset}" >&2
        else
          echo -n ": ${_ansiSecondaryLiteral}${___varValue[@]})${_ansiReset}" >&2
        fi
      fi
    elif [[ -z ${!___varName+x} ]]; then
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
  [[ -n $___endWithJoiner ]] && printf "$___joiner" >&2
  if [[ -n $___clean ]]; then
    echo >&2
  else
    _debugStack
  fi
}

_okBoolOptions=( 'stderr' )
_ok() { eval "$_funcOptions2"
  local msg="${_ansiOK}OK: $@${_ansiReset}"
  if [[ -n $stderr ]]; then
    echo "$msg" >&2
  else
    echo "$msg"
  fi
  return 0
}

_errScalarOptions=( 'showStack' 'prefix' )
_err() {
  local returnCode=$?; [[ $returnCode -eq 0 ]] && returnCode=1
  eval "$_funcOptions2"
  local opt; [[ -n $showStack && -z $noStack ]] && opt=-n
  echo $opt "$prefix${_ansiErr}ERR: $@${_ansiReset}" >&2
  [[ -n $showStack && -z $noStack ]] && ___joiner=', ' _debugStack $showStack
  return $returnCode
}

# =============================================================================

