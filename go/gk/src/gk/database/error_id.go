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

const (
	ERROR_ID_DATABASE_CONNECTION = 0x30000000 + iota
	ERROR_ID_PREPARE
	ERROR_ID_QUERY
	ERROR_ID_ROWS_SCAN
	ERROR_ID_NO_ROWS_FOUND
	ERROR_ID_EXECUTE
	ERROR_ID_UNIQUE_VIOLATION
)

