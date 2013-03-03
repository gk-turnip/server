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

package gkrand

import (
	"testing"
)

func TestRand(t *testing.T) {

	var randContext *GkRandContextDef

	randContext = NewGkRandContext()

	var s1 string
	var s2 string
	s1 = randContext.GetRandomString(10)
	if len(s1) != 10 {
		t.Logf("invalid return len: %s", s1)
		t.Fail()
	}
	s2 = randContext.GetRandomString(10)
	if len(s2) != 10 {
		t.Logf("invalid return len: %s", s2)
		t.Fail()
	}

	if s1 == s2 {
		t.Logf("invalid return %s %s", s1, s2)
		t.Fail()
	}
}

