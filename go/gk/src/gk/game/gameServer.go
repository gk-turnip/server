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
	"fmt"
	"flag"
	"time"
	"net/http"
	"code.google.com/p/go.net/websocket"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

type httpContextDef struct {
	gameConfig gameConfigDef
}

type websocketContextDef struct {
	gameConfig gameConfigDef
}

func GameServerStart() {

	var fileName *string = flag.String("config", "", "config file name")
	var gameConfig gameConfigDef
	var gkErr *gkerr.GkErrDef

	flag.Parse()

	if *fileName == "" {
		flag.PrintDefaults()
		return
	}

	gameConfig, gkErr = loadConfigFile(*fileName)
	if gkErr != nil {
		fmt.Print(gkErr.String())
		return
	}

	var httpContext httpContextDef
	httpContext.gameConfig = gameConfig

	gklog.LogInit(gameConfig.LogDir)
	gkErr = httpContext.gameInit()
	if gkErr != nil {
		gklog.LogGkErr("gameConfig.gameInit", gkErr)
		return
	}

	gklog.LogTrace("game server started")

	httpAddress := fmt.Sprintf(":%d", gameConfig.HttpPort)

	var err error

	go func() {
		err = http.ListenAndServe(httpAddress, &httpContext)
		if err != nil {
			gkErr = gkerr.GenGkErr("http.ListenAndServer http", err, ERROR_ID_HTTP_SERVER_START)
			gklog.LogGkErr("http.ListenAndServer", gkErr)
			return
		}
	}()

	websocketAddress := fmt.Sprintf(":%d", gameConfig.WebsocketPort)
	var websocketContext websocketContextDef
	websocketContext.gameConfig = gameConfig

	go func() {
		gklog.LogTrace("starting web socket listener")
		err = http.ListenAndServe(websocketAddress, websocket.Handler(websocketHandler))
		if err != nil {
			gkErr = gkerr.GenGkErr("http.ListenAndServer websocket", err, ERROR_ID_WEBSOCKET_SERVER_START)
			gklog.LogGkErr("http.ListenAndServer", gkErr)
			return
		}
	}()

	// give it time for the servers to start
	time.Sleep(time.Second * 60)
	// wait for all go routines to finish
	select {}
	gklog.LogTrace("game server ended")
}
