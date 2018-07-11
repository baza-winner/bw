
# =============================================================================

_resetBash

# =============================================================================

_runBashTestTests=(
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_runBashTest${_ansiErr} ожидает, что опция ${_ansiCmd}--varValue${_ansiErr} будет задана вместе с опцией ${_ansiCmd}--varName${_ansiReset}"
    --noErrorStack
    "_runBashTest \
      --varValue 1 \
      \"echo me my\" \
      someTest \
    "
  '
      # '$_stqq'echo me'$_stqq' \
  '
    --return "0"
    --noErrorStack
    --stdout "
      before
      immediatly
      after
      ${_ansiOK}OK: someTest${_ansiReset}: ${_ansiCmd}echo me${_ansiReset}
    "
    --stdoutParaWithIndent "0"
    "_runBashTest \
      --stdout me \
      --stdoutTstFileSpec '$_stq'/tmp/bwinstall.bash.runSelfTest.stdout'$_stq' \
      --before '$_stq'echo before'$_stq' \
      --immediatly '$_stq'echo immediatly'$_stq' \
      --after '$_stq'echo after'$_stq' \
      '$_stq'echo me'$_stq' \
      someTest \
    "
  '
)

# =============================================================================
