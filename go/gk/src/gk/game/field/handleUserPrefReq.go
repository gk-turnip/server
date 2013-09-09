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
	"gk/game/ses"
	"gk/gkerr"
	"gk/gklog"
)

type userPrefReqDef struct {
	PrefName string
	PrefValue     string
}

// websocketConnectionContext entry must be moved from old pod to new pod
func (fieldContext *FieldContextDef) handleUserPrefReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var userPrefReq userPrefReqDef
	var gkErr *gkerr.GkErrDef
	var err error

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(messageFromClient.SessionId)

	err = json.Unmarshal(messageFromClient.JsonData, &userPrefReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	gkErr = fieldContext.persistenceContext.SetUserPref(singleSession.GetUserName(), userPrefReq.PrefName, userPrefReq.PrefValue)
	if gkErr != nil {
		// inserting user preferences is non critical
		// so just log the error
		gklog.LogGkErr("fieldContext.persistenceContext.SetUserPref", gkErr)
	}

	return nil
}

