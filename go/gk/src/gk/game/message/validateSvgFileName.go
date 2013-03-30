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

package message

func validateSvgFileName(fileName string) bool {
	if len(fileName) < 1 || len(fileName) > 20 {
		return false
	}

	for i := 0;i < len(fileName); i++ {
		if (fileName[i] < '0' || fileName[i] > '9') &&
			(fileName[i] < 'a' || fileName[i] > 'z') &&
			(fileName[i] < 'A' || fileName[i] > 'Z') &&
			(fileName[i] != '-') &&
			(fileName[i] != '_') {
			return false
		}
	}
	return true
}

