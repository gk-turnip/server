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
	"gk/game/iso"
)

type newPodReqDef struct {
	PodId string
	destination iso.IsoXYZDef
}

// websocketConnectionContext entry must be moved from old pod to new pod
func (fieldContext *FieldContextDef) handleNewPodReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var newPodReq newPodReqDef
	var gkErr *gkerr.GkErrDef
	var err error

	err = json.Unmarshal(messageFromClient.JsonData, &newPodReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	var websocketConnectionContext *websocketConnectionContextDef

	websocketConnectionContext, gkErr = fieldContext.getWebsocketConnectionContextById(messageFromClient.SessionId)
	if gkErr != nil {
		return gkErr
	}

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(websocketConnectionContext.sessionId)

	var oldPodId int32 = singleSession.GetCurrentPodId()

	var newPodId int64
	newPodId, _ = strconv.ParseInt(newPodReq.PodId, 10, 32)


	if (fieldContext.isPodIdValid(int32(newPodId))) && (oldPodId != int32(newPodId)) {
		gkErr = fieldContext.moveAllAvatarBySessionId(messageFromClient.SessionId, oldPodId, int32(newPodId))
		if gkErr != nil {
			gklog.LogGkErr("", gkErr)
			return gkErr
		}
		delete(fieldContext.podMap[oldPodId].websocketConnectionMap, messageFromClient.SessionId)

		singleSession.SetCurrentPodId(int32(newPodId))

		fieldContext.podMap[int32(newPodId)].websocketConnectionMap[messageFromClient.SessionId] = websocketConnectionContext

		for objKey, fieldObject := range fieldContext.podMap[newPodId].avatarMap {
			if fieldObject.sourceSessionId == websocketConnectionContext.sessionId {
				fieldContext.podMap[newPodId].avatarMap[objKey].isoXYZ = newPodReq.destination
				break
			}
		}

		fieldContext.uploadNewPodInfo(websocketConnectionContext, int32(newPodId))
	} else {
		gkErr = gkerr.GenGkErr(fmt.Sprintf("invalid podId: %d", newPodId), nil, ERROR_ID_INVALID_POD_ID)
		gklog.LogGkErr("", gkErr)
		return gkErr
	}

	return nil
}

func (fieldContext *FieldContextDef) uploadNewPodInfo(websocketConnectionContext *websocketConnectionContextDef, podId int32) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	gkErr = fieldContext.loadTerrain(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	gkErr = fieldContext.sendAllAvatarObjects(podId, websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	return nil
}
