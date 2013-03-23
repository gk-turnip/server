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

package ses

import (
	"strings"
	"sync"
	"time"
)

import (
	"gk/gkrand"
)

type SessionContextDef struct {
	sessionMutex sync.Mutex
	sessionMap   map[string]*SingleSessionDef
	randContext  *gkrand.GkRandContextDef
}

type SingleSessionDef struct {
	sessionId       string
	remoteAddr      string
	createdTime     time.Time
	isWebsocketOpen bool
	userName        string
}

func NewSessionContext(randContext *gkrand.GkRandContextDef) *SessionContextDef {
	var sessionContext *SessionContextDef = new(SessionContextDef)

	sessionContext.sessionMap = make(map[string]*SingleSessionDef)
	sessionContext.randContext = randContext

	return sessionContext
}

func (sessionContext *SessionContextDef) NewSingleSession(userName string, remoteAddr string) *SingleSessionDef {

	var singleSession *SingleSessionDef = new(SingleSessionDef)
	singleSession.remoteAddr = remoteAddr
	singleSession.userName = userName

	sessionContext.sessionMutex.Lock()
	defer sessionContext.sessionMutex.Unlock()
	for {
		var ok bool

		singleSession.sessionId = sessionContext.genSessionString()
		_, ok = sessionContext.sessionMap[singleSession.sessionId]
		if !ok {
			singleSession.createdTime = time.Now()
			sessionContext.sessionMap[singleSession.sessionId] = singleSession
			break
		}
	}

	return singleSession
}

func (sessionContext *SessionContextDef) OpenSessionWebsocket(rawQuery string, remoteAddr string) string {
	var index int
	var sessionId string = ""

	index = strings.Index(rawQuery, "=")
	if index < 1 {
		return sessionId
	}
	if (index + 1) == len(rawQuery) {
		return sessionId
	}
	if rawQuery[:index] != "ses" {
		return sessionId
	}
	sessionId = rawQuery[index+1:]

	sessionContext.sessionMutex.Lock()
	defer sessionContext.sessionMutex.Unlock()

	var ok bool
	_, ok = sessionContext.sessionMap[sessionId]
	if !ok {
		sessionId = ""
		return sessionId
	}

	var singleSession *SingleSessionDef
	singleSession = sessionContext.sessionMap[sessionId]

	//	this comparison is good for security
	//	it would keep people from stealing a session
	//	by sniffing out the session id of an existing session
	//	but this currently has two problems:
	//	1) apache2 is currently configured to reverse proxy to the game
	//		so the addresses will never match
	//	2) both the session.remoteAddr and remoteAddr have ip and port
	//		and it is unknown at this time if the port will match
	//	so we comment it out for now
	//	if singleSession.remoteAddr != remoteAddr {
	//		sessionId = ""
	//		return sessionId
	//	}

	if singleSession.isWebsocketOpen {
		sessionContext.closeSessionWebsocketNoLock(singleSession.sessionId)
	}

	singleSession.isWebsocketOpen = true

	return sessionId
}

func (sessionContext *SessionContextDef) CloseSessionWebsocket(sessionId string) {
	sessionContext.sessionMutex.Lock()
	defer sessionContext.sessionMutex.Unlock()

	sessionContext.closeSessionWebsocketNoLock(sessionId)
}

func (sessionContext *SessionContextDef) closeSessionWebsocketNoLock(sessionId string) {
	var ok bool
	_, ok = sessionContext.sessionMap[sessionId]
	if ok {
		sessionContext.sessionMap[sessionId].isWebsocketOpen = false
	}
}

func (sessionContext *SessionContextDef) GetSessionFromId(sessionId string) *SingleSessionDef {

	sessionContext.sessionMutex.Lock()
	defer sessionContext.sessionMutex.Unlock()

	var singleSession *SingleSessionDef
	singleSession = sessionContext.sessionMap[sessionId]

	return singleSession
}

func (singleSession *SingleSessionDef) IsSessionWebsocketOpen() bool {
	return singleSession.isWebsocketOpen
}

func (singleSession *SingleSessionDef) GetSessionId() string {
	return singleSession.sessionId
}

func (singleSession *SingleSessionDef) GetUserName() string {
	return singleSession.userName
}

func (sessionContext *SessionContextDef) genSessionString() string {
	return sessionContext.randContext.GetRandomString(12)
}
