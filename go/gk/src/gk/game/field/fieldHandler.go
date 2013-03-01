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
	"time"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

func (fieldContext *FieldContextDef) StartFieldHandler() {

	var gkErr *gkerr.GkErrDef

	var timerChan <-chan time.Time
	timerChan = time.Tick(time.Second)

	for {
		select {
		case websocketOpenedMessage := <-fieldContext.WebsocketOpenedChan:
			gkErr = fieldContext.handleWebsocketOpened(websocketOpenedMessage)
			if gkErr != nil {
				gklog.LogGkErr("handleWebsocketOpened", gkErr)
			}
		case websocketClosedMessage := <-fieldContext.WebsocketClosedChan:
			gkErr = fieldContext.handleWebsocketClosed(websocketClosedMessage)
			if gkErr != nil {
				gklog.LogGkErr("handleWebsocketClosed", gkErr)
			}
		case messageFromClient := <-fieldContext.MessageFromClientChan:
			gkErr = fieldContext.handleMessageFromClient(messageFromClient)
			if gkErr != nil {
				gklog.LogGkErr("handleMessageFromClient", gkErr)
			}
		case timerMessage := <-timerChan:
			gkErr = fieldContext.handleTicker(timerMessage)
			if gkErr != nil {
				gklog.LogGkErr("handleTicker", gkErr)
			}
		}
	}
}
