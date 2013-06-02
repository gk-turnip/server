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

package database

import (
	"database/sql"
)

import (
	"gk/gkerr"
)

type DbPodDef struct {
	Id    int32
	Title string
}

func (gkDbCon *GkDbConDef) GetPodsList() ([]DbPodDef, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error
	var dbPodsList []DbPodDef = make([]DbPodDef, 0, 0)

	stmt, err = gkDbCon.sqlDb.Prepare("select id, pod_title from pods")
	if err != nil {
		return nil, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		return nil, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	for rows.Next() {
		var dbPod DbPodDef
		err = rows.Scan(&dbPod.Id, &dbPod.Title)
		if err != nil {
			return nil, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
		dbPodsList = append(dbPodsList, dbPod)
	}

	return dbPodsList, nil
}
