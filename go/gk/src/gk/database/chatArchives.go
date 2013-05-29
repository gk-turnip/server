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

type DbChatArchiveDef struct {
	id                  int32
	userId              int32
	messageCreationDate time.Time
	chatMessage         string
}

// return error if row not found
// this could be improved by saving recent "id" values
func (gkDbCon *GkDbConDef) GetLastChatArchiveEntries(count int) ([]LugChatArchiveDef, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error
	var results []LugChatArchiveDef = make([]LugChatArchiveDef, 0, count)

	stmt, err = gkDbCon.sqlDb.Prepare("select chat_archives.message_creation_date, chat_archives.chat_message, users.user_name from users, chat_archives where users.id = chat_archives.user_id order by chat_archives.message_creation_date desc")
	if err != nil {
		return nil, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		return nil, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	for rows.Next() {
		var lugChatArchive LugChatArchiveDef

		err = rows.Scan(&lugChatArchive.MessageCreationDate, &lugChatArchive.ChatMessage, &lugChatArchive.UserName)
		if err != nil {
			return nil, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
		results = append(results,lugChatArchive)
	}

	return results, nil
}

func (gkDbCon *GkDbConDef) AddNewChatMessage(userName string, chatMessage string) *gkerr.GkErrDef {

	var stmt *sql.Stmt
	var err error

	var id int32
	var gkErr *gkerr.GkErrDef

	id, gkErr = gkDbCon.getNextChatArchivesId()
	if gkErr != nil {
		return gkErr
	}

	stmt, err = gkDbCon.sqlDb.Prepare("insert into chat_archives (id, user_id, message_creation_date, chat_message) values ($1, $2, $3, $4)")
	if err != nil {
		return gkerr.GenGkErr("stmt.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var dbUser *DbUserDef
	dbUser, gkErr = gkDbCon.GetUser(userName)
	if gkErr != nil {
		return gkErr
	}

	var messageCreationDate time.Time = time.Now()

	_, err = stmt.Exec(id, dbUser.id, messageCreationDate, chatMessage)
	if err != nil {
		if isUniqueViolation(err) {
			return gkerr.GenGkErr("stmt.Exec unique violation", err, ERROR_ID_UNIQUE_VIOLATION)
		}
		return gkerr.GenGkErr("stmt.Exec"+getDatabaseErrorMessage(err), err, ERROR_ID_EXECUTE)
	}

	return nil
}

func (gkDbCon *GkDbConDef) getNextChatArchivesId() (int32, *gkerr.GkErrDef) {
	var stmt *sql.Stmt
	var err error

	stmt, err = gkDbCon.sqlDb.Prepare("select nextval('chat_archives_seq')")
	if err != nil {
		return 0, gkerr.GenGkErr("sql.Prepare"+getDatabaseErrorMessage(err), err, ERROR_ID_PREPARE)
	}

	var rows *sql.Rows

	rows, err = stmt.Query()
	if err != nil {
		return 0, gkerr.GenGkErr("stmt.Query"+getDatabaseErrorMessage(err), err, ERROR_ID_QUERY)
	}

	var id int32

	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, gkerr.GenGkErr("rows.Scan"+getDatabaseErrorMessage(err), err, ERROR_ID_ROWS_SCAN)
		}
	} else {
		return 0, gkerr.GenGkErr("select users", nil, ERROR_ID_NO_ROWS_FOUND)
	}

	return id, nil
}
