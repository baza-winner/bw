
# =============================================================================

_resetBash

# =============================================================================

# _autoHelpParams=(
#   '--bool-opt/a'
#   '--scalar-opt0='
#   '--scalar-opt:..'
#   '--scalar-opta/b:..='
#   '--scalar-opt2:0..'
#   '--scalar-opt2a/c:0..=1'
#   '--scalar-opt3:1..=2'
#   '--scalar-opt3a/d:1..'
#   '--scalar-opt4:..-1'
#   '--scalar-opt4a/e:..-1=-2'
#   '--scalar-opt5:..0'
#   '--scalar-opt5a/f:..0=-1'
#   '--scalar-opt6:2..4'
#   '--scalar-opt6a/g:2..4=3'
#   '--scalar-opt7/z:(x y z)=y'
#   '--scalar-opt8/i:( $(echo x y z) )=$(echo x)'
#   '@1..4--list-opt/l:..=(2 3)'
#   '@1..4--list-opt2/L:..:unique=( $(echo 2 3) )'
#   'argA:..=0'
#   'argB:(x y z)'
#   'argC:(a b c)=b'
#   'argD:(a b c):?=b'
#   '@args:?'
# )
# _autoHelpDescription='autoHelp test function'
# _autoHelpDescriptionOfScalarOptA='optA option single line description'
# _autoHelpDescriptionOfScalarOpt7='
#   scalar-opt7
#   multi line description
# '
# _autoHelpDescriptionOfScalarOpt8_x='x value single line description'
# _autoHelpDescriptionOfScalarOpt8_y='
#   y value
#   multi line description
# '
# _autoHelp() { eval "$_funcParams"
#   _debugVar ${__varNames[@]}
# }
# _autoHelpTests=(
#   '
#     --return "2"
#     "--stdout:
#       ${_ansiHeader}Использование:${_ansiReset} ${_ansiCmd}_autoHelp${_ansiReset} [${_ansiOutline}Опции${_ansiReset}] ${_ansiOutline}argA${_ansiReset} ${_ansiOutline}argB${_ansiReset} [${_ansiOutline}argC${_ansiReset} [${_ansiOutline}argD${_ansiReset} [${_ansiOutline}args${_ansiReset}...]]]
#       ${_ansiHeader}Описание:${_ansiReset} autoHelp test function
#       ${_ansiOutline}Опции${_ansiReset}
#         ${_ansiCmd}--help${_ansiReset} или ${_ansiCmd}-?${_ansiReset} или ${_ansiCmd}-h${_ansiReset} Выводит справку
#         ${_ansiCmd}--bool-opt${_ansiReset} или ${_ansiCmd}-a${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfBoolOpt${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt0=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} должно быть непустым
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt0${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opta=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-b=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpta${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt2=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - неотрицательное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt2${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt2a=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-c=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - неотрицательное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt2a${_ansiErr} действия опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}1${_ansiReset}
#         ${_ansiCmd}--scalar-opt3=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - положительное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt3${_ansiErr} действия опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}2${_ansiReset}
#         ${_ansiCmd}--scalar-opt3a=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-d=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - положительное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt3a${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt4=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - отрицательное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt4${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt4a=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-e=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - отрицательное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt4a${_ansiErr} действия опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}-2${_ansiReset}
#         ${_ansiCmd}--scalar-opt5=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - неположительное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt5${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt5a=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-f=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - неположительное целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt5a${_ansiErr} действия опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}-1${_ansiReset}
#         ${_ansiCmd}--scalar-opt6=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число из диапазона ${_ansiSecondaryLiteral}2..4${_ansiReset}
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt6${_ansiErr} действия опции${_ansiReset}
#         ${_ansiCmd}--scalar-opt6a=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-g=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число из диапазона ${_ansiSecondaryLiteral}2..4${_ansiReset}
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt6a${_ansiErr} действия опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}3${_ansiReset}
#         ${_ansiCmd}--scalar-opt7=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-z=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}x y z${_ansiReset}
#           scalar-opt7
#           multi line description
#           ${_ansiOutline}значение${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}y${_ansiReset}
#         ${_ansiCmd}--scalar-opt8=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-i=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}x y z${_ansiReset}
#           ${_ansiPrimaryLiteral}x${_ansiReset} x value single line description
#           ${_ansiPrimaryLiteral}y${_ansiReset}
#             y value
#             multi line description
#           ${_ansiPrimaryLiteral}z${_ansiReset} ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfScalarOpt8_z${_ansiErr} значения опции${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} по умолчанию - значение выражения ${_ansiOutline}$ (echo x)${_ansiReset}
#         ${_ansiCmd}--list-opt=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-l=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfListOpt${_ansiErr} действия опции${_ansiReset}
#           Опция предназначена для того, чтобы сформировать непустой список (не более 4 элементов) ${_ansiOutline}значений${_ansiReset} путем многократного использования этой опции
#           Значение ${_ansiOutline}списка${_ansiReset} по умолчанию: ${_ansiSecondaryLiteral}( 2 3 )${_ansiReset}
#         ${_ansiCmd}--list-opt2=${_ansiOutline}значение${_ansiReset} или ${_ansiCmd}-L=${_ansiOutline}значение${_ansiReset}
#           ${_ansiOutline}значение${_ansiReset} - целое число
#           ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfListOpt2${_ansiErr} действия опции${_ansiReset}
#           Опция предназначена для того, чтобы сформировать непустой список (не более 4 элементов) ${_ansiUnderline}уникальных${_ansiReset} ${_ansiOutline}значений${_ansiReset} путем многократного использования этой опции
#           Значение ${_ansiOutline}списка${_ansiReset} по умолчанию - значение выражения ${_ansiOutline}( $ (echo 2 3) )${_ansiReset}
#       ${_ansiOutline}argA${_ansiReset}
#         ${_ansiOutline}argA${_ansiReset} - целое число
#         ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfArgA${_ansiErr} аргумента${_ansiReset}
#         ${_ansiOutline}argA${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}0${_ansiReset}
#       ${_ansiOutline}argB${_ansiReset}
#         ${_ansiOutline}argB${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}x y z${_ansiReset}
#         ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfArgB${_ansiErr} аргумента${_ansiReset}
#       ${_ansiOutline}argC${_ansiReset}
#         ${_ansiOutline}argC${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}a b c${_ansiReset}
#         ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfArgC${_ansiErr} аргумента${_ansiReset}
#         ${_ansiOutline}argC${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}b${_ansiReset}
#       ${_ansiOutline}argD${_ansiReset}
#         ${_ansiOutline}argD${_ansiReset} - один из вариантов: ${_ansiSecondaryLiteral}a b c${_ansiReset}
#         ${_ansiOutline}argD${_ansiReset} может быть пустым
#         ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfArgD${_ansiErr} аргумента${_ansiReset}
#         ${_ansiOutline}argD${_ansiReset} по умолчанию ${_ansiPrimaryLiteral}b${_ansiReset}
#       ${_ansiOutline}args${_ansiReset}... возможно пустой список ${_ansiOutline}значений${_ansiReset}
#         ${_ansiOutline}значение${_ansiReset} может быть пустым
#         ${_ansiErr}Нет описания ${_ansiOutline}_autoHelpDescriptionOfArgs${_ansiErr} аргумента${_ansiReset}
#     "
#     --stdoutParaWithIndent "0"
#     --stdoutEtaPreProcess "perl -pe s/\\\$\s/\\\$/g"
#     "_autoHelp -?"
#   '
# )

# _getCodeToGetValueDescriptionTestFuncParams=( '--isDiapVarType' '--diapMin' '--diapMax' 'indentLevel' 'valueName' )
# _getCodeToGetValueDescriptionTestFunc() { eval "$_funcParams"
#   codeHolder=_codeToPrepareValueDescription eval "$_evalCode"
#   echo "$valueDescription"
# }
# _getCodeToGetValueDescriptionTests=(
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - неотрицательное целое число"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:0 --diapMax:Infinite 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 2)_indent 2; echo \"\${_ansiOutline}arg\${_ansiReset} - положительное целое число"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:1 --diapMax:Infinite 2 arg"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - целое число не менее \${_ansiPrimaryLiteral}2\${_ansiReset}"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:2 --diapMax:Infinite 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - целое число"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:Infinite --diapMax:Infinite 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - неположительное целое число"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:Infinite --diapMax:0 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - отрицательное целое число"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:Infinite --diapMax:-1 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - целое число не более \${_ansiPrimaryLiteral}-2\${_ansiReset}"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:Infinite --diapMax:-2 1 значение"
#   '
#   '
#     --stdout "\";$_nl$(_indent --base:2 1)_indent 1; echo \"\${_ansiOutline}значение\${_ansiReset} - целое число из диапазона \${_ansiSecondaryLiteral}-2..-2\${_ansiReset}"
#     "_getCodeToGetValueDescriptionTestFunc --isDiapVarType --diapMin:-2 --diapMax:-2 1 значение"
#   '
# )

# =============================================================================
