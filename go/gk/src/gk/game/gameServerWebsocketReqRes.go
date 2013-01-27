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

import (
	"gk/gkerr"
)

const _getSvgReq = "getSvgReq"
const _getSvgRes = "getSvgRes"

const _turnOnRainReq = "turnOnRainReq"
const _turnOffRainReq = "turnOffRainReq"

func dispatchWebsocketRequest(websocketReq *websocketReqDef) (*websocketResDef, *gkerr.GkErrDef) {

	var websocketRes *websocketResDef
	var gkErr *gkerr.GkErrDef

	switch websocketReq.command {
	case _getSvgReq:
		websocketRes, gkErr = doGetSvgReq(websocketReq)
		if gkErr != nil {
			return nil, gkErr
		}
	default:
		gkErr = gkerr.GenGkErr("unknonwn websocket request", nil, ERROR_ID_UNKNOWN_WEBSOCKET_COMMAND)
		return nil, gkErr
	}

	return websocketRes, nil
}
