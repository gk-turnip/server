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

package gksvg

import (
	"testing"
)

import (
	"gk/gkerr"
)

func TestSVG(t *testing.T) {

	testSvgLowLevel(t)
	testSvgMidLevel(t)
	testSvgHighLevel(t)
}

func testSvgLowLevel(t *testing.T) {

	var idIn string
	var idOut string
	var gkErr *gkerr.GkErrDef

	idIn = "#id"
	idOut, gkErr = getIdOutOfHref(idIn)
	if gkErr != nil {
		t.Logf("gkErr on getIdOutOfHref in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	if idOut != "id" {
		t.Logf("got wrong in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	idIn = "id"
	idOut, gkErr = getIdOutOfHref(idIn)
	if gkErr == nil {
		t.Logf("missing gkErr on getIdOutOfHref in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	idIn = "fill:url(#radialGradient9814);fill-opacity:1"
	idOut, gkErr = getIdOutOfStyle(idIn)
	if gkErr != nil {
		t.Logf("gkErr on getIdOutOfStyle in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	if idOut != "radialGradient9814" {
		t.Logf("got wrong in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	idIn = "opacity:0.8;fill:url(#radialGradient14051)"
	idOut, gkErr = getIdOutOfStyle(idIn)
	if gkErr != nil {
		t.Logf("gkErr on getIdOutOfStyle in: " + idIn + " out: " + idOut)
		t.Fail()
	}

	if idOut != "radialGradient14051" {
		t.Logf("got wrong in: " + idIn + " out: " + idOut)
		t.Fail()
	}
}

func testSvgMidLevel(t *testing.T) {
	var result string
	var gkErr *gkerr.GkErrDef

	var idMap map[string]string = make(map[string]string)
	var space, name, value string

	idMap["id1"] = "new_id1"
	idMap["id2"] = "new_id2"

	name="href"
	space="xlink"
	value="#id1"
	result, gkErr = substituteOneAttributeId(idMap, space, name, value)
	if gkErr != nil {
		t.Logf("gkErr on substituteOneAttributeId")
		t.Fail()
	}
	if result != "#new_id1" {
		t.Logf("invalid result on substituteOneAttibuteId")
		t.Fail()
	}

	name="style"
	space=""
	value="fill:url(#id2);fill-opacity:1"
	result, gkErr = substituteOneAttributeId(idMap, space, name, value)
	if gkErr != nil {
		t.Logf("gkErr on substituteOneAttributeId " + gkErr.String())
		t.Fail()
	}
	if result != "fill:url(#new_id2);fill-opacity:1" {
		t.Logf("invalid result on substituteOneAttibuteId " + result)
		t.Fail()
	}

	name="style"
	space=""
	value="stroke:none;fill-opacity:1"
	result, gkErr = substituteOneAttributeId(idMap, space, name, value)
	if gkErr != nil {
		t.Logf("gkErr on substituteOneAttributeId " + gkErr.String())
		t.Fail()
	}
	if result != "stroke:none;fill-opacity:1" {
		t.Logf("invalid result on substituteOneAttibuteId " + result)
		t.Fail()
	}
}

const svgInputData1 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!-- created by some editor or other...  -->

<svg
	xmlns:svg="http://www.w3.org/2000/svg"
	xmlns="http://www.w3.org/2000/svg"
   xmlns:xlink="http://www.w3.org/1999/xlink">
<linearGradient id="linearGradient1"/>
<radialGradient id="radialGradient1" xlink:href="#linearGradient1"/>
<path style="fill:url(#radialGradient1);fill-opacity:1"/>
</svg>
`

const svgOutputData1 = `<svg xmlns:svg="http://www.w3.org/2000/svg" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><g><linearGradient id="pre_linearGradient1"></linearGradient><radialGradient id="pre_radialGradient1" xlink:href="#pre_linearGradient1"></radialGradient><path style="fill:url(#pre_radialGradient1);fill-opacity:1"></path></g></svg>`

func testSvgHighLevel(t *testing.T) {
	var gkErr *gkerr.GkErrDef

	var inputData []byte = []byte(svgInputData1)
	var result []byte

	result, gkErr = FixSvgData(inputData, "pre")
	if gkErr != nil {
		t.Logf("FixSvgData failure " + gkErr.String())
		t.Fail()
	}

	if string(result) != svgOutputData1 {
		t.Logf("FixSvgData did not match in: " + svgInputData1 + "\n out: " + string(result) + "\n exp: " + svgOutputData1 + "\n")
		t.Fail()
	}
}

