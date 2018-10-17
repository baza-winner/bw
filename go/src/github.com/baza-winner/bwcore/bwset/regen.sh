#!/bin/sh
bwsetterPackagePath=github.com/baza-winner/bwcore/bwsetter
bwsetPackagePath=github.com/baza-winner/bwcore/bwset
go install $bwsetterPackagePath && rm -f "$GOPATH/src/$bwsetPackagePatch/*_set*.go" && go generate $bwsetPackagePath && go test -v $bwsetPackagePath 2>&1 | pp
