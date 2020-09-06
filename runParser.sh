#!/bin/sh

BINARY=parser
PS=dockerresult.txt

#cleanup
rm bin/$BINARY
rm $PS

#build
make build

#run
if test -f "bin/$BINARY"; then
    cd bin/
    chmod 777 $BINARY
    cd ..
    ./bin/$BINARY $1
fi

