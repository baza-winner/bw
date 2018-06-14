
# =============================================================================

_resetBash

# =============================================================================

_verboseCache=
# _verboseCache='--verbosity:all'
_codeToUseCache='_useCache "$varName"; local returnCode=$?; [[ $returnCode -eq 2 ]] || return $returnCode'
_useCacheParams=( 'varName' )
_useCache() { eval "$_funcParams2"
  local msg returnCode
  local codeFileSpec; dstVarName= codeType=cache additionalSuffix= fileSpec= originalCodeDeep= codeHolder=_codeToPrepareCodeFileSpec eval "$_evalCode"
  if [[ ! -f $codeFileSpec ]]; then
    msg="${_ansiWarn}WARN: Кеш ${_ansiCmd}$codeFileSpec${_ansiWarn} не найден${_ansiReset}"
    returnCode=2
  elif [[ -n $_isBwDevelop ]] && ! _everyFileNotNewerThan "$codeFileSpec" "${BASH_SOURCE[0]}" "${BASH_SOURCE[1]}"; then
    msg="${_ansiWarn}WARN: Кеш ${_ansiCmd}$codeFileSpec${_ansiWarn} устарел${_ansiReset}"
    _rm $codeFileSpec || return $?
    returnCode=2
  elif ! grep "^$varName=" "$codeFileSpec" >/dev/null 2>&1; then
    msg="${_ansiWarn}WARN: Кеш ${_ansiCmd}$codeFileSpec${_ansiWarn} не содержит значение для ${_ansiOutline}$varName${_ansiReset}"
    returnCode=2
  else
    _fileSpec="$codeFileSpec" codeHolder=_codeSource eval "$_evalCode"
    msg="${_ansiOK}OK: Использован кеш ${_ansiCmd}$codeFileSpec${_ansiOK} для ${_ansiOutline}$varName${_ansiReset}"
    returnCode=0
  fi
  [[ -n $_verboseCache ]] && echo "$msg" >&2
  return $returnCode
}

_saveToCacheParamsOpt=(--canBeMoreParams)
_saveToCacheParams=( '--array/a' '--debug/d' 'varName' )
_saveToCache() { eval "$_funcParams2"
  local result
  if [[ $# -eq 0 ]]; then
    local typeSignature=$(declare -p $varName 2>/dev/null)
    if [[ $typeSignature =~ ^declare[[:space:]]-a ]]; then
      result=$(echo "$typeSignature" | perl -pe "s/^declare -\S* $varName='?/$varName=/; s/'$//")
    elif [[ -z ${!varName+x} ]]; then
      echo "${_ansiErr}ERR: Переменная ${_ansiOutline}$varName${_ansiErr} не определена${_ansiReset}" >&2 && return 1
    elif [[ ${!varName} =~ \' ]]; then
      result=$(echo "$typeSignature" | perl -pe "s/^declare -\S* $varName=/$varName=/")
    else
      result="$varName='${!varName}'"
    fi
  elif [[ -n $array ]]; then
    local -a varValue=( "$@" )
    result=$(declare -p varValue | perl -pe "s/^declare -\S* varValue='?/$varName=/; s/'$//")
  else
    [[ $# -gt 1 ]] && echo "${_ansiErr}ERR: ${_ansiCmd}${FUNCNAME[0]}${_ansiErr} ожидает не более одного значения для ${_ansiOutline}$varName${_ansiErr}, но получены $#: ${_ansiPrimaryLiteral}$(_quotedArgs "$@")${_ansiReset}" >&2 && return 1
    local varValue="$1"
    if [[ $varValue =~ \' ]]; then
      result=$(declare -p varValue | perl -pe "s/^declare -\S* varValue=/$varName=/")
    else
      result="$varName='$varValue'"
    fi
  fi
  if [[ -n $debug ]]; then
    echo "$result"
  else
    local codeFileSpec; dstVarName= codeType=cache additionalSuffix= fileSpec= originalCodeDeep= codeHolder=_codeToPrepareCodeFileSpec eval "$_evalCode"
    if [[ -f $codeFileSpec ]] && ! _everyFileNotNewerThan "$codeFileSpec" "${BASH_SOURCE[0]}"; then
      _rm "$codeFileSpec" || return $?
    fi
    _assureDir $(dirname "$codeFileSpec") || return $?
    echo "$result" >> "$codeFileSpec"
    [[ -z $_verboseCache ]] || echo "${_ansiOK}OK: Сохранено в кеш ${_ansiFileSpec}$codeFileSpec${_ansiOK} значение ${_ansiOutline}$varName${_ansiReset}" >&2
  fi
}

_rmCacheParams=()
_rmCache() { eval "$_funcParams2"
  local codeFileSpec; dstVarName= codeType=cache additionalSuffix= fileSpec= originalCodeDeep= codeHolder=_codeToPrepareCodeFileSpec eval "$_evalCode"
  _rm $_verboseCache "$codeFileSpec"
}

# =============================================================================

