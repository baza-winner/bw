# =============================================================================

_resetBash

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
bw_version_description='выводит номер версии bw.bash'
bw_version() { eval "$_funcParams2"
  local suffix="$BW_SELF_UPDATE_SOURCE"
  [[ -z $suffix ]] || suffix="@$suffix"
  echo "$_bwVersion$suffix"
}

# =============================================================================

bw_update_description="обновляет bw"
bw_update() { eval "$_funcParams2"
  . "$_bwFileSpec" -p -
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

_enumAnsi='(Reset Bold Dim Italic Underline Black Red Green Yellow Blue Magenta Cyan LightGray LightGrey DarkGray DarkGrey LightRed LightGreen LightYellow LightBlue LightMagenta LightCyan White)'
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

_bw_project_billgate() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}
_bw_project_billcore() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}


_mlsGitOrigin='github.com:baza-winner/mls-pm.git'
_billgateGitOrigin='github.com:baza-winner/billing-gate.git'
_billcoreGitOrigin='github.com:baza-winner/billingcore.git'
bw_projectParams=(
  '!--uninstall/u'
  '!--force/f'
  '--branch=develop'
  '--info/i'
  '--all/a'
  "${_verbosityParams[@]}"
  'projName:?:(mls billgate billcore)'
  'projDir:?'
)
bw_projectParamsOpt=(--canBeMixedOptionsAndArgs)
_uninstall_description='включает режим удаления'
_force_description='игнорировать факт предыдущей установки'
bw_project_projName_description=Имя-проекта
bw_project_description='разворачивание/удаление проекта'
bw_projectShortcuts=( 'p' )
bw_project() { eval "$_funcParams2"
  codeHolder=_codeToInitSubOPT eval "$_evalCode"

  if [[ -n $info ]]; then
    _bwProjectInfo
  else
    [[ -n $projName ]] || return $(_err "Не указано ${_ansiOutline}$bw_project_projName_description")
    local gitOrigin=; _bwPrepareGitOrigin || return $?

    local profileLineRegExp="^\s*\.\s+\"?(.+?)\/bin\/$projName\.bash\"?\s*$"
    local alreadyProjDir=
      ! grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1 || alreadyProjDir=$(cat "$_profileFileSpec" | perl -ne "print \$1 if /$profileLineRegExp/" | tail -n 1)
    if [[ -n $alreadyProjDir ]]; then
      if [[ -n $uninstall ]]; then
        projDir="$alreadyProjDir"
      elif [[ $verbosity != dry && -z $force ]]; then
        if [[ $verbosity != none ]]; then
          local cmd="${FUNCNAME[0]//_/ } $projName"
          local msg=
          msg+="Проект ${_ansiPrimaryLiteral}$projName${_ansiWarn} уже установлен в ${_ansiFileSpec}$(_shortenFileSpec "$alreadyProjDir")${_ansiWarn}$_nl"
          msg+="Перед повторной установкой его необходимо удалить командой$_nl"
          msg+="  ${_ansiCmd}$cmd -u${_ansiWarn}$_nl"
          msg+="или установить с опцией ${_ansiCmd}--force${_ansiWarn}:$_nl"
          msg+="  ${_ansiCmd}$cmd -f${_ansiWarn}"
          _warn "$msg"
        fi
        return 4
      fi
    elif [[ -n $uninstall ]]; then
      [[ $verbosity == none ]] || _err "Проект ${_ansiPrimaryLiteral}$projName${_ansiWarn} не обнаружен"
      return 5
    fi

    [[ -n $projDir ]] || projDir="$HOME/$(basename "$gitOrigin" .git)"


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
            if ! _inDir -v none "$projDir" _prepareGitDirty "$gitOrigin"; then
              if [[ -z $(ls -A "$projDir") ]]; then
                _rm "${sub_OPT[@]}" -d "$projDir" \
                  || { returnCode=$?; break; }
              else
                # _exec "${sub_OPT[@]}" cd "$projDir"
                _err "Папка ${_ansiCmd}$projDir${_ansiErr} существует и непуста. Ее надо предварительно удалить вручную" \
                  || { returnCode=$?; break; }
              fi
            else
              _warn "Папка ${_ansiCmd}$projDir${_ansiWarn} уже содержит репозиторий проекта ${_ansiPrimaryLiteral}$projName"
              _exec "${sub_OPT[@]}" cd "$projDir" || { returnCode=$?; break; }
              _hasItem "$gitDirty" '?' '*' '+' '^' || _exec "${sub_OPT[@]}" git pull
              break
            fi
          fi
          _mkDir "${sub_OPT[@]}" "$projDir" || { returnCode=$?; break; }
          _exec "${sub_OPT[@]}" git clone git@$gitOrigin "$projDir" || { returnCode=$?; break; }
          _exec "${sub_OPT[@]}" cd "$projDir"
          _exec "${sub_OPT[@]}" git checkout "$branch" || { returnCode=$?; break; }
          local funcName="_${FUNCNAME[0]}_$projName"
          ! _funcExists $funcName || $funcName || { returnCode=$?; break; }
        else
          if [[ -d $projDir ]]; then
            local gitDirty=
            if ! _inDir -v none "$projDir" _prepareGitDirty "$gitOrigin"; then
              if [[ -z $(ls -A "$projDir") ]]; then
                _rm "${sub_OPT[@]}" -d "$projDir" \
                  || { returnCode=$?; break; }
              else
                _warn "Папка ${_ansiCmd}$projDir${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$projName${_ansiWarn}; оставлена для ручного удаления"
              fi
            elif _hasItem "$gitDirty" '?' '*' '+'; then
              _err "Репозиторий проекта ${_ansiCmd}$projName${_ansiWarn} содержит изменения, проверьте ${_ansiCmd}git status" \
                || { returnCode=$?; break; }
            elif [[ $gitDirty == '^' ]]; then
              _err "Репозиторий проекта ${_ansiCmd}$projName${_ansiWarn} содержит изменения, отсутствующие на сервере, проверьте ${_ansiCmd}git log --branches --not --remotes" \
                || { returnCode=$?; break; }
            elif  [[ $gitDirty == '$' ]]; then
              _err "Репозиторий проекта ${_ansiCmd}$projName${_ansiWarn} stashed-изменения,проверьте ${_ansiCmd}git stash list" \
                || { returnCode=$?; break; }
            fi
          fi
          local needChangePwd=
          [[ $(cd "$projDir" && pwd) != $(pwd) ]] || needChangePwd=true
          _exec "${sub_OPT[@]}" rm -rf "$projDir" \
            || { returnCode=$?; break; }
          [[ -z $needChangePwd ]] || _exec "${sub_OPT[@]}" cd $HOME \
            || { returnCode=$?; break; }
        fi
        break
      done; [[ $returnCode -eq 0 ]] || break

      local cmdFileSpec="$projDir/bin/${projName}.bash"
      local profileLine=". $(_quotedArgs "$cmdFileSpec")"
      local perlCode=
      if [[ $verbosity == dry ]]; then
        echo "${_ansiCmd}echo \"$profileLine\" >> \"$_profileFileSpec\"${_ansiReset}"
      else
        if grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1; then
          [[ -n $uninstall ]] \
            && perlCode="print unless /$profileLineRegExp/" \
            || perlCode="if (! /$profileLineRegExp/) { print } elsif (! \$state) { print $(_quotedArgs --quote:all "$profileLine") . \"\n\"; \$state=1 }"
        elif [[ -z $uninstall ]]; then
          echo "$profileLine" >> "$_profileFileSpec"
        fi
        if [[ -n $perlCode ]]; then
          local newFileSpec="$_profileFileSpec.new"
          cat "$_profileFileSpec" | perl -ne "$perlCode" > "$_profileFileSpec.new"
          mv "$_profileFileSpec.new" "$_profileFileSpec"
        fi
      fi

      break
    done

    folder="$projDir" name="$projName" _showProjectResult

    if [[ $returnCode -eq 0 && -z $uninstall && $verbosity != dry ]]; then
      if [[ ! -f "$cmdFileSpec" ]]; then
        local msg=
        msg+="Не найден файл ${_ansiFileSpec}bin/$projName.bash${_ansiErr}$_nl"
        msg+="Не удалось инициализировать команду ${_ansiCmd}$cmdFileSpec"
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
        [[ $verbosity == none  ]] || echo "${_ansiWarn}Теперь доступна команда ${_ansiCmd}$projName${_ansiReset}"
        _exec "${sub_OPT[@]}" --treatAsOK 3 $projName -?
      fi
    fi

    return $returnCode
  fi
}

_bwPrepareGitOrigin() {
  gitOrigin=
  local gitOriginHolder="_${projName}GitOrigin"
  if [[ -n ${!gitOriginHolder} ]]; then
    gitOrigin=${!gitOriginHolder}
  else
    return $(_err "Не задана переменная ${_ansiOutline}$gitOriginHolder")
  fi
}

_bwProjectInfo() {
  local skipNonExistent=
  if [[ -n $projName ]]; then
    _bwProjectInfoHelper
  else
    eval local -a enumValues="$__ENUM_projName"
    [[ -n $all ]] || skipNonExistent=true
    for projName in "${enumValues[@]}"; do
      _bwProjectInfoHelper
    done
  fi
}

_bwProjectInfoHelper() {
  local profileLineRegExp="^\s*\.\s+\"?(.+?)\/bin\/$projName\.bash\"?\s*$"
  if grep -E "$profileLineRegExp" "$_profileFileSpec" >/dev/null 2>&1; then
    local alreadyProjDir=$(cat "$_profileFileSpec" | perl -ne "print \$1 if /$profileLineRegExp/" | tail -n 1)
    if [[ ! -d $alreadyProjDir ]]; then
      _warn "Папка ${_ansiFileSpec}$alreadyProjDir${_ansiWarn} проекта ${_ansiPrimaryLiteral}$projName${_ansiWarn} не обнаружена"
      return 7
    else
      local gitOrigin=; _bwPrepareGitOrigin || return $?
      if ! _inDir "$alreadyProjDir" _prepareGitDirty "$gitOrigin"; then
        _warn "Папка ${_ansiFileSpec}$alreadyProjDir${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$projName${_ansiWarn}"
        return 6
      else
        local gitBranchName=; _inDir "$alreadyProjDir" _prepareGitBranchName
        _ok "Ветка ${_ansiSecondaryLiteral}$gitBranchName${_ansiOK} проекта ${_ansiPrimaryLiteral}$projName${_ansiOK} обнаружена в ${_ansiFileSpec}$alreadyProjDir"
        return 0
      fi
    fi
  elif [[ -z $skipNonExistent ]]; then
    _warn "Проект ${_ansiPrimaryLiteral}$projName${_ansiWarn} не обнаружен"
    return 5
  fi
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

bw_installParams=(
  "${_verbosityParams[@]}"
  # '!--uninstall/u'
)
bw_installParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_install_cmd_name=Имя-приложения
bw_install_description='устанавливает приложения'
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
  name=git codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_gitLinux() {
  _exec "${sub_OPT[@]}" --sudo apt install -y git || returnCode=$?
}

# =============================================================================

bw_install_dockerParams=()
bw_install_docker_description="Устанавливает DockerCE ${_ansiUrl}https://www.docker.com/community-edition${_ansiReset}"
bw_install_docker() { eval "$_funcParams2"
  name=Docker codeHolder=_codeToInstallApp eval "$_evalCode"
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
