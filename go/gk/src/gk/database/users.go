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
	"time"
)

import (
	"gk/gkerr"
)

type DbUserDef struct {
	id                  int32
	UserName            string
	PasswordHash        string
	PasswordSalt        string
	Email               string
	accountCreationDate time.Time
	lastLoginDate       time.Time
}

type DbContextUserDef struct {
	id            int32
	LastPositionX int16
	LastPositionY int16
	LastPositionZ int16
	LastPod       string
}

// return error if row not found
func (gkDbCon *GkDbConDef) GetUser(userName string) (*DbUserDef, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error
	var dbUser *DbUserDef = new(DbUserDef)

	stmt, err = gkDbCon.sqlDb.Prepare("select id, user_name, password_hash, password_salt, email, account_creation_date, last_login_date from users where user_name = $1")
	if err != nil {
		return nil, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query(userName)
	if err != nil {
		return nil, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	if rows.Next() {
		err = rows.Scan(&dbUser.id, &dbUser.UserName, &dbUser.PasswordHash, &dbUser.PasswordSalt, &dbUser.Email, &dbUser.accountCreationDate, &dbUser.lastLoginDate)
		if err != nil {
			return nil, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
	} else {
		return nil, gkerr.GenGkErr("select users", nil, ERROR_ID_NO_ROWS_FOUND)
	}

	return dbUser, nil
}

// no error return if row is not found
func (gkDbCon *GkDbConDef) GetContextUser(userName string) (*DbContextUserDef, bool, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error
	var dbContextUser *DbContextUserDef = new(DbContextUserDef)
	var rowFound bool

	stmt, err = gkDbCon.sqlDb.Prepare("select context_users.id, context_users.last_position_x, context_users.last_position_y, context_users.last_position_z, context_users.last_pod from users, context_users where users.id = context_users.id and users.user_name = $1")
	if err != nil {
		return nil, false, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query(userName)
	if err != nil {
		return nil, false, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	if rows.Next() {
		err = rows.Scan(&dbContextUser.id, &dbContextUser.LastPositionX, &dbContextUser.LastPositionY, &dbContextUser.LastPositionZ, &dbContextUser.LastPod)
		if err != nil {
			return nil, false, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
		rowFound = true
	} else {
		rowFound = false
	}

	return dbContextUser, rowFound, nil
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

	stmt, err = gkDbCon.sqlDb.Prepare("insert into users (id, user_name, password_hash, password_salt, email, account_creation_date, last_login_date) values ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var accountCreationDate time.Time = time.Now()
	var lastLoginDate time.Time = time.Now()

	_, err = stmt.Exec(id, userName, passwordHash, passwordSalt, email, accountCreationDate, lastLoginDate)
	if err != nil {
		if isUniqueViolation(err) {
			return gkerr.GenGkErr("stmt.Exec unique violation", err, ERROR_ID_UNIQUE_VIOLATION)
		}
		return gkerr.GenGkErr("stmt.Exec"+getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) UpdateUserLoginDate(userName string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	//var gkErr *gkerr.GkErrDef

	stmt, err = gkDbCon.sqlDb.Prepare("update users set last_login_date = $1 where user_name = $2")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var lastLoginDate time.Time = time.Now()

	_, err = stmt.Exec(lastLoginDate, userName)
	if err != nil {
		return gkerr.GenGkErr("stmt.Exec"+getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) ChangePassword(userName string, passwordHash string, passwordSalt string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	//var gkErr *gkerr.GkErrDef

	stmt, err = gkDbCon.sqlDb.Prepare("update users set password_hash = $1, password_salt = $2 where user_name = $3")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	_, err = stmt.Exec(passwordHash, passwordSalt, userName)
	if err != nil {
		return gkerr.GenGkErr("stmt.Exec"+getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) getNextUsersId() (int64, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("select nextval('users_seq')")
	if err != nil {
		return 0, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		return 0, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	var userId int64

	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return 0, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
	} else {
		return 0, gkerr.GenGkErr("select users", nil, ERROR_ID_NO_ROWS_FOUND)
	}

	return userId, nil
}
