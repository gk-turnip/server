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

package message

import (
	"os"
	"bytes"
	"text/template"
)

import (
	"gk/gkerr"
	"gk/gkcommon"
	"gk/game/iso"
)

type MessageToClientDef struct {
	Command string
	JsonData []byte
	Data []byte
}

type SvgJsonDataDef struct {
	Id string
	IsoXYZ iso.IsoXYZDef
}

func (messageToClient *MessageToClientDef) BuildSvgMessageToClient(svgDir string, command string, fileName string, svgJsonData *SvgJsonDataDef) *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	messageToClient.Command = command

	jsonFileName := svgDir + string(os.PathSeparator) + fileName + ".json"
	svgFileName := svgDir + string(os.PathSeparator) + fileName+ ".svg"

	messageToClient.JsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
	if gkErr != nil {
		return gkErr
	}
	messageToClient.Data, gkErr = gkcommon.GetFileContents(svgFileName)
	if gkErr != nil {
		return gkErr
	}
	messageToClient.Data = fixSvgData(messageToClient.Data)

	if svgJsonData != nil {
		messageToClient.JsonData, gkErr = templateTranslateJsonData(messageToClient.JsonData, svgJsonData)
		if gkErr != nil {
			return gkErr
		}
	}

	return nil
}

func templateTranslateJsonData(inputData []byte, svgJsonData *SvgJsonDataDef) ([]byte, *gkerr.GkErrDef) {
	var tmpl *template.Template
	var err error
	var gkErr *gkerr.GkErrDef

	tmpl, err = template.New("fieldObject").Parse(string(inputData))
	if err != nil {
		gkErr = gkerr.GenGkErr("template.New.Parse", err, ERROR_ID_TEMPLATE_PARSE)
		return nil, gkErr
	}

	result := make([]byte, 0, 0)
	var writer *bytes.Buffer = bytes.NewBuffer(result)

	err = tmpl.Execute(writer, svgJsonData)
	if err != nil {
		gkErr = gkerr.GenGkErr("tmpl.Execute", err, ERROR_ID_TEMPLATE_EXECUTE)
		return nil, gkErr
	}

	return writer.Bytes(), nil
}

