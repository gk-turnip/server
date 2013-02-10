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

var _lastSessionId int64 = 1
var _sessionMutex sync.Mutex
var _sessionMap map[string]*sessionDef = make(map[string]*sessionDef)

type sessionDef struct {
	sessionId string
	connectionId int32
	remoteAddr string
	createdTime time.Time
}

func newSession(remoteAddr string) *sessionDef {
	var session *sessionDef = new(sessionDef)

	session.remoteAddr = remoteAddr
	session.connectionId = -1

	_sessionMutex.Lock()
	defer _sessionMutex.Unlock()

	for {
		var ok bool

		_lastSessionId += 1
		_lastSessionId = _lastSessionId & 0x7fffff // 23 bits for counter
		session.sessionId = genSessionString(_lastSessionId)
		_, ok = _sessionMap[session.sessionId]
		if !ok {
			session.createdTime = time.Now()
			_sessionMap[session.sessionId] = session
			break
		}
	}

	return session
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

func openSessionWebsocket(rawQuery string, remoteAddr string, connectionId int32) string {
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

	_sessionMutex.Lock()
	defer _sessionMutex.Unlock()

	var ok bool
	_, ok = _sessionMap[sessionId]
	if !ok {
		sessionId = ""
		return sessionId
	}

	var session *sessionDef
	session = _sessionMap[sessionId]

//	this comparison is good for security
//	it would keep people from stealing a session
//	by sniffing out the session id of an existing session
//	but this currently has two problems:
//	1) apache2 is currently configured to reverse proxy to the game
//		so the addresses will never match
//	2) both the session.remoteAddr and remoteAddr have ip and port
//		and it is unknown at this time if the port will match
//	so we comment it out for now
//	if session.remoteAddr != remoteAddr {
//		sessionId = ""
//		return sessionId
//	}

	if session.connectionId != -1 {
		closeSessionWebsocketNoLock(session.connectionId, session.sessionId)
	}

	session.connectionId = connectionId

	return sessionId
}

func closeSessionWebsocket(connectionId int32, sessionId string) {
	_sessionMutex.Lock()
	defer _sessionMutex.Unlock()

	closeSessionWebsocketNoLock(connectionId, sessionId)
}

func closeSessionWebsocketNoLock(connectionId int32, sessionId string) {
	var ok bool
	_, ok = _sessionMap[sessionId]
	if ok {
		_sessionMap[sessionId].connectionId = -1
	}
}

func getConnectionIdFromSession(sessionId string) int32 {
	var connectionId int32 = -1

	_sessionMutex.Lock()
	defer _sessionMutex.Unlock()

	var ok bool
	_, ok = _sessionMap[sessionId]
	if ok {
		connectionId = _sessionMap[sessionId].connectionId
	}
	
	return connectionId
}

