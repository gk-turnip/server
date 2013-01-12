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

// documentation on go websocket: http://godoc.org/code.google.com/p/go.net/websocket
// getting go websocket: go get code.google.com/p/go.net/websocket

import (
	"code.google.com/p/go.net/websocket"
)

import (
	"gk/gkerr"
	"gk/gklog"
)


func websocketHandler(ws *websocket.Conn) {

	var err error
	var gkErr *gkerr.GkErrDef

	gklog.LogTrace("handleWebSocket start")
	defer gklog.LogTrace("handleWebSocket finished")

	buf := make([]byte, 1024, 1024)

	for {
		var readCount int
		var writeCount int

		readCount, err = ws.Read(buf)
		if err != nil {
			gkErr = gkerr.GenGkErr("ws.Read", err, ERROR_ID_WEBSOCKET_READ)
			gklog.LogGkErr("ws.Read", gkErr)
			return
		}

		gklog.LogTrace("data: " + string(buf[:readCount]))

		var o []byte = []byte("test")

		writeCount, err = ws.Write([]byte(o))
		if err != nil {
			gkErr = gkerr.GenGkErr("ws.Write", err, ERROR_ID_WEBSOCKET_WRITE)
			gklog.LogGkErr("ws.Read", gkErr)
			return
		}
		if writeCount != len(o) {
			gkErr = gkerr.GenGkErr("write count short len", err, ERROR_ID_WEBSOCKET_SHORT_WRITE)
			gklog.LogGkErr("ws.Read", gkErr)
			return
		}
	}
}

