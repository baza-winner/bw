
# =============================================================================

[[ $(type -t _resetBash) != function  ]] || _resetBash

# =============================================================================

# http://ezprompt.net/
export _psCaughtErrorCodeFileSpec='/tmp/bw.ps.errorCode'
export _psGitBranchCodeFileSpec='/tmp/bw.ps.gitBranch'

_psPrepare_error() {
  local errorCode=$?
  if [[ $errorCode -ne 0 ]]; then
    echo "local returnCode=$errorCode" > "$_psCaughtErrorCodeFileSpec"
  elif [[ -f $_psCaughtErrorCodeFileSpec ]]; then
    rm "$_psCaughtErrorCodeFileSpec" >/dev/null 2>&1
  fi
}

_codeToEchoOutput='
  if [[ $1 == both ]]; then
    output=" $output "
  elif [[ $1 == before ]]; then
    output=" $output"
  elif [[ $1 == after ]]; then
    output="$output "
  fi
  echo -n "$output"
'
_psIf_error() {
  if [[ -f $_psCaughtErrorCodeFileSpec ]]; then
    local output="$2"
    eval "$_codeToEchoOutput"
  fi
}

_ps_errorCode() {
  if [[ -f $_psCaughtErrorCodeFileSpec ]]; then
    . "$_psCaughtErrorCodeFileSpec" >/dev/null 2>&1
    local output="$returnCode"
    eval "$_codeToEchoOutput"
  fi
}

_psPrepare_git() {
  local branch=$(_gitBranch)
  if [[ -n $branch ]]; then
    echo "local branch=$branch" > "$_psGitBranchCodeFileSpec"
  elif [[ -f $_psGitBranchCodeFileSpec ]]; then
    rm "$_psGitBranchCodeFileSpec" >/dev/null 2>&1
  fi
}

_ps_gitBranch() {
  if [[ -f $_psGitBranchCodeFileSpec ]]; then
    . "$_psGitBranchCodeFileSpec" >/dev/null 2>&1
    local output="$branch"
    eval "$_codeToEchoOutput"
  fi
}

_ps_gitDirty() {
  if [[ -f $_psGitBranchCodeFileSpec ]]; then
    local status=$(_gitDirty)
    if [[ -n $status ]]; then
      local output="$status"
      eval "$_codeToEchoOutput"
    fi
  fi
}

_psIf_git() {
  if [[ -f $_psGitBranchCodeFileSpec ]]; then
    local output="$2"
    eval "$_codeToEchoOutput"
  fi
}

export _psColorSegmentBeginPrefix='\['
export _psColorSegmentBeginSuffix='\]'
export _psColorSegmentEnd='\[\e[m\]'
export _ps_user='\u'
# export _ps_user_description='Имя пользователя'
export _ps_host='\h'
export _psFullQualifiedDomainName='\H'
export _ps_fqdn="$_psFullQualifiedDomainName"
export _ps_fqdn_description='FullQualifiedDomainName'
export _ps_shell='\s'
export _ps_shellVersion='\v'
export _ps_shellRelease='\V'
export _psPathToCurrentDirectory='\w'
export _ps_ptcd="$_psPathToCurrentDirectory"
export _ps_ptcd_description="PathToCurrentDirectory"
export _psCurrentDirectory='\W'
export _ps_cd="$_psCurrentDirectory"
export _ps_cd_description="CurrentDirectory"

# =============================================================================

