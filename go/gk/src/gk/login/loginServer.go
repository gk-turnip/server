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
	"flag"
	"fmt"
	"net/http"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

func LoginServerStart() {

	var fileName *string = flag.String("config", "", "config file name")
	var loginConfig loginConfigDef
	var gkErr *gkerr.GkErrDef

	flag.Parse()

	if *fileName == "" {
		flag.PrintDefaults()
		return
	}

	loginConfig, gkErr = loadConfigFile(*fileName)
	if gkErr != nil {
		fmt.Printf("error before log setup %s\n", gkErr.String())
		return
	}

	gklog.LogInit(loginConfig.LogDir)
	gkErr = loginConfig.loginInit()
	if gkErr != nil {
		gklog.LogGkErr("loginConfig.loginInit", gkErr)
		return
	}

	address := fmt.Sprintf(":%d", loginConfig.Port)

	http.ListenAndServe(address, &loginConfig)
}
