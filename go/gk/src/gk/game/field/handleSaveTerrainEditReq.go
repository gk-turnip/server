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
	"fmt"
	"encoding/json"
)

import (
	"gk/game/message"
	"gk/game/ses"
	"gk/gkerr"
	"gk/gklog"
)

type saveTerrainEditReqDef struct {
	TerrainMapMap	map[string] mapEntryDef
	PodId int
//	TerrainMapMap terrainMapMapDef
//	TerrainSvgMap string
//	TerrainWallMap string
//	TerrainAUdioMap string
}

//type terrainMapMapDef struct {
//	MapEntries map[string] mapEntryDef
//}

type mapEntryDef struct {
	X int
	Y int
	Zlist []int
	TerrainName string
	Field bool
}

func (fieldContext *FieldContextDef) handleSaveTerrainEditReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var saveTerrainEditReq saveTerrainEditReqDef
	var gkErr *gkerr.GkErrDef
	var err error

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(messageFromClient.SessionId)

	gklog.LogTrace("handleSaveTerrainEditReq")
	gklog.LogTrace(fmt.Sprintf("singleSession: %+v",singleSession))
	gklog.LogTrace(fmt.Sprintf("messageFromClient.JsonData: %s",string(messageFromClient.JsonData)))

	err = json.Unmarshal(messageFromClient.JsonData, &saveTerrainEditReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	for k, v := range saveTerrainEditReq.TerrainMapMap {
		gklog.LogTrace(fmt.Sprintf("k: %+v",k))
		gklog.LogTrace(fmt.Sprintf("v: %+v",v))
// v has:
// x, y, zlist, terrainName, Field

	}

//
//	gkErr = fieldContext.persistenceContext.SetSaveTerrainEdit(singleSession.GetUserName(), saveTerrainEditReq.PrefName, saveTerrainEditReq.PrefValue)
//	if gkErr != nil {
//		// inserting user preferences is non critical
//		// so just log the error
//		gklog.LogGkErr("fieldContext.persistenceContext.SetSaveTerrainEdit", gkErr)
//	}

	return nil
}

