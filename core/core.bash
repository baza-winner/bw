
# =============================================================================

_resetBash

# =============================================================================
_bwCoreFileNames=(
  'coreUtils' 
  'funcOptionsSupport2' 
  'coreUtilsWithOptions' 
  'funcParamsSupport2' 
  'coreUtilsWithParams' 
  'coreFileUtils' 
  'cacheSupport'
  'spinnerSupport'
)
for _fileSpec in "${_bwCoreFileNames[@]}"; do
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
