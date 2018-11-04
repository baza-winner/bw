#/bin/sh
find . -name '*_test.go' ! -path '*defparse*' ! -path '*defvalid*' ! -path '*pfa*' | perl -ne '@pp = split "/", $_; shift @pp; pop @pp; print join "/", @pp; print "\n"' | xargs -n 1 sh -c 'go test -v github.com/baza-winner/bwcore/$0 || exit 255'
# find . -name '*_test.go' | perl -ne '@pp = split "/", $_; shift @pp; pop @pp; print join "/", @pp; print "\n"' | xargs -n 1 sh -c 'go test -v github.com/baza-winner/bwcore/$0 || exit 255'
