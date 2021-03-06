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

	go install pq/pq
	go install code.google.com/p/go.net/dict
	go install code.google.com/p/go.net/html
	go install code.google.com/p/go.net/html/atom
	go install code.google.com/p/go.net/idna
	go install code.google.com/p/go.net/ipv4
	go install code.google.com/p/go.net/proxy
	go install code.google.com/p/go.net/spdy
	go install code.google.com/p/go.net/websocket
	go install code.google.com/p/go.net/publicsuffix
	go install gk/gkerr
	go install gk/gklog
	go install gk/gktmpl
	go install gk/gkjson
	go install gk/gksvg
	go install gk/sec
	go install gk/wf
	go install gk/login
	go install gk/game/persistence
	go install gk/game/field
	go install gk/game/iso
	go install gk/game/config
	go install gk/game/message
	go install gk/game/ses
	go install gk/game/ws
	go install gk/game
	go install gk/gknet
	go install gk/database
	go install gk/gkrand

	go test gk/gkerr
	go test gk/gklog
	go test gk/gktmpl
	go test gk/gkjson
	go test gk/gksvg
	go test gk/sec
	go test gk/wf
	go test gk/login
	go test gk/game/field
	go test gk/game/iso
	go test gk/game/config
	go test gk/game/message
	go test gk/game/ses
	go test gk/game/ws
	go test gk/game
	go test gk/gknet
	go test gk/database
	go test gk/gkrand

	#go build -o gk/bin/wfToJsMain gk/src/gk/main/wfToJsMain.go

	go build -o gk/bin/loginServerMain gk/src/gk/main/loginServerMain.go
	go build -o gk/bin/gameServerMain gk/src/gk/main/gameServerMain.go
	go build -o gk/bin/fixSvgMain gk/src/gk/main/fixSvgMain.go

fi

