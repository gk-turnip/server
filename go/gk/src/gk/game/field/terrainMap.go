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

package field

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
)

import (
//	"gk/game/persistence"
	"gk/gkerr"
)

type terrainJsonDef struct {
	jsonMapData struct {
		TileList []struct {
			Terrain string
			X       int
			Y       int
			Z       []int
		}

		ObjectList []struct {
			Object string
		}
	}
}

func (fieldContext *FieldContextDef) newTerrainMap(podId int32) (*terrainJsonDef, *gkerr.GkErrDef) {
	//	var terrainMap *terrainMapDef = new(terrainMapDef)
//	var terrainMap map[int32]*terrainJsonDef = make(map[int32]*terrainJsonDef)

	var gkErr *gkerr.GkErrDef

	var jsonFileName string

	jsonFileName = fieldContext.terrainSvgDir + string(os.PathSeparator) + "map_terrain_" + strconv.FormatInt(int64(podId), 10) + ".json"

	var terrainJson *terrainJsonDef

	terrainJson, gkErr = getSingleTerrainMap(jsonFileName)
	if gkErr != nil {
		return nil, gkErr
	}

	return terrainJson, nil
}

func getSingleTerrainMap(jsonFileName string) (*terrainJsonDef, *gkerr.GkErrDef) {
	var terrainJson *terrainJsonDef = new(terrainJsonDef)
	var jsonFile *os.File
	var err error
	var jsonData []byte = make([]byte, 0, 256)
	var gkErr *gkerr.GkErrDef

	jsonFile, err = os.Open(jsonFileName)
	if err != nil {
		gkErr = gkerr.GenGkErr("could not open: "+jsonFileName, err, ERROR_ID_OPEN_TERRAIN_MAP)
		return nil, gkErr
	}
	defer jsonFile.Close()

	var buf []byte
	var readCount int
	for {
		buf = make([]byte, 128, 128)
		readCount, err = jsonFile.Read(buf)
		if readCount > 0 {
			jsonData = append(jsonData, buf[:readCount]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			gkErr = gkerr.GenGkErr("could not read "+jsonFileName, err, ERROR_ID_READ_TERRAIN_MAP)
			return nil, gkErr
		}
	}

	err = json.Unmarshal(jsonData, &terrainJson.jsonMapData)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return nil, gkErr
	}

	return terrainJson, nil
}
