
# =============================================================================

_resetBash

# =============================================================================

verbosityDefault=none silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
bwdev_buildParams=(
  '!--projDir/p='
  '--moreDebugInfo/d'
  "${_noPregen_params[@]}"
  "${_verbosityParams[@]}"
  'scenario:(byTest only afterAllTests )=byTest'
)
bwdev_buildParamsOpt=( '--canBeMixedOptionsAndArgs' )
bwdev_build_scenario_afterAllTests_description='Сборка после полного тестирования'
bwdev_build_scenario_byTest_description='Сборка после тестирования'
bwdev_build_scenario_only_description='Сборка без тестирования'
bwdev_build_description='Тестирует и собирает bw.bash'
bwdev_build_justBuild_description='Только собирает bw.bash'
bwdev_build_moreDebugInfo_description='Больше отладочной информации'
bwdev_build() { eval "$_funcParams2"
  _prepareProjDir bwdev || return $?
  local bwBashFileSpec="$projDir/bw.bash"
  local -a OPT=( "${OPT_verbosity[@]}" "${OPT_silent[@]}" )
  local -a OPT_noPregen=()
  if [[ -n $noPregen ]]; then
    OPT_noPregen=( -p - )
  fi
  if [[ $scenario == only ]]; then
    _exec "${OPT[@]}" . "$bwBashFileSpec" "${OPT_noPregen[@]}" _buildBw
  else
    local bwTestAll=
    [[ $scenario != afterAllTests ]] || bwTestAll=true
    local dirSpec; for dirSpec in "$_bwDir/tests/$_generatedDir" "$_bwDir/$_generatedDir" ; do
      _rm "${OPT[@]}" -d "$dirSpec"
    done
    BW_VERBOSE="$moreDebugInfo" BW_TEST_ALL="$bwTestAll" _exec "${OPT[@]}" . "$bwBashFileSpec" "${OPT_noPregen[@]}" bw bt
  fi
}

bwdev_test_args_name="$bw_bashTests_args_name"
bwdev_test_args_description="$bw_bashTests_args_description"
bwdev_test_list_description="$bw_bashTests_list_description"
bwdev_testParamsOpt=( '--canBeMixedOptionsAndArgs' )
bwdev_testParams=( 
  '!--projDir/p='
  "${_noPregen_params[@]}"
  '--list' 
  '@..args' 
)
bwdev_test_description='Запускает тест(ы)'
bwdev_test() { eval "$_funcParams2"
  _prepareProjDir bwdev || return $?
  local -a OPT=() 
  if [[ -n $noPregen ]]; then
    OPT=( -p - )
  fi
  . "$projDir/bw.bash" "${OPT[@]}" && bw bt "${OPT_list[@]}" "${args[@]}"
}
bwdev_testComplete() {
  bw_bashTestsComplete
}

bwdev_profileParamsOpt=( '--canBeMoreParams' )
bwdev_profileParams=( 'cmd' )
bwdev_profile_cmd_name='Команда'
bwdev_profile_cmd_description='Профилируемая команда'
bwdev_profile_description='Профилирование кода'
bwdev_profile() { eval "$_funcParams2"
  if [[ $OSTYPE =~ ^darwin ]]; then 
    bw_install --silentIfAlreadyInstalled gdate
  fi
  _profileInit
  BW_PROFILE=true $cmd "$@"
  _profileResult
}

_initBwProjCmd

# =============================================================================
