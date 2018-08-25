
# =============================================================================

_resetBash

# =============================================================================

_OPT_profileFileSpec=( '--profileFileSpec' 'profile' )
_dollarSign='\$'
_setAtBashProfileTests=(
	'
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last thing >> profile"

    --before "touch profile.bak0"
    --before "rm -f profile.bak1"
    --before "cp profile profile.bak1.eta"

    --before "echo first thing > profile.eta"
    --before "echo ${_dollarSign}mid thing >> profile.eta"
    --before "echo last thing >> profile.eta"

    --fileSpec profile
    --file2Spec profile.bak1

    "_setAtBashProfile \"\\\$mid thing\" \"some\" ${_OPT_profileFileSpec[*]}"

  '
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last thing >> profile"

    --before "touch profile.bak0 profile.bak1 profile.bak2 profile.bak3 profile.bak4 profile.bak5 profile.bak6 profile.bak7 profile.bak8 profile.bak9"
    --before "rm -f profile.bak10"
    --before "cp profile profile.bak10.eta"

    --before "echo mid thing > profile.eta"
    --fileSpec profile
    --file2Spec profile.bak10

    "_setAtBashProfile \"mid thing\" \"thing\" ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last thing >> profile"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --before "cp profile profile.eta"
    --before "echo one more thing >> profile.eta"

    --fileSpec profile
    --file2Spec profile.bak0

    "_setAtBashProfile \"one more thing\" \"another\" ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last thing >> profile"

    --before "rm -f profile.bak0"

    --before "cp profile profile.eta"

    --fileSpec profile

    --fileNotExist profile.bak0

    "_setAtBashProfile \"some thing\" \"some thing\" ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo \". ~/lab/crm/bin/crm.bash; crm update completionOnly\" > profile.eta"

    --before "echo \". ~/crm/bin/crm.bash; crm update completionOnly\" > profile"
    --before "echo \". ~/work/crm/bin/crm.bash; crm update completionOnly\" >> profile"
    --before "echo \". ~/billing-gate/bin/bgate.bash; bgate update completionOnly\" | tee -a profile >> profile.eta"

    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0

    "_setAtBashProfile '\''. ~/lab/crm/bin/crm.bash; crm update completionOnly'\'' '\''${_sourceMatchRegexp}bin\\\/crm\\\.bash'\'' ${_OPT_profileFileSpec[*]}"
  '

 # --uninstall 
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last thing >> profile"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --before "grep -v some profile > profile.eta"

    --fileSpec profile
    --file2Spec profile.bak0

    "_setAtBashProfile \"some thing\" \"some thing\" --uninstall ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo some thing >> profile"
    --before "echo last Thing >> profile"

    --before "echo last Thing >> profile.eta"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0

    "_setAtBashProfile \"some thing\" \"thing\" --uninstall ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"

    --before "rm -f profile.bak0"

    --before "cp profile profile.eta"

    --fileSpec profile

    --fileNotExist profile.bak0

    "_setAtBashProfile \"some thing\" \"some thing\" --uninstall ${_OPT_profileFileSpec[*]}"
  '
  '
    --inTmpDir

    --before "echo \". ~/crm/bin/crm.bash; crm update completionOnly\" > profile"
    --before "echo \". ~/work/crm/bin/crm.bash; crm update completionOnly\" >> profile"
    --before "echo \". ~/billing-gate/bin/bgate.bash; bgate update completionOnly\" | tee -a profile > profile.eta"

    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0

    "_setAtBashProfile '\''. ~/lab/crm/bin/crm.bash; crm update completionOnly'\'' '\''${_sourceMatchRegexp}bin\\\/crm\\\.bash'\'' --uninstall ${_OPT_profileFileSpec[*]}"
  '
)

_hasAtBashProfileTests=(
  '
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile thing thing ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 2

    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile thing thing --differ ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 1
    
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile some some --differ ${_OPT_profileFileSpec[*]}"
  '
# --no
  '
    --return 1

    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile thing thing --no ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 2

    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile thing thing --no --differ ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 0
    
    --inTmpDir

    --before "echo first thing > profile"
    --before "echo last thing >> profile"
    
    "_hasAtBashProfile some some --no --differ ${_OPT_profileFileSpec[*]}"
  '
)

_exportVarAtBashProfileTests=(
  '
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile.eta"

    --before "rm -f profile profile.bak0"

    --before "SOME_VAR=\"some value\""

    --fileSpec profile
    
    "_exportVarAtBashProfile SOME_VAR ${_OPT_profileFileSpec[*]}"

  '
  '
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"

    --before "SOME_VAR=newValue"
    --before "echo \"export SOME_VAR=newValue # by bw.bash\" > profile.eta"

    --before "echo \"export ANOTHER_VAR=\\\"Another value\\\" # by bw.bash\" | tee -a profile >> profile.eta"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0
    --noErrorStack

    "_exportVarAtBashProfile SOME_VAR ${_OPT_profileFileSpec[*]}"

  '
  '
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\"\" | tee -a profile > profile.eta"

    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    --before "SOME_VAR=newValue"
    --before "echo \"export SOME_VAR=newValue # by bw.bash\" >> profile.eta"

    --before "echo \"export ANOTHER_VAR=\\\"Another value\\\" # by bw.bash\" | tee -a profile  >> profile.eta"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0
    --noErrorStack

    "_exportVarAtBashProfile SOME_VAR ${_OPT_profileFileSpec[*]}"

  '

 # --uninstall 

  '
    --inTmpDir

    --before "SOME_VAR=newValue"

    --before "echo \"export SOME_VAR=\"some value\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=newValue # by bw.bash\" > profile"
    --before "echo \"export ANOTHER_VAR=\"another value\" # by bw.bash\" | tee -a profile  > profile.eta"

    --before "rm -f profile.bak0"
    --before "cp profile profile.bak0.eta"

    --fileSpec profile
    --file2Spec profile.bak0
    --noErrorStack

    "_exportVarAtBashProfile SOME_VAR --uninstall ${_OPT_profileFileSpec[*]}"

  '
)

_hasExportVarAtBashProfileTests=(
  '
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile SOME_VAR ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 2

    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile SOME_VAR --differ ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 1
    
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile ANOTHER_VAR ${_OPT_profileFileSpec[*]}"
  '
# --no
  '
    --return 1

    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile SOME_VAR --no ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 2

    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile SOME_VAR --no --differ ${_OPT_profileFileSpec[*]}"
  '
  '
    --return 0
    
    --inTmpDir

    --before "echo \"export SOME_VAR=\\\"some value\\\" # by bw.bash\" > profile"
    --before "echo \"export SOME_VAR=\\\"another value\\\" # by bw.bash\" >> profile"
    
    "_hasExportVarAtBashProfile ANOTHER_VAR --no ${_OPT_profileFileSpec[*]}"
  '
)

# =============================================================================
