/*
	Copyright 2012-2013 1620469 Ontario Limited.

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package game

// documentation on go websocket: http://godoc.org/code.google.com/p/go.net/websocket
// getting go websocket: go get code.google.com/p/go.net/websocket

import (
	"fmt"
	"bytes"
	"io"
	"code.google.com/p/go.net/websocket"
	"net"
	"net/url"
	"sync"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

// some ugly global variables
// but websocket does not seem to have a way around it
var _websocketGameConfig gameConfigDef
var _nextConnectionId int32 = 1
var _nextConnectionMutex sync.Mutex

type websocketReqDef struct {
	gameConfig   *gameConfigDef
	connectionId int32
	sessionId	string
	remoteAddr   net.Addr
	command      string
	jsonData     []byte
	data         []byte
}

type websocketResDef struct {
	command  string
	jsonData []byte
	data     []byte
}

type receiveWebsocketDef struct {
	message []byte
	err error
}

func websocketSetConfig(gameConfig gameConfigDef) {
	_websocketGameConfig = gameConfig
}

func websocketHandler(ws *websocket.Conn) {

	var websocketReq websocketReqDef

	var url *url.URL = ws.Request().URL
	var connectionId int32
	var gkErr *gkerr.GkErrDef

	defer ws.Close()

	connectionId = getNextConnectionId()

	var config *websocket.Config
	config = ws.Config()

	gklog.LogTrace(fmt.Sprintf("new weboscket connection id: %d\nlocalAddr: %+v\nremoteAddr: %+v\npath: %+v\nquery: %+v\nLocation: %+v\nOrigin: %+v\nreq_adr: %+v",connectionId, ws.LocalAddr(), ws.RemoteAddr(), url.Path, url.RawQuery, config.Location.Host, config.Origin.Host, ws.Request().RemoteAddr))

	if url.Path != _websocketGameConfig.WebsocketPath {
		gkErr = gkerr.GenGkErr("invalid websocket path: " + url.Path, nil, ERROR_ID_WEBSOCKET_INVALID_PATH)
		gklog.LogGkErr("", gkErr)
		return
	}

	var sessionId string

	sessionId = openSessionWebsocket(url.RawQuery, ws.Request().RemoteAddr, connectionId) 
	if sessionId == "" {
		gkErr = gkerr.GenGkErr("session not valid", nil, ERROR_ID_WEBSOCKET_INVALID_SESSION)
		gklog.LogGkErr("", gkErr)
		return
	}

	defer func() {
		closeSessionWebsocket(connectionId, sessionId)
	}()

	var runtimeChan chan runtimeWebsocketReqDef

	runtimeChan, gkErr = addNewWebsocketLink(connectionId, sessionId)
	if gkErr != nil {
		gklog.LogGkErr("addNewWebsocketLink", gkErr)
		return
	}

	defer func() {
		go goRemoveWebsocketLink(connectionId)
	}()

	var receiveWebsocketChan chan *receiveWebsocketDef = make(chan *receiveWebsocketDef)

	go goGetMessage(ws, receiveWebsocketChan)

	for {
		var runtimeWebsocketReq runtimeWebsocketReqDef
		var receiveWebsocket *receiveWebsocketDef

		select {
		case receiveWebsocket = <- receiveWebsocketChan:
			gklog.LogTrace("got websocket message from client: " + string(receiveWebsocket.message))

			if receiveWebsocket.err != nil {
				if receiveWebsocket.err == io.EOF {
					gklog.LogTrace(fmt.Sprintf("closing websocket got eof connectionId: %d",connectionId))
					break
				}
				gkErr = gkerr.GenGkErr(fmt.Sprintf("got websocket input error connectionId: %d",connectionId), receiveWebsocket.err, ERROR_ID_WEBSOCKET_RECEIVE)
				gklog.LogGkErr("websocket error", gkErr)
				return
			} else {
				websocketReq.gameConfig = &_websocketGameConfig
				websocketReq.connectionId = connectionId
				websocketReq.remoteAddr = ws.RemoteAddr()
				websocketReq.command, websocketReq.jsonData, websocketReq.data, gkErr =
					getCommandJsonData(receiveWebsocket.message)
				gklog.LogTrace("got websocket command: " + websocketReq.command + " json: " + string(websocketReq.jsonData) + " data: " + string(websocketReq.data))

				var websocketRes *websocketResDef

				websocketRes, gkErr = dispatchWebsocketRequest(&websocketReq)
				if gkErr != nil {
					gklog.LogGkErr(fmt.Sprintf("websocketReq.doRequest connectionId: %d",connectionId), gkErr)
					return
				}

				gkErr = sendWebsocketRes(ws, websocketRes.command, websocketRes.jsonData, websocketRes.data)
				if gkErr != nil {
					gklog.LogGkErr(fmt.Sprintf("sendWebsocketRes connectionId: %d",connectionId), gkErr)
					return
				}
			}

		case runtimeWebsocketReq = <- runtimeChan:
			gklog.LogTrace("got message from websocket runtime context")

			gkErr = sendWebsocketRes(ws, runtimeWebsocketReq.command, runtimeWebsocketReq.jsonData, runtimeWebsocketReq.data)
			if gkErr != nil {
				gklog.LogGkErr(fmt.Sprintf("sendWebsocketRes connectionId: %d",connectionId), gkErr)
				return
			}
		}
	}
}

func sendWebsocketRes(ws *websocket.Conn, command string, jsonData []byte, data []byte) *gkerr.GkErrDef {
	gklog.LogTrace("sending websocket response: " + fmt.Sprintf("c: %s j: %s d: %s", command, jsonData, data))

	var websocketResponse []byte
	var err error
	var gkErr *gkerr.GkErrDef

	websocketResponse = make([]byte, 0, 0)
	websocketResponse = append(websocketResponse, []byte(command)...)
	websocketResponse = append(websocketResponse, '~')
	websocketResponse = append(websocketResponse, jsonData...)
	websocketResponse = append(websocketResponse, '~')
	websocketResponse = append(websocketResponse, data...)

	err = websocket.Message.Send(ws, string(websocketResponse))
	if err != nil {
		gkErr = gkerr.GenGkErr("websocket.Message.Send", err, ERROR_ID_WEBSOCKET_SEND)
		gklog.LogGkErr("websocket.Message.Send", gkErr)
		return gkErr
	}

	return nil
}

func getNextConnectionId() int32 {

	var connectionId int32

	_nextConnectionMutex.Lock()
	defer _nextConnectionMutex.Unlock()

	connectionId = _nextConnectionId
	_nextConnectionId += 1

	return connectionId
}

func getCommandJsonData(message []byte) (string, []byte, []byte, *gkerr.GkErrDef) {
	var index1, index2 int

	index1 = bytes.IndexByte(message, '~')
	if index1 == -1 {
		return "", nil, nil, gkerr.GenGkErr("missing ~ from websocket message", nil, ERROR_ID_UNKNOWN_WEBSOCKET_INPUT)
	}

	index2 = bytes.IndexByte(message[index1+1:], '~')
	if index2 == -1 {
		return "", nil, nil, gkerr.GenGkErr("missing second ~ from websocketMessage", nil, ERROR_ID_UNKNOWN_WEBSOCKET_INPUT)
	}

	index2 += index1 + 1

	return string(message[:index1]), message[index1+1 : index2], message[index2+1:], nil
}

func goGetMessage(ws *websocket.Conn, ch chan *receiveWebsocketDef) {

	var receiveWebsocket *receiveWebsocketDef
	var err error

	for {
		var message []byte
		message = make([]byte,0,0)
		err = websocket.Message.Receive(ws, &message)
		receiveWebsocket = new(receiveWebsocketDef)
		receiveWebsocket.message = message
		receiveWebsocket.err = err
gklog.LogTrace(fmt.Sprintf("got websocket message before chan write err: %+v",err))
		ch <- receiveWebsocket
		if err != nil {
			gklog.LogTrace("exit goGetMessage due to error")
			break
		}
	}
}


