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
	"pq/pq"
	"gk/gkerr"
	"gk/gklog"
)

type GkDbConDef struct {
	sqlDb *sql.DB
}

const pgErrorUniqueViolation = "23505"

func NewGkDbCon(userName string, password string, host string, port int, database string) (*GkDbConDef, *gkerr.GkErrDef) {
	var gkDbCon *GkDbConDef = new(GkDbConDef)
	var err error
	var connectionString string

	connectionString = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",userName, password, host, port, database)

	gklog.LogTrace(connectionString)

	gkDbCon.sqlDb, err = sql.Open("postgres",connectionString)

	if err != nil {
		return nil, gkerr.GenGkErr("sql.Open " + getDatabaseErrorMessage(err),err,ERROR_ID_DATABASE_CONNECTION)
	}

	return gkDbCon, nil
}

func (gkDbCon *GkDbConDef) Close() {
	gkDbCon.sqlDb.Close()
}

func getDatabaseErrorMessage(err error) string {
	result := "unknown"

	if err != nil {
		var ok bool
		var pge pq.PGError
		pge, ok = err.(pq.PGError)
		if ok {
			var r, l, m, c, s, f string
			r = pge.Get('R')
			l = pge.Get('L')
			m = pge.Get('M')
			c = pge.Get('C')
			s = pge.Get('S')
			f = pge.Get('F')

			result = fmt.Sprintf("postgres error r: %s l: %s m: %s c: %s s: %s f: %s",r,l,m,c,s,f)
		}
	}

	return result
}

func isUniqueViolation(err error) bool {
	if err != nil {
		var ok bool
		var pge pq.PGError
		pge, ok = err.(pq.PGError)
		if ok {
			if pge.Get('C') == pgErrorUniqueViolation {
				return true
			}
		}
	}
	return false
}

