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

type DbUserPrefDef struct {
	Id    int32
	PrefName string
	PrefValue string
}

func (gkDbCon *GkDbConDef) GetUserPrefsList(userName string) ([]DbUserPrefDef, *gkerr.GkErrDef) {
	var dbUser *DbUserDef
	var gkErr *gkerr.GkErrDef

	var stmt *sql.Stmt
	var err error
	var dbUserPrefsList []DbUserPrefDef = make([]DbUserPrefDef, 0, 0)

	dbUser, gkErr = gkDbCon.GetUser(userName)
	if gkErr != nil {
		return nil, gkErr
	}

	stmt, err = gkDbCon.sqlDb.Prepare("select user_id, pref_name, pref_value from user_prefs where user_id = $1")
	if err != nil {
		return nil, gkerr.GenGkErr("sql.Prepare " + getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	defer stmt.Close()

	var rows *sql.Rows

	rows, err = stmt.Query(dbUser.id)
	if err != nil {
		return nil, gkerr.GenGkErr("stmt.Query " + getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	defer rows.Close()

	for rows.Next() {
		var dbUserPref DbUserPrefDef
		err = rows.Scan(&dbUserPref.Id, &dbUserPref.PrefName, &dbUserPref.PrefValue)
		if err != nil {
			return nil, gkerr.GenGkErr("rows.Scan " + getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
		dbUserPrefsList = append(dbUserPrefsList, dbUserPref)
	}

	return dbUserPrefsList, nil
}

func (gkDbCon *GkDbConDef) SetUserPref(userName string, prefName string, prefValue string) *gkerr.GkErrDef {
	var dbUser *DbUserDef
	var gkErr *gkerr.GkErrDef

	dbUser, gkErr = gkDbCon.GetUser(userName)
	if gkErr != nil {
		return gkErr
	}

	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("select pref_value from user_prefs where user_id = $1 and pref_name = $2")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare " + getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	defer stmt.Close()

	var rows *sql.Rows

	rows, err = stmt.Query(dbUser.id, prefName)
	if err != nil {
		return  gkerr.GenGkErr("stmt.Query" + getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	defer rows.Close()

	if rows.Next() {
		gkErr = gkDbCon.UpdateUserPref(dbUser.id, prefName, prefValue)
		if gkErr != nil {
			return gkErr
		}
	} else {
		gkErr = gkDbCon.InsertUserPref(dbUser.id, prefName, prefValue)
		if gkErr != nil {
			return gkErr
		}
	}

	return nil
}

func (gkDbCon *GkDbConDef) UpdateUserPref(id int32, prefName string, prefValue string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("update user_prefs set pref_value = $1 where user_id = $2 and pref_name = $3")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare " + getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	defer stmt.Close()

	_, err = stmt.Exec(prefValue, id, prefName)
	if err != nil {
		return gkerr.GenGkErr("stmt.Exec " + getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) InsertUserPref(id int32, prefName string, prefValue string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("insert into user_prefs (user_id, pref_value, pref_name) values ($1, $2, $3)")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare " + getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	defer stmt.Close()

	_, err = stmt.Exec(id, prefValue, prefName)
	if err != nil {
		return gkerr.GenGkErr("stmt.Exec " + getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

