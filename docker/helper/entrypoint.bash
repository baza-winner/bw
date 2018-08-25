#!/usr/local/bin/dumb-init /bin/bash

# =============================================================================
# =============================================================================

_init() {
  if [[ ! -f /tmp/owned ]]; then
    _chown # dev -D 3 -P 8 -L 50
    
    # local cmdTitle="${_ansiCmd}sudo chown -R dev $HOME${_ansiReset}"
    # _spinner "$cmdTitle" sudo chown -R dev "$HOME"; local returnCode=$?
    # _spinner -t "Выполнение команды ${_ansiCmd}sudo chown -R dev $HOME${_ansiReset} заняло" $cmdTitle" _chown; local returnCode=$?

    if [[ $returnCode -eq 0 ]]; then
      _ok "$cmdTitle"
      touch /tmp/owned
    else
      _err "$cmdTitle"
      return $returnCode
    fi
  fi

  if [[ ! -f /tmp/sshd ]]; then
    sudo /etc/init.d/ssh start
    touch /tmp/sshd
  fi

  if [[ $- =~ i ]]; then 
    . "$HOME/.bashrc" || return $?
  else
    . "$HOME/._bashrc" || return $?
  fi

  [[ -n $_bwFileSpec ]] || . "$HOME/bw.bash" -p - || return $?
  . "$HOME/proj/bin/${_bwProjShortcut}.bash" || return $?

  local funcName="_${_bwProjShortcut}_init"
  if _funcExists "$funcName" && [[ ! -f /tmp/init ]]; then
    $funcName || return $?
    touch /tmp/init
  fi
}

# _chownParamsOpt=( --canBeMixedOptionsAndArgs )
# _chownParams=( 
#   '--homeSubdir/d='
#   '--maxProcesses/P:1..100=1'
#   '--maxFilesPerLine/L:1..1000=20'
#   '--maxdepth/D=0..'
#   '--verbose/v'
#   'user'
# )
_chown() { 
  local homeSubdir=
  # local maxProcesses=32
  local maxProcesses=8
  local maxFilesPerLine=500
  # local maxdepth=12
  local maxdepth=5
  local user=dev
  local verbose
  verbose=true
  # eval "$_funcParams2"

  local root="$HOME/$homeSubdir"
  [[ $root =~ /$ ]] && root=${root:0:-1}
  local title="${_ansiCmd}chown -R $user '$root'${_ansiReset}"
  echo "$title . . ."
  timeStart=$(date +%s)
  rm -f _chown.stdout _chown.stderr
  local batchFileSpec="/tmp/_chown.batch"
  if [[ $maxdepth -eq 0 ]]; then
    echo "sudo chown -R $user \"$root\"" > "$batchFileSpec"
    if [[ -n $verbose ]]; then
      printf "recursive: %s\n" 1
    fi
  else
    local -a maxdepth_OPT=()
    if [[ $maxdepth -gt 0 ]]; then
      maxdepth_OPT=( -maxdepth $maxdepth )
    fi
    local awkFileSpec;_prepareAwkFileSpec
    local -a awk_OPT=(
      -v "root=$root" 
      -v "maxdepth=$maxdepth" 
      -v "user=$user" 
      -v "maxProcesses=$maxProcesses"
      -v "maxFilesPerLine=$maxFilesPerLine"
      -v "verbose=$verbose"
      -f "$awkFileSpec"
    )
    { 
      sudo find "$HOME/$homeSubdir" "${maxdepth_OPT[@]}" -type d 
      echo '========'
      sudo find "$HOME/$homeSubdir" "${maxdepth_OPT[@]}" ! -type d 
    } | awk "${awk_OPT[@]}" > "$batchFileSpec"
  fi
  ( . "$batchFileSpec" )
  timeEnd=$(date +%s)
  timeElapsed=$(( timeEnd - timeStart ))
  printf "Выполнение $title заняло %ss\n" $timeElapsed
}
_prepareAwkFileSpec() { 
  # eval "$_funcParams2"
  # [[ -z $infix ]] || infix=".$infix"
  local infix=
  awkFileSpec="$(dirname "${BASH_SOURCE[1]}")/${FUNCNAME[1]}$infix.awk"
}
_entrypoint() {
  if [[ $# -gt 0 ]]; then
    local pidFileSpec="$HOME/proj/docker/$1.pid"; shift
    echo $PPID > "$pidFileSpec"
  fi

  _init || return $?

  if [[ $# -gt 0 ]]; then
    eval "$(_quotedArgs "$@")"
  elif [[ $- =~ i ]]; then
    alias q='exit 0'
    PS1="$_prompt"

    bw set horizontal-scroll-mode  # https://superuser.com/questions/848516/long-commands-typed-in-bash-overwrite-the-same-line/862341#862341

    if [[ $_hostUser == yurybikuzin ]]; then
      bw set vi
    else
      bw set vi -u
    fi

    echo "
${_ansiPrimaryLiteral}$_hostUser${_ansiOK}, Вы вошли в Docker-контейнер ${_ansiPrimaryLiteral}$_bwProjShortcut${_ansiOK} проекта ${_ansiSecondaryLiteral}$_bwProjName
${_ansiReset}В Docker-контейнере доступна команда ${_ansiCmd}$_bwProjShortcut${_ansiReset}"

  "$_bwProjShortcut" update -c
  "$_bwProjShortcut" -?

    echo "
${_ansiWarn}ВНИМАНИЕ! Для выхода из docker-контейнера выполните команду ${_ansiCmd}q"
  fi
}

# =============================================================================
# =============================================================================

_ansiReset=$'\e[0m' # https://superuser.com/questions/33914/why-doesnt-echo-support-e-escape-when-using-the-e-argument-in-macosx/33950#33950
_ansiBold=$'\e[1m'
_ansiRed=$'\e[31m'
_ansiGreen=$'\e[32m'
_ansiWhite=$'\e[97m'
_ansiCmd="${_ansiWhite}${_ansiBold}"
_ansiOK="${_ansiGreen}${_ansiBold}"
_ansiErr="${_ansiRed}${_ansiBold}"

_spinner() { eval "$_funcParams2"
  local title="$1"; shift
  ( 
    ( 
      "$@" 
      exit $?
    ) & 
    _spinnerProcessPid=$!

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

    [[ $returnCode -ne 130 ]] \
      || _err "Процесс ${_ansiHeader}$title${_ansiErr} прерван пользователем командой ${_ansiCmd}CTRL+C"

    exit $returnCode
  )
}

_bs=$'\b'
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

_err() {
  local returnCode=$?; [[ $returnCode -eq 0 ]] && returnCode=1
  echo "${_ansiErr}ERR: $*${_ansiReset}" >&2
  return $returnCode
}

_ok() {
  echo "${_ansiOK}OK: $*${_ansiReset}"
}

# =============================================================================

_entrypoint "$@"

# =============================================================================
