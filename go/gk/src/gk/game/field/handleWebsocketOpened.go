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
	"gk/gkerr"
)

func (fieldContext *FieldContextDef) handleWebsocketOpened(websocketOpenedMessage WebsocketOpenedMessageDef) *gkerr.GkErrDef {

	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef
	var ok bool

	websocketConnectionContext, ok = fieldContext.websocketConnectionMap[websocketOpenedMessage.SessionId]
	if ok {
		gkErr = gkerr.GenGkErr("opening already opened session", nil, ERROR_ID_OPENING_ALREADY_OPEN_SESSION)
		return gkErr
	}

	websocketConnectionContext = new(websocketConnectionContextDef)

	websocketConnectionContext.sessionId = websocketOpenedMessage.SessionId
	websocketConnectionContext.messageToClientChan = websocketOpenedMessage.MessageToClientChan
	websocketConnectionContext.initQueue()

	fieldContext.websocketConnectionMap[websocketOpenedMessage.SessionId] = websocketConnectionContext

	gkErr = fieldContext.loadTerrain(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	gkErr = fieldContext.sendAllFieldObjects(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	return nil
}
