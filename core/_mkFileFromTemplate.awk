#/bin/awk
BEGIN { 
  exitCode = 0 
}
END { 
  if ( exitCode != 0 ) {
    exit exitCode
  } else {
    for (key in varNames) print key
  }
}
{
  pos = 1 
  while (match($0, /\${/ )) {
    $0 = substr($0, RSTART + RLENGTH)
    pos += RSTART + RLENGTH - 1
    if ( match($0, /}/) ) {
      varName = substr($0, 1, RSTART - 1)
      $0 = substr($0, RSTART + 1)
      pos += RSTART
      if (length(varName)) varNames[varName] += 1
    } else {
      print "ERR: " funcName " expects } somewhere at " FILENAME " line " NR " after pos " pos > "/dev/tty"
      exitCode=1
      exit exitCode
    }
  }
}