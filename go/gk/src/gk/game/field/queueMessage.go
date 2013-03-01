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

package field

import (
	"gk/game/message"
	"gk/gkerr"
	"gk/gklog"
)

func (websocketConnectionContext *websocketConnectionContextDef) initQueue() {
	websocketConnectionContext.toClientQueue.queueSize = 0
	websocketConnectionContext.toClientQueue.messagesChan = make(chan *message.MessageToClientDef, MAX_MESSAGES_TO_CLIENT_QUEUE+1)
	websocketConnectionContext.toClientQueue.doneChan = make(chan bool)

	go websocketConnectionContext.runQueue()
}

func (websocketConnectionContext *websocketConnectionContextDef) closeQueue() {
	// put this in go routine
	// since the message to the client may be stuck
	go func() {
		websocketConnectionContext.toClientQueue.doneChan <- true
	}()
}

func (fieldContext *FieldContextDef) queueMessageToClient(sessionId string, messageToClient *message.MessageToClientDef) {

	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef

	gklog.LogTrace("queu up message " + messageToClient.Command)

	websocketConnectionContext, gkErr =
		fieldContext.getWebsocketConnectionContextById(sessionId)

	if gkErr != nil {
		gklog.LogGkErr("", gkErr)
	} else {
		var localSize int

		websocketConnectionContext.toClientQueue.mutex.Lock()
		localSize = websocketConnectionContext.toClientQueue.queueSize
		websocketConnectionContext.toClientQueue.mutex.Unlock()

		if localSize > MAX_MESSAGES_TO_CLIENT_QUEUE {
			gkErr = gkerr.GenGkErr("messageToClient queue overflow, dropping message", nil, ERROR_ID_MESSAGE_TO_CLIENT_QUEUE_OVERFLOW)
			gklog.LogGkErr("", gkErr)
		} else {
			websocketConnectionContext.toClientQueue.mutex.Lock()
			websocketConnectionContext.toClientQueue.queueSize += 1
			websocketConnectionContext.toClientQueue.mutex.Unlock()
			websocketConnectionContext.toClientQueue.messagesChan <- messageToClient
		}
	}
}

func (websocketConnectionContext *websocketConnectionContextDef) runQueue() {

	var done bool
	done = false
	for !done {
		var messageToClient *message.MessageToClientDef

		select {
		case messageToClient = <-websocketConnectionContext.toClientQueue.messagesChan:
		case done = <-websocketConnectionContext.toClientQueue.doneChan:
		}
		if !done {
			gklog.LogTrace("got message to send: " + messageToClient.Command)
			select {
			case websocketConnectionContext.messageToClientChan <- messageToClient:
			case done = <-websocketConnectionContext.toClientQueue.doneChan:
			}
			if !done {
				websocketConnectionContext.toClientQueue.mutex.Lock()
				websocketConnectionContext.toClientQueue.queueSize -= 1
				websocketConnectionContext.toClientQueue.mutex.Unlock()
			}
		}
	}
}
