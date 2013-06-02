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

const (
	ERROR_ID_UNKNOWN_WEBSOCKET_COMMAND = 0xb000 + iota
	ERROR_ID_JSON_UNMARSHAL
	ERROR_ID_COULD_NOT_GET_WEBSOCKET_CONNECTION_CONTEXT
	ERROR_ID_MESSAGE_TO_CLIENT_QUEUE_OVERFLOW
	ERROR_ID_OPENING_ALREADY_OPEN_SESSION
	ERROR_ID_CLOSING_ALREADY_CLOSED_SESSION
	ERROR_ID_SVG_DIR_OPEN
	ERROR_ID_SVG_DIR_READ
	ERROR_ID_COULD_NOT_FIND_OBJECT_TO_MOVE
	ERROR_ID_OPEN_TERRAIN_MAP
	ERROR_ID_READ_TERRAIN_MAP
	ERROR_ID_INVALID_POD_ID
)
