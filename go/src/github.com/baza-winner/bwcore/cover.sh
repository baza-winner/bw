#!/bin/sh
go test -coverprofile=coverage.out $1 && go tool cover -html=coverage.out
