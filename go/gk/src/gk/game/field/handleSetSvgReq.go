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
	"fmt"
	"strconv"
)

import (
	"gk/game/message"
	"gk/game/ses"
	"gk/gkerr"
	"gk/gklog"
)

// handleMoveAvatarsSvgReq and handleSetAvatarSvgReq should be refactor to use common code

type setSvgDef struct {
	Id string
	X  string
	Y  string
	Z  string
}

func (fieldContext *FieldContextDef) handleSetAvatarSvgReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef
	var err error
	var setSvg setSvgDef

	gklog.LogTrace("json raw: " + string(messageFromClient.JsonData))

	err = json.Unmarshal(messageFromClient.JsonData, &setSvg)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(messageFromClient.SessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	var fieldObject *fieldObjectDef
	var ok bool
	fieldObject, ok = fieldContext.podMap[podId].avatarMap[setSvg.Id]
	if ok {
		var cord int
		cord, _ = strconv.Atoi(setSvg.X)
		fieldObject.isoXYZ.X = int16(cord)
		cord, _ = strconv.Atoi(setSvg.Y)
		fieldObject.isoXYZ.Y = int16(cord)
		cord, _ = strconv.Atoi(setSvg.Z)
		fieldObject.isoXYZ.Z = int16(cord)

		gklog.LogTrace("one")
		fieldContext.setAllAvatars(messageFromClient.SessionId, fieldObject)
	} else {
		gkErr = gkerr.GenGkErr("move object", nil, ERROR_ID_COULD_NOT_FIND_OBJECT_TO_MOVE)
		gklog.LogGkErr("", gkErr)
	}

	return nil
}

func (fieldContext *FieldContextDef) setAllAvatars(sessionId string, fieldObject *fieldObjectDef) {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

	messageToClient.Command = message.SetSvgReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\", \"x\": %d, \"y\": %d, \"z\": %d }", fieldObject.id, fieldObject.isoXYZ.X, fieldObject.isoXYZ.Y, fieldObject.isoXYZ.Z))

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(sessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	for _, websocketConnectionContext := range fieldContext.podMap[podId].websocketConnectionMap {
		gklog.LogTrace("compare session " + websocketConnectionContext.sessionId + " " + sessionId)

		if websocketConnectionContext.sessionId != sessionId {
			gklog.LogTrace("Trace about to queue up move command")
			fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
		}
	}
}
