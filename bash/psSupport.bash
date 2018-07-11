
# =============================================================================

_resetBash

# =============================================================================

# http://ezprompt.net/
_psCaughtErrorCodeFileSpec='/tmp/bw.ps.errorCode'
_psGitBranchCodeFileSpec='/tmp/bw.ps.gitBranch'

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

_psColorSegmentBeginPrefix='\['
_psColorSegmentBeginSuffix='\]'
_psColorSegmentEnd='\[\e[0m\]'
_ps_user='\u'
# export _ps_user_description='Имя пользователя'
_ps_host='\h'
_psFullQualifiedDomainName='\H'
_ps_fqdn="$_psFullQualifiedDomainName"
_ps_fqdn_description='FullQualifiedDomainName'
_ps_shell='\s'
_ps_shellVersion='\v'
_ps_shellRelease='\V'
_psPathToCurrentDirectory='\w'
_ps_ptcd="$_psPathToCurrentDirectory"
_ps_ptcd_description="PathToCurrentDirectory"
_psCurrentDirectory='\W'
_ps_cd="$_psCurrentDirectory"
_ps_cd_description="CurrentDirectory"

# =============================================================================

