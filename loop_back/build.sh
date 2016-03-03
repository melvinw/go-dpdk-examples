#!/bin/bash
if [ -z "$RTE_SDK" ]
then
	echo "RTE_SDK undefined"
    exit 1
fi

if [ -z "$RTE_TARGET" ]
then
    RTE_TARGET=x86_64-native-linuxapp-gcc
fi

OLD_LDFLAGS=$CGO_LDFLAGS
OLD_CFLAGS=$CGO_CFLAGS

export CGO_LDFLAGS="-L${RTE_SDK}/${RTE_TARGET}/lib -Wl,--whole-archive -ldpdk \
    -lz -Wl,--start-group -lrt -lm -ldl -Wl,--end-group -Wl,--no-whole-archive"
export CGO_CFLAGS="-m64 -pthread -Ofast -march=native \
    -I${RTE_SDK}/${RTE_TARGET}/include"

go get github.com/melvinw/go-dpdk

export CGO_LDFLAGS=$OLD_LDFLAGS
export CGO_CFLAGS=$OLD_CFLAGS

go build main.go
