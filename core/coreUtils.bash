#!/bin/bash

{
  # shellcheck disable=SC2154
  if [[ -n $shellcheck ]]; then 
    . bw.bash
    . core/profileSupport.bash
  fi
}

# =============================================================================

_resetBash

# =============================================================================

_noop() { true; }

# =============================================================================

_hasItem() {
  #shellcheck disable=SC2119
  _profileBegin
  local testItem="$1"; shift
  local returnCode=1
  local item; for item in "$@"; do
    if [[ "$item" == "$testItem" ]]; then
      returnCode=0
      break
    fi
  done
  #shellcheck disable=SC2119
  _profileEnd
  return $returnCode
}

_joinBy() { 
  # https://stackoverflow.com/questions/1527049/join-elements-of-an-array
  local separator="$1"; shift
  local needSeparator=
  while [[ $# -gt 0 ]]; do
    [[ -z $needSeparator ]] || echo -n "$separator"
    if [[ $# -eq 1 ]]; then
      echo "$1"
    else
      echo -n "$1"
      needSeparator=true
    fi
    shift
  done
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
      local quotedArg; quotedArg=$(declare -p arg)
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
  local string="$*"
  echo "$(tr '[:lower:]' '[:upper:]' <<< "${string:0:1}")${string:1}"
}

_lowerFirst() {
  local string="$*"
  echo "$(tr '[:upper:]' '[:lower:]' <<< "${string:0:1}")${string:1}"
}

_kebabCaseToCamelCase() { # https://stackoverflow.com/questions/34420091/spinal-case-to-camel-case/34420162#34420162
  if [[ ! $1 =~ - ]]; then
    eval "${dstVarName?}"'="$1"'
  else
    #shellcheck disable=SC2119
    _profileBegin
    local result; result=$(echo "$1" | perl -pe 's/(-)(\w)/\U$2/g')
    eval "${dstVarName?}"'=$result'
    #shellcheck disable=SC2119
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
    eval "${dstVarName?}=$1"
  else
    local holder=__upperCamelCaseToKebabCase_$1_
    if [[ -n ${!holder} ]]; then
      eval "${dstVarName?}"'=${!holder}'
    else
      #shellcheck disable=SC2119
      _profileBegin
      local result; result=$(echo "$1" | perl -pe 's/(?<=.)([A-Z])/-\L$1/g; s/^(.)/\L$1/')
      eval "${dstVarName?}"'=$result'
      eval "$holder"'=$result'
      #shellcheck disable=SC2119
      _profileEnd
    fi
  fi
}

_funcExists() {
  [[ $(type -t "$1") == 'function' ]]
}

_getUniqArray() {
  # https://superuser.com/questions/461981/how-do-i-convert-a-bash-array-variable-to-a-string-delimited-with-newlines/462400#462400
  # https://stackoverflow.com/questions/13648410/how-can-i-get-unique-values-from-an-array-in-bash/13648438#13648438
  local unique
  # shellcheck disable=SC2034
  unique=$(IFS=$'\n'; echo "$*" | sort -u)
  # https://stackoverflow.com/questions/11393817/bash-read-lines-in-file-into-an-array/11393884#11393884
  local -a arr
  # shellcheck disable=SC2016
  IFS=$'\r\n' GLOBIGNORE='*' command eval 'arr=( $(echo "$unique") )'
  _quotedArgs "${arr[@]}"
}

_everyFileNotNewerThan() {
  local etaFileSpec="$1"; shift
  local tstFileSpec; for tstFileSpec in "$@"; do
    [[ $tstFileSpec != "$_bwFileName" ]] || tstFileSpec="$_bwFileSpec"
    [[ ! $tstFileSpec -nt $etaFileSpec ]] || return $?
  done
  return 0
}

_i=
_spaceContainer=
for ((_i=0; _i<1024; _i++)); do _spaceContainer+=" "; done
_indent() {
  echo -n "${_spaceContainer:0:${1:0}}"
}

_substituteValueSuffix='_VALUE'
_substituteIdxSuffix='_IDX'
_restore() {
  #shellcheck disable=SC2119
  _profileBegin
  local varName; for varName in "$@"; do
    local typeOfVar; typeOfVar=$(_getTypeOfVar "$varName")
    # shellcheck disable=SC2046
    [[ $typeOfVar != 'none' ]] \
      || return $(_err "could not resolve type of ${_ansiOutline}$varName${_ansiErr}, first declare it with initial value")
    local idxVarName="${_substitutePrefix}${varName}${_substituteIdxSuffix}"
    # shellcheck disable=SC2046
    [[ -n ${!idxVarName} ]] \
      || return $(_err "can not restore ${_ansiOutline}$varName${_ansiErr} which is not substituted before")
    local valueToRestoreVarName="${_substitutePrefix}${varName}${_substituteValueSuffix}${!idxVarName}"
    local codeToEval="$varName="
    if [[ $typeOfVar == 'scalar' ]]; then
      codeToEval+='$'"$valueToRestoreVarName"
    elif [[ $typeOfVar == 'array' ]]; then
      # shellcheck disable=SC2016
      codeToEval+='( "${'"$valueToRestoreVarName"'[@]}" )'
    fi
    codeToEval+="$_nl$idxVarName="
    # shellcheck disable=SC2016
    [[ ${!idxVarName} -le 0 ]] || codeToEval+='$(( idxVarName - 1 ))'
    eval "$codeToEval"
  done
  #shellcheck disable=SC2119
  _profileEnd
}

_getTypeOfVar() {
  local __varName="$1"
  local __typeSignature; __typeSignature=$(declare -p "$__varName" 2>/dev/null)
  if [[ $__typeSignature =~ ^declare[[:space:]]-a ]]; then
    echo array
  elif [[ -z ${!__varName+x} ]]; then
    echo none
  else
    echo scalar
  fi
}

# shellcheck disable=SC2034
_lcp_description="longest common prefix, from ${_ansiUrl}https://rosettacode.org/wiki/Longest_common_prefix#Perl${_ansiReset}"
_lcp() {
  perl -e "print scalar ((join(\"\\0\", @ARGV) =~ /^([^\\0]*)[^\\0]*(?:\\0\\1[^\\0]*)*\$/s)[0])" "$@"
}

_noStack=
# shellcheck disable=SC2120
_debugStack() {
  if [[ -n $_noStack ]]; then
    echo >&2
  else
    local ofs=$1; [[ $ofs =~ ^[1-9][0-9]* ]] || ofs=2
    # shellcheck disable=SC2154
    echo -n "$___joiner${_ansiOutline}STACK${_ansiReset}($(( ${#FUNCNAME[@]} - ofs ))):" >&2
    local idx; for (( idx=ofs; idx<${#FUNCNAME[@]}; idx++ )); do
      echo -n " ${_ansiCmd}${FUNCNAME[$idx]}${_ansiReset}${_ansiDim}@${_ansiResetDim}${_ansiCmd}${BASH_SOURCE[$idx]}${_ansiReset}" >&2
    done
    echo "${_ansiReset}" >&2
  fi
}

_warn() {
  echo "${_ansiWarn}WARN: $*${_ansiReset}" >&2
}

_todo() {
  # shellcheck disable=SC2154
  echo -n "$prefix${_ansiErr}TODO in ${_ansiCmd}${FUNCNAME[1]}${_ansiErr}: $*${_ansiReset}" >&2
  # shellcheck disable=SC2119
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
  if command -v hostname >/dev/null 2>&1; then
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
