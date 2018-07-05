
# =============================================================================

_resetBash

# =============================================================================

_profileFileName=
if [[ $OSTYPE =~ ^darwin ]]; then
  _profileFileName=".bash_profile"
elif [[ $OSTYPE =~ ^linux ]]; then
  _profileFileName=".bashrc"
else
  return _err "Неожиданный тип OS ${_ansiPrimaryLiteral}$OSTYPE"
fi
_profileFileSpec="$HOME/$_profileFileName"
_profileCorrentionsVerbose=()
_profileCorrections=()

_commonProfileParams=()
_codeToPrepareCommonProfileParams='
  codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
  _commonProfileParams=(
    "${_verbosityParams[@]}"
    "--exportVar/e"
    "arg"
  )
'
# =============================================================================

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareCommonProfileParams eval "$_evalCode"
_profileUnlessParams=(
  '--uninstall/u'
  '--force/f'
  '--checkOnly/c'
  '--didNotChangeReturnCode/r:0..=0'
  "${_commonProfileParams[@]}"
  'unlessCondition:?'
)
_profileUnless() { eval "$_funcParams2"
  _profileBegin
  local didChange
  local -a OPT_verbosityOfExist=(-v none); [[ $verbosity != 'all' ]] || OPT_verbosityOfExist=(-v ok)
  if [[ -z $uninstall ]]; then
    if ( [[ -n $force ]] || [[ -z $unlessCondition ]] || ! _silent "$unlessCondition" ) && ( \
      [[ $verbosity == 'dry' ]] || \
      [[ -n $exportVar ]] || ! _hasLineAtProfile ${OPT_verbosityOfExist[@]} "$arg" || \
    false ); then
      [[ $verbosity == 'dry' ]] || didChange=true
      [[ -n $checkOnly ]] || _addLineToProfile ${OPT_verbosity[@]} ${OPT_silent[@]} ${OPT_exportVar[@]} "$arg" || return $?
    fi
  else
    if \
      [[ $verbosity == 'dry' ]] || \
      ! _hasLineAtProfile -n ${OPT_verbosityOfExist[@]} ${OPT_exportVar[@]} "$arg" || \
    false; then
      [[ $verbosity == 'dry' ]] || didChange=true
      [[ -n $checkOnly ]] || _removeLineFromProfile ${OPT_verbosity[@]} ${OPT_silent[@]} ${OPT_exportVar[@]} "$arg" || return $?
    fi
  fi
  local returnCode=0
  if [[ -z $didChange && -n $didNotChangeReturnCode ]]; then
    returnCode=$didNotChangeReturnCode
  fi
  _profileEnd
  return $returnCode
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareCommonProfileParams eval "$_evalCode"
_addLineToProfileParams=(
  "${_commonProfileParams[@]}"
)
_addLineToProfile() { eval "$_funcParams2"
  _profileBegin
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  _removeLineFromProfile -v err ${OPT_exportVar[@]} "$arg" || return $?
  local profileLine
  if [[ -z $exportVar ]]; then
    profileLine="$arg"
  else
    local varName="$arg"
    local typeSignature=$(declare -p $varName 2>/dev/null)
    local profileLine
    if [[ $typeSignature =~ ^(declare[[:space:]]-(a)?[^[:space:]]*[[:space:]]) ]]; then
      profileLine="export ${typeSignature:${#BASH_REMATCH[1]}}"
    else
      return $(_err --showStack 2 "Переменная ${_ansiOutline}$varName${_ansiErr} не объявлена")
    fi
  fi
  _exec ${OPT_verbosity[@]} --append "$_profileFileSpec" echo "$profileLine" || return $?
  _profileCorrections+=( "$profileLine" )
  # _profileCorrentionsVerbose+=( "Добавлена строка ${_ansiCmd}$profileLine" )
  # _didChangeProfile=true
  _profileEnd
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareCommonProfileParams eval "$_evalCode"
_removeLineFromProfileParams=(
  "${_commonProfileParams[@]}"
)
_removeLineFromProfile() { eval "$_funcParams2"
  _profileBegin
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  local OPT_verbosityOfExist=(-v none); [[ $verbosity != 'all' ]] || OPT_verbosityOfExist=(-v ok)
  if [[ $verbosity == 'dry' ]] || ! _hasLineAtProfile -n ${OPT_verbosityOfExist[@]} ${OPT_exportVar[@]} "$arg"; then
    local -a grepParams=(); _prepareGrepParamsForProfileHelper
    _exec ${OPT_verbosity[@]} --stdout "$_profileFileSpec.removed" --treatAsOK 1 grep "${grepParams[@]}" "$_profileFileSpec"; local returnCode=$?; [[ $returnCode -eq 0 || $returnCode -eq 1 ]] || return $returnCode
    _exec ${OPT_verbosity[@]} --stdout "$_profileFileSpec.new" --treatAsOK 1 grep -v "${grepParams[@]}" "$_profileFileSpec"; local returnCode=$?; [[ $returnCode -eq 0 || $returnCode -eq 1 ]] || return $returnCode
    _mvFile ${OPT_verbosity[@]} ${OPT_silent[@]} "$_profileFileSpec.new" "$_profileFileSpec" || return $?
    if [[ -n $exportVar ]]; then
      _profileCorrections+=( "unset $arg" )
    else
      IFS=$'\r\n' GLOBIGNORE='*' command eval 'local -a lines=( $(cat "$_profileFileSpec.removed") )'
      local line; for line in "${lines[@]}"; do
        _profileCorrentionsVerbose+=( "Удалена строка ${_ansiCmd}$profileLine" )
      done
    fi
  fi
  _profileEnd
}

verbosityDefault=none silentDefault=yes codeHolder=_codeToPrepareCommonProfileParams eval "$_evalCode"
_hasLineAtProfileParams=(
  '--no/n'
  "${_commonProfileParams[@]}"
)
_hasLineAtProfile() { eval "$_funcParams2"
  _profileBegin
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  local -a grepParams=(); _prepareGrepParamsForProfileHelper
  local returnCode=0
  if [[ $verbosity == 'dry' ]]; then
    echo "${_ansiCmd}grep $(_quotedArgs "${grepParams[@]}") \"$_profileFileSpec\"${_ansiReset}"
  else
    _exec -v none ${OPT_silent[@]} grep "${grepParams[@]}" "$_profileFileSpec"; returnCode=$?
    if [[ -n $no ]]; then
      [[ $returnCode -eq 0 ]] && returnCode=1 || returnCode=0
    fi
    if [[ $verbosity != 'none' ]]; then
      local prefix="Строка "
      [[ -z $exportVar ]] \
        && prefix+="${_ansiCmd}$arg" \
        || prefix+="начинающаяся с ${_ansiCmd}export $arg="
      local suffix=' '; [[ -z $no && $returnCode -ne 0 || -n $no && $returnCode -eq 0 ]] && suffix+='не '
      suffix+="обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName"
      if [[ $returnCode -ne 0 ]]; then
        [[ $verbosity == 'ok' ]] || _err "$prefix${_ansiErr}$suffix"
      else
        [[ $verbosity == 'err' ]] || _ok "$prefix${_ansiOK}$suffix"
      fi
    fi
  fi
  _profileEnd
  return $returnCode
}

_prepareGrepParamsForProfileHelper() {
  [[ -z $exportVar ]] && grepParams=( '-F' "$arg" ) || grepParams=( "^export $arg=" )
}

# =============================================================================

verbosityDefault=none silentDefault=yes codeHolder=_codeToPrepareCommonProfileParams eval "$_evalCode"

_exportVarAtBashProfileParams=(
  "${_verbosityParams[@]}"
  '--uninstall/u'
  'varName'
)
_exportVarAtBashProfile() { eval "$_funcParams2"
  local profileLine="export $varName=$(_quotedArgs "${!varName}")"
  local profileLineRegExp="^\s*(export\s+)?BW_SELF_UPDATE_SOURCE="
  _setAtBashProfile "${OPT_verbosity[@]}" "${OPT_silent[@]}" "$profileLine" "$profileLineRegExp"
}

_setAtBashProfileParams=(
  "${_verbosityParams[@]}"
  '--uninstall/u'
  'profileLine'
  'profileLineRegExp'
)
_setAtBashProfile() { eval "$_funcParams2"
  if [[ $verbosity == dry ]]; then
    echo "${_ansiCmd}echo \"$profileLine\" >> \"$_profileFileSpec\""
  else
    if _exec "${OPT_verbosity[@]}" "${OPT_silent[@]}" grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1; then
      [[ -n $uninstall ]] \
        && perlCode="print unless /$profileLineRegExp/" \
        || perlCode="if (! /$profileLineRegExp/) { print } elsif (! \$state) { print $(_quotedArgs --quote:all "$profileLine") . \"\n\"; \$state=1 }"
    elif [[ -z $uninstall ]]; then
      echo "$profileLine" >> "$_profileFileSpec"
    fi
    if [[ -n $perlCode ]]; then
      local newFileSpec="$_profileFileSpec.new"
      _exec "${OPT_verbosity[@]}" --cmdAsIs "cat $(_quotedArgs "$_profileFileSpec") | perl -ne $(_quotedArgs --quote:dollarSign "$perlCode") > $(_quotedArgs "$_profileFileSpec.new")"
      _backupProfileFile
      _mvFile "${OPT_verbosity[@]}" "${OPT_silent[@]}" "$_profileFileSpec.new" "$_profileFileSpec"
    fi
  fi
}

_backupProfileFile() {
  local num=0
  while [[ -f "$_profileFileSpec.bak$num" ]]; do
    num=$(( num + 1 ))
  done
  cp "$_profileFileSpec" "$_profileFileSpec.bak$num"
}

# =============================================================================

