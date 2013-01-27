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

package game

const (
	ERROR_ID_OPEN_CONFIG = 0x40000000 + iota
	ERROR_ID_DECODE_CONFIG
	ERROR_ID_HTTP_SERVER_START
	ERROR_ID_WEBSOCKET_SERVER_START
	ERROR_ID_WEBSOCKET_SEND
	ERROR_ID_UNKNOWN_WEBSOCKET_INPUT
	ERROR_ID_UNKNOWN_WEBSOCKET_COMMAND
	ERROR_ID_JSON_UNMARSHAL
	ERROR_ID_DUPLICATE_WEBSOCKET_ID
	ERROR_ID_WEBSOCKET_RECEIVE
)
