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

package gksvg

const (
	ERROR_ID_DATABASE_CONNECTION = 0xd000 + iota
	ERROR_ID_SVG_PARSE
	ERROR_ID_INVALID_SVG
	ERROR_ID_INTERNAL_ID_MAP
	ERROR_ID_HREF_PARSE
	ERROR_ID_STYLE_PARSE
	ERROR_ID_WRITE_SVG
)
