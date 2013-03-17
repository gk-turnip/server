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

package gkjson

func JsonEscape(message string) string {
	result := make([]byte,0,len(message) + 4)

	for i := 0;i < len(message); i++ {
		switch message[i] {
		case '~':
			result = append(result, '.')
		case '"':
			result = append(result, '\\', '"')
		default:
			result = append(result, message[i])
		}
	}

	return string(result)
}

