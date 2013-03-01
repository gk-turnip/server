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

package ws

import (
	"gk/game/config"
	"gk/game/field"
	"gk/game/message"
	"gk/game/ses"
)

type WsContextDef struct {
	contextMap     map[string]*singleWsDef
	gameConfig     *config.GameConfigDef
	sessionContext *ses.SessionContextDef
	fieldContext   *field.FieldContextDef
}

type singleWsDef struct {
	singleSession       *ses.SingleSessionDef
	messageToClientChan chan *message.MessageToClientDef
}

func NewWsContext(gameConfig *config.GameConfigDef, sessionContext *ses.SessionContextDef, fieldContext *field.FieldContextDef) *WsContextDef {
	var wsContext *WsContextDef = new(WsContextDef)

	wsContext.contextMap = make(map[string]*singleWsDef)
	wsContext.gameConfig = gameConfig
	wsContext.sessionContext = sessionContext
	wsContext.fieldContext = fieldContext

	return wsContext
}

func (wsContext *WsContextDef) newSingleWs(singleSession *ses.SingleSessionDef) *singleWsDef {
	var singleWs *singleWsDef = new(singleWsDef)

	singleWs.messageToClientChan = make(chan *message.MessageToClientDef)
	singleWs.singleSession = singleSession
	wsContext.contextMap[singleWs.singleSession.GetSessionId()] = singleWs

	return singleWs
}
