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

import (
	"testing"
)

func TestGkNet(t *testing.T) {
	testRequestMatches(t)
}

func testRequestMatches(t *testing.T) {
	if !RequestMatches("a", "a") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if RequestMatches("a", "b") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if !RequestMatches("/one", "/one") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if !RequestMatches("/one", "/one/") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if !RequestMatches("/one/", "/one") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if RequestMatches("/one/", "/two/") {
		t.Logf("RequestMatches")
		t.Fail()
	}
	if RequestMatches("/one", "/two/") {
		t.Logf("RequestMatches")
		t.Fail()
	}
}
