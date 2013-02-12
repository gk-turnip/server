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
	"os"
	"strings"
)

import (
	"gk/gkerr"
	"gk/gklog"
	"gk/gkcommon"
)

var _websocketMap map[int32]websocketEntryDef = make(map[int32]websocketEntryDef)
var _websocketMutex sync.Mutex

type websocketEntryDef struct {
	connectionId int32
	sessionId string
	websocketChan chan runtimeWebsocketReqDef
}

type runtimeWebsocketReqDef struct {
	command string
	jsonData []byte
	data []byte
}

func addNewWebsocketLink (connectionId int32, sessionId string) (chan runtimeWebsocketReqDef, *gkerr.GkErrDef) {

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
	websocketEntry.sessionId = sessionId
	websocketEntry.websocketChan = make(chan runtimeWebsocketReqDef)
	_websocketMap[connectionId] = websocketEntry

	return websocketEntry.websocketChan, nil
}

func sendWebsocketTerrainLoad(gameConfig *gameConfigDef, connectionId int32) {

	var ok bool
	var gkErr *gkerr.GkErrDef

	_websocketMutex.Lock()
	_, ok = _websocketMap[connectionId]
	_websocketMutex.Unlock()

	if !ok {
		gkErr = gkerr.GenGkErr("could not find connection id", nil, ERROR_ID_COULD_NOT_FIND_CONNECTION_ID)
		gklog.LogGkErr("", gkErr)
		return
	}

	var websocketEntry websocketEntryDef

	_websocketMutex.Lock()
	websocketEntry = _websocketMap[connectionId]
	_websocketMutex.Unlock()

	var file *os.File
	var err error

	file, err = os.Open(gameConfig.SvgDir)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open svg dir: " + gameConfig.SvgDir, err, ERROR_ID_SVG_DIR_OPEN)
		gklog.LogGkErr("",gkErr)
		return
	}

	var fileNames []string
	fileNames, err = file.Readdirnames(0)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open svg dir: " + gameConfig.SvgDir, err, ERROR_ID_SVG_DIR_READ)
		gklog.LogGkErr("",gkErr)
		return
	}

	defer file.Close()

	var runtimeWebsocketReq *runtimeWebsocketReqDef

	for i := 0;i < len(fileNames); i++ {
		if strings.HasPrefix(fileNames[i],"terrain_") {
			if strings.HasSuffix(fileNames[i],".json") {
				runtimeWebsocketReq = createNewLoadTerrainEntry(gameConfig.SvgDir, fileNames[i][:len(fileNames[i]) - 5])
gklog.LogTrace("load Terrain file " + fileNames[i] + " json data: " + string(runtimeWebsocketReq.jsonData))
				websocketEntry.websocketChan <- *runtimeWebsocketReq
			}
		}
	}

	runtimeWebsocketReq = new(runtimeWebsocketReqDef)
	runtimeWebsocketReq.command = _setTerrainReq
// hard coded train for testing
	runtimeWebsocketReq.jsonData = []byte(`
{
	"setList": [
	{ "terrain": "sand", "x": 0, "y": 0 },
	{ "terrain": "sand", "x": 0, "y": 1 },
	{ "terrain": "sand", "x": 0, "y": 2 },
	{ "terrain": "grass", "x": 0, "y": 3 },
	{ "terrain": "grass", "x": 0, "y": 4 },
	{ "terrain": "sand", "x": 1, "y": 0 },
	{ "terrain": "sand", "x": 1, "y": 1 },
	{ "terrain": "grass", "x": 1, "y": 2 },
	{ "terrain": "grass", "x": 1, "y": 3 },
	{ "terrain": "grass", "x": 1, "y": 4 },
	{ "terrain": "grass", "x": 2, "y": 0 },
	{ "terrain": "grass", "x": 2, "y": 1 },
	{ "terrain": "grass", "x": 2, "y": 2 },
	{ "terrain": "grass", "x": 2, "y": 3 },
	{ "terrain": "grass", "x": 2, "y": 4 },
	{ "terrain": "grass", "x": 3, "y": 0 },
	{ "terrain": "grass", "x": 3, "y": 1 },
	{ "terrain": "grass", "x": 3, "y": 2 },
	{ "terrain": "grass", "x": 3, "y": 3 },
	{ "terrain": "grass", "x": 3, "y": 4 },
	{ "terrain": "grass", "x": 4, "y": 0 },
	{ "terrain": "grass", "x": 4, "y": 1 },
	{ "terrain": "grass", "x": 4, "y": 2 },
	{ "terrain": "grass", "x": 4, "y": 3 },
	{ "terrain": "grass", "x": 4, "y": 4 }
	]
}
`)
gklog.LogTrace("set Terrain " + string(runtimeWebsocketReq.jsonData))
	websocketEntry.websocketChan <- *runtimeWebsocketReq
}

func createNewLoadTerrainEntry(dir string, fileName string) *runtimeWebsocketReqDef {
	var runtimeWebsocketReq runtimeWebsocketReqDef
	//var fileContents []byte
	var gkErr *gkerr.GkErrDef
	var jsonFileName, svgFileName string

	jsonFileName = dir + string(os.PathSeparator) + fileName + ".json"
	runtimeWebsocketReq.command = _loadTerrainReq

	runtimeWebsocketReq.jsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
	if gkErr != nil {
		gklog.LogGkErr("terrain GetFileContents",gkErr)
		return nil
	}

	svgFileName = dir + string(os.PathSeparator) + fileName + ".svg"
	runtimeWebsocketReq.data, gkErr = gkcommon.GetFileContents(svgFileName)
	if gkErr != nil {
		gklog.LogGkErr("terrain GetFileContents",gkErr)
		return nil
	}

	runtimeWebsocketReq.data = fixSvgData(runtimeWebsocketReq.data)

	return &runtimeWebsocketReq
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

