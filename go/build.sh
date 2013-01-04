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
	rm -f gk/bin/*

	go install gk/gkerr
	go install gk/gklog
	go install gk/gktmpl
	go install pq/pq
	go install gk/sec
	go install gk/wf
	go install gk/login
	go install gk/database

	go test gk/sec

	go build -o gk/bin/wfToJsMain gk/src/gk/main/wfToJsMain.go
	go build -o gk/bin/loginServerMain gk/src/gk/main/loginServerMain.go

fi

