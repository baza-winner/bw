# =============================================================================

_resetBash

# =============================================================================

_unsetBash() {
  local verbosity="$1"
  local fileSpec="${fileSpec:-${BASH_SOURCE[1]}}"
  local unsetFileSpec; unsetFileSpec="$(_getUnsetFileSpecFor "$fileSpec")"
  local -a unsetFileSpecs=()
  local -a rmFileSpecs=()
  if [[ -f $unsetFileSpec ]]; then
    unsetFileSpecs+=( "$unsetFileSpec" )
    rmFileSpecs+=( "$unsetFileSpec" )
  fi
  local funcName; for funcName in $(_getFuncNamesOfScriptToUnset "$fileSpec"); do
    local completionCodeFileSpec
    dstVarName=completionCodeFileSpec codeType=completion additionalSuffix='' eval "$_codeToPrepareCodeFileSpec"
    if [[ -f $completionCodeFileSpec ]]; then
      rmFileSpecs+=( "$completionCodeFileSpec" )
    fi
    local unsetFileSpec="${completionCodeFileSpec:0:$(( ${#completionCodeFileSpec} - ${#_codeBashExt} ))}$_unsetFileExt"
    if [[ -f $unsetFileSpec ]]; then
      unsetFileSpecs+=( "$unsetFileSpec" )
      rmFileSpecs+=( "$unsetFileSpec" )
    fi
  done
  for unsetFileSpec in "${unsetFileSpecs[@]}"; do
    if [[ $verbosity == dry ]]; then
      echo "${_ansiCmd}. \"$unsetFileSpec\"${_ansiReset}"
    else
      . "$unsetFileSpec"
      if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
        echo "${_ansiOK}OK: ${_ansiCmd}. $unsetFileSpec${_ansiReset}"
      fi
    fi
  done
  if [[ ${#rmFileSpecs[@]} -gt 0 ]]; then
    if [[ $verbosity == dry ]]; then
      echo "${_ansiCmd}rm $(_quotedArgs "${rmFileSpecs[@]}")${_ansiReset}"
    else
      rm "${rmFileSpecs[@]}"
      if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
        echo "${_ansiOK}OK: ${_ansiCmd}rm $(_quotedArgs "${rmFileSpecs[@]}")${_ansiReset}"
      fi
    fi
  fi
}
# =============================================================================

_codeToInitSubOPT='
  local -a sub_OPT_silent=( "${OPT_silent[@]}" )
  local -a sub_OPT_verbosity=()
  [[ $verbosity =~ ^(none|dry|all)$ ]] \
    && local -a sub_OPT_verbosity=( "${OPT_verbosity[@]}" ) \
    || local -a sub_OPT_verbosity=( --verbosity err )
  local -a sub_OPT=( "${sub_OPT_silent[@]}" "${sub_OPT_verbosity[@]}" )
'
verbosityDefault=allBrief silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"

# =============================================================================

bwParams=(
  # "${_verbosityParams[@]}"
)
bwParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_cmd_name=Команда
bw_description='Базовая утилита проектов baza-winner'
bw() { eval "$_funcParams2"
}

# =============================================================================

bw_versionParams=()
bw_version_description='Выводит номер версии bw.bash'
bw_version() { eval "$_funcParams2"
  local suffix="$BW_SELF_UPDATE_SOURCE"
  [[ -z $suffix ]] || suffix="@$suffix"
  echo "$_bwVersion$suffix"
}

# =============================================================================

bw_updateParams=( '--remove/r' )
bw_update_remove_description="Удалить прегенеренные файлы перед обновлением"
bw_update_description="Обновляет bw"
bw_update() { eval "$_funcParams2"
  if [[ -n $remove ]]; then
    "$_bwFileSpec" rm -y
  fi
  . "$_bwFileSpec"
  echo "Current version: $(bw_version)"
}

# =============================================================================

bw_removeParams=(
  --yes/y --completely/c
  "${_verbosityParams[@]}"
)
bw_remove_yes_description='подтверждает удаление'
bw_remove_completely_description='удаляет не только все связанное с $_bwFileName, но и сам bw.bash'
bw_remove_silent_description='подавляет сообщение об успешном удалении'
bw_removeShortcuts=( 'rm' )
bw_remove_description="удаляет $_bwFileName и все связанное с ним"
bw_removeCondition='[[ -z $_isBwDevelopInherited ]]'
bw_remove() { eval "$_funcParams2"
  if [[ -z $yes ]]; then
    _warn "Чтобы произвести удаление, запустите эту команду с опцией ${_ansiCmd}--yes${_ansiWarn} или ${_ansiCmd}-y"
    return 2
  else
    codeHolder=_codeToInitSubOPT eval "$_evalCode"

    if [[ -z $_isBwDevelop ]]; then
      _profileUnless "${sub_OPT[@]}" -u ". $(_quotedArgs "$(_shortenFileSpec "$_bwFileSpec")")"
    fi

    local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
      [[ -d $dirSpec ]] || continue
      for _fileSpec in "$dirSpec/"*$_unsetFileExt; do
        if _hasItem $(basename "$_fileSpec" "$_unsetFileExt") "${_bw_removeRequiredFiles[@]}"; then
          _mvFile "${sub_OPT[@]}" "$_fileSpec" "/tmp/$(basename "$_fileSpec")"
        else
          codeHolder=_codeSource eval "$_evalCode"
        fi
      done
    done

    if [[ -z $_isBwDevelop ]]; then
      [[ ! -d $_bwDir ]] || _rm "${sub_OPT[@]}" -d "$_bwDir"
      [[ -z $completely ]] || [[ ! -f $_bwFileSpec ]] || _rm "${sub_OPT[@]}" "$_bwFileSpec"
    else
      _rm "${sub_OPT[@]}" "$_bwDir/old.bw.bash"
      _rm "${sub_OPT[@]}" -d "$_bwDir/tgz"
      _rm "${sub_OPT[@]}" -pd "$_bwDir/tmp"
    fi

    local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
      _rm "${sub_OPT[@]}" -d "$dirSpec"
    done

    if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
      if [[ -n $_isBwDevelop ]]; then
        echo "${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все прегенеренные вспомогательные файлы. Все основные ${_ansiFileSpec}*.bash${_ansiWarn}-файлы оставлены нетронутыми${_ansiReset}"
      else
        local suffix;
        if [[ -z $completely ]]; then
          suffix=", кроме ${_ansiFileSpec}$_bwFileSpec${_ansiWarn}. Для повторной установки выполните команду ${_ansiCmd}. $_bwFileSpec${_ansiWarn}"
        else
          suffix=", включая ${_ansiFileSpec}$_bwFileSpec${_ansiWarn}"
        fi
        echo "${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все связанное с ней (содержимое директории ${_ansiFileSpec}$_bwDir${_ansiWarn})$suffix${_ansiReset}"
      fi
    fi

    local -a varNamesToUnset=( $( compgen -v | perl -ne "print if /^(_specialVars|$_specialVars)\$/" ) )
    if [[ $verbosity == dry ]]; then
      for _fileSpec in "${_bw_removeRequiredFiles[@]}"; do
        _fileSpec="/tmp/$_fileSpec$_unsetFileExt"
        echo "${_ansiCmd}. \"$_fileSpec\"${_ansiReset}"
        codeHolder=_codeSource eval "$_evalCode"
      done
      echo "${_ansiCmd}unset ${varNamesToUnset[@]}${_ansiReset}"
    else
      for _fileSpec in "${_bw_removeRequiredFiles[@]}"; do
        _fileSpec="/tmp/$_fileSpec$_unsetFileExt"
        codeHolder=_codeSource eval "$_evalCode"
      done
      unset ${varNamesToUnset[@]}
    fi
  fi
}

# =============================================================================

_testPlurals='тест теста тестов'
_funcPlurals='функция функции функций'

bw_bashTests_list_description='
  Выводит список bash-функций, для которых существуют тесты,
  с указанием количества тестов для каждой функции
'
bw_bashTests_noTiming_description='Отключает замер времени работы тестов'
bw_bashTests_args_name=Аргумент
bw_bashTests_args_description='!
  ${_ansiOutline}Значение${_ansiReset} - ${_ansiOutline}Имя-Функции${_ansiReset} или ${_ansiOutline}Диапазон-Номеров${_ansiReset}
  ${_ansiOutline}Имя-Функции${_ansiReset} -- имя функции, для которой надо прогнать тест
  ${_ansiOutline}Диапазон-Номеров${_ansiReset} -- диапазон номеров тестов, которые нужно прогнать
    Значения ${_ansiOutline}Диапазона-Номеров${_ansiReset}
      ${_ansiSecondaryLiteral}0..2${_ansiReset} - 0, 1 и 2 тесты
      ${_ansiPrimaryLiteral}-2..-1${_ansiReset} - предпоследний и последний тест
      ${_ansiSecondaryLiteral}4..${_ansiReset} - все тесты, начиная с 4-го
      ${_ansiSecondaryLiteral}..3${_ansiReset} - все тесты до 3-го
      ${_ansiSecondaryLiteral}6${_ansiReset} - 6-ой тест
      ${_ansiSecondaryLiteral}-1${_ansiReset} - последний тест
  После каждого ${_ansiOutline}Имени-Функции${_ansiReset} могут следовать один или несколько ${_ansiOutline}Диапазонов-Номеров${_ansiReset}
  Если ${_ansiOutline}Диапазон-Номеров${_ansiReset} не задан, то прогоняет все тесты для функции ${_ansiOutline}Имя-Функции${_ansiReset}
  Если список ${_ansiOutline}Аргументов${_ansiReset} пуст, то прогоняет все существующие тесты
'
bw_bashTestsShortcuts=( 'bt' )
bw_bashTestsParams=( '--noTiming' '--list' '@..args' )
bw_bashTestsParamsOpt=( --canBeMixedOptionsAndArgs )
bw_bashTestsSuppressInheritedOptionVarNames=( 'verbosity' 'silent' )
bw_bashTests_description='запускает тесты bash-функций'
bw_bashTests() {
  eval "$_funcParams2"
  local testsDirSpec="$_bwDir/tests"
  if [[ -z $_isBwDevelop && -z $_isBwDevelopInherited ]]; then
    local testsSupportFileSpec="$testsDirSpec/testsSupport.bash"
    # if [[ ! -f $testsSupportFileSpec || $_bwFileSpec -nt $testsSupportFileSpec ]]; then
    _getBwTar "$_bwFileSpec" tests | tar xf - -C "$_bwDir" \
      || return $(_err "Не удалось извлечь архив tests из ${_ansiFileSpec}$_bwFileSpec${_ansiErr} в ${_ansiFileSpec}$testsDirSpec")
    # fi
    _fileSpec="$testsSupportFileSpec" codeHolder=_codeSource eval "$_evalCode"
  fi
  _runInBackground bw_bashTestsHelper
}
bw_bashTestsComplete() {
  if [[ $__varName == args ]]; then
    local -a funcsWithTests=(); _prepareFuncWithTests
    local __candidate; for __candidate in "${funcsWithTests[@]}"; do
      eval "$_codeToAddCandidateToCompletion"
    done
  fi
  return 0
}
_prepareFuncWithTests() {
  local testsDirSpec="$_bwDir/tests"
  for _fileSpec in "$testsDirSpec/"*Tests.bash; do
    codeHolder=_codeSource eval "$_evalCode"
  done
  local -a allFuncsWithTests=()
  _prepareAllFuncWithTests
  local funcTestFor; for funcTestFor in ${allFuncsWithTests[@]}; do
    local conditionHolder="${funcTestFor}TestsCondition"
    [[ -z ${!conditionHolder} ]] || eval "${!conditionHolder}" || continue
    funcsWithTests+=( "$funcTestFor" )
  done
}
_prepareAllFuncWithTests() {
  [[ -n $_isBwDevelop ]] && _rmCache \
    || varName=allFuncsWithTests codeHolder=_codeToUseCache eval "$_evalCode"
  _prepareAllFuncWithTestsHelper
  _saveToCache allFuncsWithTests
}
_prepareAllFuncWithTestsHelper() {
  _profileBegin
  allFuncsWithTests=()
  local testsVarNameSuffix=Tests
  local testsVarName; for testsVarName in $(compgen -v | grep -E $testsVarNameSuffix'$' | grep -v -E '^(allFuncsWith|funcsWith|self|_succeed|_failed)' ); do
    dstVarName=selfTests srcVarName=$testsVarName codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
    local count="${#selfTests[@]}"
    if [[ $count -gt 0 ]]; then
      local funcTestFor=${testsVarName:0:$(( ${#testsVarName} - ${#testsVarNameSuffix} ))}
      allFuncsWithTests+=( "$funcTestFor" )
    fi
  done
  _profileEnd
}

# =============================================================================

bw_setParams=(
  '!--uninstall/u'
)
bw_setParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_set_cmd_name=Имя-настройки
bw_set_description='включает/отключает настройку'
bw_set() { eval "$_funcParams2"
}

# =============================================================================


_enumAnsi=
_initEnumAnsi() {
  _enumAnsi=
  _enumAnsi+='('
  local _ansi; for _ansi in $(compgen -v _ansi); do
    _enumAnsi+=" ${_ansi:5}"
  done
  _enumAnsi+=' )'
}
_initEnumAnsi

_enumSpace='(before after both none -)'
_enumOnOff='(yes no on off -)'
_ps_user_defaultAnsi='Green'
_ps_ptcd_defaultAnsi='White'
_ps_cd_defaultAnsi='White'

_preparePromptParams=()
_preparePromptParams() {
  varName=_preparePromptParams codeHolder=_codeToUseCache eval "$_evalCode"
  local _psItems=( $( compgen -v | perl -ne "print if s/^_ps_([^_]+)\$/\$1/" ) )
  local _preparePromptParams_itemsDefaults=( user ptcd git error separator )
  local _preparePromptParams_itemsEnum=( ${_preparePromptParams_itemsDefaults[@]} )
  local _preparePromptParamsAddon=()
  local _psItem; for _psItem in ${_psItems[@]}; do
    local defaultAnsiHolder="_ps_${_psItem}_defaultAnsi"
    local defaultValueHolder="_ps_${_psItem}"
    local defaultAnsi="Reset"
    if [[ -n ${!defaultAnsiHolder} ]]; then
      defaultAnsi="${!defaultAnsiHolder}"
    fi
    _preparePromptParamsAddon+=(
      "--${_psItem}=$(_quotedArgs --strip ${!defaultValueHolder})"
      "--${_psItem}Space:${_enumSpace}=after"
      "@1..--${_psItem}Ansi:${_enumAnsi}=( $defaultAnsi )"
    )
    if ! _hasItem ${_psItem} "${_preparePromptParams_itemsEnum[@]}"; then
      _preparePromptParams_itemsEnum+=( ${_psItem} )
    fi
  done

  _preparePromptParams=(
    "@1..--items:( $( echo ${_preparePromptParams_itemsEnum[@]} ) )=( $( echo ${_preparePromptParams_itemsDefaults[@]} ) )"

    "${_preparePromptParamsAddon[@]}"

    "--git:$_enumOnOff=on"

    '--gitPrefix=\('
    "--gitPrefixSpace:$_enumSpace=none"
    "@1..--gitPrefixAnsi:$_enumAnsi=( DarkGrey )"

    "--gitBranch:$_enumOnOff=on"
    "--gitBranchSpace:$_enumSpace=none"
    "@1..--gitBranchAnsi:$_enumAnsi=( DarkGrey )"

    "--gitDirty:$_enumOnOff=on"
    "--gitDirtySpace:$_enumSpace=before"
    "@1..--gitDirtyAnsi:$_enumAnsi=( Yellow Bold )"

    '--gitSuffix=\)'
    "--gitSuffixSpace:$_enumSpace=after"
    "@1..--gitSuffixAnsi:$_enumAnsi=( DarkGrey )"

    "--error:$_enumOnOff=on"

    '--errorPrefix=\\\$?'
    "--errorPrefixSpace:$_enumSpace=none"
    "@1..--errorPrefixAnsi:$_enumAnsi=( DarkGrey )"

    '--errorInfix=='
    "--errorInfixSpace:$_enumSpace=none"
    "@1..--errorInfixAnsi:$_enumAnsi=( Reset )"

    "--errorCode:$_enumOnOff=on"
    "--errorCodeSpace:$_enumSpace=after"
    "@1..--errorCodeAnsi:$_enumAnsi=( Red Bold )"

    "--errorSuffix:$_enumOnOff=off"
    "--errorSuffixSpace:$_enumSpace=after"
    "@1..--errorSuffixAnsi:$_enumAnsi=( Reset )"

    '--separator=\\\$'
    "--separatorSpace:$_enumSpace=after"
    "@1..--separatorAnsi:$_enumAnsi=( White )"
  )
  _saveToCache '_preparePromptParams'
}
_preparePrompt() {
  eval "$_funcParams2"
  local userAnsiAsString=
  local optVarName; for optVarName in ${__optVarNames[@]}; do
    if [[ ${optVarName:$(( ${#optVarName} - 4 ))} == Ansi  ]]; then
      local item=${optVarName:0:$(( ${#optVarName} - 4 ))}
      dstVarName=ansiSrc srcVarName=${item}Ansi eval "$_codeToInitLocalCopyOfArray"
      local ansiAsStringHolder=${item}AnsiAsString
      local $ansiAsStringHolder=
      local ansi; for ansi in ${ansiSrc[@]}; do
        local ansiHolder="_ansi${ansi}"
        eval $ansiAsStringHolder+=\"\${!ansiHolder}\"
      done
      local ansiAsString=${!ansiAsStringHolder}
    fi
  done

  prompt=
  local -a groups=(error git)
  local -a realItems=()
  local item; for item in ${items[@]}; do
    local group; for group in ${groups[@]}; do
      if [[ ${!group} != off && ${!group} != - && $item == $group ]]; then
        eval local $group=true
      fi
    done
    if [[ $item == error ]]; then
      realItems+=( errorPrefix errorInfix errorCode errorSuffix )
    elif [[ $item == git ]]; then
      realItems+=( gitPrefix gitBranch gitDirty gitSuffix )
    else
      realItems+=( $item )
    fi
  done
  local group; for group in ${groups[@]}; do
    if [[ ${!group} == true ]]; then
      prompt+='`_psPrepare_'$group'`'
    fi
  done
  local item; for item in ${realItems[@]}; do
    if ! [[ ${!item} == - || ${!item} == off || ${!item} == no ]]; then
      local spaceHolder="${item}Space"
      local foundGroup= group; for group in ${groups[@]}; do
        if [[ ${item:0:${#group}} == $group ]]; then
          foundGroup=$group
          break
        fi
      done
      local psFuncName="_ps_$item"
      local promptItem=
      if _funcExists $psFuncName; then
        if [[ -z $foundGroup || ${!foundGroup} == true ]]; then
          promptItem='`'$psFuncName' '${!spaceHolder}'`'
        fi
      elif [[ -z $foundGroup ]]; then
        local output="${!item}"
        if [[ ${!spaceHolder} == both ]]; then
          promptItem=" $output "
        elif [[ ${!spaceHolder} == before ]]; then
          promptItem=" $output"
        elif [[ ${!spaceHolder} == after ]]; then
          promptItem="$output "
        else
          promptItem="$output"
        fi
      elif [[ ${!foundGroup} == true ]]; then
        promptItem='`'_psIf_$group' '${!spaceHolder}' '$(_quotedArgs --quote:all "${!item}")'`'
      fi
      if [[ -n $promptItem ]]; then
        local ansiAsStringHolder=${item}AnsiAsString
        if [[ -z ${!ansiAsStringHolder} ]]; then
          prompt+="$promptItem"
        else
          prompt+="${_psColorSegmentBeginPrefix}${!ansiAsStringHolder}${_psColorSegmentBeginSuffix}$promptItem${_psColorSegmentEnd}"
        fi
      fi
    fi
  done
}

# =============================================================================

bw_set_promptParams=()
bw_set_promptParams() {
  _preparePromptParams
  bw_set_promptParams=(
    "${_preparePromptParams[@]}"
  )
}
bw_set_prompt_description='Настраивает prompt'
bw_set_prompt() { eval "$_funcParams2"
  local -a supportModules=( git ansi ps )
  local supportFileNameSuffix='Support.bash'
  local newFileSpec="$_profileFileSpec.new"

  if [[ -n $uninstall ]]; then
    local dstVarName=PS1
    local srcVarName=OLD_PS1
    local -a fileNames=()
    local moduleName; for moduleName in ${supportModules[@]}; do
      fileNames+=( "${moduleName}${supportFileNameSuffix}")
    done
    local -a fileSpecs=()
    local fileName; for fileName in "${fileNames[@]}"; do
      local fileSpec="$HOME/$fileName"
      [[ ! -f $fileSpec ]] || fileSpecs+=( "$fileSpec" )
    done
    [[ ${#fileSpecs[@]} -eq 0 ]] || rm "${fileSpecs[@]}"
    [[ -z $OLD_PS1 ]] || PS1="$OLD_PS1"
    if grep -E "^(export )?(OLD_)?PS1=" "$_profileFileSpec" >/dev/null 2>&1; then
      grep -v -E "^(export )?(OLD_)?PS1=" "$_profileFileSpec" >"$newFileSpec"
      mv "$newFileSpec" "$_profileFileSpec"
    fi
  else

    local -a subOPT=(); local optVarName; for optVarName in ${__OPTVarNames[@]}; do
      [[ $optVarName != OPT_uninstall && $optVarName != OPT_help ]] || continue
      dstVarName=OPT srcVarName=$optVarName eval "$_codeToInitLocalCopyOfArray"
      subOPT+=( ${OPT[@]} )
    done
    local prompt; _preparePrompt "${subOPT[@]}" || return $?
    local profileLine="export PS1="\'"$prompt"\'
    local moduleName; for moduleName in ${supportModules[@]}; do
      profileLine+="; . ~/${moduleName}${supportFileNameSuffix}"
    done

    local moduleName; for moduleName in ${supportModules[@]}; do
      local fileName="${moduleName}${supportFileNameSuffix}"
      local srcFileSpec="$_bwDir/core/$fileName"
      local dstFileSpec="$HOME/$fileName"
      if [[ -f "$srcFileSpec" ]]; then
        cp "$srcFileSpec" "$dstFileSpec"
      else
        echo '[[ $(type -t _resetBash) != function ]] || _resetBash' >"$dstFileSpec"
        _getFileChunk "$_bwFileSpec" "# ==${moduleName} start" "# ==${moduleName} end" >>"$dstFileSpec"
      fi
    done
    if ! grep -x "$profileLine" "$_profileFileSpec" >/dev/null 2>&1; then
      local reOld=
      local profileLineOld=
      if [[ -z $OLD_PS1 ]]; then
        export OLD_PS1="$PS1"
        profileLineOld="export OLD_PS1=\"$PS1\""
        reOld='(OLD_)?'
      fi
      grep -v -E "^(export )?${reOld}PS1=" "$_profileFileSpec" >"$newFileSpec"
      [[ -z $profileLineOld ]] || echo "$profileLineOld" >>"$newFileSpec"
      echo "$profileLine" >>"$newFileSpec"
      mv "$newFileSpec" "$_profileFileSpec"
      export PS1="$prompt"
    fi
  fi
}

_codeToPrepareDescriptionsOf_bw_set_prompt='
  local bw_set_prompt_items_description="Задает состав и порядок элементов bash-prompt"
  local _psItems=( $( compgen -v | perl -ne "print if s/^_ps_([^_]+)\$/\$1/" ) )
  local howToTurnOff="Любое из значений ${_ansiSecondaryLiteral}off no -${_ansiReset} \"выключает\" элемент"
  local _psItem; for _psItem in ${_psItems[@]}; do
    local descriptionHolder="_ps_${_psItem}_description"
    local description=${_psItem}
    if [[ -n ${!descriptionHolder} ]]; then
      description=${!descriptionHolder}
    fi
    local bw_set_prompt_${_psItem}_description="Задает значение для элемента ${_ansiPrimaryLiteral}${description}${_ansiReset}. $howToTurnOff"
    local bw_set_prompt_${_psItem}Space_description="Определяет пробелы, окружающие элемент ${_ansiPrimaryLiteral}${description}${_ansiReset}"
    local bw_set_prompt_${_psItem}Ansi_description="Задает ansi-настройки (в т.ч.цвет) элемента ${_ansiPrimaryLiteral}${description}${_ansiReset}"
  done

  local bw_set_prompt_git_description="Задает присутствие в bash-prompt группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}. Любое из значений ${_ansiSecondaryLiteral}off no -${_ansiReset} \"выключает\" всю группу"

  local bw_set_prompt_gitPrefix_description="Задает префикс группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_gitPrefixSpace_description="Определяет пробелы, окружающие префикс группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}"
  local bw_set_prompt_gitPrefixAnsi_description="Задает ansi-настройки (в т.ч.цвет) префикса группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}"

  local bw_set_prompt_gitBranch_description="Задает присутствие элемента ${_ansiPrimaryLiteral}branch${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_gitBranchSpace_description="Определяет пробелы, окружающие элемент ${_ansiPrimaryLiteral}branch${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}"
  local bw_set_prompt_gitBranchAnsi_description="Задает ansi-настройки (в т.ч.цвет) элемента ${_ansiPrimaryLiteral}branch${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}"

  local bw_set_prompt_gitDirty_description="Задает присутствие элемента ${_ansiPrimaryLiteral}dirty${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_gitDirtySpace_description="Определяет пробелы, окружающие элемент ${_ansiPrimaryLiteral}dirty${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}"
  local bw_set_prompt_gitDirtyAnsi_description="Задает ansi-настройки (в т.ч.цвет) элемента ${_ansiPrimaryLiteral}dirty${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}"

  local bw_set_prompt_gitSuffix_description="Задает суффикс группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_gitSuffixSpace_description="Определяет пробелы, окружающие суффикс группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}"
  local bw_set_prompt_gitSuffixAnsi_description="Задает ansi-настройки (в т.ч.цвет) суффикса группы элементов ${_ansiPrimaryLiteral}git${_ansiReset}"

  local bw_set_prompt_error_description="Задает присутствие в bash-prompt группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}. Любое из значений ${_ansiSecondaryLiteral}off no -${_ansiReset} \"выключает\" всю группу"

  local bw_set_prompt_errorPrefix_description="Задает префикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_errorPrefixSpace_description="Определяет пробелы, окружающие префикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"
  local bw_set_prompt_errorPrefixAnsi_description="Задает ansi-настройки (в т.ч.цвет) префикса группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"

  local bw_set_prompt_errorInfix_description="Задает инфикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_errorInfixSpace_description="Определяет пробелы, окружающие инфикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"
  local bw_set_prompt_errorInfixAnsi_description="Задает ansi-настройки (в т.ч.цвет) инфикса группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"

  local bw_set_prompt_errorCode_description="Задает инфикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_errorCodeSpace_description="Определяет пробелы, окружающие инфикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"
  local bw_set_prompt_errorCodeAnsi_description="Задает ansi-настройки (в т.ч.цвет) инфикса группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"

  local bw_set_prompt_errorCode_description="Задает присутствие элемента ${_ansiPrimaryLiteral}code${_ansiReset} группы ${_ansiSecondaryLiteral}git${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_errorCodeSpace_description="Определяет пробелы, окружающие элемент ${_ansiPrimaryLiteral}code${_ansiReset} группы ${_ansiSecondaryLiteral}error${_ansiReset}"
  local bw_set_prompt_errorCodeAnsi_description="Задает ansi-настройки (в т.ч.цвет) элемента ${_ansiPrimaryLiteral}code${_ansiReset} группы ${_ansiSecondaryLiteral}error${_ansiReset}"

  local bw_set_prompt_errorSuffix_description="Задает суффикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}. $howToTurnOff"
  local bw_set_prompt_errorSuffixSpace_description="Определяет пробелы, окружающие суффикс группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"
  local bw_set_prompt_errorSuffixAnsi_description="Задает ansi-настройки (в т.ч.цвет) суффикса группы элементов ${_ansiPrimaryLiteral}error${_ansiReset}"

  local bw_set_prompt_separator_description="Задает окончание bash-prompt. $howToTurnOff"
  local bw_set_prompt_separatorSpace_description="Определяет пробелы, окружающие окончание bash-prompt"
  local bw_set_prompt_separatorAnsi_description="Задает ansi-настройки (в т.ч.цвет) окончания bash-prompt"
'

# =============================================================================

bw_set_viParams=()
bw_set_vi_description='Включает vi-режим для readline'
bw_set_vi() { eval "$_funcParams2"
  local -a opts
  if [[ -n $uninstall ]]; then
    opts=(
      --editing-mode emacs
      --show-mode-in-prompt off
    )
  else
    opts=(
      --editing-mode vi
      --show-mode-in-prompt on
      --vi-cmd-mode-string "\1${_ansiDarkGray}\2CMD \1${_ansiReset}\2"
      --vi-ins-mode-string "\1${_ansiMagenta}${_ansiBold}\2INS \1${_ansiReset}\2"
    )
  fi
  _inputrcSetProps "${opts[@]}"
}

# =============================================================================

_showResult() {
  local msg=
  local title="${_ansiPrimaryLiteral}$name"
  if [[ $returnCode -eq 0 ]]; then
    if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
      if [[ -n $titlePrefix ]]; then
        msg+="$(_upperFirst "$titlePrefix") "
      fi
      msg+="$title"
      msg+="${_ansiOK} "
      [[ -z $uninstall ]] \
        && msg+="$didInstall" \
        || msg+="$didUninstall"
      _ok "$msg"
    fi
  else
    if [[ $verbosity =~ ^(err|all.*)$ ]]; then
      msg="Не удалось "
      [[ -z $uninstall ]] \
        && msg+="$toInstall" \
        || msg+="$toUninstall"
      if [[ -n $titlePrefix ]]; then
        msg+=" ${_ansiErr}$(_lowerFirst "$titlePrefix")"
      fi
      msg+=" $title"
      _err "$msg"
    fi
  fi
}

_showProjectResult() {
  toFolder="в папку ${_ansiFileSpec}$folder"
  fromFolder="из папки ${_ansiFileSpec}$folder"
  titlePrefix=проект didInstall="установлен $toFolder" didUninstall="удален $fromFolder" toInstall="установить $toFolder" toUninstall="полностью удалить $fromFolder" _showResult
}

_showInstallResult() {
  titlePrefix= didInstall="установлен" didUninstall="удален" toInstall="установить" toUninstall="полностью удалить" _showResult
}

_showRunResult() {
  titlePrefix= didInstall="запущен" didUninstall="остановлен" toInstall="запустить" toUninstall="остановить" _showResult
}

# =============================================================================

_bwProjDefs=(
  'mls' '
    --gitOrigin github.com:baza-winner/mls-pm.git
  '
  'bgate' '
    --gitOrigin github.com:baza-winner/billing-gate.git
    --branch feature/docker
    --http 8086
    --https 8087
  '
  'bcore' '
    --gitOrigin github.com:baza-winner/billingcore.git
  '
  'crm' '
    --gitOrigin git@github.com:baza-winner/crm.git
  '
)

_bw_project_bgate() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}
_bw_project_bcore() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}

_getBwProjShortcuts() {
  local -a result=()
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    result+=( "${_bwProjDefs[$i]}" )
  done
  echo "${result[@]}"
}

_prepareBwProjVars() {
  [[ -n $bwProjShortcut ]] \
    || return $(_throw "ожидает, что переменная ${_ansiOutline}bwProjShortcut${_ansiErr} будет иметь непустое значение")
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    if [[ $bwProjShortcut == ${_bwProjDefs[$i]} ]]; then
      codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
      break
    fi
  done
  [[ -n $bwProjGitOrigin ]] \
    || return $(_err "not found gitOrigin for bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle")
}

_codeToDeclareLocalBwProjVars='
  local bwProjGitOrigin="" 
  local bwProjDefaultHttp=""
  local bwProjDefaultHttps="" 
  local bwProjDefaultBranch="" 
  local bwProjName=""
  local bwProjTitle=""
'
_codeToCallPrepareBwProjVarsHelper='
  bwProjShortcut="${_bwProjDefs[$i]}"
  local codeToGetBwProjDef
  _prepareBwProjVarsHelper ${_bwProjDefs[$((i + 1))]} || return $?
  codeHolder=codeToGetBwProjDef eval "$_evalCode"
  bwProjName=$(basename "$bwProjGitOrigin" .git)
  bwProjTitle="$bwProjName ($bwProjShortcut)"
'

_tcpPortDiap='1024..65535'
_prepareBwProjVarsHelperParams=(
  '--gitOrigin='
  '--http:'$_tcpPortDiap
  '--https:'$_tcpPortDiap
  '--branch='
)
_prepareBwProjVarsHelper() { eval "$_funcParams2"
  codeToGetBwProjDef='
    bwProjGitOrigin='$(_quotedArgs "$gitOrigin")'
    bwProjDefaultHttp='$(_quotedArgs "$http")'
    bwProjDefaultHttps='$(_quotedArgs "$https")'
    bwProjDefaultBranch='$(_quotedArgs "$branch")'
  '
}

_codeToPrepareDescriptionsOf_bw_project='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_project_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'

bw_projectParams=()
bw_projectParams() {
  varName="${FUNCNAME[0]}" codeHolder=_codeToUseCache eval "$_evalCode"
  local -a bwProjShortcuts=()
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    if [[ ! $bwProjShortcut =~ $_isValidVarNameRegExp ]]; then
      return $(_err "bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle${_ansiErr} $_mustBeValidVarName")
    elif _hasItem "$bwProjShortcut" "${bwProjShortcuts[@]}"; then
      return $(_err "Duplicate bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle")
    else
      bwProjShortcuts+=( "$bwProjShortcut" )
    fi
  done
  local bwProjShortcutsAsString="${bwProjShortcuts[@]}"
  verbosityDefault=allBrief silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
  bw_projectParams=(
    '!--uninstall/u'
    '!--force/f'
    '--branch='
    "${_verbosityParams[@]}"
    "bwProjShortcut:( $bwProjShortcutsAsString )"
    'projDir:?'
  )
  _saveToCache "${FUNCNAME[0]}"
}

bwProjGlobalDefaultBranch=develop
bwProjGlobalDefaultHttp=8080
bwProjGlobalDefaultHttps=8081
bw_projectParamsOpt=(--canBeMixedOptionsAndArgs)
_uninstall_description='Включить режим удаления'
_force_description='Игнорировать факт предыдущей установки'
bw_project_bwProjShortcut_name='Сокращенное-имя-проекта'
bw_project_projDir_name='Папка-проекта'
bw_project_projDir_description='
  Папка, куда будет установлен проект
  По умолчанию, в качестве папки проекта используется ${_ansiCmd}~/${_ansiOutline}Полное-имя-проекта${_ansiReset}
  ${_ansiOutline}Полное-имя-проекта${_ansiReset} - имя проекта на github'\''е
'
bw_project_branch_description='Ветка, на которую следует переключиться после установки проекта'
bw_project_description='Разворачивает/удаляет проект'
bw_projectShortcuts=( 'p' )
bw_project() { eval "$_funcParams2"
  codeHolder=_codeToInitSubOPT eval "$_evalCode"

  [[ -n $bwProjShortcut ]] || return $(_err "Не указано ${_ansiOutline}$bw_project_bwProjShortcut_description")
  
  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  [[ -n $branch ]] || branch=${bwProjDefaultBranch:-$bwProjGlobalDefaultBranch}

  local profileLineRegExp="^\s*\.\s+\"?(.+?)\/bin\/$bwProjShortcut\.bash\"?\s*$"
  local alreadyProjDir=
    ! grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1 || alreadyProjDir=$(cat "$_profileFileSpec" | perl -ne "print \$1 if /$profileLineRegExp/" | tail -n 1)
  if [[ -n $alreadyProjDir ]]; then
    if [[ -n $uninstall ]]; then
      projDir="$alreadyProjDir"
    elif [[ $verbosity != dry && -z $force ]]; then
      if [[ $verbosity != none ]]; then
        local cmd="${FUNCNAME[0]//_/ } $bwProjShortcut"
        local msg=
        msg+="Проект ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} уже установлен в ${_ansiFileSpec}$(_shortenFileSpec "$alreadyProjDir")${_ansiWarn}$_nl"
        msg+="Перед повторной установкой его необходимо удалить командой$_nl"
        msg+="  ${_ansiCmd}$cmd -u${_ansiWarn}$_nl"
        msg+="или установить с опцией ${_ansiCmd}--force${_ansiWarn}:$_nl"
        msg+="  ${_ansiCmd}$cmd -f${_ansiWarn}"
        _warn "$msg"
      fi
      return 4
    fi
  elif [[ -n $uninstall ]]; then
    [[ $verbosity == none ]] || _err "Проект ${_ansiPrimaryLiteral}$bwProjTitle${_ansiErr} не обнаружен"
    return 5
  fi

  [[ -n $projDir ]] || projDir="$HOME/$bwProjName"


  local returnCode=0
  while true; do

    if [[ -n $alreadyProjDir ]]; then
      local cmdFileSpec="$alreadyProjDir/bin/cmd.bash"
      if [[ -f "$cmdFileSpec" ]]; then
        fileSpec="$cmdFileSpec" _unsetBash ${sub_OPT_verbosity[1]}
      fi
    fi

    while true; do
      if [[ -z $uninstall ]]; then
        if [[ -d $projDir ]]; then
          local gitDirty
          if ! _inDir -v none "$projDir" _prepareGitDirty "$bwProjGitOrigin"; then
            if [[ -z $(ls -A "$projDir") ]]; then
              _rm "${sub_OPT[@]}" -d "$projDir" \
                || { returnCode=$?; break; }
            else
              _err "Папка ${_ansiCmd}$projDir${_ansiErr} существует и непуста. Ее надо предварительно удалить вручную" \
                || { returnCode=$?; break; }
            fi
          else
            _warn "Папка ${_ansiCmd}$projDir${_ansiWarn} уже содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle"
            _exec "${sub_OPT[@]}" cd "$projDir" || { returnCode=$?; break; }
            _hasItem "$gitDirty" '?' '*' '+' '^' || _exec "${sub_OPT[@]}" git pull
            break
          fi
        fi
        _mkDir "${sub_OPT[@]}" "$projDir" || { returnCode=$?; break; }
        _exec "${sub_OPT[@]}" git clone git@$bwProjGitOrigin "$projDir" || { returnCode=$?; break; }
        _exec "${sub_OPT[@]}" cd "$projDir"
        _exec "${sub_OPT[@]}" git checkout "$branch" || { returnCode=$?; break; }
        local funcName="_${FUNCNAME[0]}_$bwProjShortcut"
        ! _funcExists $funcName || $funcName || { returnCode=$?; break; }
      else
        if [[ -d $projDir ]]; then
          local gitDirty=
          if ! _inDir -v none "$projDir" _prepareGitDirty "$bwProjGitOrigin"; then
            if [[ -z $(ls -A "$projDir") ]]; then
              _rm "${sub_OPT[@]}" -d "$projDir" \
                || { returnCode=$?; break; }
            else
              _warn "Папка ${_ansiCmd}$projDir${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn}; оставлена для ручного удаления"
            fi
          elif _hasItem "$gitDirty" '?' '*' '+'; then
            _err "Репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} содержит изменения, проверьте ${_ansiCmd}git status" \
              || { returnCode=$?; break; }
          elif [[ $gitDirty == '^' ]]; then
            _err "Репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} содержит изменения, отсутствующие на сервере, проверьте ${_ansiCmd}git log --branches --not --remotes" \
              || { returnCode=$?; break; }
          elif  [[ $gitDirty == '$' ]]; then
            _err "Репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} stashed-изменения,проверьте ${_ansiCmd}git stash list" \
              || { returnCode=$?; break; }
          fi
          local needChangePwd=
          [[ $(cd "$projDir" && pwd) != $(pwd) ]] || needChangePwd=true
          _exec "${sub_OPT[@]}" rm -rf "$projDir" \
            || { returnCode=$?; break; }
          [[ -z $needChangePwd ]] || _exec "${sub_OPT[@]}" cd $HOME \
            || { returnCode=$?; break; }
          fi
      fi
      break
    done; [[ $returnCode -eq 0 ]] || break

    local cmdFileSpec="$projDir/bin/${bwProjShortcut}.bash"
    local profileLine=". $(_quotedArgs "$cmdFileSpec")"
    _setAtBashProfile ${OPT_uninstall[@]} "$profileLine" "$profileLineRegExp"
    break
  done

  folder="$projDir" name="$bwProjTitle" _showProjectResult

  if [[ $returnCode -eq 0 && -z $uninstall && $verbosity != dry ]]; then
    if [[ ! -f "$cmdFileSpec" ]]; then
      local msg=
      msg+="Не найден файл ${_ansiFileSpec}bin/$bwProjShortcut.bash${_ansiErr}$_nl"
      msg+="Не удалось инициализировать команду ${_ansiCmd}$bwProjShortcut"
      [[ $verbosity == none  ]] || _err "$msg"
      returnCode=1
    else
      _exec "${sub_OPT[@]}" . "$cmdFileSpec"
      local -a __completions=();
      local -a funcNames=( $(_getFuncNamesOfScriptToUnset "$cmdFileSpec") )
      _exec "${sub_OPT[@]}" _pregen "${funcNames[@]}"
      for _fileSpec in "${__completions[@]}"; do
        _exec "${sub_OPT[@]}" . "$_fileSpec"
      done
      [[ $verbosity == none  ]] || echo "${_ansiWarn}Теперь доступна команда ${_ansiCmd}$bwProjShortcut${_ansiReset}"
      _exec "${sub_OPT[@]}" --treatAsOK 3 $bwProjShortcut -?
    fi
  fi

  return $returnCode
}

_prepareGitDirtyParams=( 'originSuffix' )
_prepareGitDirty() { eval "$_funcParams2"
  local returnCode=0
  gitDirty=
  local gitOrigin=$(_gitOrigin)
  if [[ ${#gitOrigin} -ge ${#originSuffix} && $originSuffix == ${gitOrigin:$(( ${#gitOrigin} - ${#originSuffix} ))} ]]; then
    gitDirty=$(_gitDirty)
  else
    returnCode=1
  fi
  return $returnCode
}

# =============================================================================

_loopbackAlias='10.254.254.254'

_getDefaultXdebugRemoteHost() {
  if [[ $OSTYPE =~ ^darwin ]]; then
    echo host.docker.internal
  else
    echo "$_loopbackAlias"
  fi
}

_initBwProjCmd() {
  local fileSpec=$(_getSelfFileSpec 2)
  local bwProjShortcut=$(basename "$fileSpec" .bash)
  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  eval _${bwProjShortcut}FileSpec="$fileSpec"
  eval _${bwProjShortcut}Dir="$(dirname "$fileSpec")/.."
  export _${bwProjShortcut}DockerImageName="bazawinner/dev-${bwProjShortcut}"
  export _${bwProjShortcut}DockerContainerName="dev-${bwProjShortcut}"

  eval ${bwProjShortcut}_description=\"'Базовая утилита проекта ${_ansiPrimaryLiteral}'$bwProjName' ${_ansiUrl}'$bwProjGitOrigin'${_ansiReset}'\"
  eval ${bwProjShortcut}Params='()'
  eval ${bwProjShortcut}ParamsOpt='(--canBeMixedOptionsAndArgs --isCommandWrapper)'
  eval $bwProjShortcut'() { eval "$_funcParams2"; }'

  eval ${bwProjShortcut}_updateParams='( "--howDeep/d:(pregen source)" )'
  eval ${bwProjShortcut}_update_howDeep_description=\"'Определяет "глубину" обновления'\"
  eval ${bwProjShortcut}_update_howDeep_pregen_description=\"'только completion'\"
  eval ${bwProjShortcut}_update_howDeep_source_description=\"'source и completion'\"
  eval ${bwProjShortcut}_update_description=\"'Обновляет команду ${_ansiCmd}'$bwProjShortcut'${_ansiReset}'\"
  eval $bwProjShortcut'_update() { eval "$_funcParams2"
    local sourceFileSpec
    if _isInDocker; then
      sourceFileSpec=~/proj/docker/bin/cmd.bash
    else
      sourceFileSpec="$_'$bwProjShortcut'FileSpec"
    fi
    if [[ -z $howDeep ]]; then
      local sourceFileName=$(basename "$sourceFileSpec")
      if [[ $(basename "${BASH_SOURCE[1]}") == "$sourceFileName" || $(basename "${BASH_SOURCE[3]}") == "$sourceFileName" ]]; then
        howDeep=pregen
      else
        howDeep=source
      fi
    fi
    if [[ $howDeep == source ]]; then
      . "$_bwFileSpec" -p -
      . "$sourceFileSpec"
      rm -f "$_bwDir/generated/"'$bwProjShortcut'*
    fi
    local -a __completions=()
    _pregen $(compgen -c '$bwProjShortcut')
    for fileSpec in "${__completions[@]}"; do
      . "$fileSpec"
    done
  }'

  eval ${bwProjShortcut}_dockerParams='()'
  eval ${bwProjShortcut}_dockerParamsOpt='(--canBeMixedOptionsAndArgs --isCommandWrapper)'
  eval ${bwProjShortcut}_docker_description=\"'Docker-операции'\"
  eval ${bwProjShortcut}_dockerCondition=\""! _isInDocker"\"
  eval ${bwProjShortcut}'_docker() { eval "$_funcParams2"; }'

  local dockerImageTitle='docker-образ ${_ansiCmd}$_'$bwProjShortcut'DockerImageName${_ansiReset}'
  eval ${bwProjShortcut}_docker_buildParams='()'
  eval ${bwProjShortcut}_docker_build_description=\"'Собирает '$dockerImageTitle\"
  eval ${bwProjShortcut}'_docker_build() { eval "$_funcParams2"
    _isInDocker && return 4
    docker build -t "$_'$bwProjShortcut'DockerImageName" $_'$bwProjShortcut'Dir/docker
  }'

  eval ${bwProjShortcut}_docker_pushParams='()'
  eval ${bwProjShortcut}_docker_push_description=\"'Push'\''ит '$dockerImageTitle\"
  eval ${bwProjShortcut}'_docker_push() { eval "$_funcParams2"
    _isInDocker && return 4
    docker push "$_'$bwProjShortcut'DockerImageName"
  }'

  eval ${bwProjShortcut}_docker_upParamsOpt='(
    --additionalDependencies "$_'$bwProjShortcut'Dir/bin/'$bwProjShortcut'.bash"
  )'
  eval ${bwProjShortcut}_docker_upParams='( 
    --http:'$_tcpPortDiap'=${bwProjDefaultHttp:-$bwProjGlobalDefaultHttp} 
    --https:'$_tcpPortDiap'=${bwProjDefaultHttps:-$bwProjGlobalDefaultHttps} 
    @--restart/r
    --force-recreate/f
    "${'$bwProjShortcut'_docker_upParamsAddon[@]}"
  )'
  eval ${bwProjShortcut}_docker_up_http_description=\"'http-порт по которому будет доступен контейнер'\"
  eval ${bwProjShortcut}_docker_up_https_description=\"'https-порт по которому будет доступен контейнер'\"
  eval ${bwProjShortcut}_docker_up_restart_description=\"'Restart'\''ит (${_ansiCmd}docker-compose restart${_ansiReset}) указанные сервисы'\"
  eval ${bwProjShortcut}_docker_up_forceRecreate_description=\"'Up'\''ит указанные сервисы с опцией --force-recreate, передаваемой ${_ansiCmd}docker-compose up${_ansiReset}'\"
  eval ${bwProjShortcut}_docker_up_description=\"'Up'\''ит '$dockerImageTitle\"
  eval ${bwProjShortcut}'_docker_up() { eval "$_funcParams2"
    _isInDocker && return 4

    if [[ $https -eq $http ]]; then
      return $(_throw "ожидает, что значение ${_ansiOutline}http${ansiReset} ${_ansiPrimaryLiteral}$http${_ansiReset} будет отличаться от значения ${_ansiOutline}https")
    fi

    export _'$bwProjShortcut'DockerHttp="$http"
    export _'$bwProjShortcut'DockerHttps="$https"

    if [[ -z $'$bwProjShortcut'Prompt ]]; then
      local prompt; _preparePrompt --user '$bwProjShortcut' --userAnsi PrimaryLiteral --gitDirty -
      export _'$bwProjShortcut'Prompt="$prompt"
    fi

    local returnCode=0

    _inDir "$_'$bwProjShortcut'Dir/docker" '$bwProjShortcut'_docker_upHelper

    return $returnCode
  }'
  eval ${bwProjShortcut}'_docker_upHelper() { 
    local dockerComposeEnvFileSpec=".env"
    {
      echo "# file generated by '$bwProjShortcut'_docker_up"
      echo "_'$bwProjShortcut'DockerHttp=$_'$bwProjShortcut'DockerHttp"
      echo "_'$bwProjShortcut'DockerHttps=$_'$bwProjShortcut'DockerHttps"
      echo "_'$bwProjShortcut'DockerContainerName=$_'$bwProjShortcut'DockerContainerName"
      echo "_'$bwProjShortcut'DockerImageName=$_'$bwProjShortcut'DockerImageName"
    } > "$dockerComposeEnvFileSpec"

    local dockerContainerEnvFileSpec="main.env"
    { 
      echo "# file generated by '$bwProjShortcut'_docker_up"
      echo "_isBwDevelopInherited=$_isBwDevelop"
      echo "BW_SELF_UPDATE_SOURCE=$BW_SELF_UPDATE_SOURCE"
      echo "_bwProjName='$bwProjName'"
      echo "_bwProjShortcut='$bwProjShortcut'"
      echo "_hostUser=$(whoami)"
      echo "_dockerHttp=$http"
      echo "_dockerHttps=$https"
      echo "_prompt=$_'$bwProjShortcut'Prompt"

      local addonFuncName=_'$bwProjShortcut'_docker_upAddon
      if _funcExists $addonFuncName; then 
        $addonFuncName
      else
        _debugVar addonFuncName
      fi
    } > "$dockerContainerEnvFileSpec"

    local stderrFileSpec="/tmp/docker-compose.stderr"
    {
      if [[ "${#restart[@]}" -eq 0 ]]; then
        [[ -z $forceRecreate ]] || OPT_forceRecreate=( --force-recreate )
        _dockerCompose up -d "${OPT_forceRecreate[@]}" --remove-orphans
      else
        _dockerCompose restart "${restart[@]}"
      fi
    } 2> >(tee "$stderrFileSpec"); returnCode=$?
    if [[ $returnCode -ne 0 ]] && grep -P -o '\''(?<=:)\d+(?= failed: port is already allocated)'\'' "$stderrFileSpec" >/dev/null; then
      _err "port is already allocated"
    fi
  }'

  eval '_'$bwProjShortcut'_dockerContainerNames() {
    local -a containerNames=()
    local containerNameHolder; for containerNameHolder in $(cat "$_'$bwProjShortcut'Dir/docker/docker-compose.yml" | perl -ne "print if s/^\s*container_name:\s*(.+?)\s*/\$1/"); do
      eval containerNames+=\( $containerNameHolder \)
    done
    echo "${containerNames[@]}"
  }'
  # eval ${bwProjShortcut}_docker_bashParamsOpt='( 
  #   --additionalDependencies "'$(_getSelfFileSpec)'"
  # )'
  eval ${bwProjShortcut}_docker_bashParams='( 
    '\''containerName:( $(_'$bwProjShortcut'_dockerContainerNames) )='\'\$_${bwProjShortcut}DockerContainerName'
  )'
  eval ${bwProjShortcut}'_docker_bash_description='\"'Запускает bash в Docker-контейнере'\"
  eval ${bwProjShortcut}'_docker_bash_containerName_name='\"'Имя-контейнера'\"
  eval ${bwProjShortcut}'_docker_bash() { eval "$_funcParams2"
    _docker exec -it "$containerName" bash
  }'


  eval ${bwProjShortcut}_docker_down_description=\"'Down'\''ит '$dockerImageTitle\"
  eval ${bwProjShortcut}_docker_downParams='()'
  eval ${bwProjShortcut}'_docker_down() { eval "$_funcParams2"
    _isInDocker && return 4
    _inDir "$_'$bwProjShortcut'Dir/docker" _dockerCompose down
  }'


}

# =============================================================================

bw_projectInfoShortcuts=( 'pi' )
bw_projectInfoParamsOpt=( --canBeMixedOptionsAndArgs )
bw_projectInfoParams=()
bw_projectInfoParams() {
  varName="${FUNCNAME[0]}" codeHolder=_codeToUseCache eval "$_evalCode"
  local -a bwProjShortcuts=()
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    if [[ ! $bwProjShortcut =~ $_isValidVarNameRegExp ]]; then
      return $(_err "bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle${_ansiErr} $_mustBeValidVarName")
    elif _hasItem "$bwProjShortcut" "${bwProjShortcuts[@]}"; then
      return $(_err "Duplicate bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle")
    else
      bwProjShortcuts+=( "$bwProjShortcut" )
    fi
  done
  local bwProjShortcutsAsString="${bwProjShortcuts[@]}"
  bw_projectInfoParams=(
    '--all/a'
    '--def/d'
    "bwProjShortcut:?:( $bwProjShortcutsAsString )"
  )
  _saveToCache "${FUNCNAME[0]}"
}
_codeToPrepareDescriptionsOf_bw_projectInfo='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_projectInfo_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'
bw_projectInfo_bwProjShortcut_name="Сокращенное-имя-проекта"
bw_projectInfo_bwProjShortcut_description='
  Если ${_ansiOutline}$bw_projectInfo_bwProjShortcut_name${_ansiReset} не задано, то выводит информация обо всех ${_ansiOutline}обнаруженных${_ansiReset} проектах
  С опцией ${_ansiCmd}--all${_ansiReset} -- обо всех проектах
'
bw_projectInfo_all_description='
  Вывести информацию обо всех проектах
  Без опции ${_ansiCmd}--all${_ansiReset} -- обо всех ${_ansiOutline}обнаруженных${_ansiReset} проектах
'
bw_projectInfo_def_description='
  Вывести ${_ansiOutline}определение${_ansiReset} проекта/проектов
'
bw_projectInfo_description="Выводит информацию о проекте/проектах"
bw_projectInfo() { eval "$_funcParams2"
  local skipNonExistent=
  if [[ -n $def ]]; then
    local -a duplicatePorts; _prepareDuplicatePorts
  fi
  if [[ -n $bwProjShortcut ]]; then
    eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
    _bwProjectInfoHelper
  else
    [[ -n $all ]] || skipNonExistent=true
    local -a found=()
    eval "$_codeToDeclareLocalBwProjVars"
    local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
      local bwProjShortcut="${_bwProjDefs[$i]}"
      codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
      _bwProjectInfoHelper; local returnCode=$?
      if [[ $returnCode -ne 5 ]]; then
        found+=( "$bwProjTitle" )
      fi
    done
    if [[ -n $def ]]; then
      echo "Всего $(_getPluralWord ${#found[@]} определен определено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}${found[@]}${_ansiReset}"
    else
      if [[ ${#found[@]} -gt 0 ]]; then
        _ok "Всего $(_getPluralWord ${#found[@]} обнаружен обнаружено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}${found[@]}"
      else
        _warn "Не обнаружено ни одного проекта"
      fi
    fi
  fi
}

_prepareDuplicatePorts() {
  unset $(compgen -v __usedPort)
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    if [[ -n $bwProjDefaultHttp ]]; then
      eval local __usedPort$bwProjDefaultHttp
      eval __usedPort$bwProjDefaultHttp=\$\(\( __usedPort$bwProjDefaultHttp + 1 \)\)
    fi
    if [[ -n $bwProjDefaultHttps ]]; then
      eval local __usedPort$bwProjDefaultHttps
      eval __usedPort$bwProjDefaultHttps=\$\(\( __usedPort$bwProjDefaultHttps + 1 \)\)
    fi
  done
  duplicatePorts=()
  local varName; for varName in $(compgen -v __usedPort); do
    if [[ $varName -gt 1 ]]; then
      duplicatePorts+=( ${varName:10})
    fi
  done
}

_bwProjectInfoHelper() {
  if [[ -n $def ]]; then
    echo "$bwProjShortcut:"
    echo "  name: $bwProjName"
    echo "  gitOrigin: ${_ansiUrl}"$(_quotedArgs "$bwProjGitOrigin")"${_ansiReset}"

    echo -n "  branch: "
    local branch=${bwProjDefaultBranch:-$bwProjGlobalDefaultBranch}
    if [[ $branch == "$bwProjGlobalDefaultBranch" ]]; then
      echo "$branch"
    else
      echo "${_ansiWarn}$branch${_ansiReset}"
    fi

    echo -n "  http: "
    if [[ -z $bwProjDefaultHttp || $bwProjDefaultHttp -eq "$bwProjGlobalDefaultHttp" ]]; then
      echo "${_ansiWarn}$bwProjGlobalDefaultHttp${_ansiReset}"
    elif _hasItem "$bwProjDefaultHttp" "${duplicatePorts[@]}"; then
      echo "${_ansiWarn}$bwProjDefaultHttp${_ansiReset}"
    else
      echo "$bwProjDefaultHttp"
    fi

    echo -n "  https: "
    if [[ -z $bwProjDefaultHttps || $bwProjDefaultHttps -eq "$bwProjGlobalDefaultHttps" ]]; then
      echo "${_ansiWarn}$bwProjGlobalDefaultHttps${_ansiReset}"
    elif _hasItem "$bwProjDefaultHttps" "${duplicatePorts[@]}"; then
      echo "${_ansiWarn}$bwProjDefaultHttps${_ansiReset}"
    else
      echo "$bwProjDefaultHttps"
    fi

  else
    local profileLineRegExp="^\s*\.\s+\"?(.+?)\/bin\/$bwProjShortcut\.bash\"?\s*$"
    if grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1; then
      local alreadyProjDir=$(cat "$_profileFileSpec" | perl -ne "print \$1 if /$profileLineRegExp/" | tail -n 1)
      if [[ ! -d $alreadyProjDir ]]; then
        _warn "Папка ${_ansiFileSpec}$alreadyProjDir${_ansiWarn} проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} не обнаружена"
        return 7
      else
        if ! _inDir "$alreadyProjDir" _prepareGitDirty "$bwProjGitOrigin"; then
          _warn "Папка ${_ansiFileSpec}$alreadyProjDir${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn}"
          return 6
        else
          local gitBranchName=; _inDir "$alreadyProjDir" _prepareGitBranchName
          _ok "Ветка ${_ansiSecondaryLiteral}$gitBranchName${_ansiOK} проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiOK} обнаружена в ${_ansiFileSpec}$alreadyProjDir"
          return 0
        fi
      fi
    else
      [[ -n $skipNonExistent ]] || _warn "Проект ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} не обнаружен"
      return 5
    fi
  fi
}

# =============================================================================

bw_installParams=(
  "${_verbosityParams[@]}"
  "--force/f"
  # '!--uninstall/u'
)
bw_installParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_install_cmd_name=Имя-приложения
bw_install_description='устанавливает приложения'
bw_install_force_description='устанавливает приложения, даже если оно уже установлено'
bw_install() { eval "$_funcParams2"
}

_codeToInstallApp='
  # uninstallTitle=Удаление
  showResult=_showInstallResult codeHolder=_codeToInstallOrRunApp eval "$_evalCode"
'
_codeToRunApp='
  # uninstallTitle=Останов
  showResult=_showRunResult codeHolder=_codeToInstallOrRunApp eval "$_evalCode"
'
_codeToInstallOrRunApp='
  codeHolder=_codeToInitSubOPT eval "$_evalCode"
  local returnCode=0
  # [[ -z $uninstall ]] || return $(_err "$uninstallTitle ${_ansiPrimaryLiteral}name${_ansiErr} не поддерживается")
  _osSpecific || return $?
  $showResult
  return $returnCode
'
_osSpecific() {
  local funcName=${FUNCNAME[1]}
  local osSpecificFuncName=
  if [[ $OSTYPE =~ ^darwin ]]; then
    local osSpecificFuncName="_${funcName}Darwin"
  elif [[ $OSTYPE =~ ^linux ]]; then
    local osSpecificFuncName="_${funcName}Linux"
  fi
  if [[ -z $osSpecificFuncName ]] || ! _funcExists $osSpecificFuncName; then
    _err "Неподдерживамая ОС ${_ansiPrimaryLiteral}$OSTYPE"
    return 1
  else
    $osSpecificFuncName
  fi
}

# =============================================================================

bw_install_brewParams=()
bw_install_brew_description='Устанавливает Homebrew'
bw_install_brew() { eval "$_funcParams2"
  name=git codeHolder=_codeToInstallApp eval "$_evalCode"
}

_bw_install_brewDarwin() {
  _exec "${sub_OPT[@]}" /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
}

# =============================================================================

# Install Bash version 4 on MacOS X: https://gist.github.com/samnang/1759336

# =============================================================================

bw_install_gitParams=()
bw_install_git_description='Устанавливает git'
bw_install_git() { eval "$_funcParams2"
  if [[ -z $force ]] && which git >/dev/null 2>&1; then
    [[ $verbosity =~ ^(ok|all)$ ]] && _ok "${_ansiCmd}git${_ansiOK} уже установлен"
  else
    name=git codeHolder=_codeToInstallApp eval "$_evalCode"
  fi
}
_bw_install_gitLinux() {
  _exec "${sub_OPT[@]}" --sudo apt install -y --force-yes git || returnCode=$?
}

# =============================================================================

bw_install_githubKeygenParams=()
bw_install_githubKeyGen_description='Устанавливает github-keygen (${_ansiUrl}https://github.com/dolmen/github-keygen${_ansiReset})'
bw_install_githubKeygen() { eval "$_funcParams2"
  if [[ -z $force && -f "$_bwDir/github-keygen" ]]; then
    [[ $verbosity =~ ^(ok|all)$ ]] && _ok "${_ansiCmd}github-keygen${_ansiOK} уже установлен"
  else
    name=git codeHolder=_codeToInstallApp eval "$_evalCode"
  fi
}
_bw_install_githubKeygenDarwin() {
  _exec git clone https://github.com/dolmen/github-keygen.git "$_bwDir/github-keygen"
}
_bw_install_githubKeygenLinux() {
  _bw_install_githubKeygenDarwin || return $?
  _exec --sudo apt-get update
  _exec --sudo apt-get install -y --force-yes xclip
}

bw_githubKeygenParams=( 'username' )
bw_githubKeygen() { eval "$_funcParams2"
  bw_install_githubKeygen
  local stdoutFileSpec="/tmp/github-keygen.stdout"
  "$_bwDir/github-keygen" "$username" | tee "$stdoutFileSpec"
  if [[ $(tail -n 1 "$stdoutFileSpec") =~ ^Paste ]]; then
    read -p "Press enter to continue" # https://unix.stackexchange.com/questions/293940/bash-how-can-i-make-press-any-key-to-continue
    _osSpecific || return $?
  fi
  
}
_githubKeysUrl='https://github.com/settings/keys'
_bw_githubKeygenDarwin() {
  local appGoogleChrome="/Applications/Google Chrome.app"
  if [[ -d "$appGoogleChrome" ]]; then
    /usr/bin/open -a "$appGoogleChrome" "$_githubKeysUrl"
  else
    /usr/bin/open "$_githubKeysUrl"
  fi
}
_bw_githubKeygenLinux() {
  if which google-chrome >/dev/null 2>&1; then
    google-chrome "$_githubKeysUrl"
  else
    xdg-open "$_githubKeysUrl"
  fi
}

# =============================================================================

bw_install_dockerParams=()
bw_install_docker_description="Устанавливает DockerCE ${_ansiUrl}https://www.docker.com/community-edition${_ansiReset}"
bw_install_docker() { eval "$_funcParams2"
  if [[ -z $force ]] || bw_install_dockerHelper; then
    [[ $verbosity =~ ^(ok|all)$ ]] && _ok "${_ansiCmd}github-keygen${_ansiOK} уже установлен"
  else
    name=Docker codeHolder=_codeToInstallApp eval "$_evalCode"
  fi
}
bw_install_dockerHelper() {
  which docker >/dev/null 2>&1 && [[ $(docker --version | perl -e '$\=undef; $_=<STDIN>; printf("%d%04d", $1, $2) if m/(\dd+)/'\d) -ge 180003 ]]
}
_bw_install_dockerDarwin() {
  while true; do
    local -r appName=Docker
    local -r applicationsPath='/Applications'
    local -r appDir="${applicationsPath}/${appName}.app"
    local -r dmgFileSpec=~/Downloads/Docker.dmg
    if [[ -d $appDir && -z $force ]]; then
      _warn "${_ansiCmd}$appName${_ansiWarn} уже установлен"
    else
      local -r sourceUrlOfDockerDmg='https://download.docker.com/mac/stable/Docker.dmg'
      _download -c etag -r 3 "$sourceUrlOfDockerDmg" "$dmgFileSpec" || retunCode=$?
      local -r volumePath="/Volumes/$(basename "$dmgFileSpec" .dmg)"
      if [[ $returnCode -eq 3 || ! -d $appDir || -n $force ]]; then
        returnCode=0
        [[ ! -d $appDir ]] || _rm -d $appDir || { returnCode=$?; break; }
        _exec -v err --silent hdiutil attach "$dmgFileSpec" || { returnCode=$?; break; }
        _exec -v err cp -R "${volumePath}/${appName}.app" "$applicationsPath" || { returnCode=$?; break; }
        _exec -v none hdiutil detach "$volumePath"
      fi
    fi
    break
  done
}
_bw_install_dockerLinux() {
  while true; do
    # https://docs.docker.com/install/linux/docker-ce/ubuntu/#supported-storage-drivers
    _exec "${sub_OPT[@]}" --sudo apt-get update || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --cmdAsIs 'curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -' || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo apt-get update || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo apt-get install -y docker-ce || { returnCode=$?; break; }
    # https://docs.docker.com/compose/install/#install-compose
    _exec "${sub_OPT[@]}" --sudo curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo chmod +x /usr/local/bin/docker-compose || { returnCode=$?; break; }
    break
  done
}

# =============================================================================

bw_runParams=(
  "${_verbosityParams[@]}"
)
bw_runParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_run_cmd_name=Имя-приложения
bw_run_description='запускает приложение'
bw_run() { eval "$_funcParams2"
}

# =============================================================================

bw_run_dockerParams=()
bw_run_docker_description="Запускает Docker"
bw_run_docker() { eval "$_funcParams2"
  name=Docker codeHolder=_codeToRunApp eval "$_evalCode"
}
_bw_run_dockerLinux() {
  _warn "В Linux специальный запуск ${_ansiPrimaryLiteral}Docker'а${_ansiWarn} не требуется"
}
_bw_run_dockerDarwin() {
  while true; do
    local -r appName=Docker
    local -r applicationsPath='/Applications'
    local -r appDir="${applicationsPath}/${appName}.app"
    [[ -d $appDir ]] || _installDocker || { returnCode=$?; break; }
    _exec -v err --silent --untilSuccessSleep=1 open -a "$appName" || { returnCode=$?; break; } # http://osxdaily.com/2007/02/01/how-to-launch-gui-applications-from-the-terminal/
    _spinner "Ожидание запуска ${_ansiPrimaryLiteral}$appName${_ansiReset}" _bw_run_dockerDarwinHelper || { returnCode=$?; break; }
    break
  done
}
_bw_run_dockerDarwinHelper() {
  while ! _silent docker ps; do
    # https://stackoverflow.com/questions/32041674/linux-how-to-kill-sleep/32049811#32049811
    sleep 1 &
    wait
  done
}

# =============================================================================

_redDir="$_bwDir/red"
_redForMacUrl="https://static.red-lang.org/dl/mac/red-063"
_redForMacFileSpec="$_redDir/red_063_for_mac"
_redForLinuxUrl="https://static.red-lang.org/dl/linux/red-063"
_redForLinuxFileSpec="$_redDir/red_063_for_linux"
red() {
  local returnCode=0
  _osSpecific || return $?
  if [[ $returnCode -eq 99 ]]; then
    _err "Не удалось сделать файл ${_ansiCmd}$_redForMacFileSpec${_ansiErr} исполняемым"
  fi
  return $returnCode
}
_redDarwin() {
  while true; do
    _download -s $_redForMacUrl $_redForMacFileSpec || { returnCode=$?; break; }
    chmod u+x $_redForMacFileSpec || { returnCode=99; break; }
    $_redForMacFileSpec "$@" || { returnCode=$?; break; }
    break
  done
}
_redLinux() {
  while true; do
    _download -s $_redForLinuxUrl $_redForLinuxFileSpec || { returnCode=$?; break; }
    chmod u+x $_redForLinuxFileSpec || { returnCode=99; break; }
    $_redForLinuxFileSpec "$@" || { returnCode=$?; break; }
    break
  done
}

# =============================================================================
