#!/usr/bin/awk

# expects following vars ( -v ):
#   root
#   maxProcesses: 2..=2
#   maxLineLength: 1000...=115000 # 115000 - эмпирическая величина, на ней Ubuntu14 еще не "валится", а на 120000 уже
#   verbose: "true" or other
#   user: string of username

# example:
  # awk -f ~/bw/core/_chown.awk -v root=/home/dev -v user=dev -v verbose=true some.list
# stdout:
  # #!/bin/bash
  # pushd "/home/dev/"
  # pwd
  # ( # process #0
  # timeStart=$(date +%s)
  # sudo chown dev 'proj/one' 'proj/two' 'proj/three'
  # # files in above line: 3
  # timeEnd=$(date +%s)
  # timeElapsed=$(( timeEnd - timeStart ))
  # printf 'process #%s: filesInProcess: %s, linesInProcess: %s, timeElapsed: %ss\n' 0  3 1 $timeElapsed
  # ) &
  # _chownPid0=$!
  # ( # process #1
  # timeStart=$(date +%s)
  # sudo chown dev 'proj/four' 'proj/five' 'proj/six' 'proj/seven'
  # # files in above line: 4
  # timeEnd=$(date +%s)
  # timeElapsed=$(( timeEnd - timeStart ))
  # printf 'process #%s: filesInProcess: %s, linesInProcess: %s, timeElapsed: %ss\n' 1  4 1 $timeElapsed
  # ) &
  # _chownPid1=$!
  # wait $_chownPid0
  # wait $_chownPid1
  # popd
# stderr:
  # root: /home/dev/, rootLen: 10
  # fileSpecs_qt: 7, maxProcesses: 2, maxLineLength: 115000
  # lines_qt: 2, foundMaxLineLength: 62, foundMaxFilesInLine: 4
  # process_qt: 2, maxFilesPerProcess: 3,5, foundMaxFilesInProcess: 4, foundMaxLinesInProcess: 1
# where some.list:
  # /home/dev/proj/one
  # /home/dev/proj/two
  # /home/dev/proj/three
  # /home/dev/proj/four
  # /home/dev/proj/five
  # /home/dev/proj/six
  # /home/dev/proj/seven

# =============================================================================

BEGIN {
  sub(/\/+$/, "", root)
  root = root "/"
  rootLen = length(root)
  if (length(user) == 0) {
    print "`-v user=<some username>` should be specified" > "/dev/tty"
    exitCode = 1
    exit exitCode
  }
  linePrefix = sprintf("sudo chown %s", user)
  if (maxLineLength <= length(linePrefix)) { maxLineLength = 115000 }
  if (maxProcesses < 2) maxProcesses = 2

  nettoMaxLineLength = maxLineLength - length(linePrefix) 
  fileSpecs_qt = 0
  foundMaxLineLength = 0
  foundMaxFilesInLine = 0
  foundMaxLinesPerProcess = 0
  foundMaxFilesInProcess = 0
}

{ 
  if (substr($0, 1, rootLen) == root) {
    $0 = substr($0, rootLen + 1)
  }
  if (length($0) > 0 ) fileSpecs[fileSpecs_qt++]=$0
}

# =============================================================================

END {
  if (exitCode != 0) {
    exit exitCode
  } else {
    maxFilesPerProcess = fileSpecs_qt / maxProcesses
    if (maxFilesPerProcess < 1) maxFilesPerProcess = 1
    maxFilePerProcessThreshold = maxFilesPerProcess - 0.5

    process_qt = 0
    filesInProcess = 0
    linesInProcess = 0

    isFilledLine = 0
    line = ""
    lineSuffix = ""
    filesInLine = 0
    for (i = 0; i < fileSpecs_qt; i++) {
      fileSpec=fileSpecs[i]
      gsub("'", "'\\''", fileSpec)
      lineSuffix = " '" fileSpec "'"
      # sub(/^ '\/home\/dev\//, " ~/'", lineSuffix)
      isLineFilled = length(line) + length(lineSuffix) > nettoMaxLineLength
      isProcessFilled = filesInProcess >= maxFilePerProcessThreshold && process_qt < maxProcesses - 1 
      if (filesInLine > 0 && (isLineFilled || isProcessFilled )) {
        lineFilled()
        if (isProcessFilled) {
          processFilled()
        }
      } 
      line = line lineSuffix
      filesInLine++
      filesInProcess++
    }
    lineFilled()
    processFilled()

    for (i = 0; i < process_qt; i++) {
      print "wait $_chownPid" i
    }
    print "popd >/dev/null 2>&1"

    if (verbose == "true") {
      print "root: " root ", " "rootLen: " rootLen > "/dev/stderr"
      print "fileSpecs_qt: " fileSpecs_qt ", " "maxProcesses: " maxProcesses ", " "maxLineLength: " maxLineLength > "/dev/stderr"
      print "lines_qt: " lines_qt ", " "foundMaxLineLength: " foundMaxLineLength ", " "foundMaxFilesInLine: " foundMaxFilesInLine > "/dev/stderr"
      print "process_qt: " process_qt ", " "maxFilesPerProcess: " maxFilesPerProcess ", " "foundMaxFilesInProcess: " foundMaxFilesInProcess  ", " "foundMaxLinesInProcess: " foundMaxLinesInProcess > "/dev/stderr"
    }
  }
}

# =============================================================================
# =============================================================================

function lineFilled() {
  if (linesInProcess == 0) {
    if (process_qt == 0) {
      print "#!/bin/bash"
      print "pushd \"" root "\" >/dev/null 2>&1"
    }
    header = "("
    if (verbose == "true") header = header " # process #" process_qt 
    print header
    if (verbose == "true") {
      print "timeStart=$(date +%s)"
    }
  }
  line = linePrefix line
  print line
  lines_qt++
  linesInProcess++
  if (verbose == "true") {
    print "# files in above line: " filesInLine
  }
  isFilledLine = 1
  if (length(line) > foundMaxLineLength) foundMaxLineLength = length(line)
  if (filesInLine > foundMaxFilesInLine) foundMaxFilesInLine = filesInLine
  line = ""
  filesInLine = 0
}

function processFilled() {
  if (verbose == "true") {
    print "timeEnd=$(date +%s)"
    print "timeElapsed=$(( timeEnd - timeStart ))"
    print "printf 'process #%s: filesInProcess: %s, linesInProcess: %s, timeElapsed: %ss\\n'" " " process_qt " " " " filesInProcess " " linesInProcess " " "$timeElapsed"
  }
  print ") &"
  print "_chownPid" process_qt++ "=$!"

  if (filesInProcess > foundMaxFilesInProcess) foundMaxFilesInProcess = filesInProcess
  if (linesInProcess > foundMaxLinesInProcess) foundMaxLinesInProcess = linesInProcess
  filesInProcess = 0
  linesInProcess = 0
}
