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
	"encoding/xml"
	"fmt"
	"os"
)

import (
	"gk/gkerr"
)

type loginConfigDef struct {
	XMLName                xml.Name `xml:"config"`
	Port                   int      `xml:"port"`
	LogDir                 string   `xml:"logDir"`
	TemplateDir            string   `xml:"templateDir"`
	LoginWebAddressPrefix  string   `xml:"loginWebAddressPrefix"`
	GameWebAddressPrefix   string   `xml:"gameWebAddressPrefix"`
	GameTokenAddressPrefix string   `xml:"gameTokenAddressPrefix"`
	DatabaseHost           string   `xml:"databaseHost"`
	DatabasePort           int      `xml:"databasePort"`
	DatabaseUserName       string   `xml:"databaseUserName"`
	DatabasePassword       string   `xml:"databasePassword"`
	DatabaseDatabase       string   `xml:"databaseDatabase"`
	ServerFromEmail        string   `xml:"serverFromEmail"`
	EmailServer            string   `xml:"emailServer"`
}

func loadConfigFile(fileName string) (loginConfigDef, *gkerr.GkErrDef) {
	var err error
	var loginConfig loginConfigDef

	var file *os.File
	file, err = os.Open(fileName)
	if err != nil {
		return loginConfig, gkerr.GenGkErr(fmt.Sprintf("os.Open file: %s", fileName), err, ERROR_ID_OPEN_CONFIG)
	}
	defer file.Close()

	err = xml.NewDecoder(file).Decode(&loginConfig)
	if err != nil {
		return loginConfig, gkerr.GenGkErr(fmt.Sprintf("xml.NewDecoder file: %s", fileName), err, ERROR_ID_DECODE_CONFIG)
	}

	return loginConfig, nil
}
