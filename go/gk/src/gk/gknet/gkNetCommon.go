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
package gknet

// check if the request matches the path
// even if the request or path have an extra trailing slash
func RequestMatches(request string, path string) bool {
	if request == path {
		return true
	}

	if len(request) < 1 || len(path) < 1 {
		return false
	}

	if request[len(request) - 1:] == "/" &&
		path[len(path) - 1:] == "/" {
		return false
	}

	if request[len(request) - 1:] != "/" &&
		path[len(path) - 1:] != "/" {
		return false
	}

	if path[len(path) - 1:] != "/" {
		if request[:len(request) - 1] == path {
			return true
		}
	}

	if request[len(request) - 1:] != "/" {
		if path[:len(path) - 1] == request {
			return true
		}
	}

	return false
}

