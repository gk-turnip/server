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
	"testing"
)

import (
	"gk/gkerr"
)

func TestGameServer(t *testing.T) {
	testTrimBetweenMarkers(t)
	testTrimCrLf(t)
	testGetCommandJsonData(t)
}

func testTrimBetweenMarkers(t *testing.T) {
	if string(trimBetweenMarkers(
		[]byte("one two three four five"), "two", "three")) !=
		"one  four five" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("one two three two four two five two"), "two", "two")) !=
		"one  four " {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(
		trimBetweenMarkers([]byte("<a\nb\nc\nd>"), "<", ">")) != "" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("a<\nb\nc>\nd>"), "<", ">")) !=
		"a\nd>" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("one<two>three<four"), "<", ">")) !=
		"onethree<four" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("one|two|three|four"), "|", "|")) !=
		"onethree|four" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("onexxxtwoxxxthreexxxfour"), "xxx", "xxx")) !=
		"onethreexxxfour" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("onexxxxtwoxxxthreexxxfour"), "xxx", "xxx")) !=
		"onethreexxxfour" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("Xone two threeX"), "X", "X")) != "" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("Xone two threeXX"), "X", "X")) != "X" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("XYZone two threeAB"), "XYZ", "AB")) != "" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("XYone two threeABC"), "XY", "ABC")) != "" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

	if string(trimBetweenMarkers(
		[]byte("zero XY one two three ABC four"), "XY", "ABC")) != "zero  four" {

		t.Logf("trimBetweenMarkers")
		t.Fail()
	}

}

func testTrimCrLf(t *testing.T) {
	if string(trimCrLf([]byte("one two"))) != "one two" {
		t.Logf("trimCrLf")
		t.Fail()
	}

	if string(trimCrLf([]byte("one\ntwo"))) != "onetwo" {
		t.Logf("trimCrLf")
		t.Fail()
	}
	if string(trimCrLf([]byte("one\ntwo\n"))) != "onetwo" {
		t.Logf("trimCrLf")
		t.Fail()
	}
	if string(trimCrLf([]byte("one\r\ntwo\n"))) != "onetwo" {
		t.Logf("trimCrLf")
		t.Fail()
	}
	if string(trimCrLf([]byte("one\r\ntwo\r\n"))) != "onetwo" {
		t.Logf("trimCrLf")
		t.Fail()
	}
}

func testGetCommandJsonData(t *testing.T) {
	var command string
	var jsonData []byte
	var data []byte
	var gkErr *gkerr.GkErrDef
	var message []byte

	message = []byte("com~{\"name\":\"value\"}~data")
	command, jsonData, data, gkErr = getCommandJsonData(message)
	if gkErr != nil {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if command != "com" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(jsonData) != "{\"name\":\"value\"}" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(data) != "data" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	message = []byte("com~{\"name\":\"value\"}~")
	command, jsonData, data, gkErr = getCommandJsonData(message)
	if gkErr != nil {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if command != "com" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(jsonData) != "{\"name\":\"value\"}" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(data) != "" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	message = []byte("commandOnly~~")
	command, jsonData, data, gkErr = getCommandJsonData(message)
	if gkErr != nil {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if command != "commandOnly" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(jsonData) != "" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(data) != "" {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data) + " gkErr: " + gkErr.String())
		t.Fail()
	}

	message = []byte("com~{\"name\":\"value\"}data")
	command, jsonData, data, gkErr = getCommandJsonData(message)
	if gkErr == nil {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data))
		t.Fail()
	}
	message = []byte("com{\"name\":\"value\"}data")
	command, jsonData, data, gkErr = getCommandJsonData(message)
	if gkErr == nil {
		t.Logf("getCommandJsonData message: " + string(message) + " jsonData: " + string(jsonData) + " data: " + string(data))
		t.Fail()
	}

}

