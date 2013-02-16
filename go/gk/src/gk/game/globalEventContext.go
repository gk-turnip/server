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
	"time"
	"math/rand"
	"sync"
	"strconv"
)

import (
	"gk/gklog"
)

var	_randContext *rand.Rand
var _globalEventContext globalEventContextDef
var _globalEventMutex sync.Mutex
var _lastObjectId int64
var _lastObjectIdMutex sync.Mutex

type globalEventContextDef struct {
	rainOn bool
	fieldObjectList map[string]fieldObjectDef
}

type fieldObjectDef struct {
	Id string
	fileName string
	X,Y,Z int16
}

func init() {
	_lastObjectId = 0
	var source rand.Source

	source = rand.NewSource(time.Now().UnixNano())
	_randContext = rand.New(source)

	_globalEventContext.fieldObjectList = make(map[string]fieldObjectDef)
}

func goGlobalEventLoop(globalEventChan chan globalEventContextDef) {
	var nextRain time.Time

	nextRain = getNextRainTime()

	for {
		time.Sleep(time.Second)
		if time.Now().After(nextRain) {
			nextRain = getNextRainTime()
			_globalEventMutex.Lock()
			if _globalEventContext.rainOn {
				_globalEventContext.rainOn = false
			} else {
				_globalEventContext.rainOn = true
			}
			_globalEventMutex.Unlock()
			sendGlobalEvent(globalEventChan)
		}
		if _globalEventContext.rainOn {
			if rand.Int31n(5) == 2 {
gklog.LogTrace("trying to add dandelion")
				addNewFieldObject("dandelion", int16(rand.Int31n(50)), int16(rand.Int31n(50)), 0)
				sendGlobalEvent(globalEventChan)
			}
		} else{
			if rand.Int31n(4) == 2 {
gklog.LogTrace("trying to del dandelion")
				delNewFieldObject("dandelion")
				sendGlobalEvent(globalEventChan)
			}
		}
	}
}

func sendGlobalEvent(globalEventChan chan globalEventContextDef) {
	_globalEventMutex.Lock()
	var localEventContext globalEventContextDef
	populateLocalEventContextNoLock(&localEventContext)
	_globalEventMutex.Unlock()
	globalEventChan <- localEventContext
}

func addNewFieldObject(fileName string, x int16, y int16, z int16) {
	var fieldObject fieldObjectDef

	fieldObject.Id = getNextObjectId()
	fieldObject.fileName = fileName
	fieldObject.X = x
	fieldObject.Y = y
	fieldObject.Z = z
	_globalEventContext.fieldObjectList[fieldObject.Id] = fieldObject
}

func delNewFieldObject(fileName string) {
	for _, fieldObject := range _globalEventContext.fieldObjectList {
		if fileName == fieldObject.fileName {
			delete(_globalEventContext.fieldObjectList,fieldObject.Id)
			break
		}
	}
}

func populateLocalEventContext(localEventContext *globalEventContextDef) {
	_globalEventMutex.Lock()
	defer _globalEventMutex.Unlock()

	populateLocalEventContextNoLock(localEventContext)
}

func populateLocalEventContextNoLock(localEventContext *globalEventContextDef) {
	localEventContext.rainOn = _globalEventContext.rainOn
	localEventContext.fieldObjectList = make(map[string]fieldObjectDef)
	for _, fieldObject := range _globalEventContext.fieldObjectList {
		localEventContext.fieldObjectList[fieldObject.Id] = fieldObject
	}
}

func getNextRainTime() time.Time {
	return  time.Now().Add(time.Second * (15 + time.Duration(rand.Int31n(15))))
}

func getNextObjectId() string {

	_lastObjectIdMutex.Lock()
	defer _lastObjectIdMutex.Unlock()

	_lastObjectId += 1
	
	return strconv.FormatInt(_lastObjectId, 36)
}

