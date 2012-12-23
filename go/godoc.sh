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
	godoc -http=:6060 -goroot=$GOBASE -templates=$GOBASE/godoc_templates
fi

