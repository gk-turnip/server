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
	remoteAddr   net.Addr
	requestPath  string
	requestQuery string
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

	var runtimeChan chan runtimeWebsocketReqDef

	runtimeChan, gkErr = addNewWebsocketLink(connectionId)
	if gkErr != nil {
		gklog.LogGkErr("addNewWebsocketLink", gkErr)
		return
	}

	gklog.LogTrace(fmt.Sprintf("new weboscket connection id: %d addr: %+v path: %+v query: %+v",connectionId, ws.RemoteAddr(), url.Path, url.RawQuery))

	var receiveWebsocketChan chan *receiveWebsocketDef = make(chan *receiveWebsocketDef)

	go goGetMessage(ws, receiveWebsocketChan)

	for {
//		var message []byte
//
//		err = websocket.Message.Receive(ws, &message)
//		if err != nil {
//			gkErr = gkerr.GenGkErr("websocket.Message.Receive", err, ERROR_ID_WEBSOCKET_RECEIVE)
//			gklog.LogGkErr("websocket.Message.Receive", gkErr)
//			return
//		}

		var runtimeWebsocketReq runtimeWebsocketReqDef
		var receiveWebsocket *receiveWebsocketDef

		select {
		case receiveWebsocket = <- receiveWebsocketChan:
			gklog.LogTrace("got websocket message from client: " + string(receiveWebsocket.message))

			websocketReq.gameConfig = &_websocketGameConfig
			websocketReq.connectionId = connectionId
			websocketReq.remoteAddr = ws.RemoteAddr()
			websocketReq.requestPath = url.Path
			websocketReq.requestQuery = url.RawQuery
			websocketReq.command, websocketReq.jsonData, websocketReq.data, gkErr =
				getCommandJsonData(receiveWebsocket.message)
			gklog.LogTrace("got websocket command: " + websocketReq.command + " json: " + string(websocketReq.jsonData) + " data: " + string(websocketReq.data))

			var websocketRes *websocketResDef

			websocketRes, gkErr = dispatchWebsocketRequest(&websocketReq)
			if gkErr != nil {
				gklog.LogGkErr("websocketReq.doRequest", gkErr)
				return
			}

			gkErr = sendWebsocketRes(ws, websocketRes.command, websocketRes.jsonData, websocketRes.data)
			if gkErr != nil {
				gklog.LogGkErr("sendWebsocketRes", gkErr)
				return
			}

		case runtimeWebsocketReq = <- runtimeChan:
			gklog.LogTrace("got message from websocket runtime context")

			gkErr = sendWebsocketRes(ws, runtimeWebsocketReq.command, runtimeWebsocketReq.jsonData, runtimeWebsocketReq.data)
			if gkErr != nil {
				gklog.LogGkErr("sendWebsocketRes", gkErr)
				return
			}
		}

//		gklog.LogTrace("sending websocket response: " + fmt.Sprintf("c: %s j: %s d: %s", websocketRes.command, websocketRes.jsonData, websocketRes.data))
//
//		var websocketResponse []byte
//
//		websocketResponse = make([]byte, 0, 0)
//		websocketResponse = append(websocketResponse, []byte(websocketRes.command)...)
//		websocketResponse = append(websocketResponse, '~')
//		websocketResponse = append(websocketResponse, websocketRes.jsonData...)
//		websocketResponse = append(websocketResponse, '~')
//		websocketResponse = append(websocketResponse, websocketRes.data...)
//
//		err = websocket.Message.Send(ws, string(websocketResponse))
//		if err != nil {
//			gklog.LogGkErr("websocket.Message.Send", gkerr.GenGkErr("websocket.Message.Send", err, ERROR_ID_WEBSOCKET_SEND))
//			return
//		}
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


