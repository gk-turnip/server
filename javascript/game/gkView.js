
// this is the view into the pod handler

var gkViewContext = new gkViewContextDef();

function gkViewContextDef() {
	this.viewOffsetIsoXYZ = new GkIsoXYZDef(-20, 20, 0);
	this.scale = 1.0;
	this.fernWidth = 6;
	this.fernHeight = 6;
	this.svgWidth = 800;
	this.svgHeight = 400;
//	this.viewMap = new Object();
	this.marginX = 5;
	this.marginY = 10;
	this.scrollEdgeX = 100;
	this.scrollEdgeY = 100;
	this.lowScale = 0.2;
	this.highScale = 3.0;
	this.lowPanX = -3000;
	this.highPanX = 3000;
	this.lowPanY = -2000;
	this.highPanY = 2000;
}

// change the offset and scale of the view
function gkViewRender() {

//	console.log("gkViewRender");

	gkViewContext.fernWidth = Math.round((gkViewContext.svgWidth / 100) / gkViewContext.scale);
	gkViewContext.fernHeight = Math.round((gkViewContext.svgHeight / 50) / gkViewContext.scale);
	
	var gkField = document.getElementById("gkField");

	gkField.setAttribute("width",gkViewContext.svgWidth);
	gkField.setAttribute("height",gkViewContext.svgHeight);
	
	var scaledOffset = new GkIsoXYZDef(
		gkViewContext.viewOffsetIsoXYZ.x,
		gkViewContext.viewOffsetIsoXYZ.y,
		gkViewContext.viewOffsetIsoXYZ.z)

	scaledOffset.x *= gkViewContext.scale;
	scaledOffset.y *= gkViewContext.scale;
	scaledOffset.z *= gkViewContext.scale;

	var winXY = scaledOffset.convertToWin();

	var gkView = document.getElementById("gkView");
	gkView.setAttribute("transform","translate(" + (-winXY.x) + "," + (-winXY.y) + ") scale(" + gkViewContext.scale + ")");

}

function gkViewConvertWinToIso(x, marginX, y, marginY, z) {
	var winXY = new GkWinXYDef(x - marginX, y - marginY);
	var isoXYZ = winXY.convertToIso(z);

	isoXYZ.x /= gkViewContext.scale;
	isoXYZ.y /= gkViewContext.scale;
	isoXYZ.z /= gkViewContext.scale;

	isoXYZ.x += gkViewContext.viewOffsetIsoXYZ.x;
	isoXYZ.y += gkViewContext.viewOffsetIsoXYZ.y;
	isoXYZ.z += gkViewContext.viewOffsetIsoXYZ.z;

	isoXYZ.x = Math.round(isoXYZ.x);
	isoXYZ.y = Math.round(isoXYZ.y);
	isoXYZ.z = Math.round(isoXYZ.z);

	return isoXYZ;
}

