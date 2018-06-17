
# =============================================================================

_resetBash

# =============================================================================

_bs=$'\b'

_spinnerShowPid=
_spinnerProcessPid=
_spinnerParamsOpt=( --canBeMoreParams )
_spinnerParams=( 
  '--timing/t=' 
  '--timingTreshold=0' 
  'title' 
  'cmd' 
)
_spinner() { eval "$_funcParams2"
  [[ -z $_spinnerShowPid ]] || return $(_err "Spinner уже запущен")
  local _didNotStartSpinner
  [[ -z $timing ]] || local timeStart=$(date +%s)

  _profileInitTransfer
  ( profileTmpFileSpec="$_profileTmpFileSpec"
    _profileInitTransfer
    ( profileTmpFileSpec="$_profileTmpFileSpec"
      $cmd "$@"; returnCode=$?
      _profileDoTransfer
      exit $returnCode
    ) & 
    _spinnerProcessPid=$!
    _profileGetTransfer

    _showSpinner "$title" & 
    _spinnerShowPid=$!

    trap 'trap - SIGINT; trap - EXIT; kill -SIGTERM $_spinnerProcessPid' SIGINT
    trap 'trap - EXIT; kill -SIGTERM $_spinnerProcessPid; kill -SIGTERM $_spinnerShowPid' EXIT # https://unix.stackexchange.com/questions/17314/what-is-signal-0-in-a-trap-command/17315#17315

    wait $_spinnerProcessPid 2>/dev/null; local returnCode=$?
    unset _spinnerProcessPid

    kill -SIGTERM $_spinnerShowPid
    wait $_spinnerShowPid 2>/dev/null # https://stackoverflow.com/questions/5719030/bash-silently-kill-background-function-process/5722850#5722850
    unset _spinnerShowPid

    trap - SIGINT
    trap - EXIT

    if [[ $returnCode -eq 130 ]]; then
      _err "Процесс ${_ansiHeader}$title${_ansiErr} прерван пользователем командой ${_ansiCmd}CTRL+C"
    elif [[ -n $timing && $returnCode -eq 0 ]]; then
      local timeEnd=$(date +%s)
      local timeElapsed=$(( timeEnd - timeStart ))
      if [[ $timeElapsed -ge $timingTreshold ]]; then
        echo "$timing $timeElapsed $(_getPluralWord $timeElapsed секунду секунды секунд)"
      fi
    fi

    exit $returnCode
  ); local returnCode=$?
  _profileGetTransfer
  return $returnCode
}

_spinnerSleepTime=0.1
_showSpinner() { # https://unix.stackexchange.com/questions/11498/how-to-trap-a-suspend-a-resume-from-a-bash-script
  local title; [[ $1 == '-' ]] || title="$1 . . . "
  local waitChars='|/-\|/-\'
  eval _spinnerTitle=\${title:--}
  eval _spinnerTitleLen=\${#title}
  eval _spinnerWill=start
  trap '_spinnerWill=pause' SIGTSTP
  trap '_spinnerWill=continue' SIGCONT
  trap '[[ $_spinnerWill == pause || $_spinnerWill == start ]] && _spinnerWill=exit || _spinnerWill=term' SIGTERM
  while true; do
    local idx; for ((idx=0; idx<${#waitChars}; idx++)); do
      sleep $_spinnerSleepTime # https://serverfault.com/questions/469247/how-do-i-sleep-for-a-millisecond-in-bash-or-ksh/469249#469249
      if [[ $_spinnerWill == 'start' ]]; then
        echo -n "$title "
        _spinnerWill=
      elif [[ $_spinnerWill == 'term' ]]; then
        local i; for (( i=0; i<=_spinnerTitleLen; i++)); do echo -n $_bs; done
        local i; for (( i=0; i<=_spinnerTitleLen; i++)); do echo -n ' '; done
        local i; for (( i=0; i<=_spinnerTitleLen; i++)); do echo -n $_bs; done
        exit
      elif [[ $_spinnerWill == 'exit' ]]; then
        exit
      elif [[ $_spinnerWill == 'undo' ]]; then
        echo -n ' '
        _spinnerWill=
      elif [[ $_spinnerWill == 'pause' ]]; then
        echo -n "$_bs $_bs"
        _spinnerWill=beContinued
      elif [[ -z $_spinnerWill ]]; then
        echo -n $_bs${waitChars:$idx:1}
      fi
    done
  done
}

_runInBackgroundParamsOpt=( --canBeMoreParams )
_runInBackgroundParams=( 'cmd' )
_runInBackground() { eval "$_funcParams2"
  _profileInitTransfer
  ( profileTmpFileSpec="$_profileTmpFileSpec"
    _profileInitTransfer
    ( profileTmpFileSpec="$_profileTmpFileSpec"
      $cmd "$@"; returnCode=$?
      _profileDoTransfer
      exit $returnCode
    ) & 
    local pid=$!
    _profileGetTransfer

    trap 'trap - SIGINT; kill -SIGTERM $pid' SIGINT
    trap 'trap - EXIT; kill -SIGTERM $pid' EXIT

    wait $pid 2>/dev/null; local returnCode=$?
    
    trap - SIGINT
    trap - EXIT

    _profileDoTransfer
    exit $returnCode
  ); local returnCode=$?
  _profileGetTransfer
  return $returnCode
}

_pauseSpinnerBeforeSudo() {
  [[ -n $_spinnerShowPid ]] || return
  kill -SIGTSTP $_spinnerShowPid
  sleep $_spinnerSleepTime
}

_resumeSpinnerAfterSudo() {
  [[ -n $_spinnerShowPid ]] || return
  kill -SIGCONT $_spinnerShowPid
}

# =============================================================================
