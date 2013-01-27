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
	"encoding/json"
	"os"
)

import (
	"gk/gkcommon"
	"gk/gkerr"
)

type getSvgDef struct {
	SvgName string
}

func doGetSvgReq(websocketReq *websocketReqDef) (*websocketResDef, *gkerr.GkErrDef) {
	var websocketRes websocketResDef
	var getSvg getSvgDef
	var gkErr *gkerr.GkErrDef
	var err error

	err = json.Unmarshal(websocketReq.jsonData, &getSvg)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return nil, gkErr
	}

	websocketRes.command = _getSvgRes
	jsonFileName := websocketReq.gameConfig.SvgDir + string(os.PathSeparator) + getSvg.SvgName + ".json"
	svgFileName := websocketReq.gameConfig.SvgDir + string(os.PathSeparator) + getSvg.SvgName + ".svg"

	websocketRes.jsonData, gkErr = gkcommon.GetFileContents(jsonFileName)
	if gkErr != nil {
		return nil, gkErr
	}
	websocketRes.data, gkErr = gkcommon.GetFileContents(svgFileName)
	if gkErr != nil {
		return nil, gkErr
	}
	websocketRes.data = fixSvgData(websocketRes.data)

	return &websocketRes, nil
}
