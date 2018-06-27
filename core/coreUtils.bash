
# =============================================================================

[[ $(type -t _resetBash) != function  ]] || _resetBash

# =============================================================================

_noop() { true; }

_gdate=date; [[ $OSTYPE =~ ^darwin ]] && _gdate=gdate

_profileStack=()
_profileStackIdx=0
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
    [[ -z $needProfile ]] || eval _startTime_$profileName=$($_gdate +%s%3N)
  fi
}

_gdateTiming=0

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
      local endTime=$($_gdate +%s%3N)
      local elapsedTime=$(( endTime - ${!startTimeHolder} ))
      if [[ $_profileStackIdx -gt 0 ]]; then
        local parentProfileStackIdx=$(( _profileStackIdx - 1 ))
        _profileStack[$parentProfileStackIdx]=$(( _profileStack[$parentProfileStackIdx] + elapsedTime + _gdateTiming ))
      fi
      elapsedTime=$(( elapsedTime - _profileStack[$_profileStackIdx] - _gdateTiming ))
      [[ $elapsedTime -ge 0 ]] || elapsedTime=0
      eval _profileTotal_$profileName=$(( _profileTotal_$profileName + elapsedTime ))
      eval _profileCount_$profileName=$(( _profileCount_$profileName + 1 ))
    fi
  fi
}

_profileInit() {
  local prefix=_profileTotal_
  local varName; for varName in $(compgen -v | grep "^$prefix" ); do # https://unix.stackexchange.com/questions/3510/how-to-print-only-defined-variables-shell-and-or-environment-variables-in-bash/5691#5691
    [[ -n $varName ]] || continue
    eval $varName=0
    local funcName=${varName:${#prefix}}
    eval _profileCount_$funcName=0
  done

  _gdateTiming=$($_gdate +%s%3N)
  for ((_i=0; _i<30; _i++)); do
    $_gdate +%s%3N >/dev/null
  done
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
        echo $funcName: ${!varName}ms ${!countHolder} $(( ${!varName} / ${!countHolder} ))ms
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
    . "$_profileTmpFileSpec"
  fi
}

_noSleepWhileProfile=true
_funcToProfile() { _profileBegin;
  [[ -n $_noSleepWhileProfile ]] || sleep 1
  _subFuncToProfile
  _subFuncToProfile2
_profileEnd; }

_subFuncToProfile() { _profileBegin;
  [[ -n $_noSleepWhileProfile ]] || sleep .3
  _subSubFuncToProfile
_profileEnd; }

_subSubFuncToProfile() { _profileBegin;
  [[ -n $_noSleepWhileProfile ]] || sleep .1
  _subSubSubFuncToProfile
_profileEnd; }

_subSubSubFuncToProfile() { _profileBegin;
  [[ -n $_noSleepWhileProfile ]] || sleep .1
_profileEnd; }

_subFuncToProfile2() { _profileBegin;
  [[ -n $_noSleepWhileProfile ]] || sleep .5
_profileEnd; }

_useOptimizedHasItem=
# _useOptimizedHasItem=true

_hasItemAccumTime=0
_hasItem() {
  _profileBegin
  local testItem="$1"; shift
  local returnCode=1
  local item; for item in "$@"; do
    if [[ "$item" == "$testItem" ]]; then
      returnCode=0
      break
    fi
  done
  _profileEnd
  return $returnCode
}

# _hasSuffix() {
#   _profileBegin
#   local testItem="$1"; shift
#   local testItemLen=${#testItem}
#   local returnCode=1
#   local etaSuffix; for etaSuffix in "$@"; do
#     local len=${#etaSuffix}
#     [[ $len -lt $testItemLen ]] || continue
#     local tstSuffix="${testItem:$((testItemLen - len)):$len}"
#     [[ "$etaSuffix" == "$tstSuffix" ]] || continue 
#     returnCode=0
#     break
#   done
#   _profileEnd
#   return $returnCode
# }

# _hasPrefix() {
#   _profileBegin
#   local testItem="$1"; shift
#   local testItemLen=${#testItem}
#   local returnCode=1
#   local etaPrefix; for etaPrefix in "$@"; do
#     local len=${#etaPrefix}
#     [[ $len -lt $testItemLen ]] || continue
#     local tstPrefix="${testItem:0:$len}"
#     [[ "$etaPrefix" == "$tstPrefix" ]] || continue 
#     returnCode=0
#     break
#   done
#   _profileEnd
#   return $returnCode
# }

_quotedArgs() {
  local strip quote
  while [[ $1 =~ ^--(strip|quote[:=](dollarSign|all)) ]]; do
    if [[ ${BASH_REMATCH[1]} == 'strip' ]]; then
      strip=true
    else
      quote=${BASH_REMATCH[2]}
    fi
    shift
  done
  local arg joiner result; for arg in "$@"; do
    if ! [[
      -z $arg ||
      $arg =~ [[:space:]] ||
      $arg =~ \" ||
      $arg =~ \' ||
      $arg =~ \\ ||
      $arg =~ [*?+] ||
      $quote == 'dollarSign' && $arg =~ \$ ||
      $quote == 'all'
    ]]; then
      result+="$joiner$arg"
    else
      local quotedArg=$(declare -p arg )
      [[ $quotedArg =~ ^([^=]*=) ]] \
        && quotedArg=${quotedArg:${#BASH_REMATCH[0]}}
      [[ -n $strip ]] \
        && quotedArg=${quotedArg:1:$((${#quotedArg} - 2))}
      result+="$joiner$quotedArg"
    fi
    joiner=' '
  done
  echo "$result"
}

_upperFirst() {
  echo "$*" | perl -C -pe 's/^(.)/\U$1/'
}

_lowerFirst() {
  echo "$*" | perl -C -pe 's/^(.)/\L$1/'
}

_kebabCaseToCamelCase() { # https://stackoverflow.com/questions/34420091/spinal-case-to-camel-case/34420162#34420162
  if [[ ! $1 =~ - ]]; then
    # echo dstVarName: $dstVarName, \$1: $1
    eval $dstVarName=\"\$1\"
  else
    _profileBegin
    local result=$(echo $1 | perl -pe 's/(-)(\w)/\U$2/g')
    eval $dstVarName=\$result
    _profileEnd
  fi
}


# _kebabCaseToUpperCamelCase() { # https://stackoverflow.com/questions/34420091/spinal-case-to-camel-case/34420162#34420162
#   _profileBegin
#   echo $1 | perl -pe 's/(^|-)(\w)/\U$2/g'
#   _profileEnd
# }

_upperCamelCaseToKebabCase() { # https://stackoverflow.com/questions/28795479/awk-sed-script-to-convert-a-file-from-camelcase-to-underscores/28795550#28795550
  if [[ ! $1 =~ [A-Z] ]]; then
    eval $dstVarName=$1
  else
    local holder=__upperCamelCaseToKebabCase_$1_
    if [[ -n ${!holder} ]]; then
      eval $dstVarName=\${!holder}
    else
      _profileBegin
      local result=$(echo $1 | perl -pe 's/(?<=.)([A-Z])/-\L$1/g; s/^(.)/\L$1/')
      eval $dstVarName=\$result
      eval $holder=\$result
      _profileEnd
    fi
  fi
}

_funcExists() {
  [[ $(type -t $1) == 'function' ]]
}

_getUniqArray() {
  # https://superuser.com/questions/461981/how-do-i-convert-a-bash-array-variable-to-a-string-delimited-with-newlines/462400#462400
  # https://stackoverflow.com/questions/13648410/how-can-i-get-unique-values-from-an-array-in-bash/13648438#13648438
  local unique=$(IFS=$'\n'; echo "$*" | sort -u)
  # https://stackoverflow.com/questions/11393817/bash-read-lines-in-file-into-an-array/11393884#11393884
  IFS=$'\r\n' GLOBIGNORE='*' command eval 'local -a arr=( $(echo "$unique") )'
  echo $(_quotedArgs "${arr[@]}")
}

_everyFileNotNewerThan() {
  local etaFileSpec="$1"; shift
  local tstFileSpec; for tstFileSpec in "$@"; do
    [[ $tstFileSpec != $_bwFileName ]] || tstFileSpec="$_bwFileSpec"
    [[ ! $tstFileSpec -nt $etaFileSpec ]] || return $?
  done
  return 0
}

_spaceContainer=; for ((_i=0; _i<1024; _i++)); do _spaceContainer+=" "; done
_indent() {
  echo -n "${_spaceContainer:0:${1:0}}"
}

_restore() {
  _profileBegin;
  local varName; for varName in "$@"; do
    local typeOfVar; eval "$_codeToPrepareTypeOfVar"
    if [[ $typeOfVar == 'none' ]]; then
      return $(_err "could not resolve type of ${_ansiOutline}$varName${_ansiErr}, first declare it with initial value")
    else
      local idxVarName="${_substitutePrefix}${varName}${_substituteIdxSuffix}"
      if [[ -z ${!idxVarName} ]]; then
        return $(_err "can not restore ${_ansiOutline}$varName${_ansiErr} which is not substituted before")
      fi
      local valueToRestoreVarName="${_substitutePrefix}${varName}${_substituteValueSuffix}${!idxVarName}"
      if [[ $typeOfVar == 'scalar' ]]; then
        eval "$varName=\$$valueToRestoreVarName"
      elif [[ $typeOfVar == 'array' ]]; then
        eval $varName=\( \"\${$valueToRestoreVarName[@]}\" \)
      fi
      if [[ ${!idxVarName} -le 0 ]]; then
        eval "$idxVarName="
      else
        eval "$idxVarName=\$(( $idxVarName - 1 ))"
      fi
    fi
  done
  _profileEnd;
}

_codeToPrepareTypeOfVar='
  local __typeSignature=$(declare -p $varName 2>/dev/null)
  if [[ $__typeSignature =~ ^declare[[:space:]]-a ]]; then
    typeOfVar=array
  elif [[ -z ${!varName+x} ]]; then
    typeOfVar=none
  else
    typeOfVar=scalar
  fi
'

_lcpDescription="longest common prefix, from ${_ansiUrl}https://rosettacode.org/wiki/Longest_common_prefix#Perl${_ansiReset}"
_lcp() {
  perl -e "print scalar ((join(\"\\0\", @ARGV) =~ /^([^\\0]*)[^\\0]*(?:\\0\\1[^\\0]*)*\$/s)[0])" "$@"
}

noStack=
_debugStack() {
  if [[ -n $noStack ]]; then
    echo >&2
  else
    local ofs=$1; [[ $ofs =~ ^[1-9][0-9]* ]] || ofs=2
    echo -n "$___joiner${_ansiOutline}STACK${_ansiReset}($(( ${#FUNCNAME[@]} - $ofs ))):" >&2
    local idx; for (( idx=$ofs; idx<${#FUNCNAME[@]}; idx++ )); do
      echo -n " ${_ansiCmd}${FUNCNAME[$idx]}${_ansiReset}${_ansiDim}@${_ansiResetDim}${_ansiCmd}${BASH_SOURCE[$idx]}${_ansiReset}" >&2
    done
    echo "${_ansiReset}" >&2
  fi
}

_warn() {
  echo "${_ansiWarn}WARN: $@${_ansiReset}" >&2
}

_todo() {
  echo -n "$prefix${_ansiErr}TODO in ${_ansiCmd}${FUNCNAME[1]}${_ansiErr}: $@${_ansiReset}" >&2
  _debugStack
  return 3
}

# _existsAnyFileOfTemplate() {
#   local fileSpecTemplate="$1"
#   # https://stackoverflow.com/questions/6363441/check-if-a-file-exists-with-wildcard-in-shell-script/7702334#7702334
#   test -n "$(find "$(dirname "$fileSpecTemplate")" -maxdepth 1 -name "$(basename  "$fileSpecTemplate")" -print -quit)"
# }



_getExternalIp() {
  curl ipecho.net/plain 
}      

_getOwnIpList() {
  # https://stackoverflow.com/questions/13322485/how-to-get-the-primary-ip-address-of-the-local-machine-on-linux-and-os-x
  if which hostname >/dev/null 2>&1; then
    hostname -I
  else
    local useSedVersion=
    if [[ -z $useSedVersion ]]; then
      ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1'
    else
      ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p'
    fi
  fi
}

# _getOwnIp() {
  # https://stackoverflow.com/questions/38252963/how-can-i-debug-php-mounted-to-a-container-running-on-docker-beta-for-mac
  # ipconfig getifaddr en0

  # https://stackoverflow.com/questions/13322485/how-to-get-the-primary-ip-address-of-the-local-machine-on-linux-and-os-x/13322667#13322667
  # local ip myip line
  # while IFS=$': \t' read -a line ;do
  #   _debugVar line
  #   [ -z "${line%inet}" ] && ip=${line[${#line[1]}>4?1:2]} && [ "${ip#127.0.0.1}" ] && myip=$_ip
  # done< <(LANG=C /sbin/ifconfig)
  # echo $myip
  # printf ${1+-v} $1 "%s${_nl:0:$[${#1}>0?0:1]}" $myip
# }


# =============================================================================

_isInDocker() {
  [[ -f /.dockerenv ]]
}

# =============================================================================

