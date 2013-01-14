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

import (
	"bytes"
)

func fixSvgData(svgData []byte) []byte {

	var result []byte
	result = trimBetweenMarkers(svgData, "<?", "?>")
	result = trimBetweenMarkers(result, "<!--", "-->")
	result = trimCrLf(result)

	return result
}

func trimBetweenMarkers(data []byte, start string, end string) []byte {

	result := make([]byte, 0, len(data))

	var curSlice []byte
	curSlice = data
	var index1 int
	var index2 int

	for {
		index1 = bytes.Index(curSlice, []byte(start))
		if index1 == -1 {
			result = append(result, curSlice...)
			break
		}

		index2 = bytes.Index(curSlice[index1+len(start):], []byte(end))
		if index2 == -1 {
			result = append(result, curSlice...)
			break
		}
		index2 += index1 + len(start)

		result = append(result, curSlice[:index1]...)

		if (index2 + len(end)) >= len(curSlice) {
			break
		}
		curSlice = curSlice[index2+len(end):]
	}

	return result
}

func trimCrLf(data []byte) []byte {
	result := make([]byte, 0, len(data))

	var curSlice []byte
	var index int
	var index1 int
	var index2 int

	curSlice = data
	for {
		index1 = bytes.IndexByte(curSlice, '\n')
		index2 = bytes.IndexByte(curSlice, '\r')
		if index1 == -1 && index2 == -1 {
			result = append(result, curSlice...)
			break
		}
		index = index1

		if index1 == -1 {
			index = index2
		} else {
			if index2 != -1 && index2 < index1 {
				index = index2
			}
		}

		result = append(result, curSlice[:index]...)

		if index >= len(curSlice) {
			break
		}
		curSlice = curSlice[index+1:]
	}

	return result
}
