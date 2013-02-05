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

package game

import (
	"sync"
	"time"
)

import (
	"gk/gkerr"
)

var _websocketMap map[int32]websocketEntryDef = make(map[int32]websocketEntryDef)
var _websocketMutex sync.Mutex

type websocketEntryDef struct {
	connectionId int32
	websocketChan chan runtimeWebsocketReqDef
}

type runtimeWebsocketReqDef struct {
	command string
	jsonData []byte
	data []byte
}

func addNewWebsocketLink (connectionId int32) (chan runtimeWebsocketReqDef, *gkerr.GkErrDef) {

	_websocketMutex.Lock()
	defer _websocketMutex.Unlock()

	var ok bool
	var gkErr *gkerr.GkErrDef

	_, ok = _websocketMap[connectionId]
	if ok {
		gkErr = gkerr.GenGkErr("duplicate websocket id", nil, ERROR_ID_DUPLICATE_WEBSOCKET_ID)
		return nil, gkErr
	}

	var websocketEntry websocketEntryDef

	websocketEntry.connectionId = connectionId
	websocketEntry.websocketChan = make(chan runtimeWebsocketReqDef)
	_websocketMap[connectionId] = websocketEntry

	return websocketEntry.websocketChan, nil
}

func goRemoveWebsocketLink(connectionId int32) {
	_websocketMutex.Lock()
	defer _websocketMutex.Unlock()

	delete(_websocketMap,connectionId)
}

func goRuntimeContextLoop(gameConfig *gameConfigDef) {
	for {
		time.Sleep(time.Second * 10)
		turnOnRain(true)
		time.Sleep(time.Second * 10)
		turnOnRain(false)
	}
}

func turnOnRain(on bool) {
	var runtimeWebsocketReq runtimeWebsocketReqDef

	for _, websocketEntry := range _websocketMap {
		if on {
			runtimeWebsocketReq.command = _turnOnRainReq
		} else {
			runtimeWebsocketReq.command = _turnOffRainReq
		}
		go goSendRain(websocketEntry.websocketChan, runtimeWebsocketReq);
		//websocketEntry.websocketChan <- runtimeWebsocketReq
	}
}

func goSendRain(websocketChan chan runtimeWebsocketReqDef, runtimeWebsocketReq runtimeWebsocketReqDef) {
	websocketChan <- runtimeWebsocketReq
}

