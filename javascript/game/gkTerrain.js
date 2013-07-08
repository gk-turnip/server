
// this is the terrain handler

var gkTerrainContext = new gkTerrainContextDef();

function gkTerrainContextDef() {
	// map holding the position of each terrain diamond
	this.terrainMapMap = new Object();

	// map holding the svg data for each terrain diamond
	this.terrainSvgMap = new Object();

	// map holding the places an avatar cannot go
	this.terrainWallMap = new Object();

	// array holding the enviromental audio sources
	this.terrainAudioMap = new Array();

	this.terrainDiamondOffsetX = 50;
	this.terrainDiamondOffsetY = 0;

	this.terrainUndefinedZ = -30000;

	this.moveMarker = null;
}

function gkTerrainMapMapEntryDef(x, y, zList, terrainName) {
	this.x = x
	this.y = y
	this.zList = zList
	if (terrainName != undefined) {
		this.terrainName = terrainName
	}
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

function gkTerrainClearTerrainGridListLayer() {
	var layer = document.getElementById("gkTerrainGridListLayer");

	gkTerrainClearLayer(layer);
}

function gkTerrainClearLayer(layer) {
	while (layer.firstChild) {
		layer.removeChild(layer.firstChild);
	}
}

function gkTerrainClearTerrain(jsonData) {
	// clear out old terrain map
	gkFieldRemoveOtherAvatars();
	gkTerrainClearMoveMarker();
	gkTerrainClearTerrainBaseLayer();
	gkTerrainClearTerrainDandelionLayer();
	gkTerrainClearTerrainObjectLayer();
	gkTerrainClearTerrainGridListLayer();
	gkTerrainContext.terrainMapMap = new Object();
	gkTerrainContext.terrainSvgMap = new Object();
	gkTerrainContext.terrainWallMap = new Object();
	gkTerrainContext.terrainAudioMap = new Array();
}

// called as a request from the server
// to set the entire pod terran map
// the terrain svg must be done before the terrain map
function gkTerrainSetTerrainMap(jsonData) {
//console.log("gkTerrainSetTerrainMap");
	var i;

	var baseLayer = document.getElementById("gkTerrainBaseLayer");

console.log("tileList.length: " + jsonData.tileList.length);
	// list of unique grid list names
	var gridListMap = new Object();

	for (i = 0;i < jsonData.tileList.length; i++) {
		var x, y, z, zList, terrainName;
		x = parseInt(jsonData.tileList[i].x);
		y = parseInt(jsonData.tileList[i].y);
		if (jsonData.tileList[i].t != undefined) {
			terrainName = jsonData.tileList[i].t;
		} else {
			terrainName = undefined;
		}

		var gridListIndexName = gkIsoGetGridListIndexName(x,y,0)
		gridListMap[gridListIndexName] = gridListIndexName;

		zList = jsonData.tileList[i].z;
		
		var terrainMapMapEntry = new gkTerrainMapMapEntryDef(x, y, zList, terrainName)

		var mapKey = gkTerrainGetMapKey(x, y);
		gkTerrainContext.terrainMapMap[mapKey] = terrainMapMapEntry;
	}

	var gridListArray = new Array();
	for (prop in gridListMap) {
		gridListArray.push(prop);
	}

console.log("grid array size: " + gridListArray.length);
	// grid list must be added in order
	gridListArray = gridListArray.sort();
	for (var i = 0;i < gridListArray.length;i++) {
		gkFieldAddGridListEntry(gridListArray[i]);
	}

	for (i = 0;i < jsonData.tileList.length; i++) {
		var x, y, z, zList, terrainName;
		x = parseInt(jsonData.tileList[i].x);
		y = parseInt(jsonData.tileList[i].y);
		if (jsonData.tileList[i].t != undefined) {
			terrainName = jsonData.tileList[i].t;
		} else {
			terrainName = undefined;
		}

		zList = jsonData.tileList[i].z;

		if (terrainName != undefined) {		
			var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(terrainName);

			if (terrainSvgMapEntry != undefined) {
				for (var j = 0;j < zList.length;j++) {
					var isoXYZ = new GkIsoXYZDef(x,y,zList[j]);

					gkFieldAddTerrainObject(gkFieldContext.defsTerrainPrefix, gkFieldContext.terrainObjectPrefix + i, terrainName, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY, terrainSvgMapEntry.originZ);
				}
			}
		}
	}

console.log("objectList.length: " + jsonData.oList.length);
	for (i = 0;i < jsonData.oList.length; i++) {
		var x, y, z, objectName, podId, destination;
		x = jsonData.oList[i].x;
		y = jsonData.oList[i].y;
		z = jsonData.oList[i].z;
		if (jsonData.oList[i].podId != undefined) {
			podId = jsonData.oList[i].podId;
			destination = jsonData.oList[i].destination;
		}
		objectName = jsonData.oList[i].o;

		var objectLayer = document.getElementById("gkTerrainObjectLayer");

		var isoXYZ = new GkIsoXYZDef(x,y,z);
		var terrainSvgMapEntry = gkTerrainGetSvgMapEntry(objectName);

		if (terrainSvgMapEntry != undefined) {
			gkFieldAddObjectToGridList(gkFieldContext.defsTerrainPrefix, gkFieldContext.terrainTilePrefix + i, objectName, isoXYZ, terrainSvgMapEntry.originX, terrainSvgMapEntry.originY, terrainSvgMapEntry.originZ, podId, destination);
		}
	}

	for (i = 0;i < jsonData.audioList.length; i++) {
		var clip, x, y, z;
		clip = jsonData.audioList[i].clip;
		x = jsonData.audioList[i].x;
		y = jsonData.audioList[i].y;
		z = jsonData.audioList[i].z;

		var audioMapEntry = new gkTerrainAudioMapEntryDef(clip, x, y, z);
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

function gkTerrainTestMoveElevation(x, y, z, maxOffset) {
	testMove = new Object();
	testMove.canMove = false;

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

	localX = gkIsoGetFernFromDecifern(localX);
	localY = gkIsoGetFernFromDecifern(localY);

	var mapKey = gkTerrainGetMapKey(localX, localY);

	var terrainWallEntry = gkTerrainContext.terrainWallMap[mapKey];
	if (terrainWallEntry == undefined) {
		var terrainMapMapEntry = gkTerrainContext.terrainMapMap[mapKey];
		if (terrainMapMapEntry != undefined) {
			for (var i = 0;i < terrainMapMapEntry.zList.length;i++) {
				if (Math.abs(z - terrainMapMapEntry.zList[i]) <= maxOffset) {
					testMove.canMove = true;
					testMove.z = terrainMapMapEntry.zList[i];
					break;
				}
			}
		}
	}

	return testMove;
}

function gkTerrainSetSvgObjectOnClick(ref, objectName, isoXYZ, originX, originY, originZ, podId, destination) {
	ref.onclick = function() { gkTerrainSvgObjectClick(objectName, isoXYZ.x, isoXYZ.y, isoXYZ.z, originX, originY, originZ, podId, destination) };
}

function gkTerrainSvgObjectClick(id, x, y, z, originX, originY, originZ, podId, destination) {
	console.log("svgObjectClick id: " + id + " xyz: " + x + "," + y + "," + z + " origin: " + originX + "," + originY);

	if (podId != undefined) {
		var x = destination.x;
		var y = destination.y;
		var z = destination.z;
		
		x = x + (((Math.floor(Math.random() * 2) * 2) - 1) * ((Math.floor(Math.random() * 20) + 1)));
		y = y + (((Math.floor(Math.random() * 2) * 2) - 1) * ((Math.floor(Math.random() * 20) + 1)));

console.log("destination: " + x + "," + y + "," + z);

		gkWsSendMessage("newPodReq~{ \"podId\":\"" + podId + "\", \"x\":\"" + x + "\", \"y\":\"" + y + "\", \"z\": \"" + z + "\" }~");

	}
}

// called as a request from the server
// to set all the svg files for the require terrain
function gkTerrainSetTerrainSvg(jsonData, rawSvgData) {
	if (jsonData.terrain != undefined) {
		var terrainName, originX, originY, originZ, layer

		terrainName = jsonData.terrain;
		originX = jsonData.originX;
		originY = jsonData.originY;
		originZ = jsonData.originZ;
		layer = jsonData.layer;
		var terrainSvgMapEntry = new gkTerrainSvgMapEntryDef(terrainName, originX, originY, originZ, layer);
		gkTerrainContext.terrainSvgMap[jsonData.terrain] = terrainSvgMapEntry;

		var gkDefs = document.getElementById("gkDefs");

		g = gkIsoCreateSvgObject(rawSvgData);
		g.setAttribute("id",gkFieldContext.defsTerrainPrefix + jsonData.terrain);

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

	svgMapEntry = gkTerrainContext.terrainSvgMap[terrainName]

	return svgMapEntry
}

