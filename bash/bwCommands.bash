#!/bin/bash

# shellcheck disable=SC2154,SC1090,SC2016
true

# =============================================================================

_resetBash

# =============================================================================

# shellcheck disable=SC2016,SC2034
_bwProjDefs=()

# =============================================================================

_bwProjDefs+=(
  'dubu' '
    --gitOrigin github.com:baza-winner/dev-ubuntu.git
    --branch master
    --docker-image-name bazawinner/dev-ubuntu
    --ssh 2201
    --docker-compose "docker-compose.main.yml"
  '
)

# shellcheck disable=SC2154
_bw_project_dubu() {
  _exec "${sub_OPT[@]}" git submodule update --init --remote --recursive # https://stackoverflow.com/questions/1777854/git-submodules-specify-a-branch-tag/15782629#15782629
}

# =============================================================================

_bwProjDefs+=(
  'dubu18' '
    --gitOrigin github.com:baza-winner/dev-ubuntu18.git
    --branch master
    --docker-image-name bazawinner/dev-ubuntu18
    --ssh 2206
    --docker-compose "docker-compose.main.yml"
  '
)

# shellcheck disable=SC2154
_bw_project_dubu18() {
  _exec "${sub_OPT[@]}" git submodule update --init --remote --recursive # https://stackoverflow.com/questions/1777854/git-submodules-specify-a-branch-tag/15782629#15782629
}

# =============================================================================

_bwProjDefs+=(
  'bwdev' '
    --gitOrigin github.com:baza-winner/bw.git
    --branch develop
    --http 8002
    --https 4402
    --no-docker-build
    --docker-compose "docker-compose.nginx.yml"
  '
)

# =============================================================================

_bwProjDefs+=(
  'bgate' '
    --gitOrigin github.com:baza-winner/billing-gate.git
    --branch develop
    --ssh 2203
    --http 8003
    --https 4403
    --upstream 6783
    --mysql 3303
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
  '
)

# shellcheck disable=SC2154
_bw_project_bgate() {
  _exec "${sub_OPT[@]}" git submodule update --init --recursive
}

# =============================================================================

_bwProjDefs+=(
  'crm' '
    --gitOrigin github.com:baza-winner/crm.git
    --branch develop
    --ssh 2204
    --http 8004
    --https 4404
    --mysql 3304
    --upstream 3000
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
  '
)

# =============================================================================

_bwProjDefs+=(
  'agate' '
    --gitOrigin github.com:baza-winner/agate.git
    --branch develop
    --ssh 2205
    --http 8005
    --https 4405
    --mysql 3305
    --upstream 3005
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
  '
)

# =============================================================================

_bwProjDefs+=(
  'mls' '
    --gitOrigin github.com:baza-winner/mls-pm.git
    --branch develop
    --ssh 2207
    --http 8007
    --https 4407
    --mysql 3307
    --postgresql 5407
    --elastic 9207
    --upstream 3000
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
  '
)

# =============================================================================

_bwProjDefs+=(
  'dip' '
    --gitOrigin github.com:baza-winner/dip2.git
    --branch develop
    --upstream 3000
    --ssh 2208
    --http 8008
    --https 4408
    --mysql 3308
    --redis 6308
    --webdis 7308
    --rabbitmq 5608
    --rabbitmqManagement 15608
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
    --docker-compose "docker-compose.upstream.yml"
    --docker-compose "docker-compose.mysql.yml"
  '
)

# shellcheck disable=SC2154
_bw_project_dip() {
  _exec "${sub_OPT[@]}" git submodule update --init --remote --recursive # https://stackoverflow.com/questions/1777854/git-submodules-specify-a-branch-tag/15782629#15782629
}

# =============================================================================

_bwProjDefs+=(
  'wcraw' '
    --gitOrigin github.com:baza-winner/wcrawler.git
    --ssh 2209
    --branch develop
    --docker-compose "docker-compose.main.yml"
  '
)

# =============================================================================

_bwProjDefs+=(
  'geo' '
    --gitOrigin github.com:baza-winner/geo2-api.git
    --branch feature/docker
    --ssh 2210
    --http 8010
    --https 4410
    --mysql 3310
    --elastic 9210
    --upstream 3000
    --docker-compose "docker-compose.nginx.yml"
    --docker-compose "docker-compose.main.yml"
  '
)

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

# shellcheck disable=SC2034
_codeToInitSubOPT='
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

# shellcheck disable=SC2034
{
bwParamsOpt=( '--canBeMixedOptionsAndArgs' '--isCommandWrapper' )
bwParams=()
bw_description='Базовая утилита проектов baza-winner'
}
bw() { eval "$_funcParams2"
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_versionParams=()
bw_version_description='Выводит номер версии bw.bash'
}
bw_version() { eval "$_funcParams2"
  local suffix="$BW_SELF_UPDATE_SOURCE"
  [[ -z $suffix ]] || suffix="@$suffix"
  echo "$_bwVersion$suffix"
}

# =============================================================================

# shellcheck disable=SC2034
{
_noPregen_params=( '!--noPregen/n' )
_noPregen_description="Исключить прегенерацию"
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_updateParams=(
  '--remove/r'
  "${_noPregen_params[@]}"
)
bw_update_remove_description="Удалить прегенеренные файлы перед обновлением"
bw_update_description="Обновляет bw"
}
bw_update() { eval "$_funcParams2"
  if [[ -n $remove ]]; then
    "$_bwFileSpec" rm -y
  fi
  local -a OPT=()
  if [[ -n $noPregen ]]; then
    OPT=( -p - )
  fi
  . "$_bwFileSpec" "${OPT[@]}"; local returnCode=$?
  echo "Current version: $(bw_version)"
  return $returnCode
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_removeParams=(
  '--yes/y'
  '--completely/c'
  "${_verbosityParams[@]}"
)
bw_remove_yes_description='Подтверждить удаление'
bw_remove_completely_description='Удалить не только все связанное с bw.bash, но и сам bw.bash'
bw_removeShortcuts=( 'rm' )
bw_remove_description="Удаляет bw.bash и все связанное с ним"
bw_removeCondition='[[ -z $_isBwDevelop ]] || ! _isInDocker'
}
bw_remove() { eval "$_funcParams2"
  if [[ -z $yes ]]; then
    _warn "Чтобы произвести удаление, запустите эту команду с опцией ${_ansiCmd}--yes${_ansiWarn} или ${_ansiCmd}-y"
    return 2
  else
    codeHolder=_codeToInitSubOPT eval "$_evalCode"

    if [[ -z $_isBwDevelop ]]; then
      local exactLine='. ~/bw.bash'
      if [[ $verbosity != dry ]]; then
        _setAtBashProfile -u "$exactLine" "$_bwMatchRegexp"
        # _exportVarAtBashProfile -u BW_SELF_UPDATE_SOURCE
        unset BW_SELF_UPDATE_SOURCE
      fi
    fi

    local totalUnsetFileSpec="/tmp/bw.remove.unset.bash"
    _rm "$totalUnsetFileSpec"
    find "$_bwDir" -name "*$_unsetFileExt" -exec sh -c 'cat $1 >> "$totalUnsetFileSpec"' _ {} \; # https://github.com/koalaman/shellcheck/wiki/SC2227
    local -a varNames; mapfile -t varNames < <(compgen -v | grep __upperCamelCaseToKebabCase_)
    echo "unset ${varNames[*]}" >> "$totalUnsetFileSpec"

    local msg=''
    if [[ $verbosity =~ ^(ok|all.*)$ ]]; then
      if [[ -n $_isBwDevelop ]]; then
        msg+="${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все прегенеренные вспомогательные файлы. Все основные ${_ansiFileSpec}*.bash${_ansiWarn}-файлы оставлены нетронутыми${_ansiReset}"
      else
        local suffix;
        local shortendBwFileSpec; shortendBwFileSpec="$(_shortenFileSpec "$_bwFileSpec")"
        if [[ -z $completely ]]; then
          suffix=", кроме ${_ansiFileSpec}$shortendBwFileSpec${_ansiWarn}. Для повторной установки выполните команду ${_ansiCmd}. $shortendBwFileSpec${_ansiWarn}"
        else
          suffix=", включая ${_ansiFileSpec}$shortendBwFileSpec${_ansiWarn}"
        fi
        local shortendBwDir; shortendBwDir="$(_shortenFileSpec "$_bwDir")"
        msg+="${_ansiWarn}Удалены команда ${_ansiCmd}bw${_ansiWarn} и все связанное с ней (содержимое директории ${_ansiFileSpec}$shortendBwDir${_ansiWarn})$suffix${_ansiReset}"
      fi
    fi

    if [[ -z $_isBwDevelop ]]; then
      [[ ! -d $_bwDir ]] || _rm "${sub_OPT[@]}" -d "$_bwDir"
      [[ -z $completely ]] || [[ ! -f $_bwFileSpec ]] || _rm "${sub_OPT[@]}" "$_bwFileSpec"
    else
      _rm "${sub_OPT[@]}" "$_bwDir/old.bw.bash"
      _rm "${sub_OPT[@]}" -d "$_bwDir/tgz"
      _rm "${sub_OPT[@]}" -pd "$_bwDir/tmp"
      local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
        if [[ $verbosity == dry ]]; then
          echo "${_ansiCmd}rm -rf \"$dirSpec\"${_ansiReset}"
        else
          rm -rf "$dirSpec"
        fi
      done
    fi

    if [[ $verbosity == dry ]]; then
      echo "${_ansiCmd}. \"$totalUnsetFileSpec\"${_ansiReset}"
      echo "${_ansiCmd}rm \"$totalUnsetFileSpec\"${_ansiReset}"
    else
      . "$totalUnsetFileSpec"
      rm "$totalUnsetFileSpec"
    fi

    echo "$msg"
  fi
}

# =============================================================================

# shellcheck disable=SC2034
{
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
bw_bashTestsParamsOpt=( '--canBeMixedOptionsAndArgs' )
bw_bashTestsParams=(
  '--noTiming'
  '--list'
  '@..args'
)
bw_bashTests_description='запускает тесты bash-функций'
}
bw_bashTests() {
  eval "$_funcParams2"
  local testsDirSpec="$_bwDir/tests"
  if [[ -z $_isBwDevelop ]]; then
    local testsSupportFileSpec="$testsDirSpec/testsSupport.bash"
    _getBwTar "$_bwFileSpec" tests | tar xf - -C "$_bwDir" \
      || { _err "Не удалось извлечь архив tests из ${_ansiFileSpec}$_bwFileSpec${_ansiErr} в ${_ansiFileSpec}$testsDirSpec"; return $?; }
    fileSpec="$testsSupportFileSpec" codeHolder=_codeSource eval "$_evalCode"
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
  for fileSpec in "$testsDirSpec/"*Tests.bash; do
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

# shellcheck disable=SC2034
{
bw_setParamsOpt=( '--canBeMixedOptionsAndArgs' '--isCommandWrapper' )
bw_setParams=(
  '!--uninstall/u'
)
bw_set_cmd_name=Имя-настройки
bw_set_description='включает/отключает настройку'
}
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

# shellcheck disable=SC2034
{
_enumSpace='(before after both none -)'
_enumOnOff='(yes no on off -)'
_ps_user_defaultAnsi='Green'
_ps_ptcd_defaultAnsi='White'
_ps_cd_defaultAnsi='White'

_preparePromptParams=()
}
_preparePromptParams() {
  varName=_preparePromptParams codeHolder=_codeToUseCache eval "$_evalCode"
  local -a _psItems; mapfile -t _psItems < <( compgen -v | sed -nEe 's/^_ps_([^_]+)$/\1/p' )
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
{
_getAnsiAsStringParamsOpt=( '--canBeMixedOptionsAndArgs' )
_getAnsiAsStringParams=(
  '--varName='
  "ansi:$_enumAnsi"
)
}
_getAnsiAsString() { eval "$_funcParams2"
  local resultHolder="_ansi${ansi}AsString"
  _getAnsiAsStringHelper "$resultHolder" || return $?
  if [[ -z $varName ]]; then
    echo "${!resultHolder}"
  else
    eval "$varName"'="${!resultHolder}"'
  fi
}

# shellcheck disable=SC2034
_getAnsiAsStringHelperParams=( 'varName' )
_getAnsiAsStringHelper() { eval "$_funcParams2"
  local awkFileSpec; _prepareAwkFileSpec || return $?
  _useCache --additionalDependencies "$awkFileSpec" "$varName"; local returnCode=$?
  [[ $returnCode -eq 2 ]] || return $returnCode
  local -a awk_OPT=(
    -f "$awkFileSpec"
    -v "searchFor=$varName"
    -v "funcName=${FUNCNAME[0]}"
  )
  local -a varNames
  local code; code="$(awk "${awk_OPT[@]}" "$_bwFileSpec")"
  codeHolder=code eval "$_evalCode"
  _saveToCache "${varNames[@]}"
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
        # _getAnsiAsString "$ansi" --varName $ansiAsStringHolder
        eval "$ansiAsStringHolder"'+="$(_getAnsiAsString "$ansi")"'
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

# shellcheck disable=SC2034
bw_set_promptParams=()
bw_set_promptParams() {
  _preparePromptParams
# shellcheck disable=SC2034
  bw_set_promptParams=( "${_preparePromptParams[@]}" )
}
# shellcheck disable=SC2034
bw_set_prompt_description='Настраивает prompt'
bw_set_prompt() { eval "$_funcParams2"
  local -a supportModules=( git ansi ps )
  local supportFileNameSuffix='Support.bash'
  # local newFileSpec="$_profileFileSpec.new"

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
    _exportVarAtBashProfile --uninstall -q OLD_PS1
    _exportVarAtBashProfile --uninstall -q PS1
  else
    local -a subOPT=(); local optVarName; for optVarName in "${__OPTVarNames[@]}"; do
      [[ $optVarName != OPT_uninstall && $optVarName != OPT_help ]] || continue
      dstVarName=OPT srcVarName=$optVarName eval "$_codeToInitLocalCopyOfArray"
      subOPT+=( "${OPT[@]}" )
    done
    local prompt; _preparePrompt "${subOPT[@]}" || return $?
    if _hasExportVarAtBashProfile -nq OLD_PS1; then
      export OLD_PS1="$PS1"
      _exportVarAtBashProfile -q OLD_PS1
    fi
    export PS1="$prompt"
    _exportVarAtBashProfile -q PS1
  fi
}

# shellcheck disable=SC2034
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

# shellcheck disable=SC2034
{
bw_set_viParams=()
bw_set_vi_description='Включает vi-режим для readline'
}
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

# shellcheck disable=SC2034
{
bw_set_horizontalScrollModeParams=()
bw_set_horizontalScrollMode_description='
  В${_ansiOutline}ы${_ansiReset}ключает horizontal-scroll-mode для readline
  Подробнее см. ${_ansiUrl}https://superuser.com/questions/848516/long-commands-typed-in-bash-overwrite-the-same-line/862341#862341${_ansiReset}
'
}
bw_set_horizontalScrollMode() { eval "$_funcParams2"
  local -a opts
  if [[ -n $uninstall ]]; then
    opts=(
      --horizontal-scroll-mode on
    )
  else
    opts=(
      --horizontal-scroll-mode off
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
  toFolder="в папку ${_ansiFileSpec}$(_shortenFileSpec "$folder")"
  fromFolder="из папки ${_ansiFileSpec}$(_shortenFileSpec "$folder")"
  titlePrefix=проект didInstall="установлен $toFolder" didUninstall="удален $fromFolder" toInstall="установить $toFolder" toUninstall="полностью удалить $fromFolder" _showResult
}

_showInstallResult() {
  titlePrefix='' didInstall="установлен" didUninstall="удален" toInstall="установить" toUninstall="полностью удалить" _showResult
}

_showRunResult() {
  titlePrefix='' didInstall="запущен" didUninstall="остановлен" toInstall="запустить" toUninstall="остановить" _showResult
}

# =============================================================================

_processProjDefs() {
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
  # local bwProjShortcutsAsString="${bwProjShortcuts[*]}"
  verbosityDefault=allBrief silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
  # shellcheck disable=SC2034
  bw_projectParams=(
    '!--uninstall/u'
    '!--force/f'
    '--branch='
    '--projDir/p='
    "${_verbosityParams[@]}"
    "bwProjShortcut:( ${bwProjShortcuts[*]} )"
  )
  # shellcheck disable=SC2034
  bw_projectInfoParams=(
    '--all/a'
    '--def/d'
    "bwProjShortcut:?:( ${bwProjShortcuts[*]} )"
  )
  # shellcheck disable=SC2034
  bw_prepareParams=(
    "bwProjShortcut:( ${bwProjShortcuts[*]} )"
    "projDir:?"
  )
}
_processProjDefs

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

# shellcheck disable=SC2034
{
_codeToCallPrepareBwProjVarsHelper='
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

_prepareBwProjVarsHelperParams() {
  _prepareBwProjVarsHelperParams=(
    '--gitOrigin='
    '--branch='
    '--docker-image-name='
    '--no-docker-build'
    '@--docker-compose'
  )
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    _prepareBwProjVarsHelperParams+=( "--$port:$_tcpPortDiap" )
  done
}

# shellcheck disable=SC2034
_prepareBwProjVarsHelper() { eval "$_funcParams2"
  codeToGetBwProjDef="
    bwProjGitOrigin=\"$gitOrigin\"
    bwProjDefaultBranch=\"$branch\"
    bwProjDockerImageName=\"$dockerImageName\"
    bwProjNoDockerBuild=\"$noDockerBuild\"
    bwProjDockerCompose=( $(_quotedArgs "${dockerCompose[@]}") )
  "
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    local Port=$(_upperFirst "$port")
    local bwProjDefaultPort="bwProjDefault$Port"
    codeToGetBwProjDef+="$bwProjDefaultPort=\"${!port}\";"
  done
}

_bwProjPorts=(
  ssh # 22
  http # 80
  https # 443
  mysql # 3306
  postgresql # 5432
  elastic # 9200
  redis # 6379
  webdis # 7379
  rabbitmq # 5672
  rabbitmqManagement # 15672
)
_tcpPortDiap='1024..65535'

_codeToDeclareLocalBwProjVars='
  local bwProjName=""
  local bwProjTitle=""
  local bwProjGitOrigin=""
  local bwProjDefaultBranch=""
  local bwProjDockerImageName=""
  local bwProjNoDockerBuild=""
  local -a bwProjDockerCompose=()
'
}
_updateCodeToDeclareLocalBwProjVars() {
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    local Port=$(_upperFirst "$port")
    local bwProjDefaultPort="bwProjDefault$Port"
    _codeToDeclareLocalBwProjVars+="local $bwProjDefaultPort=;"
  done
}
_updateCodeToDeclareLocalBwProjVars

# shellcheck disable=SC2034
_codeToPrepareDescriptionsOf_bw_project='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_project_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'

# shellcheck disable=SC2034
{
bw_projectParamsOpt=( '--canBeMixedOptionsAndArgs' )
bwProjGlobalDefaultBranch=develop
_uninstall_description='Включить режим обратного действия'
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
}
bw_project() { eval "$_funcParams2"
  codeHolder=_codeToInitSubOPT eval "$_evalCode"

  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  [[ -n $branch ]] || branch=${bwProjDefaultBranch:-$bwProjGlobalDefaultBranch}

  if [[ -n $uninstall ]]; then
    _prepareProjDir "$bwProjShortcut" || return $?
  else
    [[ -n $projDir ]] || projDir="$HOME/$bwProjName"
    if [[ $verbosity != dry && -z $force ]]; then
      local projDirs; _prepareProjDirs '"$bwProjShortcut"' || return $?
      local __projDir; for __projDir in "${projDirs[@]}"; do
        if [[ $(_shortenFileSpec "$__projDir") == $(_shortenFileSpec "$projDir") ]]; then
          if [[ $verbosity != none ]]; then
            projDir="$(_shortenFileSpec "$projDir")"
            local cmd="${FUNCNAME[0]//_/ } $bwProjShortcut"
            local msg=
            msg+="Проект ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} уже установлен в ${_ansiFileSpec}$(_shortenFileSpec "$projDir")${_ansiWarn}$_nl"
            msg+="Перед повторной установкой его необходимо удалить командой$_nl"
            msg+="  ${_ansiCmd}$cmd -u $(_quotedArgs "projDir")${_ansiWarn}$_nl"
            msg+="или установить с опцией ${_ansiCmd}--force${_ansiWarn}:$_nl"
            msg+="  ${_ansiCmd}$cmd -f $(_quotedArgs "projDir")${_ansiWarn}"
            _warn "$msg"
          fi
          return 4
        fi
      done
    fi
  fi

  local returnCode=0
  while true; do

    if [[ -n $alreadyProjDir ]]; then
      local cmdFileSpec="$alreadyProjDir/bin/$bwProjShortcut.bash"
      if [[ -f "$cmdFileSpec" ]]; then
        fileSpec="$cmdFileSpec" _unsetBash "${sub_OPT_verbosity[1]}"
      fi
    fi

    bw_install --silentIfAlreadyInstalled git || { returnCode=$?; break; }

    while true; do
      if [[ -z $uninstall ]]; then
        if [[ -d $projDir ]]; then
          local gitDirty
          if ! _inDir -v none "$projDir" _prepareGitDirty "$bwProjGitOrigin"; then
            if [[ -z $(ls -A "$projDir") ]]; then
              _rm "${sub_OPT[@]}" -d "$projDir" \
                || { returnCode=$?; break; }
            else
              _err "Папка ${_ansiCmd}$(_shortenFileSpec "$projDir")${_ansiErr} существует и непуста. Ее надо предварительно удалить вручную" \
                || { returnCode=$?; break; }
            fi
          else
            _warn "Папка ${_ansiCmd}$(_shortenFileSpec "$projDir")${_ansiWarn} уже содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle"
            _exec "${sub_OPT[@]}" cd "$projDir" || { returnCode=$?; break; }
            _hasItem "$gitDirty" '?' '*' '+' '^' || _exec "${sub_OPT[@]}" git pull
            break
          fi
        fi
        _mkDir "${sub_OPT[@]}" "$projDir" || { returnCode=$?; brancheak; }
        local cloneStderrFileSpec="/tmp/$bwProjShortcut.clone.stderr"
        if ! _exec "${OPT_verbosity[@]}" --stderr "$cloneStderrFileSpec" git ls-remote -t "git@$bwProjGitOrigin" no-ref; then
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
        # _exec -v all -s no --cmdAsIs "git clone git@$bwProjGitOrigin \"$projDir\"" || { returnCode=$?; break; }
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
              _warn "Папка ${_ansiCmd}$(_shortenFileSpec "$projDir")${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn}; оставлена для ручного удаления"
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
          _exec --sudo "${sub_OPT[@]}" rm -rf "$projDir" \
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
    local cmdFileSpec matchRegexp exactLine; _bw_projectHelper || return $?
    if [[ -n $uninstall ]]; then
      _setAtBashProfile -u "$exactLine" "$matchRegexp" || return $?
    else
      if [[ ! -f "$cmdFileSpec" ]]; then
        local msg=
        msg+="Не найден файл ${_ansiFileSpec}bin/$bwProjShortcut.bash${_ansiErr}$_nl"
        msg+="Не удалось инициализировать команду ${_ansiCmd}$bwProjShortcut"
        [[ $verbosity == none  ]] || _err "$msg"
        returnCode=1
      elif ! BW_NEED_REGEN=true _exec "${sub_OPT[@]}" . "$cmdFileSpec"; then
        local msg=
        msg+="Не удалось инициализировать команду ${_ansiCmd}$bwProjShortcut"
        [[ $verbosity == none  ]] || _err "$msg"
        returnCode=1
      else
        _setAtBashProfile "$exactLine" "$matchRegexp" || return $?
        rm -f "$_bwDir/generated/$bwProjShortcut"*
        "$bwProjShortcut" update -c
        [[ $verbosity == none  ]] || echo "${_ansiWarn}Теперь доступна команда ${_ansiCmd}$bwProjShortcut${_ansiReset}"
        _exec "${sub_OPT[@]}" "$bwProjShortcut" -? || return $?
      fi
    fi
  fi

  return $returnCode
}
_bw_projectHelper() {
  cmdFileSpec="$projDir/bin/${bwProjShortcut}.bash"
  matchRegexp='^[ \t]*\.[ \t]*'"$(_quotedArgs "$(_shortenFileSpec "$cmdFileSpec")")"
  exactLine="bw prepare $bwProjShortcut $(_quotedArgs "$(_shortenFileSpec "$projDir")")"
}

# shellcheck disable=SC2034
_prepareGitDirtyParams=( 'originSuffix' )
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

# shellcheck disable=SC2034
{
bw_prepare_description='Предназначено для использования в ~/$_profileFileName'
bw_prepare_bwProjShortcut_name='Сокращенное-имя-проекта'
bw_prepare_projDir_name='Папка-проекта'
}
bw_prepare() { eval "$_funcParams2"
  _prepareProjDir "$bwProjShortcut" || return $?
  fileSpec="$projDir/bin/$bwProjShortcut.bash" eval "$_codeSource"
  $bwProjShortcut update -c -p "$projDir"
  if [[ $(basename "${BASH_SOURCE[-1]}") != "$_profileFileName" ]]; then
    local cmdFileSpec matchRegexp exactLine; _bw_projectHelper || return $?
    _setAtBashProfile "$exactLine" "$matchRegexp"
  fi
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

# shellcheck disable=SC2034
{
_bwDevDockerEntryPointFileSpec="/home/dev/.bw/docker/helper/entrypoint.bash"
# _bwDevDockerBwDir="/home/dev/.bw"
_bwSslFileSpecPrefix="$_bwDir/ssl/server."

_projDir_description='Папка проекта'
_completionOnly_description='Обновить только completion-определения'
}

_initBwProjCmd() {
  local fileSpec; fileSpec=$(_getSelfFileSpec 2)
  local bwProjShortcut; bwProjShortcut=$(basename "$fileSpec" .bash)
  eval "$_codeToDeclareLocalBwProjVars" && _prepareBwProjVars || return $?
  eval export "_$bwProjShortcut"'DockerImageName="${bwProjDockerImageName:-bazawinner/dev-$bwProjShortcut}"'
  local -a funcNamesToRegen=()
  mapfile -t funcNamesToRegen < <( _getFuncNamesOfScriptToUnset "${BASH_SOURCE[1]}")

  eval "$bwProjShortcut"'_description='\''Базовая утилита проекта ${_ansiPrimaryLiteral}'"$bwProjName"' ${_ansiUrl}$bwProjGitOrigin${_ansiReset}'\'
  eval "$bwProjShortcut"'Params=()'
  eval "$bwProjShortcut"'ParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)'
  eval "$bwProjShortcut()"' { eval "$_funcParams2"
  }'

  funcNamesToRegen+=( "$bwProjShortcut"'_update' )
  eval "$bwProjShortcut"'_updateParams=(
    "${_noPregen_params[@]}"
    "!--completionOnly/c"
    !--projDir/p=
  )'
  eval "$bwProjShortcut"'_update_description='\''Обновляет команду ${_ansiCmd}'"$bwProjShortcut"'${_ansiReset}'\'
  eval "$bwProjShortcut"'_update() { eval "$_funcParams2"
    local sourceFileSpec
    if _isInDocker; then
      projDir="$HOME/proj"
    else
      _prepareProjDir '"$bwProjShortcut"' || return $?
    fi
    _cmd_update "${OPT_noPregen[@]}" "${OPT_completionOnly[@]}" '"$bwProjShortcut"' "$projDir/bin/'"$bwProjShortcut"'.bash"
  }'

  funcNamesToRegen+=( "$bwProjShortcut"'_docker' )
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
    dockerUpParams+=(
      '--no-check/n'
      '--portIncrement/i:0..=0'
    )
    local dockerImageTitle='${_ansiPrimaryLiteral}$_'"$bwProjShortcut"'DockerImageName${_ansiReset}'
    eval "$bwProjShortcut"'_docker_up_noCheck_description='\''Не проверять актуальность docker-образа '\''"$dockerImageTitle"'
    codeToPrepareOPTForDockerUp+='
      OPT+=(
        --dockerImageName "$_'"$bwProjShortcut"'DockerImageName"
      )
    '
    codeToPrepareOPTForSelfTest+='
      OPT+=(
        --dockerImageName "$_'"$bwProjShortcut"'DockerImageName"
      )
    '
    funcNamesToRegen+=( "$bwProjShortcut"'_docker_build' )
    eval "$bwProjShortcut"'_docker_buildParams=(
      --force
      !--projDir/p=
    )'
    eval "$bwProjShortcut"'_docker_build_force_description='\''Невзирая на возможное отсутствие изменений в docker/Dockerfile'\'
    eval "$bwProjShortcut"'_docker_build_description='\''Собирает docker-образ '\''"$dockerImageTitle"'
    eval "$bwProjShortcut"'_docker_build() { eval "$_funcParams2"
      _prepareProjDir '"$bwProjShortcut"' || return $?

      local -a params=(
        ${OPT_force[@]}
        '"$bwProjShortcut"'
        "$_'"$bwProjShortcut"'DockerImageName"
        "'"$dockerImageTitle"'"
      )
      _inDir "$projDir/docker" _docker_build "${params[@]}"
    }'

    funcNamesToRegen+=( "$bwProjShortcut"'_docker_push' )
    eval "$bwProjShortcut"'_docker_pushParams=()'
    eval "$bwProjShortcut"'_docker_push_description='\''Push-ит docker-образ '\''"$dockerImageTitle"'
    eval "$bwProjShortcut"'_docker_push() { eval "$_funcParams2"
      docker push "$_'"$bwProjShortcut"'DockerImageName"
    }'
  fi

  eval local -a 'additionalDependencies=(
    --additionalDependencies "$_'"$bwProjShortcut"'Dir/bin/'"$bwProjShortcut"'.bash"
  )'
  local needNoTestAccessMessage=''
  local port; for port in "${_bwProjPorts[@]:0:3}"; do
    _prepareVarsForDefaultPort $port && needNoTestAccessMessage=true
  done
  local port; for port in "${_bwProjPorts[@]:3}"; do
    _prepareVarsForDefaultPort $port
  done
  _prepareVarsForDefaultPort upstream 'для сервиса ${_ansiPrimaryLiteral}nginx${_ansiReset}'
  if [[ -n $needNoTestAccessMessage ]]; then
    dockerUpParams+=( '--noTestAccessMessage/m' )
    eval "${bwProjShortcut}_docker_up_noTestAccessMessage_description='Не выводить сообщение о проверке доступности docker-приложения'"
    codeToPrepareOPTForDockerUp+='
      OPT+=( ${OPT_noTestAccessMessage[@]} )
    '
  fi
  eval "$bwProjShortcut"'_selfTestParamsOpt=(
    --canBeMixedOptionsAndArgs
    "${additionalDependencies[@]}"
  )'
  funcNamesToRegen+=( "$bwProjShortcut"'_selfTest' )
  eval "$bwProjShortcut"'_selfTestParams=(
    "${selfTestParams[@]}"
    !--projDir/p=
    --auto/a
  )'
  eval "$bwProjShortcut"'_selfTestCondition='\''! _isInDocker'\'
  eval "$bwProjShortcut"'_selfTestShortcuts=( st )'
  eval "$bwProjShortcut"'_selfTest_auto_description='\''В автоматическом режиме (с использованием ${_ansiCmd}expect${_ansiReset})'\'
  eval "$bwProjShortcut"'_selfTest_description='\''Самопроверка'\'
  eval "$bwProjShortcut"'_selfTest() { eval "$_funcParams2"
    _prepareProjDir '"$bwProjShortcut"' || return $?

    local -a OPT=( ${OPT_auto[@]} )
    '"$codeToPrepareOPTForSelfTest"'
    _inDir -v none "$projDir/docker" _cmd_selfTest "${OPT[@]}" '"$bwProjShortcut"' "$projDir"
  }'

  eval "$bwProjShortcut"'_docker_upParamsOpt=(
    --canBeMixedOptionsAndArgs
    "${additionalDependencies[@]}"
  )'
  funcNamesToRegen+=( "$bwProjShortcut"'_docker_up' )
  local funcName="_${bwProjShortcut}_clean"
  local forceRecreateAndCleanMappedVolumesDef=
  local forceRecreateAndCleanMappedVolumesInclude=
  dstVarName=cleanDockerDirs srcVarName="_${bwProjShortcut}_clean_docker_dirs" codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
  if [[ -n $bwProjDefaultMysql ]]; then
    cleanDockerDirs=(
      "mysql/.var_lib"
      "mysql/log"
      "${cleanDockerDirs[@]}"
    )
  fi
  cleanDockerDirs=(
    ".nvm"
    "${cleanDockerDirs[@]}"
  )

  if _funcExists $funcName; then
    forceRecreateAndCleanMappedVolumesDef='--force-recreate-and-clean-mapped-volumes/F'
    forceRecreateAndCleanMappedVolumesInclude='"${OPT_forceRecreateAndCleanMappedVolumes[@]}"'
    eval "$bwProjShortcut"'_docker_up_forceRecreateAndCleanMappedVolumes_description='\''
      Поднимает ${_ansiPrimaryLiteral}$(_quotedArgs ${containerNames[@]})${_ansiReset}
      с опцией ${_ansiCmd}--force-recreate${_ansiReset}, передаваемой ${_ansiCmd}docker-compose up${_ansiReset},
      предварительно очищая комадой ${_ansiCmd}'_${bwProjShortcut}'_clean${_ansiReset} подключаемые к образу папки
      Что именно делаeт ${_ansiCmd}'_${bwProjShortcut}'_clean${_ansiReset} можно проверить, запустив ее с опцией ${_ansiCmd}-v dry${_ansiReset}:
        ${_ansiCmd}'_${bwProjShortcut}'_clean -v dry${_ansiReset}
      Примечание: текущей папкой на время выполнения ${_ansiCmd}'_${bwProjShortcut}'_clean${_ansiReset} становится ${_ansiCmd}$(_prepareProjDir '${bwProjShortcut}' && _shortenFileSpec $projDir)/docker${_ansiReset}
    '\'
  fi
  eval "$bwProjShortcut"'_docker_upParams=(
    "${dockerUpParams[@]}"
    !--projDir/p=
    --force-recreate/f
    '$forceRecreateAndCleanMappedVolumesDef'
    '\''@--restart/r:( $(_dockerComposeServiceNames '"$bwProjShortcut"' $projDir) )'\''
    "${'"$bwProjShortcut"'_docker_upParamsAddon[@]}"
  )'
  eval '_codeToPrepareVarsForDescriptionsOf_'"$bwProjShortcut"'_docker_up='\''
    local -a containerNames; _dockerComposeContainerNames -v containerNames '"$bwProjShortcut"' $projDir || return $?
  '\'
  eval "$bwProjShortcut"'_docker_up_portIncrement_description='\''
    Смещение относительно базовых значений для всех портов
    Полезно для старта второго экземпляра docker-приложения ${_ansiSecondaryLiteral}'${bwProjShortcut}'${_ansiReset}
    Примечание: второй экземпляр следует запускать из второй копии проекта, которую можно установить командой:
      ${_ansiCmd}bw p '${bwProjShortcut}' -p ${_ansiOutline}Абсолютный-путь-к-папке-второй-копии-проекта${_ansiReset}
  '\'
  eval "$bwProjShortcut"'_docker_up_restart_description='\''Останавливает и поднимает указанные сервисы'\'
  eval "$bwProjShortcut"'_docker_up_forceRecreate_description='\''
    Поднимает ${_ansiPrimaryLiteral}$(_quotedArgs ${containerNames[@]})${_ansiReset}
    с опцией ${_ansiCmd}--force-recreate${_ansiReset}, передаваемой ${_ansiCmd}docker-compose up${_ansiReset}
  '\'
  eval "$bwProjShortcut"'_docker_up_description='\''Поднимает (docker-compose up) следующие контейнеры: ${_ansiPrimaryLiteral}$(_quotedArgs ${containerNames[@]})${_ansiReset}'\'
  eval "$bwProjShortcut"'_docker_up() { eval "$_funcParams2"
    _prepareProjDir '"$bwProjShortcut"' || return $?

    local -a OPT=(
      "${OPT_noCheck[@]}"
      "${OPT_forceRecreate[@]}"
      '$forceRecreateAndCleanMappedVolumesInclude'
      "${OPT_restart[@]}"
      --bwProjShortcut '"$bwProjShortcut"'
      --bwProjName "'"$bwProjName"'"
    )

    '"$codeToPrepareOPTForDockerUp"'

    _inDir -v none "$projDir/docker" _docker_up "${OPT[@]}" '"${bwProjDockerCompose[*]}"'
  }'

  eval '_'"$bwProjShortcut"'_prepareDockerComposeYmlParams=(
    "${dockerUpParams[@]}"
    !--projDir/p=
    "${'"$bwProjShortcut"'_docker_upParamsAddon[@]}"
  )'
  eval '_'"$bwProjShortcut"'_prepareDockerComposeYml() {
    _prepareProjDir '"$bwProjShortcut"' || return $?
    eval "$_funcParams2"

    local -a OPT=(
      --bwProjShortcut '"$bwProjShortcut"'
    )
    '"$codeToPrepareOPTForDockerUp"'
    _inDir -v none "$projDir/docker" _prepareDockerComposeYml "${OPT[@]}" '"${bwProjDockerCompose[*]}"'
  }'

  eval '_codeToPrepareVarsForDescriptionsOf_'"$bwProjShortcut"'_docker_down='\''
    local -a containerNames; _dockerComposeContainerNames -v containerNames '"$bwProjShortcut"' $projDir || return $?
  '\'
  eval "$bwProjShortcut"'_docker_down_description='\''Останавливает (${_ansiCmd}docker-compose down${_ansiReset}) следующие контейнеры: ${_ansiPrimaryLiteral}${containerNames[*]}${_ansiReset}'\'
  funcNamesToRegen+=( "$bwProjShortcut"'_docker_down' )
  eval "$bwProjShortcut"'_docker_downParams=( !--projDir/p= )'
  eval "$bwProjShortcut"'_docker_down() { eval "$_funcParams2"
    _prepareProjDir '"$bwProjShortcut"' || return $?
    _inDir -v none "$projDir/docker" _docker_down "$projDir"
  }'

  eval "$bwProjShortcut"'_docker_shellParamsOpt=( --canBeMixedOptionsAndArgs )'
  funcNamesToRegen+=( "$bwProjShortcut"'_docker_shell' )
  eval "$bwProjShortcut"'_docker_shellParams=(
    !--projDir/p=
    '\''serviceName:( $(_dockerComposeServiceNames '"$bwProjShortcut"' $projDir) )=$(_defaultDockerComposeServiceName '"$bwProjShortcut"' $projDir)'\''
    '\''shell=$(_getDefaultShellOfDockerService "$serviceName")'\''
  )'
  eval "$bwProjShortcut"'_docker_shell_description='\''Запускает bash в docker-контейнере'\'
  eval "$bwProjShortcut"'_docker_shell_serviceName_name='\''Имя-сервиса'\'
  eval "$bwProjShortcut"'_docker_shell_shell_name='\''Командная-оболочка'\'
  eval "$bwProjShortcut"'_docker_shell() { eval "$_funcParams2"
    _prepareProjDir '"$bwProjShortcut"' || return $?
    _inDir -v none "$projDir/docker" _docker_shell "$serviceName" "$shell"
  }'

  if [[ -n $BW_NEED_REGEN ]]; then
    local -a fileNamesToDelete=()
    local funcName; for funcName in "${funcNamesToRegen[@]}"; do
      local codeType; for codeType in funcParams help; do
        fileNamesToDelete+=( "$funcName.$codeType.code.bash" )
      done
    done
    _inDir "$_bwDir/generated" _rm "${fileNamesToDelete[@]}"
  fi
}

# shellcheck disable=SC2034
_prepareProjDirParams=( 'bwProjShortcut' )
_prepareProjDir() { eval "$_funcParams2"
  if [[ -z $projDir ]]; then
    local projDirs; _prepareProjDirs "$bwProjShortcut" || return $?
    local count=${#projDirs[@]}
    if [[ $count -eq 0 ]]; then
      # shellcheck disable=SC2046
      return $(_throw "не обнаружил в ${_ansiFileSpec}$_profileFileSpec${_ansiErr} строку, содержащую ${_ansiPrimaryLiteral}bin/$bwProjShortcut.bash")
    elif [[ $count -eq 1 ]]; then
      projDir=${projDirs[0]}
    else
      local curDir; curDir="$(_shortenFileSpec "$(pwd)")/"
      local __projDir; for __projDir in "${projDirs[@]}"; do
        __projDir="$(_shortenFileSpec "$__projDir")/"
        if [[ "${curDir:0:${#__projDir}}" == "$__projDir" ]]; then
          projDir="$__projDir"
          break
        fi
      done
      if [[ -z $projDir ]]; then
        if [[ -z $____isCompletionMode ]]; then
          local msg="${_ansiWarn}Переместитесь в одну из папок проекта:${_ansiReset}"
          local __projDir; for __projDir in "${projDirs[@]}"; do
            msg+=$_nl"  ${_ansiCmd}cd $(_quotedArgs "$(_shortenFileSpec "$__projDir")")${_ansiReset}"
          done
          msg+=$_nl"${_ansiWarn}Или укажите ее явно в качестве опции: ${_ansiCmd}-p ${_ansiOutline}Папка${_ansiReset}"
          echo "$msg"
        fi
        return 5
      fi
    fi
  fi
  local tilda="~"
  projDir=${projDir/$tilda/$HOME}
  projDir=${projDir%/} # https://stackoverflow.com/questions/1848415/remove-slash-from-the-end-of-a-variable/1848456#1848456
  local dockerDirSpec="$projDir/docker"
  if [[ ! -d "$dockerDirSpec" ]]; then
    if [[ -n $____isCompletionMode ]]; then
      return 1
    else
      # shellcheck disable=SC2046
      return $(_throw "Папка ${_ansiFileSpec}$(_shortenFileSpec "$projDir")${_ansiErr} не содержит проект ${_ansiPrimaryLiteral}$bwProjShortcut${_ansiErr}: не обнаружена папка ${_ansiFileSpec}$dockerDirSpec")
    fi
  fi
}

# shellcheck disable=SC2034
_runInDockerContainer='
  if ! _isInDocker; then
    _runInDockerContainer_outOfDocker
    return $?
  elif [[ $- =~ i && -z $wrapped ]]; then
    _runInDockerContainer_nonWrapped
    return $?
  else
    _runInDockerContainer_main
  fi
  local projDir=/home/dev/proj
'
_runInDockerContainer_outOfDocker() {
  local bwProjShortcut="${FUNCNAME[1]%%_*}"
  _prepareProjDir "$bwProjShortcut" || return $?
  _inDir -v none "$projDir/docker" _runInDockerContainer_outOfDockerHelper ${FUNCNAME[1]//_/ } "${__params[@]}"
}
_runInDockerContainer_nonWrapped() {
  local bwProjShortcut="${FUNCNAME[1]%%_*}"
  local code='
    [[ $- =~ i ]] && . "$HOME/.bashrc"
    . "$HOME/bw.bash" -p -
    . "$HOME/proj/bin/'"$bwProjShortcut"'.bash"
    wrapped=true '"${FUNCNAME[1]//_/ } $(_quotedArgs "${__params[@]}")"
  /usr/local/bin/dumb-init /bin/bash -c "$code"
}
_runInDockerContainer_outOfDockerHelper() {
  local -a dockerCompose_OPT; _prepareDockerComposeOpt
  local -a OPT_T; [[ -n $it ]] || OPT_T=( -T )
  local -a cmd_base=( _dockerCompose "${dockerCompose_OPT[@]:2}" exec "${OPT_T[@]}" main )
  local -a cmd=( "${cmd_base[@]}" "$_bwDevDockerEntryPointFileSpec" "$$" "$@" )
  if [[ -n $it ]]; then
    eval "$(_quotedArgs "${cmd[@]}")"
  else
    (
      eval "$(_quotedArgs "${cmd[@]}")" &
      local subPid=$!
      local pidFileSpec="$$.pid"
      # shellcheck disable=SC2064
      trap "trap - SIGINT; _runInDockerCtrlC $subPid \"$pidFileSpec\" \"$*\" $(_quotedArgs "${cmd_base[@]}")" SIGINT
      wait $subPid; local returnCode=$?
      trap - SIGINT
      return $returnCode
    )
  fi
}
_runInDockerCtrlC() {
  local subPid=$1; shift
  local pidFileSpec=$1; shift
  local title=$1; shift
  local wasPidFileSpec

  if _exist -v err "$pidFileSpec"; then
    wasPidFileSpec=true
    local pid; read -r -d $'\x04' pid < "$pidFileSpec"
    _spinner \
      -t "Отправка сигнала SIGINT в docker-контейнер заняла" \
      "Отправка сигнала SIGINT в docker-контейнер" \
      "$@" \
        kill -SIGINT "$pid"
  fi
  wait "$subPid"
  if [[ -n $wasPidFileSpec ]]; then
    echo "Остановлено выполнение ${_ansiCmd}$title${_ansiReset}"
  fi
}
_runInDockerCtrlCHelper() {
  while [[ -f "$pidFileSpec" ]]; do
    sleep 1
  done
}
_runInDockerContainer_main() {
  local queueFileSpec="$HOME/proj/docker/${queue:-default}.queue"; shift
  if [[ -s $queueFileSpec ]]; then
    local pid; read -r -d $'\x04' pid < "$queueFileSpec"
    if [[ $pid -ne $PPID ]] && ps axo pid | awk '{print $1}' | grep -x $pid > /dev/null; then
      _silent kill -SIGINT $pid
      sleep 2
      _silent kill -SIGTERM $pid
    fi
  fi
  if [[ $PPID -gt 8 ]]; then # защита от самоуничтожения
    echo $PPID > "$queueFileSpec"
  fi
}

# shellcheck disable=SC2034
_doOnceInContainer='
  local returnCode=0
  while true; do
    [[ -n $doWhat ]] || { returnCode=$?; _throw "expects non empty ${_ansiOutline}doWhat${_ansiErr}"; break; }
    local -a __doWhat=( $doWhat )
    _funcExists "${__doWhat[0]}" || { returnCode=$?; _throw "expects function ${_ansiOutline}${__doWhat[0]}${_ansiErr} to be defined"; break; }
    _mkDir "$HOME/did" || { returnCode=$?; break; }
    local markerFileSpec="$HOME/did/$doWhat"
    if [[ ! -f "$markerFileSpec" ]]; then
      eval "$doWhat" || { returnCode=$?; break; }
      touch "$markerFileSpec" || { returnCode=$?; break; }
    fi
    break
  done
  [[ $returnCode -eq 0 ]]
'

# shellcheck disable=SC2034
_prepareProjDirsParams=(
  '--profileFileSpec=$_profileFileSpec'
  'bwProjShortcut'
)
_prepareProjDirs() { eval "$_funcParams2"
  local sedCode
  # sedCode+='s/'"$_sourceMatchRegexp"'bin\/'"$bwProjShortcut"'\.bash''.*$/\1/p;'
  sedCode+='s/^[ \t]*\.[ \t]+"?([ a-zA-Z0-9\/~_-]+)\/bin\/'"$bwProjShortcut"'\.bash''.*$/\1/p;'

  # sedCode+='s/^bw[ \t]*prepare[ \t]*'"$bwProjShortcut"'[ \t]*"?([ a-zA-Z0-9\/~_-]+)"?.*/\1/p;'
  sedCode+='s/^bw[ ]prepare[ ]'"$bwProjShortcut"'[ ]"?([ a-zA-Z0-9\/~_-]+)"?.*/\1/p;'
  mapfile -t projDirs < <(sed -nEe "$sedCode" "$profileFileSpec")
}

# shellcheck disable=SC2034
_prepareVarsForDefaultPortParams=(
  'portName'
  'descriptionSuffix:?'
)
_prepareVarsForDefaultPort() { eval "$_funcParams2"
  local defaultHolder; defaultHolder="bwProjDefault$(_upperFirst "$portName")"
  if [[ -z ${!defaultHolder} ]]; then
    unset "${bwProjShortcut}_docker_up_${portName}_description"
    unset "${bwProjShortcut}_selfTest_${portName}_description"
    return 2
  else
    local paramDef="--$portName:$_tcpPortDiap=\$(( ${!defaultHolder} + portIncrement ))"
    local description="$portName-порт по которому будет доступно docker-приложение"
    if [[ -n $descriptionSuffix ]]; then
      description+=" $descriptionSuffix"
    fi
    local code='
      OPT+=( ${OPT_'"$portName"'[@]} )
    '
    dockerUpParams+=( "$paramDef" )
    selfTestParams+=( "$paramDef" )
    eval "${bwProjShortcut}_docker_up_${portName}_description='$description'"
    eval "${bwProjShortcut}_selfTest_${portName}_description='$description'"
    codeToPrepareOPTForDockerUp+="$code"
    codeToPrepareOPTForSelfTest+="$code"
    return 0
  fi

}

_getDefaultShellOfDockerContainer() {
  if [[ $1 =~ nginx$ ]]; then
    echo '/bin/sh'
  else
    echo '/bin/bash'
  fi
}

_getDefaultShellOfDockerService() {
  if [[ $1 == nginx ]]; then
    echo '/bin/sh'
  else
    echo '/bin/bash'
  fi
}

# shellcheck disable=SC2034
_cmd_updateParams=(
  "${_noPregen_params[@]}"
  '!--completionOnly/c'
  'bwProjShortcut'
  'sourceFileSpec'
)
_cmd_update() { eval "$_funcParams2"
  if [[ -z $completionOnly ]]; then
    local -a OPT=()
    [[ -z $noPregen ]] || OPT=( -p - )
    . "$_bwFileSpec" "${OPT[@]}" || return $?
    . "$sourceFileSpec" || return $?
    rm -f "$_bwDir/generated/$bwProjShortcut"*
  fi
  local -a __completions=()
  local -a funcNamesToPregen; mapfile -t funcNamesToPregen < <(compgen -c "$bwProjShortcut")
  _pregen "${funcNamesToPregen[@]}" || return $?
  for fileSpec in "${__completions[@]}"; do
    . "$fileSpec"
  done
}

_noMysql_description='Без подготовки MySQL'
_noPostgresql_description='Без подготовки PostgreSQL'
_noElasticsearch_description='Без подготовки Elasticsearch'

_declare_prepareLocal_jsParams=()
_declare_prepareLocal_js() { eval "$_funcParams2"
  local bwProjShortcut=$(basename "${BASH_SOURCE[1]}" .bash)
  eval "${bwProjShortcut}_prepareLocal"'ParamsOpt=( --canBeMixedOptionsAndArgs )'
  eval "${bwProjShortcut}_prepareLocal"'Params=(
    !--projDir/p=
    "${'"$bwProjShortcut"'_prepareLocalParamsAddon[@]}"
  )'
  eval "${bwProjShortcut}_prepareLocal"'_description='\''Производит подготовку локальных инструментов (баз данных и т.п)'\'
  eval "${bwProjShortcut}_prepareLocal"'() { eval "$_funcParams2"
    eval "$_runInDockerContainer"
    local bwProjShortcut='"$bwProjShortcut"'
    local funcName="_'"$bwProjShortcut"'_prepareLocal"
    _funcExists "$funcName" || return $(_throw "функция ${_ansiCmd}$funcName${_ansiErr} не определена")
    _'"$bwProjShortcut"'_prepareLocal
  }'
}

_declare_jsParamsOpt=( --canBeMixedOptionsAndArgs )
_declare_jsParams=(
  '--pm2Dir='
  '--configFile=.json'
  'kind:(api worker)'
  'whereConfig'
  '@2..envsAndYaml'
)
_declare_js() { eval "$_funcParams2"
  local bwProjShortcut=$(basename "${BASH_SOURCE[1]}" .bash)
  local -a envs=( "${envsAndYaml[@]::${#envsAndYaml[@]}-1}" )
  local yaml=${envsAndYaml[@]: -1:1}
  eval "${bwProjShortcut}_${kind}"'ParamsOpt=( --canBeMixedOptionsAndArgs )'
  eval "${bwProjShortcut}_${kind}"'Params=(
    !--projDir/p=
    "${'"$bwProjShortcut"'_prepareLocalParamsAddon[@]}"
    '\''applicationEnv:('"${envs[*]}"')='"${envs[0]}"''\''
  )'
  eval "${bwProjShortcut}_${kind}"'_applicationEnv_name='\''Имя-среды'\'
  eval "${bwProjShortcut}_${kind}"'_applicationEnv_description='\''Совпадает с именем одного из ${_ansiFileSpec}'"$configFile"'${_ansiReset}-файлов из ${_ansiFileSpec}'"$whereConfig"'${_ansiReset}'\'
  eval "${bwProjShortcut}_${kind}"'_description='\''(Пере)Запускает '"$kind"''\'
  eval "${bwProjShortcut}_${kind}"'() { eval "$_funcParams2"
    eval "$_runInDockerContainer"
    local bwProjShortcut='"$bwProjShortcut"'
    local kind="'"$kind"'"
    local yaml
    read -d "" yaml <<EOF'"$yaml"'EOF
    _cmd_js "$kind" "$bwProjShortcut" "'"$pm2Dir"'" "$yaml"
  }'
}

_cmd_jsParams=( 'kind' 'bwProjShortcut' 'pm2Dir:?' 'yaml' )
_cmd_js() { eval "$_funcParams2"
  # _npm_install "$bwProjShortcut" || return $?

  local funcName="${bwProjShortcut}_prepare_local"
  if _funcExists "$funcName"; then
    dstVarName=localEnvs srcVarName="_${bwProjShortcut}_localEnvs" eval "$_codeToInitLocalCopyOfArray"
    if [[ ${#localEnvs[@]} -eq 0 ]]; then
      localEnvs=( local testing )
    fi
    if _hasItem "$applicationEnv" "${localEnvs[@]}"; then
      eval "$funcName"
    fi
  fi

  local bwStartProcessFilename="process.yml"
  echo "$yaml" > "$HOME/proj/$pm2Dir/$bwStartProcessFilename"

  echo "Сервис (среда ${_ansiPrimaryLiteral}$applicationEnv${_ansiReset}) будет доступен по следущим URL'ам:"
  echo "  ${_ansiUrl}http://localhost:${_http}${_ansiReset}"
  echo "  ${_ansiUrl}https://localhost:${_https}${_ansiReset}"

  FORCE_COLOR=1 _inDir "$HOME/proj/$pm2Dir" pm2-dev "$bwStartProcessFilename"
}

_declare_test_jsParams=( 'testVars:?' )
_declare_test_js() { eval "$_funcParams2"
  local bwProjShortcut=$(basename "${BASH_SOURCE[1]}" .bash)
  eval "$bwProjShortcut"'_testParamsOpt=(--canBeMixedOptionsAndArgs)'
  eval "$bwProjShortcut"'_testParams=(
    --noCoverage/c
    !--projDir/p=
    "${'"$bwProjShortcut"'_prepareLocalParamsAddon[@]}"
    '\''@testNames:( $(_'"$bwProjShortcut"'_testNames "$projDir") )'\''
  )'
  eval "$bwProjShortcut"'_test_description='\''Прогоняет тест(ы)'\'
  eval "$bwProjShortcut"'_test_noCoverage_description='\''Без покрытия (без istanbul'\'\\\'\''а)'\'
  eval "$bwProjShortcut"'_test() { eval "$_funcParams2"
    eval "$_runInDockerContainer"
    local testContainerDir=""
    local nodeEnv=travis
    local recursive=
    local path="$projDir/node_modules/.bin"
    '"${testVars}"'
    testContainerDir="$HOME/proj/$testContainerDir"
    PATH="$PATH:$path" _cmd_test_js '"$bwProjShortcut"'
  }'
}

_cmd_test_jsParams=( 'bwProjShortcut' )
_cmd_test_js() { eval "$_funcParams2"
  local funcName="_${bwProjShortcut}_prepareLocal"
  ! _funcExists "$funcName" || $funcName || return $?

  local -a testsToRun=()
  if [[ ${#testNames[@]} -eq 0 ]]; then
    [[ -z $recursive ]] || testsToRun+=( --recursive )
    testsToRun+=( 'test' )
  else
    local testName; for testName in "${testNames[@]}"; do
      testsToRun+=( "test/$testName.js" )
    done
  fi
  local -a mocha_OPT=(-t "${timeout:-10000}" -c) # https://www.npmjs.com/package/supports-color
  local -a params=( FORCE_COLOR=1 istanbul cover --report html _mocha -- )
  [[ -z $noCoverage ]] || params=( mocha )
  _inDir "$testContainerDir" _exec -v all NODE_ENV="$nodeEnv" "${params[@]}" "${mocha_OPT[@]}" "${testsToRun[@]}"
}

_cmd_mysql() {
  local returnCode=0
  while true; do
    echo "${_ansiWarn}ВНИМАНИЕ! Для выхода из ${_ansiCmd}mysql${_ansiWarn} используйте команду ${_ansiCmd}quit;${_ansiReset}"
    MYSQL_PWD=root mysql -u root || { returnCode=$?; break; }
    break
  done
  return $returnCode
}

 _declare_mysqlParams=()
 _declare_mysql() { eval "$_funcParams2"
   local bwProjShortcut=$(basename "${BASH_SOURCE[1]}" .bash)
   eval "$bwProjShortcut"'_mysql_description='\''Запускает mysql-клиент'\'
   eval "$bwProjShortcut"'_mysql() { eval "$_funcParams2"
    it=true eval "$_runInDockerContainer"
    _cmd_mysql
   }'
 }

_declare_mysql_mysqldumpParams=()
_declare_mysql_mysqldump() { eval "$_funcParams2"
  local bwProjShortcut=$(basename "${BASH_SOURCE[1]}" .bash)
  local codeToPrepareMysqlHostHolder="_${bwProjShortcut}_codeToPrepareMysqlHost"
  [[ -n ${!codeToPrepareMysqlHostHolder} ]] || return $(_throw "Значение ${_ansiOutline}$codeToPrepareMysqlHostHolder${_ansiErr} не определено")

  eval "$bwProjShortcut"'_mysqlParams=(
    !--projDir/p=
    '\''--env:(local alpha stable)=local'\''
    "${'"$bwProjShortcut"'_mysqlParamsAddon[@]}"
    --execute/e=
    --database/D=
    --skip-column-names/N
    --batch/B
    --auto-rehash/a
    --write
  )'
  eval "$bwProjShortcut"'_mysql_env_description='\''Среда'\'
  eval "$bwProjShortcut"'_mysql_execute_description='\''Опция ${_ansiCmd}--execute${_ansiReset} команды ${_ansiCmd}mysql${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysql_database_description='\''Опция ${_ansiCmd}--database${_ansiReset} команды ${_ansiCmd}mysql${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysql_skipColumnNames_description='\''Опция ${_ansiCmd}--skip-column-names${_ansiReset} команды ${_ansiCmd}mysql${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysql_batch_description='\''Опция ${_ansiCmd}--batch${_ansiReset} команды ${_ansiCmd}mysql${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysql_autoRehash_description='\''Отключает опцию ${_ansiCmd}--no-auto-rehash${_ansiReset} команды ${_ansiCmd}mysql${_ansiReset}, по умолчанию, включенную'\'
  eval "$bwProjShortcut"'_mysql_write_description='\''Зайти под логином с полномочиями на модификацию'\'
  eval "$bwProjShortcut"'_mysql_description='\''Запускает mysql-клиент'\'
  eval "$bwProjShortcut"'_mysql() { eval "$_funcParams2"
    local it; [[ -n $execute ]] || it=true
    queue=$$ it=$it eval "$_runInDockerContainer"
    if [[ $env != local ]]; then
      codeHolder=_'"$bwProjShortcut"'_codeToPrepareMysqlHost eval "$_evalCode"
      codeHolder=_'"$bwProjShortcut"'_codeToPrepareMysqlLogin eval "$_evalCode"
      [[ -n $host ]] || return $(_throw "ожидает, что ${_ansiOutline}_'"$bwProjShortcut"'_codeToPrepareMysqlLogin${_ansiErr} определит значение переменной ${_ansiOutline}host")
    fi
    _mysqlHelper
  }'
  eval "_$bwProjShortcut"'_codeToPrepareMysqlLogin='\''
    local accessLevel
    if [[ -z $write ]]; then
      accessLevel=readonly
    else
      accessLevel=readwrite
    fi
    local srcVarName="_'"$bwProjShortcut"'_${env}_${accessLevel}_login"
    dstVarName=login eval "$_codeToInitLocalCopyOfArray"
    if ! [[ ${#login[@]} -eq 3 && ( ${login[0]} == --user || ${login[0]} == -u ) && ${login[2]} =~ ^-p ]]; then
      _debugVar $srcVarName
      return $(_throw "Ожидает, что массив ${_ansiOutline}$srcVarName${_ansiErr} будет содержать три элемента, первый из которых будет ${_ansiPrimaryLiteral}--user${_ansiErr} или ${_ansiSecondaryLiteral}-u${_ansiErr}, а третий начинается на ${_ansiPrimaryLiteral}-p")
    fi
    local pwd="${login[@]: -1}"
  '\'

  eval "$bwProjShortcut"'_mysqldumpParams=( --canBeMoreParams )'
  eval "$bwProjShortcut"'_mysqldumpParams=(
    !--projDir/p=
    '\''--env:(local alpha stable)=alpha'\''
    "${'"$bwProjShortcut"'_mysqlParamsAddon[@]}"
    --write
    outputFileSpec
    @1..mysqldumpParams
  )'
  eval "$bwProjShortcut"'_mysqldump_env_description='\''Среда'\'
  eval "$bwProjShortcut"'_mysqldump_server_description='\''Сервер объявлений (adv) или брокеров (broker), актуально, если ${_ansiCmd}--env${_ansiReset} не ${_ansiPrimaryLiteral}local${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysqldump_shard_description='\''Номер шарда, актуально только для ${_ansiCmd}--server adv${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysqldump_write_description='\''Зайти под логином с полномочиями на модификацию'\'
  eval "$bwProjShortcut"'_mysqldump_outputFileSpec_description='\''Путь относительно корня проекта к файлу, куда будет записан dump.'\'
  eval "$bwProjShortcut"'_mysqldump_mysqldumpParams_description='\''параметров ${_ansiCmd}mysqldump${_ansiReset}'\'
  eval "$bwProjShortcut"'_mysqldump_description='\''Запускает mysqldump'\'
  eval "$bwProjShortcut"'_mysqldump() { eval "$_funcParams2"
    queue=$$ eval "$_runInDockerContainer"
    if [[ $env != local ]]; then
      codeHolder=_'"$bwProjShortcut"'_codeToPrepareMysqlHost eval "$_evalCode"
      codeHolder=_'"$bwProjShortcut"'_codeToPrepareMysqlLogin eval "$_evalCode"
    fi
    _mysqldumpHelper
  }'
}
_mysqlHelper() {
  local -a OPT_noAutoRehash=(); [[ -z $execute && -z $autoRehash ]] && OPT_noAutoRehash=( --no-auto-rehash )
  local -a mysqlParams=( "${OPT_database[@]}" "${OPT_skipColumnNames[@]}" "${OPT_noAutoRehash[@]}" "${OPT_batch[@]}" "${OPT_execute[@]}" )

  if [[ -z $execute ]]; then
    echo -n "${_nl}${_ansiWarn}ВНИМАНИЕ! Вы входите на сервер ${_ansiCmd}mysql:"
    if [[ $env == local ]]; then
      echo "root@localhost"
    else
      echo "${login[1]}@$host"
      echo -n "${_ansiWarn}У Вас есть полномочия "
      if [[ -z $write ]]; then
        echo -n "только на ${_ansiOK}чтение"
      else
        echo -n "на ${_ansiOK}чтение${_ansiWarn} и ${_ansiErr}модификацию"
      fi
      echo "${_ansiReset}"
    fi
    echo "${_ansiWarn}Для выхода используйте команду ${_ansiCmd}quit;${_ansiReset}"
  fi

  local -a mysql_OPT
  local mysqlPwd
  if [[ $env == local ]]; then
    mysqlPwd=root
    mysql_OPT=( -u root )
  else
    mysqlPwd=${pwd:2}
    mysql_OPT=( --host $host "${login[@]:0:2}")
  fi

  MYSQL_PWD="$mysqlPwd" _exec -v err mysql "${mysql_OPT[@]}" "${mysqlParams[@]}"
}
_mysqldumpHelper() {
  if [[ $outputFileSpec =~ ^/ ]]; then
    return $(errOrigin=1 _throw "ожидает, что параметр ${_ansiOutline}outputFileSpec${_ansiErr} будет путем ${_ansiOutline}относительно${_ansiErr} корня проекта")
  else
    local outputFileDir=$(dirname "$projDir/$outputFileSpec")
    _mkDir -v err "$outputFileDir" || return $?
  fi
  local -a mysqldump_OPT
  local mysqlPwd
  if [[ $env == local ]]; then
    mysqlPwd=root
    mysqldump_OPT=( -u root )
  else
    mysqlPwd="${pwd:2}"
    mysqldump_OPT=( --host $host "${login[@]:0:2}" )
  fi
  local -a mysqldump=(mysqldump --skip-lock-tables "${mysqldump_OPT[@]}" "${mysqldumpParams[@]}")
  MYSQL_PWD="$mysqlPwd" _spinner \
    -t "Выполнение заняло" \
    "$(_quotedArgs "${mysqldump[@]}" ) > $(_quotedArgs "$outputFileSpec")" \
    _exec -v err --stdout "$projDir/$outputFileSpec" "${mysqldump[@]}"
}

_exec_sqlParamsOpt=( '--canBeMoreParams' )
_exec_sqlParams=(
  'scriptName:( $(find "$HOME/proj/docker/mysql" -maxdepth 1 -name "*.sql" -exec basename {} .sql \; ) )'
)
_exec_sql() { eval "$_funcParams2"
  local fileSpec="$HOME/proj/docker/mysql/$scriptName.sql"
  if [[ -s $fileSpec ]]; then
    local sqlCommands
    read -d $'\x04' sqlCommands < "$fileSpec"
    MYSQL_PWD=root _inDir "$HOME/proj/docker" _exec -v allBrief mysql -u root -e "${_nl}$sqlCommands${_nl}" "$@"
  fi
}

_clear_mysql_dbParams=( 'dbName' )
_clear_mysql_db() { eval "$_funcParams2"
  # https://dev.mysql.com/doc/refman/8.0/en/drop-database.html
  local scriptName="clear_database.temp"
  local fileSpec="$HOME/proj/docker/mysql/$scriptName.sql"
  MYSQL_PWD=root mysql -u root -D "$dbName" -BNe 'show tables' | awk '{print "set foreign_key_checks=0; drop table if exists `" $1 "`;"}' > "$fileSpec"
  if [[ -s $fileSpec ]]; then
    _exec_sql "$scriptName" -D "$dbName"
  fi
  rm -f "$fileSpec"
}

_start_mysql() {
  local returnCode=0
  while true; do
    if [[ ! -d /var/lib/mysql/mysql ]]; then
      _exec -v all --sudo cp -r /var/lib/_mysql/. /var/lib/mysql/ || { returnCode=$?; break; }
      _exec -v all --sudo chown -R mysql:mysql /var/lib/mysql /var/run/mysqld || { returnCode=$?; break; }
    fi
    _exec -v all --sudo chown -R mysql:adm /var/log/mysql || { returnCode=$?; break; }
    _exec --sudo service mysql start || { returnCode=$?; break; } # https://askubuntu.com/questions/2075/whats-the-difference-between-service-and-etc-init-d
    local sqlCommands
    read -d $'\x04' sqlCommands < "$_bwDir/docker/helper/mysql/mysql_secure_installation.sql"
    MYSQL_PWD=root _inDir "$HOME/proj/docker" _exec -v allBrief mysql -u root -e "${_nl}$sqlCommands${_nl}" "$@"
    mysql_tzinfo_to_sql /usr/share/zoneinfo | MYSQL_PWD=root mysql -u root mysql
    [[ ! -f $HOME/proj/docker/mysql/init_database.sql ]] || _exec_sql init_database || { returnCode=$?; break; }
    break
  done
  return $returnCode
}

_npm_install() {
  if [[ $# -eq 0 ]]; then
    _npm_installHelper
  else
    local dir; for dir in "$@"; do
      _npm_installHelper "$dir" || return $?
    done
  fi
}
_npm_installHelper() {
  local dir="$1"
  npm set color always
  local cmdTitle="${_ansiCmd}"
  [[ -z $dir ]] || cmdTitle+="$dir: "
  cmdTitle+="npm install${_ansiReset}"
  dir="$HOME/proj/$dir"
  echo "$cmdTitle . . ."
  _inDir "$dir" npm install; local returnCode=$?
  [[ $returnCode -ne 0 ]] || _ok "$cmdTitle"
  return $returnCode
}

# shellcheck disable=SC2034
_docker_buildParams=(
  '--force'
  'bwProjShortcut'
  'dockerImageName'
  'dockerImageTitle'
)
_docker_build() { eval "$_funcParams2"
  local msg=
  if [[ -z $force ]] && ! _gitIsChangedFile 'docker/Dockerfile'; then
    msg+="Перед сборкой образа необходимо внести изменения в ${_ansiFileSpec}$(_shortenFileSpec "$(pwd)/Dockerfile")${_ansiReset} или выполнить команду с опцией ${_ansiCmd}--force"
    _warn "$msg"
  else
    _docker -v all build --pull -t "$dockerImageName" - < Dockerfile; local returnCode=$?
    if [[ $returnCode -eq 0 ]]; then
      local dockerImageIdFileName="dev-$bwProjShortcut.id"
      _docker inspect --format "{{json .Id}}" "$dockerImageName:latest" > "$dockerImageIdFileName"
      if _gitIsChangedFile "docker/$dockerImageIdFileName"; then
        msg+="Обновлен docker-образ $dockerImageTitle${_ansiWarn}"$_nl
        msg+="${_ansiWarn}Не забудьте поместить его в docker-репозиторий командой"$_nl
        msg+="    ${_ansiCmd}$bwProjShortcut docker push${_ansiReset}"
        _warn "$msg"
      fi
    fi
  fi
  return $returnCode
}

# shellcheck disable=SC2034
_docker_downParams=( 'projDir' )
_docker_down() { eval "$_funcParams2"
  if _exist "$projDir/docker/docker-compose.yml"; then
    local -a dockerCompose_OPT; _prepareDockerComposeOpt
    _dockerCompose "${dockerCompose_OPT[@]}" down --remove-orphans
    _bwCleanTempDockerFiles
  fi
}

_bwCleanTempDockerFiles() {
  rm -f *.pid *.queue
  rm -rf whoami
  rm -f mysql/*.temp.sql
}

_prepareDockerComposeOpt() {
  local dockerComposeProjectName; dockerComposeProjectName="$(_shortenFileSpec "$projDir")"
  [[ $dockerComposeProjectName =~ ^~ ]] && dockerComposeProjectName=${dockerComposeProjectName:1}
  dockerCompose_OPT=( -v all -p "$dockerComposeProjectName" )
}

# shellcheck disable=SC2034
_cmd_selfTestParams() {
  _cmd_selfTestParams=()
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    _cmd_selfTestParams+=( "--$port=" )
  done
  _cmd_selfTestParams+=(
    '--auto'
    '--dockerImageName='
    'bwProjShortcut'
    'projDir'
  )
}
_cmd_selfTest() { eval "$_funcParams2"
  if [[ -n $auto ]]; then
    local -a OPT=()
    local port; for port in "${_bwProjPorts[@]}" upstream; do
      dstVarName=_OPT srcVarName="OPT_${port}" eval "$_codeToInitLocalCopyOfArray"
      OPT+=( "${_OPT[@]}" )
    done
    bw_install --silentIfAlreadyInstalled expect
    local patch=''
    if [[ $bwProjShortcut == agate ]]; then
      patch='
      -ex "passing" {
      sleep 1
      send \x03
    }'
    fi
    expect -c '
set timeout -1

spawn bash -c ". '"$_bwFileSpec"' -p -; . ../bin/'"$bwProjShortcut"'.bash; '"$bwProjShortcut"' self-test '"$(_quotedArgs "${OPT[@]}")"'"

while {1} {
  expect {
    eof { break }
    -ex "\[97m\[1mquit;\[0m" {
      sleep 3
      send -- "quit;\r"
    }
    -ex "\[97m\[1mexit 0\[0m" {
      sleep 10
      send -- "YeS\r"
      sleep 1
      send -- "exit 0\r"
    }
    -ex "\[97m\[1mq\r" {
      sleep 1
      send -- "q\r"
    }
    -ex "Password:" {
      if { [info exists env(ROOT_PWD)] } {
        send "$env(ROOT_PWD)\r"
      } else {
        stty -echo
        expect_user -timeout -1 -re "(.*)\[\r\n]"
        stty echo
        send "$expect_out(1,string)\r"
      }
    }
    -ex "sudo] password for" {
      if { [info exists env(ROOT_PWD)] } {
        send "$env(ROOT_PWD)\r"
      } else {
        stty -echo
        expect_user -timeout -1 -re "(.*)\[\r\n]"
        stty echo
        send "$expect_out(1,string)\r"
      }
    }
    '"$patch"'
  }
}'
    return $?
  fi
  local returnCode=0
  while true; do

    _cmd_selfTestHelper -F || { returnCode=$?; break; }
    _cmd_selfTestHelper -f || { returnCode=$?; break; }

    # local -a OPT=()
    # local port; for port in "${_bwProjPorts[@]}" upstream; do
    #   dstVarName=_OPT srcVarName="OPT_${port}" eval "$_codeToInitLocalCopyOfArray"
    #   OPT+=( "${_OPT[@]}" )
    # done
    # _exec -v all "$bwProjShortcut" docker up -f "${OPT[@]}" || { returnCode=$?; break; }
    # _mkDir $projDir/tmp
    # local tstFileSpec="$projDir/tmp/$bwProjShortcut.selfTest"
    # if [[ -n $http ]]; then
    #   _exec -v err curl -s "http://localhost:${http}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
    #   _exec -v all diff "$tstFileSpec" "nginx/whoami/index.html" || { returnCode=$?; break; }
    #   rm "$tstFileSpec"
    # fi
    # if [[ -n $https ]]; then
    #   _exec -v err curl -s "https://localhost:${https}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
    #   _exec -v all diff "$tstFileSpec" "nginx/whoami/index.html" || { returnCode=$?; break; }
    #   rm "$tstFileSpec"
    # fi
    # local skipFail=true # for debug purpose only
    # if [[ -n $mysql ]]; then
    #   _exec -v all "${bwProjShortcut}" mysql --env local -BN -e 'create database if not exists temp;' || { returnCode=$?; break; }
    #   _exec -v all "${bwProjShortcut}" mysqldump --env local tmp/temp.dump temp || { returnCode=$?; break; }
    #   _exec -v all "${bwProjShortcut}" mysql --env alpha -BN -e 'show databases;' || [[ -n $skipFail ]] || { returnCode=$?; break; }
    # fi
    # local funcName="_${bwProjShortcut}_selfTestAddon"
    # if _funcExists "$funcName"; then
    #   "$funcName" || { returnCode=$?; break; }
    # fi
    # if [[ -n $mysql ]]; then
    #   _exec -v all "${bwProjShortcut}" mysql || { returnCode=$?; break; }
    # fi
    # if [[ -n $ssh ]]; then
    #   echo "${_nl}${_ansiWarn}ВНИМАНИЕ! Для выхода из ssh-сессии выполните команду ${_ansiCmd}exit 0${_ansiReset}"
    #   _exec -v all ssh -p "$ssh" dev@localhost || { returnCode=$?; break; }
    # fi
    # if [[ -n $dockerImageName ]]; then
    #   _exec -v all "$bwProjShortcut" docker shell || { returnCode=$?; break; }
    # fi
    # if [[ -n $http || -n $https ]]; then
    #   echo "${_nl}${_ansiWarn}ВНИМАНИЕ! Для выхода из docker-контейнера выполните команду ${_ansiCmd}exit 0${_ansiReset}"
    #   _exec -v all "$bwProjShortcut" docker shell "nginx" || { returnCode=$?; break; }
    # fi
    # _exec -v all "$bwProjShortcut" docker down || { returnCode=$?; break; }

    break
  done
  if [[ $returnCode -eq 0 ]]; then
    _ok "self-test complete"
  else
    _err "self-test failed"
  fi
  return $returnCode
}
_cmd_selfTestHelper() {
  local docker_up_OPT=$1
  local returnCode=0
  while true; do
    local -a OPT=()
    local port; for port in "${_bwProjPorts[@]}" upstream; do
      dstVarName=_OPT srcVarName="OPT_${port}" eval "$_codeToInitLocalCopyOfArray"
      OPT+=( "${_OPT[@]}" )
    done
    _exec -v all "$bwProjShortcut" docker up $docker_up_OPT "${OPT[@]}" || { returnCode=$?; break; }
    _mkDir $projDir/tmp
    local tstFileSpec="$projDir/tmp/$bwProjShortcut.selfTest"
    if [[ -n $http ]]; then
      _exec -v err curl -s "http://localhost:${http}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
      _exec -v all diff "$tstFileSpec" "nginx/whoami/index.html" || { returnCode=$?; break; }
      rm "$tstFileSpec"
    fi
    if [[ -n $https ]]; then
      _exec -v err curl -s "https://localhost:${https}/whoami/" > "$tstFileSpec" || { returnCode=$?; break; }
      _exec -v all diff "$tstFileSpec" "nginx/whoami/index.html" || { returnCode=$?; break; }
      rm "$tstFileSpec"
    fi
    local skipFail=true # for debug purpose only
    if [[ -n $mysql ]]; then
      _exec -v all "${bwProjShortcut}" mysql --env local -BN -e 'create database if not exists temp;' || { returnCode=$?; break; }
      _exec -v all "${bwProjShortcut}" mysqldump --env local tmp/temp.dump temp || { returnCode=$?; break; }
      _exec -v all "${bwProjShortcut}" mysql --env alpha -BN -e 'show databases;' || [[ -n $skipFail ]] || { returnCode=$?; break; }
    fi
    local funcName="_${bwProjShortcut}_selfTestAddon"
    if _funcExists "$funcName"; then
      "$funcName" || { returnCode=$?; break; }
    fi
    if [[ -n $mysql ]]; then
      _exec -v all "${bwProjShortcut}" mysql || { returnCode=$?; break; }
    fi
    if [[ -n $ssh ]]; then
      echo "${_nl}${_ansiWarn}ВНИМАНИЕ! Для выхода из ssh-сессии выполните команду ${_ansiCmd}exit 0${_ansiReset}"
      _exec -v all ssh -p "$ssh" dev@localhost || { returnCode=$?; break; }
    fi
    if [[ -n $dockerImageName ]]; then
      _exec -v all "$bwProjShortcut" docker shell || { returnCode=$?; break; }
    fi
    if [[ -n $http || -n $https ]]; then
      echo "${_nl}${_ansiWarn}ВНИМАНИЕ! Для выхода из docker-контейнера выполните команду ${_ansiCmd}exit 0${_ansiReset}"
      _exec -v all "$bwProjShortcut" docker shell "nginx" || { returnCode=$?; break; }
    fi
    _exec -v all "$bwProjShortcut" docker down || { returnCode=$?; break; }
    break
  done
  return $returnCode
}


_selfTestParamsOpt=( --canBeMixedOptionsAndArgs )
_selfTestParams=(
  '--queue=default'
  '--sleep:0..=10'
  'runCommand'
  'testCommands'
)
_selfTest() { eval "$_funcParams2"
  local queueFileName="$queue.queue"
  rm -f "$queueFileName"
  ( _exec -v all $runCommand & )
  sleep $sleep
  _exist "$(pwd)/$queueFileName" || return $?
  local returnCode=0
  while true; do
    eval "$testCommands"
    break
  done
  local pid; read -r -d $'\x04' pid < "$queueFileName"
  local -a dockerCompose_OPT; _prepareDockerComposeOpt
  _spinner \
    -t "Отправка сигнала SIGINT в docker-контейнер заняла" \
    "Отправка сигнала SIGINT в docker-контейнер" \
    _dockerCompose "${dockerCompose_OPT[@]:2}" exec -T main \
      kill -SIGINT "$pid"
  return $returnCode
}

# shellcheck disable=SC2034
_docker_upParams() {
  _docker_upParams=()
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    _docker_upParams+=( "--$port=" )
  done
  _docker_upParams+=(
    '--force-recreate/f'
    '--force-recreate-and-clean-mapped-volumes/F'
    '--dockerImageName='
    '--bwProjShortcut='
    '--bwProjName='
    '@--restart'
    '--noCheck'
    '--noTestAccessMessage/m'
    '@1..dockerComposeFileNames'
  )
}
_docker_up() { eval "$_funcParams2"
  local returnCode=0
  while true; do

    if [[ -n $https ]]; then
      bw_install --silentIfAlreadyInstalled root-cert || { returnCode=$?; break; }
    fi

    if [[ $OSTYPE =~ ^linux ]]; then
      # https://github.com/guard/listen/wiki/Increasing-the-amount-of-inotify-watchers
      local line=fs.inotify.max_user_watches=524288
      local fileSpec=/etc/sysctl.conf
      if ! _silent grep -Fx "$line" "$fileSpec"; then
        _exec -v all --cmdAsIs "echo \"$line\" | sudo tee -a \"$fileSpec\" && sudo sysctl -p" || { returnCode=$?; break; }
      fi
    fi

    if [[ -n $https && -n $http && $https -eq $http ]]; then
      errOrigin=1 _throw "ожидает, что значение ${_ansiOutline}http ${_ansiPrimaryLiteral}$http${_ansiErr} будет отличаться от значения ${_ansiOutline}https"
      returnCode=$?; break
    fi

    bw_install --silentIfAlreadyInstalled docker || { returnCode=$?; break; }

    if [[ -n $dockerImageName && -z $noCheck ]]; then
      local dockerImageLs; dockerImageLs=$(_docker image ls "$dockerImageName:latest" -q) || { returnCode=$?; break; }
      if [[ -z $dockerImageLs ]]; then
        _docker -v all image pull "$dockerImageName:latest" || { returnCode=$?; break; }
      fi
      local tstImageIdFileSpec="/tmp/$bwProjShortcut.image.id"
      _docker inspect --format "{{json .Id}}" "$dockerImageName:latest" > "$tstImageIdFileSpec"
      while true; do
        local dockerImageIdFileName="dev-$bwProjShortcut.id"
        if ! _silent cmp "$tstImageIdFileSpec" "$dockerImageIdFileName"; then
          local needWarn="" msg=""
          if _gitIsChangedFile 'docker/Dockerfile'; then
            needWarn=true
          else
            _docker -v all image pull "$dockerImageName:latest" || { returnCode=$?; break; }
            _docker inspect --format "{{json .Id}}" "$dockerImageName:latest" > "$tstImageIdFileSpec"
            if ! _silent cmp "$tstImageIdFileSpec" "$dockerImageIdFileName"; then
              msg+="Идентификатор ${_ansiPrimaryLiteral}$(cat "$tstImageIdFileSpec")${_ansiReset} образа ${_ansiSecondaryLiteral}$dockerImageName:latest${_ansiReset} "
              msg+="отличается от эталонного значения ${_ansiPrimaryLiteral}$(cat "$dockerImageIdFileName")${_ansiReset}, хранящегося в файле ${_ansiFileSpec}$(_shortenFileSpec "$(pwd)/$dockerImageIdFileName")"$_nl
              needWarn=true
            fi
          fi
          if [[ -n $needWarn ]]; then
            msg+="Чтобы избавиться от этого сообщения, необходимо выполнить:"$_nl
            msg+="  ${_ansiCmd}$bwProjShortcut docker build --force${_ansiReset}"$_nl
            msg+="  ${_ansiCmd}$bwProjShortcut docker push${_ansiReset}"$_nl
            msg+="  ${_ansiCmd}git add $(pwd)/Dockerfile $(pwd)/$dockerImageIdFileName${_ansiReset}"$_nl
            msg+="  ${_ansiCmd}git commit${_ansiReset}"
            _warn "$msg"
          fi
        fi
        break
      done
      rm "$tstImageIdFileSpec"
      [[ $returnCode -eq 0 ]] || break
    fi

    local -a _prepareDockerComposeYml_OPT=(
      "${OPT_dockerImageName[@]}"
      "${OPT_bwProjShortcut[@]}"
      "${OPT_dockerImageName[@]}"
    )
    local port; for port in "${_bwProjPorts[@]}" upstream; do
      dstVarName=_OPT srcVarName="OPT_${port}" eval "$_codeToInitLocalCopyOfArray"
      _prepareDockerComposeYml_OPT+=( "${_OPT[@]}" )
    done
    _prepareDockerComposeYml "${_prepareDockerComposeYml_OPT[@]}" "${dockerComposeFileNames[@]}" || { returnCode=$?; break; }

    local -a dockerCompose_OPT; _prepareDockerComposeOpt
    if [[ "${#restart[@]}" -gt 0 ]]; then
      _dockerCompose "${dockerCompose_OPT[@]}" stop "${restart[@]}" || { returnCode=$?; break; }
    fi

    [[ -z $forceRecreate && -z $forceRecreateAndCleanMappedVolumes ]] || OPT_forceRecreate=( '--force-recreate' )
    if [[ -n $forceRecreateAndCleanMappedVolumes ]]; then
      _dockerCompose "${dockerCompose_OPT[@]}" down
      local funcName="_${bwProjShortcut}_clean"
      $funcName || { returnCode=$?; break; }
    fi
    # https://stackoverflow.com/questions/692000/how-do-i-write-stderr-to-a-file-while-using-tee-with-a-pipe
    local fileSpec="/tmp/$$.docker_up.log"
    _dockerCompose "${dockerCompose_OPT[@]}" up -d "${OPT_forceRecreate[@]}" --remove-orphans 2> >(tee "$fileSpec" >&2) || { returnCode=$?; }
    # sudo docker-compose "${dockerCompose_OPT[@]}" up -d "${OPT_forceRecreate[@]}" --remove-orphans 2> >(tee "$fileSpec" >&2) || { returnCode=$?; }
    if [[ $returnCode -ne 0 ]]; then
      local awkFileSpec; _prepareAwkFileSpec
      local -a awk_OPT=(
        -f "$awkFileSpec"
      )
      local alreadyAllocatedPort
      if alreadyAllocatedPort=$(awk "${awk_OPT[@]}" "$fileSpec"); then
        local optName=
        local port; for port in "${_bwProjPorts[@]}"; do
          if [[ $alreadyAllocatedPort -eq ${!port} ]]; then
            optName=$port
          fi
        done
        if [[ -n $optName ]]; then
          local msg=''
          msg+="${_nl}${_ansiWarn}ВНИМАНИЕ!"
          msg+="${_nl}  ${_ansiSecondaryLiteral}$optName${_ansiWarn}-порт ${_ansiPrimaryLiteral}${!port}${_ansiWarn} занят"
          msg+="${_nl}  Пользуясь опцией ${_ansiCmd}--$optName${_ansiWarn}, укажите другоe значение для этого порта:"
          msg+="${_nl}    ${_ansiCmd}--$optName ${_ansiOutline}Порт${_ansiWarn}"
          msg+="${_nl}  или опцией ${_ansiCmd}--port-increment${_ansiWarn} ( ${_ansiCmd}-i${_ansiWarn} ), укажите смещение относительно базовых значений для всех портов:"
          msg+="${_nl}    ${_ansiCmd}-i ${_ansiOutline}Смещение"
          msg+="${_nl}  ${_ansiReset}"
          echo "$msg"
        fi
      fi
    elif [[ -n $forceRecreate ]]; then
      _bwCleanTempDockerFiles
    fi
    rm -f "$fileSpec"
    [[ $returnCode -eq 0 ]] || break

    if [[ -n $dockerImageName ]]; then
      local -a dockerCompose_OPT; _prepareDockerComposeOpt
      local -a OPT_T; [[ -n $it ]] || OPT_T=( -T )
      local -a cmd_base=( _dockerCompose "${dockerCompose_OPT[@]:2}" exec "${OPT_T[@]}" main )
      # local -a cmd=( "${cmd_base[@]}" "$_bwDevDockerEntryPointFileSpec" "$$" "$@" )

      local -a cmd=( "${cmd_base[@]}" sudo usermod -u $(id -u) dev)
      eval "$(_quotedArgs "${cmd[@]}")" || { returnCode=$?; break; }

      # if [[ -n $it ]]; then
      #   eval "$(_quotedArgs "${cmd[@]}")"
      # else
        # # (
        # #   # local subPid=$!
        # #   # local pidFileSpec="$$.pid"
        # #   # # shellcheck disable=SC2064
        # #   # trap "trap - SIGINT; _runInDockerCtrlC $subPid \"$pidFileSpec\" \"$*\" $(_quotedArgs "${cmd_base[@]}")" SIGINT
        # #   # wait $subPid; local returnCode=$?
        # #   # trap - SIGINT
        # #   # return $returnCode
        # ) || { returnCode=$?; break; }
      # fi
      # _inDir -v none "$projDir/docker"
      _runInDockerContainer_outOfDockerHelper || { returnCode=$?; break; }
    fi

    if [[ -n $http || -n $https ]]; then
      _mkDir nginx/whoami
      {
      echo -n "<pre>"
      echo -n "
  projName: ${bwProjName}
  projShortcut: ${bwProjShortcut}
  projDir: $(_shortenFileSpec $(cd .. && pwd))"
  local port; for port in "${_bwProjPorts[@]}"; do
       [[ -z ${!port} ]] || echo -n "
  $port: ${!port}"
  done
      echo "${_nl}</pre>"
      } > nginx/whoami/index.html
    fi

    if [[ ( -n $http || -n $https || -n $ssh ) && -z $noTestAccessMessage ]]; then
      echo "Проверка доступности docker-приложения:"
      [[ -z $http ]] || echo "  ${_ansiUrl}http://localhost:$http/whoami/${_ansiReset}"
      [[ -z $https ]] || echo "  ${_ansiUrl}https://localhost:$https/whoami/${_ansiReset}"
      [[ -z $ssh ]] || echo "  ${_ansiCmd}ssh -p $ssh dev@localhost${_ansiReset}"
    fi

    if [[ -n $ssh ]]; then
      _prepareSsh
    fi

    break
  done

  return $returnCode
}
_prepareProjPrompt() {
  local -a OPT=(
    --additionalDependencies "$_bwFileSpec"
    --additionalDependencies "$_bwDir/bash/psSupport.bash"
  )
  _useCache "${OPT[@]}" "$promptHolder"; local returnCode=$?; [[ $returnCode -eq 2 ]] || return $returnCode
  local prompt; _preparePrompt --user "$bwProjShortcut" --userAnsi PrimaryLiteral --gitDirty - || return $?
  eval "$promptHolder"'="$prompt"'
  _saveToCache "$promptHolder"
}
_prepareSsh() {
  if [[ ! -f ~/.ssh/bw_dev_id_rsa ]]; then
    local returnCode=0
    while true; do
      _exec -v err cp "$_bwDir/ssh/bw_dev_id_rsa" ~/.ssh || { returnCode=$?; break; }
      _exec -v err chmod 600 ~/.ssh/bw_dev_id_rsa || { returnCode=$?; break; }
      _exec -v err --append ~/.ssh/config echo "
# added by bw.bash (_prepareSsh@$_bwDir/bash/bwCommands.bash)
Host localhost
User dev
IdentityFile ~/.ssh/bw_dev_id_rsa" || { returnCode=$?; break; }
      break
    done
    if [[ $returnCode -ne 0 ]]; then
      _warn "Могут быть проблемы с ssh-доступом к docker-контейнеру"
    fi
  fi
}

# shellcheck disable=SC2034
_docker_shellParams=(
  'serviceName'
  'shell'
)
_docker_shell() { eval "$_funcParams2"
  local -a OPT=()
  if [[ $serviceName == "main" && $shell == /bin/bash ]]; then
    OPT+=( --init-file "$_bwDevDockerEntryPointFileSpec" )
  fi
  local -a dockerCompose_OPT; _prepareDockerComposeOpt
  # _dockerCompose "${dockerCompose_OPT[@]}" exec "$serviceName" "$shell" "${OPT[@]}"
  _dockerCompose "${dockerCompose_OPT[@]}" exec "$serviceName" "$shell" "${OPT[@]}"
}

# shellcheck disable=SC2034
{
_dockerComposeContainerNamesParamsOpt=( '--canBeMoreParams' )
_dockerComposeContainerNamesParams=(
  '--varName/v='
  'bwProjShortcut'
)
}
_dockerComposeContainerNames() { eval "$_funcParams2"
  local projDir="$*"
  _prepareProjDir "$bwProjShortcut" || return $?
  [[ -f $projDir/docker/docker-compose.yml ]] \
    || "_${bwProjShortcut}_prepareDockerComposeYml" -p "$projDir" \
    || return $?
  if [[ -n $varName ]]; then
    mapfile -t "$varName" < <(_dockerComposeContainerNamesHelper)
  else
    local -a result; mapfile -t result < <(_dockerComposeContainerNamesHelper)
    echo "${result[@]}"
  fi
}
_dockerComposeContainerNamesHelper() {
  sed -nEe 's/^[ \t]*container_name:[ \t]*([^ \t]*)/\1/p' "$projDir/docker/docker-compose.yml"
}

# shellcheck disable=SC2034
{
_dockerComposeServiceNamesParamsOpt=( '--canBeMoreParams' )
_dockerComposeServiceNamesParams=(
  '--varName/v='
  'bwProjShortcut'
)
}
_dockerComposeServiceNames() { eval "$_funcParams2"
  local projDir="$*"
  _prepareProjDir "$bwProjShortcut" || return $?
  [[ -f $projDir/docker/docker-compose.yml ]] \
    || "_${bwProjShortcut}_prepareDockerComposeYml" -p "$projDir" \
    || return $?
  if [[ -n $varName ]]; then
    mapfile -t "$varName" < <(_dockerComposeServiceNamesHelper)
  else
    local -a result; mapfile -t result < <(_dockerComposeServiceNamesHelper)
    echo "${result[@]}"
  fi
}
_dockerComposeServiceNamesHelper() {
  _dockerCompose -f  "$projDir/docker/docker-compose.yml" config --services
}

_defaultDockerComposeServiceNameParamsOpt=( '--canBeMoreParams' )
_defaultDockerComposeServiceNameParams=( 'bwProjShortcut' )
_defaultDockerComposeServiceName() { eval "$_funcParams2"
  local projDir="$*"
  local -a serviceNames
  _dockerComposeServiceNames -v serviceNames "$bwProjShortcut" "$projDir" || return $?
  if _hasItem main serviceNames; then
    echo main
  else
    echo "${serviceNames[0]}"
  fi
}

# shellcheck disable=SC2034
_prepareDockerComposeYmlParams() {
  _prepareDockerComposeYmlParams=()
  local port; for port in "${_bwProjPorts[@]}" upstream; do
    _prepareDockerComposeYmlParams+=( "--$port=" )
  done
  _prepareDockerComposeYmlParams+=(
    '--dockerImageName='
    '--bwProjShortcut='
    '@1..dockerComposeFileNames'
  )
}
_prepareDockerComposeYml() { eval "$_funcParams2"
  local returnCode=0
  while true; do

    if [[ -n $dockerImageName ]]; then
      local promptHolder="_${bwProjShortcut}Prompt"
      _prepareProjPrompt || return $?
      local dockerContainerEnvFileName="main.env"
      {
        echo "#!/bin/bash"
        echo "# file generated by ${bwProjShortcut}_docker_up"
        echo "export BW_SELF_UPDATE_SOURCE=$BW_SELF_UPDATE_SOURCE"
        echo "export _bwProjName=$bwProjName"
        echo "export _bwProjShortcut=$bwProjShortcut"
        echo "export _hostUser=$(whoami)"
        local port; for port in "${_bwProjPorts[@]}"; do
          [[ -n ${!port} ]] || continue
          echo "export _${port}=${!port}"
        done
        echo "export _prompt='${!promptHolder}'"

        local addonFuncName="_${bwProjShortcut}_docker_upAddon"
        if _funcExists "$addonFuncName"; then
          $addonFuncName
        fi
      } > "$dockerContainerEnvFileName"
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
      local nginxDir="nginx"
      if [[ -n $needProcessNginxConfTemplate ]]; then
        local separatorLine='# !!! you SHOULD put all custom rules above this line; all lines below will be autocleaned'
        if [[ -f "$nginxDir/.gitignore" ]]; then
          awk "!a;/$separatorLine/{a=1}" "$nginxDir/.gitignore" > "$nginxDir/.gitignore.new"
        elif [[ -f "$nginxDir/.gitignore.new" ]]; then
          echo "$separatorLine" > "$nginxDir/.gitignore.new"
        fi
        local projName="$bwProjName"
        local projShortcut="$bwProjShortcut"
        local templateExt='.template'
        local -a templateFileSpecs; mapfile -t templateFileSpecs < <(find "$nginxDir" -name "*$templateExt" | sort)
        local templateFileSpec; for templateFileSpec in "${templateFileSpecs[@]}"; do
          fileSpec=${templateFileSpec:0:$(( ${#templateFileSpec} - ${#templateExt} ))}
          local relativeFileSpec="${fileSpec:$(( ${#nginxDir} + 1 ))}" # должно быть именно до вызова _mkFileFromTemplate, потому что в _mkFileFromTemplate есть баг, который приводит к тому, что после вызова  _mkFileFromTemplate templateFileSpec становится пустым
          if _mkFileFromTemplate "$fileSpec"; then
            echo "$relativeFileSpec" >> "$nginxDir/.gitignore.new"
          fi
        done
        if [[ -f "$nginxDir/.gitignore.new" ]]; then
          mv "$nginxDir/.gitignore.new" "$nginxDir/.gitignore"
        fi
      fi
    fi

    local -a config_OPT=()
    dockerComposeFileNames+=( "docker-compose.proj.yml" )
    local dockerComposeFileName; for dockerComposeFileName in "${dockerComposeFileNames[@]}"; do
      if [[ $dockerComposeFileName != 'docker-compose.proj.yml' ]]; then
        local srcFileSpec="$_bwDir/docker/helper/$dockerComposeFileName"
        _exec -v err cp -f "$srcFileSpec" "$dockerComposeFileName" || { returnCode=$?; break; }
        if [[ $dockerComposeFileName == 'docker-compose.mysql.yml' ]]; then
          local myCnfFileSpec="${_bwDir}/docker/helper/mysql/my.cnf"
          [[ -f mysql/my.cnf ]] && myCnfFileSpec=$(pwd)/mysql/my.cnf
          _debugVar myCnfFileSpec
        fi
      fi
      dockerComposeFileName="$(basename "$dockerComposeFileName" .yml)"

      # shellcheck disable=SC2034,SC2086
      {
      local mainContainerName; mainContainerName=$(_getDockerMainContainerName $bwProjShortcut $projDir)
      local nginxContainerName="$mainContainerName-nginx"
      local mainImageName="$dockerImageName"
      local hostUserId=$(id -u)
      }

      _mkFileFromTemplate -n -e .yml "$dockerComposeFileName"

      config_OPT+=( -f "$projDir/docker/$dockerComposeFileName" )
    done

    _dockerCompose -v err "${config_OPT[@]}" config > 'docker-compose.yml' || { returnCode=$?; break; }
    break
  done
  return $returnCode
}

# shellcheck disable=SC2034
{
_getDockerMainContainerNameParamsOpt=( '--canBeMoreParams' )
_getDockerMainContainerNameParams=(
  'bwProjShortcut'
)
}
_getDockerMainContainerName() { eval "$_funcParams2"
  local projDir="$*"
  _prepareProjDir "$bwProjShortcut" || return $?
  echo -n dev-
  _shortenFileSpec "$projDir" | tr -dC 'a-zA-Z0-9_.-'
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_projectInfoParamsOpt=( '--canBeMixedOptionsAndArgs' )
bw_projectInfoShortcuts=( 'pi' )
bw_projectInfo_bwProjShortcut_name='Сокращенное-имя-проекта'
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
bw_projectInfo_description='Выводит информацию о проекте/проектах'
}
bw_projectInfo() { eval "$_funcParams2"
  local skipNonExistent=
  if [[ -n $def ]]; then
    local -a duplicatePorts; _prepareDuplicatePorts || return $?
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
      echo "Всего $(_getPluralWord ${#found[@]} определен определено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}$(_joinBy ", " "${found[@]}")${_ansiReset}"
    else
      if [[ ${#found[@]} -gt 0 ]]; then
        _ok "Всего $(_getPluralWord ${#found[@]} обнаружен обнаружено) ${#found[@]} $(_getPluralWord ${#found[@]} проект проекта проектов): ${_ansiSecondaryLiteral}$(_joinBy ", " "${found[@]}")"
      else
        _warn "Не обнаружено ни одного проекта"
      fi
    fi
  fi
}

# shellcheck disable=SC2034
_codeToPrepareDescriptionsOf_bw_projectInfo='
  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    eval local bw_projectInfo_bwProjShortcut_${bwProjShortcut}_description=\"Сокращение для проекта \${_ansiSecondaryLiteral}\$bwProjName \${_ansiUrl}\$bwProjGitOrigin\${_ansiReset}\"
  done
'

_prepareDuplicatePorts() {
  local -a varNames; mapfile -t varNames < <(compgen -v __usedPort)
  unset "${varNames[@]}"

  eval "$_codeToDeclareLocalBwProjVars"
  local i; for ((i=0; i<${#_bwProjDefs[@]}; i+=2)); do
    local bwProjShortcut="${_bwProjDefs[$i]}"
    codeHolder=_codeToCallPrepareBwProjVarsHelper eval "$_evalCode"
    local port; for port in "${_bwProjPorts[@]}"; do
      local Port=$(_upperFirst $port)
      local bwProjDefaultPort="bwProjDefault$Port"
      eval local "__usedPort${!bwProjDefaultPort}"
      eval "__usedPort${!bwProjDefaultPort}"'=$(( __usedPort'"${!bwProjDefaultPort}"' + 1 ))'
    done
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

    local port; for port in "${_bwProjPorts[@]}"; do
      local Port=$(_upperFirst "$port")
      local bwProjDefaultPort="bwProjDefault$Port"
      if [[ -n ${!bwProjDefaultPort} ]]; then
        echo -n "  $port: "
        if _hasItem "${!bwProjDefaultPort}" "${duplicatePorts[@]}"; then
          echo "${_ansiWarn}${!bwProjDefaultPort}${_ansiReset}"
        else
          echo "${!bwProjDefaultPort}"
        fi
      fi
    done
  else
    local -a projDirs; _prepareProjDirs "$bwProjShortcut" || return $?
    if [[ "${#projDirs[@]}" -gt 0 ]]; then
      local projDir; for projDir in "${projDirs[@]}"; do
        local tilda="~"
        local expandedProjDir=${projDir/$tilda/$HOME}
        local projDir; projDir=$(_shortenFileSpec "$projDir")
        if [[ ! -d $expandedProjDir ]]; then
          _warn "Папка ${_ansiFileSpec}$(_shortenFileSpec "$projDir")${_ansiWarn} проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} не обнаружена"
        else
          if ! _inDir "$expandedProjDir" _prepareGitDirty "$bwProjGitOrigin"; then
            _warn "Папка ${_ansiFileSpec}$(_shortenFileSpec "$projDir")${_ansiWarn} не содержит репозиторий проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn}"
          else
            local gitBranchName; gitBranchName=$(_inDir "$expandedProjDir" _gitBranch) || return $?
            _ok "Ветка ${_ansiSecondaryLiteral}$gitBranchName${_ansiOK} проекта ${_ansiPrimaryLiteral}$bwProjTitle${_ansiOK} обнаружена в ${_ansiFileSpec}$(_shortenFileSpec "$projDir")"
          fi
        fi
      done
    else
      [[ -n $skipNonExistent ]] || _warn "Проект ${_ansiPrimaryLiteral}$bwProjTitle${_ansiWarn} не обнаружен"
      return 5
    fi
  fi
}

# =============================================================================

bw_projectTestParamsOpt=( '--canBeMixedOptionsAndArgs' )
bw_projectTestParams=(
  '--githubUser/u='
  '--rootPwd/p='
  '--containerDir/d=/tmp/bw_projectTest'
  '@bwProjShortcuts:( '"$(_getBwProjShortcuts)"' )=( '"$(_getBwProjShortcuts)"' )'
)
bw_projectTest_bwProjShortcuts_name='Короткое-имя-проекта'
bw_projectTest_containerDir_description='Папка, куда будут разворачиваться проекты для тестирования'
bw_projectTest_githubUser_description='Имя пользователя на ${_ansiCmd}github${_ansiReset} для настройки ssh-доступа'
bw_projectTest_rootPwd_description='Требуется для автоматического запуска ${_ansiCmd}docker${_ansiReset} в Ubuntu'
bw_projectTest_description='Тестирование команды ${_ansiCmd}bw project${_ansiReset}'
bw_projectTestShortcuts=( 'pt' )
bw_projectTest() { eval "$_funcParams2"
  [[ $OSTYPE =~ ^darwin ]] || [[ -n $rootPwd ]] || return $(_throw "значение ${_ansiCmd}--rootPwd${_ansiReset} должно быть задано для прогона тестов в автоматическом режиме")
  bw_install --silentIfAlreadyInstalled expect || return $?
  if ! _bw_install_gitCheck; then
    [[ -n $rootPwd ]] || return $(_throw "значение ${_ansiCmd}--rootPwd${_ansiReset} должно быть задано для установки ${_ansiCmd}git${_ansiReset} в автоматическом режиме")
    ROOT_PWD="$rootPwd" expect -c '
set timeout -1
spawn bash -c ". '"$_bwFileSpec"' -p -; bw_install --silentIfAlreadyInstalled git;"
while {1} {
  expect {
    eof { break }
    -ex "sudo] password for" {
      send "$env(ROOT_PWD)\r"
    }
    '"$patch"'
  }
}'
    _bw_install_gitCheck || return $?
  fi
  if [[ -n $githubUser && ! -f ~/.ssh/id_${githubUser}@github.pub ]]; then
    [[ -n $rootPwd ]] || return $(_throw "значение ${_ansiCmd}--rootPwd${_ansiReset} должно быть задано для работы ${_ansiCmd}github-keygen${_ansiReset} в автоматическом режиме")

    if ! _bw_install_chromeCheck; then
      ROOT_PWD="$rootPwd" expect -c '
set timeout -1
spawn bash -c ". '"$_bwFileSpec"' -p -; bw_install --silentIfAlreadyInstalled chrome;"
while {1} {
  expect {
    eof { break }
    -ex "sudo] password for" {
      send "$env(ROOT_PWD)\r"
    }
    '"$patch"'
  }
}'
      _bw_install_chromeCheck || return $?
    fi

    ROOT_PWD="$rootPwd" expect -c '
    set timeout -1
spawn bash -c ". '"$_bwFileSpec"' -p -; bw github-keygen '"$githubUser"';"
while {1} {
  expect {
    eof { break }
    -ex "sudo] password for" {
      send "$env(ROOT_PWD)\r"
    }
    -ex "Enter passphrase" {
      send "$env(ROOT_PWD)\r"
    }
    -ex "Enter same passphrase again" {
      send "$env(ROOT_PWD)\r"
    }
    -ex "errors while creating key" {
      \x03
    }
    -ex "\[36m\[1mEnter\[33m\[1m" {
      sleep 1
      send "\r"
    }
    '"$patch"'
  }
}'
    if [[ ! -f ~/.ssh/id_${githubUser}@github.pub ]]; then
      return $(_throw "Failed to create key for github")
    else
      read -r -p "${_ansiWarn}Press ${_ansiPrimaryLiteral}Enter${_ansiWarn} when finished${_ansiReset}" # https://unix.stackexchange.com/questions/293940/bash-how-can-i-make-press-any-key-to-continue
    fi
  fi
  _rm -v all -d "$containerDir" || return $?
  _mkDir -v allBrief "$containerDir" || return $?
  local bwProjShortcut; for bwProjShortcut in "${bwProjShortcuts[@]}"; do
    local projDir="$containerDir/$bwProjShortcut"
    bw p $bwProjShortcut -u -p "$projDir" >/dev/null 2>&1
  done
  local bwProjShortcut; for bwProjShortcut in "${bwProjShortcuts[@]}"; do
    local projDir="$containerDir/$bwProjShortcut"
    bw p $bwProjShortcut -p "$projDir" 2>&1 | tee "$containerDir/$bwProjShortcut.p.log"
    bw prepare $bwProjShortcut "$projDir"
    ROOT_PWD="$rootPwd" "$bwProjShortcut" st -a -p "$projDir" 2>&1 | tee "$containerDir/$bwProjShortcut.st.log"
    # one more time:
    ROOT_PWD="$rootPwd" "$bwProjShortcut" st -a -p "$projDir" 2>&1 | tee -a "$containerDir/$bwProjShortcut.st.log"
  done
  echo "================================================================================"
  echo "================================ Сводка ========================================"
  local bwProjShortcut; for bwProjShortcut in "${bwProjShortcuts[@]}"; do
    echo "==== $bwProjShortcut"
    local stLogFileSpec="$containerDir/$bwProjShortcut.st.log"
    echo "${stLogFileSpec}:"
    tail -n 2 "$stLogFileSpec"
  done
  . "$_profileFileSpec"
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_installParams=(
  "${_verbosityParams[@]}"
  "--force/f"
  "--silentIfAlreadyInstalled"
)
bw_installParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_install_cmd_name=Имя-приложения
bw_install_description='Устанавливает приложения'
bw_install_force_description='Устанавливает приложения, даже если оно уже установлено'
bw_install_silentIfAlreadyInstalled_description='Для служебного пользования'
}
bw_install() { eval "$_funcParams2"
}

# shellcheck disable=SC2034
{
_codeToInstallApp='
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
_codeToRunApp='
  showResult=_showRunResult codeHolder=_codeToInstallOrRunApp eval "$_evalCode"
'
_codeToInstallOrRunApp='
  codeHolder=_codeToInitSubOPT eval "$_evalCode"
  local returnCode=0
  _osSpecific || return $?
  $showResult
  return $returnCode
'
}
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

# shellcheck disable=SC2034
{
bw_install_brewParams=()
bw_install_brew_description='Устанавливает ${_ansiPrimaryLiteral}Homebrew${_ansiReset}'
bw_install_brewCondition='[[ $OSTYPE =~ ^darwin ]]'
}
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

# shellcheck disable=SC2034
{
bw_install_gdateParams=()
bw_install_gdate_description='Устанавливает ${_ansiPrimaryLiteral}gdate${_ansiReset} (только macOS; нужен для работы профайлера bwdev)'
# bw_install_gdateCondition='[[ $OSTYPE =~ ^darwin ]]'
}
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

# shellcheck disable=SC2034
{
bw_install_gitParams=()
bw_install_git_description='Устанавливает ${_ansiPrimaryLiteral}git${_ansiReset}'
}
bw_install_git() { eval "$_funcParams2"
  name=git codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_gitCheck() {
  _which git
}
_bw_install_gitDarwin() {
  while true; do
    bw_install --silentIfAlreadyInstalled brew || { returnCode=$?; break; }
    _exec brew install "${OPT_force[@]}" git || returnCode=$?
    break
  done
}
_bw_install_gitLinux() {
  _exec "${sub_OPT[@]}" --sudo apt-get install -y --force-yes git || returnCode=$?
}

# =============================================================================

_bfg_fileName='bfg-1.13.0.jar'
bw_install_bfg_description="${_ansiUrl}https://rtyley.github.io/bfg-repo-cleaner/${_ansiReset}"
bw_install_bfg() { eval "$_funcParams2"
  name=bfg codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_bfgCheck() {
  [[ -f $HOME/bfg/"$_bfg_fileName" ]]
}
_bw_install_bfg() {
  _mkDir $HOME/bfg || return $?
  _pushd $HOME/bfg || return $?
  local returnCode=0
  while true; do
    _download -s no "http://repo1.maven.org/maven2/com/madgag/bfg/1.13.0/$_bfg_fileName" || { returnCode=$?; break; }
    local cmd="$(_quotedArgs alias bfg='java -jar "$HOME/bfg/'"$_bfg_fileName"'"')"
    _setAtBashProfile "$cmd" || { returnCode=$?; break; }
    eval "$cmd" || { returnCode=$?; break; }
    break
  done
  _popd
  return $returnCode
}


# =============================================================================

# shellcheck disable=SC2034
{
bw_install_expectParams=()
bw_install_expect_description='Устанавливает ${_ansiPrimaryLiteral}expect${_ansiReset}'
}
bw_install_expect() { eval "$_funcParams2"
  name=expect codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_expectCheck() {
  _which expect
}
_bw_install_expectDarwin() {
  while true; do
    bw_install --silentIfAlreadyInstalled brew || { returnCode=$?; break; }
    _exec brew install "${OPT_force[@]}" expect || returnCode=$?
    break
  done
}
_bw_install_expectLinux() {
  _exec "${sub_OPT[@]}" --sudo apt-get install -y expect || returnCode=$?
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_install_githubKeygenParams=()
bw_install_githubKeygen_description='Устанавливает ${_ansiPrimaryLiteral}github-keygen${_ansiReset} (${_ansiUrl}https://github.com/dolmen/github-keygen${_ansiReset})'
}
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

# shellcheck disable=SC2034
{
bw_install_xclipParams=()
# bw_install_xclipCondition='[[ $OSTYPE =~ ^darwin ]]' # вынужден был закоменнтировать, потому что bw при сборке на macOS не включает xclip в список доступных вариантов bw_install, и потом на Linux xclip недоступен
bw_install_xclip_description='Устанавливает ${_ansiPrimaryLiteral}xclip${_ansiReset} (только Linux; нужен для работы ${_ansiCmd}bw github-keygen${_ansiReset})'
}
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

# shellcheck disable=SC2034
bw_githubKeygenParams=( 'username' )
bw_githubKeygen_username_description='Имя пользователя на ${_ansiCmd}github${_ansiReset}'
bw_githubKeygen_description='Устанавливает ключ для доступа к ${_ansiCmd}github${_ansiReset}'
bw_githubKeygen() { eval "$_funcParams2"
  bw_install --silentIfAlreadyInstalled github-keygen || return $?
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

# shellcheck disable=SC2034
{
bw_install_rootCertParams=()
bw_install_rootCert_description='Устанавливает ${_ansiPrimaryLiteral}корневой сертификат dev.baza-winner.ru${_ansiReset} для localhost'
}
bw_install_rootCert() { eval "$_funcParams2"
  name="rootCert" codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_rootCertCheck() {
  if [[ $OSTYPE =~ ^linux ]]; then
    [[ -f /usr/local/share/ca-certificates/bw/root.cert.crt ]]
  elif [[ $OSTYPE =~ ^darwin ]]; then
    _silent security find-certificate -c dev.baza-winner.ru
  else
    return 1
  fi
}
_bwCertFile="$_bwDir/ssl/rootCA.pem"
_bwCertName="dev.baza-winner.ru - WinNER"
_bwDarwinCertuilFileSpec="/usr/local/opt/nss/bin/certutil"
_bw_install_rootCertDarwin() {
  while true; do
    _exec "${sub_OPT[@]}" --sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "$_bwDir/ssl/rootCA.pem" || { returnCode=$?; break; }
    _bw_install_rootCert_usingCertutil --certutil "$_bwDarwinCertuilFileSpec" || { returnCode=$?; break; }
    break
  done
}
_bw_install_rootCertLinux() {
  while true; do
    #https://thomas-leister.de/en/how-to-import-ca-root-certificate/
    [[ -d /usr/local/share/ca-certificates/bw ]] || _exec "${sub_OPT[@]}" --sudo mkdir /usr/local/share/ca-certificates/bw || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo cp "$_bwCertFile" /usr/local/share/ca-certificates/bw/root.cert.crt || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo update-ca-certificates || { returnCode=$?; break; }

    _bw_install_rootCert_usingCertutil --findFrom "$HOME" || { returnCode=$?; break; }
    break
  done
}

# shellcheck disable=SC2034
_bw_install_rootCert_usingCertutilParams=(
  '--findFrom=/'
  '--certutil=certutil'
)
_bw_install_rootCert_usingCertutil() { eval "$_funcParams2"
  local cert8DirsFileSpec=/tmp/cert8Dirs
  local cert9DirsFileSpec=/tmp/cert9Dirs
  _spinner \
    -t "Поиск файлов cert8.db cert9.db занял" \
    "Поиск файлов cert8.db cert9.db" \
    _bw_install_rootCert_prepareCertDirs
  local -a cert8Dirs
  local -a cert9Dirs
  . "$cert8DirsFileSpec"
  . "$cert9DirsFileSpec"

  if [[ ${#cert8Dirs[@]} -gt 0 || ${#cert9Dirs[@]} -gt 0 ]]; then
    bw_install --silentIfAlreadyInstalled certutil || return $?
    local -a certutil_OPT=(-A -n "${_bwCertName}" -t "TC,C,T" -i "${_bwCertFile}") # https://blogs.oracle.com/meena/about-trust-flags-of-certificates-in-nss-database-that-can-be-modified-by-certutil
    local certDir; for certDir in "${cert8Dirs[@]}"; do
      _exec "${sub_OPT[@]}" "$certutil" "${certutil_OPT[@]}" -d "dbm:${certDir}" || return $?
    done
    local certDir; for certDir in "${cert9Dirs[@]}"; do
      _exec "${sub_OPT[@]}" "$certutil" "${certutil_OPT[@]}" -d "sql:${certDir}" || return $?
    done
  fi
}
_bw_install_rootCert_prepareCertDirs() {
  _bw_install_rootCert_prepareCertDirsHelper cert8 &
  local pid8=$!
  _bw_install_rootCert_prepareCertDirsHelper cert9 &
  local pid9=$!
  wait $pid8
  wait $pid9
}
_bw_install_rootCert_prepareCertDirsHelper() {
  local certName=$1
  local certDirsFileSpecHolder="${certName}DirsFileSpec"
  local certDirs
  mapfile -t certDirs < <(find "$findFrom" -not -path '/dev/*' -name "$certName.db" -exec dirname {} \; 2>/dev/null)
  echo "${certName}Dirs=( $(_quotedArgs "${certDirs[@]}") )" > "${!certDirsFileSpecHolder}"
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_install_certutilParams=()
bw_install_certutil_description='Устанавливает ${_ansiPrimaryLiteral}certutil${_ansiReset}'
}
bw_install_certutil() { eval "$_funcParams2"
  name="certutil" codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_certutilCheck() {
  if [[ $OSTYPE =~ ^linux ]]; then
    _which certutil
  elif [[ $OSTYPE =~ ^darwin ]]; then
    [[ -f $_bwDarwinCertuilFileSpec ]]
  else
    return 1
  fi
}
_bw_install_certutilLinux() {
  while true; do
    _exec "${sub_OPT[@]}" --sudo apt install -y --force-yes libnss3-tools || { returnCode=$?; break; }
    break
  done
}
_bw_install_certutilDarwin() {
  while true; do
    bw_install --silentIfAlreadyInstalled brew || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" brew install nss || { returnCode=$?; break; }
    break
  done
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_install_chromeParams=()
bw_install_chrome_description='Устанавливает ${_ansiPrimaryLiteral}Google Chrome${_ansiReset} (пока только Linux: Ubuntu)'
}
bw_install_chrome() { eval "$_funcParams2"
  name="Google Chrome" codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_chromeCheck() {
  _which google-chrome
}
_bw_install_chromeDarwin() {
  while true; do
    # https://superuser.com/questions/602680/how-to-install-google-chrome-from-the-command-line/602692#602692
    _exec "${sub_OPT[@]}" brew install cask
    _exec "${sub_OPT[@]}" brew cask install google-chrome
    break
  done
}
_bw_install_chromeLinux() {
  while true; do
    # https://askubuntu.com/questions/642758/installing-chrome-on-ubuntu-14-04/642765#642765
    _exec "${sub_OPT[@]}" --cmdAsIs 'wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add -' || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo apt-get update || { returnCode=$?; break; }
    _exec "${sub_OPT[@]}" --sudo apt-get install -y --force-yes google-chrome-stable || { returnCode=$?; break; }
    local ubuntuVer; _exec "${sub_OPT[@]}" _getUbuntuVersion --varName ubuntuVer  || { returnCode=$?; break; }
    if [[ $ubuntuVer -le 15010 ]]; then
      # https://askubuntu.com/questions/954918/cant-launch-google-chrome-stable-getting-out-of-date-nss-libnss3/976937#976937
      # _exec "${sub_OPT[@]}" --sudo sudo dpkg -i libnspr4_4.13.1-0ubuntu0.16.10.1_amd64.deb || { returnCode=$?; break; }
      # _exec "${sub_OPT[@]}" --sudo dpkg -i libnss3_3.26.2-1ubuntu1_amd64.deb || { returnCode=$?; break; }
      # _exec "${sub_OPT[@]}" --sudo dpkg --force-all -i libnss3-nssdb_3.26.2-0ubuntu0.16.04.2_all.deb || { returnCode=$?; break; }

      # https://askubuntu.com/questions/954918/cant-launch-google-chrome-stable-getting-out-of-date-nss-libnss3/979610#979610
      _exec "${sub_OPT[@]}" --sudo apt-get update
      _exec "${sub_OPT[@]}" --sudo apt-get install -y --force-yes libnss3
    fi
    break
  done
}

# shellcheck disable=SC2034
_getUbuntuVersionParams=( '--varName=' )
_getUbuntuVersion() { eval "$_funcParams2"
  local awkFileSpec; _prepareAwkFileSpec
  local -a awk_OPT=(
    -f "$awkFileSpec"
  )
  if [[ -z $varName ]]; then
    # https://askubuntu.com/questions/686239/how-do-i-check-the-version-of-ubuntu-i-am-running/686249#686249
    lsb_release -a |  awk "${awk_OPT[@]}"
  else
    eval "$varName"'=$(lsb_release -a |  awk "${awk_OPT[@]}")'
  fi
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_install_dockerParams=()
bw_install_docker_description='Устанавливает ${_ansiPrimaryLiteral}DockerCE${_ansiReset} ${_ansiUrl}https://www.docker.com/community-edition${_ansiReset}'
}
bw_install_docker() { eval "$_funcParams2"
  name=DockerCE codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_dockerCheck() {
  _which docker && [[ $(docker --version | perl -e '$\=undef; $_=<STDIN>; printf("%d%04d", $1, $2) if m/(\d+).(\d+)/') -ge 180003 ]]
}
# shellcheck disable=SC2034
_bw_install_dockerDarwin() {
  # TODO: try `brew install docker'
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

# shellcheck disable=SC2034
{
bw_install_dockerComposeParams=()
bw_install_dockerCompose_description="Устанавливает ${_ansiPrimaryLiteral}docker-compose${_ansiReset}"
}
bw_install_dockerCompose() { eval "$_funcParams2"
  bw_install --silentIfAlreadyInstalled docker || return $?
  name=docker-compose codeHolder=_codeToInstallApp eval "$_evalCode"
}
_bw_install_dockerComposeCheck() {
  _which docker-compose
}
_bw_install_dockerComposeDarwin() {
  # try `brew install docker-compose`
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

# shellcheck disable=SC2034
{
bw_runParams=(
  "${_verbosityParams[@]}"
)
bw_runParamsOpt=(--canBeMixedOptionsAndArgs --isCommandWrapper)
bw_run_cmd_name=Имя-приложения
bw_run_description='запускает приложение'
}
bw_run() { eval "$_funcParams2"
}

# =============================================================================

# shellcheck disable=SC2034
{
bw_run_dockerParams=()
bw_run_docker_description="Запускает Docker"
}
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
    _spinner \
      "Ожидание запуска ${_ansiPrimaryLiteral}$appName${_ansiReset}" \
      _bw_run_dockerDarwinHelper \
      || { returnCode=$?; break; }
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
_redForDarwinUrl="$_bwGithubSource/master/red/red-063.darwin"
# _redForDarwinUrl="https://static.red-lang.org/dl/mac/red-063"
# _redForLinuxUrl="https://static.red-lang.org/dl/linux/red-063"
_redForDarwinFileSpec="$_redDir/red-063.darwin"
_redForLinuxUrl="$_bwGithubSource/master/red/red-063.linux"
_redForLinuxFileSpec="$_redDir/red-063.linux"
red() {
  local returnCode=0
  _osSpecific || return $?
  if [[ $returnCode -eq 99 ]]; then
    _err "Не удалось сделать файл ${_ansiCmd}$_redForDarwinFileSpec${_ansiErr} исполняемым"
  fi
  return $returnCode
}
_redDarwin() {
  while true; do
    _download -s yes "$_redForDarwinUrl" "$_redForDarwinFileSpec" || { returnCode=$?; break; }
    chmod u+x "$_redForDarwinFileSpec" || { returnCode=99; break; }
    "$_redForDarwinFileSpec" "$@" || { returnCode=$?; break; }
    break
  done
}
_redLinux() {
  while true; do
    _download -s yes "$_redForLinuxUrl" "$_redForLinuxFileSpec" || { returnCode=$?; break; }
    chmod u+x "$_redForLinuxFileSpec" || { returnCode=99; break; }
    if ! _silent dpkg -l libc6:i386 || ! _silent dpkg -l libcurl3:i386; then
      _exec --sudo dpkg --add-architecture i386
      _exec --sudo apt-get update
      _exec --sudo apt-get install -y libc6:i386 libcurl3:i386
    fi
    "$_redForLinuxFileSpec" "$@" || { returnCode=$?; break; }
    break
  done
}

# =============================================================================
