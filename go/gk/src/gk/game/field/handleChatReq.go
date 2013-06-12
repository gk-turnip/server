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
	"strings"
	"time"
)

import (
	"gk/game/message"
	"gk/gkerr"
	"gk/gkjson"
	"gk/gklog"
	"gk/database"
)

type chatReqDef struct {
	UserName string
	Message  string
}

const maxSavedChatLines = 50

type chatMessageDef struct {
	time     time.Time
	message  string
	userName string
}

func (fieldContext *FieldContextDef) handleChatReq(messageFromClient *message.MessageFromClientDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var chatReq chatReqDef
	var gkErr *gkerr.GkErrDef
	var err error

	err = json.Unmarshal(messageFromClient.JsonData, &chatReq)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	var chatMessage chatMessageDef
	chatMessage.time = time.Now()
	chatMessage.userName = chatReq.UserName
	chatMessage.message = chatReq.Message

	gklog.LogTrace(fmt.Sprintf("chat: t: %v u: %s m: %s", chatMessage.time, chatMessage.userName, chatMessage.message))

	chatMessage.message = strings.Replace(chatMessage.message, "<", "&lt;", -1)
	chatMessage.message = strings.Replace(chatMessage.message, ">", "&gt;", -1)
	chatMessage.message = strings.Replace(chatMessage.message, "\\", "&#92;", -1)
	chatMessage.message = strings.Replace(chatMessage.message, "&nbsp;", " ", -1)
	chatMessage.message = strings.Replace(chatMessage.message, "&", "&amp;", -1)
	if !strings.ContainsAny(strings.ToLower(chatMessage.message), "abcdefghijklmnopqrstuvwxyz1234567890") {
		// drop message since it has nothing of value
		return nil
	}

	messageToClient.Command = message.ChatReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"userName\": \"%s\", \"message\": \"%s\" }", chatMessage.userName, gkjson.JsonEscape(chatMessage.message)))

	fieldContext.sendChatMessageToAll(messageToClient)

	fieldContext.savedChatMutex.Lock()
	defer fieldContext.savedChatMutex.Unlock()

	fieldContext.savedChat.PushFront(chatMessage)

	for fieldContext.savedChat.Len() > maxSavedChatLines {
		fieldContext.savedChat.Remove(fieldContext.savedChat.Back())
	}

	gkErr = fieldContext.persistenceContext.AddNewChatMessage(chatMessage.userName, chatMessage.message)
	if gkErr != nil {
		// inserting chat is non critical
		// so just log the error
		gklog.LogGkErr("fieldContext.persistenceContext.AddNewChatMessage",gkErr)
	}

	return nil
}

func (fieldContext *FieldContextDef) getPastChatJsonData() ([]byte, *gkerr.GkErrDef) {
	var returnValue []byte = make([]byte, 0, 0)
	var gkErr *gkerr.GkErrDef
	returnValue = append(returnValue, []byte("{ \"pastChat\": [")...)

	var lugChatArchiveList []database.LugChatArchiveDef

	lugChatArchiveList, gkErr = fieldContext.persistenceContext.GetLastChatArchiveEntries(maxSavedChatLines)
	if gkErr != nil {
		return nil, gkErr
	}

	gklog.LogTrace(fmt.Sprintf("len lugChatArchiveList: %d",len(lugChatArchiveList)))

	for i := 0;i < len(lugChatArchiveList);i++ {
		var line string

		lugChatArchive := lugChatArchiveList[i]

		if i > 0 {
			returnValue = append(returnValue, []byte(",")...)
		}

		var escapedUserName, escapedMessage []byte
		var err error

		escapedUserName, err = json.Marshal(lugChatArchive.UserName)
		if err != nil {
			escapedUserName = []byte(fmt.Sprintf("%v", err))
		}
		escapedMessage, err = json.Marshal(lugChatArchive.ChatMessage)
		if err != nil {
			escapedMessage = []byte(fmt.Sprintf("%v", err))
		}

		line = fmt.Sprintf(
			"{ \"userName\": %s,\"message\":%s,\"time\":\"%d\"}",
			string(escapedUserName), string(escapedMessage), lugChatArchive.MessageCreationDate.Unix()*1000)

		returnValue = append(returnValue, []byte(line)...)
	}
	returnValue = append(returnValue, []byte("]}")...)

	gklog.LogTrace("chat history: " + string(returnValue))

	return returnValue, nil
}
