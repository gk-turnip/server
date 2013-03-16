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
	"encoding/json"
)

import (
	"gk/game/message"
	"gk/gkerr"
)

type getSvgDef struct {
	SvgName string
}

func (fieldContext *FieldContextDef) handleGetAvatarSvgReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var getSvg getSvgDef
	var gkErr *gkerr.GkErrDef
	var err error

	err = json.Unmarshal(messageFromClient.JsonData, &getSvg)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	var svgJsonData *message.SvgJsonDataDef = new(message.SvgJsonDataDef)
	svgJsonData.Id = fieldContext.getNextObjectId()
	svgJsonData.IsoXYZ.X = 40
	svgJsonData.IsoXYZ.Y = 40

	gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.svgDir, message.GetAvatarSvgRes, getSvg.SvgName, svgJsonData)
	if gkErr != nil {
		gkErr = gkerr.GenGkErr("messageToClient.BuildSvgMessageToClient", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	fieldContext.queueMessageToClient(messageFromClient.SessionId, messageToClient)

	var fieldObject *fieldObjectDef = new(fieldObjectDef)
	fieldObject.id = svgJsonData.Id
	fieldObject.fileName = getSvg.SvgName
	fieldObject.isoXYZ = svgJsonData.IsoXYZ
	fieldObject.sourceSessionId = messageFromClient.SessionId
	fieldContext.addFieldObject(fieldObject)

	var websocketConnectionContext *websocketConnectionContextDef

	websocketConnectionContext, gkErr = fieldContext.getWebsocketConnectionContextById(messageFromClient.SessionId)

	if gkErr != nil {
		return gkErr
	}

	websocketConnectionContext.avatarId = fieldObject.id

	gkErr = fieldContext.sendNewAvatarToAll(messageFromClient.SessionId, fieldObject.id)
	if gkErr != nil {
		return gkErr
	}

	return nil
}
