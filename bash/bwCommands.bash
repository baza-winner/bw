#!/bin/bash

# shellcheck disable=SC2154,SC1090,SC2016
true

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

export _codeToInitSubOPT='
  local -a sub_OPT_silent=( "${OPT_silent[@]}" )
  local -a sub_OPT_verbosity=()
  if [[ $verbosity =~ ^(none|dry|all)$ ]]; then
    sub_OPT_verbosity=( "${OPT_verbosity[@]}" )
  else
    sub_OPT_verbosity=( --verbosity err )
  fi
  local -a sub_OPT=( "${sub_OPT_silent[@]}" "${sub_OPT_verbosity[@]}" )
'
verbosityDefault=allBrief silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"

# =============================================================================

export bwParams=()
export bwParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
export bw_description='Базовая утилита проектов baza-winner'
bw() { eval "$_funcParams2"
}

# =============================================================================

export bw_versionParams=()
export bw_version_description='Выводит номер версии bw.bash'
bw_version() { eval "$_funcParams2"
  local suffix="$BW_SELF_UPDATE_SOURCE"
  [[ -z $suffix ]] || suffix="@$suffix"
  echo "$_bwVersion$suffix"
}

# =============================================================================

export bw_updateParams=( 
  '--remove/r' 
  '--noPregen/n'
)
export bw_update_remove_description="Удалить прегенеренные файлы перед обновлением"
export bw_update_noPregen_description="Пропустить прегенерацию файлов после обновления"
export bw_update_description="Обновляет bw"
bw_update() { eval "$_funcParams2"
  if [[ -n $remove ]]; then
    "$_bwFileSpec" rm -y
  fi
  local -a OPT=()
  if [[ -n $noPregen ]]; then
    OPT=( -p - )
  fi
  . "$_bwFileSpec" "${OPT[@]}"
  echo "Current version: $(bw_version)"
}

# =============================================================================

export bw_removeParams=(
  --yes/y --completely/c
  "${_verbosityParams[@]}"
)
export bw_remove_yes_description='подтверждает удаление'
export bw_remove_completely_description='удаляет не только все связанное с $_bwFileName, но и сам bw.bash'
export bw_remove_silent_description='подавляет сообщение об успешном удалении'
export bw_removeShortcuts=( 'rm' )
export bw_remove_description="удаляет $_bwFileName и все связанное с ним"
export bw_removeCondition='[[ -z $_isBwDevelopInherited ]]'
bw_remove() { eval "$_funcParams2"
  if [[ -z $yes ]]; then
    _warn "Чтобы произвести удаление, запустите эту команду с опцией ${_ansiCmd}--yes${_ansiWarn} или ${_ansiCmd}-y"
    return 2
  else
    codeHolder=_codeToInitSubOPT eval "$_evalCode"

    if [[ -z $_isBwDevelop ]]; then
      local exactLine='. ~/bw.bash'
      # local matchRegexp='^[ \t]*\.[ \t]+\"?(.+?)\/bw\.bash'
      [[ $verbosity == dry ]] || _setAtBashProfile -u "$exactLine" "$_bwMatchRegexp"
      # _profileUnless "${sub_OPT[@]}" -u ". $(_quotedArgs "$(_shortenFileSpec "$_bwFileSpec")")"
    fi

    local -a requiredFiles; mapfile -t requiredFiles < <(find "$_bwDir/core" -name '*.bash' -exec basename {} .bash \;)
    requiredFiles+=( bw )

    local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
      [[ -d $dirSpec ]] || continue
      for _fileSpec in "$dirSpec/"*"$_unsetFileExt"; do
        if _hasItem "$(basename "$_fileSpec" "$_unsetFileExt")" "${requiredFiles[@]}"; then
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

    local msg=''
    if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
      if [[ -n $_isBwDevelop ]]; then
        msg+="${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все прегенеренные вспомогательные файлы. Все основные ${_ansiFileSpec}*.bash${_ansiWarn}-файлы оставлены нетронутыми${_ansiReset}"
      else
        local suffix;
        if [[ -z $completely ]]; then
          suffix=", кроме ${_ansiFileSpec}$_bwFileSpec${_ansiWarn}. Для повторной установки выполните команду ${_ansiCmd}. $_bwFileSpec${_ansiWarn}"
        else
          suffix=", включая ${_ansiFileSpec}$_bwFileSpec${_ansiWarn}"
        fi
        msg+="${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все связанное с ней (содержимое директории ${_ansiFileSpec}$_bwDir${_ansiWarn})$suffix${_ansiReset}"
      fi
    fi

    local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
      if [[ $verbosity == dry ]]; then
        echo "${_ansiCmd}rm -rf \"$dirSpec\"${_ansiReset}"
      else
        rm -rf "$dirSpec"
      fi
    done

    # local -a varNamesToUnset=( $( compgen -v | perl -ne "print if /^(_specialVars|$_specialVars)\$/" ) )
    local -a varNamesToUnset; mapfile -t varNamesToUnset < <(compgen -v | perl -ne "print if /^(_specialVars|$_specialVars)\$/")
    for _fileSpec in "${requiredFiles[@]}"; do
      _fileSpec="/tmp/$_fileSpec$_unsetFileExt"
      if [[ $verbosity == dry ]]; then
        echo "${_ansiCmd}. \"$_fileSpec\"${_ansiReset}"
      else
        codeHolder=_codeSource eval "$_evalCode"
      fi
    done
    if [[ $verbosity == dry ]]; then
      echo "${_ansiCmd}unset ${varNamesToUnset[*]}${_ansiReset}"
    else
      unset "${varNamesToUnset[@]}"
    fi

    echo "$msg"
  fi
}

# =============================================================================

export _testPlurals='тест теста тестов'
export _funcPlurals='функция функции функций'

export bw_bashTests_list_description='
  Выводит список bash-функций, для которых существуют тесты,
  с указанием количества тестов для каждой функции
'
export bw_bashTests_noTiming_description='Отключает замер времени работы тестов'
export bw_bashTests_args_name=Аргумент
export bw_bashTests_args_description='!
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
export bw_bashTestsShortcuts=( 'bt' )
export bw_bashTestsParams=( '--noTiming' '--list' '@..args' )
export bw_bashTestsParamsOpt=( --canBeMixedOptionsAndArgs )
export bw_bashTestsSuppressInheritedOptionVarNames=( 'verbosity' 'silent' )
export bw_bashTests_description='запускает тесты bash-функций'
bw_bashTests() {
  eval "$_funcParams2"
  local testsDirSpec="$_bwDir/tests"
  if [[ -z $_isBwDevelop && -z $_isBwDevelopInherited ]]; then
    local testsSupportFileSpec="$testsDirSpec/testsSupport.bash"
    # if [[ ! -f $testsSupportFileSpec || $_bwFileSpec -nt $testsSupportFileSpec ]]; then
    _getBwTar "$_bwFileSpec" tests | tar xf - -C "$_bwDir" \
      || { _err "Не удалось извлечь архив tests из ${_ansiFileSpec}$_bwFileSpec${_ansiErr} в ${_ansiFileSpec}$testsDirSpec"; return $?; }
    # fi
    _fileSpec="$testsSupportFileSpec" codeHolder=_codeSource eval "$_evalCode"
  fi
  _runInBackground bw_bashTestsHelper
}

# shellcheck disable=SC2034
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
  local funcTestFor; for funcTestFor in "${allFuncsWithTests[@]}"; do
    local conditionHolder="${funcTestFor}TestsCondition"
    [[ -z ${!conditionHolder} ]] || eval "${!conditionHolder}" || continue
    funcsWithTests+=( "$funcTestFor" )
  done
}
_prepareAllFuncWithTests() {
  if [[ -n $_isBwDevelop ]]; then
    _rmCache
  else
    varName=allFuncsWithTests codeHolder=_codeToUseCache eval "$_evalCode"
  fi
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

export bw_setParams=(
  '!--uninstall/u'
)
export bw_setParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
export bw_set_cmd_name=Имя-настройки
export bw_set_description='включает/отключает настройку'
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
export _ps_user_defaultAnsi='Green'
export _ps_ptcd_defaultAnsi='White'
export _ps_cd_defaultAnsi='White'

_preparePromptParams=()
_preparePromptParams() {
  varName=_preparePromptParams codeHolder=_codeToUseCache eval "$_evalCode"
  # local -a _psItems=( $( compgen -v | perl -ne "print if s/^_ps_([^_]+)\$/\$1/" ) )
  local -a _psItems; mapfile -t _psItems < <( compgen -v | perl -ne "print if s/^_ps_([^_]+)\$/\$1/" )
  local -a _preparePromptParams_itemsDefaults=( user ptcd git error separator )
  local -a _preparePromptParams_itemsEnum=( "${_preparePromptParams_itemsDefaults[@]}" )
  local -a _preparePromptParamsAddon=()
  local _psItem; for _psItem in "${_psItems[@]}"; do
    local defaultAnsiHolder="_ps_${_psItem}_defaultAnsi"
    local defaultValueHolder="_ps_${_psItem}"
    local defaultAnsi="Reset"
    if [[ -n ${!defaultAnsiHolder} ]]; then
      defaultAnsi="${!defaultAnsiHolder}"
    fi
    _preparePromptParamsAddon+=(
      "--${_psItem}=$(_quotedArgs --strip "${!defaultValueHolder}")"
      "--${_psItem}Space:${_enumSpace}=after"
      "@1..--${_psItem}Ansi:${_enumAnsi}=( $defaultAnsi )"
    )
    if ! _hasItem "$_psItem" "${_preparePromptParams_itemsEnum[@]}"; then
      _preparePromptParams_itemsEnum+=( "$_psItem" )
    fi
  done

  _preparePromptParams=(
    "@1..--items:( ${_preparePromptParams_itemsEnum[*]} )=( ${_preparePromptParams_itemsDefaults[*]} )"

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

# shellcheck disable=SC2034
_preparePrompt() {
  eval "$_funcParams2"
  local optVarName; for optVarName in "${__optVarNames[@]}"; do
    if [[ ${optVarName:$(( ${#optVarName} - 4 ))} == Ansi  ]]; then
      local item=${optVarName:0:$(( ${#optVarName} - 4 ))}
      dstVarName=ansiSrc srcVarName=${item}Ansi eval "$_codeToInitLocalCopyOfArray"
      local ansiAsStringHolder=${item}AnsiAsString
      eval local "$ansiAsStringHolder"=
      local ansi; for ansi in "${ansiSrc[@]}"; do
        local ansiHolder="_ansi${ansi}"
        eval "$ansiAsStringHolder"+=\"\$\{!ansiHolder\}\"
      done
    fi
  done

  prompt=
  local -a groups=(error git)
  local -a realItems=()
  local item; for item in "${items[@]}"; do
    local group; for group in "${groups[@]}"; do
      if [[ ${!group} != off && ${!group} != - && $item == "$group" ]]; then
        eval local "$group"=true
      fi
    done
    if [[ $item == error ]]; then
      realItems+=( errorPrefix errorInfix errorCode errorSuffix )
    elif [[ $item == git ]]; then
      realItems+=( gitPrefix gitBranch gitDirty gitSuffix )
    else
      realItems+=( "$item" )
    fi
  done
  local group; for group in "${groups[@]}"; do
    if [[ ${!group} == true ]]; then
      prompt+='`_psPrepare_'$group'`'
    fi
  done
  local item; for item in "${realItems[@]}"; do
    if ! [[ ${!item} == - || ${!item} == off || ${!item} == no ]]; then
      local spaceHolder="${item}Space"
      local foundGroup='' group; for group in "${groups[@]}"; do
        if [[ ${item:0:${#group}} == "$group" ]]; then
          foundGroup=$group
          break
        fi
      done
      local psFuncName="_ps_$item"
      local promptItem=
      if _funcExists "$psFuncName"; then
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

export bw_set_promptParams=()
bw_set_promptParams() {
  _preparePromptParams
  bw_set_promptParams=( "${_preparePromptParams[@]}" )
}
export bw_set_prompt_description='Настраивает prompt'
bw_set_prompt() { eval "$_funcParams2"
  local -a supportModules=( git ansi ps )
  local supportFileNameSuffix='Support.bash'
  local newFileSpec="$_profileFileSpec.new"

  if [[ -n $uninstall ]]; then
    local -a fileNames=()
    local moduleName; for moduleName in "${supportModules[@]}"; do
      fileNames+=( "${moduleName}${supportFileNameSuffix}" )
    done
    local -a fileSpecs=()
    local fileName; for fileName in "${fileNames[@]}"; do
      local fileSpec="$HOME/$fileName"
      [[ ! -f $fileSpec ]] || fileSpecs+=( "$fileSpec" )
    done
    [[ ${#fileSpecs[@]} -eq 0 ]] || rm "${fileSpecs[@]}"
    [[ -z $OLD_PS1 ]] || PS1="$OLD_PS1"
    _exportVarAtBashProfile --uninstall OLD_PS1
    _exportVarAtBashProfile --uninstall PS1

    # if grep -E "^(export )?(OLD_)?PS1=" "$_profileFileSpec" >/dev/null 2>&1; then
    #   grep -v -E "^(export )?(OLD_)?PS1=" "$_profileFileSpec" >"$newFileSpec"
    #   _backupBashProfileFile
    #   mv "$newFileSpec" "$_profileFileSpec"
    # fi
  else

    local -a subOPT=(); local optVarName; for optVarName in "${__OPTVarNames[@]}"; do
      [[ $optVarName != OPT_uninstall && $optVarName != OPT_help ]] || continue
      dstVarName=OPT srcVarName=$optVarName eval "$_codeToInitLocalCopyOfArray"
      subOPT+=( "${OPT[@]}" )
    done
    local prompt; _preparePrompt "${subOPT[@]}" || return $?
    # local exactLine="export PS1=\\'$prompt\\'"
    # local moduleName; for moduleName in "${supportModules[@]}"; do
    #   exactLine+="; . ~/${moduleName}${supportFileNameSuffix}"
    # done

    # local moduleName; for moduleName in "${supportModules[@]}"; do
    #   local fileName="${moduleName}${supportFileNameSuffix}"
    #   local srcFileSpec="$_bwDir/core/$fileName"
    #   local dstFileSpec="$HOME/$fileName"
    #   if [[ -f "$srcFileSpec" ]]; then
    #     cp "$srcFileSpec" "$dstFileSpec"
    #   else
    #     echo '[[ $(type -t _resetBash) != function ]] || _resetBash' >"$dstFileSpec"
    #     _getFileChunk "$_bwFileSpec" "# ==${moduleName} start" "# ==${moduleName} end" >>"$dstFileSpec"
    #   fi
    # done
    if _hasExportVarAtBashProfile -n OLD_PS1; then
      export OLD_PS1="$PS1"
      _exportVarAtBashProfile OLD_PS1
    fi
    export PS1="$prompt"
    _exportVarAtBashProfile PS1
    # if ! grep -x "$exactLine" "$_profileFileSpec" >/dev/null 2>&1; then
    #   local reOld=
    #   local profileLineOld=
    #   if [[ -z $OLD_PS1 ]]; then
    #     export OLD_PS1="$PS1"
    #     profileLineOld="export OLD_PS1=\"$PS1\""
    #     reOld='(OLD_)?'
    #   fi
    #   grep -v -E "^(export )?${reOld}PS1=" "$_profileFileSpec" >"$newFileSpec"
    #   [[ -z $profileLineOld ]] || echo "$profileLineOld" >>"$newFileSpec"
    #   echo "$exactLine" >>"$newFileSpec"
    #   _backupBashProfileFile
    #   mv "$newFileSpec" "$_profileFileSpec"

    #   export PS1="$prompt"
    # fi
  fi
}

export _codeToPrepareDescriptionsOf_bw_set_prompt='
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

export bw_set_viParams=()
export bw_set_vi_description='Включает vi-режим для readline'
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
      --vi-cmd-mode-string "\\1${_ansiDarkGray}\\2CMD \\1${_ansiReset}\\2"
      --vi-ins-mode-string "\\1${_ansiMagenta}${_ansiBold}\\2INS \\1${_ansiReset}\\2"
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
      if [[ -z $uninstall ]]; then
        msg+="$didInstall"
      else
        msg+="$didUninstall"
      fi
      _ok "$msg"
    fi
  else
    if [[ $verbosity =~ ^(err|all.*)$ ]]; then
      msg="Не удалось "
      if [[ -z $uninstall ]]; then
        msg+="$toInstall"
      else
        msg+="$toUninstall"
      fi
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
  titlePrefix='' didInstall="установлен" didUninstall="удален" toInstall="установить" toUninstall="полностью удалить" _showResult
}

_showRunResult() {
  titlePrefix='' didInstall="запущен" didUninstall="остановлен" toInstall="запустить" toUninstall="остановить" _showResult
}

# =============================================================================

_bwProjDefs=(
  # 'mls' '
  #   --gitOrigin github.com:baza-winner/mls-pm.git
  # '
  'bwdev' '
    --gitOrigin github.com:baza-winner/bw.git
    --branch master
    --http 8998
    --https 8999
    --no-docker-build
    --docker-compose "$_bwDir/docker/docker-compose.nginx.yml"
  '
  'dubu' '
    --gitOrigin github.com:baza-winner/dev-ubuntu.git
    --branch master
    --docker-image-name bazawinner/dev-ubuntu
    --http 8996
    --https 8997
    --docker-compose "$_bwDir/docker/docker-compose.nginx.yml"
    --docker-compose "$_bwDir/docker/docker-compose.main.yml"
  '
  'dnvm' '
    --gitOrigin github.com:baza-winner/dev-nvm.git
    --branch master
    --docker-image-name bazawinner/dev-nvm
    --http 8994
    --https 8995
    --docker-compose "$_bwDir/docker/docker-compose.nginx.yml"
    --docker-compose "$_bwDir/docker/docker-compose.main.yml"
  '
  'bgate' '
    --gitOrigin github.com:baza-winner/billing-gate.git
    --branch develop
    --http 8086
    --https 8087
    --upstream 6783
    --docker-compose "$_bwDir/docker/docker-compose.nginx.yml"
    --docker-compose "$_bwDir/docker/docker-compose.main.yml"
  '
  'crm' '
    --gitOrigin github.com:baza-winner/crm.git
    --branch develop
    --http 8088
    --https 8089
    --upstream 3000
    --docker-compose "$_bwDir/docker/docker-compose.nginx.yml"
    --docker-compose "$_bwDir/docker/docker-compose.main.yml"
  '
  # 'wcraw' '
  #   --gitOrigin github.com:baza-winner/billing-gate.git
  #   --branch feature/docker
  #   --http 8086
  #   --https 8087
  # '
  # 'bcore' '
  #   --gitOrigin github.com:baza-winner/billingcore.git
  # '
  # 'crm' '
  #   --gitOrigin git@github.com:baza-winner/crm.git
  # '
)

_bw_project_bgate() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}
_bw_project_bcore() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}
_bw_project_dnvm() {
  _exec "${sub_OPT[@]}" git submodule update --remote --recursive # https://stackoverflow.com/questions/1777854/git-submodules-specify-a-branch-tag/15782629#15782629
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
    || { _throw "ожидает, что переменная ${_ansiOutline}bwProjShortcut${_ansiErr} будет иметь непустое значение"; return $?; }
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    if [[ $bwProjShortcut == "${_bwProjDefs[$i]}" ]]; then
      codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
      break
    fi
  done
  [[ -n $bwProjGitOrigin ]] \
    || { _throw "not found gitOrigin for bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle"; return $?; }
}

export _codeToCallPrepareBwProjVarsHelper='
  bwProjShortcut="${_bwProjDefs[$i]}"
  local codeToGetBwProjDef
  eval local -a params=\( ${_bwProjDefs[$((i + 1))]} \)
  _prepareBwProjVarsHelper "${params[@]}" || return $?
  codeHolder=codeToGetBwProjDef eval "$_evalCode"
  bwProjName=$(basename "$bwProjGitOrigin" .git)
  if [[ $bwProjName == "$bwProjShortcut" ]]; then
    bwProjTitle="$bwProjName"
  else
    bwProjTitle="$bwProjName ($bwProjShortcut)"
  fi
'

export _tcpPortDiap='1024..65535'
export _prepareBwProjVarsHelperParams=(
  '--gitOrigin='
  "--http:$_tcpPortDiap"
  "--https:$_tcpPortDiap"
  "--upstream:$_tcpPortDiap"
  '--branch='
  '--docker-image-name='
  '--no-docker-build'
  '@--docker-compose'
)
export _codeToDeclareLocalBwProjVars='
  local bwProjName=""
  local bwProjTitle=""
  local bwProjGitOrigin="" 
  local bwProjDefaultHttp=""
  local bwProjDefaultHttps="" 
  local bwProjDefaultUpstream=""
  local bwProjDefaultBranch="" 
  local bwProjDockerImageName="" 
  local bwProjNoDockerBuild=""
  local -a bwProjDockerCompose=()
'

# shellcheck disable=SC2034
_prepareBwProjVarsHelper() { eval "$_funcParams2"
  codeToGetBwProjDef='
    bwProjGitOrigin="'"$gitOrigin"'"
    bwProjDefaultHttp='"$http"'
    bwProjDefaultHttps='"$https"'
    bwProjDefaultUpstream="'"$upstream"'"
    bwProjDefaultBranch="'"$branch"'"
    bwProjDockerImageName="'"$dockerImageName"'"
    bwProjNoDockerBuild="'"$noDockerBuild"'"
    bwProjDockerCompose=( '$(_quotedArgs "${dockerCompose[@]}")' )
  '
}

export _codeToPrepareDescriptionsOf_bw_project='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_project_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'

export bw_projectParams=()
bw_projectParams() {
  varName="${FUNCNAME[0]}" codeHolder=_codeToUseCache eval "$_evalCode"
  local -a bwProjShortcuts=()
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    if [[ ! $bwProjShortcut =~ $_isValidVarNameRegExp ]]; then
      _throw "bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle${_ansiErr} $_mustBeValidVarName"; return $?
    elif _hasItem "$bwProjShortcut" "${bwProjShortcuts[@]}"; then
      _throw "Duplicate bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle"; return $?
    else
      bwProjShortcuts+=( "$bwProjShortcut" )
    fi
  done
  local bwProjShortcutsAsString="${bwProjShortcuts[*]}"
  verbosityDefault=allBrief silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
  export bw_projectParams=(
    '!--uninstall/u'
    '!--force/f'
    '--branch='
    "${_verbosityParams[@]}"
    "bwProjShortcut:( $bwProjShortcutsAsString )"
    'projDir:?'
  )
  _saveToCache "${FUNCNAME[0]}"
}

export bwProjGlobalDefaultBranch=develop
export bw_projectParamsOpt=(--canBeMixedOptionsAndArgs)
export _uninstall_description='Включить режим удаления'
export _force_description='Игнорировать факт предыдущей установки'
export bw_project_bwProjShortcut_name='Сокращенное-имя-проекта'
export bw_project_projDir_name='Папка-проекта'
export bw_project_projDir_description='
  Папка, куда будет установлен проект
  По умолчанию, в качестве папки проекта используется ${_ansiCmd}~/${_ansiOutline}Полное-имя-проекта${_ansiReset}
  ${_ansiOutline}Полное-имя-проекта${_ansiReset} - имя проекта на github'\''е
'
export bw_project_branch_description='Ветка, на которую следует переключиться после установки проекта'
export bw_project_description='Разворачивает/удаляет проект'
export bw_projectShortcuts=( 'p' )
bw_project() { eval "$_funcParams2"
  codeHolder=_codeToInitSubOPT eval "$_evalCode"

  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  [[ -n $branch ]] || branch=${bwProjDefaultBranch:-$bwProjGlobalDefaultBranch}

  local alreadyProjDir='';  _getAlreadyProjDir "$bwProjShortcut" --varName alreadyProjDir
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
      local cmdFileSpec="$alreadyProjDir/bin/$bwProjShortcut.bash"
      if [[ -f "$cmdFileSpec" ]]; then
        fileSpec="$cmdFileSpec" _unsetBash "${sub_OPT_verbosity[1]}"
      fi
    fi

    bw_install git --silentIfAlreadyInstalled || { returnCode=$?; break; }

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
        local cloneStderrFileSpec="/tmp/$bwProjShortcut.clone.stderr"
        if ! git ls-remote -t "git@$bwProjGitOrigin" no-ref >"$cloneStderrFileSpec" 2>&1; then
          if grep 'Permission denied (publickey)' "$cloneStderrFileSpec"; then
            local msg=
            msg+="Похоже, Вы не настроили ssh-ключи для доступа к ${_ansiPrimaryLiteral}git@$bwProjGitOrigin${_ansiWarn}"$_nl
            msg+="Воспользуйтесь следующей командой:"$_nl
            msg+="    ${_ansiCmd}bw github-keygen ${_ansiOutline}Имя-пользователя-на-github${_ansiWarn}"
            _warn "$msg"
          else
            cat "$cloneStderrFileSpec"
          fi
          rm -f "$cloneStderrFileSpec"
          { returnCode=1; break; }
        fi
        _exec "${sub_OPT[@]}" --cmdAsIs "git clone git@$bwProjGitOrigin \"$projDir\"" || { returnCode=$?; break; }
        _exec "${sub_OPT[@]}" cd "$projDir" || { returnCode=$?; break; }
        _exec "${sub_OPT[@]}" git checkout "$branch" || { returnCode=$?; break; }
        local funcName="_${FUNCNAME[0]}_$bwProjShortcut"
        ! _funcExists "$funcName" || $funcName || { returnCode=$?; break; }
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
          [[ -z $needChangePwd ]] || _exec "${sub_OPT[@]}" cd "$HOME" \
            || { returnCode=$?; break; }
          fi
      fi
      break
    done; [[ $returnCode -eq 0 ]] || break

    break
  done

  folder="$projDir" name="$bwProjTitle" _showProjectResult

  if [[ $returnCode -eq 0 && $verbosity != dry ]]; then
    local cmdFileSpec="$projDir/bin/${bwProjShortcut}.bash"
    local exactLine; exactLine=". $(_quotedArgs $(_shortenFileSpec "$cmdFileSpec")); $bwProjShortcut update completionOnly"
    # local matchRegxp="^[[:space:]]*\\.\\s+\"?(.+?)\\/bin\\/$bwProjShortcut\\.bash"
    # local matchRegxp='^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~].+)\/bin\/'"$bwProjShortcut"'\.bash'
    local matchRegxp="$_sourceMatchRegexp"'bin\/'"$bwProjShortcut"'\.bash'
    if [[ -n $uninstall ]]; then
        _setAtBashProfile -u "$exactLine" "$matchRegxp"
    else
      if [[ ! -f "$cmdFileSpec" ]]; then
        local msg=
        msg+="Не найден файл ${_ansiFileSpec}bin/$bwProjShortcut.bash${_ansiErr}$_nl"
        msg+="Не удалось инициализировать команду ${_ansiCmd}$bwProjShortcut"
        [[ $verbosity == none  ]] || _err "$msg"
        returnCode=1
      elif ! _exec "${sub_OPT[@]}" . "$cmdFileSpec"; then
        local msg=
        msg+="Не удалось инициализировать команду ${_ansiCmd}$bwProjShortcut"
        [[ $verbosity == none  ]] || _err "$msg"
        returnCode=1
      else
        _setAtBashProfile "${OPT_uninstall[@]}" "$exactLine" "$matchRegxp"
        "$bwProjShortcut" update completionOnly
        [[ $verbosity == none  ]] || echo "${_ansiWarn}Теперь доступна команда ${_ansiCmd}$bwProjShortcut${_ansiReset}"
        _exec "${sub_OPT[@]}" --treatAsOK 3 "$bwProjShortcut" -?
      fi
    fi
  fi

  return $returnCode
}

export _prepareGitDirtyParams=( 'originSuffix' )
_prepareGitDirty() { eval "$_funcParams2"
  local returnCode=0
  gitDirty=
  local gitOrigin; gitOrigin=$(_gitOrigin)
  if [[ ${#gitOrigin} -ge ${#originSuffix} && $originSuffix == "${gitOrigin:$(( ${#gitOrigin} - ${#originSuffix} ))}" ]]; then
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

export _bwDevDockerEntryPointFileSpec="/home/dev/.bw/docker/entrypoint.bash"
export _bwDevDockerBwFileSpec="/home/dev/bw.bash"
export _bwDevDockerBwDir="/home/dev/.bw"
export _bwDevDockerProjDir="/home/dev/proj"
export _bwSslFileSpecPrefix="$_bwDir/ssl/server."
export _bwNginxConfDir="$_bwDir/docker/nginx/conf.bw"

# shellcheck disable=SC2034
_initBwProjCmd() {
  local fileSpec; fileSpec=$(_getSelfFileSpec 2)
  local bwProjShortcut; bwProjShortcut=$(basename "$fileSpec" .bash)
  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  eval "_$bwProjShortcut"'FileSpec="$fileSpec"'
  local bwProjDirHolder="_${bwProjShortcut}Dir"
  _realpath -v "$bwProjDirHolder" "$(dirname "$fileSpec")/.."
  eval export "_$bwProjShortcut"'DockerImageName="${bwProjDockerImageName:-bazawinner/dev-$bwProjShortcut}"'
  eval export "_$bwProjShortcut"'DockerImageIdFileSpec="${!bwProjDirHolder}/docker/dev-$bwProjShortcut.id"'
  eval export "_$bwProjShortcut"'DockerContainerName="dev-$bwProjShortcut"'

  eval "$bwProjShortcut"'_description='\''Базовая утилита проекта ${_ansiPrimaryLiteral}'"$bwProjName"' ${_ansiUrl}$bwProjGitOrigin${_ansiReset}'\'
  eval "$bwProjShortcut"'Params=()'
  eval "$bwProjShortcut"'ParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)'
  eval "$bwProjShortcut()"' { eval "$_funcParams2" 
  }'

  eval "$bwProjShortcut"'_updateParams=( 
    "mode:(completionOnly sourceWithoutPregen pregenSourced)=pregenSourced" 
  )'
  eval "$bwProjShortcut"'_update_mode_completionOnly_description='\''Только обновить completion-определения'\'
  eval "$bwProjShortcut"'_update_mode_sourceWithoutPregen_description='\''Перечитать исходные файлы (без прегенерации) перед обновлением completion-определения'\'
  eval "$bwProjShortcut"'_update_mode_full_description='\''Перечитать исходные файлы перед обновлением completion-определения'\'
  eval "$bwProjShortcut"'_update_description='\''Обновляет команду ${_ansiCmd}'"$bwProjShortcut"'${_ansiReset}'\'
  eval "$bwProjShortcut"'_update() { eval "$_funcParams2"
    _cmd_update "$mode" "'"$bwProjShortcut"'"
  }'

  eval "$bwProjShortcut"'_dockerParams=()'
  eval "$bwProjShortcut"'_dockerParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)'
  eval "$bwProjShortcut"'_docker_description='\''docker-операции'\'
  eval "$bwProjShortcut"'_dockerCondition='\''! _isInDocker'\'
  eval "$bwProjShortcut"'_docker() { eval "$_funcParams2"; }'

  local codeToPrepareOPTForDockerUp=''
  local -a dockerUpParams=()

  local codeToPrepareOPTForSelfTest=''
  local -a selfTestParams=()

  if [[ -n $bwProjNoDockerBuild ]]; then
    unset "$bwProjShortcut"_docker_buildParams "$bwProjShortcut"_docker_build_description
    unset -f "$bwProjShortcut"_docker_build
    unset "$bwProjShortcut"_docker_pushParams "$bwProjShortcut"_docker_push_description
    unset -f "$bwProjShortcut"'_docker_push'
    unset "$bwProjShortcut"_docker_up_noCheck_description
  else
    dockerUpParams+=( '--no-check/n' )
    local dockerImageTitle='${_ansiPrimaryLiteral}$_'$bwProjShortcut'DockerImageName${_ansiReset}'
    eval "$bwProjShortcut"'_docker_up_noCheck_description='\''Не проверять актуальность docker-образа '\''"$dockerImageTitle"'
    codeToPrepareOPTForDockerUp+='
      OPT+=( 
        --dockerImageName "$_'$bwProjShortcut'DockerImageName"
        --dockerImageIdFileSpec "$_'$bwProjShortcut'DockerImageIdFileSpec"
      )
    '
    codeToPrepareOPTForSelfTest+='
      OPT+=( 
        --dockerImageName "$_'$bwProjShortcut'DockerImageName"
      )
    '
    eval "$bwProjShortcut"'_docker_buildParams=()'
    eval "$bwProjShortcut"'_docker_build_description='\''Собирает docker-образ '\''"$dockerImageTitle"'
    eval "$bwProjShortcut"'_docker_build() { eval "$_funcParams2"
      local -a params=(
        "'"$bwProjShortcut"'"
        "$_'"$bwProjShortcut"'Dir"
        "$_'"$bwProjShortcut"'DockerImageName"
        "$_'"$bwProjShortcut"'DockerImageIdFileSpec"
        "'"$dockerImageTitle"'"
      )
      _docker_build "${params[@]}"
    }'

    eval "$bwProjShortcut"'_docker_pushParams=()'
    eval "$bwProjShortcut"'_docker_push_description='\''Push-ит docker-образ '\''"$dockerImageTitle"'
    eval "$bwProjShortcut"'_docker_push() { eval "$_funcParams2"
      docker push "$_'"$bwProjShortcut"'DockerImageName"
    }'
  fi

  _prepareVarsForDefaultPort http
  _prepareVarsForDefaultPort https
  _prepareVarsForDefaultPort upstream 'для сервиса ${_ansiPrimaryLiteral}nginx${_ansiReset}'

  eval "$bwProjShortcut"'_selfTestParamsOpt=(
    --additionalDependencies "$_'"$bwProjShortcut"'Dir/bin/'"$bwProjShortcut"'.bash"
  )'
  eval "$bwProjShortcut"'_selfTestParams=( 
    "${selfTestParams[@]}"
  )'
  eval "$bwProjShortcut"'_selfTestCondition='\''! _isInDocker'\'
  eval "$bwProjShortcut"'_selfTestShortcuts=( st )'
  eval "$bwProjShortcut"'_selfTest_description='\''Самопроверка'\'
  eval "$bwProjShortcut"'_selfTest() { eval "$_funcParams2"
    local -a OPT=()
    '"$codeToPrepareOPTForSelfTest"'
    _cmd_selfTest "${OPT[@]}" "'"$bwProjShortcut"'" "'"${!bwProjDirHolder}"'"
  }'

  local dockerComposeFileSpec="$_${bwProjShortcut}Dir/docker/docker-compose.yml"

  eval "$bwProjShortcut"'_docker_upParamsOpt=(
    --additionalDependencies "$_'"$bwProjShortcut"'Dir/bin/'"$bwProjShortcut"'.bash"
  )'

  eval "$bwProjShortcut"'_docker_upParams=( 
    "${dockerUpParams[@]}"
    '\''@--restart/r:( $(_'"$bwProjShortcut"'_dockerServiceNames) )'\''
    --force-recreate/f
    "${'"$bwProjShortcut"'_docker_upParamsAddon[@]}"
  )'
  eval "$bwProjShortcut"'_docker_up_restart_description='\''Останавливает и поднимает указанные сервисы'\'
  eval "$bwProjShortcut"'_docker_up_forceRecreate_description='\''Поднимает ${_ansiPrimaryLiteral}$(_'"$bwProjShortcut"'_dockerContainerNames)${_ansiReset} с опцией ${_ansiCmd}--force-recreate${_ansiReset}, передаваемой ${_ansiCmd}docker-compose up${_ansiReset}'\'
  eval "$bwProjShortcut"'_docker_up_description='\''Поднимает (up) ${_ansiPrimaryLiteral}$(_'"$bwProjShortcut"'_dockerContainerNames)${_ansiReset}'\'
  eval "$bwProjShortcut"'_docker_up() { eval "$_funcParams2"
    local -a OPT=( 
      "${OPT_noCheck[@]}" 
      "${OPT_forceRecreate[@]}" 
      "${OPT_restart[@]}" 
      --dockerContainerName "$_'"$bwProjShortcut"'DockerContainerName"
      --bwProjShortcut "'"$bwProjShortcut"'"
      --bwProjName "'"$bwProjName"'"
    )

    '"$codeToPrepareOPTForDockerUp"'

    _inDir "$_'"$bwProjShortcut"'Dir/docker" _docker_up "${OPT[@]}" '"$(_quotedArgs "${bwProjDockerCompose[@]}")"' "$_'"$bwProjShortcut"'Dir/docker/docker-compose.yml"
  }'

  eval "$bwProjShortcut"'_docker_down_description='\''Останавливает (${_ansiCmd}docker-compose down${_ansiReset}) следующие контейнеры: ${_ansiPrimaryLiteral}$(_'"$bwProjShortcut"'_dockerContainerNames)${_ansiReset}'\'
  eval "$bwProjShortcut"'_docker_downParams=()'
  eval "$bwProjShortcut"'_docker_down() { eval "$_funcParams2"
    _docker_down '"$(_quotedArgs "${bwProjDockerCompose[@]}")"' "$_'"$bwProjShortcut"'Dir/docker/docker-compose.yml"
  }'

  eval "_$bwProjShortcut"'_dockerContainerNames() {
    _dockerComposeContainerNames "$_'"$bwProjShortcut"'Dir/docker/docker-compose.yml"
  }'
  eval "_$bwProjShortcut"'_dockerServiceNames() {
   _dockerComposeServiceNames "$_'"$bwProjShortcut"'Dir/docker/docker-compose.yml"
  }'
  eval "$bwProjShortcut"'_docker_shellParams=( 
    '\''containerName:( $(_'"$bwProjShortcut"'_dockerContainerNames) )='\''"$_'"$bwProjShortcut"'DockerContainerName"
    '\''shell=$(_getDefaultShellOfDockerContainer "$containerName")'\''
  )'
  eval "$bwProjShortcut"'_docker_shell_description='\''Запускает bash в docker-контейнере'\'
  eval "$bwProjShortcut"'_docker_shell_containerName_name='\''Имя-контейнера'\'
  eval "$bwProjShortcut"'_docker_shell() { eval "$_funcParams2"
    _docker_shell "$containerName" "$_'"$bwProjShortcut"'DockerContainerName" "$shell"
  }'
}

_prepareVarsForDefaultPortParams=(
  'portName'
  'descriptionSuffix:?'
)
_prepareVarsForDefaultPort() { eval "$_funcParams2"
  local defaultHolder="bwProjDefault$(_upperFirst "$portName")"
  if [[ -z ${!defaultHolder} ]]; then
    unset "${bwProjShortcut}_docker_up_${portName}_description"
    unset "${bwProjShortcut}_selfTest_${portName}_description"
  else
    local paramDef="--$portName:$_tcpPortDiap=${!defaultHolder}"
    local description="$portName-порт по которому будет доступен контейнер" 
    if [[ -n $descriptionSuffix ]]; then
      description+=" $descriptionSuffix"
    fi
    local code='
      OPT+=( ${OPT_'"$portName"'[@]} )
    '
    dockerUpParams+=( "$paramDef" )
    selfTestParams+=( "$paramDef" )
    eval "${bwProjShortcut}_docker_up_${portName}_description"=\'"$description"\'
    eval "${bwProjShortcut}_selfTest_${portName}_description"=\'"$description"\'
    codeToPrepareOPTForDockerUp+="$code"
    codeToPrepareOPTForSelfTest+="$code"
  fi
}

_getDefaultShellOfDockerContainer() {
  if [[ $1 =~ nginx$ ]]; then
    echo '/bin/sh'
  else
    echo '/bin/bash'
  fi
}

export _cmd_updateParams=( 
  'mode:(completionOnly sourceWithoutPregen pregenSourced)'
  'bwProjShortcut' 
)
_cmd_update() { eval "$_funcParams2"
  local sourceFileSpecHolder
  if _isInDocker; then
    sourceFileSpecHolder='_bwDevDockerEntryPointFileSpec'
  else
    sourceFileSpecHolder="_${bwProjShortcut}FileSpec"
  fi
  if [[ $mode != completionOnly ]]; then
    local -a OPT=()
    if [[ $mode != pregenSourced ]]; then
      OPT=( -p - )
    fi
    . "$_bwFileSpec" "${OPT[@]}" || return $?
    . "${!sourceFileSpecHolder}" || return $?
    rm -f "$_bwDir/generated/$bwProjShortcut"*
  fi
  local -a __completions=()
  local -a funcNamesToPregen; mapfile -t funcNamesToPregen < <(compgen -c "$bwProjShortcut")
  _pregen "${funcNamesToPregen[@]}" || return $?
  for fileSpec in "${__completions[@]}"; do
    . "$fileSpec"
  done
}

export _docker_buildParams=( 
  'bwProjShortcut' 
  'bwProjDir' 
  'dockerImageName'
  'dockerImageIdFileSpec'
  'dockerImageTitle'
)
_docker_build() { eval "$_funcParams2"
  if _docker -v all build --pull -t "$dockerImageName" "$bwProjDir/docker"; then
    _docker inspect --format "{{json .Id}}" "$dockerImageName:latest" > "$dockerImageIdFileSpec"
    if (git update-index -q --refresh && git diff-index --name-only HEAD -- | grep "$_bwFileName" >/dev/null 2>&1); then
      local msg=
      msg+="Обновлен docker-образ $dockerImageTitle${_ansiWarn}"$_nl
      msg+="${_ansiWarn}Не забудьте поместить его в docker-репозиторий командой"$_nl
      msg+="    ${_ansiCmd}$bwProjShortcut docker push${_ansiReset}"
      _warn "$msg"
    fi
  fi
}

export _docker_downParams=( 
  '@1..dockerComposeFileSpecs'
)
_docker_down() { eval "$_funcParams2"
  local dockerDir; dockerDir=$(dirname "${dockerComposeFileSpecs[ $(( ${#dockerComposeFileSpecs[@]} - 1 )) ]}")
  local -a OPT=()
  local dockerComposeFileSpec; for dockerComposeFileSpec in "${dockerComposeFileSpecs[@]}"; do
    if [[ $(dirname "$dockerComposeFileSpec") != "$dockerDir" ]]; then
      local srcFileSpec="$dockerComposeFileSpec"
      dockerComposeFileSpec="$dockerDir/$(basename "$dockerComposeFileSpec")"
    fi
    OPT+=( -f "$dockerComposeFileSpec" )
  done
  _dockerCompose "${OPT[@]}" down --remove-orphans
}

export _cmd_selfTestParams=(
  '--http='
  '--https='
  '--upstream='
  '--dockerImageName='
  'bwProjShortcut'
  'bwProjDir'
)
_cmd_selfTest() { eval "$_funcParams2"
  local returnCode=0
  while true; do
    _exec -v all "$bwProjShortcut" docker up -f "${OPT_http[@]}" "${OPT_https[@]}" "${OPT_upstream[@]}" || { returnCode=$?; break; }
    local tstFileSpec="/tmp/$bwProjShortcut.selfTest"
    if [[ -n $http ]]; then
      _exec -v err curl -s "http://localhost:${http}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
      _exec -v all diff "$tstFileSpec" "$bwProjDir/docker/nginx/whoami/index.html" || { returnCode=$?; break; }
      rm "$tstFileSpec"
    fi
    if [[ -n $https ]]; then
      _exec -v err curl -s "https://localhost:${https}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
      _exec -v all diff "$tstFileSpec" "$bwProjDir/docker/nginx/whoami/index.html" || { returnCode=$?; break; }
      rm "$tstFileSpec"
    fi
    local funcName="_${bwProjShortcut}_selfTestAddOn"
    if _funcExists "$funcName"; then
      "$funcName"  || { returnCode=$?; break; }    
    fi
    if [[ -n $dockerImageName ]]; then
      _exec -v all $bwProjShortcut docker shell || { returnCode=$?; break; }    
    fi
    if [[ -n $http || -n $https ]]; then
      echo "${_nl}${_ansiWarn}ВНИМАНИЕ! Чтобы выйти из docker-контейнера, выполните команду ${_ansiCmd}exit 0${_ansiReset}"
      _exec -v all $bwProjShortcut docker shell "dev-${bwProjShortcut}-nginx" || { returnCode=$?; break; }    
    fi
    break
  done
  if [[ $returnCode -eq 0 ]]; then
    _ok "self-test complete"
  else
    _err "self-test failed"
  fi
  return $returnCode
}

export _docker_upParams=(
  '--force-recreate/f'
  '--http='
  '--https='
  '--upstream='
  '--dockerImageName='
  '--dockerImageIdFileSpec='
  '--dockerContainerName='
  '--bwProjShortcut='
  '--bwProjName='
  '@--restart'
  '--noCheck'
  '@1..dockerComposeFileSpecs'
)
_docker_up() { eval "$_funcParams2"

  if [[ -n $https && -n $http && $https -eq $http ]]; then
    errOrigin=1 _throw "ожидает, что значение ${_ansiOutline}http ${_ansiPrimaryLiteral}$http${_ansiErr} будет отличаться от значения ${_ansiOutline}https"
    return $?
  fi

  local dockerDir; dockerDir=$(dirname "${dockerComposeFileSpecs[ $(( ${#dockerComposeFileSpecs[@]} - 1 )) ]}")

  bw_install --silentIfAlreadyInstalled docker || return $?

  if [[ -n $dockerImageName && -z $noCheck ]]; then
    if [[ -z $(_docker image ls "$dockerImageName:latest" -q) ]]; then
      _docker -v all image pull "$dockerImageName:latest" || return $?
    fi
    local tstImageIdFileSpec="/tmp/$bwProjShortcut.image.id"
    _docker inspect --format "{{json .Id}}" "$dockerImageName:latest" > "$tstImageIdFileSpec"
    if [[ ! -f "$dockerImageIdFileSpec" ]] || ! _silent cmp "$tstImageIdFileSpec" "$dockerImageIdFileSpec"; then
      rm "$tstImageIdFileSpec"
      "${bwProjShortcut}_docker_build" || return $?
    else
      rm "$tstImageIdFileSpec"
    fi
  fi

  local dockerComposeEnvFileSpec="$dockerDir/.env"
  {
    echo "# file generated by $bwProjShortcut_docker_up"
    local varName; for varName in \
      _bwNginxConfDir \
      _bwSslFileSpecPrefix \
      _bwDir \
      _bwDevDockerBwDir \
      _bwFileSpec \
      _bwDevDockerBwFileSpec \
      _bwDevDockerProjDir \
    ; do
      echo "$varName=${!varName}"
    done
    echo "_${bwProjShortcut}DockerContainerName=$dockerContainerName"
    [[ -z $http ]] || echo "_${bwProjShortcut}DockerHttp=$http"
    [[ -z $https ]] || echo "_${bwProjShortcut}DockerHttps=$https"
    [[ -z $upstream ]] || echo "_${bwProjShortcut}DockerUpstream=$upstream"
    [[ -z $dockerImageName ]] || echo "_${bwProjShortcut}DockerImageName=$dockerImageName"
  } > "$dockerComposeEnvFileSpec"

  if [[ -n $dockerImageName ]]; then
    local promptHolder="_${bwProjShortcut}Prompt"
    if [[ -z ${!promptHolder} ]]; then
      local prompt; _preparePrompt --user "$bwProjShortcut" --userAnsi PrimaryLiteral --gitDirty -
      eval export "$promptHolder"'="$prompt"'
    fi
    local dockerContainerEnvFileSpec="$dockerDir/main.env"
    { 
      echo "# file generated by '$bwProjShortcut'_docker_up"
      echo "_isBwDevelopInherited=$_isBwDevelop"
      echo "BW_SELF_UPDATE_SOURCE=$BW_SELF_UPDATE_SOURCE"
      echo "_bwProjName=$bwProjName"
      echo "_bwProjShortcut=$bwProjShortcut"
      echo "_hostUser=$(whoami)"
      [[ -z $http ]] || echo "_dockerHttp=$http"
      [[ -z $https ]] || echo "_dockerHttps=$https"
      [[ -z $upstream ]] || echo "_dockerUpstream=$upstream"
      echo "_prompt=${!promptHolder}"

      local addonFuncName="_${bwProjShortcut}_docker_upAddon"
      if _funcExists "$addonFuncName"; then 
        $addonFuncName
      fi
    } > "$dockerContainerEnvFileSpec"
  fi

  if [[ -n $http || -n $https ]]; then
    local needProcessNginxConfTemplate=
    if [[ "${#restart[@]}" -eq 0 ]]; then
      needProcessNginxConfTemplate=true
    else
      local serviceNameToRestart; for serviceNameToRestart in "${restart[@]}"; do
        if [[ $serviceNameToRestart =~ nginx ]]; then
          needProcessNginxConfTemplate=true
          break
        fi
      done
    fi
    local nginxDir="$dockerDir/nginx"
    if [[ -n $needProcessNginxConfTemplate ]]; then
      local separatorLine='# !!! you SHOULD put all custom rules above this line; all lines below will be autocleaned'
      if [[ -f "$nginxDir/.gitignore" ]]; then
        awk "!a;/$separatorLine/{a=1}" "$nginxDir/.gitignore" > "$nginxDir/.gitignore.new"
      elif [[ -f "$nginxDir/.gitignore.new" ]]; then
        echo "$separatorLine" > "$nginxDir/.gitignore.new"
      fi
      local templateExt='.template'
      local -a templateFileSpecs; mapfile -t templateFileSpecs < <(find "$nginxDir" -name "*$templateExt" 2>/dev/null)
      local templateFileSpec; for templateFileSpec in "${templateFileSpecs[@]}"; do
        fileSpec=${templateFileSpec:0:$(( ${#templateFileSpec} - ${#templateExt} ))}
        local relativeFileSpec="${fileSpec:$(( ${#nginxDir} + 1 ))}"
        if \
          http="$http" \
          https="$https" \
          upstream="$upstream" \
          projName="$bwProjName" \
          projShortcut="$bwProjShortcut" \
          _mkFileFromTemplate "$fileSpec"
        then
          echo "$relativeFileSpec" >> "$nginxDir/.gitignore.new"
        fi
      done
      if [[ -f "$nginxDir/.gitignore.new" ]]; then
        mv "$nginxDir/.gitignore.new" "$nginxDir/.gitignore"
      fi
    fi
  fi
  [[ -z $http ]] || eval export "_${bwProjShortcut}DockerHttp"'="$http"'
  [[ -z $https ]] || eval export "_${bwProjShortcut}DockerHttps"'="$https"'
  [[ -z $upstream ]] || eval export "_${bwProjShortcut}DockerUpstream"'="$upstream"'

  local -a OPT=()
  local dockerComposeFileSpec; for dockerComposeFileSpec in "${dockerComposeFileSpecs[@]}"; do
    if [[ $(dirname "$dockerComposeFileSpec") != "$dockerDir" ]]; then
      local srcFileSpec="$dockerComposeFileSpec"
      dockerComposeFileSpec="$dockerDir/$(basename "$dockerComposeFileSpec")"
      cp "$srcFileSpec" "$dockerComposeFileSpec"
    fi
    OPT+=( -f "$dockerComposeFileSpec" )
  done

  if [[ "${#restart[@]}" -gt 0 ]]; then
    _dockerCompose -v all "${OPT[@]}" stop "${restart[@]}"
  fi

  local stderrFileSpec="/tmp/docker-compose.stderr"
  {
    [[ -z $forceRecreate ]] || OPT_forceRecreate=( --force-recreate )
    _dockerCompose -v all "${OPT[@]}" up -d "${OPT_forceRecreate[@]}" --remove-orphans
  } 2> >(tee "$stderrFileSpec"); local returnCode=$?

  if [[ $returnCode -ne 0 ]] && grep -P -o '(?<=:)\d+(?= failed: port is already allocated)' "$stderrFileSpec" >/dev/null; then
    _err "port is already allocated"
  fi
  return $returnCode
}

export _docker_shellParams=( 
  'containerName' 
  'mainContainerName' 
  'shell'
)
_docker_shell() { eval "$_funcParams2"
  if [[ $containerName == "$mainContainerName" && $shell == /bin/bash ]]; then
    _docker exec -it "$containerName" "$shell" --init-file "$_bwDevDockerEntryPointFileSpec"
  else
    _docker exec -it "$containerName" "$shell"
  fi
}

export _dockerComposeContainerNamesParams=( 'dockerComposeFileSpec' )
_dockerComposeContainerNames() { eval "$_funcParams2"
  local -a result=()
  local nameHolder; for nameHolder in $(perl -ne 'print if s/^\s*container_name:\s*(.+?)\s*/$1/' "$dockerComposeFileSpec"); do
    eval result+=\( "$nameHolder" \)
  done
  echo "${result[@]}"
}

export _dockerComposeServiceNamesParams=( 'dockerComposeFileSpec' )
_dockerComposeServiceNames() { eval "$_funcParams2"
  perl -ne 'print if s/^\s\s(\w[\w\d]*):\s*$/$1\n/' "$dockerComposeFileSpec"
}

# =============================================================================

export bw_projectInfoShortcuts=( 'pi' )
export bw_projectInfoParamsOpt=( --canBeMixedOptionsAndArgs )
export bw_projectInfoParams=()
bw_projectInfoParams() {
  varName="${FUNCNAME[0]}" codeHolder=_codeToUseCache eval "$_evalCode"
  local -a bwProjShortcuts=()
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    if [[ ! $bwProjShortcut =~ $_isValidVarNameRegExp ]]; then
      _throw "bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle${_ansiErr} $_mustBeValidVarName"; return $?
    elif _hasItem "$bwProjShortcut" "${bwProjShortcuts[@]}"; then
      _throw "Duplicate bwProjShortcut ${_ansiPrimaryLiteral}$bwProjTitle"; return $?
    else
      bwProjShortcuts+=( "$bwProjShortcut" )
    fi
  done
  local bwProjShortcutsAsString="${bwProjShortcuts[*]}"
  bw_projectInfoParams=(
    '--all/a'
    '--def/d'
    "bwProjShortcut:?:( $bwProjShortcutsAsString )"
  )
  _saveToCache "${FUNCNAME[0]}"
}
export _codeToPrepareDescriptionsOf_bw_projectInfo='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_projectInfo_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'
export bw_projectInfo_bwProjShortcut_name='Сокращенное-имя-проекта'
export bw_projectInfo_bwProjShortcut_description='
  Если ${_ansiOutline}$bw_projectInfo_bwProjShortcut_name${_ansiReset} не задано, то выводит информация обо всех ${_ansiOutline}обнаруженных${_ansiReset} проектах
  С опцией ${_ansiCmd}--all${_ansiReset} -- обо всех проектах
'
export bw_projectInfo_all_description='
  Вывести информацию обо всех проектах
  Без опции ${_ansiCmd}--all${_ansiReset} -- обо всех ${_ansiOutline}обнаруженных${_ansiReset} проектах
'
export bw_projectInfo_def_description='
  Вывести ${_ansiOutline}определение${_ansiReset} проекта/проектов
'
export bw_projectInfo_description='Выводит информацию о проекте/проектах'
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
      echo "Всего $(_getPluralWord ${#found[@]} определен определено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}${found[*]}${_ansiReset}"
    else
      if [[ ${#found[@]} -gt 0 ]]; then
        _ok "Всего $(_getPluralWord ${#found[@]} обнаружен обнаружено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}${found[*]}"
      else
        _warn "Не обнаружено ни одного проекта"
      fi
    fi
  fi
}

_prepareDuplicatePorts() {
  local -a varNames; mapfile -t varNames < <(compgen -v __usedPort)
  unset "${varNames[@]}"

  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    if [[ -n $bwProjDefaultHttp ]]; then
      eval local "__usedPort$bwProjDefaultHttp"
      eval "__usedPort$bwProjDefaultHttp"'=$(( __usedPort$bwProjDefaultHttp + 1 ))'
    fi
    if [[ -n $bwProjDefaultHttps ]]; then
      eval local "__usedPort$bwProjDefaultHttps"
      eval "__usedPort$bwProjDefaultHttps"'=$(( __usedPort$bwProjDefaultHttps + 1 ))'
    fi
  done
  duplicatePorts=()
  local prefix=__usedPort
  local varName; for varName in $(compgen -v $prefix); do
    if [[ $varName -gt 1 ]]; then
      duplicatePorts+=( "${varName:${#prefix}}" )
    fi
  done
}

_bwProjectInfoHelper() {
  if [[ -n $def ]]; then
    echo "$bwProjShortcut:"
    echo "  name: $bwProjName"
    echo "  gitOrigin: ${_ansiUrl}$(_quotedArgs "$bwProjGitOrigin")${_ansiReset}"

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
    local matchRegxp='^\s*\.\s+\"?(.+?)\/bin\/'"$bwProjShortcut"'\.bash\"?\s*$'
    if grep -E "$matchRegxp" "$_profileFileSpec" >/dev/null 2>&1; then
      local alreadyProjDir; alreadyProjDir=$(perl -ne "print \$1 if /$matchRegxp/" "$_profileFileSpec" | tail -n 1)
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

export bw_installParams=(
  "${_verbosityParams[@]}"
  "--force/f"
  "--silentIfAlreadyInstalled"
)
export bw_installParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
export bw_install_cmd_name=Имя-приложения
export bw_install_description='устанавливает приложения'
export bw_install_force_description='устанавливает приложения, даже если оно уже установлено'
bw_install() { eval "$_funcParams2"
}

export _codeToInstallApp='
  local checkFuncName="_${FUNCNAME[0]}Check"
  if ! _funcExists $checkFuncName; then
    _throw "ожидает, что функция ${_ansiCmd}$checkFuncName${_ansiReset} будет определена"
    return $?
  fi
  if [[ -n $force ]] || ! $checkFuncName; then
    echo "${_ansiWarn}Установка ${_ansiCmd}$name${_ansiReset} . . ."
    showResult=_showInstallResult codeHolder=_codeToInstallOrRunApp eval "$_evalCode"
  elif [[ -z $silentIfAlreadyInstalled ]]; then
    _ok "${_ansiCmd}$name${_ansiOK} уже ${alreadyInstalled:-установлен}"
  fi
'
export _codeToRunApp='
  showResult=_showRunResult codeHolder=_codeToInstallOrRunApp eval "$_evalCode"
'
export _codeToInstallOrRunApp='
  codeHolder=_codeToInitSubOPT eval "$_evalCode"
  local returnCode=0
  _osSpecific || return $?
  $showResult
  return $returnCode
'
_osSpecific() {
  local funcName=${FUNCNAME[1]}
  local osSpecificFuncName=
  if [[ $OSTYPE =~ ^darwin ]]; then
    osSpecificFuncName="_${funcName}Darwin"
  elif [[ $OSTYPE =~ ^linux ]]; then
    osSpecificFuncName="_${funcName}Linux"
  fi
  if [[ -n $osSpecificFuncName ]] && _funcExists "$osSpecificFuncName"; then
    $osSpecificFuncName
  else
    osSpecificFuncName="_${funcName}"
    if [[ -n $osSpecificFuncName ]] && _funcExists "$osSpecificFuncName"; then
      $osSpecificFuncName
    else
      _throw "Неподдерживамая ОС ${_ansiPrimaryLiteral}$OSTYPE"
      return 1
    fi
  fi
}

# =============================================================================

export bw_install_brewParams=()
export bw_install_brew_description='Устанавливает Homebrew'
export bw_install_brewCondition='[[ $OSTYPE =~ ^darwin ]]'
bw_install_brew() { eval "$_funcParams2"
  name=brew codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_brewCheck() {
  _which brew
}
_bw_install_brewDarwin() {
  _exec "${sub_OPT[@]}" /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
}

# =============================================================================

export bw_install_gdateParams=()
export bw_install_gdate_description='Устанавливает gdate'
export bw_install_gdateCondition='[[ $OSTYPE =~ ^darwin ]]'
bw_install_gdate() { eval "$_funcParams2"
  name=gdate codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_gdateCheck() {
  _which gdate
}
_bw_install_gdateDarwin() {
  bw_install --silentIfAlreadyInstalled brew
  _exec "${sub_OPT[@]}" brew install coreutils
}

# =============================================================================

# Install Bash version 4 on MacOS X: https://gist.github.com/samnang/1759336

# =============================================================================

export bw_install_gitParams=()
export bw_install_git_description='Устанавливает git'
bw_install_git() { eval "$_funcParams2"
  name=git codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_gitCheck() {
  _which git
}
_bw_install_gitDarwin() {
  while true; do
    bw_install brew --silentIfAlreadyInstalled || { returnCode=$?; break; }
    _exec brew install "${OPT_force[@]}" git || returnCode=$?
    break
  done
}
_bw_install_gitLinux() {
  _exec "${sub_OPT[@]}" --sudo apt install -y --force-yes git || returnCode=$?
}

# =============================================================================

export bw_install_githubKeygenParams=()
export bw_install_githubKeyGen_description='Устанавливает github-keygen (${_ansiUrl}https://github.com/dolmen/github-keygen${_ansiReset})'
bw_install_githubKeygen() { eval "$_funcParams2"
  [[ ! $OSTYPE =~ ^linux  ]] || bw_install --silentIfAlreadyInstalled xclip || return $?
  name=github-keygen codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_githubKeygenCheck() {
  [[ -d "$_bwDir/github-keygen" ]]
}
_githubKeygenPastedMarkFileSpec='/tmp/github-keygen-pasted'
_githubKeygenFixFileSpec='/tmp/github-keygen-fix.bash'
_bw_install_githubKeygen() {
   while true; do
    _exec git clone https://github.com/dolmen/github-keygen.git "$_bwDir/github-keygen" || { returnCode=$?; break; }
    # patch github-keygen
    perl -i.bak -pe '
      $_="$_        `touch \"'"$_githubKeygenPastedMarkFileSpec"'\"`; # patched by bw.bash\n" if /^\s*close \$clip;\s*$/;
      $_="        open(my \$fh, \">\", \"'"$_githubKeygenFixFileSpec"'\"); printf \$fh \"chmod u-x,og-wx %s\", SSH_CONFIG_FILE; close \$fh; # patched by bw.bash\n$_" if /^\s*die\s+sprintf\("%s:\s+bad\s+file\s+permissions/;
    ' "$_bwDir/github-keygen/github-keygen" || { returnCode=$?; break; }
    break
  done 
}

# =============================================================================

export bw_install_xclipParams=()
# bw_install_xclipCondition='[[ $OSTYPE =~ ^darwin ]]' # вынужден был закоменнтировать, потому что bw при сборке на macOS не включает xclip в список доступных вариантов bw_install, и потом на Linux xclip недоступен
export bw_install_xclip_description='Устанавливает xclip (нужен для работы ${_ansiCmd}bw github-keygen${_ansiReset})'
bw_install_xclip() { eval "$_funcParams2"
  name=xclip codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_xclipCheck() {
  _which xclip
}
_bw_install_xclipLinux() {
  while true; do
    _exec --sudo apt-get update || { returnCode=$?; break; }
    _exec --sudo apt-get install -y --force-yes xclip || { returnCode=$?; break; }
    break
  done
}

# =============================================================================

export bw_githubKeygenParams=( 'username' )
bw_githubKeygen() { eval "$_funcParams2"
  bw_install github-keygen --silentIfAlreadyInstalled || return $?
  _rm "$_githubKeygenPastedMarkFileSpec"
  _rm "$_githubKeygenFixFileSpec"
  "$_bwDir/github-keygen/github-keygen" "$username" 
  if [[ -f $_githubKeygenFixFileSpec ]]; then
    . "$_githubKeygenFixFileSpec" || return $?
    _rm "$_githubKeygenFixFileSpec"
    "$_bwDir/github-keygen/github-keygen" "$username" 
  fi
  if [[ -f $_githubKeygenPastedMarkFileSpec ]]; then
    _rm "$_githubKeygenPastedMarkFileSpec"
    read -r -p "${_ansiWarn}Press ${_ansiPrimaryLiteral}Enter${_ansiWarn} to open browser${_ansiReset}" # https://unix.stackexchange.com/questions/293940/bash-how-can-i-make-press-any-key-to-continue
    _osSpecific || return $?
  fi
}
# _githubKeysUrl='https://github.com/settings/keys'
_githubKeysUrl='https://github.com/baza-winner/bw/wiki/%D0%A3%D1%81%D1%82%D0%B0%D0%BD%D0%BE%D0%B2%D0%BA%D0%B0-ssh-%D0%BA%D0%BB%D1%8E%D1%87%D0%B0-%D0%B4%D0%BB%D1%8F-%D0%B4%D0%BE%D1%81%D1%82%D1%83%D0%BF%D0%B0-%D0%BA-github@...'
_bw_githubKeygenDarwin() {
  local appGoogleChrome="/Applications/Google Chrome.app"
  if [[ -d "$appGoogleChrome" ]]; then
    /usr/bin/open -a "$appGoogleChrome" "$_githubKeysUrl"
  else
    /usr/bin/open "$_githubKeysUrl"
  fi
}
_bw_githubKeygenLinux() {
  if _which google-chrome; then
    google-chrome "$_githubKeysUrl"
  else
    xdg-open "$_githubKeysUrl"
  fi
}

# =============================================================================

export bw_install_dockerParams=()
export bw_install_docker_description='Устанавливает DockerCE ${_ansiUrl}https://www.docker.com/community-edition${_ansiReset}'
bw_install_docker() { eval "$_funcParams2"
  name=DockerCE codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_dockerCheck() {
  _which docker && [[ $(docker --version | perl -e '$\=undef; $_=<STDIN>; printf("%d%04d", $1, $2) if m/(\d+).(\d+)/') -ge 180003 ]]
}
# shellcheck disable=SC2034
_bw_install_dockerDarwin() {
  while true; do
    local -r appName=Docker
    local -r applicationsPath='/Applications'
    local -r appDir="${applicationsPath}/${appName}.app"
    local -r dmgFileSpec=~/Downloads/Docker.dmg
    if [[ -n $force || ! -d $appDir ]]; then
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
    elif [[ -z $silentIfAlreadyInstalled ]]; then
      _ok "${_ansiCmd}$appName${_ansiOK} уже установлен"
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
    _bw_install_dockerComposeLinux || { returnCode=$?; break; }
    break
  done
}

# =============================================================================

export bw_install_dockerComposeParams=()
export bw_install_dockerCompose_description="Устанавливает docker-compose"
bw_install_dockerCompose() { eval "$_funcParams2"
  bw_install docker --silentIfAlreadyInstalled || return $?
  name=docker-compose codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_dockerComposeCheck() {
  _which docker-compose
}
_bw_install_dockerComposeDarwin() {
  true
}
_bw_install_dockerComposeLinux() {
  while true; do
    # https://docs.docker.com/compose/install/#install-compose
    _exec "${sub_OPT[@]}" --sudo curl -L "https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo chmod +x /usr/local/bin/docker-compose || { returnCode=$?; break; }
    break
  done
}

# =============================================================================

export bw_runParams=(
  "${_verbosityParams[@]}"
)
export bw_runParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
export bw_run_cmd_name=Имя-приложения
export bw_run_description='запускает приложение'
bw_run() { eval "$_funcParams2"
}

# =============================================================================

export bw_run_dockerParams=()
export bw_run_docker_description="Запускает Docker"
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

export _redDir="$_bwDir/red"
export _redForMacUrl="https://static.red-lang.org/dl/mac/red-063"
export _redForMacFileSpec="$_redDir/red_063_for_mac"
export _redForLinuxUrl="https://static.red-lang.org/dl/linux/red-063"
export _redForLinuxFileSpec="$_redDir/red_063_for_linux"
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
    _download -s "$_redForMacUrl" "$_redForMacFileSpec" || { returnCode=$?; break; }
    chmod u+x "$_redForMacFileSpec" || { returnCode=99; break; }
    "$_redForMacFileSpec" "$@" || { returnCode=$?; break; }
    break
  done
}
_redLinux() {
  while true; do
    _download -s "$_redForLinuxUrl" "$_redForLinuxFileSpec" || { returnCode=$?; break; }
    chmod u+x "$_redForLinuxFileSpec" || { returnCode=99; break; }
    "$_redForLinuxFileSpec" "$@" || { returnCode=$?; break; }
    break
  done
}

# =============================================================================
