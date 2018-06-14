
# =============================================================================

_resetBash

# =============================================================================

bwTests=(
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}bw${_ansiErr} вместо ${_ansiPrimaryLiteral}unexpected${_ansiErr} ожидает одну из следующих команд: ${_ansiSecondaryLiteral}bash-tests bt install project remove rm run set update version${_ansiReset}"
    --before "_substitute noStack true"
    --after "_restore noStack"
    "bw unexpected"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}bw${_ansiErr} в качестве первого аргумента ожидает одну из следующих команд: ${_ansiSecondaryLiteral}bash-tests bt install project remove rm run set update version${_ansiReset}"
    --before "_substitute noStack true"
    --after "_restore noStack"
    "bw"
  '
)

bw_removeTests=(
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}bw remove${_ansiReset} [${_ansiOutline}Опции${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset} удаляет bw.bash и все связанное с ним
      ${_ansiOutline}Опции${_ansiReset}
          ${_ansiCmd}--yes${_ansiReset} или ${_ansiCmd}-y${_ansiReset}
              подтверждает удаление
          ${_ansiCmd}--completely${_ansiReset} или ${_ansiCmd}-c${_ansiReset}
              удаляет не только все связанное с bw.bash, но и сам bw.bash
          ${_ansiCmd}--verbosity${_ansiReset} ${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-v${_ansiReset} ${_ansiOutline}значение${_ansiReset}
              Варианты ${_ansiOutline}значения${_ansiReset}:
                ${_ansiPrimaryLiteral}dry${_ansiReset} режим \"dry run\"
                ${_ansiPrimaryLiteral}none${_ansiReset} молчаливый режим
                ${_ansiPrimaryLiteral}err${_ansiReset} только вывод ошибок
                ${_ansiPrimaryLiteral}ok${_ansiReset} только вывод OK
                ${_ansiPrimaryLiteral}allBrief${_ansiReset} вывод результатов в краткой форме
                ${_ansiPrimaryLiteral}all${_ansiReset} полный вывод
              ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}allBrief${_ansiReset}
          ${_ansiCmd}--silent${_ansiReset} ${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-s${_ansiReset} ${_ansiOutline}значение${_ansiReset}
              подавляет сообщение об успешном удалении
              Варианты ${_ansiOutline}значения${_ansiReset}:
                ${_ansiPrimaryLiteral}yes${_ansiReset} не подавлять вывод вспомогательных команд
                ${_ansiPrimaryLiteral}ok${_ansiReset} подавлять вывод вспомогательных команд в случае отсутствия ошибок
                ${_ansiPrimaryLiteral}err${_ansiReset} подавлять вывод вспомогательных команд в случае наличия ошибок
                ${_ansiPrimaryLiteral}no${_ansiReset} подавлять вывод вспомогательных команд
              ${_ansiOutline}Значение${_ansiReset} по умолчанию: ${_ansiPrimaryLiteral}no${_ansiReset}
          ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
              Выводит справку
    "
    --stdoutParaWithIndent 0
    --stdoutParaWithIndentBase 2
    --before "_substitute noStack true"
    --after "_restore noStack"
    "bw remove -?"
  '
)

bw_updateTests=(
  '
    --return "3"
    --stdout "
      ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}bw update${_ansiReset} [${_ansiOutline}Опции${_ansiReset}]
      ${_ansiHeader}Описание:${_ansiReset} $bw_update_description
      ${_ansiOutline}Опции${_ansiReset}
          ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset}
              Выводит справку
    "
    --stdoutParaWithIndent 0
    --stdoutParaWithIndentBase 2
    --before "_substitute noStack true"
    --after "_restore noStack"
    "bw update -?"
  '
)

_bw_bashTestTestFunc() {
  echo "    some
    multiline
    output"
  return 3
}

bw_bashTestTests=(
  '
    --return "3"
    --stdout "
      some
      multiline
      output
    "
    --stdoutParaWithIndent "1"
    "_bw_bashTestTestFunc"
  '
  '
    --stdout "
      ${_ansiOK}OK: _test 0${_ansiReset}: ${_ansiCmd}_ok test success${_ansiReset}
      ${_ansiOK}OK: _test2 1${_ansiReset}: ${_ansiCmd}_ok test2 success${_ansiReset}
      ${_ansiOK}OK: Все тесты (${_ansiPrimaryLiteral}2${_ansiOK}) пройдены успешно${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    --stdoutTstFileSpec "/tmp/'$(basename "${BASH_SOURCE[0]}")'.stdout"
    --stderrTstFileSpec "/tmp/'$(basename "${BASH_SOURCE[0]}")'.stderr"
    "__allowTestsForTest=true bw bt --noTiming _test 0 _test2 1"
  '
  '
    --return "1"
    --before "_substitute noStack true"
    --after "_restore noStack"
    --stdout "
      ${_ansiOK}OK: _test 0${_ansiReset}: ${_ansiCmd}_ok test success${_ansiReset}
      ${_ansiOK}OK: _test2 1${_ansiReset}: ${_ansiCmd}_ok test2 success${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    --stderr "
      ${_ansiErr}ERR: _test 1${_ansiReset}: ${_ansiCmd}_err test fail${_ansiReset}
      ${_ansiOutline}stderr${_ansiReset} (${_ansiCmd}diff /tmp/testsSupport.bash.stderr /tmp/testsSupport.bash.stderr.eta${_ansiReset}):
      1d0
      < ${_ansiErr}ERR: test fail${_ansiReset}
      ${_ansiErr}ERR: _test 2${_ansiReset}: ${_ansiCmd}_err another test fail${_ansiReset}
      ${_ansiOutline}stderr${_ansiReset} (${_ansiCmd}diff /tmp/testsSupport.bash.stderr /tmp/testsSupport.bash.stderr.eta${_ansiReset}):
      1d0
      < ${_ansiErr}ERR: another test fail${_ansiReset}
      ${_ansiErr}ERR: _test2 2${_ansiReset}: ${_ansiCmd}_err another test2 fail${_ansiReset}
      ${_ansiOutline}stderr${_ansiReset} (${_ansiCmd}diff /tmp/testsSupport.bash.stderr /tmp/testsSupport.bash.stderr.eta${_ansiReset}):
      1d0
      < ${_ansiErr}ERR: another test2 fail${_ansiReset}
      ${_ansiErr}ERR: _test2 0${_ansiReset}: ${_ansiCmd}_err test2 fail${_ansiReset}
      ${_ansiOutline}return${_ansiErr} is expected to be ${_ansiOK}0${_ansiErr}, but got ${_ansiWarn}1${_ansiReset}
      ${_ansiErr}ERR: Завершились ошибкой 4 тестов из 6:${_ansiReset}
      ${_ansiErr}ERR:#SPACE#_test 1..2${_ansiReset}
      ${_ansiErr}ERR:#SPACE#_test2 0 2${_ansiReset}
    "
    --stderrEtaPreProcess "perl -pe \"s/#SPACE#/\ \ \ /g\""
    --stderrParaWithIndent "0"
    --stdoutTstFileSpec "/tmp/'$(basename "${BASH_SOURCE[0]}")'.stdout"
    --stderrTstFileSpec "/tmp/'$(basename "${BASH_SOURCE[0]}")'.stderr"
    "__allowTestsForTest=true bw bt --noTiming _test -3.. _test2 2 ..1"
  '
)
_test() { true; }
_testTestsCondition='[[ -n $__allowTestsForTest ]]'
_testTests=(
  '
    --stdout "${_ansiOK}OK: test success${_ansiReset}"
    "_ok test success"
  '
  '
    --return "1"
    "_err test fail"
  '
  '
    --return "1"
    "_err another test fail"
  '
)
_test2() { true; }
_test2TestsCondition='[[ -n $__allowTestsForTest ]]'
_test2Tests=(
  '
  --stderr "${_ansiErr}ERR: test2 fail${_ansiReset}"
    "_err test2 fail"
  '
  '
    --stdout "${_ansiOK}OK: test2 success${_ansiReset}"
    "_ok test2 success"
  '
  '
    --return "1"
    "_err another test2 fail"
  '
)

# bw_installTests=(
#   '
#     --return "1"
#     --stderr "${_ansiErr}ERR: ${_ansiCmd}bw install${_ansiErr} вместо ${_ansiPrimaryLiteral}unexpected${_ansiErr} ожидает одну из следующих команд: ${_ansiSecondaryLiteral}docker${_ansiReset}"
#     --before "_substitute noStack true"
#     --after "_restore noStack"
#     "bw install unexpected"
#   '
#   '
#     --return "1"
#     --stderr "${_ansiErr}ERR: ${_ansiCmd}bw install${_ansiErr} в качестве первого аргумента ожидает одну из следующих команд: ${_ansiSecondaryLiteral}docker${_ansiReset}"
#     --before "_substitute noStack true"
#     --after "_restore noStack"
#     "bw install"
#   '
# )

# =============================================================================
