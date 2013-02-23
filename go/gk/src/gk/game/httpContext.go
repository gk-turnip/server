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
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

import (
	"gk/gkerr"
	"gk/gklog"
	"gk/gknet"
	"gk/gktmpl"
	"gk/game/ses"
	"gk/game/config"
)

const _methodGet = "GET"
const _methodPost = "POST"

const _gameRequest = "/gk/gameServer"
const _websocketRequest = "/ws"

const _actParam = "act"
const _submitParam = "submit"
const _registerParam = "register"
const _userNameParam = "userName"
const _passwordParam = "password"
const _emailParam = "email"
const _tokenParam = "token"

var _gameTemplate *gktmpl.TemplateDef
var _gameTemplateName string = "game"

type httpContextDef struct {
	sessionContext *ses.SessionContextDef
	gameConfig *config.GameConfigDef
	tokenContext *tokenContextDef
}

type gameDataDef struct {
	Title string
	WebAddressPrefix string
	WebsocketAddressPrefix string
	AudioAddressPrefix string
	WebsocketPath string
	SessionId string
}

var _errorTemplate *gktmpl.TemplateDef
var _errorTemplateName string = "error"

type errorDataDef struct {
	Title   string
	Message string
	WebAddressPrefix string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewHttpContext(gameConfig *config.GameConfigDef, sessionContext *ses.SessionContextDef, tokenContext *tokenContextDef) *httpContextDef {
	var httpContext *httpContextDef = new(httpContextDef)

	httpContext.gameConfig = gameConfig
	httpContext.sessionContext = sessionContext
	httpContext.tokenContext = tokenContext

	return httpContext
}

func (httpContext *httpContextDef) gameInit() *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	_gameTemplate, gkErr = gktmpl.NewTemplate(httpContext.gameConfig.TemplateDir, _gameTemplateName)
	if gkErr != nil {
		return gkErr
	}

	_errorTemplate, gkErr = gktmpl.NewTemplate(httpContext.gameConfig.TemplateDir, _errorTemplateName)
	if gkErr != nil {
		return gkErr
	}

	return nil
}

func (httpContext *httpContextDef) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if _gameTemplate == nil {
		gklog.LogError("missing call to gameInit")
	}

	path := req.URL.Path

	gklog.LogTrace(req.Method)
	gklog.LogTrace(path)

	if req.Method == _methodGet || req.Method == _methodPost {
		if gknet.RequestMatches(path, _gameRequest) {
			httpContext.handleGameRequest(res, req)
		} else {
			http.NotFound(res, req)
		}
	} else {
		http.NotFound(res, req)
	}
}

func (httpContext *httpContextDef) handleGameRequest(res http.ResponseWriter, req *http.Request) {
	var act string

	req.ParseForm()

	act = req.Form.Get(_actParam)

	switch act {
	case "":
		httpContext.handleGameInitial(res, req)
		return
	default:
		gklog.LogError("unknown act")
		httpContext.redirectToError("unknown act", res, req)
		return
	}
}

func (httpContext *httpContextDef) handleGameInitial(res http.ResponseWriter, req *http.Request) {
	var gameData gameDataDef
	var gkErr *gkerr.GkErrDef
	var singleSession *ses.SingleSessionDef
	var token string

	token = req.Form.Get(_tokenParam)
gklog.LogTrace("got token: " + token)
	var userName string
	userName = httpContext.tokenContext.getUserFromToken(token)
gklog.LogTrace("got username: " + userName)

	if len(userName) < 3 {
		httpContext.redirectToError("not valid token", res, req)
		return
	}

	singleSession = httpContext.sessionContext.NewSingleSession(userName, req.RemoteAddr)

	gameData.Title = "game"
	gameData.WebAddressPrefix = httpContext.gameConfig.WebAddressPrefix
	gameData.WebsocketAddressPrefix = httpContext.gameConfig.WebsocketAddressPrefix
	gameData.AudioAddressPrefix = httpContext.gameConfig.AudioAddressPrefix
	gameData.WebsocketPath = httpContext.gameConfig.WebsocketPath
	gameData.SessionId = singleSession.GetSessionId()

	gkErr = _gameTemplate.Build(gameData)
	if gkErr != nil {
		gklog.LogGkErr("_gameTemplate.Build", gkErr)
		httpContext.redirectToError("_gameTemplate.Build", res, req)
		return
	}

	gkErr = _gameTemplate.Send(res, req)
	if gkErr != nil {
		gklog.LogGkErr("_gameTemplate.Send", gkErr)
		return
	}
}

func genErrorMarker() template.HTML {
	return template.HTML("<span class=\"errorMarker\">*</span>")
}

func (httpContext *httpContextDef) redirectToError(message string, res http.ResponseWriter, req *http.Request) {
	var errorData errorDataDef
	var gkErr *gkerr.GkErrDef

	errorData.Title = "Error"
	errorData.Message = message
	errorData.WebAddressPrefix = httpContext.gameConfig.WebAddressPrefix

	gkErr = _errorTemplate.Build(errorData)
	if gkErr != nil {
		gklog.LogGkErr("_errorTemplate.Build", gkErr)
	}

	gkErr = _errorTemplate.Send(res, req)
	if gkErr != nil {
		gklog.LogGkErr("_errorTemplate.Send", gkErr)
	}
}
