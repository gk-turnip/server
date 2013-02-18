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
	"testing"
)

func testSession(t *testing.T) {
	var sessionContext *SessionContextDef
	var session1 *SingleSessionDef
	var session2 *SingleSessionDef
	var id1, id2, id3 string
//	var id4 int32

	sessionContext = NewSessionContext()

	session1 = sessionContext.NewSingleSession("1.1.1.1")
	id1 = session1.GetSessionId()
	session2 = sessionContext.NewSingleSession("1.1.1.2")
	id2 = session2.GetSessionId()
	if id1 == id2 {
		t.Logf("duplicate sessionId")
		t.Fail()
	}

	if len(id1) < 8 {
		t.Logf("short sessionId len: " + id1)
		t.Fail()
	}

	if len(id2) < 8 {
		t.Logf("short sessionId len: " + id2)
		t.Fail()
	}

	id3 = sessionContext.OpenSessionWebsocket("ses=" + session1.GetSessionId(), "1.1.1.1")
	if id3 == "" {
		t.Logf("getSessionFromQuery failed")
		t.Fail()
	}
	if id3 != session1.GetSessionId() {
		t.Logf("getSessionFromQuery failed")
		t.Fail()
	}

//	remove this test for now, see session.go for the reason
//	id3 = openSessionWebsocket("ses=" + session1.sessionId, "1.1.1.2", 1)
//	if id3 != "" {
//		t.Logf("getSessionFromQuery invalid result")
//		t.Fail()
//	}

	id3 = sessionContext.OpenSessionWebsocket("sesx=" + session1.GetSessionId(), "1.1.1.1")
	if id3 != "" {
		t.Logf("getSessionFromQuery invalid result")
		t.Fail()
	}

	id3 = sessionContext.OpenSessionWebsocket("ses=" + session1.sessionId + "x", "1.1.1.1")
	if id3 != "" {
		t.Logf("getSessionFromQuery invalid result")
		t.Fail()
	}

	var ses1 *SingleSessionDef
	ses1 = sessionContext.GetSessionFromId(id3)

	if !ses1.IsSessionWebsocketOpen() {
		t.Logf("session should be open")
		t.Fail()
	}

	sessionContext.CloseSessionWebsocket(session1.sessionId)

	if ses1.IsSessionWebsocketOpen() {
		t.Logf("session should be closed")
		t.Fail()
	}
}
