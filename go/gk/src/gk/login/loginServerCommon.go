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

import "strings"

const MIN_USER_NAME_LENGTH = 3
const MAX_USER_NAME_LENGTH = 20
const MIN_PASSWORD_LENGTH = 6
const MAX_PASSWORD_LENGTH = 60

func isNewUserNameValid(userName string) bool {
	if len(userName) < MIN_USER_NAME_LENGTH {
		return false
	}
	if len(userName) > MAX_USER_NAME_LENGTH {
		return false
	}
	if userName[0] == ' ' {
		return false
	}
	if userName[len(userName) - 1] == ' ' {
		return false
	}
	if strings.Index(userName,"  ") != -1 {
		return false
	}

	for i := 0;i < len(userName); i++ {
		if (userName[i] < '0' || userName[i] > '9') &&
			(userName[i] < 'a' || userName[i] > 'z') &&
			(userName[i] < 'A' || userName[i] > 'Z') &&
			(userName[i] != ' ') &&
			(userName[i] != '-') &&
			(userName[i] != '_') &&
			(userName[i] != '.') &&
			(userName[i] != '+') {
			return false
		}
	}

	return true;
}

func isPasswordValid(password string) bool {
	if len(password) < MIN_PASSWORD_LENGTH {
		return false
	}
	if len(password) > MAX_PASSWORD_LENGTH {
		return false
	}
	if password[0] == ' ' {
		return false
	}
	if password[len(password) - 1] == ' ' {
		return false
	}
	for i := 0;i < len(password); i++ {
		if (password[i] < '0' || password[i] > '9') &&
			(password[i] < 'a' || password[i] > 'z') &&
			(password[i] < 'A' || password[i] > 'Z') {
			if strings.IndexRune("`~!@ #$%^&*()-=_+[]{};':\",./<>?\\|",rune(password[i])) == -1 {
				return false
			}
		}
	}

	var digitsCount int = 0
	var lowerCount int = 0
	var upperCount int = 0
	var specialCount int = 0
	
	for i := 0;i < len(password); i++ {
		if password[i] >= '0' && password[i] <= '0' {
			digitsCount += 1
		} else {
			if password[i] >= 'a' && password[i] <= 'z' {
				lowerCount += 1
			} else {
				if password[i] >= 'A' && password[i] <= 'Z' {
					upperCount += 1
				} else {
					specialCount += 1
				}
			}
		}
	}

	if lowerCount == 0 && upperCount == 0 {
		return false
	}

	if digitsCount == 0 && specialCount == 0 {
		return false
	}

	return true
}

func isEmailValid(email string) bool {
	if len(email) < 3 {
		return false
	}

	if email[0] == ' ' {
		return false
	}
	if email[0] == '.' {
		return false
	}
	if email[len(email) - 1] == ' ' {
		return false
	}
	if email[len(email) - 1] == '-' {
		return false
	}

	var atIndex int
	atIndex = strings.Index(email,"@")

	if atIndex == -1 {
		return false
	}

	if (atIndex + 2) >= len(email) {
		return false
	}

	if strings.Index(email[atIndex + 1:],"@") != -1 {
		return false
	}

	if strings.Index(email[atIndex + 1:],".") == -1 {
		return false
	}

	if strings.Index(email,"..") != -1 {
		return false
	}

	if strings.Index(email,".@") != -1 {
		return false
	}

	for i := 0;i < len(email); i++ {
		if i > atIndex {
			if (email[i] < '0' || email[i] > '9') &&
				(email[i] < 'a' || email[i] > 'z') &&
				(email[i] < 'A' || email[i] > 'Z') &&
				(email[i] != '.') &&
				(email[i] != '-') {
				return false
			}
		} else {
			if (email[i] < '0' || email[i] > '9') &&
				(email[i] < 'a' || email[i] > 'z') &&
				(email[i] < 'A' || email[i] > 'Z') &&
				(email[i] != '@') &&
				(email[i] != '.') &&
				(email[i] != '!') &&
				(email[i] != '#') &&
				(email[i] != '$') &&
				(email[i] != '%') &&
				(email[i] != '&') &&
				(email[i] != '\'') &&
				(email[i] != '*') &&
				(email[i] != '+') &&
				(email[i] != '-') &&
				(email[i] != '/') &&
				(email[i] != '=') &&
				(email[i] != '?') &&
				(email[i] != '^') &&
				(email[i] != '_') &&
				(email[i] != '`') &&
				(email[i] != '{') &&
				(email[i] != '}') &&
				(email[i] != '~') {
				return false
			}
		}
	}

	return true;
}

