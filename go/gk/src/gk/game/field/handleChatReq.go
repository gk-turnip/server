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
	"container/list"
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

const maxSavedChatLines = 24

//var savedChatMutex *sync.Mutex = new(sync.Mutex)
//var savedChat *list.List = list.New()

type chatMessageDef struct {
	time time.Time
	message string
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

	messageToClient.Command = message.ChatReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"userName\": \"%s\", \"message\": \"%s\" }", chatMessage.userName, gkjson.JsonEscape(chatMessage.message)))

	fieldContext.sendMessageToAll(messageToClient)

	fieldContext.savedChatMutex.Lock()
	defer fieldContext.savedChatMutex.Unlock()

	fieldContext.savedChat.PushFront(chatMessage)

	for fieldContext.savedChat.Len() > maxSavedChatLines {
		fieldContext.savedChat.Remove(fieldContext.savedChat.Back())
	}

	return nil
}

func (fieldContext *FieldContextDef) getPastChatJsonData() []byte {
	var element *list.Element
	var returnValue []byte = make([]byte,0,0)
	var firstElement bool = true
	returnValue = append(returnValue,[]byte("{ \"pastChat\": [")...)

	for element = fieldContext.savedChat.Front(); element != nil; element = element.Next() {
		var line string
		var chatMessage chatMessageDef

		chatMessage = element.Value.(chatMessageDef)

		if !firstElement {
			returnValue = append(returnValue,[]byte(",")...)
		}
		firstElement = false

		var escapedUserName, escapedMessage []byte
		var err error

		escapedUserName, err = json.Marshal(chatMessage.userName)
		if err != nil {
			escapedUserName = []byte(fmt.Sprintf("%v",err))
		}
		escapedMessage, err = json.Marshal(chatMessage.message)
		if err != nil {
			escapedMessage = []byte(fmt.Sprintf("%v",err))
		}

		line = fmt.Sprintf(
			"{ \"userName\": %s,\"message\":%s,\"time\":\"%d\"}",
			string(escapedUserName), string(escapedMessage), chatMessage.time.Unix() * 1000)

		returnValue = append(returnValue,[]byte(line)...)
	}
	returnValue = append(returnValue,[]byte("]}")...)

	gklog.LogTrace("chat history: " + string(returnValue))

	return returnValue
}

