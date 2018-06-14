
# =============================================================================

_resetBash

# =============================================================================

_inputrcFileName=".inputrc"
_inputrcFileSpec="$HOME/$_inputrcFileName"

# =============================================================================

_inputrcSetPropsParams=(
  '--didNotChangeReturnCode/r:0..'
	'--editingMode:(vi emacs)'
  '--showModeInPrompt:(on off)'
  '--viCmdModeString='
  '--viInsModeString='
)
_inputrcSetProps() { eval "$_funcParams2"
  local returnCode=0
  local -a optVarNames=()
  local optVarName; for optVarName in ${__optVarNames[@]}; do
    if [[ -n ${!optVarName} ]]; then
      local propNameHolder=__OPTNAME_$optVarName
      local propName="${!propNameHolder}"
      if [[ -f $_inputrcFileSpec ]]; then
        local propValue="${!optVarName}"
        if [[ $propValue =~ [[:space:]] ]]; then
          propValue="\"$propValue\""
        fi
        if grep -F "set ${!propNameHolder:2} $propValue" "$_inputrcFileSpec" >/dev/null 2>&1; then
          continue
        fi
      fi
      optVarNames+=( "$optVarName" )
    fi
  done
  if [[ ${#optVarNames[@]} -gt 0 ]]; then
    local newFileSpec="$_inputrcFileSpec.new"
    if [[ -f $_inputrcFileSpec ]]; then
      local propNameRegExp=
      local optVarName; for optVarName in "${optVarNames[@]}"; do
        [[ -z $propNameRegExp ]] || propNameRegExp+='|'
        local propNameHolder=__OPTNAME_$optVarName
        propNameRegExp+="${!propNameHolder:2}"
      done
      grep -v -E "^set ($propNameRegExp) " "$_inputrcFileSpec" >"$newFileSpec"
    elif [[ -f $newFileSpec ]]; then
      rm "$newFileSpec"
    fi
    local optVarName; for optVarName in "${optVarNames[@]}"; do
      local propNameHolder=__OPTNAME_$optVarName
      local propValue="${!optVarName}"
      if [[ $propValue =~ [[:space:]] ]]; then
        propValue="\"$propValue\""
      fi
      echo "set ${!propNameHolder:2} $propValue" >> "$newFileSpec"
    done
    mv "$newFileSpec" "$_inputrcFileSpec"
    bind -f  "$_inputrcFileSpec"
  elif [[ -n $didNotChangeReturnCode ]]; then
    returnCode=$didNotChangeReturnCode
  fi
  return $returnCode
}
