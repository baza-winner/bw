#/bin/sh
find . -depth 1 -type d -exec basename {} \; | xargs -n 1 sh -c 'go test -v github.com/baza-winner/bwcore/$0 || exit 255'