
# =============================================================================

_resetBash

# =============================================================================

verbosityDefault=none silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
export bwdev_buildParams=(
  '--moreDebugInfo/d'
  "${_verbosityParams[@]}"
  'scenario:(byTest only afterAllTests )=byTest'
)
export bwdev_build_scenario_afterAllTests_description='Сборка после полного тестирования'
export bwdev_build_scenario_byTest_description='Сборка после тестирования'
export bwdev_build_scenario_only_description='Сборка без тестирования'
export bwdev_build_description='Тестирует и собирает bw.bash'
export bwdev_build_justBuild_description='Только собирает bw.bash'
export bwdev_build_moreDebugInfo_description='Больше отладочной информации'
bwdev_build() { eval "$_funcParams2"
  local BW_VERBOSE="$moreDebugInfo"
  local bwBashFileSpec="$_bwdevDir/bw.bash"
  local -a OPT=( "${OPT_verbosity[@]}" "${OPT_silent[@]}" )
  if [[ $scenario == only ]]; then
    _exec "${OPT[@]}" --cmdAsIs '. "'"$bwBashFileSpec"'" _buildBw'
  else
    if [[ $scenario == afterAllTests ]]; then
      local BW_TEST_ALL=true
    fi
    _exec "${OPT[@]}" --cmdAsIs '. "'"$bwBashFileSpec"'" -p - bw rm -y && . "'"$bwBashFileSpec"'" bw bt'
  fi
}

export bwdev_test_args_name="$bw_bashTests_args_name"
export bwdev_test_args_description="$bw_bashTests_args_description"
export bwdev_test_list_description="$bw_bashTests_list_description"
export bwdev_test_noPregen_description="Исключить прегенерация"
export bwdev_testParams=( 
  '--noPregen/n'
  '--list' 
  '@..args' 
)
export bwdev_testParamsOpt=( --canBeMixedOptionsAndArgs )
export bwdev_test_description='Запускает тест(ы)'
bwdev_test() { eval "$_funcParams2"
  local -a OPT=() 
  if [[ -n $noPregen ]]; then
    OPT=( -p - )
  fi
  . "$_bwdevDir/bw.bash" "${OPT[@]}" && bw bt "${OPT_list[@]}" "${args[@]}"
}
bwdev_testComplete() {
  bw_bashTestsComplete
}

bwdev_profileParamsOpt=( --canBeMoreParams )
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
