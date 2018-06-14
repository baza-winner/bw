
# =============================================================================

_resetBash

# =============================================================================

# _getCodeToGetListDescriptionTestFuncParams=( 'minListCount' 'maxListCount' 'uniqueVarType:?')
# _getCodeToGetListDescriptionTestFunc() { eval "$_funcParams"
#   codeHolder=_codeToPrepareListDescription eval "$_evalCode"
#   echo $listDescription
# }
# _getCodeToGetListDescription() {
#   echo $_codeToPrepareListDescription
# }
# _getCodeToGetListDescriptionTests=(
#   '
#     --stdout "возможно пустой список \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 0 Infinite"
#   '
#   '
#     --stdout "возможно пустой список \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 0 Infinite unique"
#   '
#   '
#     --stdout "непустой список \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 1 Infinite"
#   '
#   '
#     --stdout "непустой список \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 1 Infinite unique"
#   '
#   '
#     --stdout "список (не менее 2 элементов) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 2 Infinite"
#   '
#   '
#     --stdout "список (не менее 2 элементов) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 2 Infinite unique"
#   '
#   '
#     --stdout "список (не менее 21 элемента) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 21 Infinite"
#   '
#   '
#     --stdout "список (не менее 21 элемента) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 21 Infinite unique"
#   '
#   '
#     --stdout "возможно пустой список (не более 2 элементов) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 0 2"
#   '
#   '
#     --stdout "возможно пустой список (не более 2 элементов) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 0 2 unique"
#   '
#   '
#     --stdout "непустой список (не более 2 элементов) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 1 2"
#   '
#   '
#     --stdout "непустой список (не более 2 элементов) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 1 2 unique"
#   '
#   '
#     --stdout "список (не менее 2 и не более 21 элемента) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 2 21"
#   '
#   '
#     --stdout "список (не менее 2 и не более 21 элемента) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 2 21 unique"
#   '
#   '
#     --stdout "список (не менее 21 и не более 30 элементов) \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 21 30"
#   '
#   '
#     --stdout "список (не менее 21 и не более 30 элементов) \${_ansiUnderline}уникальных\${_ansiReset} \${_ansiOutline}значений\${_ansiReset}"
#     "_getCodeToGetListDescriptionTestFunc 21 30 unique"
#   '
# )

# =============================================================================
