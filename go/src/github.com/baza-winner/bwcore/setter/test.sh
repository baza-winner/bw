#!/bin/sh
_test_setter() {
  while true; do
    go install github.com/baza-winner/bwcore/setter || { returnCode=?; break; }
    rm -f $GOPATH/src/github.com/baza-winner/bwcore/setter/example/*_set.go || { returnCode=?; break; }
    go generate github.com/baza-winner/bwcore/setter/example || { returnCode=?; break; }
    go test -v github.com/baza-winner/bwcore/setter/example || { returnCode=?; break; }
    break
  done
  return $returnCode
}
_test_setter
