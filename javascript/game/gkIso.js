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

// convert between windows x,y coordinates and iso x,y,z coordinates

// higher iso x coordinates go down and to the right
// higher iso y coordinates go down and to the left
// higher iso z coordinates go up

// note: all ISO coordinates are in deciferns or 1/10 of a fern 
// which is 5 wide and 2.5 high in svg units

var gkIsoContext = new gkIsoContextDef();

function gkIsoContextDef() {
	this.svgNameSpace = "http://www.w3.org/2000/svg";
	this.xlinkNameSpace = "http://www.w3.org/1999/xlink";
	this.zFactor = 5;
	this.gridListIndexNamePrefix = "gl_"
}

// create and return a single small diamond (1/10 fern sized)
function gkIsoCreateSingleDiamond(isoXYZ, colour, opacity) {
	winXY = isoXYZ.convertToWin();

	diamond = document.createElementNS(gkIsoContext.svgNameSpace,"polygon");

	x1 = winXY.x;
	y1 = winXY.y;
	x2 = winXY.x - 5;
	y2 = winXY.y + 2.5;
	x3 = winXY.x;
	y3 = winXY.y + 5;
	x4 = winXY.x + 5;
	y4 = winXY.y + 2.5;

	zOffset = isoXYZ.z;
	y1 -= zOffset;
	y2 -= zOffset;
	y3 -= zOffset;
	y4 -= zOffset;

	diamond.setAttribute("points", x1 + "," + y1 + " " + x2 + "," + y2 + " " + x3 + "," + y3 + " " + x4 + "," + y4);

	diamond.setAttribute("fill", colour);
	diamond.setAttribute("stroke", colour);
	diamond.setAttribute("fill-opacity", opacity);
	diamond.setAttribute("stroke-width", "0");

	return diamond;
}

// create an object from raw svg data
// it uses "firstChild" because the uploaded svg file has
// the <svg> tag, but the svg image in the browser also already has
// the <svg> tag, so I take the first child.
// note that sometimes the svg files must be edited to add an extra
// <g> tag around the entire contents after the <svt>
// so that the first child is the entire image, minus the <svg>
function gkIsoCreateSvgObject(rawSvgData) {
	var g
	g = document.createElementNS(gkIsoContext.svgNameSpace,"g");
	var r1 = new DOMParser().parseFromString(rawSvgData, "text/xml");
	g.appendChild(document.importNode(r1.documentElement.firstChild,true))

//	svgDiamond = document.importNode(r1.documentElement.firstChild,true)

	return g
}

// set the position of the object
function gkIsoSetSvgObjectPosition(svgDiamond, isoXYZ) {
	var winXY;
	winXY = isoXYZ.convertToWin();
console.log("win x,y: " + winXY.x + "," + winXY.y);
	svgDiamond.setAttribute("transform","translate(" + winXY.x + "," + winXY.y + ")");
}

// set the position of the object, with originX and originY offsets
function gkIsoSetSvgObjectPositionWithOffset(svgDiamond, isoXYZ, originX, originY, originZ) {
	var winXY;
	winXY = isoXYZ.convertToWin();
	winXY.x -= originX
	winXY.y -= originY
	winXY.y -= originZ;
	svgDiamond.setAttribute("transform","translate(" + winXY.x + "," + winXY.y + ")");
}

// a windows x,y object
function GkWinXYDef(x, y) {
	this.x = x;
	this.y = y;

	GkWinXYDef.prototype.convertToIso = function(z) {
		isoX = Math.floor(((this.y * 2) + this.x ) / 10);
		isoY = Math.floor((this.y - (this.x) / 2) / 5);
		isoZ = z;

		return new GkIsoXYZDef(isoX, isoY, isoZ);
	}
}

// an iso x,y,z object
function GkIsoXYZDef(x, y, z) {
	this.x = x;
	this.y = y;
	this.z = z;

	GkIsoXYZDef.prototype.convertToWin = function() {
		var winX;
		var winY;
		winX = this.x - this.y;
		winY = this.x + this.y;

		winX *= 5;
		winY *= 2.5;
		winY -= this.z * gkIsoContext.zFactor;

		return new GkWinXYDef(winX, winY);
	}
}

// this must return a string that can be orderd
// for grid list
// x and y range should be -32768 to 32767
// so a constant is added to produce a string
// that will always have a 6 digit suffix
function gkIsoGetGridListIndexName(x, y, z) {
	var localX = gkIsoGetFernFromDecifern(x);
	var localY = gkIsoGetFernFromDecifern(y);

	return gkIsoContext.gridListIndexNamePrefix + (165536 + localX + localY);
}

// get the fern coordinate from a decifern coordinate
function gkIsoGetFernFromDecifern(x) {
	localX = x / 10;
	localX = Math.floor(localX);
	localX = localX * 10;
	return localX;
}

