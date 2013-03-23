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

func (fieldContext *FieldContextDef) handleWebsocketClosed(websocketClosedMessage WebsocketClosedMessageDef) *gkerr.GkErrDef {

	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef
	var ok bool

	websocketConnectionContext, ok = fieldContext.websocketConnectionMap[websocketClosedMessage.SessionId]
	if !ok {
		gkErr = gkerr.GenGkErr("closing already closed session", nil, ERROR_ID_CLOSING_ALREADY_CLOSED_SESSION)
		return gkErr
	}

	websocketConnectionContext.closeQueue()

	delete(fieldContext.websocketConnectionMap, websocketClosedMessage.SessionId)

	fieldContext.removeAllAvatarBySessionId(websocketClosedMessage.SessionId)

	return nil
}
