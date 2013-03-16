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

package message

import (
	"testing"
)

import (
	"gk/gkerr"
)

func TestMessage(t *testing.T) {
	testTrimBetweenMarkers(t)
	testTrimCrLf(t)
	testPopulateFromMessage(t)
	testTrimLeadingSpaces(t)
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

func testTrimLeadingSpaces(t *testing.T) {
	if string(trimLeadingSpaces([]byte("one\ntwo\nthree"))) != "one\ntwo\nthree" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("one\r\ntwo\r\nthree"))) != "one\r\ntwo\r\nthree" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("one\n  two\nthree"))) != "one\n  two\nthree" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("one\ntwo\nthree\n "))) != "one\ntwo\nthree\n" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte(" one\ntwo\nthree\n "))) != "one\ntwo\nthree\n" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("  "))) != "" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte(""))) != "" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte(" one\ntwo   three\nfour\n "))) != "one\ntwo   three\nfour\n" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("one\ntwo\nthree   \n"))) != "one\ntwo\nthree   \n" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
	if string(trimLeadingSpaces([]byte("<g id=\"test\">\r\n </g>\r\n"))) != "<g id=\"test\">\r\n</g>\r\n" {
		t.Logf("trimLeadingSpaces")
		t.Fail()
	}
}

func testPopulateFromMessage(t *testing.T) {
	//	var command string
	//	var jsonData []byte
	//	var data []byte
	var gkErr *gkerr.GkErrDef
	var message []byte
	var messageFromClient *MessageFromClientDef
	var sessionId string = "test"

	message = []byte("com~{\"name\":\"value\"}~data")
	messageFromClient = new(MessageFromClientDef)
	gkErr = messageFromClient.PopulateFromMessage(sessionId, message)
	if gkErr != nil {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if messageFromClient.Command != "com" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.JsonData) != "{\"name\":\"value\"}" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.data) != "data" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	message = []byte("com~{\"name\":\"value\"}~")
	messageFromClient = new(MessageFromClientDef)
	gkErr = messageFromClient.PopulateFromMessage(sessionId, message)
	if gkErr != nil {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if messageFromClient.Command != "com" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.JsonData) != "{\"name\":\"value\"}" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.data) != "" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	message = []byte("commandOnly~~")
	messageFromClient = new(MessageFromClientDef)
	gkErr = messageFromClient.PopulateFromMessage(sessionId, message)
	if gkErr != nil {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if messageFromClient.Command != "commandOnly" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.JsonData) != "" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}
	if string(messageFromClient.data) != "" {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data) + " gkErr: " + gkErr.String())
		t.Fail()
	}

	message = []byte("com~{\"name\":\"value\"}data")
	messageFromClient = new(MessageFromClientDef)
	gkErr = messageFromClient.PopulateFromMessage(sessionId, message)
	if gkErr == nil {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data))
		t.Fail()
	}
	message = []byte("com{\"name\":\"value\"}data")
	messageFromClient = new(MessageFromClientDef)
	gkErr = messageFromClient.PopulateFromMessage(sessionId, message)
	if gkErr == nil {
		t.Logf("PopulateFromMessage message: " + string(message) + " jsonData: " + string(messageFromClient.JsonData) + " data: " + string(messageFromClient.data))
		t.Fail()
	}
}
