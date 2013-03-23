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
	"os"
	"io"
	"strconv"
	"sync"
	"time"
	"encoding/json"
)

import (
	"gk/game/iso"
	"gk/game/message"
	"gk/game/ses"
	"gk/gkcommon"
	"gk/gkerr"
	"gk/gklog"
)

type FieldContextDef struct {
	globalAvatarMap   map[string]*fieldObjectDef
	globalTerrainMap   map[string]*fieldObjectDef
	websocketConnectionMap map[string]*websocketConnectionContextDef
	sessionContext         *ses.SessionContextDef
	WebsocketOpenedChan    chan WebsocketOpenedMessageDef
	WebsocketClosedChan    chan WebsocketClosedMessageDef
	MessageFromClientChan  chan *message.MessageFromClientDef
	avatarSvgDir           string
	terrainSvgDir          string
	lastObjectId           int64
	lastObjectIdMutex      sync.Mutex
	rainContext            rainContextDef
}

type websocketConnectionContextDef struct {
	sessionId           string
	messageToClientChan chan<- *message.MessageToClientDef
	toClientQueue       toClientQueueDef
	avatarId            string
}

type fieldObjectDef struct {
	id              string
	fileName        string
	isoXYZ          iso.IsoXYZDef
	sourceSessionId string
}

type rainContextDef struct {
	rainCurrentlyOn bool
	nextRainEvent   time.Time
}

const MAX_MESSAGES_TO_CLIENT_QUEUE = 20

type toClientQueueDef struct {
	messagesChan chan *message.MessageToClientDef
	doneChan     chan bool
	mutex        sync.Mutex
	queueSize    int
}

func NewFieldContext(avatarSvgDir string, terrainSvgDir string, sessionContext *ses.SessionContextDef) *FieldContextDef {
	var fieldContext *FieldContextDef = new(FieldContextDef)

	fieldContext.avatarSvgDir = avatarSvgDir
	fieldContext.terrainSvgDir = terrainSvgDir
	fieldContext.sessionContext = sessionContext
	fieldContext.globalAvatarMap = make(map[string]*fieldObjectDef)
	fieldContext.globalTerrainMap = make(map[string]*fieldObjectDef)
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

	return "gki_" + strconv.FormatInt(fieldContext.lastObjectId, 36)
}

func (fieldContext *FieldContextDef) addAvatarObject(fieldObject *fieldObjectDef) {
	fieldContext.globalAvatarMap[fieldObject.id] = fieldObject
}

func (fieldContext *FieldContextDef) addTerrainObject(fieldObject *fieldObjectDef) {
	fieldContext.globalTerrainMap[fieldObject.id] = fieldObject
}

func (fieldContext *FieldContextDef) delAvatarObject(fieldObject *fieldObjectDef) {
	delete(fieldContext.globalAvatarMap, fieldObject.id)
}

func (fieldContext *FieldContextDef) delTerrainObject(fieldObject *fieldObjectDef) {
	delete(fieldContext.globalTerrainMap, fieldObject.id)
}

func (fieldContext *FieldContextDef) sendAllAvatarObjects(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	for _, fieldObject := range fieldContext.globalAvatarMap {

		gkErr = fieldContext.sendSingleAvatarObject(websocketConnectionContext, fieldObject)
		if gkErr != nil {
			return gkErr
		}
	}

	return nil
}

func (fieldContext *FieldContextDef) sendSingleAvatarObject(websocketConnectionContext *websocketConnectionContextDef, fieldObject *fieldObjectDef) *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	var svgJsonData *message.SvgJsonDataDef = new(message.SvgJsonDataDef)

	svgJsonData.Id = fieldObject.id
	svgJsonData.IsoXYZ = fieldObject.isoXYZ
	gklog.LogTrace("sourceSessionId: " + fieldObject.sourceSessionId)
	if fieldObject.sourceSessionId != "" {
		var singleSession *ses.SingleSessionDef
		singleSession = fieldContext.sessionContext.GetSessionFromId(fieldObject.sourceSessionId)
		svgJsonData.UserName = singleSession.GetUserName()
		gklog.LogTrace("going to send to ws userName: " + singleSession.GetUserName())
	}

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.avatarSvgDir, message.AddSvgReq, fieldObject.fileName, svgJsonData)
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

	fieldObject, ok = fieldContext.globalAvatarMap[id]

	if ok {
		for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
			if websocketConnectionContext.sessionId != sessionId {
				gkErr = fieldContext.sendSingleAvatarObject(websocketConnectionContext, fieldObject)
				if gkErr != nil {
					return gkErr
				}
			}
		}
	}

	return nil
}

func (fieldContext *FieldContextDef) removeAllAvatarBySessionId(sessionId string) {
	gklog.LogTrace("removing all object by session id")
	for _, fieldObject := range fieldContext.globalAvatarMap {
		if fieldObject.sourceSessionId == sessionId {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

			messageToClient.Command = message.DelSvgReq
			messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}", fieldObject.id))
			messageToClient.Data = make([]byte, 0, 0)

			for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}
			delete(fieldContext.globalAvatarMap, fieldObject.id)
		}
	}
}

func (fieldContext *FieldContextDef) sendAllRemoveMessageForObject(sessionId string, fieldObject *fieldObjectDef) {
	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

	messageToClient.Command = message.DelSvgReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}", fieldObject.id))
	for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
	}
}

func (fieldContext *FieldContextDef) sendMessageToAll(messageToClient *message.MessageToClientDef) {
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

	fieldObject, ok = fieldContext.globalAvatarMap[websocketConnectionContext.avatarId]
	if ok {
		fieldContext.sendAllRemoveMessageForObject(websocketConnectionContext.sessionId, fieldObject)
	}

	if websocketConnectionContext.avatarId != "" {
		delete(fieldContext.globalAvatarMap, websocketConnectionContext.avatarId)
	}

	return nil
}

// I think this should be (mostly) moved to the message package
func (fieldContext *FieldContextDef) loadTerrain(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	gkErr = fieldContext.doTerrainSvg(websocketConnectionContext);
	if gkErr != nil {
		return gkErr
	}

	gkErr = fieldContext.doTerrainMap(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	return nil
}

func (fieldContext *FieldContextDef) doTerrainSvg(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {
/*
	var dir *os.File
	var err error
	var gkErr *gkerr.GkErrDef

	dir, err = os.Open(fieldContext.terrainSvgDir)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open svg dir: "+fieldContext.terrainSvgDir, err, ERROR_ID_SVG_DIR_OPEN)
		return gkErr
	}

	defer dir.Close()

	var fileNames []string
	fileNames, err = dir.Readdirnames(0)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open svg dir: "+fieldContext.terrainSvgDir, err, ERROR_ID_SVG_DIR_READ)
		return gkErr
	}

	for i := 0; i < len(fileNames); i++ {
		if strings.HasPrefix(fileNames[i], "terrain_") {
			if strings.HasSuffix(fileNames[i], ".json") {
				var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
				gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.terrainSvgDir, message.SetTerrainSvgReq, fileNames[i][:len(fileNames[i])-5], nil)
				if gkErr != nil {
					return gkErr
				}
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}
		}
	}

;zzz
*/

	var jsonFileName string = fieldContext.terrainSvgDir + string(os.PathSeparator) + "map_terrain.json"

	var jsonFile *os.File
	var err error
	var gkErr *gkerr.GkErrDef
	var jsonData []byte = make([]byte, 0, 256)

	jsonFile, err = os.Open(jsonFileName)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open: " + jsonFileName, err, ERROR_ID_OPEN_TERRAIN_MAP)
		return gkErr
	}
	defer jsonFile.Close()

	var buf []byte
	var readCount int
	for {
		buf = make([]byte, 128, 128)
		readCount, err = jsonFile.Read(buf)
		if readCount > 0 {
			jsonData = append(jsonData, buf[:readCount]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			gkErr = gkerr.GenGkErr("could not read " + jsonFileName, err, ERROR_ID_READ_TERRAIN_MAP)
			return gkErr
		}
	}

	var terrainMap struct {
		TileList []struct {
			Terrain string
		}
		ObjectList []struct {
			Object string
		}
	}
	err = json.Unmarshal(jsonData, &terrainMap)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return gkErr
	}

	//gklog.LogTrace(fmt.Sprintf("terrainMap: %+v",terrainMap))

	var terrainSentMap map[string]string = make(map[string]string)

	for i := 0; i < len(terrainMap.TileList); i++ {
		var terrain string = terrainMap.TileList[i].Terrain
		var ok bool

		_, ok = terrainSentMap[terrain]
		if !ok {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
			gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.terrainSvgDir, message.SetTerrainSvgReq, terrain, nil)
			if gkErr != nil {
				return gkErr
			}
			fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

			terrainSentMap[terrain] = terrain
		}
	}

	for i := 0; i < len(terrainMap.ObjectList); i++ {
		var terrain string = terrainMap.ObjectList[i].Object
		var ok bool

		_, ok = terrainSentMap[terrain]
		if !ok {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
			gkErr = messageToClient.BuildSvgMessageToClient(fieldContext.terrainSvgDir, message.SetTerrainSvgReq, terrain, nil)
			if gkErr != nil {
				return gkErr
			}
			fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

			terrainSentMap[terrain] = terrain
		}
	}

	return nil
}

func (fieldContext *FieldContextDef) doTerrainMap(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	messageToClient.Command = message.SetTerrainMapReq
	var jsonFileName string = fieldContext.terrainSvgDir + string(os.PathSeparator) + "map_terrain.json"
	messageToClient.JsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
	if gkErr != nil {
		return gkErr
	}
	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

// I think this should be (mostly) moved to the message package
func (fieldContext *FieldContextDef) sendUserName(websocketConnectionContext *websocketConnectionContextDef, userName string) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	messageToClient.Command = message.UserNameReq;
	messageToClient.JsonData = []byte("{ \"userName\": \"" + userName + "\" }")
	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}
