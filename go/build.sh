#!/usr/bin/env bash

ok='y'

if [ -z "$GOPATH" ]
then
        echo "missing GOPATH"
        ok='n'
fi

if [ -z "$GOBASE" ]
then
        echo "missing GOBASE"
        ok='n'
fi

if [ -z "$GOROOT" ]
then
        echo "missing GOROOT"
        ok='n'
fi

if [ $ok == "y" ]
then
        if [ ! -d $GOBASE/gk/bin ]
        then
                echo "creating missing gk bin directory"
                mkdir $GOBASE/gk/bin
        fi
        if [ ! -d $GOBASE/gk/pkg ]
        then
                echo "creating missing gk pkg directory"
                mkdir $GOBASE/gk/pkg
                mkdir $GOBASE/gk/pkg/linux_amd64
        fi

        if [ ! -d $GOBASE/pq/pkg ]
        then
                echo "creating missing gk pkg directory"
                mkdir $GOBASE/gk/pkg
                mkdir $GOBASE/gk/pkg/linux_amd64
        fi

	go clean all
	rm -f */pkg/*/*.a
	rm -f nw/bin/*

	go install pq/pq
	go install gk/sec
	go test gk/sec

	#go build -o nw/bin/xMain gk/src/main/xMain.go

	#go test gk
fi

