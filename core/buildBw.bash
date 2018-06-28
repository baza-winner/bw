# =============================================================================

_resetBash

# =============================================================================

_buildBwParams=()
_buildBw() { eval "$_funcParams2"
  [[ -n $_isBwDevelop ]] || return $(_err "Команда ${_ansiCmd}${FUNCNAME[0]}${_ansiErr} доступна только в режиме _isBwDevelop")
  _inDir "$_bwDir" _buildBwHelper
}
_buildBwHelper() {
  local fileNamesToMainArchive="\
    $(_getFileListToBundle "*.bash" core) \
    $(_getFileListToBundle "*.bash" bash) \
    $(_getFileListToBundle "*$_codeBashExt" "$_generatedDir") \
    $(_getFileListToBundle "*.completion.unset.bash" "$_generatedDir") \
    docker-compose.yml $(find docker -type f | grep -v .DS_Store) \
  "
  local fileNamesToTestsArchive="\
    $(_getFileListToBundle "*.bash" tests ) \
    $(_getFileListToBundle "*$_codeBashExt" "tests" "$_generatedDir") \
    $(_getFileListToBundle "*.completion.unset.bash" "tests" "$_generatedDir") \
  "
  local needBuild
  local -r bwOldFileName="old.$_bwFileName"
  (git show "HEAD:$_bwFileName" > "$bwOldFileName") || return $(_err "Не удалось извлечь ${_ansiFileSpec}$_bwFileName${_ansiErr} в ${_ansiFileSpec}$bwOldFileName")
  echo "${_ansiHeader}Выясняем необходимость обновления ${_ansiFileSpec}$_bwFileSpec${_ansiReset} . . ."
  diff <(echo $fileNamesToMainArchive | xargs -n1 | sort) <(_getBwTar "$bwOldFileName" | tar tf - | sort) || needBuild=true
  diff <(echo $fileNamesToTestsArchive | xargs -n1 | sort) <(_getBwTar "$bwOldFileName" tests | tar tf - | sort) || needBuild=true
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
    _addBwTar "$bwNewFileName" 'main' $fileNamesToMainArchive
    _addBwTar "$bwNewFileName" 'tests' $fileNamesToTestsArchive
    _mvFile "$bwNewFileName" "$_bwFileName" || return $?
    local bwFileSize=$(wc -c "$_bwFileName" | perl -pe "s/^\s*(\d+).*\$/\$1/")
    local bwOldFileSize=$(wc -c "$bwOldFileName" | perl -pe "s/^\s*(\d+).*\$/\$1/")
    _ok "Архив в ${_ansiFileSpec}$_bwFileName${_ansiOK} обновлен:${_ansiReset} новый размер ${_ansiFileSpec}$_bwFileName${_ansiReset} $bwFileSize $(_getPluralWord $bwFileSize байт байта байт), прежний размер $bwOldFileSize $(_getPluralWord $bwOldFileSize байт байта байт)"
  fi
  if \
    (git update-index -q --refresh && git diff-index --name-only HEAD -- | grep "$_bwFileName" >/dev/null 2>&1) && \
    ! ( git diff "$_bwFileName" | grep -E '^\+export\s+_bwVersion=' >/dev/null 2>&1)
  then
    _warn "${_ansiFileSpec}$_bwFileName${_ansiWarn} изменен. Необходимо изменить номер версии ${_ansiOutline}_bwVersion${_ansiWarn}"
  fi
}
_getFileListToBundleParamsOpt=(--canBeMoreParams)
_getFileListToBundleParams=( '--deep/d' 'fileMask' 'dirName' )
_getFileListToBundle() { eval "$_funcParams2"
  local dirToFind="$dirName"
  local dirForRegexp="$dirName"
  while [[ $# -gt 0 ]]; do
    [[ -n $1 ]] || continue
    dirToFind+="/$1"
    dirForRegexp+="\/$1"
    shift
  done
  local regexp
  if [[ -n $deep ]]; then
    regexp="$dirForRegexp\/.+"
  else
    regexp="$dirForRegexp\/[^\/]+"
  fi
  local excludeRegExp=
  local bwProjShortcut; for bwProjShortcut in $(_getBwProjShortcuts); do
    if [[ -n $excludeRegExp ]]; then
      excludeRegExp+='|'
    fi
    excludeRegExp+="$bwProjShortcut"
  done
  excludeRegExp="generated/($excludeRegExp)"
  find "$_bwDir/$dirToFind" -type f -name "$fileMask" | grep -v -E "$excludeRegExp" | perl -ne "print $_ if s/^.*\/($regexp)\$/\$1/"
}

# =============================================================================

