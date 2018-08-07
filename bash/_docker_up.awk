BEGIN { 
  port=0
} 
END { 
  if (port > 0) { 
    print port 
  } else { 
    exit 1 
  } 
} 
{ 
  if (match($0, /:[0-9]+ failed: port is already allocated/)) { 
    $0=substr($0, RSTART)
    match($0, " ")
    port=substr($0, 2, RSTART - 2)
    exit 
  }
} 