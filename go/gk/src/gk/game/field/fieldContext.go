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
	"fmt"
	"sync"
	"os"
	"time"
	"strings"
	"strconv"
)

import (
	"gk/game/message"
	"gk/game/iso"
	"gk/gkcommon"
	"gk/gkerr"
	"gk/gklog"
)

type FieldContextDef struct {
	globalFieldObjectMap map[string]*fieldObjectDef
	websocketConnectionMap map[string]*websocketConnectionContextDef
	WebsocketOpenedChan chan WebsocketOpenedMessageDef
	WebsocketClosedChan chan WebsocketClosedMessageDef
	MessageFromClientChan chan *message.MessageFromClientDef
	svgDir string
	lastObjectId int64
	lastObjectIdMutex sync.Mutex
	rainContext rainContextDef
}

type websocketConnectionContextDef struct {
	sessionId string
	messageToClientChan chan<- *message.MessageToClientDef
	toClientQueue toClientQueueDef
	avatarId string
}

type fieldObjectDef struct {
	id string
	fileName string
	isoXYZ iso.IsoXYZDef
	sourceSessionId string
}

type rainContextDef struct {
	rainCurrentlyOn bool
	nextRainEvent time.Time
}

const MAX_MESSAGES_TO_CLIENT_QUEUE = 5

type toClientQueueDef struct {
	messagesChan chan *message.MessageToClientDef
	doneChan chan bool
	mutex sync.Mutex
	queueSize int
}

func NewFieldContext(svgDir string) *FieldContextDef {
	var fieldContext *FieldContextDef = new(FieldContextDef)

	fieldContext.svgDir = svgDir
	fieldContext.globalFieldObjectMap = make(map[string]*fieldObjectDef)
	fieldContext.websocketConnectionMap = make(map[string]*websocketConnectionContextDef)
	fieldContext.WebsocketOpenedChan = make(chan WebsocketOpenedMessageDef)
	fieldContext.WebsocketClosedChan = make(chan WebsocketClosedMessageDef)
	fieldContext.MessageFromClientChan = make(chan *message.MessageFromClientDef)

	return fieldContext
}

func (fieldContext *FieldContextDef) getWebsocketConnectionContextById(sessionId string) (*websocketConnectionContextDef, *gkerr.GkErrDef) {
	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef
	var ok bool

	websocketConnectionContext, ok = fieldContext.websocketConnectionMap[sessionId]
	if !ok {
		gkErr = gkerr.GenGkErr("getWebsocketConnectionContextById", nil, ERROR_ID_COULD_NOT_GET_WEBSOCKET_CONNECTION_CONTEXT)
		return nil, gkErr
	}

	return websocketConnectionContext, nil
}

func (fieldContext *FieldContextDef) getNextObjectId() string {
	fieldContext.lastObjectIdMutex.Lock()
	defer fieldContext.lastObjectIdMutex.Unlock()

	fieldContext.lastObjectId += 1

	return strconv.FormatInt(fieldContext.lastObjectId, 36)
}

func (fieldContext *FieldContextDef) addFieldObject(fieldObject *fieldObjectDef) {
	fieldContext.globalFieldObjectMap[fieldObject.id] = fieldObject
}

func (fieldContext *FieldContextDef) delFieldObject(fieldObject *fieldObjectDef) {
	delete(fieldContext.globalFieldObjectMap,fieldObject.id)
}

func (fieldContext *FieldContextDef) sendAllFieldObjects(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	for _, fieldObject := range fieldContext.globalFieldObjectMap {

		gkErr = fieldContext.sendSingleFieldObject(websocketConnectionContext, fieldObject)
		if gkErr != nil {
			return gkErr
		}
/*
		var svgJsonData *message.SvgJsonDataDef = new(message.SvgJsonDataDef)
		svgJsonData.Id = fieldObject.id
		svgJsonData.IsoXYZ = fieldObject.isoXYZ

		var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
		gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.svgDir, message.AddSvgReq, fieldObject.fileName, svgJsonData)
		if gkErr != nil {
			return gkErr
		}

		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
*/
	}

	return nil
}

func (fieldContext *FieldContextDef) sendSingleFieldObject(websocketConnectionContext *websocketConnectionContextDef, fieldObject *fieldObjectDef) *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	var svgJsonData *message.SvgJsonDataDef = new(message.SvgJsonDataDef)

	svgJsonData.Id = fieldObject.id
	svgJsonData.IsoXYZ = fieldObject.isoXYZ

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.svgDir, message.AddSvgReq, fieldObject.fileName, svgJsonData)
	if gkErr != nil {
		return gkErr
	}

	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

// all except for the current session
func (fieldContext *FieldContextDef) sendNewAvatarToAll(sessionId string, id string) *gkerr.GkErrDef {

	var fieldObject *fieldObjectDef
	var gkErr *gkerr.GkErrDef
	var ok bool

	fieldObject, ok = fieldContext.globalFieldObjectMap[id]

	if ok {
		for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
			if websocketConnectionContext.sessionId != sessionId {
				gkErr = fieldContext.sendSingleFieldObject(websocketConnectionContext, fieldObject)
				if gkErr != nil {
					return gkErr
				}
			}
		}
	}

	return nil
}

func (fieldContext *FieldContextDef) removeAllObjectsBySessionId(sessionId string) {
gklog.LogTrace("removing all object by session id")
	for _, fieldObject := range fieldContext.globalFieldObjectMap {
		if fieldObject.sourceSessionId == sessionId {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

			messageToClient.Command = message.DelSvgReq
			messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}",fieldObject.id))
			messageToClient.Data = make([]byte,0,0)

			for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}
			delete(fieldContext.globalFieldObjectMap, fieldObject.id)
		}
	}
}

func (fieldContext *FieldContextDef) sendAllRemoveMessageForObject(sessionId string, fieldObject *fieldObjectDef) {
	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

	messageToClient.Command = message.DelSvgReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}",fieldObject.id))
	for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
	}
}

func (fieldContext *FieldContextDef) removeAvatarBySessionId(sessionId string) *gkerr.GkErrDef {
	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef

	websocketConnectionContext, gkErr = fieldContext.getWebsocketConnectionContextById(sessionId)
	if gkErr != nil {
		return gkErr
	}

	var fieldObject *fieldObjectDef
	var ok bool

	fieldObject, ok = fieldContext.globalFieldObjectMap[websocketConnectionContext.avatarId]
	if ok {
		fieldContext.sendAllRemoveMessageForObject(websocketConnectionContext.sessionId, fieldObject)
	}

	if websocketConnectionContext.avatarId != "" {
		delete(fieldContext.globalFieldObjectMap, websocketConnectionContext.avatarId)
	}

	return nil
}

// I think this should be (mostly) moved to the message package
func (fieldContext *FieldContextDef) loadTerrain(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var dir *os.File
	var err error
	var gkErr *gkerr.GkErrDef

	dir, err = os.Open(fieldContext.svgDir)
	if err != nil {
        gkErr = gkerr.GenGkErr("could not open svg dir: " + fieldContext.svgDir, err, ERROR_ID_SVG_DIR_OPEN)
		return gkErr
	}

	defer dir.Close()

	var fileNames []string
	fileNames, err = dir.Readdirnames(0)
	if err != nil {
        gkErr = gkerr.GenGkErr("could not open svg dir: " + fieldContext.svgDir, err, ERROR_ID_SVG_DIR_READ)
		return gkErr
	}

	for i := 0; i < len(fileNames);i++ {
		if strings.HasPrefix(fileNames[i],"terrain_") {
			if strings.HasSuffix(fileNames[i],".json") {
				var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
				gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.svgDir, message.LoadTerrainReq, fileNames[i][:len(fileNames[i]) - 5],nil)
				if gkErr != nil {
					return gkErr
				}
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}
		}
	}

	{
		var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
		messageToClient.Command = message.SetTerrainReq
		var jsonFileName string = fieldContext.svgDir + string(os.PathSeparator) + "map_terrain.json"
		messageToClient.JsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
		if gkErr != nil {
			return gkErr
		}
		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
	}

	return nil
}

