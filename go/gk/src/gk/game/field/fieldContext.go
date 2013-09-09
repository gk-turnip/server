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
	"bytes"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

import (
	"gk/database"
	"gk/game/iso"
	"gk/game/message"
	"gk/game/persistence"
	"gk/game/ses"
	"gk/gkcommon"
	"gk/gkerr"
	"gk/gklog"
)

const firstPodId = 1

type FieldContextDef struct {
	sessionContext        *ses.SessionContextDef
	persistenceContext    *persistence.PersistenceContextDef
	WebsocketOpenedChan   chan WebsocketOpenedMessageDef
	WebsocketClosedChan   chan WebsocketClosedMessageDef
	MessageFromClientChan chan *message.MessageFromClientDef
	avatarSvgDir          string
	terrainSvgDir         string
	lastObjectId          int64
	lastObjectIdMutex     sync.Mutex
	rainContext           rainContextDef
	savedChatMutex        *sync.Mutex
	savedChat             *list.List
	podMap                map[int32]*podEntryDef
}

type podEntryDef struct {
	podId                  int32
	title                  string
	terrainJson            *terrainJsonDef
	websocketConnectionMap map[string]*websocketConnectionContextDef
	avatarMap              map[string]*fieldObjectDef
	objectMap              map[string]*fieldObjectDef
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

const MAX_MESSAGES_TO_CLIENT_QUEUE = 40

type toClientQueueDef struct {
	messagesChan chan *message.MessageToClientDef
	doneChan     chan bool
	mutex        sync.Mutex
	queueSize    int
}

func NewFieldContext(avatarSvgDir string, terrainSvgDir string, sessionContext *ses.SessionContextDef, persistenceContext *persistence.PersistenceContextDef) (*FieldContextDef, *gkerr.GkErrDef) {
	var fieldContext *FieldContextDef = new(FieldContextDef)
	var gkErr *gkerr.GkErrDef

	fieldContext.avatarSvgDir = avatarSvgDir
	fieldContext.terrainSvgDir = terrainSvgDir
	fieldContext.sessionContext = sessionContext
	fieldContext.persistenceContext = persistenceContext
	fieldContext.WebsocketOpenedChan = make(chan WebsocketOpenedMessageDef)
	fieldContext.WebsocketClosedChan = make(chan WebsocketClosedMessageDef)
	fieldContext.MessageFromClientChan = make(chan *message.MessageFromClientDef)

	var podList []database.DbPodDef

	podList, gkErr = persistenceContext.GetPodsList()
	if gkErr != nil {
		return nil, gkErr
	}

	fieldContext.podMap = make(map[int32]*podEntryDef)
	for _, dbPod := range podList {
		gklog.LogTrace(fmt.Sprintf("populate pod %+v", dbPod))
		var podEntry *podEntryDef = new(podEntryDef)
		podEntry.podId = dbPod.Id
		podEntry.title = dbPod.Title
		podEntry.websocketConnectionMap = make(map[string]*websocketConnectionContextDef)
		podEntry.avatarMap = make(map[string]*fieldObjectDef)
		podEntry.objectMap = make(map[string]*fieldObjectDef)

		podEntry.terrainJson, gkErr = fieldContext.newTerrainMap(podEntry.podId)
		if gkErr != nil {
			return nil, gkErr
		}
		fieldContext.podMap[podEntry.podId] = podEntry
	}

	fieldContext.savedChatMutex = new(sync.Mutex)
	fieldContext.savedChat = list.New()

	//	fieldContext.terrainMap, gkErr = fieldContext.newTerrainMap(fieldContext, fieldContext.terrainSvgDir, fieldContext.persistenceContext)
	//	if gkErr != nil {
	//		return nil, gkErr
	//	}

	return fieldContext, nil
}

func (fieldContext *FieldContextDef) getWebsocketConnectionContextById(sessionId string) (*websocketConnectionContextDef, *gkerr.GkErrDef) {
	var websocketConnectionContext *websocketConnectionContextDef = nil
	var gkErr *gkerr.GkErrDef
	var ok bool

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(sessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	websocketConnectionContext, ok = fieldContext.podMap[podId].websocketConnectionMap[sessionId]
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

func (fieldContext *FieldContextDef) addAvatarObject(podId int32, fieldObject *fieldObjectDef) {

	fieldContext.podMap[podId].avatarMap[fieldObject.id] = fieldObject
}

func (fieldContext *FieldContextDef) addTerrainObject(fieldObject *fieldObjectDef, podId int32) {

	if fieldContext.podMap[podId].objectMap == nil {
		fieldContext.podMap[podId].objectMap = make(map[string]*fieldObjectDef)
	}

	fieldContext.podMap[podId].objectMap[fieldObject.id] = fieldObject

	//	_, ok := fieldContext.globalObjectMap[podId]
	//	_, ok := fieldContext.podMap[podId]
	//	if ok {
	//		podEntry.globalObjectMap = make(map[string]*fieldObjectDef)
	//		fieldContext.globalObjectMap[podId] = make(map[string]*fieldObjectDef)
	//	}
	//	fieldContext.globalObjectMap[podId][fieldObject.id] = fieldObject
}

func (fieldContext *FieldContextDef) delAvatarObject(podId int32, fieldObject *fieldObjectDef) {
	delete(fieldContext.podMap[podId].avatarMap, fieldObject.id)
}

func (fieldContext *FieldContextDef) delTerrainObject(podId int32, fieldObject *fieldObjectDef) {

	delete(fieldContext.podMap[podId].objectMap, fieldObject.id)
	//	delete(fieldContext.globalObjectMap[podId], fieldObject.id)
}

func (fieldContext *FieldContextDef) sendAllAvatarObjects(podId int32, websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	for _, fieldObject := range fieldContext.podMap[podId].avatarMap {

		if fieldObject.sourceSessionId != websocketConnectionContext.sessionId {
			gkErr = fieldContext.sendSingleAvatarObject(websocketConnectionContext, fieldObject)
			if gkErr != nil {
				return gkErr
			}
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
	if fieldObject.sourceSessionId != websocketConnectionContext.sessionId {
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

func (fieldContext *FieldContextDef) sendAllPastChat(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var gkErr *gkerr.GkErrDef

	messageToClient.Command = message.SendPastChatReq
	messageToClient.JsonData, gkErr = fieldContext.getPastChatJsonData()
	if gkErr != nil {
		return gkErr
	}
	messageToClient.Data = make([]byte, 0, 0)

	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

func (fieldContext *FieldContextDef) sendUserPrefRestore(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	var gkErr *gkerr.GkErrDef

	messageToClient.Command = message.UserPrefRestoreReq
	messageToClient.JsonData, gkErr = fieldContext.getUserPrefJsonData(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}
	messageToClient.Data = make([]byte, 0, 0)

	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

func (fieldContext *FieldContextDef) getUserPrefJsonData(websocketConnectionContext *websocketConnectionContextDef) ([]byte, *gkerr.GkErrDef) {
	var singleSession *ses.SingleSessionDef
	var userPrefsList []database.DbUserPrefDef
	var gkErr *gkerr.GkErrDef

	singleSession = fieldContext.sessionContext.GetSessionFromId(websocketConnectionContext.sessionId)

	userPrefsList, gkErr = fieldContext.persistenceContext.GetUserPrefsList(singleSession.GetUserName())

	var jsonData []byte = make([]byte,0,512)

	jsonData = append(jsonData,[]byte("{\"userPrefList\":[")...)

	for i, e := range(userPrefsList) {
		if i > 0 {
			jsonData = append(jsonData,',')
		}
		jsonData = append(jsonData,[]byte(fmt.Sprintf("{\"prefName\": \"%s\", \"prefValue\": \"%s\"}", e.PrefName, e.PrefValue))...)
	}

	jsonData = append(jsonData,']')
	jsonData = append(jsonData,'}')

	return jsonData, gkErr
}

// all except for the current session
func (fieldContext *FieldContextDef) sendNewAvatarToAll(podId int32, sessionId string, id string) *gkerr.GkErrDef {

	var fieldObject *fieldObjectDef
	var gkErr *gkerr.GkErrDef
	var ok bool

	fieldObject, ok = fieldContext.podMap[podId].avatarMap[id]

	if ok {
		for _, websocketConnectionContext := range fieldContext.podMap[podId].websocketConnectionMap {
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

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(sessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	for _, fieldObject := range fieldContext.podMap[podId].avatarMap {
		if fieldObject.sourceSessionId == sessionId {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

			messageToClient.Command = message.DelSvgReq
			messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}", fieldObject.id))
			messageToClient.Data = make([]byte, 0, 0)

			//fieldContext.removeSendRemoveAvatarBySessionId(podId, messageToClient)
			for _, websocketConnectionContext := range fieldContext.podMap[podId].websocketConnectionMap {
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}
			delete(fieldContext.podMap[podId].avatarMap, fieldObject.id)
		}
	}
}

// an avatar is moving from one pod to another
// so delete any object matching by sessionId from old pod
// then add them to the new pod
func (fieldContext *FieldContextDef) moveAllAvatarBySessionId(sessionId string, oldPodId int32, newPodId int32, destinationX int16, destinationY int16, destinationZ int16) *gkerr.GkErrDef {
	gklog.LogTrace("moving all object by session id")
	var gkErr *gkerr.GkErrDef

	for _, fieldObject := range fieldContext.podMap[oldPodId].avatarMap {
		if fieldObject.sourceSessionId == sessionId {
			var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

			messageToClient.Command = message.DelSvgReq
			messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}", fieldObject.id))
			messageToClient.Data = make([]byte, 0, 0)

			for _, websocketConnectionContext := range fieldContext.podMap[oldPodId].websocketConnectionMap {
				if sessionId == websocketConnectionContext.sessionId {
					fieldObject.isoXYZ.X = destinationX
					fieldObject.isoXYZ.Y = destinationY
					fieldObject.isoXYZ.Z = destinationZ
					gklog.LogTrace(fmt.Sprintf("moveAllAvatarBySessionId new destination: %d,%d,%d", fieldObject.isoXYZ.X, fieldObject.isoXYZ.Y, fieldObject.isoXYZ.Z))
				}
				fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
			}

			for _, websocketConnectionContext := range fieldContext.podMap[newPodId].websocketConnectionMap {
				gkErr = fieldContext.sendSingleAvatarObject(websocketConnectionContext, fieldObject)
				if gkErr != nil {
					return gkErr
				}
			}

			delete(fieldContext.podMap[oldPodId].avatarMap, fieldObject.id)
			fieldContext.podMap[newPodId].avatarMap[fieldObject.id] = fieldObject
		}
	}
	return nil
}

// put the avatar back
func (fieldContext *FieldContextDef) reAddAvatarBySessionId(sessionId string, newPodId int32) *gkerr.GkErrDef {
	gklog.LogTrace("re adding an avatar by session id")
	var gkErr *gkerr.GkErrDef

	for _, fieldObject := range fieldContext.podMap[newPodId].avatarMap {
		if fieldObject.sourceSessionId == sessionId {
			websocketConnectionContext := fieldContext.podMap[newPodId].websocketConnectionMap[sessionId]
			gklog.LogTrace(fmt.Sprintf("reAddAvatarBySessionId new destination: %d,%d,%d", fieldObject.isoXYZ.X, fieldObject.isoXYZ.Y, fieldObject.isoXYZ.Z))
			gkErr = fieldContext.sendSingleAvatarObject(websocketConnectionContext, fieldObject)
			if gkErr != nil {
				return gkErr
			}

			//			for _, websocketConnectionContext := range fieldContext.podMap[newPodId].websocketConnectionMap {

			//				if (sessionId == websocketConnectionContext.sessionId) {
			//				}
			//			}
		}
	}

	return nil
}

/*
func (fieldContext *FieldContextDef) removeSendRemoveAvatarBySessionId(podIdFromDisconnect int32, messageToClient *message.MessageToClientDef) {

	for _, websocketConnectionContext := range fieldContext.websocketConnectionMap {
		var singleSession *ses.SingleSessionDef
		singleSession = fieldContext.sessionContext.GetSessionFromId(websocketConnectionContext.sessionId)
		var podIdFromWebsocket int32 = singleSession.GetCurrentPodId()

		if (podIdFromWebsocket == podIdFromDisconnect) {
			fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
		}
	}
}
*/

func (fieldContext *FieldContextDef) sendAllRemoveMessageForObject(sessionId string, fieldObject *fieldObjectDef) {
	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(sessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	messageToClient.Command = message.DelSvgReq
	messageToClient.JsonData = []byte(fmt.Sprintf("{ \"id\": \"%s\"}", fieldObject.id))
	for _, websocketConnectionContext := range fieldContext.podMap[podId].websocketConnectionMap {
		fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
	}
}

func (fieldContext *FieldContextDef) sendChatMessageToAll(messageToClient *message.MessageToClientDef) {
	for _, podEntry := range fieldContext.podMap {
		for _, websocketConnectionContext := range podEntry.websocketConnectionMap {
			fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)
		}
	}
}

func (fieldContext *FieldContextDef) removeAvatarBySessionId(sessionId string) *gkerr.GkErrDef {
	var websocketConnectionContext *websocketConnectionContextDef
	var gkErr *gkerr.GkErrDef
	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(sessionId)
	var podId int32 = singleSession.GetCurrentPodId()

	websocketConnectionContext, gkErr = fieldContext.getWebsocketConnectionContextById(sessionId)
	if gkErr != nil {
		return gkErr
	}

	var fieldObject *fieldObjectDef
	var ok bool

	fieldObject, ok = fieldContext.podMap[podId].avatarMap[websocketConnectionContext.avatarId]
	if ok {
		fieldContext.sendAllRemoveMessageForObject(websocketConnectionContext.sessionId, fieldObject)
	}

	if websocketConnectionContext.avatarId != "" {
		delete(fieldContext.podMap[podId].avatarMap, websocketConnectionContext.avatarId)
	}

	return nil
}

// I think this should be (mostly) moved to the message package
func (fieldContext *FieldContextDef) loadTerrain(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef

	// the clear terrain must be done before the terrain svg
	gkErr = fieldContext.doTerrainClear(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	// the terrain svg must be done before the terrain map
	gkErr = fieldContext.doTerrainSvg(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	gkErr = fieldContext.doTerrainMap(websocketConnectionContext)
	if gkErr != nil {
		return gkErr
	}

	return nil
}

func (fieldContext *FieldContextDef) doTerrainClear(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	messageToClient.Command = message.ClearTerrainReq
	messageToClient.JsonData = []byte("{}")
	messageToClient.Data = make([]byte, 0, 0)
	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

func (fieldContext *FieldContextDef) doTerrainSvg(websocketConnectionContext *websocketConnectionContextDef) *gkerr.GkErrDef {

	var terrainSentMap map[string]string = make(map[string]string)
	var gkErr *gkerr.GkErrDef

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(websocketConnectionContext.sessionId)

	var terrainJson *terrainJsonDef

	terrainJson = fieldContext.podMap[singleSession.GetCurrentPodId()].terrainJson

	for i := 0; i < len(terrainJson.jsonMapData.TileList); i++ {
		var terrain string = terrainJson.jsonMapData.TileList[i].Terrain
		var ok bool

		if terrain != "" {
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
	}

	for i := 0; i < len(terrainJson.jsonMapData.ObjectList); i++ {
		var terrain string = terrainJson.jsonMapData.ObjectList[i].Object
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

	var singleSession *ses.SingleSessionDef
	singleSession = fieldContext.sessionContext.GetSessionFromId(websocketConnectionContext.sessionId)

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	messageToClient.Command = message.SetTerrainMapReq
	var jsonFileName string = fieldContext.terrainSvgDir + string(os.PathSeparator) + "map_terrain_" + strconv.FormatInt(int64(singleSession.GetCurrentPodId()), 10) + ".json"
	messageToClient.JsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
	if gkErr != nil {
		return gkErr
	}
	var lf []byte = []byte("\n")
	var tb []byte = []byte("\t")
	var sp []byte = []byte(" ")
	var nl []byte = []byte("")
	var te []byte = []byte("errain")
	var bj []byte = []byte("bject")
	//	Not typos

	messageToClient.JsonData = bytes.Replace(messageToClient.JsonData, lf, nl, -1)
	messageToClient.JsonData = bytes.Replace(messageToClient.JsonData, sp, nl, -1)
	messageToClient.JsonData = bytes.Replace(messageToClient.JsonData, tb, nl, -1)
	messageToClient.JsonData = bytes.Replace(messageToClient.JsonData, te, nl, -1)
	messageToClient.JsonData = bytes.Replace(messageToClient.JsonData, bj, nl, -1)
	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

// I think this should be (mostly) moved to the message package
func (fieldContext *FieldContextDef) sendUserName(websocketConnectionContext *websocketConnectionContextDef, userName string) *gkerr.GkErrDef {

	var messageToClient *message.MessageToClientDef = new(message.MessageToClientDef)
	messageToClient.Command = message.UserNameReq
	messageToClient.JsonData = []byte("{ \"userName\": \"" + userName + "\" }")
	fieldContext.queueMessageToClient(websocketConnectionContext.sessionId, messageToClient)

	return nil
}

func (fieldContext *FieldContextDef) isPodIdValid(podId int32) bool {
	_, ok := fieldContext.podMap[podId]
	return ok
}
