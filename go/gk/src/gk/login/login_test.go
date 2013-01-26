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
	"testing"
)

func TestLogin(t *testing.T) {
	testValidUsers(t)
	testValidPasswords(t)
	testValidEmail(t)
}

func testValidUsers(t *testing.T) {
	var userName string

	userName = ""
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = "ab"
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = " one"
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = "one "
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = "one  two"
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = "one\ttwo"
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}
	userName = "123456789012345678901"
	if isNewUserNameValid(userName) {
		t.Logf("invalid user: %s", userName)
		t.Fail()
	}

	userName = "one"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
	userName = "one two"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
	userName = "One two"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
	userName = "One ._-+two"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
	userName = "1234567890"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
	userName = "12345678901234567890"
	if !isNewUserNameValid(userName) {
		t.Logf("valid user: %s", userName)
		t.Fail()
	}
}

// must have at least one upper or lower case letter
// must have at least one digit, or one special character
func testValidPasswords(t *testing.T) {
	var password string

	password = ""
	if isPasswordValid(password) {
		t.Logf("invalid password: %s", password)
		t.Fail()
	}
	password = "a1.%e"
	if isPasswordValid(password) {
		t.Logf("invalid password: %s", password)
		t.Fail()
	}
	password = "a1.%ef\tg"
	if isPasswordValid(password) {
		t.Logf("invalid password: %s", password)
		t.Fail()
	}
	password = "123456!@#"
	if isPasswordValid(password) {
		t.Logf("invalid password: %s", password)
		t.Fail()
	}
	password = "onetwo"
	if isPasswordValid(password) {
		t.Logf("invalid password: %s", password)
		t.Fail()
	}

	password = "a.Cdef"
	if !isPasswordValid(password) {
		t.Logf("valid password: %s", password)
		t.Fail()
	}
	password = "ab2@ef"
	if !isPasswordValid(password) {
		t.Logf("valid password: %s", password)
		t.Fail()
	}
	password = "abcD2f"
	if !isPasswordValid(password) {
		t.Logf("valid password: %s", password)
		t.Fail()
	}
	password = "a cD2f"
	if !isPasswordValid(password) {
		t.Logf("valid password: %s", password)
		t.Fail()
	}

}

func testValidEmail(t *testing.T) {
	var email string

	email = ""
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one.two.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one.two@"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one.two@x"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one.@two.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = ".one@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one..two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one@two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one\"two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one<two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one()<>[]:,\\two@three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one@two_three.com"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one@two.com-"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}
	email = "one@two"
	if isEmailValid(email) {
		t.Logf("invalid email: %s", email)
		t.Fail()
	}

	email = "one@two.edu"
	if !isEmailValid(email) {
		t.Logf("valid email: %s", email)
		t.Fail()
	}
	email = "one.two@three.gov"
	if !isEmailValid(email) {
		t.Logf("valid email: %s", email)
		t.Fail()
	}
//	email = "one@[IPv6:0001:0002:0003:0004:0005:0006:0007:0008]"
//	if !isEmailValid(email) {
//		t.Logf("valid email: %s", email)
//		t.Fail()
//	}
//	email = "\"one two\"@three.com"
//	if !isEmailValid(email) {
//		t.Logf("valid email: %s", email)
//		t.Fail()
//	}
//	email = "\"\"@three.com"
//	if !isEmailValid(email) {
//		t.Logf("valid email: %s", email)
//		t.Fail()
//	}
//	email = "\"()<>[]:,\\@\"@one.net"
//	if !isEmailValid(email) {
//		t.Logf("valid email: %s", email)
//		t.Fail()
//	}
	email = "!#$%&'*+-/=?^_`{}~@one.com"
	if !isEmailValid(email) {
		t.Logf("valid email: %s", email)
		t.Fail()
	}
	email = "test@one-two.com"
	if !isEmailValid(email) {
		t.Logf("valid email: %s", email)
		t.Fail()
	}
}

