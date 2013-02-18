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
	"io"
	"strconv"
	"strings"
	"time"
	"sync"
	"crypto/rand"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

type SessionContextDef struct {
	lastSessionId int64
	sessionMutex sync.Mutex
	sessionMap map[string]*SingleSessionDef
}

type SingleSessionDef struct {
	sessionId string
	remoteAddr string
	createdTime time.Time
	isWebsocketOpen bool
}

func NewSessionContext() *SessionContextDef {
	var sessionContext *SessionContextDef = new(SessionContextDef)

	sessionContext.sessionMap = make(map[string]*SingleSessionDef)

	return sessionContext
}

func (sessionContext *SessionContextDef) NewSingleSession(remoteAddr string) *SingleSessionDef {

	var singleSession *SingleSessionDef = new(SingleSessionDef)
	singleSession.remoteAddr = remoteAddr

	sessionContext.sessionMutex.Lock()
	defer sessionContext.sessionMutex.Unlock()
	for {
		var ok bool

		sessionContext.lastSessionId += 1
		sessionContext.lastSessionId = sessionContext.lastSessionId & 0x7fffff // 23 bits for counter
		singleSession.sessionId = genSessionString(sessionContext.lastSessionId)
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

	index = strings.Index(rawQuery,"=")
	if index < 1 {
		return sessionId
	}
	if (index + 1) == len(rawQuery) {
		return sessionId
	}
	if rawQuery[:index] != "ses" {
		return sessionId
	}
	sessionId = rawQuery[index + 1:]

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

func (singleSessionContext *SingleSessionDef) GetSessionId() string {
	return singleSessionContext.sessionId
}

func genSessionString(sessionId23 int64) string {
	var sessionId63 int64
	var readCount int
	var err error
	var gkErr *gkerr.GkErrDef

	buf := make([]byte, 5, 5)
	readCount, err = io.ReadFull(rand.Reader, buf)
	if (readCount != len(buf)) || (err != nil) {
		// just log the error
		// the system can continue on with a damaged session id
		gkErr = gkerr.GenGkErr("rand io.ReadFull", err, ERROR_ID_RAND)
		gklog.LogGkErr("rand io.ReadFull", gkErr)
	}

	sessionId63 = (sessionId23 << 40) | int64(buf[0]) | (int64(buf[1]) << 8) | (int64(buf[2]) << 16) | (int64(buf[3]) << 24) | (int64(buf[4]) << 32)

	return strconv.FormatInt(sessionId63,36)
}

