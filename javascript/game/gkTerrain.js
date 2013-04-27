
// this is the terrain handler

var gkTerrainContext = new gkTerrainContextDef();

function gkTerrainContextDef() {
	// map holding the position of each terrain diamond
	this.terrainMapMap = new Object();

	// map holding the svg data for each terrain diamond
	this.terrainSvgMap = new Object();

	this.terrainDiamondOffsetX = 50;
	this.terrainDiamondOffsetY = 0;
}

function gkTerrainMapMapEntryDef(x, y, z, terrainName) {
	this.x = x
	this.y = y
	this.z = z
	this.terrainName = terrainName
}

function gkTerrainSvgMapEntryDef(terrainName, originX, originY, layer) {
	this.terrainName = terrainName;
	this.originX = originX;
	this.originY = originY;
	this.layer = layer;
}

function gkTerrainInit() {
}

// called as a request from the server
// to set the entire pod terran map
// the terrain svg must be done before the terrain map
function gkSetTerrainMap(jsonData) {
//console.log("gkSetTerrainMap");
	var i;

console.log("tileList.length: " + jsonData.tileList.length);
	for (i = 0;i < jsonData.tileList.length; i++) {
		var x, y, z, terrainName;
		x = jsonData.tileList[i].x;
		y = jsonData.tileList[i].y;
		z = jsonData.tileList[i].z;
		terrainName = jsonData.tileList[i].t;

		var terrainMapMapEntry = new gkTerrainMapMapEntryDef(x, y, z, terrainName)

		var mapKey = gkTerrainGetMapKey(x, y);
		gkTerrainContext.terrainMapMap[mapKey] = terrainMapMapEntry;

		var baseLayer = document.getElementById("gkTerrainBaseLayer");

		var isoXYZ = new GkIsoXYZDef(x,y,z);
		var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(x,y);

		if (terrainSvgMapEntry != undefined) {
			var ref = document.createElementNS(gkIsoContext.svgNameSpace,"use");
			ref.setAttributeNS(gkIsoContext.xlinkNameSpace,"href","#t_" + terrainName);
			gkIsoSetSvgUsePositionWithOffset(ref, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY);
			baseLayer.appendChild(ref);
		}
	}

console.log("objectList.length: " + jsonData.objectList.length);
	for (i = 0;i < jsonData.objectList.length; i++) {
		var x, y, z, objectName;
		x = jsonData.objectList[i].x;
		y = jsonData.objectList[i].y;
		z = jsonData.objectList[i].z;
		objectName = jsonData.objectList[i].o;

		var objectMapMapEntry = new gkTerrainMapMapEntryDef(x, y, z, objectName)

		var mapKey = gkTerrainGetMapKey(x, y);
		gkTerrainContext.terrainMapMap[mapKey] = objectMapMapEntry;

		var objectLayer = document.getElementById("gkTerrainObjectLayer");

		var isoXYZ = new GkIsoXYZDef(x,y,z);
		var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(x,y);

		if (terrainSvgMapEntry != undefined) {
			var ref = document.createElementNS(gkIsoContext.svgNameSpace,"use");
			ref.setAttributeNS(gkIsoContext.xlinkNameSpace,"href","#t_" + objectName);
			gkIsoSetSvgUsePositionWithOffset(ref, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY);
			objectLayer.appendChild(ref);
		}
	}


	gkViewRender();
}

// called as a request from the server
// to set all the svg files for the require terrain
function gkSetTerrainSvg(jsonData, rawSvgData) {
//console.log("gkSetTerrainSvg");
	if (jsonData.terrain != undefined) {
		var terrainName, originX, originY, layer;

console.log("gkSetTerrainSvg name: " + jsonData.terrain);
		terrainName = jsonData.terrain;
		originX = jsonData.originX;
		originY = jsonData.originY;
		layer = jsonData.layer;
		var terrainSvgMapEntry = new gkTerrainSvgMapEntryDef(terrainName, originX, originY, layer);
		gkTerrainContext.terrainSvgMap[jsonData.terrain] = terrainSvgMapEntry;

		var gkDefs = document.getElementById("gkDefs");

		g = gkIsoCreateSvgObject(rawSvgData);
		g.setAttribute("id","t_" + jsonData.terrain);

		gkDefs.appendChild(g);
	}
}

// create the "key" to the map from the x, y
function gkTerrainGetMapKey(x, y) {
	return "k" + x + "," + y
}

function gkTerrainGetMapMapEntry(x, y) {
	var mapKey = gkTerrainGetMapKey(x, y);

	terrainMapMapEntry = gkTerrainContext.terrainMapMap[mapKey];

	return terrainMapMapEntry
}

function gkTerrainGetSvgMapEntry(x, y) {
	var mapKey = gkTerrainGetMapKey(x, y);

	var terrainMapMapEntry = gkTerrainContext.terrainMapMap[mapKey];
	var svgMapEntry

	if (terrainMapMapEntry == undefined) {
		console.error("missing terrainMapMapEntry mapKey: " + mapKey);
	} else {
		var terrainName = terrainMapMapEntry.terrainName

		svgMapEntry = gkTerrainContext.terrainSvgMap[terrainName]
	}

	return svgMapEntry
}

