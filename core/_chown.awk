#!/usr/bin/awk

# =============================================================================

/^========$/ {
  isFile = 1
  next
}

# =============================================================================

{ 
  if (isFile) {
    single[single_qt++]=$0
  } else {
    relative = substr($0, length(root) + 1); 
    depth = split(relative, arr, "/")
    if (depth > 0) depth -= 1
    # print $0 ", " "root: " root ", " "depth: " depth ", " "maxdepth: " maxdepth > "/dev/tty"
    if (depth < maxdepth) {
      single[single_qt++]=$0
    } else {
      recursive[recursive_qt++]=$0
    }
  }
}

# =============================================================================

END {
  for (i = 0; i < recursive_qt; i++) {
    fileSpec=recursive[i]
    gsub("'", "'\\''", fileSpec)
    lines[lines_qt++] = sprintf("sudo chown -R %s '%s'", user, fileSpec)
  }

  if (maxFilesPerLine < 1) maxFilesPerLine = 1
  filesInLine=0
  isFilledLine = 0
  line = ""
  for (i = 0; i < single_qt; i++) {
    fileSpec=single[i]
    gsub("'", "'\\''", fileSpec)
    line = line " '" fileSpec "'"
    if (++filesInLine >= maxFilesPerLine) {
      lineFilled()
    } else {
      isFilledLine = 0
    }
  }
  if (!isFilledLine) {
    lineFilled()
  }

  if (maxProcesses < 1) maxProcesses = 1
  maxLinesPerProcess = lines_qt / maxProcesses
  if (verbose == "true") {
    print "maxdepth: " maxdepth ", " "recursive_qt: " recursive_qt ", " "single_qt: " single_qt ", " "maxFilesPerLine: " maxFilesPerLine ", " "lines_qt: " lines_qt ", " "maxProcesses: " maxProcesses ", " "maxLinesPerProcess: " maxLinesPerProcess > "/dev/stderr"
  }
  shuf(lines, sh_lines)
  if (maxLinesPerProcess < 1) maxLinesPerProcess = 1

  process_qt = 0
  linesInProcess = 0
  isFilledProcess = 0
  for (i = 0; i < lines_qt; i++) {
    if (linesInProcess == 0) {
      print "("
      print "timeStart=$(date +%s)"
      isFilledProcess = 0
    }
    print sh_lines[i]
    if (++linesInProcess >= maxLinesPerProcess) {
      processFilled()
    }
  }
  if (!isFilledProcess) {
    processFilled()
  }

  for (i = 0; i < process_qt; i++) {
    print "wait $_chownPid" i
  }
}

# =============================================================================
# =============================================================================

function lineFilled() {
  lines[lines_qt++] = sprintf("sudo chown %s%s", user, line)
  filesInLine = 0
  isFilledLine = 1
  line = ""
}

function processFilled() {
  print "timeEnd=$(date +%s)"
  print "timeElapsed=$(( timeEnd - timeStart ))"
  print "printf 'process #%s: %ss\\n' " process_qt " $timeElapsed"
  print ") &"
  print "_chownPid" process_qt++ "=$!"

  isFilledProcess = 1
  linesInProcess = 0
}

# =============================================================================
# =============================================================================
# https://github.com/e36freak/awk-libs/blob/master/shuf.awk

# actual shuffle function
# shuffles the values in "array" in-place, from indices "left" to "right".
# required for all of the shuf() functions below
function __shuffle(array, left, right,    r, i, tmp) {
  # loop backwards over the elements
  for (i=right; i>left; i--) {
    # generate a random number between the start and current element
    r = int(rand() * (i - left + 1)) + left;

    # swap current element and randomly generated one
    tmp = array[i];
    array[i] = array[r];
    array[r] = tmp;
  }
}

## usage: shuf(s, d)
## shuffles the array "s", creating a new shuffled array "d" indexed with
## sequential integers starting with 0. returns the length, or -1 if an error
## occurs. leaves the indices of the source array "s" unchanged. uses the knuth-
## fisher-yates algorithm. requires the __shuffle() function.
function shuf(array, out, count, i) {
  # loop over each index, and generate a new array with the same values and
  # sequential indices
  count = 0;
  for (i in array) {
    out[count++] = array[i];
  }

  # seed the random number generator
  srand();

  # actually shuffle
  __shuffle(out, 1, count);

  # return the length
  return count;
}
