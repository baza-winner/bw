
# =============================================================================

_resetBash

# =============================================================================

_evalCode() {
  return $(_err "use variable ${_ansiOutline}_evalCode${_ansiErr} instead")
}
_evalCodeTestFunc() {
  codeHolder="$1" eval "$_evalCode"
  eval "$2"
}
_testCodeToEval='
  while true; do
    return $(_err "some error")
  done
'
_testCodeToEval2='
  if true; then
    local someVar=someValue
  fi
'
_evalCodeTests=(
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiOutline}_evalCode${_ansiErr} expects ${_ansiOutline}codeHolder${_ansiErr} to be specified${_ansiReset}"
    "_evalCodeTestFunc"
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiOutline}_evalCode${_ansiErr} expects ${_ansiOutline}_nonExistent${_ansiErr} to be defined${_ansiReset}"
    "_evalCodeTestFunc _nonExistent"
  '
  '
    --noErrorStack
    --return "1"
    --stderr "${_ansiErr}ERR: some error${_ansiReset}"
    "_evalCodeTestFunc _testCodeToEval"
  '
  '
    --noErrorStack
    --stderr "${_ansiOutline}someVar${_ansiReset}: ${_ansiPrimaryLiteral}someValue${_ansiReset}"
    "_evalCodeTestFunc _testCodeToEval2 '$_stq'_debugVar --clean someVar'$_stq'"
  '
)

# =============================================================================


# _initialInstall() {
#   true
# }
# _initialInstallTestsCondition='[[ -n $_isBwDevelop || -n $BW_TEST_ALL ]]'
# _initialInstallTests=(
#   '
#     --before "_pushd \"$_bwDir\""
#     --before "_dockerCompose up -d"
#     --after "_popd"
#   '
# )
# _initialInstallTestFunc() {
#   _inDir -t "$_bwTmpDir" _initialInstallTestFuncHelper
# }
# _initialInstallTestFuncHelper() {
#   curl -O localhost:8082/bw.bash \
#   && . bw.bash -u:localhost:8082 \
#   && _exist -d .bw \
#   && bw rm -y \
#   && _exist -d -n .bw \
#   && _exist -n \
#   true
# }

# =============================================================================
