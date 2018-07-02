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
  local fileNamesToMainArchive="\
    $(find core bash -name "*.bash") \
    $(find "$_generatedDir" -name "*$_codeBashExt" -or -name "*.completion.unset.bash" | grep -v -E "$excludeRegExp") \
    $(find docker -type f | grep -v .DS_Store) \
    https/server.crt https/server.key https/rootCA.pem \
  "
  local fileNamesToTestsArchive="\
    $(find tests -name "*.bash" -maxdepth 1) \
    $(find "tests/$_generatedDir" -name "*$_codeBashExt" -or -name "*.completion.unset.bash") \
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

# =============================================================================

