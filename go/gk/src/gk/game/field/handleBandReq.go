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
)

import (
	"gk/game/message"
	"gk/gkerr"
)

type bandReqDef struct {
	BandId string
}

func (fieldContext *FieldContextDef) handleBandReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var bandReq bandReqDef
	var gkErr *gkerr.GkErrDef
	var err error
	var dataOut = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla vitae odio eget libero eleifend commodo vel eu velit turpis duis."

	err = json.Unmarshal(messageFromClient.JsonData, &bandReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	messageToClient.Command = message.BandRes
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"bandId\": \"%s\" }", bandReq.BandId, "{ \"in\": \"%s\" }", dataOut))

	fieldContext.queueMessageToClient(messageFromClient.SessionId, messageToClient)

	return nil
}
