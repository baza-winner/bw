
# =============================================================================

_resetBash

# =============================================================================

_verbosityParams=()
_codeToPrepareVerbosityParams='
  _verbosityParams=(
    "!--verbosity/v:(dry none err ok allBrief all)=${verbosityDefault:-err}"
    "!--silent/s:(yes ok err no)=${silentDefault:-no}"
  )
'

_verbosity_dry_description='режим "dry run"'
_verbosity_none_description='молчаливый режим'
_verbosity_err_description='только вывод ошибок'
_verbosity_ok_description='только вывод OK'
_verbosity_allBrief_description='вывод результатов в краткой форме'
_verbosity_all_description='полный вывод'

_silent_yes_description='не подавлять вывод вспомогательных команд'
_silent_ok_description='подавлять вывод вспомогательных команд в случае отсутствия ошибок'
_silent_err_description='подавлять вывод вспомогательных команд в случае наличия ошибок'
_silent_no_description='подавлять вывод вспомогательных команд'

# =============================================================================

_silent() {
  "$@" >/dev/null 2>&1
}

verbosityDefault=all silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_execParams=(
  "${_verbosityParams[@]}"
  '--stdin='
  '--stdout='
  '--stderr='
  '--append='
  '@--treatAsOK:0..'
  '--preserveReturnCode'
  '--untilSuccessSleep:0..'
  '--cmdAsIs'
  '--sudo'
  'cmd'
)
_execParamsOpt=(
  --canBeMoreParams
  --treatUnknownOptionAsArg
)
_exec() { eval "$_funcParams2"
  _profileBegin
  if [[ $silent == no ]]; then
    [[ -z $append || ( -z $stdout && -z $stderr ) ]] \
      || return $(_throw "не ожидает одновременно ${_ansiCmd}${OPT_append[@]} ${OPT_stdout[@]} ${OPT_stderr[@]}")
  else
    [[ -z $append && -z $stdout && -z $stderr ]] \
      || return $(_throw "не ожидает одновременно ${_ansiCmd}${OPT_silent[@]} ${OPT_append[@]} ${OPT_stdout[@]} ${OPT_stderr[@]}")
  fi

  if [[ $silent != 'no' ]]; then
    if [[ $silent == 'yes' ]]; then
      stdout='/dev/null'
    else
      local outputFileSpec="/tmp/bw.$$.output"
      stdout="$outputFileSpec"
    fi
    stderr='&1'
  fi

  [[ -n $cmdAsIs ]] && cmd="$cmd $@" || cmd=$(_quotedArgs "$cmd" "$@")

  [[ -z $stdin ]] || stdin="<$(_quotedArgs "$stdin")"
  [[ -z $stdout ]] || stdout=">$(_quotedArgs "$stdout")"
  [[ -z $stderr ]] || stderr="2>$(_quotedArgs "$stderr")"
  [[ -z $append ]] || append=">>$(_quotedArgs "$append")"

  local -a redir=( $stdin $stdout $stderr $append )
  [[ ${#redir[@]} -eq 0 ]] || cmd+=" ${redir[@]}"

  local returnCode=0
  if [[ $verbosity == 'dry' ]]; then
    echo "${_ansiCmd}$cmd${_ansiReset}"
  else
    [[ $verbosity == 'all' ]] && echo "${_ansiCmd}$cmd${_ansiReset} . . ."
    if [[ -z $sudo ]]; then
      eval "$cmd"
    else
      sudo -n true 2>/dev/null; local ret=$?
      local noSudo; [[ $ret -ne 0 ]] && noSudo=true
      [[ -n $noSudo ]] && _pauseSpinnerBeforeSudo
      eval "sudo $cmd"
    fi
    returnCode=$?
    [[ -n $sudo && -n $noSudo ]] && _resumeSpinnerAfterSudo
    # _debugVar returnCode
    # [[ $returnCode -eq $treatAsOK ]] && returnCode=0
    local isOK=
    if [[ $returnCode -eq 0 ]]; then
      isOK=true
    elif [[ -z $isOK && ${#treatAsOK[@]} -gt 0 ]] && _hasItem $returnCode ${treatAsOK[@]}; then
      isOK=true
      [[ -n $preserveReturnCode ]] || returnCode=0
    fi
    if [[ -n $untilSuccessSleep ]]; then
      while [[ $returnCode -ne 0 ]]; do
        sleep $untilSuccessSleep
        eval "$cmd"; returnCode=$?
        [[ $returnCode -eq $treatAsOK ]] && returnCode=0
      done
    fi
    if [[ ( $silent == 'ok' && -z $isOK ) || ( $silent == 'err' && -n $isOK ) ]]; then
      [[ -f $outputFileSpec ]] && cat $outputFileSpec
    fi
    if [[ $verbosity != 'none' ]]; then
      if [[ -n $isOK ]]; then
        [[ $verbosity == 'err' ]] || _ok "${_ansiCmd}$cmd"
      else
        [[ $verbosity == 'ok' ]] || _err "${_ansiCmd}$cmd"
      fi
    fi
  fi
  _profileEnd
  return $returnCode
}

_downloadEtagFileExt=".header"
_downloadTempFileExt=".download"

_downloadDescription='
  Загружает содержимое ${_ansiOutline}url${_ansiReset}'\''а в файл ${_ansiOutline}fileSpec${_ansiReset}
  Если загрузка была прервана, то продолжает загрузку
'
_downloadDescriptionOfUrl='URL для загрузки'
_downloadDescriptionOfFileSpec='Путь к файлу, куда поместить результат загрузки'
_downloadDescriptionOfSilent='Молчаливый режим'
_downloadDescriptionOfCheck_force='
  Загружает файл с самого начала независимо ни от чего
  Игнорирует опцию ${_ansiCmd}--tolerant${_ansiReset}
'
_downloadDescriptionOfCheck_etag='
  Загружает файл, если файл ранее не был загружен,
  или если на сервере изменялся ETag файла с момента последней загрузки
'
_downloadDescriptionOfCheck_time='
  Загружает файл, если файл ранее был не загружен,
  или если на сервере он изменился с момента последней загрузки
'
_downloadDescriptionOfTolerant='
   Игнорирует ошибки, связанные с качеством интернета (его может просто не быть),
   если задана опция ${_ansiCmd}--check${_ansiReset} и файл ${_ansiOutline}fileSpec${_ansiReset} уже существует (был ранее скачан)
'
_downloadDescriptionOfReturnCodeIfActuallyUpdated='
   Устанавливает код возврата на случай,
   если файл ${_ansiOutline}fileSpec${_ansiReset} был действительно обновлен
'
verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_downloadParams=(
  "${_verbosityParams[@]}"
  '--check/c:(etag time force)'
  '--tolerant/t'
  '--returnCodeIfActuallyUpdated/r:0..=0'
  'url'
  'fileSpec=$(basename "$url")'
)
_downloadParamsOpt=(--canBeMixedOptionsAndArgs)
_download() { eval "$_funcParams2"
  local downloadTempFileSpec="$fileSpec$_downloadTempFileExt"
  [[ -f $fileSpec && ! -f $downloadTempFileSpec && -z $check ]] && return 0

  local dirName=$(dirname "$fileSpec")
  [[ -d "$dirName" ]] || mkdir "$dirName" || return $(_err "Не удалось создать директорию ${_ansiCmd}$dirName")

  local -a xOpt=()
  [[ -n $silent ]] && xOpt+=( -s )
  if [[ -f $downloadTempFileSpec && $check != 'force' ]]; then
    xOpt+=( -C - )
  elif [[ $check == 'time' ]]; then
    [[ -f $fileSpec ]] && xOpt+=( -z "$fileSpec" )
  fi
  local etagFileSpec="$fileSpec$_downloadEtagFileExt"
  local downloadEtagFileSpec="$downloadTempFileSpec$_downloadEtagFileExt"
  xOpt+=( --dump-header "$downloadEtagFileSpec" )
  if [[ $check == 'etag' && ! -f $downloadTempFileSpec && -f $etagFileSpec ]]; then
    local etag=$(grep -i "^etag: " "$etagFileSpec" | perl -pe "s/^etag: //i; s/\s*$//")
    [[ -n $etag ]] && xOpt+=( --header "If-None-Match: $etag" )
  fi
  _rm ${OPT_verbosity[@]} "$downloadEtagFileSpec"
  _exec ${OPT_verbosity[@]} curl -o "$downloadTempFileSpec" -L "$url" "${xOpt[@]}"; local returnCode=$?
  if [[ $returnCode -ge 5 && $returnCode -le 7 && -f $fileSpec && $tolerant && $check != 'force' ]]; then
    _rm ${OPT_verbosity[@]} "$downloadEtagFileSpec"
    return 0
  elif [[ ( $returnCode -eq 33 || $returnCode -eq 22 ) && $check != force ]]; then
    eval _download ${OPT_verbosity[@]} ${OPT_silent[@]} --check force "$url" "$fileSpec"
  elif [[ $returnCode -ne 0 ]]; then
    return $(_err "Не удалось скачать ${_ansiUrl}$url${_ansiReset}${_ansiErr} в ${_ansiCmd}$fileSpec${_ansiErr}. Код возврата ${_ansiCmd}curl${_ansiErr}: ${_ansiPrimaryLiteral}$returnCode")
  elif [[ -f $downloadEtagFileSpec ]] && _silent grep -i "^HTTP\S*\s304\b" $downloadEtagFileSpec; then
    _rm ${OPT_verbosity[@]} "$downloadTempFileSpec" "$downloadEtagFileSpec" || return $?
  else
    _mvFile ${OPT_verbosity[@]} "$downloadTempFileSpec" "$fileSpec" || return $?
    _mvFile ${OPT_verbosity[@]} "$downloadEtagFileSpec" "$etagFileSpec" || return $?
    return $returnCodeIfActuallyUpdated
  fi
}

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_inDirParams=(
  "${_verbosityParams[@]}"
  '--tmp/t'
  '--noCleanOnFail/n'
  '@--treatAsOK:0..'
  '--preserveReturnCode'
  'dirSpec'
  'inDirCommnand'
)
_inDirParamsOpt=(--canBeMoreParams)
_inDir() { eval "$_funcParams2"
  _profileBegin
  local -a OPT_subVerbosity;
  [[ $verbosity =~ ^all ]] \
    && OPT_subVerbosity=(-v allBrief) \
    || OPT_subVerbosity=( "${OPT_verbosity[@]}" )
  local returnCode=0
  _mkDir ${OPT_subVerbosity[@]} ${OPT_tmp[@]} "$dirSpec"; returnCode=$?
  if [[ $returnCode -eq 0 ]]; then
    _pushd ${OPT_subVerbosity[@]} "$dirSpec"; returnCode=$?
    if [[ $returnCode -eq 0 ]]; then
      _exec ${OPT_verbosity[@]} ${OPT_silent[@]} ${OPT_treatAsOK[@]} ${OPT_preserveReturnCode[@]} "$inDirCommnand" "$@"; local returnCode=$?
      _popd ${OPT_subVerbosity[@]}
      [[ -z $tmp || -n $noCleanOnFail ]] || _rm ${OPT_subVerbosity[@]} -pd "$dirSpec"
    fi
  fi
  _profileEnd
  return $returnCode
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_pushdParams=(
  "${_verbosityParams[@]}"
  'dirSpec'
)
_pushd() { eval "$_funcParams2"
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  _exec ${OPT_verbosity[@]} ${OPT_silent[@]} pushd "$dirSpec" || return $?
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_popdParams=(
  "${_verbosityParams[@]}"
)
_popd() { eval "$_funcParams2"
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  _exec ${OPT_verbosity[@]} ${OPT_silent[@]} popd || return $?
}

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_rmParamsOpt=(--canBeMoreParams)
_rmParams=(
  "${_verbosityParams[@]}"
  '--dir/d'
  '--prepare/p'
)
_rm() { eval "$_funcParams2"
  _profileBegin
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  if [[ -z $dir ]]; then
    _exec ${OPT_verbosity[@]} ${OPT_silent[@]} rm -f "$@" || return $?
  else
    local spec; for spec in "$@"; do
      [[ -n $dir && ! $spec =~ /$ ]] && spec+='/'
      if [[ -n $prepare && -n $dir && -d $spec ]]; then
        _exec ${OPT_verbosity[@]} ${OPT_silent[@]} chmod 777 $(find "$spec" -type d) || return $?
      fi
      _exec ${OPT_verbosity[@]} ${OPT_silent[@]} rm -rfd "$spec" || return $?
    done
  fi
  _profileEnd
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_existParamsOpt=(--canBeMoreParams)
_existParams=(
  "${_verbosityParams[@]}"
  '--logic/l:(every any)=every'
  '--dir/d'
  '--no/n'
)
_exist() { eval "$_funcParams2"
  _profileBegin
  if [[ $verbosity == 'dry' ]]; then
    verbosity=none
  elif [[ $verbosity == 'allBrief' ]]; then
    verbosity=all
  fi
  local testOp; [[ -z $dir ]] && testOp=-f || testOp=-d
  local verboseType
  local spec; for spec in "$@"; do
    test $testOp "$spec"; local testReturnCode=$?
    local msgCmd
    if [[ $testReturnCode -ne 0 && -z $no || $testReturnCode -eq 0 && -n $no ]]; then
      [[ $verbosity == 'none' || $verbosity == 'ok' ]] || msgCmd=_err eval "$_existHelperCode"
      [[ $logic == 'every' ]] && return 1
    else
      [[ $verbosity == 'all' || $verbosity == 'ok' ]] && msgCmd=_ok eval "$_existHelperCode"
      [[ $logic == 'any' ]] && return 0
    fi
  done
  local returnCode; [[ $logic == 'every' ]] && returnCode=0 || returnCode=1
  _profileEnd
  return $returnCode
}
_existHelperCode='
  [[ -z $dir ]] && verboseType=Файл || verboseType=Директория
  local ansiHolder infix=
  if [[ $msgCmd == _ok ]]; then
    ansiHolder=_ansiOK
    [[ -n $no ]] && infix="не "
  else
    ansiHolder=_ansiErr
    [[ -n $no ]] || infix="не "
  fi
  $msgCmd "$verboseType ${_ansiFileSpec}$spec${!ansiHolder} ${infix}существует"
'

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_mkDirParamsOpt=( --canBeMoreParams )
_mkDirParams=(
  "${_verbosityParams[@]}"
  '--tmp/t'
)
_mkDir() { eval "$_funcParams2"
  _profileBegin
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  local dirSpec; for dirSpec in "$@"; do
    [[ -z $tmp ]] || _rm ${OPT_verbosity[@]} ${OPT_silent[@]} -pd "$dirSpec" || return $?
    local -a OPT_verbosityOfExist=(-v none); [[ $verbosity == 'none' || $verbosity == 'err' ]] || OPT_verbosityOfExist=(-v ok)
    [[ $verbosity != dry ]] && _exist ${OPT_verbosityOfExist[@]} -d "$dirSpec" \
      || _exec ${OPT_verbosity[@]} ${OPT_silent[@]} mkdir -p "$dirSpec" || return $?
  done
  _profileEnd
}

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_mvFileParams=(
  "${_verbosityParams[@]}"
  'from'
  'to'
)
_mvFile() { eval "$_funcParams2"
  _profileBegin
  if [[ $verbosity != 'dry' ]]; then
    local -a OPT_subVerbosity;
    if [[ $verbosity == 'none' || $verbosity == 'ok' ]]; then
      OPT_subVerbosity=(-v none)
    fi
    _exist ${OPT_subVerbosity[@]} $from || return $?
    _rm ${OPT_subVerbosity[@]} $to || return $?
    _mkDir ${OPT_subVerbosity[@]} $(dirname "$to") || return $?
  fi
  [[ $verbosity != 'all' ]] || OPT_verbosity=(-v allBrief)
  _exec ${OPT_verbosity[@]} ${OPT_silent[@]} mv "$from" "$to" || return $?
  _profileEnd
}

verbosityDefault=err silentDefault=yes codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_killAppParams=(
  "${_verbosityParams[@]}"
  'appName'
)
_killApp() { eval "$_funcParams2"
  local pid; for pid in $(ps -A | grep "$appName" | grep -v grep | awk '{print $1}'); do
    if _silent kill -0 $pid; then
      _exec ${OPT_verbosity[@]} ${OPT_silent[@]} kill $pid
    elif _silent _sudo kill -0 $pid; then
      _exec ${OPT_verbosity[@]} ${OPT_silent[@]} --sudo kill $pid
    fi
  done
}

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_dockerComposeParams=(
  "${_verbosityParams[@]}"
)
_dockerComposeParamsOpt=(--canBeMoreParams --treatUnknownOptionAsArg)
_dockerCompose() { eval "$_funcParams2"
  bw_install docker-compose --silentIfAlreadyInstalled || return $?
  local -a OPT=()
  if [[ $OSTYPE =~ ^linux ]]; then
    OPT+=( --sudo )
  elif [[ ! $OSTYPE =~ ^darwin ]]; then
    return $(_err "Неожиданный тип OS ${_ansiPrimaryLiteral}$OSTYPE")
  fi
  _exec ${OPT_verbosity[@]} ${OPT_silent[@]} "${OPT[@]}" docker-compose "$@"
}

verbosityDefault=err silentDefault=no codeHolder=_codeToPrepareVerbosityParams eval "$_evalCode"
_dockerParams=(
  "${_verbosityParams[@]}"
)
_dockerParamsOpt=(--canBeMoreParams --treatUnknownOptionAsArg)
_docker() { eval "$_funcParams2"
  bw_install docker --silentIfAlreadyInstalled || return $?
  local -a OPT=()
  if [[ $OSTYPE =~ ^linux ]]; then
    OPT+=( --sudo )
  elif [[ ! $OSTYPE =~ ^darwin ]]; then
    return $(_err "Неожиданный тип OS ${_ansiPrimaryLiteral}$OSTYPE")
  fi
  _exec ${OPT_verbosity[@]} ${OPT_silent[@]} "${OPT[@]}" docker "$@"; local returnCode=$?
  [[ $returnCode -eq 0 ]] || _debugVar returnCode
  return $returnCode
}

# =============================================================================

_whichParams=( 'cmd' )
_which() { eval "$_funcParams2"
  _silent command -v "$cmd"
}

# =============================================================================
