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

package ws

// documentation on go websocket: http://godoc.org/code.google.com/p/go.net/websocket
// getting go websocket: go get code.google.com/p/go.net/websocket

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	//"net"
	"net/url"
	//"sync"
)

import (
	"gk/game/field"
	"gk/game/message"
	"gk/game/ses"
	"gk/gkerr"
	"gk/gklog"
)

var _wsContext *WsContextDef

type receiveWebsocketDef struct {
	message []byte
	err     error
}

func SetGlobalWsContext(wsContext *WsContextDef) {
	_wsContext = wsContext
}

func WebsocketHandler(ws *websocket.Conn) {

	var url *url.URL = ws.Request().URL
	var gkErr *gkerr.GkErrDef

	defer ws.Close()

	gklog.LogTrace("WebsocketHandler start")
	defer gklog.LogTrace("WebsocketHandler end")

	//	var websocketConfig *websocket.Config
	//	websocketConfig = ws.Config()

	if url.Path != _wsContext.gameConfig.WebsocketPath {
		gkErr = gkerr.GenGkErr("invalid websocket path: "+url.Path, nil, ERROR_ID_WEBSOCKET_INVALID_PATH)
		gklog.LogGkErr("", gkErr)
		return
	}

	var sessionId string

	sessionId = _wsContext.sessionContext.OpenSessionWebsocket(url.RawQuery, ws.Request().RemoteAddr)
	if sessionId == "" {
		gkErr = gkerr.GenGkErr("session not valid", nil, ERROR_ID_WEBSOCKET_INVALID_SESSION)
		gklog.LogGkErr("", gkErr)
		return
	}

	defer func() {
		_wsContext.sessionContext.CloseSessionWebsocket(sessionId)
	}()

	var singleSession *ses.SingleSessionDef
	var singleWs *singleWsDef

	singleSession = _wsContext.sessionContext.GetSessionFromId(sessionId)
	singleWs = _wsContext.newSingleWs(singleSession)

	var websocketOpenedMessage field.WebsocketOpenedMessageDef
	websocketOpenedMessage.SessionId = sessionId
	websocketOpenedMessage.MessageToClientChan = singleWs.messageToClientChan
	_wsContext.fieldContext.WebsocketOpenedChan <- websocketOpenedMessage

	defer func() {
		var websocketClosedMessage field.WebsocketClosedMessageDef
		websocketClosedMessage.SessionId = sessionId
		_wsContext.fieldContext.WebsocketClosedChan <- websocketClosedMessage
	}()

	var receiveWebsocketChan chan *receiveWebsocketDef = make(chan *receiveWebsocketDef)

	go goGetMessage(ws, receiveWebsocketChan)

	var done bool = false
	for !done {
		var receiveWebsocket *receiveWebsocketDef

		select {
		case receiveWebsocket = <-receiveWebsocketChan:

			if receiveWebsocket.err != nil {
				if receiveWebsocket.err == io.EOF {
					gklog.LogTrace(fmt.Sprintf("closing websocket got eof sessionId: %s", sessionId))
					done = true
					break
				}
				gkErr = gkerr.GenGkErr(fmt.Sprintf("got websocket input error sessionId %s", sessionId), receiveWebsocket.err, ERROR_ID_WEBSOCKET_RECEIVE)
				gklog.LogGkErr("websocket error", gkErr)
				return
			} else {
				var messageFromClient *message.MessageFromClientDef = new(message.MessageFromClientDef)
				messageFromClient.PopulateFromMessage(sessionId, receiveWebsocket.message)

				_wsContext.fieldContext.MessageFromClientChan <- messageFromClient
			}

		case messageToClient := <-singleWs.messageToClientChan:

			gkErr = sendWebsocketMessage(ws, messageToClient)
			if gkErr != nil {
				gklog.LogGkErr(fmt.Sprintf("sendWebsocketMessage sessionId: %s", sessionId), gkErr)
				return
			}
		}
	}
}

func sendWebsocketMessage(ws *websocket.Conn, messageToClient *message.MessageToClientDef) *gkerr.GkErrDef {

	var websocketMessage []byte
	var err error
	var gkErr *gkerr.GkErrDef

	websocketMessage = make([]byte, 0, 0)
	websocketMessage = append(websocketMessage, []byte(messageToClient.Command)...)
	websocketMessage = append(websocketMessage, '~')
	websocketMessage = append(websocketMessage, messageToClient.JsonData...)
	websocketMessage = append(websocketMessage, '~')
	websocketMessage = append(websocketMessage, messageToClient.Data...)

	err = websocket.Message.Send(ws, string(websocketMessage))
	if err != nil {
		gkErr = gkerr.GenGkErr("websocket.Message.Send", err, ERROR_ID_WEBSOCKET_SEND)
		gklog.LogGkErr("websocket.Message.Send", gkErr)
		return gkErr
	}

	return nil
}

func goGetMessage(ws *websocket.Conn, ch chan *receiveWebsocketDef) {

	var receiveWebsocket *receiveWebsocketDef
	var err error

	for {
		var message []byte
		message = make([]byte, 0, 0)
		err = websocket.Message.Receive(ws, &message)
		receiveWebsocket = new(receiveWebsocketDef)
		receiveWebsocket.message = message
		receiveWebsocket.err = err
		ch <- receiveWebsocket
		if err != nil {
			gklog.LogTrace("exit goGetMessage due to error")
			break
		}
	}
}
