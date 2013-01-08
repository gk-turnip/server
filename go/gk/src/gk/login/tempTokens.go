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

package login


import (
	"fmt"
	"time"
	"sync"
)

import (
	"gk/gklog"
)

var _tokenExpiry time.Duration = time.Second * 600 // ten minutes

type tokenEntryDef struct {
	userName string
	createdDate time.Time
}

var _tokenMap map[string]tokenEntryDef = make(map[string]tokenEntryDef)
var _tokenMapMutex sync.Mutex

// add a new token / userName to the list of tokens
func AddNewToken(token string, userName string) {
	_tokenMapMutex.Lock()
	defer _tokenMapMutex.Unlock()

	checkTokenExpire()

	var tokenEntry tokenEntryDef

	tokenEntry.userName = userName
	tokenEntry.createdDate = time.Now()
	_tokenMap[token] = tokenEntry
gklog.LogTrace(fmt.Sprintf("add map entry k: %+v v: %+v",token,tokenEntry))
}

// check if the token / userName is valid
func CheckToken(token string, userName string) bool {
	_tokenMapMutex.Lock()
	defer _tokenMapMutex.Unlock()

	checkTokenExpire()

	var tokenEntry tokenEntryDef
	var ok bool

	tokenEntry, ok = _tokenMap[token]
gklog.LogTrace(fmt.Sprintf("check map entry k: %+v v: %+v",token,ok))
	if ok {
gklog.LogTrace(fmt.Sprintf("check map entry k: %+v v: %+v",token,tokenEntry))
		if tokenEntry.userName == userName {
			return true
		}
	}

	return false
}

// purge any expired tokens
func checkTokenExpire() {
	expireTime := time.Now().Add(time.Duration(-1) * _tokenExpiry)

	for k,v := range _tokenMap {
		if expireTime.After(v.createdDate) {
gklog.LogTrace(fmt.Sprintf("removing map entry (timeout) k: %+v v: %+v",k,v))
			delete(_tokenMap,k)
		}
	}
}

