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

var GK_SVG_NAMESPACE = "http://www.w3.org/2000/svg";
var GK_SVG_MARGIN_X = 5;
var GK_SVG_MARGIN_Y = 5;
var GK_SVG_WIDTH = 600;
var GK_SVG_HEIGHT = 300;

// draw a diamond grid
// not called anymore
function gkIsoDrawGridDiamond() {
	var line;
	var gridColour = "#afafaf";

	field = document.getElementById("gkField");


	for (i = 0;i < ((GK_SVG_WIDTH / 10) + 1);i += 10) {
		line = document.createElementNS(GK_SVG_NAMESPACE,"line");
		line.setAttribute("x1",(GK_SVG_WIDTH / 2) + (i * 5));
		line.setAttribute("y1",0 + (i * 2.5));
		line.setAttribute("x2",0 + (i * 5));
		line.setAttribute("y2",(GK_SVG_WIDTH / 4) + (i * 2.5));
		line.setAttribute("stroke", gridColour);
		line.setAttribute('stroke-width', 1);
		field.appendChild(line)
	}

	for (i = 0;i < ((GK_SVG_WIDTH / 10) + 1);i += 10) {
		line = document.createElementNS(GK_SVG_NAMESPACE,"line");
		line.setAttribute("x1",0 + (i * 5));
		line.setAttribute("y1",(GK_SVG_WIDTH / 4) - (i * 2.5));
		line.setAttribute("x2",(GK_SVG_WIDTH / 2) + (i * 5));
		line.setAttribute("y2",(GK_SVG_WIDTH / 2) - (i * 2.5));
		line.setAttribute('stroke', gridColour);
		line.setAttribute('stroke-width', 1);
		field.appendChild(line)
	}
}

// draw a full grid
// not called anymore
function gkIsoDrawGridFull() {
	var line;
	var gridColour = "#afafaf";

	field = document.getElementById("gkField");

	var width;
	if (GK_SVG_WIDTH > GK_SVG_HEIGHT) {
		width = GK_SVG_WIDTH;
	} else {
		width = GK_SVG_HEIGHT;
	}

	for (i = 0;i < ((width / 2.5) + 1);i += 20) {
		line = document.createElementNS(GK_SVG_NAMESPACE,"line");
		line.setAttribute("x1",i * 5);
		line.setAttribute("y1",0);
		line.setAttribute("x2",(width / -0.5) + (i * 5));
		line.setAttribute("y2",width / 1);
		line.setAttribute("stroke", gridColour);
		line.setAttribute('stroke-width', 1);
		field.appendChild(line)
	}

	for (i = 0;i < ((width / 2.5) + 1);i += 20) {
		line = document.createElementNS(GK_SVG_NAMESPACE,"line");
		line.setAttribute("x1",(width / -1) + (i * 5));
		line.setAttribute("y1",0);
		line.setAttribute("x2",(width / 1) + (i * 5));
		line.setAttribute("y2",width / 1);
		line.setAttribute('stroke', gridColour);
		line.setAttribute('stroke-width', 1);
		field.appendChild(line)
	}
}

// create and return a single small diamond (1/10 fern sized)
function gkIsoCreateSingleDiamond(isoXYZ, colour, opacity) {
	winXY = isoXYZ.convertToWin();

	diamond = document.createElementNS(GK_SVG_NAMESPACE,"polygon");

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
function gkIsoCreateSvgDiamond(rawSvgData) {
	var g
	g = document.createElementNS(GK_SVG_NAMESPACE,"g");
	var r1 = new DOMParser().parseFromString(rawSvgData, "text/xml");
	g.appendChild(document.importNode(r1.documentElement.firstChild,true))

//	svgDiamond = document.importNode(r1.documentElement.firstChild,true)

	return g
}

// set the position of the object
function gkIsoSetSvgDiamondPosition(svgDiamond, isoXYZ) {
	var winXY;
	winXY = isoXYZ.convertToWin();
	svgDiamond.setAttribute("transform","translate(" + winXY.x + "," + winXY.y + ")");
}

// set the position of the object, with originX and originY offsets
function gkIsoSetSvgPositionWithOffset(svgDiamond, isoXYZ, originX, originY) {
	var winXY;
	winXY = isoXYZ.convertToWin();
	winXY.x -= originX
	winXY.y -= originY
	svgDiamond.setAttribute("transform","translate(" + winXY.x + "," + winXY.y + ")");
}

// a windows x,y object
function GkWinXYDef(x, y) {
	this.x = x;
	this.y = y;

	GkWinXYDef.prototype.convertToIso = function(z) {
		isoX = Math.floor(((this.y * 2) + this.x - (GK_SVG_WIDTH / 2)) / 10);
		isoY = Math.floor((this.y - (this.x - (GK_SVG_WIDTH / 2)) / 2) / 5);
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
		winX += GK_SVG_WIDTH / 2;
		winY -= this.z * 5;

		return new GkWinXYDef(winX, winY);
	}
}

