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
	"strings"
	"encoding/json"
	"fmt"
	"time"
)

import (
	"gk/game/message"
	"gk/gkerr"
	"gk/gkjson"
	"gk/gklog"
)

type chatReqDef struct {
	UserName string
	Message  string
}

func (fieldContext *FieldContextDef) handleChatReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var chatReq chatReqDef
	var gkErr *gkerr.GkErrDef
	var err error
	var msgIn string

	err = json.Unmarshal(messageFromClient.JsonData, &chatReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	msgIn = chatReq.Message
	gklog.LogTrace(fmt.Sprintf("chat: t: %v u: %s m: %s", time.Now(), chatReq.UserName, chatReq.Message))
	chatReq.Message = strings.Replace(chatReq.Message, "<", "&lt;", -1)
	chatReq.Message = strings.Replace(chatReq.Message, ">", "&gt;", -1)
	chatReq.Message = strings.Replace(chatReq.Message, "\\", "&#92;", -1)
	if msgIn != chatReq.Message {
		gklog.LogTrace("Previous message scrubbed")
	}

	messageToClient.Command = message.ChatReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"userName\": \"%s\", \"message\": \"%s\" }", chatReq.UserName, gkjson.JsonEscape(chatReq.Message)))

	fieldContext.sendMessageToAll(messageToClient)

	return nil
}
