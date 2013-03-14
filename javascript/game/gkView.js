
// this is the view into the pod handler

var gkViewContext = new gkViewContextDef();

function gkViewContextDef() {
	this.viewOffsetIsoXYZ = new GkIsoXYZDef(-20, 20, 0);
	this.scale = 1.0;
	this.fernWidth = 6;
	this.fernHeight = 6;
	this.svgWidth = 600;
	this.svgHeight = 300;
	this.viewMap = new Object();
	this.marginX = 5;
	this.marginY = 10;
	this.scrollEdgeX = 100;
	this.scrollEdgeY = 100;
}

function gkViewObjectMapEntryDef(terrainMapMapEntry, terrainSvgMapEntry, g, inUse) {
	this.terrainMapMapEntry = terrainMapMapEntry
	this.terrainSvgMapEntry = terrainSvgMapEntry
	this.g = g
	this.inUse = inUse;
}

// mark every object as not in use
// traverse every viwable coordiate, fern units at a time
// display the entry and mark it as used
// and then remove any viewMap objects that were not used
// since they should be objects that have scrolled off screen
function gkViewRender() {

	console.log("gkViewRender");

	gkViewContext.fernWidth = Math.round((gkViewContext.svgWidth / 100) / gkViewContext.scale);
	gkViewContext.fernHeight = Math.round((gkViewContext.svgHeight / 50) / gkViewContext.scale);
	
	gkViewSetObjectMapNotUsed();

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

	var start_x, start_y;

	start_x = gkViewContext.viewOffsetIsoXYZ.x;
	start_y = gkViewContext.viewOffsetIsoXYZ.y;

	for (i = 0;i < gkViewContext.fernHeight; i++) {
		var x, y;
		x = start_x;
		y = start_y;
		gkViewRenderSingleFern(x,y);
		for (j = 0;j < gkViewContext.fernWidth; j++) {
			x += 10;
			gkViewRenderSingleFern(x,y);
			y -= 10;
			gkViewRenderSingleFern(x,y);
		}

		start_x += 10;
		start_y += 10;
	}

	// remove objects still marked as un used
	gkViewRemoveNotUsed();
}

// called for each position in the viewable area
// keeps track of which viewMap entries are in use
function gkViewRenderSingleFern(rawX, rawY) {

	x = Math.round(rawX / 10) * 10;
	y = Math.round(rawY / 10) * 10;
	var mapKey = gkTerrainGetMapKey(x, y);

	var viewObjectMapEntry = gkViewContext.viewMap[mapKey]

	if (viewObjectMapEntry == null) {
		var terrainMapMapEntry = gkTerrainGetMapMapEntry(x, y);
		// possibly undefined because the render loop went outside the terrain map
		if (terrainMapMapEntry == undefined) {
//			console.error("terrainMapMapEntry undefined for x,y: " + x + "," + y);
		} else {
			var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(x, y);

			if (terrainSvgMapEntry == undefined) {
//				console.error("terranSvgMapEntry undefined for x,y: " + x + "," + y);
			} else {
				gkViewAddObjectMapEntry(mapKey, terrainMapMapEntry, terrainSvgMapEntry)
				viewObjectMapEntry = gkViewContext.viewMap[mapKey]
			}
		}
	} else {
		viewObjectMapEntry.inUse = true;
	}
}

function gkViewAddObjectMapEntry(mapKey, terrainMapMapEntry, terrainSvgMapEntry) {
	g = gkIsoCreateSvgObject(terrainSvgMapEntry.svgSegment);
	var viewObjectMapEntry = new gkViewObjectMapEntryDef(terrainMapMapEntry, terrainSvgMapEntry, g, true);
	gkViewContext.viewMap[mapKey] = viewObjectMapEntry;
	var layer = document.getElementById(terrainSvgMapEntry.layer);
	var x = viewObjectMapEntry.terrainMapMapEntry.x
	var y = viewObjectMapEntry.terrainMapMapEntry.y
	var z = viewObjectMapEntry.terrainMapMapEntry.z
	var isoXYZ = new GkIsoXYZDef(x,y,z);
	gkIsoSetSvgObjectPositionWithOffset(g, isoXYZ, gkTerrainContext.terrainDiamondOffsetX, gkTerrainContext.terrainDiamondOffsetY);
	layer.appendChild(g);
}

function gkViewRemoveNotUsed() {
	for (var mapKey in gkViewContext.viewMap) {
		if (!gkViewContext.viewMap[mapKey].inUse) {
			// delete from gkField
			var g = gkViewContext.viewMap[mapKey].g
			g.parentNode.removeChild(g);
			// delete from viewMap
			delete gkViewContext.viewMap[mapKey];
		}
	}
}

function gkViewSetObjectMapNotUsed() {
	for (var mapKey in gkViewContext.viewMap) {
		var objectMapEntry = gkViewContext.viewMap[mapKey]
		objectMapEntry.inUse = false;
	}
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

