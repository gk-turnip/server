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
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"sync"
	"net"
	"net/url"
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
	gameConfig *gameConfigDef
	connectionId int32
	remoteAddr	net.Addr
	requestPath string
	requestQuery string
	command string
	jsonData []byte
	data []byte
}

type websocketResDef struct {
	command string
	jsonData []byte
	data []byte
}

func websocketSetConfig(gameConfig gameConfigDef) {
	_websocketGameConfig = gameConfig
}

func websocketHandler(ws *websocket.Conn) {

	var websocketReq websocketReqDef

	var url *url.URL = ws.Request().URL
	var connectionId int32
	var gkErr *gkerr.GkErrDef
	var err error

	connectionId = getNextConnectionId()

	for {
		var message []byte

		websocket.Message.Receive(ws, &message)

gklog.LogTrace("got websocket message: " + string(message))

		websocketReq.gameConfig = &_websocketGameConfig
		websocketReq.connectionId = connectionId
		websocketReq.remoteAddr = ws.RemoteAddr()
		websocketReq.requestPath = url.Path
		websocketReq.requestQuery = url.RawQuery
		websocketReq.command, websocketReq.jsonData, websocketReq.data, gkErr =
			getCommandJsonData(message)
gklog.LogTrace("got websocket command: " + websocketReq.command + " json: " + string(websocketReq.jsonData) + " data: " + string(websocketReq.data))

		var websocketRes *websocketResDef

		websocketRes, gkErr = dispatchWebsocketRequest(&websocketReq)
		if gkErr != nil {
			gklog.LogGkErr("websocketReq.doRequest", gkErr)
			return
		}

gklog.LogTrace("sending websocket response: " + fmt.Sprintf("c: %s j: %s d: %s",websocketRes.command, websocketRes.jsonData, websocketRes.data))

		var websocketResponse []byte

		websocketResponse = make([]byte,0,0)
		websocketResponse = append(websocketResponse,[]byte(websocketRes.command)...)
		websocketResponse = append(websocketResponse,'~')
		websocketResponse = append(websocketResponse,websocketRes.jsonData...)
		websocketResponse = append(websocketResponse,'~')
		websocketResponse = append(websocketResponse,websocketRes.data...)

		err = websocket.Message.Send(ws, string(websocketResponse))
		if err != nil {
			gklog.LogGkErr("websocket.Message.Send", gkerr.GenGkErr("websocket.Message.Send",err, ERROR_ID_WEBSOCKET_SEND))
			return
		}
	}
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

	index2 = bytes.IndexByte(message[index1 + 1:], '~')
	if index2 == -1 {
		return "", nil, nil, gkerr.GenGkErr("missing second ~ from websocketMessage", nil, ERROR_ID_UNKNOWN_WEBSOCKET_INPUT)
	}

	index2 += index1 + 1

gklog.LogTrace(fmt.Sprintf("i1: %d i2: %d len: %d",index1, index2, len(message)))

	return string(message[:index1]), message[index1 + 1: index2], message[index2 + 1:], nil
}

