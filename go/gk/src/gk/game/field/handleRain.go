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
	"math/rand"
	"time"
)

import (
	"gk/game/message"
	"gk/gkerr"
)

func (fieldContext *FieldContextDef) handleRain() *gkerr.GkErrDef {
	if time.Now().After(fieldContext.rainContext.nextRainEvent) {
		fieldContext.rainContext.rainCurrentlyOn =
			!fieldContext.rainContext.rainCurrentlyOn

		fieldContext.sendEveryoneRainEvent()
		fieldContext.rainContext.nextRainEvent =
			time.Now().Add((time.Second * time.Duration(15)) + time.Duration(rand.Int31n(15)))
	}

	return nil
}

func (fieldContext *FieldContextDef) sendEveryoneRainEvent() {
	var rainCommand string

	if fieldContext.rainContext.rainCurrentlyOn {
		rainCommand = message.TurnOnRainReq
	} else {
		rainCommand = message.TurnOffRainReq
	}
	var podId int32 = firstPodId // rain only in the first pod

	for _, websocketConnectionContext := range fieldContext.podMap[podId].websocketConnectionMap {
		var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
		messageToClient.Command = rainCommand
		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
	}
}
