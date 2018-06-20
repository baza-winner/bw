
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

# _codeOfSelfUpdateSource='${BW_SELF_UPDATE_SOURCE:-$_bwGithubSource/master}'
_defaultBwUpdateSource="$_bwGithubSource/master"
_bwMainParamsOpt=(--canBeMoreParams)
_bwMainParams=(
  '--pregenOnly/p:?=$BW_PREGEN_ONLY'
  '--force/f'
  '--selfUpdateSource/u=${BW_SELF_UPDATE_SOURCE:-'"$_defaultBwUpdateSource"'}'
)
_bwMain_pregenOnly_description="Ограничивает прегенерацию только указанными функциями, значение ${_ansiPrimaryLiteral}-${_ansiReset} имеет смысл \"отключить прегенерацию\""
_bwMain_force_description="Форсирует прегенерацию, независимо от значения ${_ansiOutline}_isBwDevelop${_ansiReset}"
_bwMain_selfUpdateSource_description="Устанавливает URL источника обновления"
_bwMain() { eval "$_funcParams2"
  _profileBegin
  local returnCode=0
  while true; do
    [[ ! $selfUpdateSource =~ ^- ]] || selfUpdateSource="$_bwGithubSource/master"

    if [[ -z $_isBwDevelop && -z $_isBwDevelopInherited ]] ; then
      _inDir --treatAsOK 3 --preserveReturnCode "$_bwDir" _selfUpdate "$selfUpdateSource"; local returnCode=$?
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
    else
      _inDir "$_bwDir" _prepareGitBranchName \
        || { returnCode=$?; break; }
      selfUpdateSource="$gitBranchName"
    fi

    if [[ ! $pregenOnly =~ ^- && ( -n $_isBwDevelop || -n $OPT_force ) ]]; then
      _spinner \
        -t "${_ansiOK}OK: ${_ansiCmd}$_bwFileSpec${_ansiReset} обработан за" \
        "${_ansiHeader}Прегенерация${_ansiReset}" \
        _bwMainHelper \
        || { returnCode=$?; break; }
    fi

    if [[ $selfUpdateSource == $_defaultBwUpdateSource ]]; then
      export BW_SELF_UPDATE_SOURCE=
    else
      export BW_SELF_UPDATE_SOURCE="$selfUpdateSource"
    fi

    _exportVarAtBashProfile BW_SELF_UPDATE_SOURCE

    local profileLine=". $(_quotedArgs "$(_shortenFileSpec "$_bwFileSpec")")"
    if [[ -n $_isBwDevelop || -n $_isBwDevelopInherited ]]; then
      profileLine+=" -p -"
    fi
    _setAtBashProfile "$profileLine" "^\s*\.\s+\"?(.+?)\/bw\.bash\"?"

    _cmdToExecute=( "$@" )
    break
  done
  _profileEnd
  return $returnCode
}

_prepareGitBranchName() {
  gitBranchName=$(_gitBranch)
}

_bwMainHelper() {
  _profileBegin
  local -a __completions=();
  _pregen ${pregenOnly:-$(declare -F | perl -pe "s/^declare -f //")} \
    && echo "_completions=( $(_quotedArgs "${__completions[@]}") )" > "$_generatedCompletionsFileSpec"
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
