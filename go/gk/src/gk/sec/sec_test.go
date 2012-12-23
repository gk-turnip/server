/*
	Copyright 2012 1620469 Ontario Limited.

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

package sec

import (
	"testing"
)

func TestSec(t *testing.T) {
	testSalt(t)
	testPasswordHash(t)
}

func testPasswordHash(t *testing.T) {
	var hash []byte
	var password []byte = []byte("one")
	var salt []byte = []byte("two")

	hash = GenPasswordHash(password, salt)

	if string(hash) != "kdjcaicmghakemkiflkcjfhnncafmpedmaieboemekbeieeamfbchebmljkapalcfciljdgeialpiiogkgcabgdlaonbibpdnhdmaelafajimfdgngdagjadhgnbjklc" {
		t.Logf("invalid hash: %s", string(hash))
		t.Fail()
	}
}

func testSalt(t *testing.T) {
	var salt1 []byte
	var salt2 []byte
	var err error

	salt1, err = GenSalt()
	if err != nil {
		t.Logf("error: %v", err)
		t.Fail()
		return
	}
	if len(salt1) != 20 {
		t.Logf("invalid salt: %s", salt1)
		t.Fail()
		return
	}

	salt2, err = GenSalt()
	if err != nil {
		t.Logf("error: %v", err)
		t.Fail()
		return
	}
	if len(salt2) != 20 {
		t.Logf("invalid salt: %s", salt2)
		t.Fail()
		return
	}

	if string(salt1) == string(salt2) {
		t.Logf("salt not different %s %s", salt1, salt2)
		t.Fail()
		return
	}
}
