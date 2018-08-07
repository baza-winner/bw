
# =============================================================================

_resetBash

# =============================================================================

_getPluralWordTests=(
  '
    --return 0
    --stdout Завершился
    "_getPluralWord 1 Завершился Завершились"
  '
  '
    --return 0
    --stdout Завершились
    "_getPluralWord 2 Завершился Завершились"
  '
  '
    --return 0
    --stdout Завершились
    "_getPluralWord 5 Завершился Завершились"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 0 тест теста тестов"
  '
  '
    --return 0
    --stdout тест
    "_getPluralWord 1 тест теста тестов"
  '
  '
    --return 0
    --stdout теста
    "_getPluralWord 2 тест теста тестов"
  '
  '
    --return 0
    --stdout теста
    "_getPluralWord 4 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 5 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 10 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 11 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 14 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 20 тест теста тестов"
  '
  '
    --return 0
    --stdout тест
    "_getPluralWord 21 тест теста тестов"
  '
  '
    --return 0
    --stdout теста
    "_getPluralWord 22 тест теста тестов"
  '
  '
    --return 0
    --stdout теста
    "_getPluralWord 24 тест теста тестов"
  '
  '
    --return 0
    --stdout тестов
    "_getPluralWord 25 тест теста тестов"
  '
)

_substituteTests=(
  '
    --before "local __substituteTestVarA"
    --varName "__substituteTestVarA"
    --varValue "declare -- __substituteTestVarA=\"2 3\""
    --var2Name "${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}"
    --var2Value "declare -- ${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}=\"0\""
    --var3Name "${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0"
    --var3Value "declare -- ${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0=\"1\""
    --before "__substituteTestVarA=1 ${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}= ${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0="
    "_substitute __substituteTestVarA '$_stqq'2 3'$_stqq'"
  '
  '
    --before "local __substituteTestVarA"
    --varName "__substituteTestVarA"
    --varValue "declare -- __substituteTestVarA=\"\""
    --var2Name "${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}"
    --var2Value "declare -- ${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}=\"0\""
    --var3Name "${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0"
    --var3Value "declare -- ${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0=\"1\""
    --before "__substituteTestVarA=1 ${_substitutePrefix}__substituteTestVarA${_substituteIdxSuffix}= ${_substitutePrefix}__substituteTestVarA${_substituteValueSuffix}0="
    "_substitute __substituteTestVarA"
  '
  # TODO: реанимировать следующий тест
  '
    --before "local __substituteTestVarB"
    --varName "__substituteTestVarB"
    --varValue "declare -a __substituteTestVarB=([0]=\"d\" [1]=\"e f\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    --var2Name "${_substitutePrefix}__substituteTestVarB${_substituteIdxSuffix}"
    --var2Value "declare -- ${_substitutePrefix}__substituteTestVarB${_substituteIdxSuffix}=\"1\""
    --var3Name "${_substitutePrefix}__substituteTestVarB${_substituteValueSuffix}1"
    --var3Value "declare -a ${_substitutePrefix}__substituteTestVarB${_substituteValueSuffix}1=([0]=\"a\" [1]=\"b c\")"
    --before "__substituteTestVarB=( a \"b c\" ) ${_substitutePrefix}__substituteTestVarB${_substituteIdxSuffix}=0 ${_substitutePrefix}__substituteTestVarB${_substituteValueSuffix}1="
    --var3TstPostProcess "perl -pe \"s/='\''\\(/=(/s; s/'\''$//s\""
    "_substitute __substituteTestVarB d '$_stqq'e f'$_stqq'"
  '
    # --var3TstPostProcess "perl -pe \"s/='\''\\(/=(/s; s/'\''$//s\""
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_substitute${_ansiErr} could not resolve type of ${_ansiOutline}__substituteTestVarC${_ansiErr}, first declare it with initial value${_ansiReset}"
    --before "unset __substituteTestVarC"
    "_substitute __substituteTestVarC 1"
  '
)

_lcpTests=(
  '
    --stdout "som"
    --stdoutEchoOptions "-n"
    "_lcp somA somE somEA"
  '
  '
    --stdout ""
    --stdoutEchoOptions "-n"
    "_lcp a b"
  '
)

_shortenFileSpecTests=(
  '
    --stdout "~/bw.bash"
    "_shortenFileSpec '$_stqq'$HOME/bw.bash'$_stqq'"
  '
  '
    --stdout "~/.bw/tests"
    "_shortenFileSpec '$_stqq'$HOME/.bw/tests'$_stqq'"
  '
  '
    --stdout "/var/www/some file"
    "_shortenFileSpec '$_stqq'/var/www/some file'$_stqq'"
  '
)

# =============================================================================
