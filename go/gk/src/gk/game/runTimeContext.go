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
	"fmt"
	"sync"
	"os"
	"strings"
	"bytes"
	"text/template"
)

import (
	"gk/gkerr"
	"gk/gklog"
	"gk/gkcommon"
)

var _websocketMap map[int32]*websocketEntryDef
var _websocketMutex sync.Mutex

type websocketEntryDef struct {
	connectionId int32
	sessionId string
	websocketChan chan runtimeWebsocketReqDef
	localEventContext globalEventContextDef
}

type runtimeWebsocketReqDef struct {
	command string
	jsonData []byte
	data []byte
}

func init() {
	_websocketMap = make(map[int32]*websocketEntryDef)
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

	var websocketEntry *websocketEntryDef = new(websocketEntryDef)

	websocketEntry.connectionId = connectionId
	websocketEntry.sessionId = sessionId
	websocketEntry.websocketChan = make(chan runtimeWebsocketReqDef)

gklog.LogTrace("addNewWebsocketLink about to populate context")
	populateLocalEventContext(&websocketEntry.localEventContext)
gklog.LogTrace("addNewWebsocketLink populate context done")

	setContext(websocketEntry)

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

	var websocketEntry *websocketEntryDef

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

	var globalEventChan chan globalEventContextDef = make(chan globalEventContextDef)

	go goGlobalEventLoop(globalEventChan)

	for {
		var globalEventContext globalEventContextDef

		globalEventContext = <- globalEventChan
		handleNewGlobalEvent(gameConfig, globalEventContext)
	}
}

func handleNewGlobalEvent(gameConfig *gameConfigDef, globalEventContext globalEventContextDef) {
	_websocketMutex.Lock()
	defer _websocketMutex.Unlock()

	gklog.LogTrace(fmt.Sprintf("size of field objects: %d",len(globalEventContext.fieldObjectList)))
	for _, websocketEntry := range _websocketMap {
		handleSingleEventChange(gameConfig, globalEventContext, websocketEntry)
	}
}

func handleSingleEventChange(gameConfig *gameConfigDef, globalEventContext globalEventContextDef, websocketEntry *websocketEntryDef) {
	if websocketEntry.localEventContext.rainOn != globalEventContext.rainOn {
		websocketEntry.localEventContext.rainOn = globalEventContext.rainOn
		setContext(websocketEntry)
	}
	for _, fieldObject := range globalEventContext.fieldObjectList {
		var ok bool
		_, ok = websocketEntry.localEventContext.fieldObjectList[fieldObject.Id]
gklog.LogTrace(fmt.Sprintf("ok? %v",ok))
		if !ok {
gklog.LogTrace("new object ready")
			setContextAddFieldObject(gameConfig, websocketEntry, &fieldObject)
			websocketEntry.localEventContext.fieldObjectList[fieldObject.Id] = fieldObject
		}
	}
	for _, fieldObject := range websocketEntry.localEventContext.fieldObjectList {
		var ok bool
		_, ok = globalEventContext.fieldObjectList[fieldObject.Id]
		if !ok {
gklog.LogTrace("object delete")
			setContextDelFieldObject(websocketEntry, &fieldObject)
			delete(websocketEntry.localEventContext.fieldObjectList, fieldObject.Id)
		}
	}
}

func setContext(websocketEntry *websocketEntryDef) {
	var runtimeWebsocketReq runtimeWebsocketReqDef
	if websocketEntry.localEventContext.rainOn {
		runtimeWebsocketReq.command = _turnOnRainReq
	} else {
		runtimeWebsocketReq.command = _turnOffRainReq
	}
	go goSendCommand(websocketEntry.websocketChan, runtimeWebsocketReq);
}

func goSendCommand(websocketChan chan runtimeWebsocketReqDef, runtimeWebsocketReq runtimeWebsocketReqDef) {
	websocketChan <- runtimeWebsocketReq
}

func setContextAddFieldObject(gameConfig *gameConfigDef, websocketEntry *websocketEntryDef, fieldObject *fieldObjectDef) {
	var runtimeWebsocketReq runtimeWebsocketReqDef
	var gkErr *gkerr.GkErrDef

	runtimeWebsocketReq.command = _addSvgReq

    jsonFileName := gameConfig.SvgDir + string(os.PathSeparator) + fieldObject.fileName + ".json"
    svgFileName := gameConfig.SvgDir + string(os.PathSeparator) + fieldObject.fileName + ".svg"

	var jsonRawData []byte
    jsonRawData, gkErr = gkcommon.GetFileContents(jsonFileName)
    if gkErr != nil {
		gklog.LogGkErr("gkcommon.GetFileContents", gkErr)
        return
    }
    runtimeWebsocketReq.data, gkErr = gkcommon.GetFileContents(svgFileName)
    if gkErr != nil {
		gklog.LogGkErr("gkcommon.GetFileContents", gkErr)
        return
    }
    runtimeWebsocketReq.data = fixSvgData(runtimeWebsocketReq.data)

	runtimeWebsocketReq.jsonData, gkErr = templateTranslateFieldObject(jsonRawData, fieldObject)
	if gkErr != nil {
		gklog.LogGkErr("templateTranslateFieldObject", gkErr);
	}

	go goSendCommand(websocketEntry.websocketChan, runtimeWebsocketReq)
}

func setContextDelFieldObject(websocketEntry *websocketEntryDef, fieldObject *fieldObjectDef) {

	var runtimeWebsocketReq runtimeWebsocketReqDef

	runtimeWebsocketReq.command = _delSvgReq
	runtimeWebsocketReq.jsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}",fieldObject.Id))

	go goSendCommand(websocketEntry.websocketChan, runtimeWebsocketReq)
}

func templateTranslateFieldObject(jsonRawData []byte, fieldObject *fieldObjectDef) ([]byte, *gkerr.GkErrDef) {
	var tmpl *template.Template
	var err error
	var gkErr *gkerr.GkErrDef

	tmpl, err = template.New("fieldObject").Parse(string(jsonRawData))
	if err != nil {
		gkErr = gkerr.GenGkErr("template.New.Parse", err, ERROR_ID_TEMPLATE_PARSE)
		return nil, gkErr
	}

	result := make([]byte, 0, 0)
	var writer *bytes.Buffer = bytes.NewBuffer(result)

	err = tmpl.Execute(writer, fieldObject)
	if err != nil {
		gkErr = gkerr.GenGkErr("tmpl.Execute", err, ERROR_ID_TEMPLATE_EXECUTE)
		return nil, gkErr
	}

	return writer.Bytes(), nil
}


