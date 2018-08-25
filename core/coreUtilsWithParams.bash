
# =============================================================================

_resetBash

# =============================================================================

_getPluralWordDescription='
  Выводит в stdout форму слова, в зависимости от числа
'
_getPluralWordDescriptionOfCount='число, от которого зависит форма слова'
_getPluralWordDescriptionOfWord1='единственная форма слова'
_getPluralWordDescriptionOfWord2_4='форма слова для 2..4'
_getPluralWordDescriptionOfWord5more='множественная форма слова'
_getPluralWordParams=( 'count:0..' 'word1' 'word2_4' 'word5more=$word2_4' )
_getPluralWord() { eval "$_funcParams2"
  local result="$word5more"
  local decimal=$(( count / 10 % 10 ))
  if [[ $decimal != 1 ]]; then
    local unit=$(( count % 10 ))
    if [[ $unit -eq 1 ]]; then
      result="$word1"
    elif [[ $unit -ge 2 && $unit -le 4 ]]; then
      result="$word2_4"
    fi
  fi
  echo $result
}

_substituteParamsOpt=(--canBeMoreParams)
_substituteParams=( '--silent/s' 'varName' )
_substitute() { eval "$_funcParams2"
  local typeOfVar; eval "$_codeToPrepareTypeOfVar"
  local typeOfVar; typeOfVar=$(_getTypeOfVar "$varName")
  if [[ $typeOfVar == 'none' ]]; then
    if [[ -z $silent ]]; then
      _throw "could not resolve type of ${_ansiOutline}$varName${_ansiErr}, first declare it with initial value"
    fi
    return 1
  else
    local idxVarName="${_substitutePrefix}${varName}${_substituteIdxSuffix}"
    if [[ -z ${!idxVarName} ]]; then
      eval "$idxVarName=0"
    else
      eval "$idxVarName=\$(( $idxVarName + 1 ))"
    fi
    local valueToRestoreVarName="${_substitutePrefix}${varName}${_substituteValueSuffix}${!idxVarName}"
    if [[ $typeOfVar == 'scalar' ]]; then
      eval $valueToRestoreVarName=\$$varName
      eval "$varName=\"$1\""
    else
      eval $valueToRestoreVarName=\( \"\${$varName[@]}\" \)
      eval $varName=\( \"\$@\" \)
    fi
  fi
}

# _shortenFileSpecParams=( 'fileSpec' )
_shortenFileSpec() { 
  # eval "$_funcParams2"
  local fileSpec="$*"
  [[ -z $_isBwDevelop ]] || _debugVar fileSpec 2>> "$_bwDir/_shortenFileSpec.log"
  local lcp=$(_lcp "$fileSpec" "$HOME")
  local result
  if [[ ${#lcp} -ne ${#HOME} ]]; then 
    result="$fileSpec"
  else
    result="~/${fileSpec:$(( ${#HOME} + 1 ))}"
  fi
  echo "$result"
}

# =============================================================================

_prepareAwkFileSpecParams=( 'infix:?' )
_prepareAwkFileSpec() { eval "$_funcParams2"
  [[ -z $infix ]] || infix=".$infix"
  awkFileSpec="$(dirname "${BASH_SOURCE[1]}")/${FUNCNAME[1]}$infix.awk"
}

# =============================================================================
