
# =============================================================================

_resetBash

# =============================================================================

_selfUpdateParams=( 'selfUpdateSource' )
_selfUpdate() { eval "$_funcParams2"
  _download -t -c etag -r 3 "$selfUpdateSource/$_bwFileName" "$_bwFileName"; local returnCode=$?
  _silent cmp "$_bwFileName" "$_bwFileSpec" && returnCode=0
  if [[ $returnCode -eq 3 ]]; then
    _getBwTar "$_bwFileName" | tar xf -  || return $(_err "Не удалось распаковать ${_ansiFileSpec}$_bwDir/$_buildFileName")
    cp "$_bwFileName" "$_bwFileSpec" || return $(_err "Не удалось заменить ${_ansiFileSpec}$_bwFileSpec${_ansiErr} новой версией")
    _bwInstalled=
  fi
  return $returnCode
}

_defaultBwUpdateSource="$_bwGithubSource/master"
_bwMainParamsOpt=(--canBeMoreParams)
_bwMainParams=(
  '--pregenOnly/p:?=$BW_PREGEN_ONLY'
  '--force/f'
  '--noSelfUpdate/n'
  '--selfUpdateSource/u=${BW_SELF_UPDATE_SOURCE:-'"$_defaultBwUpdateSource"'}'
)
_bwMain_pregenOnly_description='Ограничивает прегенерацию только указанными функциями, значение ${_ansiPrimaryLiteral}-${_ansiReset} имеет смысл "отключить прегенерацию"'
_bwMain_force_description='Форсирует прегенерацию, независимо от значения ${_ansiOutline}_isBwDevelop${_ansiReset}'
_bwMain_noSelfUpdate_description='Блокирует самобновления из источника обновления'
_bwMain_selfUpdateSource_description='Устанавливает URL источника обновления'
# _sourceMatchRegexp='^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~_-]+)\/'
# _sourceMatchRegexp='^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~_-]+)\/'
# _bwMatchRegexp="$_sourceMatchRegexp"'bw\.bash'
# _bwMatchRegexp='^[ \t]*.[ \t]+"?([ a-zA-Z0-9/~_-]+)/bw.bash'
_bwMatchRegexp='^[ \t]*\\.[ \t]+"?([ a-zA-Z0-9/~_-]+)/bw.bash'
_bwMain() { eval "$_funcParams2"
  _profileBegin

  local returnCode=0
  while true; do

    [[ ! $selfUpdateSource =~ ^- ]] || selfUpdateSource="$_defaultBwUpdateSource"

    if [[ -n $_isBwDevelop || -n $_isBwDevelopInherited ]] ; then
      selfUpdateSource=$(_inDir "$_bwDir" _gitBranch) || { returnCode=$?; break; }
      _export_BW_SELF_UPDATE_SOURCE
    elif [[ -z $noSelfUpdate ]]; then
      _inDir --treatAsOK 3 --preserveReturnCode "$_bwDir" _selfUpdate "$selfUpdateSource"; local returnCode=$?
      _export_BW_SELF_UPDATE_SOURCE
      if [[ $returnCode -eq 3 ]]; then
        . "$_bwFileSpec" "$@"; local returnCode=$?
        [[ $returnCode -eq 0 ]] && _ok "${_ansiFileSpec}$_bwFileName${_ansiOK} обновлен до версии ${_ansiPrimaryLiteral}$(bw_version)"
        break
      elif [[ $returnCode -eq 0 && -n $_bwInstalled ]]; then
        _ok "Установлен ${_ansiFileSpec}$_bwFileName${_ansiOK} версии ${_ansiPrimaryLiteral}$(bw_version)"
        _bwInstalled=
      elif [[ $returnCode -ne 0 ]]; then
        break
      fi
    fi

    if [[ ! $pregenOnly =~ ^- && ( -n $_isBwDevelop || -n $OPT_force ) ]]; then
      local -a __completions=()
      local generatedCompletionsFileSpec="/tmp/bw.generated.completions.bash"
      _rm "$_generatedCompletionsFileSpec" || { returnCode=$?; break; }
      _spinner \
        -t "${_ansiOK}OK: ${_ansiCmd}$_bwFileSpec${_ansiReset} обработан за" \
        "${_ansiHeader}Прегенерация${_ansiReset}" \
        _bwMainHelper \
        || { returnCode=$?; break; }
      [[ ! -f $generatedCompletionsFileSpec ]] || . "$generatedCompletionsFileSpec"
      local fileSpec; for fileSpec in "${_completions[@]}"; do
        . "$fileSpec"
      done
    fi

    local exactLine=". $(_quotedArgs "$(_shortenFileSpec "$_bwFileSpec")" -n)"
    if [[ -n $_isBwDevelop || -n $_isBwDevelopInherited ]]; then
      exactLine+=" -p -"
    fi
    _setAtBashProfile "$exactLine" "$_bwMatchRegexp"

    break
  done
  _profileEnd
  if [[ $returnCode -ne 0 ]]; then
    return $returnCode
  elif [[ $# -gt 0 ]]; then
    eval "$@"
  fi
}
_export_BW_SELF_UPDATE_SOURCE() {
  if [[ $selfUpdateSource == "$_defaultBwUpdateSource" ]]; then
    export BW_SELF_UPDATE_SOURCE=
  else
    export BW_SELF_UPDATE_SOURCE="$selfUpdateSource"
  fi
  _exportVarAtBashProfile BW_SELF_UPDATE_SOURCE
}

_bwMainHelper() {
  _profileBegin
  local -a __completions=();
  _pregen ${pregenOnly:-$(declare -F | perl -pe "s/^declare -f //")} \
    && echo "__completions=( $(_quotedArgs "${__completions[@]}") )" > "$generatedCompletionsFileSpec"
  _profileEnd
}

# =============================================================================

_pregen() {
  local funcName; for funcName in "$@"; do
    local needPrepare=
    if ! [[ $funcName =~ (Complete|Params)$ || $funcName =~ ^_debug ]]; then
      local holder; for holder in \
        ${funcName}Params \
        ${funcName}BoolOptions \
        ${funcName}ScalarOptions \
        ${funcName}ListOptions \
      ; do
        [[ $holder =~ Params$ ]] && _funcExists $holder && needPrepare=Params2 && break
        if [[ $(declare -p $holder 2>/dev/null) =~ ^declare[[:space:]]-a ]]; then
          [[ $holder =~ Params$ ]] && needPrepare=Params2 || needPrepare=Options2
          break
        fi
      done
    fi
    if [[ -n $needPrepare ]]; then
      local hasWrapper=
      [[ ! $funcName =~ [^_]_.* ]] || hasWrapper=true
      codeHolder=_codeToPregen eval "$_evalCode"
    fi
  done
}

# =============================================================================
