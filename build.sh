#!/usr/bin/env bash
go mod tidy
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o otteralert .
strip --strip-unneeded otteralert
upx --lzma otteralert