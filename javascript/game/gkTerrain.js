
// this is the terrain handler

var gkTerrainContext = new gkTerrainContextDef();

function gkTerrainContextDef() {
	// map holding the position of each terrain diamond
	this.terrainMapMap = new Object();

	// map holding the svg data for each terrain diamond
	this.terrainSvgMap = new Object();

	// map holding the elevation data for each terrain diamond
	this.terrainElevationMap = new Object();

	// map holding the places an avatar cannot go
	this.terrainWallMap = new Object();

	// array holding the enviromental audio sources
	this.terrainAudioMap = new Array();

	this.terrainDiamondOffsetX = 50;
	this.terrainDiamondOffsetY = 0;

	this.terrainUndefinedZ = -30000;

	this.moveMarker = null;
}

function gkTerrainMapMapEntryDef(x, y, z, terrainName) {
	this.x = x
	this.y = y
	this.z = z
	this.terrainName = terrainName
}

function gkTerrainElevationMapEntryDef(x, y, z) {
	this.x = x
	this.y = y
	this.z = z
}

function gkTerrainWallMapEntryDef(x, y, z) {
	this.x = x
	this.y = y
	this.z = z
}

function gkTerrainSvgMapEntryDef(terrainName, originX, originY, originZ, layer) {
	this.terrainName = terrainName;
	this.originX = originX;
	this.originY = originY;
	this.originZ = originZ;
	this.layer = layer;
}

function gkTerrainAudioMapEntryDef(clip, x, y, z) {
	this.clip = clip;
	this.x = x;
	this.y = y;
	this.z = z;
}

function gkTerrainInit() {
}

function gkTerrainClearTerrainBaseLayer() {
	var layer = document.getElementById("gkTerrainBaseLayer");

	gkTerrainClearLayer(layer);
}

function gkTerrainClearTerrainDandelionLayer() {
	var layer = document.getElementById("gkTerrainDandelionLayer");

	gkTerrainClearLayer(layer);
}

function gkTerrainClearTerrainObjectLayer() {
	var layer = document.getElementById("gkTerrainObjectLayer");

	gkTerrainClearLayer(layer);
}

function gkTerrainClearLayer(layer) {
	while (layer.firstChild) {
		layer.removeChild(layer.firstChild);
	}
}

// called as a request from the server
// to set the entire pod terran map
// the terrain svg must be done before the terrain map
function gkTerrainSetTerrainMap(jsonData) {
//console.log("gkTerrainSetTerrainMap");
	var i;

	// clear out old terrain map
	gkTerrainClearMoveMarker();
	gkTerrainClearTerrainBaseLayer();
	gkTerrainClearTerrainDandelionLayer();
	gkTerrainClearTerrainObjectLayer();
	this.terrainMapMap = new Object();

	var baseLayer = document.getElementById("gkTerrainBaseLayer");

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

		var isoXYZ = new GkIsoXYZDef(x,y,z);
		var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(terrainName);

		if (terrainSvgMapEntry != undefined) {
			var ref = document.createElementNS(gkIsoContext.svgNameSpace,"use");
			ref.setAttributeNS(gkIsoContext.xlinkNameSpace,"href","#t_" + terrainName);
			gkIsoSetSvgUsePositionWithOffset(ref, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY, terrainSvgMapEntry.originZ);
			baseLayer.appendChild(ref);
		}
	}

console.log("objectList.length: " + jsonData.oList.length);
	for (i = 0;i < jsonData.oList.length; i++) {
		var x, y, z, objectName, podId;
		x = jsonData.oList[i].x;
		y = jsonData.oList[i].y;
		z = jsonData.oList[i].z;
		if (jsonData.oList[i].podId != undefined) {
			podId = jsonData.oList[i].podId;
		}
		objectName = jsonData.oList[i].o;

		var objectMapMapEntry = new gkTerrainMapMapEntryDef(x, y, z, objectName)

//		var mapKey = gkTerrainGetMapKey(x, y);
//		gkTerrainContext.terrainMapMap[mapKey] = objectMapMapEntry;

		var objectLayer = document.getElementById("gkTerrainObjectLayer");

		var isoXYZ = new GkIsoXYZDef(x,y,z);
		var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(objectName);

		if (terrainSvgMapEntry != undefined) {
			var ref = document.createElementNS(gkIsoContext.svgNameSpace,"use");
			ref.setAttributeNS(gkIsoContext.xlinkNameSpace,"href","#t_" + objectName);
			gkIsoSetSvgUsePositionWithOffset(ref, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY, terrainSvgMapEntry.originZ);

			gkTerrainSetSvgObjectOnClick(ref, objectName, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY, terrainSvgMapEntry.originZ, podId);

			objectLayer.appendChild(ref);
		}
	}

	for (i = 0;i < jsonData.elevationList.length; i++) {
		var x, y, z;

		x = jsonData.elevationList[i].x;
		y = jsonData.elevationList[i].y;
		z = jsonData.elevationList[i].z;

		var mapKey = gkTerrainGetMapKey(x, y);

		var elevationMapEntry = new gkTerrainElevationMapEntryDef(x, y, z);

		gkTerrainContext.terrainElevationMap[mapKey] = elevationMapEntry;
	}

	for (i = 0;i < jsonData.audioList.length; i++) {
		var clip, x, y, z;
		clip = jsonData.audioList[i].clip;
		x = jsonData.audioList[i].x;
		y = jsonData.audioList[i].y;
		z = jsonData.audioList[i].z;

		var audioMapEntry = new gkTerrainAudioMapEntryDef(clip, x, y, z);

		//gkTerrainContext.terrainAudioMap.append(audioMapEntry);
	}

	for (i = 0;i < jsonData.wallList.length; i++) {
		var x, y, z;

		x = jsonData.wallList[i].x;
		y = jsonData.wallList[i].y;
		z = jsonData.wallList[i].z;

		var mapKey = gkTerrainGetMapKey(x, y);
		var wallMapEntry = new gkTerrainWallMapEntryDef(x, y, z);
		gkTerrainContext.terrainWallMap[mapKey] = wallMapEntry;
	}

	gkViewRender();
}

function gkTerrainClearMoveMarker() {
	if (gkTerrainContext.moveMarker != null) {
		var objectLayer = document.getElementById("gkTerrainObjectLayer");
		objectLayer.removeChild(gkTerrainContext.moveMarker);
		gkTerrainContext.moveMarker = null;
	}
}

function gkTerrainSetMoveMarker(g) {
	var objectLayer = document.getElementById("gkTerrainObjectLayer");
	objectLayer.appendChild(g);

	gkTerrainContext.moveMarker = g;
}

function gkTerrainGetElevation1(x, y) {
	var localX;
	var localY;
	if (x.substring) {
		localX = parseInt(x);
	} else {
		localX = x;
	}
	if (y.substring) {
		localY = parseInt(y);
	} else {
		localY = y;
	}

	localX = localX / 10;
	localY = localY / 10;

	localX = Math.floor(localX);
	localY = Math.floor(localY);

	localX = localX * 10;
	localY = localY * 10;

	var z = gkTerrainContext.terrainUndefinedZ;

	var mapKey = gkTerrainGetMapKey(localX, localY);

	var terrainWallEntry = gkTerrainContext.terrainWallMap[mapKey];
	if (terrainWallEntry == undefined) {
		var terrainMapMapEntry = gkTerrainContext.terrainMapMap[mapKey];
		if (terrainMapMapEntry != undefined) {
			z = terrainMapMapEntry.z;
		}
	}

	return z;
}

function gkTerrainGetElevation2(x, y) {
	var localX;
	var localY;
	if (x.substring) {
		localX = parseInt(x);
	} else {
		localX = x;
	}
	if (y.substring) {
		localY = parseInt(y);
	} else {
		localY = y;
	}

	localX = localX / 10;
	localY = localY / 10;

	localX = Math.floor(localX);
	localY = Math.floor(localY);

	localX = localX * 10;
	localY = localY * 10;

	var mapKey = gkTerrainGetMapKey(localX, localY);

	var elevationMapEntry;

	elevationMapEntry = gkTerrainContext.terrainElevationMap[mapKey];

	var z = gkTerrainContext.terrainUndefinedZ;

	if (elevationMapEntry != undefined) {
		z = elevationMapEntry.z;
	}

	return z;
}

function gkTerrainSetSvgObjectOnClick(ref, objectName, isoXYZ, originX, originY, originZ, podId) {
	ref.onclick = function() { gkTerrainSvgObjectClick(objectName, isoXYZ.x, isoXYZ.y, isoXYZ.z, originX, originY, originZ, podId) };
}

function gkTerrainSvgObjectClick(id, x, y, z, originX, originY, originZ, podId) {
	console.log("svgObjectClick id: " + id + " xyz: " + x + "," + y + "," + z + " origin: " + originX + "," + originY);

	if (podId != undefined) {
console.log("podId: " + podId);
		gkWsSendMessage("newPodReq~{ \"podId\":\"" + podId + "\" }~");
		var a = new GkIsoXYZDef(0, 0, 0);
		gkFieldSetNewAvatarDestination(a)
	}
}

// called as a request from the server
// to set all the svg files for the require terrain
function gkTerrainSetTerrainSvg(jsonData, rawSvgData) {
//console.log("gkTerrainSetTerrainSvg");
	if (jsonData.terrain != undefined) {
		var terrainName, originX, originY, originZ, layer;

//console.log("gkTerrainSetTerrainSvg name: " + jsonData.terrain);
		terrainName = jsonData.terrain;
		originX = jsonData.originX;
		originY = jsonData.originY;
		originZ = jsonData.originZ;
		layer = jsonData.layer;
		var terrainSvgMapEntry = new gkTerrainSvgMapEntryDef(terrainName, originX, originY, originZ, layer);
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

function gkTerrainGetSvgMapEntry(terrainName) {
//	var mapKey = gkTerrainGetMapKey(x, y);
//
//	var terrainMapMapEntry = gkTerrainContext.terrainMapMap[mapKey];
//	var svgMapEntry
//
//	if (terrainMapMapEntry == undefined) {
//		console.error("missing terrainMapMapEntry mapKey: " + mapKey);
//	} else {
//	var terrainName = terrainMapMapEntry.terrainName

	svgMapEntry = gkTerrainContext.terrainSvgMap[terrainName]
//	}

	return svgMapEntry
}

