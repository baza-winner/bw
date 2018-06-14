
# =============================================================================

_resetBash

# =============================================================================

_spinnerTestFuncParams=( 'testNum:0..2' )
_spinnerTestFunc() { eval "$_funcParams"
  _fileSpec="$_bwDir/core/spinnerSupport.bash" codeHolder=_codeSource eval "$_evalCode"
  if [[ $testNum -eq 0 ]]; then
    _spinner '-t:Sleep обработан за' Sleep sleep 1 & 2>/dev/null
    local spinnerPid=$!
    wait $spinnerPid 2>/dev/null
  elif [[ $testNum -eq 1 ]]; then
    { _spinner Sleep sleep 2; } & 2>/dev/null
    local spinnerPid=$!
    sleep 1
    kill -SIGTERM $spinnerPid
    wait $spinnerPid 2>/dev/null
  fi
}
_spinnerTestsCondition='[[ -n $BW_TEST_ALL ]]'
_spinnerTests=(
  '
    --stdout "Sleep . . .  /-\|/-\|             Sleep обработан за 1 секунду"
    --catOptions "-e"
    "_spinnerTestFunc 0"
  '
  # WHY NO TEST FOR SIGINT: https://unix.stackexchange.com/questions/372541/why-doesnt-sigint-work-on-a-background-process-in-a-script
  '
    --return "143"
    --stdout "Sleep . . .  /-\|/-\|"
    --stdoutEchoOptions "-n"
    --catOptions "-e"
    "_spinnerTestFunc 1"
  '
)