
# =============================================================================

_resetBash

# =============================================================================

_hasItemTests=(
  '
    --return 0
    "_hasItem a a alpha"
  '
  '
    --return 0
    "_hasItem aa aa b b cc"
  '
  '
    --return 1
    "_hasItem a d b c"
  '
  '
    --return 0
    "_hasItem '$_stqq'a b'$_stqq' d '$_stqq'a b'$_stqq' c"
  '
  '
    --return 1
    "_hasItem '$_stqq'a'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 1
    "_hasItem '$_stqq'?'$_stqq' '$_stqq'? f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 0
    "_hasItem '$_stqq'?'$_stqq' '$_stqq'?'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 1
    "_hasItem ? '$_stqq'? f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 0
    "_hasItem ? ? '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 1
    "_hasItem '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq$_stSlashInQ'* f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 0
    "_hasItem '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 1
    "_hasItem + '$_stqq'+ f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
  '
    --return 0
    "_hasItem + + '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
  '
)

# _hasItem2Tests=(
#   '
#     --return 0
#     "_hasItem2 a a alpha"
#   '
#   '
#     --return 0
#     "_hasItem2 aa aa b b cc"
#   '
#   '
#     --return 1
#     "_hasItem2 a d b c"
#   '
#   '
#     --return 0
#     "_hasItem2 '$_stqq'a b'$_stqq' d '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 1
#     "_hasItem2 '$_stqq'a'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 1
#     "_hasItem2 '$_stqq'?'$_stqq' '$_stqq'? f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 0
#     "_hasItem2 '$_stqq'?'$_stqq' '$_stqq'?'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 1
#     "_hasItem2 ? '$_stqq'? f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 0
#     "_hasItem2 ? ? '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 1
#     "_hasItem2 '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq$_stSlashInQ'* f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 0
#     "_hasItem2 '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq$_stSlashInQ'*'$_stqq' '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 1
#     "_hasItem2 + '$_stqq'+ f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
#   '
#     --return 0
#     "_hasItem2 + + '$_stqq'd f'$_stqq' '$_stqq'a b'$_stqq' c"
#   '
# )

_quotedArgsTests=(
  '
    --stdout "\"?\" \"?\" \"?\""
    "_quotedArgs '$_stq'?'$_stq' '$_stqq'?'$_stqq' ?"
  '
  '
    --stdout "\"\\\\*\" \"\\\\*\""
    "_quotedArgs '$_stq$_stSlashInQ'*'$_stq' '$_stqq$_stSlashInQQ'*'$_stqq'"
  '
  '
    --stdout "\"+\" \"+\" \"+\""
    "_quotedArgs '$_stq'+'$_stq' '$_stqq'+'$_stqq' +"
  '
  '
    --stdout "\$some thing"
    "_quotedArgs '$_stq$_stDollarInQ'some'$_stq' thing"
  '
  '
    --stdout "\"\\\$some\" thing"
    "_quotedArgs --quote:dollarSign '$_stq$_stDollarInQ'some'$_stq' thing"
  '
  '
    --stdout "\\\$some thing"
    "_quotedArgs --strip --quote:dollarSign '$_stq$_stDollarInQ'some'$_stq' thing"
  '
  '
    --stdout "\"\\\$some\" \"thing\""
    "_quotedArgs --quote:all '$_stq$_stDollarInQ'some'$_stq' thing"
  '
  '
    --stdout "\"some${_nl}thing\" good"
    "_quotedArgs '$_stqq'some\${_nl}thing'$_stqq' good"
  '
  '
    --stdout "\"some thing\" good"
    "_quotedArgs '$_stq'some thing'$_stq' good"
  '
  '
    --stdout "\"some \\\"thing\\\"\" good"
    "_quotedArgs '$_stq'some '$_stqq'thing'$_stqq''$_stq' good"
  '
  '
    --stdout "\"some '\''thing'\''\" good"
    "_quotedArgs '$_stqq'some '$_stq'thing'$_stq''$_stqq' good"
  '
  '
    --stdout "\"s\\\\ome\" thing"
    "_quotedArgs '$_stq's'$_stSlashInQ'ome'$_stq' thing"
  '
  '
    --stdout "\"( \\\${defaultValueForOptA[@]} )\""
    "_quotedArgs '$_stq''$_stOpenBraceInQ' '$_stDollarInQ'{defaultValueForOptA[@]} '$_stCloseBraceInQ''$_stq'"
  '
  '
    --stdout "( \\\${defaultValueForOptA[@]} )"
    "_quotedArgs --strip '$_stq''$_stOpenBraceInQ' '$_stDollarInQ'{defaultValueForOptA[@]} '$_stCloseBraceInQ''$_stq'"
  '
  '
    --stdout "some \"\" thing"
    "_quotedArgs some '$_stqq''$_stqq' thing"
  '
)

_getUniqArrayTests=(
  '
    --stdout "\"a b\" \"a c\""
    "_getUniqArray '$_stqq'a b'$_stqq' '$_stqq'a c'$_stqq' '$_stqq'a b'$_stqq'"
  '
)

# _kebabCaseToCamelCaseTests=(
#   '
#     --stdout "someVar"
#     "_kebabCaseToCamelCase some-var"
#   '
# )

# _kebabCaseToUpperCamelCaseTests=(
#   '
#     --stdout "SomeVar"
#     "_kebabCaseToUpperCamelCase some-var"
#   '
# )

# _upperCamelCaseToKebabCaseTests=(
#   '
#     --stdout "some-var"
#     "_upperCamelCaseToKebabCase someVar"
#   '
#   '
#     --stdout "some-var"
#     "_upperCamelCaseToKebabCase SomeVar"
#   '
#   '
#     --stdout "aa-bb-cc"
#     "_upperCamelCaseToKebabCase AaBbCc"
#   '
#   '
#     --stdout "a-b-c"
#     "_upperCamelCaseToKebabCase ABC"
#   '
# )

_restoreTests=(
  '
    --before "local __restoreTestVarA"
    --varName "__restoreTestVarA"
    --varValue "declare -- __restoreTestVarA=\"1\""
    --var2Name "${_substitutePrefix}__restoreTestVarA${_substituteIdxSuffix}"
    --var2Value "declare -- ${_substitutePrefix}__restoreTestVarA${_substituteIdxSuffix}=\"\""
    --before "__restoreTestVarA=2 ${_substitutePrefix}__restoreTestVarA${_substituteIdxSuffix}=0 ${_substitutePrefix}__restoreTestVarA${_substituteValueSuffix}0=1"
    "_restore __restoreTestVarA"
  '
  '
    --before "local __restoreTestVarB"
    --varName "__restoreTestVarB"
    --varValue "declare -a __restoreTestVarB=([0]=\"d\" [1]=\"e f\" [2]=\"g\")"
    --varTstPostProcess '"$_postProcessDeclareArray"'
    --var2Name "${_substitutePrefix}__restoreTestVarB${_substituteIdxSuffix}"
    --var2Value "declare -- ${_substitutePrefix}__restoreTestVarB${_substituteIdxSuffix}=\"0\""
    --before "__restoreTestVarB=(a \"b c\" d) ${_substitutePrefix}__restoreTestVarB${_substituteIdxSuffix}=1 ${_substitutePrefix}__restoreTestVarB${_substituteValueSuffix}1=( d \"e f\" g)"
    "_restore __restoreTestVarB"
  '
  '
    --return "1"
    --stderr "${_ansiErr}ERR: could not resolve type of ${_ansiOutline}__restoreTestVarD${_ansiErr}, first declare it with initial value${_ansiReset}"
    --stderrEchoOptions "-e"
    "_restore __restoreTestVarD"
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

_getTypeOfVarTests=(
  '
    --before "local __getTypeOfVarTestA=0"
    --stdout "scalar"
    "_getTypeOfVar __getTypeOfVarTestA"
  '
  '
    --before "local -a __getTypeOfVarTestB=()"
    --stdout "array"
    "_getTypeOfVar __getTypeOfVarTestB"
  '
  '
    --before "unset __getTypeOfVarTestD"
    --stdout "none"
    "_getTypeOfVar __getTypeOfVarTestD"
  '
)

_upperFirstTests=(
  '
    --stdout Some
    "_upperFirst some"
  '
  '
    --stdout Слово
    "_upperFirst слово"
  '
)

_lowerFirstTests=(
  '
    --stdout some
    "_lowerFirst Some"
  '
  '
    --stdout слово
    "_lowerFirst Слово"
  '
)