
# =============================================================================

_resetBash

# =============================================================================

_stq=\'
_stDollarInQ='\$'
_stOpenBraceInQ='('
_stCloseBraceInQ=')'
_stSlashInQ='\'
_stSlashInQQ='\\\\'
_stqq='\"'

# =============================================================================

bw_bashTestsHelper() {
  local -a funcNamesToPregen=()
  local -a fileNamesToProcess=()
  local -a funcsWithTests=();

  for _fileSpec in "$testsDirSpec/testsSupport.bash"; do
    fileNamesToProcess+=( $(basename "$_fileSpec") )
    funcNamesToPregen+=( $(_getFuncNamesOfScriptToUnset "$_fileSpec") )
    codeHolder=_codeSource eval "$_evalCode"
  done

  if [[ ${#args[@]} -eq 0 ]]; then
    for _fileSpec in "$testsDirSpec/"*Tests.bash; do
      fileNamesToProcess+=( $(basename "$_fileSpec") )
      funcNamesToPregen+=( $(_getFuncNamesOfScriptToUnset "$_fileSpec") )
      codeHolder=_codeSource eval "$_evalCode"
    done
  else
    local -a funcsTestFor=()
    local arg; for arg in "${args[@]}"; do
      [[ ! $arg =~ $_isValidVarNameRegExp ]] || funcsTestFor+=( "$arg" )
    done
    funcsTestFor=( $(_getUniqArray "${funcsTestFor[@]}") )
    for funcTestFor in ${funcsTestFor[@]}; do
      for _fileSpec in "$testsDirSpec/"*Tests.bash; do
        local fileNameToProcess=$(basename "$_fileSpec")
        ! _hasItem "$fileNameToProcess" "${fileNamesToProcess[@]}" || continue
        _hasItem ${funcTestFor}Tests $(_getVarNamesOfScriptToUnset "$_fileSpec") || continue
        fileNamesToProcess+=( "$fileNameToProcess" )
        funcNamesToPregen+=( $(_getFuncNamesOfScriptToUnset "$_fileSpec") )
        codeHolder=_codeSource eval "$_evalCode"
      done
    done
  fi

  if [[ -n $_isBwDevelop ]]; then
    local -a __completions=()
    if [[ -n $noTiming ]]; then
      _pregen ${funcNamesToPregen[@]} || return $?
    else
      _spinner \
        -t "${_ansiOK}OK: ${_ansiFileSpec}$testsDirSpec${_ansiReset}: ${_ansiFileSpec}${fileNamesToProcess[*]}${_ansiReset} $(_getPluralWord ${#fileNamesToProcess[@]} обработан обработаны) за" \
        "${_ansiHeader}Прегенерация для файлов из ${_ansiFileSpec}$testsDirSpec${_ansiReset}" \
        _pregen ${funcNamesToPregen[@]} || return $?
    fi
  fi

  for _fileSpec in $(find "$testsDirSpec/$_generatedDir" -name *.completion$_codeBashExt); do
    codeHolder=_codeSource eval "$_evalCode"
  done

  local -a funcsWithTests=(); _prepareFuncWithTests
  local wasError
  local -a testsToRun=()
  if [[ ${#args[@]} -eq 0 ]]; then
    testsToRun=( "${funcsWithTests[@]}" )
  else
    local funcTestFor=
    local wasDetailsForFuncTestFor=
    local wasErrorForFuncTestFor=
    local arg; for arg in "${args[@]}"; do
      if [[ $arg =~ $_isValidVarNameRegExp ]]; then
        [[ -z $funcTestFor || -n $wasDetailsForFuncTestFor || -n $wasErrorForFuncTestFor ]] \
          || testsToRun+=( "$funcTestFor" )
        funcTestFor="$arg"
        wasDetailsForFuncTestFor=
        wasErrorForFuncTestFor=
        _hasItem $funcTestFor "${funcsWithTests[@]}" || {
          _err "Нет тестов для функции ${_ansiCmd}$funcTestFor"
          wasError=true
          wasErrorForFuncTestFor=true
          continue
        }
      else
        local diapMin= diapMax= diapErr=
        diapHolder=arg codeHolder=_codeToParseDiap eval "$_evalCode"

        if [[ -z $diapMin && -z $diapMax ]]; then
          origin=2 _throw "ожидает в качестве ${_ansiOutline}Аргумента${_ansiErr} имя функции или диапазон номеров тестов вместо ${_ansiPrimaryLiteral}$arg"
          wasError=true
        elif [[ -z $wasErrorForFuncTestFor && -z $list ]]; then
          if [[ -n $diapErr ]]; then
            _ownThrow "$diapErr"
            wasError=true
          elif [[ -z $funcTestFor ]]; then
            origin=2 _throw "ожидает в качестве ${_ansiOutline}Аргумента${_ansiErr} имя функции вместо дипазона номеров тестов ${_ansiPrimaryLiteral}$arg${_ansiErr}"
            wasError=true
          else
            local num=$diapMin
            dstVarName=selfTests srcVarName=${funcTestFor}Tests codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
            local count="${#selfTests[@]}"
            local maxIndex=$(( $count - 1 ))
            local minIndex=-$count
            local num=
            [[ $diapMax == Infinite || $diapMax -le $maxIndex ]] \
              || num=$diapMax
            [[ $diapMin == Infinite || $diapMin -ge $minIndex ]] \
              || num=$diapMin
            if [[ -n $num ]]; then
              local diap="${_ansiSecondaryLiteral}от 0 до $maxIndex${_ansiErr} или ${_ansiSecondaryLiteral}от -1 до -$count"
              [[ $maxIndex -ne 0 ]] || diap="${_ansiPrimaryLiteral}0"
              _ownThrow "ожидает, что ${_ansiOutline}Номер теста${_ansiErr} ${_ansiCmd}$funcTestFor${_ansiErr} должен быть $diap${_ansiErr}, a не ${_ansiPrimaryLiteral}$num"
              wasError=true
            else
              if [[ $diapMin == Infinite && $diapMax == Infinite ]]; then
                testsToRun+=( "$funcTestFor" )
              else
                if [[ $diapMax == Infinite ]]; then
                  [[ $diapMin -ge 0 ]] \
                    && diapMax=$maxIndex \
                    || diapMax=-1
                elif [[ $diapMin == Infinite ]]; then
                  [[ $diapMax -ge 0 ]] \
                    && diapMin=0 \
                    || diapMin=$minIndex
                fi
                for (( num=$diapMin; num<=$diapMax; num++ )); do
                  testsToRun+=( "$funcTestFor $num" )
                done
              fi
              wasDetailsForFuncTestFor=true
            fi
          fi
        fi
      fi
    done
    [[ -z $funcTestFor || -n $wasDetailsForFuncTestFor || -n $wasErrorForFuncTestFor ]] \
      || testsToRun+=( "$funcTestFor" )
  fi

  if [[ -n $list ]]; then
    local -a alreadyListed=()
    local totalCount=0
    local funcTestFor; for funcTestFor in "${testsToRun[@]}"; do
      ! _hasItem $funcTestFor ${alreadyListed[@]} || continue
      dstVarName=selfTests srcVarName=${funcTestFor}Tests codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      local count="${#selfTests[@]}"
      echo "${_ansiCmd}$funcTestFor ${_ansiPrimaryLiteral}$count${_ansiReset} $(_getPluralWord $count $_testPlurals)"
      alreadyListed+=( $funcTestFor )
      totalCount=$(( totalCount + count ))
    done
    local count="${#alreadyListed[@]}"
    echo "Всего ${_ansiPrimaryLiteral}$count${_ansiReset} $(_getPluralWord $count $_funcPlurals) и ${_ansiPrimaryLiteral}$totalCount${_ansiReset} $(_getPluralWord $totalCount $_testPlurals)"
  elif [[ -n $wasError ]]; then
    return 1
  else
    local -a _failedTests=()
    local -a _succeedTests=()
    local timeStart=$(date +%s)
    local arg; for arg in "${testsToRun[@]}"; do
      local funcTestFor="${arg% *}"
      dstVarName=selfTests srcVarName=${funcTestFor}Tests codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      local count="${#selfTests[@]}"
      local maxIndex=$(( $count - 1 ))
      local num="${arg##* }"
        if [[ $num == $funcTestFor ]]; then
          num=
        elif [[ $num -lt 0 ]]; then
          num=$(( num + count ))
        fi
      if [[ -n $num ]]; then
        _runBashTestWrapper
      else
        for ((num=0; num<$count; num++)); do
          _runBashTestWrapper
        done
      fi
    done
    local timeEnd=$(date +%s)
    local timeElapsed=$(( timeEnd - timeStart ))
    local failedCount=${#_failedTests[@]}
    local succeedCount=${#_succeedTests[@]}
    local totalCount=$(( failedCount + succeedCount ))
    if [[ $totalCount -eq 0 ]]; then
      _warn "Нет тестов для прогона"
    elif [[ $failedCount -eq 0 ]]; then
      _ok "Все тесты (${_ansiPrimaryLiteral}$succeedCount${_ansiOK}) пройдены успешно"
    elif [[ $failedCount -eq 1 ]]; then
      _err "Завершился ошибкой 1 тест из $totalCount: ${_failedTests[0]}"
    elif [[ $succeedCount -eq 0 ]]; then
      _err "Все тесты (${_ansiPrimaryLiteral}${failedCount}${_ansiErr}) завершились ошибкой"
    else
      _err "Завершились ошибкой $failedCount $(_getPluralWord $totalCount $_testPlurals) из $totalCount:"
      local -a funcsTestFor=()
      local arg; for arg in "${_failedTests[@]}"; do
        local funcTestFor="${arg% *}"
        local holder=__nums_$funcTestFor
        if ! _hasItem $funcTestFor ${funcsTestFor[@]}; then
          funcsTestFor+=( $funcTestFor )
          local $holder=
        fi
        local num="${arg##* }"
        local $holder+=" $num"
      done
      local funcTestFor; for funcTestFor in ${funcsTestFor[@]}; do
        local holder=__nums_$funcTestFor
        local -a nums=( $( printf "%s\n" ${!holder} | sort -n -u )  )
        local err="  $funcTestFor"
        local diapMin=${nums[0]}
        local diapMax=${nums[0]}
        local num; for num in ${nums[@]}; do
          [[ $num -gt $diapMax ]] || continue
          if [[ $num -eq $(( diapMax + 1 )) ]]; then
            diapMax=$num
          else
            [[ $diapMax -eq $diapMin ]] \
              && err+=" $diapMin" \
              || err+=" $diapMin..$diapMax"
            diapMin=$num
            diapMax=$num
          fi
        done
        [[ $diapMax -eq $diapMin ]] \
          && err+=" $diapMin" \
          || err+=" $diapMin..$diapMax"
        _err "$err"
      done
    fi
    local returnCode=$?
    [[ -n $noTiming ]] \
      || echo "Время выполнения: $timeElapsed $(_getPluralWord $timeElapsed секунда секунды секунд)"
    if [[ $returnCode -eq 0 ]]; then
      if [[ -n $_isBwDevelop ]]; then
        [[ ${#args[@]} -gt 0 ]] || _buildBw || returnCode=$?
      else
        local fileNameToProcess; for fileNameToProcess in "${fileNamesToProcess[@]}"; do
          _unsetBash "$testsDirSpec/$fileNameToProcess"
        done
        # local funcName; for funcName in ${funcNamesToPregen[@]}; do
        #   _fileSpec="$_bwDir/$_generatedDir/$funcName.completion.code$_unsetFileExt" codeHolder=_codeSourceIf eval "$_evalCode"
        # done
        # for _fileSpec in "$_bwDir/tests/$_generatedDir"/*$_unsetFileExt; do
        #   codeHolder=_codeSource eval "$_evalCode"
        # done
      fi
    fi
    return $returnCode
  fi
}
_codeToParseDiap='
  if [[ ${!diapHolder} =~ ^(-?[[:digit:]]+)(..(-?[[:digit:]]+)?)?$ ]]; then
    diapMin=${BASH_REMATCH[1]}
    if [[ -z ${BASH_REMATCH[2]} ]]; then
      diapMax=diapMin
    elif [[ -z ${BASH_REMATCH[3]} ]]; then
      diapMax=Infinite
    else
      diapMax=${BASH_REMATCH[3]}
    fi
  elif [[ ${!diapHolder} =~ ^..(-?[[:digit:]]+)?$ ]]; then
    diapMin=Infinite
    if [[ -z ${BASH_REMATCH[1]} ]]; then
      diapMax=Infinite
    else
      diapMax=${BASH_REMATCH[1]}
    fi
  fi
  [[ -z $diapMin || $diapMin == Infinite || -z $diapMax || $diapMax == Infinite || $diapMin -le $diapMax ]] \
    || diapErr="ожидает, что левая граница ${_ansiPrimaryLiteral}$diapMin${_ansiErr} диапазона будет не больше правой ${_ansiPrimaryLiteral}$diapMax"
'

_accumTime=
_startTime=
_runBashTestWrapper() {
  local testName="$funcTestFor"
  testName+=" $num"
  eval local -a params=( "${selfTests[$num]}" )

  _runBashTest "${params[@]}" "$testName" "$num" \
    && _succeedTests+=( "$testName" ) \
    || _failedTests+=( "$testName" )
}

# =============================================================================

_stdStreams=( 'stdout' 'stderr' )
_runBashTestOptionsMaxVarNumber=3
_runBashTestParams=()
_runBashTestParams() {
  varName=_runBashTestParams codeHolder=_codeToUseCache eval "$_evalCode"
  _runBashTestParams=(
    '@--before/b=' '@--immediatly/i=' '@--after/a='
    '--return/r=0'
  )
  local -a tstEntities=()
  local stdStream; for stdStream in ${_stdStreams[@]}; do
    _runBashTestParams+=( "--$stdStream:?=" )
    tstEntities+=( $stdStream )
  done
  local varPrefix i=0; while [[ i -lt _runBashTestOptionsMaxVarNumber ]]; do 
    i=$(( i + 1 )) 
    varPrefix=var 
    [[ $i -gt 1 ]] && varPrefix+="$i"
    _runBashTestParams+=( "--${varPrefix}Name=" "--${varPrefix}Value=" )
    tstEntities+=( $varPrefix )
  done
  local -r fileSpecBase="/tmp/$(basename "${BASH_SOURCE[0]}")"
  local tstEntity; for tstEntity in ${tstEntities[@]}; do
    _runBashTestParams+=( "--${tstEntity}TstFileSpec=$fileSpecBase.${tstEntity}" )
    for param in 'Eval:1..:?' 'DiffOptions=' 'EchoOptions=' 'ParaWithIndent:0..' 'ParaWithIndentBase:1..=4'; do
      _runBashTestParams+=( "--${tstEntity}${param}" )
    done
    _runBashTestParams+=( "@--${tstEntity}EtaPreProcess" )
    _runBashTestParams+=( "@--${tstEntity}TstPostProcess" )
  done
  _runBashTestParams+=( --catOptions "")
  _runBashTestParams+=( "--inTmpDir")
  _runBashTestParams+=( "testCmd" "testName" "testNum:0..=0")
  _saveToCache '_runBashTestParams'
}

_bwTmpDir="$_bwDir/tmp"
_runBashTest() { eval "$_funcParams2"
  local -a tstEntities=()
  local stdStream; for stdStream in ${_stdStreams[@]}; do
    local ${stdStream}Value="${!stdStream}"
    local ${stdStream}Name="$stdStream"
    tstEntities+=( "$stdStream" )
  done
  local i; for ((i=0; i<$_runBashTestOptionsMaxVarNumber; i++)); do
    local varPrefix=var; [[ $i -gt 0 ]] && varPrefix+="$i"
    local varNameHolder="${varPrefix}Name"
    local varValueHolder="${varPrefix}Value"
    [[ -z ${!varValueHolder} ]] || [[ -n ${!varNameHolder} ]] \
      || return $(_err --showStack 2 "${_ansiCmd}${FUNCNAME[0]}${_ansiErr} ожидает, что опция ${_ansiCmd}--$varValueHolder${_ansiErr} будет задана вместе с опцией ${_ansiCmd}--$varNameHolder")
  done

  local hasDifference returnDiffers
  local tstEntity; for tstEntity in "${tstEntities[@]}"; do
    local ${tstEntity}Differs=
  done
  local returnTst
  [[ -z $inTmpDir ]] \
    && _runBashTestHelper \
    || _inDir -v none -t --noCleanOnFail "$_bwTmpDir" _runBashTestHelper
  local returnCode=$?

  # remove extra backslashes from $testCmd
  local testTitle; printf -v testTitle "$testName${_ansiReset}: ${_ansiCmd}%s" "$testCmd"
  printf -v testTitle "$testTitle"

  if [[ $returnCode -ne 0 ]]; then
    _err $testTitle
    return $returnCode
  fi

  if [[ -z $hasDifference ]]; then
    _ok $testTitle
    return 0
  else
    _err $testTitle
    [[ -n $returnDiffers ]] &&
      echo -e "${_ansiOutline}return${_ansiErr} is expected to be ${_ansiOK}$return${_ansiErr}, but got ${_ansiWarn}$returnTst${_ansiReset}" >&2
    local tstEntity; for tstEntity in ${tstEntities[@]}; do
      local differsHolder="${tstEntity}Differs"
      if [[ -n ${!differsHolder} ]]; then
        local diffOptionsHolder="${tstEntity}DiffOptions"
        local diffOptions; [[ -n ${!diffOptionsHolder} ]] && diffOptions="${!diffOptionsHolder} "
        local nameHolder="${tstEntity}Name"
        local tstFileSpecHolder="${tstEntity}TstFileSpec"
        echo -e "${_ansiOutline}${!nameHolder}${_ansiReset} (${_ansiCmd}diff $diffOptions${!tstFileSpecHolder} ${!tstFileSpecHolder}.eta${_ansiReset}):" >&2
        cat $catOptions "${!tstFileSpecHolder}.diff" >&2
      fi
    done
    return 2
  fi
}

_codeToEvalItem='
  if [[ ${item:0:1} == "!" ]]; then
    eval "${item:1}"
  else
    eval $item
  fi
'

_runBashTestHelper() {
  local item; for item in "${before[@]}"; do
    codeHolder=_codeToEvalItem eval "$_evalCode" || return $?
  done

  _rm "$stdoutTstFileSpec" "$stderrTstFileSpec"
  _profileInitTransfer
  ( profileTmpFileSpec="$_profileTmpFileSpec"
    eval "$testCmd" 1>"$stdoutTstFileSpec" 2>"$stderrTstFileSpec";
    local returnTst=$?

    local i; for ((i=0; i<_runBashTestOptionsMaxVarNumber; i++)); do
      local varPrefix=var; [[ $i -gt 0 ]] && varPrefix+="$i"
      local varNameHolder="${varPrefix}Name"
      if [[ -n ${!varNameHolder} ]]; then
        local varTstFileSpecHolder="${varPrefix}TstFileSpec"
        declare -p "${!varNameHolder}" >"${!varTstFileSpecHolder}"
      fi
    done

    local -a tstEntities=()
    local stdStream; for stdStream in ${_stdStreams[@]}; do
      tstEntities+=( $stdStream )
    done
    local varPrefix i=0; while [[ i -lt _runBashTestOptionsMaxVarNumber ]]; do 
      i=$(( i + 1 )) 
      varPrefix=var         
      [[ $i -gt 1 ]] && varPrefix+="$i"
      tstEntities+=( $varPrefix )
    done
    local tstEntity; for tstEntity in ${tstEntities[@]}; do
      local tstFileSpecHolder="${tstEntity}TstFileSpec"
      local tstPostProcessHolder="${tstEntity}TstPostProcess"
      dstVarName=tstPostProcess srcVarName=$tstPostProcessHolder codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
      if [[ ${#tstPostProcess[@]} -gt 0 ]]; then
        local tstValue="$(<"${!tstFileSpecHolder}")"
        local item; for item in "${tstPostProcess[@]}"; do
          eval local -a subitems=\( $item \)
          tstValue=$(echo "$tstValue" | "${subitems[@]}" )
        done
        echo "$tstValue">"${!tstFileSpecHolder}"
      fi
    done

    _profileDoTransfer
    exit $returnTst
  ); returnTst=$?
  _profileGetTransfer

  local i; for ((i=0; i<_runBashTestOptionsMaxVarNumber; i++)); do
    local varPrefix=var; [[ $i -gt 0 ]] && varPrefix+="$i"
    local varNameHolder="${varPrefix}Name"
    [[ -z ${!varNameHolder} ]] || tstEntities+=( $varPrefix )
  done

  local item; for item in "${immediatly[@]}"; do
    codeHolder=_codeToEvalItem eval "$_evalCode" || return $?
  done

  if [[ $returnTst -ne $return ]]; then
    returnDiffers=true
    hasDifference=true
  fi
  local tstEntity; for tstEntity in "${tstEntities[@]}"; do
    local valueHolder="${tstEntity}Value"
    local evalHolder="${tstEntity}Eval"
    local paraWithIndentHolder="${tstEntity}ParaWithIndent"
    local paraWithIndentBaseHolder="${tstEntity}ParaWithIndentBase"
    local etaPreProcessHolder="${tstEntity}EtaPreProcess"
    local tstFileSpecHolder="${tstEntity}TstFileSpec"
    local diffOptionsHolder="${tstEntity}DiffOptions"
    local echoOptionsHolder="${tstEntity}EchoOptions"
    local etaValue="${!valueHolder}"
    if [[ -n ${!evalHolder} ]]; then
      local evalCount=${!evalHolder}; [[ $evalCount =~ ^[[:digit:]]+$ ]] || evalCount=1
      for ((i=0; i<evalCount; i++)); do
        eval local etaValue="$etaValue"
      done
    fi
    if [[ -n ${!paraWithIndentHolder} ]]; then
      local -a OPT_base=( --indentBase ${!paraWithIndentBaseHolder} )
      etaValue=$(_getArrangedMultilineDescription "${OPT_base[@]}" ${!paraWithIndentHolder} "$etaValue") \
        || return $?
    fi
    dstVarName=etaPreProcess srcVarName=$etaPreProcessHolder codeHolder=_codeToInitLocalCopyOfArray eval "$_evalCode"
    local item; for item in "${etaPreProcess[@]}"; do
      eval local -a subitems=\( $item \)
      etaValue=$(echo "$etaValue" | "${subitems[@]}" )
    done
    [[ -z $etaValue ]] \
      && echo -n >"${!tstFileSpecHolder}.eta" \
      || echo ${!echoOptionsHolder} "$etaValue" >"${!tstFileSpecHolder}.eta"
    $(diff -a ${!diffOptionsHolder} "${!tstFileSpecHolder}" "${!tstFileSpecHolder}.eta" >"${!tstFileSpecHolder}.diff" 2>&1)
    if [[ $? -ne 0 ]]; then
      eval ${tstEntity}Differs=true
      hasDifference=true
    fi
  done

  local item; for item in "${after[@]}"; do
    codeHolder=_codeToEvalItem eval "$_evalCode" || return $?
  done
}

_getArrangedMultilineDescriptionParams=(
  '--indentBase:1..=4'
  'indentLevel:0..'
  'description'
)
_getArrangedMultilineDescription() { eval "$_funcParams2"
  local perlCode='
    use POSIX;
    use utf8;
    use List::Util qw[max];
    $baseIndentLevel=0;
    sub initFirstIndent {
      my $spaces=shift;
      $baseIndentLevel = floor(length($spaces) / 2);
      return "";
    }
    sub replaceTo {
      my $spaces=shift;
      my $noNewLine=0;
      if ( $spaces eq "--noNewLine" ) { $spaces=""; $noNewLine=1; }
      my $indentLevel='$indentLevel' + max(0, floor(length($spaces) / 2) - $baseIndentLevel);
      my $indent=" " x ( '$indentBase' * $indentLevel );
      $noNewLine ? "$indent" : "\n$indent";
    }
    $/=undef;
    $_=<>;
    s/^(?:[ ]*\n)*(\s*)/initFirstIndent($1)/e;
    s/\s+$//;
    s/\n(\s*)/@{[replaceTo($1)]}/g;
    s/^/@{[replaceTo("--noNewLine")]}/;
    print;
  '
  echo "${description}" | perl -C -e "$perlCode"
}

# =============================================================================

_postProcessDeclareArray='"perl -e \"
    undef \\\$/; 
    \\\$_=<STDIN>; 
    if ( s/^declare -a\S* ([\w][\w\d]*)='\''\\(/declare -a \\\$1=(/s ) {
      s/'\''\\\$//s;
    }
    sub replaceTo {
      my \\\$q=shift;
      my \\\$val=shift;
      if (\\\$val =~ s/\n/\\\\\\\\n/g) {
        \\\"\\\\$'\''\\\" . \\\$val . \\\"'\''\\\"
      } else {
        \\\$q . \\\$val . \\\$q
      }
    }
    s/(\\\")([^\\\"]+)\\\"/replaceTo(\\\$1,\\\$2)/sge; 
    print;
  \""'

# =============================================================================
