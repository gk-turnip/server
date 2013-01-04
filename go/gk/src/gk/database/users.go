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
	"fmt"
	"database/sql"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

func (gkDbCon *GkDbConDef) GetPasswordHashAndSalt(userName string) (string, string, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("select password_hash, password_salt from users where user_name = $1")
	if err != nil {
		return "", "", gkerr.GenGkErr("sql.Prepare" + getDatabaseErrorMessage(err),err,ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query(userName)
	if err != nil {
		return "", "", gkerr.GenGkErr("stmt.Query" + getDatabaseErrorMessage(err),err,ERROR_ID_QUERY)
	}

	var passwordHash, passwordSalt string

	if rows.Next() {
		err = rows.Scan(&passwordHash, &passwordSalt)
		if err != nil {
			return "", "", gkerr.GenGkErr("rows.Scan" + getDatabaseErrorMessage(err),err,ERROR_ID_ROWS_SCAN)
		}
	} else {
		return "", "", gkerr.GenGkErr("select users",nil,ERROR_ID_NO_ROWS_FOUND)
	}

	return passwordHash, passwordSalt, nil

}

func (gkDbCon *GkDbConDef) AddNewUser(userName string, passwordHash string, passwordSalt string, email string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	var id int64
	var gkErr *gkerr.GkErrDef

	id, gkErr = gkDbCon.getNextUsersId()
	if gkErr != nil {
		return gkErr
	}

	stmt, err = gkDbCon.sqlDb.Prepare("insert into users (id, user_name, password_hash, password_salt, email) values ($1, $2, $3, $4, $5)")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare" + getDatabaseErrorMessage(err),err,ERROR_ID_PREPARE)
	}

	gklog.LogTrace(fmt.Sprintf("%s %s",passwordHash, passwordSalt))

	_, err = stmt.Exec(id, userName, passwordHash, passwordSalt, email)
	if err != nil {
		if isUniqueViolation(err) {
			return gkerr.GenGkErr("stmt.Exec unique violation",err,ERROR_ID_UNIQUE_VIOLATION)
		}
		return gkerr.GenGkErr("stmt.Exec" + getDatabaseErrorMessage(err),err,ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) getNextUsersId() (int64, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("select nextval('users_seq')")
	if err != nil {
		return 0, gkerr.GenGkErr("sql.Prepare" + getDatabaseErrorMessage(err),err,ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		return 0, gkerr.GenGkErr("stmt.Query" + getDatabaseErrorMessage(err),err,ERROR_ID_QUERY)
	}

	var userId int64

	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return 0, gkerr.GenGkErr("rows.Scan" + getDatabaseErrorMessage(err),err,ERROR_ID_ROWS_SCAN)
		}
	} else {
		return 0, gkerr.GenGkErr("select users",nil,ERROR_ID_NO_ROWS_FOUND)
	}

	return userId, nil
}

