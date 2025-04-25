#!/bin/sh

go install github.com/vektra/mockery/v2@v2.53

interfaces="
  ProductRepository
  CacheService
  TokenProvider
"

set -- $interfaces

while [ -n "$1" ]; do
    echo "[+] generate $1 mock..."
    $GOPATH/bin/mockery --name="$1" --recursive --output=tests/mocks --filename="${1}Mock.go" --structname="${1}Mock"
    shift
done
