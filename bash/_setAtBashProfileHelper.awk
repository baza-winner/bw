#/bin/awk

BEGIN { 
  foundExact = 0
  foundMatch = 0 
}

{ 
  if ($0 == exactLine) {
    foundExact++
    printExactLine()
  } else if ( matchRegexp != "-" && $0 ~ matchRegexp ) {
    foundMatch++ 
    printExactLine()
  } else {
    print $0 
  }
}

END { 
  if (length(uninstall) == 0 && foundExact == 1 && foundMatch == 0) 
    exitCode = 0 
  else 
    exitCode = foundExact + foundMatch == 0 ? 1 : 2
  exit exitCode
}

function printExactLine() {
  if (length(uninstall) == 0 && ( foundExact + foundMatch == 1 )) print exactLine
}

