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
                mkdir $GOBASE/pq/pkg
                mkdir $GOBASE/pq/pkg/linux_amd64
        fi

        if [ ! -d $GOBASE/gonet/pkg ]
        then
                echo "creating missing gk pkg directory"
                mkdir $GOBASE/gonet/pkg
                mkdir $GOBASE/gonet/pkg/linux_amd64
        fi

	go clean all
	find */pkg -type f -print | grep "\.a$" | xargs rm
	rm -f gk/bin/*

fi

