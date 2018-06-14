
# =============================================================================

_resetBash

# =============================================================================
_bw_removeRequiredFilesAddon=(
  'coreUtils' \
  'funcOptionsSupport2' \
  'coreUtilsWithOptions' \
  'funcParamsSupport2' \
  'coreUtilsWithParams' \
  'coreFileUtils' \
)
_bw_removeRequiredFiles+=( "${_bw_removeRequiredFilesAddon[@]}" )
for _fileSpec in \
  "${_bw_removeRequiredFilesAddon[@]}" \
  'spinnerSupport' \
  'cacheSupport' \
  'psSupport' \
  'inputrcSupport' \
  'bwMain' \
; do
  _fileSpec="$_bwDir/core/$_fileSpec.bash" codeHolder=_codeSource eval "$_evalCode"
done


if [[ -n $_isBwDevelop ]]; then
  for _fileSpec in \
    'core/buildBw' \
    'tests/testsSupport' \
  ; do
    _fileSpec="$_bwDir/$_fileSpec.bash" codeHolder=_codeSource eval "$_evalCode"
  done
fi

# =============================================================================
