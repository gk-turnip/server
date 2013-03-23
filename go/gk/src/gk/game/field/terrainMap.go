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
)

import (
	"gk/gkerr"
)

type terrainMapDef struct {
	jsonMapData struct {
		TileList []struct {
			Terrain string
			X       int
			Y       int
			Z       int
		}

		ObjectList []struct {
			Object string
		}
	}
}

func NewTerrainMap(terrainSvgDir string) (*terrainMapDef, *gkerr.GkErrDef) {
	var terrainMap *terrainMapDef = new(terrainMapDef)

	var jsonFileName string = terrainSvgDir + string(os.PathSeparator) + "map_terrain.json"

	var jsonFile *os.File
	var err error
	var gkErr *gkerr.GkErrDef
	var jsonData []byte = make([]byte, 0, 256)

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

	err = json.Unmarshal(jsonData, &terrainMap.jsonMapData)
	if err != nil {
		gkErr = gkerr.GenGkErr("json.Unmarshal", err, ERROR_ID_JSON_UNMARSHAL)
		return nil, gkErr
	}

	return terrainMap, nil
}
