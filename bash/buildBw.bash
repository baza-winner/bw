# =============================================================================

_resetBash

# =============================================================================

_buildBwParams=()
_buildBw() { eval "$_funcParams2"
  [[ -n $_isBwDevelop ]] || return $(_err "Команда ${_ansiCmd}${FUNCNAME[0]}${_ansiErr} доступна только в режиме _isBwDevelop")
  _inDir "$_bwDir" _buildBwHelper
}
_buildBwHelper() {
  local excludeRegExp='_bwCreateServerHttpsCrtKey'
  local bwProjShortcut; for bwProjShortcut in $(_getBwProjShortcuts); do
    if [[ -n $excludeRegExp ]]; then
      excludeRegExp+='|'
    fi
    excludeRegExp+="$bwProjShortcut"
  done
  excludeRegExp="generated/($excludeRegExp)"
  local -a OPT_docker=(
    -type f 
    ! -name '*.log'
    ! -name .DS_Store 
    ! -name .gitignore 
    ! -name docker-compose.yml 
    ! -name '*.env'
    ! -path '*nginx/whoami*' 
    ! -path '*nginx/conf.d*' 
    ! -path '*nginx/main.conf' 
    ! -path '*nginx/nginx.conf'
  )
  # $(find docker -type f ! -name '*.log' ! -name .DS_Store ! -name .gitignore ! -name 'docker-compose*' ! -name '*.env' ! -path '*nginx/whoami*' ! -path '*nginx/main.conf' ! -path '*nginx/nginx.conf') \
  # $(find docker -type f -name docker-compose.*.yaml) \
  local -a fileNamesToMainArchive=(
    $(find core bash -name '*.bash' -or -name '*.c' -or -name '*.awk')
    $(find "$_generatedDir" -name "*$_codeBashExt" -or -name "*.completion.unset.bash" | grep -v -E "$excludeRegExp")

    docker/nginx/conf.d/http.conf
    docker/nginx/conf.d/https.conf

    docker/helper/docker-compose.nginx.yml
    docker/helper/docker-compose.main.yml 
    docker/helper/mysql_secure_installation.sql

    docker/entrypoint.bash
    docker/docker-compose.proj.yml

    docker/nginx/conf.bw/whoami.conf
    docker/nginx/conf.bw/http.conf
    docker/nginx/conf.bw/nginx.conf
    docker/nginx/conf.bw/https.conf

    ssl/server.crt ssl/server.key ssl/rootCA.pem
    ssh/authorized_keys ssh/bw_dev_id_rsa ssh/bw_dev_id_rsa.pub
    git-completion/git-completion.bash
    git-flow-completion/git-flow-completion.bash
  )
  local -a fileNamesToTestsArchive=(
    $(find tests -maxdepth 1 -name "*.bash")
    $(find "tests/$_generatedDir" -name "*$_codeBashExt" -or -name "*.completion.unset.bash")
  )

  local needBuild
  local -r bwOldFileName="old.$_bwFileName"
  (git show "HEAD:$_bwFileName" > "$bwOldFileName") || return $(_err "Не удалось извлечь ${_ansiFileSpec}$_bwFileName${_ansiErr} в ${_ansiFileSpec}$bwOldFileName")
  echo "${_ansiHeader}Выясняем необходимость обновления ${_ansiFileSpec}$_bwFileSpec${_ansiReset} . . ."
  diff <(echo "${fileNamesToMainArchive[@]}" | xargs -n1 | sort) <(_getBwTar "$bwOldFileName" | tar tf - | sort) || needBuild=true
  diff <(echo "${fileNamesToTestsArchive[@]}" | xargs -n1 | sort) <(_getBwTar "$bwOldFileName" tests | tar tf - | sort) || needBuild=true
  local -r tgzDir=tgz
  if
    _mkDir -t "$tgzDir" && \
    _getBwTar "$bwOldFileName" | tar xf - -C "$tgzDir" && \
    _getBwTar "$bwOldFileName" tests | tar xf - -C "$tgzDir" \
  ; then
    local fileName; for fileName in $(_getBwTar "$bwOldFileName" | tar tf -) $(_getBwTar "$bwOldFileName" tests | tar tf -); do
      [[ -f $fileName ]] || continue
      cmp "$fileName" "$tgzDir/$fileName" || needBuild=true
    done
  fi

  if [[ -z $needBuild ]]; then
    _ok "Архив в ${_ansiFileSpec}$_bwFileName${_ansiOK} не нуждается в обновлении"
  else
    local -r bwNewFileName="new.$_bwFileName"
    head -n $(grep -n '^# ==BINARY[[:space:]]*$' "$_bwFileName" | cut -d ':' -f 1) "$_bwFileName" > "$bwNewFileName"
    _addBwTar "$bwNewFileName" 'main' "${fileNamesToMainArchive[@]}"
    _addBwTar "$bwNewFileName" 'tests' "${fileNamesToTestsArchive[@]}"
    _mvFile "$bwNewFileName" "$_bwFileName" || return $?
    local bwFileSize=$(wc -c "$_bwFileName" | perl -pe "s/^\s*(\d+).*\$/\$1/")
    local bwOldFileSize=$(wc -c "$bwOldFileName" | perl -pe "s/^\s*(\d+).*\$/\$1/")
    _ok "Архив в ${_ansiFileSpec}$_bwFileName${_ansiOK} обновлен:${_ansiReset} новый размер ${_ansiFileSpec}$_bwFileName${_ansiReset} $bwFileSize $(_getPluralWord $bwFileSize байт байта байт), прежний размер $bwOldFileSize $(_getPluralWord $bwOldFileSize байт байта байт)"
  fi
  if \
    (git update-index -q --refresh && git diff-index --name-only HEAD -- | grep "$_bwFileName" >/dev/null 2>&1) && \
    ! ( git diff "$_bwFileName" | grep -E '^\+_(export\s+)?bwVersion=' >/dev/null 2>&1)
  then
    _warn "${_ansiFileSpec}$_bwFileName${_ansiWarn} изменен. Необходимо изменить номер версии ${_ansiOutline}_bwVersion${_ansiWarn}"
  fi
}

# =============================================================================

_addBwTar() {
  local fileSpec="$1"; shift
  local archiveName="$1"; shift
  { 
    echo "# ==$archiveName start"
    COPYFILE_DISABLE=1 tar cf - "$@" | gzip | base64 --break=80
    echo "# ==$archiveName end" 
  } >> "$fileSpec"
}

# =============================================================================

