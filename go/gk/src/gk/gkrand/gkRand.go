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

// package to handle the context of random numbers

package gkrand

import (
	c_rand "crypto/rand"
	m_rand "math/rand"
	"sync"
	"time"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

const validDataSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type GkRandContextDef struct {
	mRandContext *m_rand.Rand
	mutex        sync.Mutex
}

func NewGkRandContext() *GkRandContextDef {
	var gkRandContext *GkRandContextDef = new(GkRandContextDef)
	var seed int64
	var err error
	var gkErr *gkerr.GkErrDef

	seed = time.Now().UnixNano()
	buf := make([]byte, 6, 6)
	_, err = c_rand.Read(buf)
	if err != nil {
		// log the error
		// but this error is not fatal
		gkErr = gkerr.GenGkErr("c_rand.Read", err, ERROR_ID_RAND_READ)
		gklog.LogGkErr("", gkErr)
	}

	seed ^= int64(buf[0]) << 16
	seed ^= int64(buf[1]) << 24
	seed ^= int64(buf[2]) << 32
	seed ^= int64(buf[3]) << 40
	seed ^= int64(buf[4]) << 48
	seed ^= int64(buf[5]) << 56

	gkRandContext.mRandContext = m_rand.New(m_rand.NewSource(seed))

	return gkRandContext
}

func (gkRandContext *GkRandContextDef) GetRandomString(length int) string {
	result := make([]byte, 0, length)

	gkRandContext.mutex.Lock()
	defer gkRandContext.mutex.Unlock()

	for len(result) < length {
		var r int64
		var err error
		var gkErr *gkerr.GkErrDef

		buf := make([]byte, 6, 6)
		_, err = c_rand.Read(buf)
		if err != nil {
			// log the error
			// but this error is not fatal
			gkErr = gkerr.GenGkErr("c_rand.Read", err, ERROR_ID_RAND_READ)
			gklog.LogGkErr("", gkErr)
		}
		r = gkRandContext.mRandContext.Int63()
		r ^= int64(buf[0]) << 16
		r ^= int64(buf[1]) << 24
		r ^= int64(buf[2]) << 32
		r ^= int64(buf[3]) << 40

		for i := 0; i < 10; i++ {
			var c int
			c = int(r & 0x3f)
			if c < len(validDataSet) {
				result = append(result, validDataSet[c])
				if len(result) >= length {
					break
				}
			}
			r = r >> 6
		}
	}

	return string(result)
}

func (gkRandContext *GkRandContextDef) GetRandomByte() byte {
	gkRandContext.mutex.Lock()
	defer gkRandContext.mutex.Unlock()

	var err error
	var gkErr *gkerr.GkErrDef

	buf := make([]byte, 1, 1)
	_, err = c_rand.Read(buf)
	if err != nil {
		// log the error
		// but this error is not fatal
		gkErr = gkerr.GenGkErr("c_rand.Read", err, ERROR_ID_RAND_READ)
		gklog.LogGkErr("", gkErr)
	}
	var r int64
	r = gkRandContext.mRandContext.Int63()
	r ^= int64(buf[0])

	return byte(r)
}
