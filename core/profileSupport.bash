#!/bin/bash

{
  # shellcheck disable=SC2154
  if [[ -n $shellcheck ]]; then 
    . bw.bash
  fi
}

# =============================================================================

_resetBash

# =============================================================================

_gdate='date'; [[ $OSTYPE =~ ^darwin ]] && _gdate='gdate'

_profileStack=()
_profileStackIdx=0
  # shellcheck disable=SC2120
_profileBegin() {
  if [[ -n $BW_PROFILE ]]; then
    local needProfile=
    local profileName=${FUNCNAME[1]}
    [[ -z $1 ]] || profileName=$1_IN_$profileName
    if [[ $_profileStackIdx -eq 0 || -n ${_profileStack[$(( _profileStackIdx - 1))]} ]]; then
      local profileHolder=_profile_$profileName
      [[ ${!profileHolder} == false ]] || needProfile=0
    fi
    _profileStack[$_profileStackIdx]=$needProfile
    _profileStackIdx=$((_profileStackIdx + 1))
    [[ -z $needProfile ]] || eval "_startTime_$profileName"'=$($_gdate +%s%3N)'
  fi
}

_gdateTiming=0

  # shellcheck disable=SC2120
_profileEnd() {
  if [[ -n $BW_PROFILE ]]; then
    _profileStackIdx=$((_profileStackIdx - 1))
    if [[ $_profileStackIdx -lt 0 ]]; then
      _debugVar _profileStackIdx
      _profileStackIdx=0
    elif [[ -n ${_profileStack[$_profileStackIdx]} ]]; then
      local profileName=${FUNCNAME[1]}
      [[ -z $1 ]] || profileName=$1_IN_$profileName
      local startTimeHolder=_startTime_$profileName
      local endTime; endTime=$($_gdate +%s%3N)
      local elapsedTime=$(( endTime - ${!startTimeHolder} ))
      if [[ $_profileStackIdx -gt 0 ]]; then
        local parentProfileStackIdx=$(( _profileStackIdx - 1 ))
        _profileStack[$parentProfileStackIdx]=$(( _profileStack[parentProfileStackIdx] + elapsedTime + _gdateTiming ))
      fi
      elapsedTime=$(( elapsedTime - _profileStack[_profileStackIdx] - _gdateTiming ))
      [[ $elapsedTime -ge 0 ]] || elapsedTime=0
      eval "_profileTotal_$profileName"'=$(( _profileTotal_'"$profileName"' + elapsedTime ))'
      eval "_profileCount_$profileName"'=$(( _profileCount_'"$profileName"' + 1 ))'
    fi
  fi
}

_profileInit() {
  local prefix=_profileTotal_
  local varName; for varName in $(compgen -v | grep "^$prefix" ); do # https://unix.stackexchange.com/questions/3510/how-to-print-only-defined-variables-shell-and-or-environment-variables-in-bash/5691#5691
    [[ -n $varName ]] || continue
    eval "$varName=0"
    local funcName=${varName:${#prefix}}
    eval "_profileCount_$funcName=0"
  done

  _gdateTiming=$($_gdate +%s%3N)
  local _i; for ((_i=0; _i<30; _i++)); do
    $_gdate +%s%3N >/dev/null
  done
  # shellcheck disable=SC2017
  _gdateTiming=$(( ( $($_gdate +%s%3N) - _gdateTiming ) / 30 * 6 / 5 ))
}

_profileResult() {
  local prefix=_profileTotal_
  local total=0
  local varName; for varName in $(compgen -v | grep "^$prefix" ); do # https://unix.stackexchange.com/questions/3510/how-to-print-only-defined-variables-shell-and-or-environment-variables-in-bash/5691#5691
    local funcName=${varName:${#prefix}}
    local profileHolder=_profile_$funcName
    if [[ ${!profileHolder} != false ]]; then
      local countHolder=_profileCount_$funcName
      if [[ ${!countHolder} -gt 0 ]]; then
        echo "$funcName: ${!varName}ms ${!countHolder} $(( ${!varName} / ${!countHolder} ))ms"
        total=$(( total + ${!varName} ))
      fi
    fi
  done
  echo total ${total}ms
}

_profileTmpFileSpec=
_profileInitTransfer() {
  if [[ -n $BW_PROFILE ]]; then 
    _profileTmpFileSpec="/tmp/bw.profile.$$.bash"
  fi
}

_profileDoTransfer() {
  if [[ -n $BW_PROFILE ]]; then 
    # shellcheck disable=SC2154
    echo >"$profileTmpFileSpec"
    local prefix=_profileTotal_
    local varName; for varName in $(compgen -v | grep "^$prefix" ); do # https://unix.stackexchange.com/questions/3510/how-to-print-only-defined-variables-shell-and-or-environment-variables-in-bash/5691#5691
      echo "$varName=${!varName}" >>"$profileTmpFileSpec"
      local funcName=${varName:${#prefix}}
      varName=_profileCount_$funcName
      echo "$varName=${!varName}" >>"$profileTmpFileSpec"
    done
  fi
}

_profileGetTransfer() {
  if [[ -n $BW_PROFILE ]]; then 
    # shellcheck disable=SC1090
    . "$_profileTmpFileSpec"
  fi
}

_noSleepWhileProfile=true
_funcToProfile() { 
  # shellcheck disable=SC2119
  _profileBegin
  [[ -n $_noSleepWhileProfile ]] || sleep 1
  _subFuncToProfile
  _subFuncToProfile2
  # shellcheck disable=SC2119
  _profileEnd 
}

_subFuncToProfile() { 
  # shellcheck disable=SC2119
  _profileBegin
  [[ -n $_noSleepWhileProfile ]] || sleep .3
  _subSubFuncToProfile
  # shellcheck disable=SC2119
  _profileEnd 
}

_subSubFuncToProfile() { 
  # shellcheck disable=SC2119
  _profileBegin
  [[ -n $_noSleepWhileProfile ]] || sleep .1
  _subSubSubFuncToProfile
  # shellcheck disable=SC2119
  _profileEnd 
}

_subSubSubFuncToProfile() { 
  # shellcheck disable=SC2119
  _profileBegin
  [[ -n $_noSleepWhileProfile ]] || sleep .1
  # shellcheck disable=SC2119
  _profileEnd 
}

_subFuncToProfile2() { 
  # shellcheck disable=SC2119
  _profileBegin
  [[ -n $_noSleepWhileProfile ]] || sleep .5
  # shellcheck disable=SC2119
	_profileEnd 
}

# =============================================================================

