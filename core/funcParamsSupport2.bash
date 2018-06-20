# =============================================================================

_resetBash

# =============================================================================

_funcParams2='local codeFileSpec && _prepareCodeToParseFuncParams2 && local __consumedParams=0 && . "$codeFileSpec" && _parseFuncParams2 "$@" && { for ((; __consumedParams>0; __consumedParams--)); do shift; done; true; } || return $?'

_help_paramDef="!--help/?h"
_help_description='Выводит справку'

_prepareCodeToParseFuncParams2BoolOptions=(
  'treatUnknownOptionAsArg'
  'isCommandWrapper'
  'canBeMixedOptionsAndArgs'
  'canBeMoreParams'
  'needPregenHelp'
)
_prepareCodeToParseFuncParams2ScalarOptions=(
  'additionalSuffix'
)
_prepareCodeToParseFuncParams2ListOptions=(
  'additionalDependencies'
)
_isCalcEnumRegExp='^\([[:space:]]*\$\('
_isValidVarNameRegExp='^[[:alpha:]_][[:alnum:]_]*$'
_prepareCodeToParseFuncParams2() {
_profileBegin
  local funcName="${FUNCNAME[1]}"  __thisFuncCommand=
  local __thisOnlyPrepareCode="$__onlyPrepareCode" __onlyPrepareCode=
  local __hasWrapper=$hasWrapper hasWrapper=
  kind=Params eval "$_codeToCheckParams"

  dstVarName=paramsOpt srcVarName=${funcName}ParamsOpt codeHolder=_codeToInitLocalCopyOfArrWithCheck eval "$_evalCode"
  set -- "${paramsOpt[@]}"
  eval "$_funcOptions2Helper"
  eval "$_codeToCheckNoArgsInOpt"

  dstVarName= codeType=funcParams fileSpec= originalCodeDeep= eval "$_codeToPrepareCodeFileSpec"
  local verbose=
  # verbose=true

  [[ -z $__thisOnlyPrepareCode || -n $needPregenHelp || $funcName =~ Complete$ || $funcName =~ ^_ ]] || needPregenHelp=true
  if [[ -n $isCommandWrapper ]]; then
    canBeMoreParams=true
    needPregenHelp=
  elif [[ $funcName =~ [[:alnum:]]_[_[:alnum:]] ]]; then
    needPregenHelp=
  fi

  local __err=
  while true; do
    if \
      [[ ! -f $codeFileSpec ]] || \
      ( [[ -n $__thisOnlyPrepareCode || -n $_isBwDevelop || -n $_isBwDevelopInherited ]] \
        && ! _everyFileNotNewerThan "$codeFileSpec" "${BASH_SOURCE[0]}" "${BASH_SOURCE[1]}" "${additionalDependencies[@]}" \
      ) \
    ; then
      [[ -z $verbose ]] || _warn "Создаем ${_ansiFileSpec}$codeFileSpec $completionCodeFileSpec"
      local paramsHolder="${funcName}Params";
      ! _funcExists $paramsHolder || $paramsHolder || return $?;
      dstVarName=params srcVarName=$paramsHolder eval "$_codeToInitLocalCopyOfArray"
      if [[ $funcName =~ Complete$ ]]; then
        if [[ ${#params[@]} -gt 0 ]]; then
          local __infix; _funcExists $paramsHolder && __infix="функция ${_ansiCmd}" || __infix="переменная ${_ansiOutline}"
          __err="не ожидает, что будет определена ${__infix}${paramsHolder}${_ansiErr}, т.к. все ${_ansiOutline}*${_ansiCmd}Complete${_ansiErr}-функции имеют предопределенный ${_ansiOutline}_completeParams${_ansiErr}: ${_ansiSecondaryLiteral}( $(_quotedArgs ${_completeParams[@]}) )"
          break
        fi
        local -a params=( "${_completeParams[@]}" );
      fi

      local -a __argVarNames=() __optVarNames=() OPTVarNames=() importantOptVarNames=()
      local allOptionsShortcuts=
      local listOfArgsParam=
      local __argIdx=0
      local code=

      local param prevParam; for param in "$_help_paramDef" "${params[@]}"; do
        [[ -n $param ]] || continue
        local isImportant= varType= varKind= __varName= __varValue= valType=
        local paramPartToParse="$param"

        if [[ $paramPartToParse =~ ^!@?-- ]]; then
          paramPartToParse=${paramPartToParse:1}
          isImportant=true
        fi

        if [[ $paramPartToParse =~ ^[[:space:]]*@$_listBordersRegExp(.*)$ ]]; then
          varType=list
          local minListCount="${BASH_REMATCH[1]}"
          local upTo="${BASH_REMATCH[2]}"
          local maxListCount="${BASH_REMATCH[3]}"
          local paramPartToParse="${BASH_REMATCH[4]}"
          if [[ -z $maxListCount ]]; then
            [[ -n $upTo || -z $minListCount ]] && maxListCount=Infinite || maxListCount="$minListCount"
          fi
          [[ -n $minListCount ]] || minListCount=0
          if [[ $maxListCount != Infinite && $minListCount -gt $maxListCount ]]; then
            __err="левая граница ${_ansiPrimaryLiteral}$minListCount${_ansiErr} не должна превосходить правую ${_ansiPrimaryLiteral}$maxListCount"
          elif [[ $maxListCount != Infinite && $maxListCount -lt 1 ]]; then
            __err="правая граница ${_ansiPrimaryLiteral}$maxListCount${_ansiErr} не должна быть меньше ${_ansiPrimaryLiteral}1"
          fi
          [[ -z $__err ]] || __err="ожидает, что в параметре ${_ansiCmd}$param${_ansiErr} $__err"
        fi
        [[ -z $__err ]] || break

        if [[ $paramPartToParse =~ ^-- ]]; then
          varKind=opt
          [[ $__argIdx -eq 0 || -n $canBeMixedOptionsAndArgs ]] \
            || __err="ожидает, что опции будут следовать строго перед аргументами, т.к. ${_ansiCmd}${FUNCNAME[0]}${_ansiErr} вызвана без ${_ansiCmd}--canBeMixedOptionsAndArgs${_ansiErr}, но обнаружено ${_ansiOutline}определение опции ${_ansiCmd}$param${_ansiErr} после ${_ansiOutline}определения аргумента ${_ansiCmd}$prevParam"
          paramPartToParse="${paramPartToParse:2}"
        else
          varKind=arg
          if [[ -n $isCommandWrapper ]]; then
            __err="не ожидает определения аргументов в режиме ${_ansiCmd}--isCommandWrapper${_ansiErr}, но получено ${_ansiCmd}$param"
          elif [[ -n $listOfArgsParam ]]; then
            __err="ожидает, что определение списочного аргумента ${_ansiCmd}$listOfArgsParam${_ansiErr} будет последним в списке определений аргументов, но после него следует ещё ${_ansiCmd}$param"
          else
            [[ $varType == list ]] \
              && listOfArgsParam="$param" \
              || varType=scalar
          fi
        fi
        [[ -z $__err ]] || break

        if [[ $paramPartToParse =~ = ]]; then
          [[ -n $varType ]] || varType=scalar
          __varValue="${paramPartToParse#*=}"
          paramPartToParse="${paramPartToParse%%=*}"
        fi
        __varName="$paramPartToParse"

        local canBeEmptyVarType= nonEmptyVarType= sortedVarType= uniqueVarType=
        if [[ $paramPartToParse =~ : ]]; then
          [[ -n $varType ]] || varType=scalar
          local varTypeDef=":${paramPartToParse#*:}:"
          local __varName="${paramPartToParse%%:*}"
          while true; do
            if [[ $varTypeDef =~ :[[:space:]]*($_varTypesInRegExp2)[[:space:]]*:$ ]]; then
              varTypeDef=${varTypeDef:0:$(( ${#varTypeDef} - ${#BASH_REMATCH[0]} + 1 ))}
            elif [[ $varTypeDef =~ ^:[[:space:]]*($_varTypesInRegExp2)[[:space:]]*: ]]; then
              varTypeDef=${varTypeDef:$(( ${#BASH_REMATCH[0]} - 1 )):$(( ${#varTypeDef} - ${#BASH_REMATCH[0]} + 1 ))}
            else
              break
            fi
            local varTypeValue="${BASH_REMATCH[1]}"
            if [[ $varTypeValue == "?" ]]; then
              varTypeValue=canBeEmpty
            elif [[ ( $varTypeValue == sorted || $varTypeValue == unique ) && $varType == scalar ]]; then
              __err="не ожидает тип ${_ansiPrimaryLiteral}$varTypeValue${_ansiErr} в определении ${_ansiOutline}скалярного${_ansiErr} параметра ${_ansiCmd}$(_quotedArgs --strip "$param")"
              break
            fi
            local ${varTypeValue}VarType=true
          done; [[ -z $__err ]] || break

          [[ ${#varTypeDef} -le 1 ]] \
            && varTypeDef= \
            || varTypeDef=${varTypeDef:1:$(( ${#varTypeDef} - 2 ))}
          if [[ $varTypeDef =~ ^[[:space:]]*\( ]]; then
            valType=enum
            local enumDef
            if ! [[ $varTypeDef =~ \)[[:space:]]*$ ]]; then
              __err="ожидает, что определение перечислимого типа ${_ansiPrimaryLiteral}$varTypeDef${_ansiErr} для параметра ${_ansiCmd}$param${_ansiErr} должно заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}$varTypeDef)"
            elif ! [[ $varTypeDef =~ $_isCalcEnumRegExp ]]; then
              enumDef="$varTypeDef"
            elif ! [[ $varTypeDef =~ \)[[:space:]]*\)[[:space:]]*$ ]]; then
              __err="ожидает, что вычисляемое определение перечислимого типа ${_ansiPrimaryLiteral}$varTypeDef${_ansiErr} для параметра ${_ansiCmd}$param${_ansiErr} должно заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}$varTypeDef)"
            else
              enumDef="$varTypeDef"
            fi
          elif [[ $varTypeDef =~ ^$_diapRegExp$ && -n ${BASH_REMATCH[0]} ]]; then
            valType=int
            local intMin="${BASH_REMATCH[1]}"
            local intMax="${BASH_REMATCH[3]}"
            if [[ -z $intMax ]]; then
              [[ -n ${BASH_REMATCH[2]} || -z $intMin ]] && intMax=Infinite || intMax="$intMin"
            fi
            [[ -n $intMin ]] || intMin=Infinite
            [[ $intMin == Infinite || $intMax == Infinite || $intMin -le $intMax ]] \
              || __err="ожидает, что в определении диапазона для параметра ${_ansiCmd}$(_quotedArgs --strip "$param")${_ansiErr} левая граница ${_ansiPrimaryLiteral}$intMin${_ansiErr} не превосходит правую ${_ansiSecondaryLiteral}$intMax"
          elif [[ -n $varTypeDef ]]; then
            __err="вместо ${_ansiPrimaryLiteral}:$varTypeDef${_ansiErr} в определении параметра ${_ansiCmd}$(_quotedArgs --strip "$param")${_ansiErr} ожидает:$(_echoVarTypeDefDescription)$_nl"
          fi; [[ -z $__err ]] || break
        fi
        [[ -n $varType ]] || varType=bool

        local __optionName; [[ $varKind == opt ]] && __optionName="--$__varName"
        local optionShortcuts=
        if [[ $varKind == opt && $__varName =~ ^(.*)/([[:alnum:]?]+)$ ]]; then
          __varName=${BASH_REMATCH[1]};
          optionShortcuts="${BASH_REMATCH[2]}"
        fi

        # [[ $varKind == opt && $__varName =~ - ]] && __varName=$(_kebabCaseToCamelCase "$__varName")
        [[ $varKind == opt ]] && dstVarName=__varName _kebabCaseToCamelCase "$__varName"
        if [[ -z $__varName ]]; then
          __err=" непустое имя переменной для параметра ${_ansiCmd}$param"
        elif [[ ! $__varName =~ $_isValidVarNameRegExp ]]; then
          __err=", что имя переменной ${_ansiPrimaryLiteral}$__varName${_ansiErr} для параметра ${_ansiCmd}$param${_ansiErr} $_mustBeValidVarName"
        fi
        [[ -z $__err ]] || { __err="ожидает$__err"; break; }
        if _hasItem $__varName ${__argVarNames[@]} ${__optVarNames[@]}; then
          local paramHolder="__PARAM_$__varName"
          __err="не может использовать то же самое имя переменной ${_ansiPrimaryLiteral}$__varName${_ansiErr} для параметра ${_ansiCmd}$param${_ansiErr}, что было использовано для параметра ${_ansiCmd}${!paramHolder}"
          break
        fi

        # $__varName is ready at this point

        if [[ $varType != bool ]]; then
          local __VALTYPE_$__varName=$valType
          local __EMPTY_$__varName="$canBeEmptyVarType"
          if [[ $valType == enum ]]; then
            local __ENUM_$__varName="$enumDef"
          elif [[ $valType == int ]]; then
            local __MIN_$__varName="$intMin"
            local __MAX_$__varName="$intMax"
          fi
        fi

        [[ -z $isImportant ]] || importantOptVarNames+=( $__varName )

        if [[ -n $__varValue ]]; then
          if [[ $varType == list ]]; then
            if [[ ! $__varValue =~ ^[[:space:]]*\( ]]; then
              __err="списочный параметр ${_ansiCmd}$param${_ansiErr} в качестве значения по умолчанию не будет иметь скалярное значение ${_ansiPrimaryLiteral}$__varValue${_ansiErr} (не заключенное в круглые скобки)"
            elif [[ ! $__varValue =~ \)[[:space:]]*$ ]]; then
              __err="списочное значение ${_ansiPrimaryLiteral}$__varValue${_ansiErr} параметра ${_ansiCmd}$param${_ansiErr} будет заканчиваться закрывающейся скобкой: ${_ansiSecondaryLiteral}$__varValue)"
            fi
          elif [[ $__varValue =~ ^[[:space:]]*\( ]]; then
            __err="скалярный параметр ${_ansiCmd}$param${_ansiErr} в качестве значения по умолчанию не будет иметь списочное значение ${_ansiPrimaryLiteral}$__varValue"
          fi
          [[ -z $__err ]] || { __err="ожидает, что $__err"; break; }
          if ! [[ $__varValue =~ \$  || ( $valType == enum && $enumDef =~ $_isCalcEnumRegExp ) ]]; then
            if [[ $varType != list ]]; then
              __infix= _validateDefaultCalculatedValue "$__varValue"
            else
              local -a listValues="$__varValue"
              local itemValue; for itemValue in "${listValues[@]}"; do
                __infix=" элемента" infix3="$__varValue" _validateDefaultCalculatedValue "$itemValue" || break
              done
            fi
            [[ -z $__err ]] || break
          fi
        fi

        if [[ $varKind == arg ]]; then
          __argVarNames+=( "$__varName" )
          __argIdx=$(( __argIdx + 1 ))
        elif [[ $varKind == opt ]]; then
          __optVarNames+=( "$__varName" )
          OPTVarNames+=( "OPT_$__varName" )
          local __i; for (( __i=0; __i<${#optionShortcuts}; __i++ )); do
            local optionShortcut="${optionShortcuts:$__i:1}"
            [[ "$allOptionsShortcuts" = *"$optionShortcut"* ]] || continue
            local optionVarName; for optionVarName in ${__optVarNames[@]}; do
              local shortcutsForOptionHolder="__SHORTCUTS_${optionVarName}"
              if [[ "${!shortcutsForOptionHolder}" == *"$optionShortcut"* ]]; then
                local optName; dstVarName=optName _upperCamelCaseToKebabCase $optionVarName
                __err="не может использовать для опции ${_ansiCmd}$__optionName${_ansiErr} то же сокращение ${_ansiPrimaryLiteral}$optionShortcut${_ansiErr}, что и для опции ${_ansiCmd}--$optName"
                break
              fi
            done; [[ -z $__err ]] || break
          done; [[ -z $__err ]] || break
          allOptionsShortcuts+="$optionShortcuts"
        fi

        code+=$_nl
        local indent=
        # [[ $funcName == _prepareCodeOfAutoHelp2TestFunc_alpha ]] && _debugVar __hasWrapper funcName
        if [[ -n $__hasWrapper ]]; then
          if [[ $varKind == arg ]]; then
            code+="! _hasItem $__varName \${__thisInheritedOptVarNames[@]} || return \$(_ownThrow \"не ожидает аргумент \${_ansiOutline}$__varName\${_ansiErr} одноименный унаследованной опции\")$_nl"
          else
            code+="if ! _hasItem $__varName \${__thisInheritedImportantOptVarNames[@]}; then$_nl"
            indent='  '
          fi
        fi

          code+="${indent}local __PARAM_$__varName="\'"$param"\'"$_nl"
          local __PARAM_$__varName="$param"
          code+="${indent}local __IMPORTANT_$__varName=$isImportant __KIND_$__varName=$varKind __TYPE_$__varName=$varType$_nl"
          local __IMPORTANT_$__varName=$isImportant __KIND_$__varName=$varKind __TYPE_$__varName=$varType
          code+="${indent}local __DEFAULT_$__varName="\'"$__varValue"\'"$_nl"
          local __DEFAULT_$__varName="$__varValue"
          if [[ $varKind == arg ]]; then
            local isOptional=
            if [[ $varType == scalar ]]; then
              [[ -z $__varValue && -z $canBeEmptyVarType ]] || isOptional=true
            else
              [[ $minListCount -gt 0 ]] || isOptional=true
            fi
            code+="${indent}local __OPTIONAL_$__varName=$isOptional$_nl"
            local __OPTIONAL_$__varName=$isOptional
          fi

          if [[ $varType != bool ]]; then
            code+="${indent}local __VALTYPE_$__varName=$valType __EMPTY_$__varName=$canBeEmptyVarType$_nl"
            if [[ $valType == int ]]; then
              code+="${indent}local __MIN_$__varName=$intMin __MAX_$__varName=$intMax$_nl"
            elif [[ $valType == enum ]]; then
              code+="${indent}local __ENUM_$__varName="\'$enumDef\'"$_nl"
              code+="${indent}local __codeIsValidEnumValue_$__varName='"
              if [[ $enumDef =~ \$ ]]; then
                code+='eval "$_codeToPrepareEnumValues" && _hasItem "$__varValue" "${enumValues[@]}"'
              else
                eval "$_codeToPrepareEnumValues"
                local codeAccum=
                local enumValue; for enumValue in "${enumValues[@]}"; do
                  [[ -z $codeAccum ]] || codeAccum+=' || '
                  codeAccum+='$__varValue == '$(_quotedArgs --quote:all "$enumValue")
                done
                code+="[[ $codeAccum ]]"
              fi
              code+="'"$_nl
            fi
          fi

          if [[ $varType != list ]]; then
            code+="${indent}local $__varName=\"\"$_nl"
          else
            code+="${indent}local -a $__varName=( )$_nl"
            code+="${indent}local __MIN_COUNT_$__varName=$minListCount __MAX_COUNT_$__varName=$maxListCount __SORTED_$__varName=$sortedVarType __UNIQUE_$__varName=$uniqueVarType$_nl"
            local __MIN_COUNT_$__varName=$minListCount __MAX_COUNT_$__varName=$maxListCount __SORTED_$__varName=$sortedVarType __UNIQUE_$__varName=$uniqueVarType
          fi

          if [[ $varKind == opt ]]; then
            code+="${indent}local -a OPT_$__varName=( )$_nl"
            local optName; dstVarName=optName _upperCamelCaseToKebabCase $__varName
            code+="${indent}local __OPTNAME_$__varName=--$optName$_nl"
            code+="${indent}local __SHORTCUTS_$__varName=\"$optionShortcuts\"$_nl"
            local __i; for (( __i=0; __i<${#optionShortcuts}; __i++ )); do
              local optionShortcut="${optionShortcuts:$__i:1}"
              if [[ $optionShortcut != \? ]]; then
                code+="${indent}local __VAR_NAME_FOR_SHORTCUT_$optionShortcut=\"$__varName\"$_nl"
                local __VAR_NAME_FOR_SHORTCUT_$optionShortcut=$__varName
              fi
            done
            local __SHORTCUTS_$__varName="$optionShortcuts"
          fi

        if [[ -n $__hasWrapper && $varKind == opt ]]; then
          code+="fi$_nl"
          indent=
        fi

        prevParam="$param"
      done; [[ -z $__err ]] || break


      if [[ -n $isCommandWrapper ]]; then
        local -a __funcNameSuffixes=(); help= _prepareSubCommandFuncSuffixes ${funcName} || return $?
        _validateSubCommands
        [[ -z $__err ]] || break
        local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
          funcName="${funcName}_$__funcNameSuffix" needPrepare=Params2 hasWrapper=true eval "$_codeToPregen"
        done
      fi

      local precode="$_nl"
      precode+=local
      for __varName in funcName treatUnknownOptionAsArg canBeMixedOptionsAndArgs canBeMoreParams isCommandWrapper additionalSuffix; do
        precode+=" __$__varName=${!__varName}"
      done
      precode+=$_nl
      precode+='local __thisFuncCommand=${__funcCommand:-$__funcName} __funcCommand= __isCompletionMode='$_nl
      precode+='local -a __additionalDependencies=( '
      local fileSpec; for fileSpec in "${additionalDependencies[@]}"; do
        [[ $fileSpec =~ ^$_bwDir ]] \
          && fileSpec='$_bwDir'${fileSpec:${#_bwDir}} \
          || fileSpec=$(_shortenFileSpec "$fileSpec")
        precode+='"'$fileSpec'"'
      done
      precode+=' )'$_nl

      precode+="local -a __argVarNames=( ${__argVarNames[@]} )$_nl"

      if [[ -n $__hasWrapper ]]; then
        local  Holder=${funcName}SuppressInheritedOptionVarNames
        dstVarName=suppressInheritedOptionVarNames srcVarName=${funcName}SuppressInheritedOptionVarNames eval "$_codeToInitLocalCopyOfArray"
        if [[ ${#suppressInheritedOptionVarNames[@]} -eq 0 ]]; then
          precode+='local -a __thisInheritedOptVarNames=( ${__inheritedOptVarNames[@]} )'$_nl
          precode+='local -a __thisInheritedOPTVarNames=( ${__inheritedOPTVarNames[@]} )'$_nl
          precode+='local -a __thisInheritedImportantOptVarNames=( ${__inheritedImportantOptVarNames[@]} )'$_nl
        else
          local hasItemCode=
          local item; for item in "${suppressInheritedOptionVarNames[@]}"; do
            if [[ ! $item =~ $_isValidVarNameRegExp ]]; then
              __err=", что имя переменной ${_ansiPrimaryLiteral}$item${_ansiErr} в ${_ansiOutline}$suppressInheritedOptionVarNamesHolder${_ansiErr} $_mustBeValidVarName"
              break 2
            fi
            [[ -z $hasItemCode ]] || hasItemCode+=' || '
            hasItemCode+='$__varName == '$item
          done
          precode+='local -a __thisInheritedOptVarNames'$_nl
          precode+='local -a __thisInheritedOPTVarNames'$_nl
          precode+='local -a __thisInheritedImportantOptVarNames'$_nl
          precode+='local __varName; for __varName in "${__inheritedOptVarNames[@]}"; do'$_nl
          precode+="  ! [[ $hasItemCode ]] || continue"$_nl
          precode+='  __thisInheritedOptVarNames+=( $__varName )'$_nl
          precode+='  __thisInheritedOPTVarNames+=( OPT_$__varName )'$_nl
          precode+='done'$_nl
          precode+='local __varName; for __varName in "${__inheritedImportantOptVarNames[@]}"; do'$_nl
          precode+="  ! [[ $hasItemCode ]] || continue"$_nl
          precode+='  __thisInheritedImportantOptVarNames+=( $__varName )'$_nl
          precode+='done'$_nl
        fi
        precode+='local -a __inheritedOptVarNames=()'$_nl
        precode+='local -a __inheritedOPTVarNames=()'$_nl
        precode+='local -a __inheritedImportantOptVarNames=()'$_nl
      fi
      precode+='local -a __optVarNames=( '
      [[ -z $__hasWrapper ]] \
        || precode+='"${__thisInheritedOptVarNames[@]}" '
      precode+=$(_quotedArgs ${__optVarNames[@]})' )'$_nl
      precode+='local -a __OPTVarNames=( '
      [[ -z $__hasWrapper ]] \
        || precode+='"${__thisInheritedOPTVarNames[@]}" '
      precode+=$(_quotedArgs ${OPTVarNames[@]})' )'$_nl

      local codeHasVarNameInOptVarNames=
      local __varName; for __varName in ${__optVarNames[@]}; do
        [[ -z $codeHasVarNameInOptVarNames ]] \
          || codeHasVarNameInOptVarNames+=' || '
        codeHasVarNameInOptVarNames+='$__varName == '$__varName
      done
      precode+="local __codeIsOptVarName='[[ $codeHasVarNameInOptVarNames ]]"
      [[ -z $__hasWrapper ]] \
        || precode+=' || _hasItem $__varName ${__thisInheritedOptVarNames[@]}'
      precode+="'"$_nl

      code="$precode$code"

      [[ -z $isCommandWrapper ]] \
        || code+="${_nl}local -a __thisInheritedImportantOptVarNames=( \${__thisInheritedImportantOptVarNames[@]} ${importantOptVarNames[@]} )"$_nl

      _assureDir $(dirname "$codeFileSpec") \
        || return $?
      echo -n "$code" > "$codeFileSpec"
      [[ -z $verbose ]] || _ok "Создан ${_ansiFileSpec}$codeFileSpec"

    fi

    if [[ -n $__thisOnlyPrepareCode ]]; then
      if [[ ! $funcName =~ Complete$ ]]; then
        local completionCodeFileSpec; dstVarName=completionCodeFileSpec codeType=completion fileSpec= originalCodeDeep= eval "$_codeToPrepareCodeFileSpec"
        if [[ ! -f $completionCodeFileSpec ]]; then
          local completeFuncName="__complete_${funcName}"
          echo '
'$completeFuncName'() {
local __consumedParams=0
. "$_bwDir'${codeFileSpec:${#_bwDir}}'"
__isCompletionMode=true _parseFuncParams2
}
complete -o nospace -F '$completeFuncName' '$funcName > "$completionCodeFileSpec"

          local unsetFileSpec="${completionCodeFileSpec:0:$(( ${#completionCodeFileSpec} - ${#_codeBashExt} ))}$_unsetFileExt"
          echo '
complete -p | grep " '$funcName'\$" >/dev/null 2>&1 && complete -r '$funcName'
unset -f '$completeFuncName > "$unsetFileSpec"

          [[ -z $verbose ]] || _ok "Создан ${_ansiFileSpec}$completionCodeFileSpec"
        fi
        __completions+=( "$completionCodeFileSpec" )
      fi

      [[ -z $needPregenHelp ]] \
        || __ownPrefix= _prepareCodeOfAutoHelp2
    fi

    [[ -z $__thisOnlyPrepareCode ]] || return 2
    break
  done
  _profileEnd
  [[ -z $__err ]] || _err --showStack 4 "${_ansiCmd}${FUNCNAME[1]}${_ansiErr} $__err"
}

# =============================================================================

_validateSubCommands() {
  while true; do
    if [[ ${#__funcNameSuffixes[@]} -eq 0 ]]; then
      __err="ожидает, что будет определена по крайней мере одна функция ${_ansiCmd}${funcName}_${_ansiOutline}*"
      break
    fi
    local -a allCommandShortcuts=()
    local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
      dstVarName=__commandShortcuts srcVarName=${funcName}_${__funcNameSuffix}Shortcuts eval "$_codeToInitLocalCopyOfArray"
      local __commandShortcut; for __commandShortcut in "${__commandShortcuts[@]}"; do
        [[ ! $__commandShortcut =~ ^[[a-z0-9_][[a-z0-9_-]*$ ]] || continue
        local __commandName; dstVarName=__commandName _upperCamelCaseToKebabCase $__funcNameSuffix
        __err="ожидает, что сокращение ${_ansiPrimaryLiteral}$(_quotedArgs "$__commandShortcut")${_ansiErr} для команды ${_ansiCmd}$__commandName $_mustBeValidCommandShortcut"
        break
      done
      allCommandShortcuts+=( "${__commandShortcuts[@]}" )
    done
    local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
      local __commandName; dstVarName=__commandName _upperCamelCaseToKebabCase $__funcNameSuffix
      if _hasItem "$__commandName" "${allCommandShortcuts[@]}"; then
        local anotherCommandName=
        local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
          dstVarName=__commandShortcuts srcVarName=${funcName}_${__funcNameSuffix}Shortcuts eval "$_codeToInitLocalCopyOfArray"
          if _hasItem "$__commandName" "${__commandShortcuts[@]}"; then
            local anotherCommandName; dstVarName=anotherCommandName _upperCamelCaseToKebabCase $__funcNameSuffix
            break
          fi
        done
        __err="ожидает, что команда ${_ansiCmd}$__commandName${_ansiErr} не будет совпадать с одним из сокращений для команды ${_ansiCmd}$anotherCommandName${_ansiErr}: ${_ansiSecondaryLiteral}$(_quotedArgs "${__commandShortcuts[@]}")"
        break
      fi
    done
    break
  done
}

# =============================================================================

_prepareSubCommandFuncSuffixesParams=(
  '--checkCondition'
  'parentCommand'
)
_prepareSubCommandFuncSuffixes() { eval "$_funcParams2"
  _profileBegin
  __funcNameSuffixes=()
  local declare='declare -f '
  local prefixLen=$(( ${#parentCommand} + 1 ))
  local funcName; for funcName in $(compgen -c ${parentCommand}_ | grep -v -E '(Complete|Helper|Params)$' ); do
    if [[ -n $checkCondition ]]; then
      local holder=${funcName}Condition
      [[ -z ${!holder} ]] || eval "${!holder}" || continue
    fi
    local funcNameSuffix=${funcName:$prefixLen}
    [[ ! $funcNameSuffix =~ _ ]] || continue
    __funcNameSuffixes+=( $funcNameSuffix )
  done
  _profileEnd
  return 0
}

# =============================================================================

_echoVarTypeDefDescription() {
  local listItem; [[ $varType == list ]] && listItem=" элемента списка"
  local valueTypeInfo="$_nl-- (необязательно) один из следующих типов$listItem:" __item; for __item in "${_codeToParseVarValueTypeInfo[@]}"; do
    if [[ ${__item:0:4} == arg: ]]; then
      [[ $varKind != arg ]] && continue || __item=${__item:4}
    elif [[ ${__item:0:4} == opt: ]]; then
      [[ $varKind != opt ]] && continue || __item=${__item:4}
    fi
    valueTypeInfo+="$_nl$(_indent 4)$__item${_ansiErr}"
  done
  local emptyTypeInfo="$_nl-- и/или (необязательно) указание пустого значения$listItem:" __item; for __item in "${_codeToParseVarEmptyTypeInfo2[@]}"; do
    emptyTypeInfo+="$_nl$(_indent 4)$__item${_ansiErr}"
  done
  local listTypeInfo=
  if [[ $varType == list ]]; then
    listTypeInfo="$_nl-- и/или (необязательно) указание инструкции для формирования списка значений:"
    for __item in "${_codeToParseVarListTypeInfo[@]}"; do
      listTypeInfo+=$_nl$(_indent 4)"$__item${_ansiErr}"
    done
  fi
  echo "${valueTypeInfo}${emptyTypeInfo}${listTypeInfo}"
}

# =============================================================================

_completeParams=(
  '--__varName'
  '--__argIdx:0..'
  'compWord:?'
)

# =============================================================================

_codeToAddCandidateToCompletion='
  local regExp="$__curWord"
  [[ $__curWord =~ \? ]] && regExp=${__curWord//\?/\\?}
  [[ ! $__candidate =~ ^$regExp ]] || _hasItem "$__candidate " "${COMPREPLY[@]}" || COMPREPLY+=( "$__candidate " )
'
_parseFuncParams2() {
  _profileBegin
  local __argCount=${#__argVarNames[@]} __argIdx=0
  local __err __varName __varValue
  if [[ -n $__curWordIdx ]]; then
    local __thisCurWordIdx=$__curWordIdx
    __curWordIdx=
  elif [[ -z $__isCompletionMode ]]; then
    local __thisCurWordIdx=1
    local -a __words=( $__funcName "$@" )
  else
    local __thisCurWordIdx=1
    local -a __words=( "${COMP_WORDS[@]}" )
  fi
  local __initAsEmpty=; [[ -z $__isCompletionMode ]] || __initAsEmpty=true


  if [[ -n $__isCompletionMode ]]; then
    local suffix; for suffix in Reset Bold Dim ResetBold ResetDim Underline Header Url Cmd FileSpec Dir Err Warn Will OK Outline Debug PrimaryLiteral SecondaryLiteral; do
      eval local _ansi$suffix=
    done
  fi

  local __curWord="${__words[$__thisCurWordIdx]}"
  while true; do

    while [[ $__thisCurWordIdx -lt ${#__words[@]} ]]; do
      if [[ $__curWord == -- && (-z $__isCompletionMode || $__thisCurWordIdx -ne $COMP_CWORD) ]]; then
        eval "$_codeToNextWord"
        break
      fi

      if [[ $__argIdx -eq 0 || -n $__canBeMixedOptionsAndArgs ]]; then

        if [[ -n $__isCompletionMode && $__thisCurWordIdx -eq $COMP_CWORD && $__curWord =~ ^- ]]; then
          _setCompReplyToOptionNames
          break 2
        fi
        if [[ $__curWord =~ ^-- || $__curWord =~ ^-[^[:digit:]]$ ]]; then

          if [[ $__curWord =~ ^-- ]]; then
            __varName="${__curWord:2}"
            dstVarName=__varName _kebabCaseToCamelCase "${__curWord:2}"
          else
            local char="${__curWord:1}"
            eval "$_codeToPrepareVarNameByOptionShortcut"
          fi
          eval "$__codeIsOptVarName" || __varName=


          if [[ -z $__varName ]]; then
            if [[ -z $__treatUnknownOptionAsArg ]]; then
              __err="не ожидает опцию ${_ansiCmd}$__curWord"
              break
            fi
          else
            local __varTypeHolder="__TYPE_$__varName"
            if [[ ${!__varTypeHolder} == "bool" ]]; then
              eval $__varName=true
              eval OPT_$__varName=\( "\$__curWord" \)
            else
              local __optionName="$__curWord"
              eval "$_codeToNextWord"

              if [[ -n $__isCompletionMode && $__thisCurWordIdx -eq $COMP_CWORD ]]; then
                _setCompReplyToValueDescription
                break 2
              fi

              local __expects="ожидает, что опция ${_ansiCmd}$__optionName${_ansiErr} будет снабжена"
              if [[ $__thisCurWordIdx -ge ${#__words[@]} ]]; then
                __err="$__expects ${_ansiOutline}значением"
              # elif [[ $__curWord =~ ^-.+  && ! $__curWord =~ ^-[[:digit:]]+$ ]]; then
                # __err="$__expects значением (${_ansiPrimaryLiteral}$__curWord${_ansiErr}) не похожим на опцию"
              else
                local __varValue="$__curWord"
                local __errSuffixCode=
                if [[ -z $__isCompletionMode ]] && ! _validateVarValue "$__varValue" ; then
                  [[ $__errSuffixCode -eq 0 ]] \
                    && __err+=" для" \
                    || __err+=" в качестве значения"
                  __err+=" опции ${_ansiCmd}$__optionName"
                elif [[ ${!__varTypeHolder} == "list" ]]; then
                  eval $__varName+=\( \"\$__curWord\" \)
                  eval OPT_$__varName+=\( \"\$__optionName\" \"\$__varValue\" \)
                else
                  local __varValueHolder="OPT_$__varName"
                  if [[ -z ${!__varValueHolder} ]]; then
                    eval $__varName=\"\$__curWord\"
                    eval OPT_$__varName=\( \"\$__optionName\" \"\$__varValue\" \)
                  else
                    __varValueHolder="$__varName"
                    [[ ${!__varValueHolder} == $__curWord ]] \
                      || __err="не ожидает, что опция ${_ansiCmd}$__optionName${_ansiErr} будет указана повторно с другим значением ${_ansiPrimaryLiteral}$__varValue${_ansiErr} против ${_ansiPrimaryLiteral}${!__varValueHolder}${_ansiErr}, указанного первоначально"
                  fi
                fi
              fi
              [[ -z $__err ]] || break
            fi
            eval "$_codeToNextWord"
            continue
          fi

        elif [[ $__curWord =~ ^-.+  && ! $__curWord =~ ^-[[:digit:]] ]]; then
          local -a __shortcutsVarNames=()
          local __treatAsArg=
          local j; for (( j=1; j < ${#__curWord}; j++ )); do
            local char="${__curWord:j:1}"
            eval "$_codeToPrepareVarNameByOptionShortcut"
            eval "$__codeIsOptVarName" || __varName=
            if [[ -n $__varName ]]; then
              local __varTypeHolder="__TYPE_$__varName"
              if [[ ${!__varTypeHolder} != "bool" ]]; then
                local optName=; dstVarName=optName _upperCamelCaseToKebabCase "$__varName"
                __err="ожидает значение для опции ${_ansiCmd}--$optName${_ansiErr}, поэтому её краткая форма ${_ansiCmd}-$char${_ansiErr} не может быть использоваана в ${_ansiOutline}объединении опций ${_ansiCmd}$__curWord"
                break
              else
                __shortcutsVarNames+=( "$__varName" )
                local __OPT_$__varName="-$char"
              fi
            elif [[ -z $__treatUnknownOptionAsArg ]]; then
              __err="не ожидает краткую опцию ${_ansiPrimaryLiteral}$char${_ansiErr} в ${_ansiOutline}объединении опций ${_ansiCmd}$__curWord"
              break
            else
              __treatAsArg=true
              break
            fi
          done; [[ -z $__err ]] || break
          if [[ -z $__treatAsArg ]]; then
            for __varName in "${__shortcutsVarNames[@]}"; do
              eval $__varName=true
              eval OPT_$__varName=\( \"\$__OPT_$__varName\" \)
            done
            eval "$_codeToNextWord"
            continue
          fi
        fi
      fi

      if [[ -z $__isCommandWrapper ]]; then

        if [[ -n $__isCompletionMode && $__thisCurWordIdx -eq $COMP_CWORD ]]; then
          _setCompReplyToArgDescription
          break 2
        fi

        if [[ $__argCount -gt 0 ]]; then
          __argIdx=$(( __argIdx + 1 ))
          [[ $__argIdx -le $__argCount ]] \
            && local __min=$__argIdx \
            || local __min=$__argCount
          local __varName=${__argVarNames[$(( __min - 1 ))]}
          local __varTypeHolder="__TYPE_$__varName"
          if [[ $__argIdx -le $__argCount || ${!__varTypeHolder} == list ]]; then
            local __varValue="$__curWord"; eval "$_codeToNextWord"
            if [[ -z $__isCompletionMode ]] && ! _validateVarValue "$__varValue" ; then
              __err="$__err в качестве $__argIdx-го аргумента ${_ansiOutline}$__varName"
              break
            fi
            [[ ${!__varTypeHolder} == scalar ]] \
              && eval $__varName=\"\$__varValue\" \
              || eval $__varName+=\( \"\$__varValue\" \)
          elif [[ -n $__canBeMoreParams ]]; then
            break
          elif [[ -z $__isCompletionMode || -n $__curWord ]]; then
            __err="ожидает не более $__argCount $(_getPluralWord $__argCount аргумента аргументов), но обнаружен $__argIdx-й"
          fi
        elif [[ -n $__canBeMoreParams ]]; then
          break
        elif [[ -z $__isCompletionMode || -n $__curWord ]]; then
          __err="не ожидает ни одного аргумента, но обнаружен"
        fi
        [[ -z $__err ]] || { __err="$__err: ${_ansiCmd}$(_quotedArgs "$__curWord")"; break; }

      else

        if [[ -n $__isCompletionMode && $__thisCurWordIdx -eq $COMP_CWORD ]]; then
          _setCompReplyToSubCommand
          break 2
        fi

        [[ -n $__canBeMixedOptionsAndArgs ]] || _postProcessVarNames || break
        local __usedSubCommand="$__curWord"
        local __funcSuffix=
        local -a __funcNameSuffixes=(); help= _prepareSubCommandFuncSuffixes --checkCondition $__funcName || return $?
        local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
          local subCommand; dstVarName=subCommand _upperCamelCaseToKebabCase $__funcNameSuffix
          dstVarName=__subCommands srcVarName=${__funcName}_${__funcNameSuffix}Shortcuts eval "$_codeToInitLocalCopyOfArray"
          __subCommands+=( "$subCommand" )
          _hasItem "$__usedSubCommand" "${__subCommands[@]}" && __funcSuffix="$__funcNameSuffix" && break
        done; [[ -z $__err ]] || break
        if [[ -z $__funcSuffix ]]; then
          __err="вместо ${_ansiPrimaryLiteral}$(_quotedArgs "$__usedSubCommand")${_ansiErr} ожидает одну из следующих команд: ${_ansiSecondaryLiteral}$(_echoAllSubCommands)"
          break
        fi
        eval "$_codeToNextWord"

        if [[ -n $__canBeMixedOptionsAndArgs ]]; then
          local -a __inheritedOptVarNames=( "${__thisInheritedOptVarNames[@]}" "${__optVarNames[@]}" )
          local -a __inheritedOPTVarNames=( "${__thisInheritedOPTVarNames[@]}" "${__OPTVarNames[@]}" )
        fi
        local -a __inheritedImportantOptVarNames=( "${__thisInheritedImportantOptVarNames[@]}" "${importantOptVarNames[@]}"  )
        local __curWordIdx=$__thisCurWordIdx
        local __funcCommand="$__thisFuncCommand $__usedSubCommand"
        local __subFuncName="${__funcName}_$__funcSuffix"
        if [[ -z $__isCompletionMode ]]; then
          $__subFuncName
        else
          __complete_$__subFuncName
        fi
        local __returnCode=$?

        return $__returnCode
      fi

    done;
    [[ -z $__err ]] || break

    [[ -z $__isCompletionMode ]] || break

    if [[ -n $help ]]; then
      local __helpCodeFileSpec; __ownPrefix=__ help= _prepareCodeOfAutoHelp2 && cmdName="$__thisFuncCommand" . "$__helpCodeFileSpec"; local __returnCode=$?
      [[ $__returnCode -eq 0 ]] \
        && return 3 \
        || return $?
    elif [[ -n $__isCommandWrapper ]]; then
      local -a __funcNameSuffixes=(); help= _prepareSubCommandFuncSuffixes --checkCondition $__funcName
      __err="в качестве первого аргумента ожидает одну из следующих команд: ${_ansiSecondaryLiteral}$(_echoAllSubCommands)"
      break
    else
      __argIdx=$(( __argIdx + 1 ))
      while [[ $__argIdx -le $__argCount ]]; do
        local __varName=${__argVarNames[$(( __argIdx - 1 ))]}
        local __varTypeHolder="__TYPE_$__varName"
        if [[ ${!__varTypeHolder} == scalar ]]; then
          local __optionalHolder="__OPTIONAL_$__varName"
          if [[ -z ${!__optionalHolder} ]]; then
            __err="ожидает $__argIdx-й аргумент ${_ansiOutline}$__varName"
            break
          fi
        fi
        __argIdx=$(( __argIdx + 1 ))
      done; [[ -z $__err ]] || break
      _postProcessVarNames || break
      local __thisFuncCommand=
    fi
    break
  done
  _profileEnd
  [[ -n $__err ]] || return 0
  [[ -n $__isCompletionMode ]] \
    && _setCompReplyToHint "$__err" \
    || _ownThrow "$__err"
}

_codeToNextWord='
  __thisCurWordIdx=$(( __thisCurWordIdx + 1 )); __curWord="${__words[$__thisCurWordIdx]}"; __consumedParams=$(( __consumedParams + 1 ))
'

_setCompReplyToOptionNames() {
  local __count=0
  local __varName; for __varName in ${__optVarNames[@]}; do
    local __len=${#COMPREPLY[@]}
    local __varTypeHolder="__TYPE_$__varName"
    if [[ ${!__varTypeHolder} != list ]]; then
      dstVarName=__OPT srcVarName=OPT_$__varName eval "$_codeToInitLocalCopyOfArray"
      [[ ${#__OPT[@]} -eq 0 ]] || continue
    else
      local __maxCountHolder="__MAX_COUNT_$__varName"
      if [[ ${!__maxCountHolder} != Infinite ]]; then
        dstVarName=__OPT srcVarName=OPT_$__varName eval "$_codeToInitLocalCopyOfArray"
        [[ ${#__OPT[@]} -lt $(( ${!__maxCountHolder} * 2 )) ]] || continue
      fi
    fi
    __candidate=--$__varName eval "$_codeToAddCandidateToCompletion"
    local __optName; dstVarName=__optName _upperCamelCaseToKebabCase $__varName
    [[ $__optName == $__varName ]] || __candidate=--$__optName eval "$_codeToAddCandidateToCompletion"
    dstVarName=__shortcuts srcVarName=__SHORTCUTS_$__varName eval "$_codeToInitLocalCopyOfScalar"
    local __i; for (( __i=0; __i < ${#__shortcuts}; __i++ )); do
      __candidate=-${__shortcuts:$__i:1} eval "$_codeToAddCandidateToCompletion"
    done
    [[ ${#COMPREPLY[@]} -eq $__len ]] || __count=$(( __count + 1 ))
  done
  [[ $__count -ne 1 ]] || COMPREPLY=( "${COMPREPLY[0]}" )
}

_setCompReplyToHint() {
  COMPREPLY=( "<-- HINT: ${__thisFuncCommand:-${FUNCNAME[1]}} $@" "-->")
}

_setCompReplyToSubCommand() {
  local __count=0
  local -a __funcNameSuffixes=(); help= _prepareSubCommandFuncSuffixes --checkCondition $__funcName
  local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
    local __len=${#COMPREPLY[@]}
    dstVarName=__subCommands srcVarName=${__funcName}_${__funcNameSuffix}Shortcuts eval "$_codeToInitLocalCopyOfArray"
    local subCommand; dstVarName=subCommand _upperCamelCaseToKebabCase $__funcNameSuffix
    local __candidate; for __candidate in $subCommand "${__subCommands[@]}"; do
      eval "$_codeToAddCandidateToCompletion"
    done
    [[ ${#COMPREPLY[@]} -eq $__len ]] || __count=$(( __count + 1 ))
  done
  [[ $__count -ne 1 ]] || COMPREPLY=( "${COMPREPLY[0]}" )
}

_setCompReplyToArgDescription() {
  if [[ $__argCount -gt 0 ]]; then
    __argIdx=$(( __argIdx + 1 ))
    [[ $__argIdx -le $__argCount ]] \
      && local __min=$__argIdx \
      || local __min=$__argCount
    local __varName=${__argVarNames[$(( __min - 1 ))]}
    local __varTypeHolder="__TYPE_$__varName"
    if [[ $__argIdx -le $__argCount || ${!__varTypeHolder} == list ]]; then
      _setCompReplyToValueDescription
    fi
  fi
}

_setCompReplyToValueDescription() {
  local __varValtypeHolder="__VALTYPE_$__varName"
  if [[ ${!__varValtypeHolder} == int ]]; then
    local valueDescription=; _prepareIntValudDescription2Helper
    eval local valueDescription=\"$valueDescription\"
    local emptyHolder="__EMPTY_$__varName"
    [[ -z ${!emptyHolder} ]] || valueDescription+=" или пустое значение"
    _setCompReplyToHint "ожидает $valueDescription"
  elif [[ ${!__varValtypeHolder} == enum ]]; then
    eval "$_codeToPrepareEnumValues"
    local enumValue; for enumValue in "${enumValues[@]}"; do
      __candidate="$enumValue" eval "$_codeToAddCandidateToCompletion"
    done
  else
    local completeFuncName="${__funcName}Complete"
    if _funcExists $completeFuncName; then
      [[ -n $__varKindHolder ]] || local __varKindHolder="__KIND_$__varName"
      $completeFuncName
    elif [[ -n ${!emptyHolder} ]]; then
      _setCompReplyToHint "ожидает возможно пустое значение"
    fi
  fi
}

_echoAllSubCommands() {
  local -a allSubCommands=()
  local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
    dstVarName=__subCommands srcVarName=${FUNCNAME[2]}_${__funcNameSuffix}Shortcuts eval "$_codeToInitLocalCopyOfArray"
    local subCommand; dstVarName=subCommand _upperCamelCaseToKebabCase $__funcNameSuffix
    allSubCommands+=( $subCommand "${__subCommands[@]}" )
  done
  echo "${allSubCommands[@]}"
}

_codeToPrepareEnumValues='
  local __varEnumHolder="__ENUM_$__varName"
  if [[ ${!__varEnumHolder} =~ $_isCalcEnumRegExp ]]; then
    eval eval local -a enumValues="${!__varEnumHolder}"
  else
    eval local -a enumValues="${!__varEnumHolder}"
  fi
'
_postProcessVarNames() {
  local returnCode=0
  local __varName; for __varName in ${__argVarNames[@]} ${__optVarNames[@]}; do
    local __varTypeHolder="__TYPE_$__varName"
    [[ ${!__varTypeHolder} != bool ]] || continue

    local __defaultValueHolder="__DEFAULT_$__varName"
    if [[ ${!__varTypeHolder} == scalar ]]; then
      if [[ -n ${!__defaultValueHolder} ]]; then
        local __varKindHolder="__KIND_$__varName"
        local __varValueHolder
        [[ ${!__varKindHolder} == arg ]] \
          && __varValueHolder="$__varName" \
          || __varValueHolder="OPT_$__varName"
        if [[ -z ${!__varValueHolder} ]]; then
          eval $__varName="${!__defaultValueHolder}"
          local __varValtypeHolder="__VALTYPE_$__varName"
          local __varEnumHolder="__ENUM_$__varName"
          if [[ \
            -z $__isCompletionMode && ( \
              ${!__defaultValueHolder} =~ \$ \
              || ( ${!__varValtypeHolder} == enum && ${!__varEnumHolder} =~ $_isCalcEnumRegExp ) \
            ) \
          ]]; then
            local __infix2=; [[ ${!__defaultValueHolder} =~ \$ ]] && __infix2=" вычисляемого"
            __infix= _validateDefaultCalculatedValue "${!__varName}"; returnCode=$?
          fi
          if [[ $returnCode -eq 0 && ${!__varKindHolder} == opt ]]; then
            local __varValueHolder="$__varName"
            local __optValue="${!__varValueHolder}"

            local __optNameHolder=__OPTNAME_$__varName
            eval OPT_$__varName=\( \${!__optNameHolder} \"\$__optValue\" \)
          fi
        fi
      fi
    else
      local -a __list=()
      local __count __isReadyList=
      if [[ -n ${!__defaultValueHolder} ]]; then
        _prepareList
        if [[ $__count -eq 0 ]]; then
          eval $__varName="${!__defaultValueHolder}"
          __isReadyList=
          local __varValtypeHolder="__VALTYPE_$__varName"
          local __varEnumHolder="__ENUM_$__varName"
          if [[ \
            -z $__isCompletionMode && ( \
              ${!__defaultValueHolder} =~ \$ \
              || ( \
                ${!__varValtypeHolder} == enum \
                && ${!__varEnumHolder} =~ $_isCalcEnumRegExp \
              ) \
            ) \
          ]]; then
            local __infix2=; [[ ${!__defaultValueHolder} =~ \$ ]] && __infix2=" вычисляемого"
            _prepareList
            local __item; for __item in "${__list[@]}"; do
              __infix=" элемента" _validateDefaultCalculatedValue "$__item"; returnCode=$?
              [[ $returnCode -eq 0 ]] || break
            done
          fi
          if [[ $returnCode -eq 0 ]]; then
            local __varKindHolder="__KIND_$__varName"
            if [[ ${!__varKindHolder} == opt ]]; then
            local __optNameHolder=__OPTNAME_$__varName
              _prepareList
              local __optValue; for __optValue in "${__list[@]}"; do
                eval OPT_$__varName+=\( \${!__optNameHolder} \"\$__optValue\" \)
              done
            fi
          fi
        fi
      fi

      if [[ $returnCode -eq 0 ]]; then
        local __sortedHolder="__SORTED_$__varName"
        local __uniqueHolder="__UNIQUE_$__varName"
        if [[ -n ${!__sortedHolder} || -n ${!__uniqueHolder} ]]; then
          _prepareList
          if [[ $__count -gt 1 ]]; then
            local -a __result=()
            local __didChange=
            if [[ -n ${!__sortedHolder} ]]; then
              local __varValtypeHolder="__VALTYPE_$__varName"
              if [[ ${!__varValtypeHolder} != enum ]]; then
                local -a __sortOpts=()
                [[ ${!__varValtypeHolder} == int ]] \
                  && __sortOpts+=( -n ) \
                  || __sortOpts+=( -d )
                [[ -n ${!__uniqueHolder} ]] && __sortOpts+=( -u )
                IFS="$_nl" __result=( $(sort "${__sortOpts[@]}" <<<"${__list[*]}" ) ); unset IFS # https://stackoverflow.com/questions/7442417/how-to-sort-an-array-in-bash
              else
                eval "$_codeToPrepareEnumValues"
                local __value; for __value in "${enumValues[@]}"; do
                  _hasItem "$__value" "${__list[@]}" || continue
                  if [[ -n ${!__uniqueHolder} ]]; then
                    __result+=( "$__value" )
                  else
                    local __count=0
                    local __i; for (( __i=0; __i < ${#__list[@]}; __i++ )); do
                      [[ ${__list[$__i]} == "$__value" ]] && __count=$(( __count + 1 ))
                    done
                    local __i; for (( __i=0; __i < $__count; __i++ )); do
                      __result+=( "$__value" )
                    done
                  fi
                done
              fi
              __didChange=true
            else
              local __value; for __value in "${__list[@]}"; do
                _hasItem "$__value" "${__result[@]}" \
                  && __didChange=true \
                  || __result+=( "$__value" )
              done
            fi
            if [[ -n $__didChange ]]; then
              eval $__varName=\( \"\${__result[@]}\" \)
              __list=( "${__result[@]}" )
              __count=${#__list[@]}
            fi
          fi
        fi
        if [[ -z $__isCompletionMode ]]; then
          local __minCountHolder="__MIN_COUNT_$__varName"
          local __maxCountHolder="__MAX_COUNT_$__varName"
          local __needCheckMinCount=; [[ ${!__minCountHolder} -eq 0 ]] || __needCheckMinCount=true
          local __needCheckMaxCount=; [[ ${!__maxCountHolder} == Infinite ]] || __needCheckMaxCount=true
          if [[ -n $__needCheckMinCount || -n $__needCheckMaxCount ]]; then
            local __isOutOfListLimit= __got= __expects= __expectsLimitKind= __expectsLimitCount
            _prepareList
            if [[ -n $__needCheckMinCount && $__count -lt ${!__minCountHolder} ]]; then
              __isOutOfListLimit=true
              __expectsLimitKind="не менее"
              __expectsLimitCount=${!__minCountHolder}
              [[ $__count -eq 0 ]] && __got="не получено ничего";
            elif [[ -n $__needCheckMaxCount && $__count -gt ${!__maxCountHolder} ]]; then
              __isOutOfListLimit=true
              __expectsLimitKind="не более"
              __expectsLimitCount=${!__maxCountHolder}
            fi
            if [[ -n $__isOutOfListLimit ]]; then
              if [[ -z $__got ]]; then
                local __ending=; [[ $__count -gt 1 ]] && __ending=ы
                __got="получен$__ending ${_ansiSecondaryLiteral}$__count${_ansiErr}: ${_ansiCmd}$(_quotedArgs "${__list[@]}")"
              fi
              [[ -z $__expects ]] && __expects="ожидает $__expectsLimitKind ${_ansiPrimaryLiteral}$__expectsLimitCount${_ansiErr} $(_getPluralWord $__expectsLimitCount элемента элементов)"
              local __infix=
              local optName; dstVarName=optName _upperCamelCaseToKebabCase "$__varName"
              eval "$__codeIsOptVarName" \
                && __infix="опции ${_ansiCmd}--$optName" \
                || __infix="аргумента ${_ansiOutline}$__varName"
              __err="$__expects в списке значений $__infix${_ansiErr}, но $__got"
              returnCode=1
            fi
          fi
        fi
      fi
    fi
    [[ $returnCode -eq 0 ]] || break
  done
  return $returnCode
}

_validateDefaultCalculatedValue() {
  local __varValue="$1"
  _validateVarValue "$__varValue" && return 0
  local paramHolder="__PARAM_$__varName"
  [[ -n $infix3 ]] && infix3=" ${_ansiSecondaryLiteral}$infix3${_ansiErr}"
  __err="$__err в качестве${__infix}${__infix2} значения по умолчанию$infix3 параметра ${_ansiCmd}${!paramHolder:-$param}"
  return 1
}

_prepareList() {
  if [[ -z $__isReadyList ]]; then
    eval __list=\( \"\${$__varName[@]}\" \)
    __count=${#__list[@]}
    __isReadyList=true
  fi
}

_codeToPregen='
  __onlyPrepareCode=true $funcName; local __returnCode=$?
  if [[ $__returnCode -eq 0 ]]; then
    declare -f $funcName >&2
    return $(_err --showStack 2 "Похоже на то, что в ${_ansiCmd}$funcName${_ansiErr} пропущено ${_ansiCmd}eval \"\$_func${needPrepare}\"")
  elif [[ $__returnCode -ne 2 ]]; then
    # _debugVar --mark WARNING funcName needPrepare __returnCode
    return $__returnCode
  fi
'

_codeToPrepareVarNameByOptionShortcut='
  [[ $char != \? ]] || char=h
  local __varNameHolder=__VAR_NAME_FOR_SHORTCUT_$char
  __varName=${!__varNameHolder}
'

_validateVarValue() {
  local __varValue="$1"
  local returnCode=0
  __err=
  __errSuffixCode=
  if [[ -z $__varValue ]]; then
    local emptyHolder="__EMPTY_$__varName"
    if [[ -z ${!emptyHolder} ]]; then
      __err="ожидает ${_ansiOutline}непустое${_ansiErr} значение"
      __errSuffixCode=0
      returnCode=1
    fi
  else
    local __varValtypeHolder="__VALTYPE_$__varName"
    if [[ ${!__varValtypeHolder} == enum ]]; then
      local __codeIsValidEnumValueHolder=__codeIsValidEnumValue_$__varName
      local __codeIsValidEnumValue="${!__codeIsValidEnumValueHolder}"
      [[ -n $__codeIsValidEnumValue ]] \
        || __codeIsValidEnumValue='eval "$_codeToPrepareEnumValues" && _hasItem "$__varValue" "${enumValues[@]}"'
      eval "$__codeIsValidEnumValue" || {
        eval "$_codeToPrepareEnumValues"
        __err="одно из следующих значений: ${_ansiSecondaryLiteral}$(_quotedArgs "${enumValues[@]}" )"
        __errSuffixCode=1
        returnCode=1
      }
    elif [[ ${!__varValtypeHolder} == int ]]; then
      local intNumber="целое число"
      if [[ ! $__varValue =~ ^-?[[:digit:]]+$ ]]; then
        __err="$intNumber"
      else
        local intMinHolder="__MIN_$__varName"
        local intMaxHolder="__MAX_$__varName"
        if [[ ${!intMinHolder} != Infinite || ${!intMaxHolder} != Infinite ]]; then
          if [[ ${!intMaxHolder} == Infinite ]]; then
            if [[ ! $__varValue -ge ${!intMinHolder} ]]; then
              if [[ ${!intMinHolder} -eq 1 ]]; then
                __err="положительное $intNumber"
              elif [[ ${!intMinHolder} -eq 0 ]]; then
                __err="неотрицательное $intNumber"
              else
                __err="$intNumber${_ansiErr} не менее ${_ansiPrimaryLiteral}${!intMinHolder}"
              fi
            fi
          elif [[ ${!intMinHolder} == Infinite ]]; then
            if [[ ! $__varValue -le ${!intMaxHolder} ]]; then
              if [[ ${!intMaxHolder} -eq -1 ]]; then
                __err="отрицательное $intNumber"
              elif [[ ${!intMaxHolder} -eq 0 ]]; then
                __err="неположительное $intNumber"
              else
                __err="$intNumber${_ansiErr} не более ${_ansiPrimaryLiteral}${!intMaxHolder}"
              fi
            fi
          elif [[ $__varValue -lt ${!intMinHolder} || $__varValue -gt ${!intMaxHolder} ]]; then
            __err="$intNumber${_ansiErr} из диапазона ${_ansiPrimaryLiteral}${!intMinHolder}..${!intMaxHolder}"
          fi
        fi
      fi
      [[ -z $__err ]] || __err="${_ansiOutline}$__err"
    fi
    if [[ -n $__err ]]; then
      __err="вместо ${_ansiPrimaryLiteral}$__varValue${_ansiErr} ожидает $__err${_ansiErr}"
      __errSuffixCode=1
      returnCode=1
    fi
  fi
  return $returnCode
}

# =============================================================================

_descriptionSuffix=_description
_nameSuffix=_name
_cmdNameSuffix=_cmd$_nameSuffix
_moreArgsSuffix=MoreArgs
_moreArgsUsageSuffix=$_moreArgsSuffix'Usage'
_prepareCodeOfAutoHelp2() {
  local __thisFuncName="${FUNCNAME[2]}" __thisOnlyPrepareCode="$__onlyPrepareCode" __onlyPrepareCode=

  local __isCommandWrapperHolder=${__ownPrefix}isCommandWrapper
  local __canBeMoreParamsHolder=${__ownPrefix}canBeMoreParams

  dstVarName=__helpCodeFileSpec codeType=help additionalSuffix=$__additionalSuffix fileSpec= originalCodeDeep=2 eval "$_codeToPrepareCodeFileSpec"

  # [[ $__thisFuncName == bw_update ]] && _debugVar __optVarNames

  [[ -n ${!__isCommandWrapperHolder} ]] \
    || [[ $__thisFuncName =~ [[:alnum:]]_[_[:alnum:]] ]] \
    || [[ ${#__inheritedOptVarNames[@]} -gt 0 ]] \
    || [[ ! -f $__helpCodeFileSpec ]] \
    || ( [[ -n $__thisOnlyPrepareCode || -n $_isBwDevelop || -n $_isBwDevelopInherited ]] \
        && ! _everyFileNotNewerThan "$__helpCodeFileSpec" "${BASH_SOURCE[0]}" "${BASH_SOURCE[2]}" "${additionalDependencies[@]}" \
        ) \
    || return 0

  local holder="_codeToPrepareDescriptionsOf_${__thisFuncName}"
  if [[ -n ${!holder} ]]; then
    codeHolder=$holder eval "$_evalCode" || return $?
  fi
  # [[ $__thisFuncName == bw_update ]] && _debugVar --mark GEN __optVarNames needPregenHelp funcName __hasWrapper

  local __optionsDescription= __optionsTitle=; _prepareOptionsDescription
  local __usage='${_ansiCmd}$cmdName${_ansiReset}'
  [[ -n $__optionsTitle ]] && __usage+=" [$__optionsTitle]"

  local __argsDescription=

  if [[ -n ${!__isCommandWrapperHolder} ]]; then

    __indentLevel=1; eval "$_codeToPrepareHelpLinePrefixTemplate"
    local __subCmdNameHolder="$__thisFuncName$_cmdNameSuffix"
    local __subCmdName="${!__subCmdNameHolder:-Команда}"
    __usage+=' ${_ansiOutline}'$__subCmdName'${_ansiReset}'
    local -a allCommandNames=()
    local -a __funcNameSuffixes=(); help= _prepareSubCommandFuncSuffixes --checkCondition $__thisFuncName
    local __funcNameSuffix; for __funcNameSuffix in "${__funcNameSuffixes[@]}"; do
      local __commandName; dstVarName=__commandName _upperCamelCaseToKebabCase $__funcNameSuffix
      __argsDescription+=$__helpCodeLinePrefix'${_ansiCmd}$cmdName '$__commandName'${_ansiReset}'$_helpCodeLineSuffix
      dstVarName=__commandShortcuts srcVarName=${__thisFuncName}_${__funcNameSuffix}Shortcuts codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      local __commandShortcut; for __commandShortcut in "${__commandShortcuts[@]}"; do
        __argsDescription+=$__helpCodeLinePrefix'${_ansiCmd}$cmdName '$__commandShortcut'${_ansiReset}'$_helpCodeLineSuffix
      done

      local __descriptionOutput=
      \
        __indentLevel=2 \
        alwaysMultiline=--alwaysMultiline \
        descriptionHolder="${__thisFuncName}_$__funcNameSuffix$_descriptionSuffix" \
        descriptionOf=команды \
        _prepareDescriptionOutput2
      __argsDescription+="$__descriptionOutput"

      allCommandNames+=( "$__commandName" "${__commandShortcuts[@]}" )
    done
    __indentLevel=0; eval "$_codeToPrepareHelpLinePrefixTemplate"
    __argsDescription+=$__helpCodeLinePrefix'${_ansiHeader}Подробнее см.${_ansiReset}'$_helpCodeLineSuffix
    __indentLevel=1; eval "$_codeToPrepareHelpLinePrefixTemplate"
    __argsDescription+=$__helpCodeLinePrefix'${_ansiCmd}$cmdName ${_ansiOutline}'$__subCmdName' ${_ansiCmd}--help${_ansiReset}'$_helpCodeLineSuffix
    local __i; for (( __i=0; __i < ${#__SHORTCUTS_help}; __i++ )); do
      local __shortcut="${__SHORTCUTS_help:$__i:1}"
      __argsDescription+=$__helpCodeLinePrefix'${_ansiCmd}$cmdName ${_ansiOutline}'$__subCmdName' ${_ansiCmd}-'"$__shortcut"'${_ansiReset}'$_helpCodeLineSuffix
    done

  else
    __indentLevel=1

    local __minOptionalArgIdx=${#__argVarNames[@]}
    while [[ $__minOptionalArgIdx -gt 0 ]]; do
      local __idx=$(( $__minOptionalArgIdx - 1 ))
      local __varName=${__argVarNames[$__idx]}
      local __optionalHolder=__OPTIONAL_$__varName
      [[ -n ${!__optionalHolder} ]] || break
      __minOptionalArgIdx=$__idx
    done
    local __i; for (( __i=0; __i < ${#__argVarNames[@]}; __i++ )); do
      local __varName=${__argVarNames[$__i]}
      local __argNameHolder="${__thisFuncName}_$__varName$_nameSuffix"
      local __argName="${!__argNameHolder:-$__varName}"
      local __varTypeHolder=__TYPE_$__varName

      local __cmdArg='${_ansiOutline}"'\'$(_quotedArgs "$__argName")\''"${_ansiReset}'
      [[ ${!__varTypeHolder} == list ]] && __cmdArg+='...'

      __usage+=" "
      [[ $__i -ge $__minOptionalArgIdx ]] && __usage+='['
      __usage+="$__cmdArg"

      local __paramDescription=; _prepareParamDescription2; __argsDescription+="$__paramDescription"
    done
    for (( __i=__minOptionalArgIdx; __i < ${#__argVarNames[@]}; __i++ )); do
      __usage+=']'
    done
  fi
  if [[ -z ${!__isCommandWrapperHolder} && -n ${!__canBeMoreParamsHolder} ]]; then
    dstVarName=__moreArgsUsage srcVarName=$__thisFuncName$_moreArgsUsageSuffix eval "$_codeToInitLocalCopyOfScalarWithCheck"
    [[ -z $__moreArgsUsage ]] || __usage+=" $(_quotedArgs --strip "$__moreArgsUsage")"
    local __moreArgsHolder="$__thisFuncName$_moreArgsSuffix"
    dstVarName=__moreArgs srcVarName=$__moreArgsHolder eval "$_codeToInitLocalCopyOfArrWithCheck"
    __moreArgs=( "${__moreArgs[@]}" )
    local __i; for (( __i=0; __i < ${#__moreArgs[@]}; __i++ )); do
      local __varName="${__moreArgs[$__i]}"
      [[ $__varName =~ $_isValidVarNameRegExp ]] \
        || return $(_ownThrow "ожидает, что имя переменной ${_ansiPrimaryLiteral}$__varName${_ansiErr} указанное в ${_ansiOutline}$__moreArgsHolder${_ansiErr} $_mustBeValidVarName")
      ! _hasItem "$__varName" "${__argVarNames[@]}" \
        || return $(_ownThrow "${_ansiOutline}$__moreArgsHolder${_ansiErr} содержит аргумент ${_ansiPrimaryLiteral}$__varName${_ansiErr} одноименный параметру из ${_ansiOutline}${__thisFuncName}Params")
      local __argNameHolder="${__thisFuncName}_$__varName$_nameSuffix"
      local __argName="${!__argNameHolder:-$__varName}"
      local __varTypeHolder=__TYPE_$__varName

      local __cmdArg='${_ansiOutline}"'\'$(_quotedArgs "$__argName")\''"${_ansiReset}'

      [[ -n $__moreArgsUsage ]] || __usage+=" $__cmdArg"
      local __KIND_$__varName=arg

      local __paramDescription=; _prepareParamDescription2; __argsDescription+="$__paramDescription"
    done
  fi

  local __result

  __indentLevel=0; eval "$_codeToPrepareHelpLinePrefixTemplate"
  __result+=$__helpCodeLinePrefix'${_ansiHeader}Использование:${_ansiReset} '$__usage$_helpCodeLineSuffix
  __result+=$__helpCodeLinePrefix'${_ansiHeader}Описание:${_ansiReset}'

  local __descriptionOutput=
  \
    __indentLevel=1 \
    alwaysMultiline= \
    descriptionHolder="$__thisFuncName$_descriptionSuffix" \
    descriptionOf=- \
    _prepareDescriptionOutput2
  __result+="$__descriptionOutput"
  __result+="$_nl"

  [[ -z ${!__isCommandWrapperHolder} ]] \
    || __result+=$__helpCodeLinePrefix'${_ansiOutline}'$__subCmdName'${_ansiReset} - один из следующих вариантов: ${_ansiSecondaryLiteral}'${allCommandNames[@]}'${_ansiReset}'$_helpCodeLineSuffix

  __result+="$__argsDescription"

  if [[ -n $__optionsDescription ]]; then
    __result+=$__helpCodeLinePrefix$__optionsTitle$_helpCodeLineSuffix
    __result+="$__optionsDescription"
  fi

  _assureDir $(dirname "$__helpCodeFileSpec") || return $?
  echo "$__result" > "$__helpCodeFileSpec"
}

_prepareOptionsDescription() {
  local optionsCount=${#__optVarNames[@]}
  if [[ $optionsCount -gt 0 ]]; then
    __optionsTitle='${_ansiOutline}'$(_getPluralWord $optionsCount Опция Опции)'${_ansiReset}'
    __usage+=" [$__optionsTitle]"
  fi
  local -a proceedVarNames=()
  local __varName; for __varName in ${__optVarNames[@]}; do
    [[ $__varName != help ]] || continue
    ! _hasItem $__varName ${proceedVarNames[@]} || continue
    proceedVarNames+=( $__varName )
    local __paramDescription=; _prepareParamDescription2; __optionsDescription+="$__paramDescription"
  done
  __varName=help
  local __paramDescription=; _prepareParamDescription2; __optionsDescription+="$__paramDescription"
}

_codeToPrepareHelpLinePrefixTemplate='
  local __helpCodeLinePrefix="$_nl"
  [[ $__indentLevel -gt 0 ]] \
    && __helpCodeLinePrefix+="$(_indent $((2 * $__indentLevel)) )_indent $((4 * $__indentLevel)); "
  __helpCodeLinePrefix+="echo \""
'
_helpCodeLineSuffix='"'
_prepareParamDescription2() {
  local paramDescriptionHolder="${__thisFuncName}_$__varName$_descriptionSuffix"
  local importantHolder=__IMPORTANT_$__varName
  [[ -z ${!paramDescriptionHolder} && -n ${!importantHolder} ]] \
    && paramDescriptionHolder="_$__varName$_descriptionSuffix"
  local __varKindHolder=__KIND_$__varName
  local kindTitle
  local __varTypeHolder=__TYPE_$__varName
  local __varValtypeHolder=__VALTYPE_$__varName

  if [[ ${!__varValtypeHolder} == enum ]]; then
    local __varEnumHolder="__ENUM_$__varName"
    __paramDescription+="$_nl"'local -a enumValues='${!__varEnumHolder}
  fi
  if [[ ${!__varKindHolder} == arg ]]; then
    __indentLevel=0; eval "$_codeToPrepareHelpLinePrefixTemplate"
    __paramDescription+=$__helpCodeLinePrefix$__cmdArg
    kindTitle=аргумента
    alwaysMultiline=
    __indentLevel=1
  elif [[ ${!__varKindHolder} == opt ]]; then
    kindTitle=опции
    alwaysMultiline=--alwaysMultiline
    __indentLevel=2
  fi
  eval "$_codeToPrepareHelpLinePrefixTemplate"

  if [[ ${!__varValtypeHolder} == enum ]]; then
    if ! [[ ${!__varEnumHolder} =~ $_isCalcEnumRegExp ]]; then
      eval "$_codeToPrepareEnumValues"
      local -a descriptionHolders=( )
      local hasEnumValueDescription=
      local enumValue; for enumValue in "${enumValues[@]}"; do
        local normalizedEnumValue="${enumValue//[-]/_}"
        normalizedEnumValue="${normalizedEnumValue//[^[:alpha:][:digit:]]/}"
        descriptionHolder="${__thisFuncName}_${__varName}_$normalizedEnumValue$_descriptionSuffix"
        [[ -z ${!descriptionHolder} &&  -n ${!importantHolder} ]] \
          && descriptionHolder="_${__varName}_$normalizedEnumValue$_descriptionSuffix"
        descriptionHolders+=( "$descriptionHolder" )
        [[ -n ${!descriptionHolder} ]] || continue
        hasEnumValueDescription=true
      done
    fi
  fi

  local descriptionOf=$kindTitle
  if [[
    ( ${!__varKindHolder} == arg && -n ${!__argNameHolder} ) || \
    ( ${!__varValtypeHolder} == enum && -n $hasEnumValueDescription ) \
  ]]; then
    descriptionOf=
  fi

  local paramDescriptionOutput=

  local descriptionPrefix= upperFirstOfPrefixWhenMultiline=
  if [[ ${!__varKindHolder} == arg && ${!__varTypeHolder} == list ]]; then
    local listDescription=; _prepareListDescription2
    descriptionPrefix+="$listDescription"
    upperFirstOfPrefixWhenMultiline=--upperFirstOfPrefixWhenMultiline
  fi

  local __descriptionOutput;
    \
      descriptionHolder="$paramDescriptionHolder"
    _prepareDescriptionOutput2
  local paramDescriptionOutput+="$__descriptionOutput"
  descriptionPrefix= upperFirstOfPrefixWhenMultiline=

  local valueDescription=
    if [[ ${!__varValtypeHolder} == int ]]; then
    _prepareIntValudDescription2
  elif [[ ${!__varValtypeHolder} == enum ]]; then
    valueDescription='арианты ${_ansiOutline}значения${_ansiReset}:'
    [[ -z $alwaysMultiline ]] \
      && valueDescription=' в'$valueDescription \
      || valueDescription=$__helpCodeLinePrefix'В'$valueDescription
    if [[ -z $hasEnumValueDescription ]]; then
      valueDescription+=' ${_ansiSecondaryLiteral}$(_quotedArgs "${enumValues[@]}")${_ansiReset}"'
    else
      valueDescription+='"'
      local __i; for (( __i=0; __i < ${#enumValues[@]}; __i++ )); do
        local enumValue="${enumValues[$__i]}"
        valueDescription+="$_nl$(_indent $((2 * $__indentLevel)) )"'_indent '$(( $__indentLevel * 4 + 2 ))'; echo "${_ansiPrimaryLiteral}"'\'$(_quotedArgs "$enumValue")\''"${_ansiReset}'
        # __indentLevel=3; eval "$_codeToPrepareHelpLinePrefixTemplate"
        __indentLevel=$(( __indentLevel + 1 )); eval "$_codeToPrepareHelpLinePrefixTemplate"
        alwaysMultiline=
        local __descriptionOutput=
        \
          descriptionHolder="${descriptionHolders[$__i]}" \
          descriptionOf='значения опции' \
        _prepareDescriptionOutput2
        valueDescription+="$__descriptionOutput"
        # __indentLevel=2; eval "$_codeToPrepareHelpLinePrefixTemplate"
        __indentLevel=$(( __indentLevel - 1 )); eval "$_codeToPrepareHelpLinePrefixTemplate"
      done
    fi
    alwaysMultiline=--alwaysMultiline
  fi

  local emptyDescription=
  if [[ ${!__varTypeHolder} != bool ]]; then
    local emptyHolder=__EMPTY_$__varName
    if [[ -n ${!emptyHolder} ]]; then
      if [[ -n $alwaysMultiline ]]; then
        emptyDescription=$__helpCodeLinePrefix'${_ansiOutline}Значение${_ansiReset} может быть пустым'$_helpCodeLineSuffix
      else
        [[ -z $valueDescription ]] || emptyDescription+=', '
        emptyDescription+='может быть пустым'
      fi
    fi
  fi

  if [[ ${!__varKindHolder} == opt ]]; then
    local __indentLevel=1; eval "$_codeToPrepareHelpLinePrefixTemplate"
    local optionValueDeclaration=
    local __varTypeHolder=__TYPE_$__varName
    [[ ${!__varTypeHolder} == bool ]] \
      || optionValueDeclaration='${_ansiReset} ${_ansiOutline}значение'
    local __optName; dstVarName=__optName _upperCamelCaseToKebabCase $__varName
    __paramDescription+=$__helpCodeLinePrefix'${_ansiCmd}--'${__optName}$optionValueDeclaration'${_ansiReset}'

    dstVarName=__shortcuts srcVarName=__SHORTCUTS_$__varName eval "$_codeToInitLocalCopyOfScalar"
    local __i; for (( __i=0; __i < ${#__shortcuts}; __i++ )); do
      local __shortcut=${__shortcuts:$__i:1}
      __paramDescription+=' или ${_ansiCmd}-'"$__shortcut"$optionValueDeclaration'${_ansiReset}'
    done
    __paramDescription+='"'
    local __indentLevel=2
  fi
  __paramDescription+="$paramDescriptionOutput"
  __paramDescription+="$valueDescription"
  __paramDescription+="$emptyDescription"
  [[ -n $alwaysMultiline ]] || { __paramDescription+='"'; alwaysMultiline=--alwaysMultiline; }

  local __defaultValueHolder=__DEFAULT_$__varName
  eval "$_codeToPrepareHelpLinePrefixTemplate"

  if [[ ${!__varTypeHolder} == list && ${!__varKindHolder} == opt ]]; then
    local listDescription=; _prepareListDescription2
    __paramDescription+=$__helpCodeLinePrefix'Опция предназначена для того, чтобы сформировать'$_helpCodeLineSuffix
    __paramDescription+=$__helpCodeLinePrefix"$listDescription"$_helpCodeLineSuffix
    __paramDescription+=$__helpCodeLinePrefix'путем eё многократного использования'$_helpCodeLineSuffix
  fi

  if [[ -n ${!__defaultValueHolder} ]]; then
    local prefix=
    [[ ${!__varTypeHolder} == list ]] \
       && prefix='Значение ${_ansiOutline}списка' \
       || prefix='${_ansiOutline}Значение'
    local ansi=
    [[ ${!__defaultValueHolder} =~ \$ ]] \
       && ansi=_ansiOutline \
       || ansi=_ansiPrimaryLiteral
    __paramDescription+=$__helpCodeLinePrefix$prefix'${_ansiReset} по умолчанию: ${'$ansi'}"'\'"${!__defaultValueHolder}"\''"${_ansiReset}'$_helpCodeLineSuffix
  fi
}

_prepareDescriptionOutput2() {
  local quote=
  local description=
  if [[ -n ${!descriptionHolder} || -n $descriptionPrefix ]]; then
    [[ -z ${!descriptionHolder} ]] || description+="${!descriptionHolder}"
  elif [[ -n $descriptionOf ]]; then
    description+='${_ansiErr}Нет описания ${_ansiOutline}'"$descriptionHolder"
    [[ $descriptionOf != - ]] \
      && description+='${_ansiErr} '"$descriptionOf"
    description+='${_ansiReset}'
    quote='--quote double'
  fi
  __descriptionOutput=
  local __returnCode=0
  if [[ -n $description || -n $descriptionPrefix ]]; then

    if [[ -n $descriptionPrefix && $description =~ ^! ]]; then
      if [[ -n $alwaysMultiline ]]; then
        __descriptionOutput+=$__helpCodeLinePrefix$descriptionPrefix$_helpCodeLineSuffix
      else
        [[ $descriptionPrefix =~ ^[[:space:]] ]] || descriptionPrefix=" $descriptionPrefix"
        __descriptionOutput+=$descriptionPrefix$_helpCodeLineSuffix
        alwaysMultiline=--alwaysMultiline
      fi
      description=${description:1}
      descriptionPrefix=
    fi

    _prepareCodeForDescriptionOutput2 $quote $alwaysMultiline $indentBase $upperFirstOfPrefixWhenMultiline $__indentLevel "$description" "$descriptionPrefix"; local __returnCode=$?;
    if [[ $__returnCode -eq 0 ]]; then
      if [[ -z $alwaysMultiline ]]; then
        __descriptionOutput+=$_helpCodeLineSuffix
        alwaysMultiline=--alwaysMultiline
      fi
      __descriptionOutput+="$codeForDescriptionOutput"
    elif [[ $__returnCode -eq 2 ]]; then
      if [[ -n $codeForDescriptionOutput ]]; then
        __descriptionOutput=" $codeForDescriptionOutput"$_helpCodeLineSuffix
        alwaysMultiline=--alwaysMultiline
      fi
      __returnCode=0
    fi
  fi
  return $__returnCode
}

_prepareIntValudDescription2() {
  _prepareIntValudDescription2Helper
  local prefix
  [[ -z $alwaysMultiline ]] \
    && prefix= \
    || prefix="$_nl$(_indent $((2 * $__indentLevel)) )"'_indent '$(( 4 * $__indentLevel ))'; echo "${_ansiOutline}Значение${_ansiReset}'
  valueDescription=$prefix' - '$valueDescription
  alwaysMultiline=
}

_prepareIntValudDescription2Helper() {
  valueDescription="целое число";
  local intMinHolder="__MIN_$__varName"
  local intMaxHolder="__MAX_$__varName"
  if [[ ${!intMinHolder} != Infinite && ${!intMaxHolder} != Infinite ]]; then
    valueDescription+=' из диапазона ${_ansiSecondaryLiteral}'${!intMinHolder}'..'${!intMaxHolder}'${_ansiReset}';
  elif [[ ${!intMinHolder} != Infinite && ${!intMaxHolder} == Infinite ]]; then
    if [[ ${!intMinHolder} -eq 0 ]]; then
      valueDescription="неотрицательное $valueDescription";
    elif [[ ${!intMinHolder} -eq 1 ]]; then
      valueDescription="положительное $valueDescription";
    else
      valueDescription+=' не менее ${_ansiPrimaryLiteral}'${!intMinHolder}'${_ansiReset}';
    fi
  elif [[ ${!intMinHolder} == Infinite && ${!intMaxHolder} != Infinite ]]; then
    if [[ ${!intMaxHolder} -eq 0 ]]; then
      valueDescription="неположительное $valueDescription";
    elif [[ ${!intMaxHolder} -eq -1 ]]; then
      valueDescription="отрицательное $valueDescription";
    else
      valueDescription+=' не более ${_ansiPrimaryLiteral}'${!intMaxHolder}'${_ansiReset}';
    fi
  fi
}

_prepareListDescription2() {
  local __minCountHolder=__MIN_COUNT_$__varName
  local __maxCountHolder=__MAX_COUNT_$__varName
  if [[ ${!__minCountHolder} -eq 0 ]]; then
    listDescription="возможно пустой "
  elif [[ ${!__minCountHolder} -eq 1 ]]; then
    listDescription="непустой "
  fi
  listDescription+=список;
  [[ ${!__maxCountHolder} != Infinite ]] && listDescription+=" ("
  if [[ ${!__minCountHolder} -gt 1 ]]; then
    [[ ${!__maxCountHolder} == Infinite ]] \
      && listDescription+=" ("
    listDescription+="не менее ${!__minCountHolder} "
    [[ ${!__maxCountHolder} != Infinite ]] \
      && listDescription+="и " \
      || listDescription+="$(_getPluralWord ${!__minCountHolder} элемента элементов))"
  fi
  [[ ${!__maxCountHolder} != Infinite ]] \
    && listDescription+="не более ${!__maxCountHolder} $(_getPluralWord ${!__maxCountHolder} элемента элементов))"
  local __uniqueHolder=__UNIQUE_$__varName
  [[ -n ${!__uniqueHolder} && ( ${!__maxCountHolder} == Infinite || ${!__maxCountHolder} -gt 1 ) ]] \
    && listDescription+=" ${_ansiUnderline}уникальных${_ansiReset}"
  listDescription+=" ${_ansiOutline}значений${_ansiReset}"
}

_prepareCodeForDescriptionOutput2Params=(
  '--alwaysMultiline'
  '--indentBase:1..=4'
  '--upperFirstOfPrefixWhenMultiline'
  '--quote:(double single)=single'
  'indentLevel:0..'
  'description:?'
  'descriptionPrefix:?'
)
_prepareCodeForDescriptionOutput2() { eval "$_funcParams2"
  [[ -z $descriptionPrefix || ! $descriptionPrefix =~ "$_nl" ]] \
    || return $(_ownThrow "ожидает ${_ansiOutline}descriptionPrefix${_ansiErr} (${_ansiPrimaryLiteral}${!descriptionPrefix}${_ansiErr}) без перевода строки")
  [[ -n $description || -n $descriptionPrefix ]] || return 0
  if [[ -z $alwaysMultiline && ! $description =~ "$_nl" ]]; then
    if [[ $quote == single ]]; then
      if [[ -n $descriptionPrefix && -n $description ]]; then
        codeForDescriptionOutput='"'\'$descriptionPrefix' '$description\''"'
      elif [[ -n $descriptionPrefix ]]; then
        codeForDescriptionOutput='"'\'$descriptionPrefix''\''"'
      else
        codeForDescriptionOutput='"'\'$description\''"'
      fi
    else
      if [[ -n $descriptionPrefix && -n $description ]]; then
        codeForDescriptionOutput="$descriptionPrefix $description"
      elif [[ -n $descriptionPrefix ]]; then
        codeForDescriptionOutput="$descriptionPrefix"
      else
        codeForDescriptionOutput="$description"
      fi
    fi
    return 2
  fi
  local ucfirst=
  [[ -z $upperFirstOfPrefixWhenMultiline ]] || ucfirst=ucfirst
  local indentBaseOPT=
    [[ $indentBase -ne 4 ]] && indentBaseOPT=' '"--base $indentBase"
  eval "$_codeToPrepareHelpLinePrefixTemplate"
  [[ -z $descriptionPrefix || $descriptionPrefix =~ [[:space:]]$ ]] || descriptionPrefix+=' '
  local perlCode='
    use POSIX;
    use utf8;
    use List::Util qw[max];
    $baseIndentLevel=0;
    sub initFirstIndent {
      my $spaces=shift;
      $baseIndentLevel = floor(length($spaces) / 2);
      return '$ucfirst' "'$descriptionPrefix'";
    }
    sub replaceTo {
      my $spaces=shift;
      $indentLevel='$indentLevel' + max(0, floor(length($spaces) / 2) - $baseIndentLevel);
      $indent="  " x $indentLevel;
      "\n${indent}_indent \$(('$indentBase' * $indentLevel)); echo \"";
    }
    $/=undef;
    $_=<>;
    s/^(?:[ ]*\n)*(\s*)/initFirstIndent($1)/e;
    s/\s+$//;
    $exitCode=0;'
  [[ -z $alwaysMultiline ]] && perlCode+='
    if ( !/\n/ ) {
      s/^/ /;
      $exitCode=2;
    } else {'
  perlCode+='
      s/\\/\\\\/g;
      s/"/\\"/g;
      s/\n(\s*)/";@{[replaceTo($1)]}/g;
      s/^/@{[replaceTo()]}/;
      s/$/";/;'
  [[ -z $alwaysMultiline ]] && perlCode+='
    }'
  perlCode+='
    print;
    exit($exitCode);
  '
  codeForDescriptionOutput=$(echo "${description}" | perl -C -e "$perlCode")
}

# =============================================================================

_codeToParseVarValueTypeInfo=(
  "перечислимый: ${_ansiSecondaryLiteral}:( ${_ansiOutline}значение1 значение 2 ... ${_ansiSecondaryLiteral})"
  "целочисленный диапазон: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}..${_ansiOutline}max"
  "целочисленный, не менее: ${_ansiSecondaryLiteral}:${_ansiOutline}min${_ansiSecondaryLiteral}.."
  "целочисленный, не более: ${_ansiSecondaryLiteral}:..${_ansiOutline}max"
  "целочисленный: ${_ansiSecondaryLiteral}:.."
  "фиксированное целое число: ${_ansiSecondaryLiteral}:${_ansiOutline}intValue"
)

_codeToParseVarListTypeInfo=(
  "список будет содержать только уникальные значения: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}unique"
  "список будет содержать отсортированные значения: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}sorted"
)

_codeToParseVarEmptyTypeInfo2=(
  "значение может быть пустым: ${_ansiSecondaryLiteral}:${_ansiPrimaryLiteral}?"
)

_varTypesInRegExp2='sorted|unique|\?'

_diapRegExp='[[:space:]]*(-?[[:digit:]]*)([[:space:]]*\.\.[[:space:]]*(-?[[:digit:]]*))?[[:space:]]*'

_listBordersRegExp='[[:space:]]*([[:digit:]]*)([[:space:]]*\.\.[[:space:]]*([[:digit:]]*))?[[:space:]]*'

# =============================================================================
