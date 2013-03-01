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
)

func (fieldContext *FieldContextDef) handleMessageFromClient(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	switch messageFromClient.Command {
	case message.GetAvatarSvgReq:
		gkErr = fieldContext.handleGetAvatarSvgReq(messageFromClient)
		if gkErr != nil {
			return gkErr
		}
	case message.DelAvatarSvgReq:
		gkErr = fieldContext.handleDelAvatarSvgReq(messageFromClient)
		if gkErr != nil {
			return gkErr
		}
	case message.MoveAvatarSvgReq:
		gkErr = fieldContext.handleMoveAvatarSvgReq(messageFromClient)
		if gkErr != nil {
			return gkErr
		}
	default:
		gkErr = gkerr.GenGkErr("unknonwn websocket request: "+messageFromClient.Command, nil, ERROR_ID_UNKNOWN_WEBSOCKET_COMMAND)
		return gkErr
	}

	return nil
}
