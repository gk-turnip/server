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

package login

import (
	//	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"
)

import (
	"gk/database"
	"gk/gkerr"
	"gk/gklog"
	"gk/gktmpl"
	"gk/sec"
)

const _methodGet = "GET"
const _methodPost = "POST"
const _loginRequest = "/gk/loginServer/"
const _gameServer = "/gk/gameServer"

const _actParam = "act"
const _submitParam = "submit"
const _registerParam = "register"
const _userNameParam = "userName"
const _passwordParam = "password"
const _emailParam = "email"

var _loginTemplate *gktmpl.TemplateDef

var _registerTemplate *gktmpl.TemplateDef

var _errorTemplate *gktmpl.TemplateDef

type loginDataDef struct {
	Title            string
	ErrorList        []string
	UserName         string
	UserNameError    template.HTML
	PasswordError    template.HTML
	WebAddressPrefix string
}

type registerDataDef struct {
	Title            string
	ErrorList        []string
	UserName         string
	UserNameError    template.HTML
	PasswordError    template.HTML
	Email            string
	EmailError       template.HTML
	WebAddressPrefix string
}

type errorDataDef struct {
	Title   string
	Message string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (loginConfig *loginConfigDef) loginInit() *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	var fileNames []string

	fileNames = []string{"main", "head", "error_list", "login"}
	_loginTemplate, gkErr = gktmpl.NewTemplate(loginConfig.TemplateDir, fileNames)
	if gkErr != nil {
		return gkErr
	}

	fileNames = []string{"main", "head", "error_list", "register"}
	_registerTemplate, gkErr = gktmpl.NewTemplate(loginConfig.TemplateDir, fileNames)
	if gkErr != nil {
		return gkErr
	}

	fileNames = []string{"main", "head", "error_list", "error"}
	_errorTemplate, gkErr = gktmpl.NewTemplate(loginConfig.TemplateDir, fileNames)
	if gkErr != nil {
		return gkErr
	}

	return nil
}

func (loginConfig *loginConfigDef) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if _loginTemplate == nil {
		gklog.LogError("missing call to loginInit")
	}

	path := req.URL.Path

	gklog.LogTrace(req.Method)
	gklog.LogTrace(path)

	if req.Method == _methodGet || req.Method == _methodPost {
		if requestMatch(path, _loginRequest) {
			handleLogin(loginConfig, res, req)
		} else {
			http.NotFound(res, req)
		}
	} else {
		http.NotFound(res, req)
	}

}

func handleLogin(loginConfig *loginConfigDef, res http.ResponseWriter, req *http.Request) {
	var act string
	var userName string
	var password string
	var email string

	req.ParseForm()

	act = req.Form.Get(_actParam)
	userName = req.Form.Get(_userNameParam)
	password = req.Form.Get(_passwordParam)
	email = req.Form.Get(_emailParam)

	// for security
	// sleep between 10 and 19 milliseconds
	randMilliseconds := int(rand.Int31n(1000)) + 1000
	time.Sleep(time.Nanosecond * 10000 * time.Duration(randMilliseconds))

	gklog.LogTrace("act: " + act)

	var submit string
	var register string

	submit = req.Form.Get(_submitParam)
	register = req.Form.Get(_registerParam)

	if submit != "" {
		handleLoginLogin(loginConfig, res, req, userName, password)
		return
	}

	if register != "" {
		handleLoginRegisterInitial(loginConfig, res, req)
		return
	}

	switch act {
	case "":
		handleLoginInitial(loginConfig, res, req)
		return
	case "login":
		if userName == "" {
			handleLoginInitial(loginConfig, res, req)
			return
		}
		handleLoginLogin(loginConfig, res, req, userName, password)
		return
	case "register":
		handleLoginRegister(loginConfig, res, req, userName, password, email)
	default:
		gklog.LogError("unknown act")
		redirectToError("unknown act", res, req)
		return
	}
}

func handleLoginInitial(loginConfig *loginConfigDef, res http.ResponseWriter, req *http.Request) {
	var loginData loginDataDef
	var gkErr *gkerr.GkErrDef

	loginData.Title = "login"
	loginData.WebAddressPrefix = loginConfig.WebAddressPrefix

	gkErr = _loginTemplate.Build(loginData)
	if gkErr != nil {
		gklog.LogGkErr("_loginTemplate.Build", gkErr)
		redirectToError("_loginTemplate.Build", res, req)
		return
	}

	gkErr = _loginTemplate.Send(res, req)
	if gkErr != nil {
		gklog.LogGkErr("_loginTemplate.Send", gkErr)
		return
	}
}

func handleLoginLogin(loginConfig *loginConfigDef, res http.ResponseWriter, req *http.Request, userName string, password string) {
	var loginData loginDataDef
	var gkErr *gkerr.GkErrDef
	var gotError bool

	loginData.Title = "login"
	loginData.UserName = userName
	loginData.WebAddressPrefix = loginConfig.WebAddressPrefix

	if loginData.UserName == "" {
		loginData.ErrorList = append(loginData.ErrorList, "user name cannot be blank")
		loginData.UserNameError = genErrorMarker()
		gotError = true
	}

	if password == "" {
		loginData.ErrorList = append(loginData.ErrorList, "password cannot be blank")
		loginData.PasswordError = genErrorMarker()
		gotError = true
	}

	var passwordHashFromDatabase string
	var passwordHashFromUser []byte
	var passwordSalt string

	if !gotError {
		var gkDbCon *database.GkDbConDef

		gkDbCon, gkErr = database.NewGkDbCon(loginConfig.DatabaseUserName, loginConfig.DatabasePassword, loginConfig.DatabaseHost, loginConfig.DatabasePort, loginConfig.DatabaseDatabase)
		if gkErr != nil {
			gklog.LogGkErr("_registerTemplate.Build", gkErr)
			redirectToError("_registerTemplate.Build", res, req)
			return
		}

		defer gkDbCon.Close()

		passwordHashFromDatabase, passwordSalt, gkErr = gkDbCon.GetPasswordHashAndSalt(loginData.UserName)

		if gkErr != nil {
			if gkErr.GetErrorId() == database.ERROR_ID_NO_ROWS_FOUND {
				password = "one two three"
				passwordSalt = "abc123QWE."
				// make it take the same amount of time
				// between no user and invalid password
				passwordHashFromUser = sec.GenPasswordHashSlow([]byte(password), []byte(passwordSalt))
				loginData.ErrorList = append(loginData.ErrorList, "invalid username/password")
				loginData.UserNameError = genErrorMarker()
				loginData.PasswordError = genErrorMarker()
				gotError = true
			} else {
				gklog.LogGkErr("gkDbCon.GetPasswordHashAndSalt", gkErr)
				redirectToError("gkDbCon.GetPasswordhashAndSalt", res, req)
				return
			}
		}
	}

	if !gotError {
		passwordHashFromUser = sec.GenPasswordHashSlow([]byte(password), []byte(passwordSalt))

		if passwordHashFromDatabase != string(passwordHashFromUser) {
			loginData.ErrorList = append(loginData.ErrorList, "invalid username/password")
			loginData.UserNameError = genErrorMarker()
			loginData.PasswordError = genErrorMarker()
			gotError = true
		}
	}

	if gotError {
		// for security
		// sleep between 100 and 190 milliseconds
		randMilliseconds := int(rand.Int31n(10)) + 10
		time.Sleep(time.Nanosecond * 10000000 * time.Duration(randMilliseconds))

		gkErr = _loginTemplate.Build(loginData)
		if gkErr != nil {
			gklog.LogGkErr("_loginTemplate.Build", gkErr)
			redirectToError("_loginTemplate.Build", res, req)
			return
		}

		gkErr = _loginTemplate.Send(res, req)
		if gkErr != nil {
			gklog.LogGkErr("_loginTemplate.Send", gkErr)
			return
		}
	} else {
		http.Redirect(res, req, loginConfig.WebAddressPrefix+_gameServer, http.StatusFound)
	}
}

func handleLoginRegisterInitial(loginConfig *loginConfigDef, res http.ResponseWriter, req *http.Request) {
	var registerData registerDataDef
	var gkErr *gkerr.GkErrDef

	registerData.Title = "register"
	registerData.WebAddressPrefix = loginConfig.WebAddressPrefix

	gkErr = _registerTemplate.Build(registerData)
	if gkErr != nil {
		gklog.LogGkErr("_registerTemplate.Build", gkErr)
		redirectToError("_registerTemplate.Build", res, req)
		return
	}

	gkErr = _registerTemplate.Send(res, req)
	if gkErr != nil {
		gklog.LogGkErr("_registerTemplate.send", gkErr)
	}
}

func handleLoginRegister(loginConfig *loginConfigDef, res http.ResponseWriter, req *http.Request, userName string, password string, email string) {
	var registerData registerDataDef
	var gkErr *gkerr.GkErrDef
	var err error

	registerData.Title = "register"
	registerData.WebAddressPrefix = loginConfig.WebAddressPrefix
	registerData.UserName = userName
	registerData.Email = email
	registerData.ErrorList = make([]string, 0, 0)

	var gotError bool

	if userName == "" {
		registerData.ErrorList = append(registerData.ErrorList, "user name cannot be blank")
		registerData.UserNameError = genErrorMarker()
		gotError = true
	}
	if password == "" {
		registerData.ErrorList = append(registerData.ErrorList, "password cannot be blank")
		registerData.PasswordError = genErrorMarker()
		gotError = true
	}
	if email == "" {
		registerData.ErrorList = append(registerData.ErrorList, "email cannot be blank")
		registerData.EmailError = genErrorMarker()
		gotError = true
	}

	if !gotError {
		var gkDbCon *database.GkDbConDef

		gkDbCon, gkErr = database.NewGkDbCon(loginConfig.DatabaseUserName, loginConfig.DatabasePassword, loginConfig.DatabaseHost, loginConfig.DatabasePort, loginConfig.DatabaseDatabase)
		if gkErr != nil {
			gklog.LogGkErr("_registerTemplate.Build", gkErr)
			redirectToError("_registerTemplate.Build", res, req)
			return
		}

		defer gkDbCon.Close()

		var passwordHash, passwordSalt []byte

		passwordSalt, err = sec.GenSalt()
		if err != nil {
			gkErr = gkerr.GenGkErr("sec.GenSalt", err, ERROR_ID_GEN_SALT)
			gklog.LogGkErr("sec.GenSalt", gkErr)
			redirectToError("sec.GenSalt", res, req)
		}

		passwordHash = sec.GenPasswordHashSlow([]byte(password), passwordSalt)

		gkErr = gkDbCon.AddNewUser(
			registerData.UserName,
			string(passwordHash),
			string(passwordSalt),
			email)

		if gkErr != nil {
			if gkErr.GetErrorId() == database.ERROR_ID_UNIQUE_VIOLATION {
				registerData.ErrorList = append(registerData.ErrorList, "user name already in use")
				registerData.UserNameError = genErrorMarker()
				gotError = true
			} else {
				gklog.LogGkErr("gbDbCon.AddNewUser", gkErr)
				redirectToError("gbDbCon.AddNewUser", res, req)
				return
			}
		}
	}

	if gotError {
		gkErr = _registerTemplate.Build(registerData)
		if gkErr != nil {
			gklog.LogGkErr("_registerTemplate.Build", gkErr)
			redirectToError("_registerTemplate.Build", res, req)
			return
		}

		gkErr = _registerTemplate.Send(res, req)
		if gkErr != nil {
			gklog.LogGkErr("_registerTemplate.send", gkErr)
		}
	} else {
		http.Redirect(res, req, loginConfig.WebAddressPrefix+_gameServer, http.StatusFound)
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

func requestMatch(path string, request string) bool {
	if path == request {
		return true
	}
	if (path + "/") == request {
		return true
	}

	return false
}
