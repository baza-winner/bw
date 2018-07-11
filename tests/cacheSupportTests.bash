
# =============================================================================

_resetBash

# =============================================================================

_postProcessCachedArray='"perl -e \"
    undef \\\$/; 
    \\\$_=<STDIN>; 
    sub replaceTo {
      my \\\$q=shift;
      my \\\$val=shift;
      if (\\\$val =~ s/\n/\\\\\\\\n/g) {
        \\\"\\\\$'\''\\\" . \\\$val . \\\"'\''\\\"
      } else {
        \\\$q . \\\$val . \\\$q
      }
    }
    s/(\\\")([^\\\"]+)\\\"/replaceTo(\\\$1,\\\$2)/sge; 
    print
  \""'
_saveToCacheTests=(
  '
    --before "unset __someVar"
    --return 1
    --noErrorStack
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_saveToCache${_ansiErr} ожидает, что переменная ${_ansiOutline}__someVar${_ansiErr} будет определена${_ansiReset}"
    "_saveToCache __someVar"
  '
  '
    --before "! local __someVar=\"some Var${_nl}Value\""
    --stdout "__someVar='\''some Var${_nl}Value'\''"
    "_saveToCache -d __someVar"
  '
  '
    --before "! local __someVar=\"some '\''Var'\''${_nl}Value\""
    --stdout "__someVar=\"some '\''Var'\''${_nl}Value\""
    "_saveToCache -d __someVar"
  '
  '
    --before "! local -a __someArr=( \"first${_nl}item\" \"second item\" )"
    --stdout "__someArr=([0]=\$'\''first\\nitem'\'' [1]=\"second item\")"
    --stdoutTstPostProcess '"$_postProcessCachedArray"'
    "_saveToCache -d __someArr"
  '
)

# =============================================================================
