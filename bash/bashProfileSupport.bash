#!/bin/bash
# =============================================================================

_resetBash

# =============================================================================

# shellcheck disable=SC2154
if [[ -n $_isShellCheck ]]; then
  . "bw.bash"
fi

# =============================================================================

_profileFileName=".bashrc"
# if [[ $OSTYPE =~ ^darwin ]]; then
#   _profileFileName=".bash_profile"
# elif [[ $OSTYPE =~ ^linux ]]; then
#   _profileFileName=".bashrc"
# else
#   return _err "Неожиданный тип OS ${_ansiPrimaryLiteral}$OSTYPE"
# fi
_profileFileSpec="$HOME/$_profileFileName"

# =============================================================================
_exportVarAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
# shellcheck disable=SC2016
_exportVarAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--singleQuote/q'
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
  if [[ -n $singleQuote ]]; then
    exactLine="export $varName='${!varName}' # by bw.bash"
  else
    exactLine="export $varName=$(_quotedArgs "${!varName}") # by bw.bash"
  fi
  matchRegexp="^[ \t]*(export[ \t]+)?$varName=.* #[ \t]*by[ \t]+bw.bash[ \t]*$"
}

# shellcheck disable=SC2016
_hasExportVarAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
_hasExportVarAtBashProfileParams=(
  '--profileFileSpec=$_profileFileSpec'
  '--singleQuote/q'
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

_hasAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
_hasAtBashProfileParams=(
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

_setAtBashProfileParamsOpt=(
  --canBeMixedOptionsAndArgs
)
# shellcheck disable=SC2016
_setAtBashProfileParams=(
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

# =============================================================================

_setAtBashProfileHelper() {
  local awkFileSpec; _prepareAwkFileSpec 
  local -a awk_OPT=(
    -f "$awkFileSpec" 
    -v "exactLine=${exactLine?}" 
    -v "matchRegexp=${matchRegexp?}" 
    -v "uninstall=${uninstall?}" 
  )
  awk "${awk_OPT[@]}" "${profileFileSpec?}"
}

# shellcheck disable=SC2016
_backupBashProfileFileParams=(
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

