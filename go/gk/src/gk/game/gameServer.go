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
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"net/http"
	"time"
)

import (
	"gk/game/config"
	"gk/game/field"
	"gk/game/ses"
	"gk/game/ws"
	"gk/gkerr"
	"gk/gklog"
	"gk/gkrand"
)

func GameServerStart() {

	var fileName *string = flag.String("config", "", "config file name")
	var gameConfig *config.GameConfigDef
	var gkErr *gkerr.GkErrDef

	flag.Parse()

	if *fileName == "" {
		flag.PrintDefaults()
		return
	}

	gameConfig, gkErr = config.LoadConfigFile(*fileName)
	if gkErr != nil {
		fmt.Print(gkErr.String())
		return
	}

	gklog.LogInit(gameConfig.LogDir)

	var randContext *gkrand.GkRandContextDef
	var tokenContext *tokenContextDef
	var sessionContext *ses.SessionContextDef
	var httpContext *httpContextDef

	randContext = gkrand.NewGkRandContext()
	tokenContext = NewTokenContext(gameConfig, randContext, sessionContext)
	sessionContext = ses.NewSessionContext(randContext)
	httpContext = NewHttpContext(gameConfig, sessionContext, tokenContext)

	gkErr = httpContext.gameInit()
	if gkErr != nil {
		gklog.LogGkErr("httpContext.gameInit", gkErr)
		return
	}

	gkErr = tokenContext.gameInit()
	if gkErr != nil {
		gklog.LogGkErr("tokenContext.gameInit", gkErr)
		return
	}

	gklog.LogTrace("game server started")

	var wsContext *ws.WsContextDef
	var fieldContext *field.FieldContextDef

	fieldContext = field.NewFieldContext(gameConfig.SvgDir, sessionContext)
	wsContext = ws.NewWsContext(gameConfig, sessionContext, fieldContext)
	ws.SetGlobalWsContext(wsContext)

	go fieldContext.StartFieldHandler()

	httpAddress := fmt.Sprintf(":%d", gameConfig.HttpPort)

	tokenAddress := fmt.Sprintf(":%d", gameConfig.TokenPort)

	var err error

	go func() {
		err = http.ListenAndServe(tokenAddress, tokenContext)
		if err != nil {
			gkErr = gkerr.GenGkErr("http.ListenAndServer token", err, ERROR_ID_TOKEN_SERVER_START)
			gklog.LogGkErr("", gkErr)
			return
		}
		gklog.LogTrace("token listener ended, this is probably bad")
	}()

	go func() {
		err = http.ListenAndServe(httpAddress, httpContext)
		if err != nil {
			gkErr = gkerr.GenGkErr("http.ListenAndServer http", err, ERROR_ID_HTTP_SERVER_START)
			gklog.LogGkErr("", gkErr)
			return
		}
		gklog.LogTrace("http listener ended, this is probably bad")
	}()

	go func() {
		websocketAddress := fmt.Sprintf(":%d", gameConfig.WebsocketPort)
		gklog.LogTrace("starting web socket listener")
		if gameConfig.CertificatePath == "" {
			err = http.ListenAndServe(websocketAddress, websocket.Handler(ws.WebsocketHandler))
		} else {
			err = http.ListenAndServeTLS(websocketAddress, gameConfig.CertificatePath, gameConfig.PrivateKeyPath, websocket.Handler(ws.WebsocketHandler))
		}
		if err != nil {
			gkErr = gkerr.GenGkErr("http.ListenAndServer websocket", err, ERROR_ID_WEBSOCKET_SERVER_START)
			gklog.LogGkErr("", gkErr)
			return
		}
		gklog.LogTrace("websocket listener ended, this is probably bad")
	}()

	// give it time for the servers to start
	time.Sleep(time.Second * 60)
	// wait for all go routines to finish
	select {}
	gklog.LogTrace("game server ended")
}
