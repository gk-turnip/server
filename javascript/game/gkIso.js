
var GK_SVG_NAMESPACE = "http://www.w3.org/2000/svg";
var GK_SVG_MARGIN_X = 5;
var GK_SVG_MARGIN_Y = 5;
var GK_SVG_WIDTH = 600;
var GK_SVG_HEIGHT = 300;

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
	diamond.setAttribute("opacity", opacity);

	return diamond;
}

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

