#!/bin/bash
# =============================================================================

_resetBash

# =============================================================================

# shellcheck disable=SC2154
if [[ -n $_isShellCheck ]]; then
  . "bw.bash"
fi

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

# =============================================================================
export _exportVarAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
# shellcheck disable=SC2016
export _exportVarAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--noBackup'
  '--uninstall/u'
  'varName'
)
# shellcheck disable=SC2154
_exportVarAtBashProfile() { eval "$_funcParams2"
  local exactLine matchRegexp 
  _exportVarAtBashProfileHelper
  local -a OPT=(
    "${OPT_profileFileSpec[@]}" 
    "${OPT_noBackup[@]}" 
    "${OPT_uninstall[@]}"
  )
  _setAtBashProfile "${OPT[@]}" "$exactLine" "$matchRegexp"
}
_exportVarAtBashProfileHelper() {
  exactLine="export $varName=$(_quotedArgs "${!varName}") # by bw.bash"
  # matchRegexp="^[[:space:]]*(export[[:space:]]+)?$varName="
  # matchRegexp="^[ \t]*(export[ \t]+)?$varName="
  matchRegexp="^[ \t]*(export[ \t]+)?$varName=.* #[ \t]*by[ \t]+bw.bash[ \t]*$"
}

# shellcheck disable=SC2016
export _hasExportVarAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
export _hasExportVarAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--differ/d'
  '--no/n'
  'varName'
)
# shellcheck disable=SC2154
_hasExportVarAtBashProfile() { eval "$_funcParams2"
  local exactLine matchRegexp 
  _exportVarAtBashProfileHelper
  local -a OPT=(
    "${OPT_profileFileSpec[@]}" 
    "${OPT_differ[@]}" 
    "${OPT_no[@]}"
  )
  _hasAtBashProfile "${OPT[@]}" "$exactLine" "$matchRegexp"
}

export _hasAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
export _hasAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--differ/d'
  '--no/n'
  'exactLine'
  'matchRegexp'
)
_hasAtBashProfile() { eval "$_funcParams2"
  local uninstall=''
  _setAtBashProfileHelper >/dev/null; local returnCode=$?
  local treatAsOK=0
  if [[ -n $no ]]; then
    treatAsOK=1
    if [[ $returnCode -eq $treatAsOK ]]; then
      returnCode=0
    elif [[ $returnCode -eq 0 ]]; then
      returnCode="$treatAsOK"
    fi
  fi
  if [[ $returnCode -eq 2 && -z $differ ]]; then
    returnCode="$treatAsOK"
  fi
  return $returnCode
}

export _setAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
# shellcheck disable=SC2016
export _setAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--noBackup'
  '--uninstall/u'
  'exactLine'
  'matchRegexp'
)
# shellcheck disable=SC2154
_setAtBashProfile() { eval "$_funcParams2"
  if [[ $verbosity == dry ]]; then
    echo "${_ansiCmd}echo \"$exactLine\" >> \"$profileFileSpec\""
  else
    local code=''
    local innerCode=''
    local notFound=
    if [[ ! -f "$profileFileSpec" ]]; then
      notFound=true
    else
      local newFileSpec="$profileFileSpec.new"
      _setAtBashProfileHelper > "$newFileSpec"; local returnCode=$?
      if [[ $returnCode -eq 1 ]]; then
        notFound=true
      elif [[ $returnCode -eq 2 ]]; then
        [[ -n $noBackup ]] || _backupBashProfileFile "$profileFileSpec"
        mv "$newFileSpec" "$profileFileSpec"
      fi
    fi
    if [[ -n $notFound && -z $uninstall ]]; then
      [[ -n $noBackup ]] || _backupBashProfileFile "$profileFileSpec"
      echo "$exactLine" >> "$profileFileSpec"
    fi
  fi
}

export _getAlreadyProjDirParamsOpt=(
  '--canBeMixedOptionsAndArgs'
)
export _getAlreadyProjDirParams=( 
  '--profileFileSpec=$_profileFileSpec'
  '--varName='
  'bwProjShortcut' 
)
_getAlreadyProjDir() { eval "$_funcParams2"
  [[ -z $varName ]] || eval "$varName"'=""'
  if [[ -z $varName ]]; then
    _getAlreadyProjDirHelper
  else
    eval "$varName"'="$(_getAlreadyProjDirHelper)"'
  fi
}
_getAlreadyProjDirHelper() {
  # local sedCode='s/^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~]+)\/bin\/'"${bwProjShortcut?}"'\.bash''.*$/\1/p'
  # _debugVar sedCode
  local sedCode='s/'"$_sourceMatchRegexp"'bin\/'"${bwProjShortcut?}"'\.bash''.*$/\1/p'
  # _debugVar sedCode
  # local sedCode='s/^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~].+)\/bin\/'"${bwProjShortcut?}"'\.bash''.*$/\1/p'
  sed -n -E -e "$sedCode" "${profileFileSpec?}" | tail -n 1
}

# =============================================================================

_setAtBashProfileHelper() {
  local awkFileSpec; _prepareAwkFileSpec 
  local -a OPT=(
    -f "$awkFileSpec" 
    -v "exactLine=${exactLine?}" 
    -v "matchRegexp=${matchRegexp?}" 
    -v "uninstall=${uninstall?}" 
  )
  awk "${OPT[@]}" "${profileFileSpec?}"
}

_prepareAwkFileSpecParams=( 'infix:?' )
_prepareAwkFileSpec() { eval "$_funcParams2"
  [[ -z $infix ]] || infix=".$infix"
  awkFileSpec="$(dirname "${BASH_SOURCE[1]}")/${FUNCNAME[1]}$infix.awk"
}

# shellcheck disable=SC2016
export _backupBashProfileFileParams=(
  'profileFileSpec=$_profileFileSpec'
)
_backupBashProfileFile() { eval "$_funcParams2"
  if [[ -f "$profileFileSpec" ]]; then
    local num=0
    while [[ -f "$profileFileSpec.bak$num" ]]; do
      num=$(( num + 1 ))
    done
    cp "$profileFileSpec" "$profileFileSpec.bak$num"
  fi
}

# =============================================================================

