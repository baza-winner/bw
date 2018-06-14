
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
    --return "1"
    --stderr "${_ansiErr}ERR: Переменная ${_ansiOutline}__someVar${_ansiErr} не определена${_ansiReset}"
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
  '
    --before "local __someVar="
    --before "! local __someVarValueHolder=\"some Var${_nl}Value\""
    --stdout "__someVar='\''some Var${_nl}Value'\''"
    "_saveToCache -d __someVar '$_stqq''$_stDollarInQ'__someVarValueHolder'$_stqq'"
  '
  '
    --before "local __someVar="
    --before "! local __someVarValueHolder=\"some '\''Var'\''${_nl}Value\""
    --stdout "__someVar=\"some '\''Var'\''${_nl}Value\""
    "_saveToCache -d __someVar '$_stqq''$_stDollarInQ'__someVarValueHolder'$_stqq'"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: ${_ansiCmd}_saveToCache${_ansiErr} ожидает не более одного значения для ${_ansiOutline}__someVar${_ansiErr}, но получены 2: ${_ansiPrimaryLiteral}\"some Var Value\" secondItem${_ansiReset}"
    "_saveToCache -d __someVar '$_stqq'some Var Value'$_stqq' secondItem"
  '
  '
    --before "local __someArr=()"
    --before "! local __firstItemValueHolder=\"first${_nl}item\""
    --stdout "__someArr=([0]=\$'\''first\\nitem'\'' [1]=\"secondItem\")"
    --stdoutTstPostProcess '"$_postProcessCachedArray"'
    "_saveToCache -d -a __someArr '$_stqq''$_stDollarInQ'__firstItemValueHolder'$_stqq' secondItem"
  '
  '
    --before "local __someArr="
    --before "! local __firstItemValueHolder=\"first${_nl}item\""
    --stdout "__someArr=([0]=\$'\''first\\nitem'\'')"
    --stdoutTstPostProcess '"$_postProcessCachedArray"'
    "_saveToCache -d -a __someArr '$_stqq''$_stDollarInQ'__firstItemValueHolder'$_stqq'"
  '
)

# =============================================================================
