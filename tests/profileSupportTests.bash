
# =============================================================================

_resetBash

# =============================================================================

_profileUnlessTests=(
  '
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --before "local BW_TEST_VAR=\"bw test var value\""
    "_profileUnless --exportVar BW_TEST_VAR && _hasLineAtProfile -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
  '
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --before "local BW_TEST_VAR=\"bw test var value\""
    --stdout "${_ansiCmd}echo \"export BW_TEST_VAR=\\\"bw test var value\\\"\" >>$(_quotedArgs "$HOME/$_profileFileName")${_ansiReset}"
    "_profileUnless --exportVar -v dry BW_TEST_VAR && _hasLineAtProfile -n -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
  '
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --before "local BW_TEST_VAR=\"bw test var value\""
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiOK}OK: ${_ansiCmd}echo \"export BW_TEST_VAR=\\\"bw test var value\\\"\" >>$(_quotedArgs "$HOME/$_profileFileName")${_ansiReset}"
    "_profileUnless --exportVar -v all BW_TEST_VAR && _hasLineAtProfile -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
  '
    --before "local BW_TEST_VAR=\"bw test var value\""
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiOK}OK: ${_ansiCmd}echo \"export BW_TEST_VAR=\\\"bw test var value\\\"\" >>$(_quotedArgs "$HOME/$_profileFileName")${_ansiReset}"
    "_profileUnless --exportVar -v all BW_TEST_VAR && _hasLineAtProfile -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
  '
    --before "local BW_TEST_VAR=\"bw test var value\""
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiCmd}echo \"export BW_TEST_VAR=\\\"bw test var value\\\"\" >>$(_quotedArgs "$HOME/$_profileFileName")${_ansiReset}"
    "_profileUnless --exportVar -v dry BW_TEST_VAR && _hasLineAtProfile -n -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
  '
    --before "local BW_TEST_VAR=\"bw test var value\""
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_profileUnless --exportVar -u BW_TEST_VAR && _hasLineAtProfile -n -v err '$_stq'export BW_TEST_VAR='$_stqq'bw test var value'$_stqq$_stq'"
  '
)

_hasLineAtProfileTests=(
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e BW_TEST_VAR"
  '
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiCmd}grep \"^export BW_TEST_VAR=\" \"$HOME/$_profileFileName\"${_ansiReset}"
    "_hasLineAtProfile -e -v dry BW_TEST_VAR"
  '
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -v none BW_TEST_VAR"
  '
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -v err BW_TEST_VAR"
  '
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiOK}OK: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiOK} обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -v all BW_TEST_VAR"
  '

  '
    --return "1"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -n BW_TEST_VAR"
  '
  '
    --return "0"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiCmd}grep \"^export BW_TEST_VAR=\" \"$HOME/$_profileFileName\"${_ansiReset}"
    "_hasLineAtProfile -e -n -v dry BW_TEST_VAR"
  '
  '
    --return "1"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -n -v none BW_TEST_VAR"
  '
  '
    --return "1"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stderr "${_ansiErr}ERR: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiErr} обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -n -v err BW_TEST_VAR"
  '
  '
    --return "1"
    --before "local BW_TEST_VAR=bwTestVarValue"
    --before "_profileUnless --exportVar BW_TEST_VAR"
    --after "_profileUnless --exportVar -u BW_TEST_VAR"
    --stderr "${_ansiErr}ERR: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiErr} обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -n -v all BW_TEST_VAR"
  '

  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -n BW_TEST_VAR"
  '
  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiCmd}grep \"^export BW_TEST_VAR=\" \"$HOME/$_profileFileName\"${_ansiReset}"
    "_hasLineAtProfile -e -n -v dry BW_TEST_VAR"
  '
  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -n -v none BW_TEST_VAR"
  '
  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -n -v err BW_TEST_VAR"
  '
  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiOK}OK: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiOK} не обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -n -v all BW_TEST_VAR"
  '

  '
    --return "1"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e BW_TEST_VAR"
  '
  '
    --return "0"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --stdout "${_ansiCmd}grep \"^export BW_TEST_VAR=\" \"$HOME/$_profileFileName\"${_ansiReset}"
    "_hasLineAtProfile -e -v dry BW_TEST_VAR"
  '
  '
    --return "1"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    "_hasLineAtProfile -e -v none BW_TEST_VAR"
  '
  '
    --return "1"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --stderr "${_ansiErr}ERR: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiErr} не обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -v err BW_TEST_VAR"
  '
  '
    --return "1"
    --before "_profileUnless --exportVar -u BW_TEST_VAR"
    --stderr "${_ansiErr}ERR: Строка начинающаяся с ${_ansiCmd}export BW_TEST_VAR=${_ansiErr} не обнаружена в ${_ansiFileSpec}$HOME/$_profileFileName${_ansiReset}"
    "_hasLineAtProfile -e -v all BW_TEST_VAR"
  '
)

# =============================================================================
