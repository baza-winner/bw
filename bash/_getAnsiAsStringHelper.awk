#/bin/awk
BEGIN { 
  exitCode = 0 
}
END { 
  if ( exitCode != 0 ) {
    exit exitCode
  } else {
    for (key in ansi) print key "AsString='" ansi[key] "'"
    printf "%s", "varNames=( " 
    for (key in ansi) printf "%sAsString ", key 
    print ")"
  }
}
/^_ansi/ {
  pos = 0
  match($0, "=")
  key = substr($0, 1, RSTART - 1 )
  $0 = substr($0, RSTART + 1)
  pos += RSTART + 1
  if (substr($0, 1, 2) == "$'") {
    $0 = substr($0, 3)
    match($0, "'")
    value = substr($0, 1, RSTART - 1)
  } else if (substr($0, 1, 1) != "\"" ) {
    print "ERR: " funcName " expects $' or \" after " key "= at " FILENAME " line " NR > "/dev/tty"
    exitCode=1
    exit exitCode
  } else {
    $0 = substr($0, 2)
    pos += 1
    match($0, "\"")
    $0 = substr($0, 1, RSTART - 1)
    value = ""
    while (length($0)) {
      if ( substr($0, 1, 2) != "${" ) {
        print "ERR: " funcName " expects ${ at " FILENAME " line " NR " pos " pos > "/dev/tty"
        exitCode=1
        exit exitCode
      }
      $0 = substr($0, 3)
      pos += 2
      match($0, "}")
      varName = substr($0, 1, RSTART - 1)
      if ( ! length(ansi[varName]) ) {
        print "ERR: " funcName " got unknown " varName " at " FILENAME " line " NR " pos " pos > "/dev/tty"
        exitCode=1
        exit exitCode
      }
      value = value ansi[varName]
      $0 = substr($0, RSTART + 1)
    }
  }
  ansi[ key ] = value
}
