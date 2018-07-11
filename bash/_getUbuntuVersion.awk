#/bin/awk
/^Release/ {
	majorVer = match($2, /[0-9]+/) ? substr($2, RSTART, RLENGTH) : 0
  minorVer = match($2, /\.[0-9]+/) ? substr($2, RSTART + 1, RLENGTH - 1) : 0
  printf("%d%03d\n", majorVer, minorVer)
}