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

// security package
package sec

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"hash"
	"io"
)

var passwordHashConstant = []byte("jvk56j3Bu") // this value must not change
const _hashLoopCount = 5000                    // this value must not change
const _saltLength = 10
const _forgotPasswordTokenLength = 12

var saltValues = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generate a hash for the given password and salt
// intentionally slow for better security
func GenPasswordHashSlow(password []byte, salt []byte) []byte {
	var h hash.Hash
	var lastHash []byte

	h = sha512.New()
	h.Write(passwordHashConstant)
	h.Write(password)
	h.Write(salt)
	lastHash = h.Sum(nil)

	for i := 0; i < _hashLoopCount; i++ {
		h.Reset()
		h.Write(salt)
		h.Write(lastHash)
		h.Write(password)
		h.Write(passwordHashConstant)
		lastHash = h.Sum(nil)
	}

	r := make([]byte, 0, 2*len(lastHash))
	for i := 0; i < len(lastHash); i++ {
		r = append(r, (lastHash[i]&0x0f)+'a')
		r = append(r, ((lastHash[i]>>4)&0x0f)+'a')
	}

	return r
}

// generate a hash for the given password and salt
// not as secure as GenPasswordHashSlow
func GenPasswordHashFast(password []byte, salt []byte) []byte {
	var h hash.Hash
	var lastHash []byte

	h = sha512.New()
	h.Write(passwordHashConstant)
	h.Write(password)
	h.Write(salt)
	lastHash = h.Sum(nil)

	r := make([]byte, 0, 2*len(lastHash))
	for i := 0; i < len(lastHash); i++ {
		r = append(r, (lastHash[i]&0x0f)+'a')
		r = append(r, ((lastHash[i]>>4)&0x0f)+'a')
	}

	return r
}

func GenSalt() ([]byte, error) {
	return genToken(_saltLength)
}

func GenForgotPasswordToken() ([]byte, error) {
	return genToken(_forgotPasswordTokenLength)
}

// geneate a new random salt
func genToken(tokenLen int) ([]byte, error) {
	salt := make([]byte, tokenLen, tokenLen)

	readCount, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	if readCount != len(salt) {
		err = errors.New("GenSalt: could not get random salt")
		return nil, err
	}

	for i := 0; i < len(salt); i++ {
		salt[i] = saltValues[salt[i]%byte(len(saltValues))]
	}

	return salt, nil
}
