#!/bin/bash

# =============================================================================

if [[ ! -f /tmp/owned ]] && sudo chown -R dev "$HOME"; then
	touch /tmp/owned
fi

. ~/.bashrc

. ~/bw.bash -p -

. "$HOME/proj/bin/$_bwProjShortcut.bash"

# =============================================================================

if [[ -n $1 ]]; then
  $_bwProjShortcut "$@"
else
  alias q='exit 0'
  PS1="$_prompt"

  if [[ $_hostUser == yurybikuzin ]]; then
    bw set vi
  else
    bw set vi -u
  fi

  echo "
${_ansiPrimaryLiteral}$_hostUser${_ansiOK}, Вы вошли в Docker-контейнер ${_ansiPrimaryLiteral}$_bwProjShortcut${_ansiOK} проекта ${_ansiSecondaryLiteral}$_bwProjName
${_ansiYellow}Подсказка: Чтобы выйти из Docker-container'а, выполните команду ${_ansiCmd}q
${_ansiReset}В Docker-контейнере доступна команда ${_ansiCmd}$_bwProjShortcut${_ansiReset}"

  $_bwProjShortcut update -m completionOnly
  $_bwProjShortcut -?
fi

# =============================================================================
