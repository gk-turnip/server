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
	"gk/gktmpl"
	"gk/gknet"
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

var _gameTemplate *gktmpl.TemplateDef
var _gameTemplateName string = "game"

type gameDataDef struct {
	Title   string
}

var _errorTemplate *gktmpl.TemplateDef

type errorDataDef struct {
	Title   string
	Message string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (httpContext *httpContextDef) gameInit() *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	_gameTemplate, gkErr = gktmpl.NewTemplate(httpContext.gameConfig.TemplateDir, _gameTemplateName)
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
		redirectToError("unknown act", res, req)
		return
	}
}

func (httpContext *httpContextDef) handleGameInitial(res http.ResponseWriter, req *http.Request) {
	var gameData gameDataDef
	var gkErr *gkerr.GkErrDef

	gameData.Title = "game"
	//gameData.WebAddressPrefix = gameConfig.WebAddressPrefix

	gkErr = _gameTemplate.Build(gameData)
	if gkErr != nil {
		gklog.LogGkErr("_gameTemplate.Build", gkErr)
		redirectToError("_gameTemplate.Build", res, req)
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

func redirectToError(message string, res http.ResponseWriter, req *http.Request) {
	var errorData errorDataDef
	var gkErr *gkerr.GkErrDef

	errorData.Title = "Error"
	errorData.Message = message

	gkErr = _errorTemplate.Build(errorData)
	if gkErr != nil {
		gklog.LogGkErr("_errorTemplate.Build", gkErr)
	}

	gkErr = _errorTemplate.Send(res, req)
	if gkErr != nil {
		gklog.LogGkErr("_errorTemplate.Send", gkErr)
	}
}

