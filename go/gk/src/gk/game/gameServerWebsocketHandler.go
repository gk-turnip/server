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
	"bytes"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"os"
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
	rawWebsocketData := make([]byte, 0, 0)

	//	for {
	var readCount int

	readCount, err = ws.Read(buf)

	if readCount > 0 {
		rawWebsocketData = append(rawWebsocketData, buf[0:readCount]...)
	}
	if err != nil {
		//			if err == io.EOF {
		//				break
		//			}
		gkErr = gkerr.GenGkErr("ws.Read", err, ERROR_ID_WEBSOCKET_READ)
		gklog.LogGkErr("ws.Read", gkErr)
		return
	}
	//	}
	var command string
	var websocketData []byte

	gklog.LogTrace(fmt.Sprintf("rawWebsocketData: <%s>\n", rawWebsocketData))

	command, websocketData, gkErr = getCommandAndData(rawWebsocketData)
	if gkErr != nil {
		gklog.LogGkErr("getCommandAndData", gkErr)
		return
	}

	gklog.LogTrace(fmt.Sprintf("command: <%s> data: <%s>\n", command, websocketData))

	switch command {
	case "getSvg":
		sendGetSvgResponse(ws, websocketData)
		//	gkErr = sendResponse(ws)
		//	if gkErr != nil {
		//		return
		//	}
	default:
		gkErr = gkerr.GenGkErr("unknown websocket command: "+command, err, ERROR_ID_INVALID_WEBSOCKET_COMMAND)
		gklog.LogGkErr("unknown websocket command", gkErr)
	}
	//	gkErr = sendResponse(ws)
	//	if gkErr != nil {
	//		return
	//	}
}

func sendGetSvgResponse(ws *websocket.Conn, websocketData []byte) *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	if websocketData[len(websocketData)-1] != '\n' {
		gkErr = gkerr.GenGkErr("invalid getSvg", nil, ERROR_ID_INVALID_WEBSOCKET_GET_SVG)
		gklog.LogGkErr("invalid getSvg", gkErr)
	}

	var svgData []byte

	svgData, gkErr = getSvgData("/home/diver2/gameServer/svg", string(websocketData[0:len(websocketData)-1]))
	if gkErr != nil {
		gklog.LogGkErr("getSvgData", gkErr)
		return gkErr
	}

	gklog.LogTrace("svgData: " + string(svgData))

	svgData = fixSvgData(svgData)

	gklog.LogTrace("fixed svgData: " + string(svgData))

	//	var writeCount int
	//	var err error

	svgResponse := []byte("svg\n")
	svgResponse = append(svgResponse, svgData...)

	websocket.Message.Send(ws, string(svgResponse))
	/*
		writeCount, err = ws.Write(svgResponse)
		if err != nil {
			gkErr = gkerr.GenGkErr("ws.Write", err, ERROR_ID_WEBSOCKET_WRITE)
			gklog.LogGkErr("ws.Read", gkErr)
			return gkErr
		}
		if writeCount != len(svgResponse) {
			gkErr = gkerr.GenGkErr("write count short len", err, ERROR_ID_WEBSOCKET_SHORT_WRITE)
			gklog.LogGkErr("ws.Read", gkErr)
			return gkErr
		}
	*/

	return nil
}

/*
func sendResponse(ws *websocket.Conn) *gkerr.GkErrDef {
	var svgResponse []byte = []byte(`svg
<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg"><g id="box"><title>Layer 1</title><path fill="#ffffff" stroke="#000000" d="m0,25l50,-25l50,25l-50,25l-50,-25z" id="diam" fill-opacity="0.04"/><path fill="#ffffff" stroke="#000000" d="m0,75l50,-25l50,25l-50,25l-50,-25z" fill-opacity="0.04" id="svg_1"/><line id="svg_3" y2="75" x2="0" y1="25" x1="0" fill-opacity="0.04" stroke="#000000" fill="none"/><line id="svg_4" y2="100" x2="50" y1="50" x1="50" fill-opacity="0.04" stroke="#000000" fill="none"/><line id="svg_6" y2="75" x2="100" y1="25" x1="100" fill-opacity="0.04" stroke="#000000" fill="none"/><line id="svg_7" y2="50" x2="50" y1="0" x1="50" fill-opacity="0.04" stroke="#000000" fill="none"/></g></svg>`)

	var writeCount int
	var err error
	var gkErr *gkerr.GkErrDef

	writeCount, err = ws.Write([]byte(svgResponse))
	if err != nil {
		gkErr = gkerr.GenGkErr("ws.Write", err, ERROR_ID_WEBSOCKET_WRITE)
		gklog.LogGkErr("ws.Read", gkErr)
		return gkErr
	}
	if writeCount != len(svgResponse) {
		gkErr = gkerr.GenGkErr("write count short len", err, ERROR_ID_WEBSOCKET_SHORT_WRITE)
		gklog.LogGkErr("ws.Read", gkErr)
		return gkErr
	}

	return nil
}
*/

func getCommandAndData(rawWebsocketData []byte) (string, []byte, *gkerr.GkErrDef) {
	var command []byte
	var websocketData []byte
	var gkErr *gkerr.GkErrDef

	index := bytes.IndexByte(rawWebsocketData, '\n')
	if index == -1 {
		gkErr = gkerr.GenGkErr("getCommandAndData", nil, ERROR_ID_INVALID_WEBSOCKET_DATA)
		return "", nil, gkErr
	}
	command = rawWebsocketData[:index]
	if (index + 1) < len(rawWebsocketData) {
		websocketData = rawWebsocketData[index+1:]
	}

	return string(command), websocketData, nil
}

func getSvgData(svgDir string, svgFile string) ([]byte, *gkerr.GkErrDef) {
	var file *os.File
	var err error
	var gkErr *gkerr.GkErrDef
	var svgData []byte = make([]byte, 0, 0)

	file, err = os.Open(svgDir + string(os.PathSeparator) + svgFile)
	if err != nil {
		gkErr = gkerr.GenGkErr("os.Open", err, ERROR_ID_OPEN_SVG_FILE)
		return nil, gkErr
	}

	defer file.Close()

	buf := make([]byte, 1024, 1024)
	var readCount int

	for {
		readCount, err = file.Read(buf)
		if readCount > 0 {
			svgData = append(svgData, buf[0:readCount]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			gkErr = gkerr.GenGkErr("file.Read", err, ERROR_ID_READ_SVG_FILE)
			return nil, gkErr
		}
	}

	return svgData, nil
}
