#/bin/awk

function printExactLine() {
  if (length(uninstall) == 0 && ( foundExact + foundMatch == 1 )) print exactLine
}

BEGIN { 
  foundExact = 0
  foundMatch = 0 
  # print "exactLine: ", exactLine > "/dev/tty"
  # print "matchRegexp: ", matchRegexp > "/dev/tty"
}

END { 
  if (length(uninstall) == 0 && foundExact == 1 && foundMatch == 0) 
    exitCode = 0 
  else 
    exitCode = foundExact + foundMatch == 0 ? 1 : 2
  # print "foundExact: ", foundExact > "/dev/tty"
  # print "foundMatch: ", foundMatch > "/dev/tty"
  # print "exitCode: ", exitCode > "/dev/tty"
  exit exitCode
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

# $0 == exactLine { 
#   foundExact++
#   printExactLine()
#   next 
# }

# $0 ~ matchRegexp  { 
#   foundMatch++ 
#   printExactLine()
#   next 
# }

# { print $0 }
