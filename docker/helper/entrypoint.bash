#!/usr/local/bin/dumb-init /bin/bash

# =============================================================================
# =============================================================================

_init() {
  if [[ ! -f /tmp/owned ]]; then
    sudo chown dev "$HOME/bw.bash"
    sudo chown -R dev "$HOME/.bw"
  fi

  if [[ $- =~ i ]]; then
    . "$HOME/.bashrc" || return $?
  else
    . "$HOME/bw.bash" -p - || return $?
    . "$HOME/._bashrc" || return $?
  fi

  . "$HOME/proj/bin/${_bwProjShortcut}.bash" || return $?

  if [[ ! -f /tmp/owned ]]; then

    dstVarName=notPath srcVarName=_${_bwProjShortcut}_chown_notPath codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
    local -a OPT_notPath=()
    local notPathItem; for notPathItem in "${notPath[@]}"; do
      OPT_notPath+=( --notPath "$HOME/$notPathItem" )
    done

    local homeSubdir=proj
    local title="${_ansiCmd}chown -R dev $HOME/$homeSubdir${_ansiReset}"
    _spinner \
      -t "Выполнение $title заняло" \
      "$title" \
      _chown dev "$homeSubdir" -P 8 "${OPT_notPath[@]}" # -v

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


  local funcName="_${_bwProjShortcut}_init"
  if _funcExists "$funcName" && [[ ! -f /tmp/init ]]; then
    $funcName || return $?
    touch /tmp/init
  fi
}

# =============================================================================

_entrypoint() {
  if [[ $# -gt 0 ]]; then
    if [[ ! -f /tmp/owned ]]; then
      sudo chown -R dev "$HOME/proj/docker"
    fi
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

_entrypoint "$@"

# =============================================================================
# =============================================================================
