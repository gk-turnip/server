
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

function gkTerrainSvgMapEntryDef(terrainName, layer, svgSegment) {
	this.terrainName = terrainName;
	this.layer = layer;
	this.svgSegment = svgSegment;
	this.subTerrainSvgArray = new Array();
}

function gkSubTerrainSvgArrayEntryDef(terrainName, fillCount, layer, svgSegment) {
	this.terrainName = terrainName;
	this.fillCount = fillCount;
	this.layer = layer;
	this.svgSegment = svgSegment;
}

function gkTerrainInit() {
}

// called as a request from the server
// to set the entire pod terran map
function gkSetTerrainMap(jsonData) {
//console.log("gkSetTerrainMap");
	var i;

console.log("tileList.length: " + jsonData.tileList.length);
	for (i = 0;i < jsonData.tileList.length; i++) {
		var x, y, z, terrainName;
		x = jsonData.tileList[i].x;
		y = jsonData.tileList[i].y;
		z = jsonData.tileList[i].z;
		terrainName = jsonData.tileList[i].terrain;

		var terrainMapMapEntry = new gkTerrainMapMapEntryDef(x, y, z, terrainName)

		var mapKey = gkTerrainGetMapKey(x, y);
		gkTerrainContext.terrainMapMap[mapKey] = terrainMapMapEntry;
	}

console.log("objectList.length: " + jsonData.objectList.length);
	for (i = 0;i < jsonData.objectList.length; i++) {
		var x, y, z, objectName;
		x = jsonData.objectList[i].x;
		y = jsonData.objectList[i].y;
		z = jsonData.objectList[i].z;
		objectName = jsonData.objectList[i].object;

		var objectMapMapEntry = new gkTerrainMapMapEntryDef(x, y, z, objectName)

		var mapKey = gkTerrainGetMapKey(x, y);
		gkTerrainContext.terrainMapMap[mapKey] = objectMapMapEntry;
	}

	gkViewRender();
}

// called as a request from the server
// to set all the svg files for the require terrain
// this function assumes that the terrain entry is processed
// before any of that terrain's sub terrain entry
function gkSetTerrainSvg(jsonData, rawSvgData) {
//console.log("gkSetTerrainSvg");
	if (jsonData.terrain != undefined) {
		var terrainName, layer, svgSegment

		terrainName = jsonData.terrain
		layer = jsonData.layer
		var terrainSvgMapEntry = new gkTerrainSvgMapEntryDef(terrainName, layer, rawSvgData)
		gkTerrainContext.terrainSvgMap[jsonData.terrain] = terrainSvgMapEntry
	}

	if (jsonData.subTerrain != undefined) {
		var terrainName, layer, svgSegment

		terrainName = jsonData.subTerrain
		fillCount = jsonData.fillCount
		layer = jsonData.layer

		var subTerrainSvgArrayEntry = new gkSubTerrainSvgArrayEntryDef(terrainName, fillCount, layer, rawSvgData)

		gkTerrainContext.terrainSvgMap[terrainName].subTerrainSvgArray.push(subTerrainSvgArrayEntry);
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

