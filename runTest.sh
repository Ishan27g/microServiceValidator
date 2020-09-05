#!/bin/sh

BINARY=bin/test
PS=dockerresult.txt
rm bin/test
rm $PS
make build
if test -f "$BINARY"; then
    cd bin/
    chmod 777 test
    cd ..
    ./bin/test
fi

