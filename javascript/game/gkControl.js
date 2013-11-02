
// this program handles the in-game controls

var gkControlContext = new gkControlContextDef();

function gkControlContextDef() {
	this.controlLayer = "gkControlLayer"
	this.controlUrlPrefix = "/assets/gk/controls/"
	this.loadMap = new Object();
	this.menuMap = null;
	this.menuItemHeight = 50;
	this.menuWidth = 50;
	this.menuStack = new Array();
	this.mouseDown = false;
	this.elevationFern = 0;
	this.elevationDeciFern = 0;
}

function gkControlInit() {
	//gkControlLoad("start", gkControlHandleLoadStart);
	gkControlLoadMenuMap();
}

function gkControlLoadMenuMap() {
	var xmlhttp = new XMLHttpRequest();
	var fullUrl = gkControlContext.controlUrlPrefix + "menuMap.json";
console.log("fullUrl: " + fullUrl);

	xmlhttp.onreadystatechange=function() {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				gkControlHandleLoadMenuMap(xmlhttp.responseText);
			} else {
				console.error("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
			}
		}
	}

	xmlhttp.open("GET", fullUrl, true);
	xmlhttp.send();
}

function gkControlHandleLoadMenuMap(menuMapText) {
//console.log("menuMapText: " + menuMapText);
	gkControlContext.menuMap = JSON.parse(menuMapText, null);

	gkControlContext.menuStack.push("menu");
	for (var i = 0;i < gkControlContext.menuMap.menu.length;i++) {
		gkControlLoad(gkControlContext.menuMap.menu[i].display, i, gkControlHandleLoadMenuItem);
	}
}

function gkControlHandleLoadMenuItem(controlId, index) {
	//console.log("got menu item control loaded controlId: " + controlId + " index: " + index);
	//console.log(gkControlContext.loadMap[controlId].rawSvg);
	//console.log(gkControlContext.loadMap[controlId].rawJson);
	gkControlAddSvg(controlId, index);
}

function gkControlLoad(controlId, index, controlFunction) {
	gkControlContext.loadMap[controlId] = new Object();
	gkControlContext.loadMap[controlId].controlFunction = controlFunction;
	gkControlContext.loadMap[controlId].svgLoaded = false;
	gkControlContext.loadMap[controlId].jsonLoaded = false;

	gkControlLoadSvg(controlId, index, controlFunction);
	gkControlLoadJson(controlId, index, controlFunction);
}

function gkControlLoadSvg(controlId, index, controlFunction) {
	var xmlhttp = new XMLHttpRequest();
	var fullUrl = gkControlContext.controlUrlPrefix + controlId + ".svg";

	xmlhttp.onreadystatechange=function() {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				gkControlHandleLoadSvg(xmlhttp.responseText, index, controlId);
			} else {
				console.error("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
			}
		}
	}

	xmlhttp.open("GET", fullUrl, true);
	xmlhttp.send();
}

function gkControlLoadJson(controlId, index, controlFunction) {
	var xmlhttp = new XMLHttpRequest();
	var fullUrl = gkControlContext.controlUrlPrefix + controlId + ".json";

	xmlhttp.onreadystatechange=function() {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				gkControlHandleLoadJson(xmlhttp.responseText, index, controlId);
			} else {
				console.error("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
			}
		}
	}

	xmlhttp.open("GET", fullUrl, true);
	xmlhttp.send();
}

function gkControlHandleLoadSvg(rawSvg, index, controlId) {
	gkControlContext.loadMap[controlId].rawSvg = rawSvg;
	gkControlContext.loadMap[controlId].svgLoaded = true;
	gkControlCheckLoaded(controlId, index);
}

function gkControlHandleLoadJson(rawJson, index, controlId) {
	gkControlContext.loadMap[controlId].rawJson = rawJson;
	gkControlContext.loadMap[controlId].jsonLoaded = true;
	gkControlCheckLoaded(controlId, index);
}

function gkControlCheckLoaded(controlId, index) {
	if ((gkControlContext.loadMap[controlId].svgLoaded) && (gkControlContext.loadMap[controlId].jsonLoaded)) {
		gkControlContext.loadMap[controlId].controlFunction(controlId, index);
	}
}

/*
function gkControlHandleLoadStart(controlId) {
	console.log("got start control loaded");
	console.log(gkControlContext.loadMap[controlId].rawSvg);
	console.log(gkControlContext.loadMap[controlId].rawJson);
	gkControlAddSvg(controlId);
}
*/

function gkControlAddSvg(controlId, index) {
	var g = gkIsoCreateSvgObject(gkControlContext.loadMap[controlId].rawSvg);
	g.setAttribute("id", controlId);
	g.setAttribute("transform","translate(" + 0 + "," + (index * gkControlContext.menuItemHeight) + ")");
	g.onclick = function() {
		gkControlMenuItemClick(controlId, index);
	};

console.log("need to add mouse events to: " + controlId);
	if ((controlId == "widthHeightPad") || (controlId == "zoomPad") || (controlId == "panPad") || (controlId == "backgroundVolumePad") || (controlId == "effectsVolumePad") || (controlId == "terrainElevationPad") || (controlId == "terrainAttributePad")) {
console.log("adding mouse event handler for control: " + controlId);
		g.onmousedown = function(evt) {
			gkControlMenuItemMouseDown(evt, controlId, index);
		};
		g.onmouseup = function(evt) {
			gkControlMenuItemMouseUp(evt, controlId, index);
		};
		g.onmousemove = function(evt) {
			gkControlMenuItemMouseMove(evt, controlId, index);
		};
	}
	var layer = document.getElementById(gkControlContext.controlLayer);
	layer.appendChild(g);
	//console.log("added to layer: " + gkControlContext.controlLayer);

	//console.log("gkControlAddSvg controlId: " + controlId);
	if (controlId == "widthHeightPad") {
		gkControlSetWidthHeightPadText();
	}
	if (controlId == "zoomPad") {
		gkControlSetZoomPadText();
	}
	if (controlId == "panPad") {
		gkControlSetPanPadText();
	}
	if (controlId == "backgroundVolumePad") {
		gkControlSetBackgroundVolumePadText();
	}
	if (controlId == "effectsVolumePad") {
		gkControlSetEffectsVolumePadText();
	}
	if (controlId == "terrainAttibutePad") {
		gkControlSetAttributePadText();
	}

	if (controlId == "start") {
		var podNameText = document.createElementNS(gkIsoContext.svgNameSpace,"text");
		podNameText.setAttribute("id", "menuPodNameText");
		podNameText.setAttribute("transform","translate(" + 0 + "," + (((index + 1) * gkControlContext.menuItemHeight) + 25 ) + ")");
		podNameText.setAttribute("font-family","sans-serif");
		podNameText.setAttribute("font-wize","24");
		layer.appendChild(podNameText);
console.log("id: " + podNameText.id);
		gkControlSetPodTitle();

		var podPosText = document.createElementNS(gkIsoContext.svgNameSpace,"text");
		podPosText.setAttribute("id", "menuPodPosText");
		podPosText.setAttribute("transform","translate(" + 0 + "," + (((index + 2) * gkControlContext.menuItemHeight) + 25 ) + ")");
		podPosText.setAttribute("font-family","sans-serif");
		podPosText.setAttribute("font-wize","24");
		layer.appendChild(podPosText);
		gkControlSetPos();

		var podFPSText = document.createElementNS(gkIsoContext.svgNameSpace,"text");
		podFPSText.setAttribute("id", "menuPodFPSText");
		podFPSText.setAttribute("transform","translate(" + 0 + "," + (((index + 3) * gkControlContext.menuItemHeight) + 25 ) + ")");
		podFPSText.setAttribute("font-family","sans-serif");
		podFPSText.setAttribute("font-wize","24");
		layer.appendChild(podFPSText);
		gkControlSetFPS();

		var inFocusText = document.createElementNS(gkIsoContext.svgNameSpace,"text");
		inFocusText.setAttribute("id", "menuInFocusText");
		inFocusText.setAttribute("transform","translate(" + 0 + "," + (((index + 4) * gkControlContext.menuItemHeight) + 25 ) + ")");
		inFocusText.setAttribute("font-family","sans-serif");
		inFocusText.setAttribute("font-wize","24");
		layer.appendChild(inFocusText);
		gkControlSetInFocus();

	}
}

function gkControlSetPodTitle() {
	var podNameText = document.getElementById("menuPodNameText");
	if (podNameText != undefined) {
		podNameText.textContent = "current pod: " + gkFieldContext.podTitle;
	}
}

function gkControlSetPos() {
	if (gkFieldContext.currentPosX != undefined) {
		var podPosText = document.getElementById("menuPodPosText");
		if (podPosText != undefined) {
			podPosText.textContent = "position: " + gkFieldContext.currentPosX + "," + gkFieldContext.currentPosY + "," + gkFieldContext.currentPosZ;
		}
	}
}

function gkControlSetFPS() {
	var podFPSText = document.getElementById("menuPodFPSText");
	if (podFPSText != undefined) {
		podFPSText.textContent = "FPS: " + gkFieldContext.frameRate;
	}
}

function gkControlSetInFocus() {
	var inFocusText = document.getElementById("menuInFocusText");
	if (inFocusText != undefined) {
		if (gkFieldContext.inFocus) {
			inFocusText.textContent = "in focus";
		} else {
			inFocusText.textContent = "not in focus";
		}
	}
}

function gkControlMenuItemClick(controlId, index) {
	console.log("gkControlMenuItemClick: " + controlId + " " + index);

	var nextLevelControlId = controlId;

	if (controlId == "addTile") {
		gkControlHandleAddTileSelect();
	} else {
		if (controlId == "removeTile") {
			gkControlHandleRemoveTileSelect();
		} else {
			if (controlId == "terrainSaveEdit") {
				gkControlHandleTerrainSaveEdit();
			} else {
				if (controlId == "close") {
					gkControlContext.mouseDown = false;

					gkControlHandleClose(gkControlContext.menuStack[gkControlContext.menuStack.length - 1]);
					gkControlContext.menuStack.pop();
					nextLevelControlId = gkControlContext.menuStack[gkControlContext.menuStack.length - 1];
				} else {
					if (gkControlContext.menuMap[nextLevelControlId] != undefined) {
						gkControlContext.menuStack.push(controlId);
					}
				}

				//console.log("new menu controlId: " + nextLevelControlId);
				if (gkControlContext.menuMap[nextLevelControlId] != undefined) {
					gkControlClearCurrentMenu();
					for (var i = 0;i < gkControlContext.menuMap[nextLevelControlId].length;i++) {
						gkControlLoad(gkControlContext.menuMap[nextLevelControlId][i].display, i, gkControlHandleLoadMenuItem);
					}
				}
			}
		}
	}

	if (controlId == "terrainEdit") {
		gkControlHandleTerrainEditSelect();
	}
}

function gkControlHandleClose(closeMenu) {
	console.log("closeMenu: " + closeMenu);

	if (closeMenu == "zoom") {
		gkControlHandleCloseZoom();
	} else {
		if (closeMenu == "widthHeight") {
			gkControlHandleCloseWidthHeight();
		} else {
			if (closeMenu == "backgroundVolume") {
				gkControlHandleCloseBackgroundVolume();
			} else {
				if (closeMenu == "effectsVolume") {
					gkControlHandleCloseEffectsVolume();
				} else {
					if (closeMenu == "terrainTileEdit") {
						gkControlHandleCloseTerrainEdit();
					}
				}
			}
		}
	}
}

function gkControlHandleCloseZoom() {
	gkControlSendUserPref("screenZoom",gkViewContext.scale);
}

function gkControlHandleCloseWidthHeight() {
	gkControlSendUserPref("screenWidth",gkViewContext.svgWidth);
	gkControlSendUserPref("screenHeight",gkViewContext.svgHeight);
}

function gkControlHandleCloseBackgroundVolume() {
	var volumeLevel = gkAudioGetVolume(gkAudioContext.backgroundVolumeSelect);
	gkControlSendUserPref("backgroundVolume", volumeLevel);
}

function gkControlHandleCloseEffectsVolume() {
	var volumeLevel = gkAudioGetVolume(gkAudioContext.effectsVolumeSelect);
	gkControlSendUserPref("effectsVolume", volumeLevel);
}

function gkControlSendUserPref(prefName, prefValue) {
	console.log("sending " + prefName + " " + prefValue);
	gkWsSendMessage("userPrefSaveReq~{\"prefName\":\"" + prefName + "\",\"prefValue\":\"" + prefValue + "\"}~");
}

function gkControlMenuItemMouseDown(evt, controlId, index) {
	evt.preventDefault();
	console.log("gkControlMenuItemMouseDown " + controlId + " " + index);
	gkControlContext.mouseDown = true;
	gkControlMenuItemMouseMove(evt, controlId, index);
}

function gkControlMenuItemMouseUp(evt, controlId, index) {
	evt.preventDefault();
	console.log("gkControlMenuItemMouseUp " + controlId + " " + index);
	gkControlContext.mouseDown = false;
}

function gkControlMenuItemMouseMove(evt, controlId, index) {
	evt.preventDefault();
	console.log("gkControlMenuItemMouseMove " + controlId + " " + index);

	if (gkControlContext.mouseDown) {
		var x,y;

		x = evt.clientX - gkViewContext.marginX;
		y = evt.clientY - gkViewContext.marginY;

		y -= 50 * index;

		x = x / 200;
		y = y / 200;

		if (x < 0) {
			x = 0;
		}
		if (x > 1) {
			x = 0.99;
		}
		if (y < 0) {
			y = 0;
		}
		if (y > 1) {
			y = 0.99;
		}

		console.log("gkControlMenuItemMouseMove " + x + "," + y + " " + controlId + " " + index);

		if (controlId == "widthHeightPad") {
			gkControlHandleWidthHeightPad(x, y);
		}
		if (controlId == "zoomPad") {
			gkControlHandleZoomPad(x);
		}
		if (controlId == "panPad") {
			gkControlHandlePanPad(x, y);
		}
		if (controlId == "backgroundVolumePad") {
			gkControlHandleBackgroundVolumePad(x);
		}
		if (controlId == "effectsVolumePad") {
			gkControlHandleEffectsVolumePad(x);
		}
		if (controlId == "terrainElevationPad") {
			gkControlHandleTerrainElevationPad(x, y);
		}
		if (controlId == "terrainAttributePad") {
			gkControlHandleTerrainAttributePad(x);
		}
	}
}

function gkControlHandleWidthHeightPad(x, y) {
	var width, height

//console.log("x,y: " + x + "," + y);
	width = Math.floor(300 + (x * 2260));
	height = Math.floor(300 + (y * 2260));

	gkViewContext.svgWidth = width;
	gkViewContext.svgHeight = height;

	gkControlSetWidthHeightPadText();

	gkViewRender();

	//console.log("new width: " + width + " height: " + height);
}

function gkControlHandleZoomPad(x) {
	var zoomLevel;

console.log("zoom x: " + x);

	zoomLevel = ((gkViewContext.highScale - gkViewContext.lowScale) * x) + gkViewContext.lowScale;

	zoomLevel = Math.floor(zoomLevel * 10) / 10;

	gkViewContext.scale = zoomLevel;

	gkControlSetZoomPadText();

	gkViewRender();
	//console.log("new zoom level: " + zoomLevel);
}

function gkControlHandlePanPad(x, y) {
	var offsetX, offsetY

console.log("x,y: " + x + "," + y);
console.log("view: " + gkViewContext.viewOffsetIsoXYZ.x + "," + gkViewContext.viewOffsetIsoXYZ.y);

//	var WinXY = gkViewContext.viewOffsetIsoXYZ.convertToWin();

	offsetX = Math.floor((x * (gkViewContext.highPanX - gkViewContext.lowPanX)) + gkViewContext.lowPanX);
	offsetY = Math.floor((y * (gkViewContext.highPanY - gkViewContext.lowPanY)) + gkViewContext.lowPanY);

	var newOffsetWinXY = new GkWinXYDef(offsetX, offsetY);
	gkViewContext.viewOffsetIsoXYZ = newOffsetWinXY.convertToIso(0);

//	gkViewContext.viewOffsetIsoXYZ.x = offsetX;
//	gkViewContext.viewOffsetIsoXYZ.y = offsetY;

//	width = Math.floor(300 + (x * 2260));
//	height = Math.floor(300 + (y * 2260));

//	gkViewContext.svgWidth = width;
//	gkViewContext.svgHeight = height;

	gkControlSetPanPadText();

	gkViewRender();

console.log("new view: " + gkViewContext.viewOffsetIsoXYZ.x + "," + gkViewContext.viewOffsetIsoXYZ.y);
}

function gkControlHandleBackgroundVolumePad(x) {
	var volumeLevel;

	console.log("new background volume x: " + x);
	volumeLevel = x;

	gkAudioVolumeChange(gkAudioContext.backgroundVolumeSelect, volumeLevel);

	gkControlSetBackgroundVolumePadText();

	console.log("new background volume level: " + volumeLevel);
}

function gkControlHandleEffectsVolumePad(x) {
	var volumeLevel;

	console.log("new effects volume x: " + x);
	volumeLevel = x;

	gkAudioVolumeChange(gkAudioContext.effectsVolumeSelect, volumeLevel);

	gkControlSetEffectsVolumePadText();

	console.log("new effects volume level: " + volumeLevel);
}

function gkControlHandleTerrainElevationPad(x, y) {
	console.log("gkControlHandleTerrainElevation: " + x + "," + y);

	gkControlSetElevationPadText();
}

function gkControlHandleTerrainAttributePad(x) {
	var attributeIndex;

console.log("attribute x: " + x);

	attributeIndex = x * 2;

	attributeIndex = Math.floor(attributeIndex);

	gkTerrainEditSetAttributeIndex(attributeIndex);

	gkControlSetAttributePadText();

	gkViewRender();

	console.log("new attributeIndex: " + attributeIndex);
}

function gkControlSetZoomPadText() {
	var zoomText = document.getElementById("zoomPad_zoomText");
	zoomText.textContent = "zoom: " + gkViewContext.scale;

	var zoomRect = document.getElementById("zoomPad_zoomRect");
//	var transX = (gkViewContext.scale - gkViewContext.lowScale) * 133;

	var transX = (gkViewContext.scale - gkViewContext.lowScale) / (gkViewContext.highScale - gkViewContext.lowScale) * 195;
	zoomRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	//console.log("transX: " + transX);
}

function gkControlSetElevationPadText() {
	var elevationText = document.getElementById("elevationPad_elevationText");
	elevationText.textContent = "elevation: " + (gkControlContext.elevationFern * 10) + gkControlContext.elevationDeciFern;

	var elevationRect = document.getElementById("elevationPad_elevationRect");
//	var transX = (gkViewContext.scale - gkViewContext.lowScale) * 133;

	var transX = (gkViewContext.scale - gkViewContext.lowScale) / (gkViewContext.highScale - gkViewContext.lowScale) * 195;
	elevationRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	//console.log("transX: " + transX);
}

function gkControlSetWidthHeightPadText() {
	var widthHeightText = document.getElementById("widthHeightPad_widthHeightText");
	widthHeightText.textContent = gkViewContext.svgWidth + " X " + gkViewContext.svgHeight;

	var widthHeightRect = document.getElementById("widthHeightPad_widthHeightRect");
	var transX = (gkViewContext.svgWidth - 300) / 11.3;
	var transY = (gkViewContext.svgHeight - 300) / 11.3;
	widthHeightRect.setAttribute("transform","translate(" + transX + "," + transY +")");
	//console.log("transX: " + transX + " transY: " + transY);
}

function gkControlSetPanPadText() {
	var panText = document.getElementById("panPad_panText");

	panText.textContent = gkViewContext.viewOffsetIsoXYZ.x + "," + gkViewContext.viewOffsetIsoXYZ.y;

	var panRect = document.getElementById("panPad_panRect");

	isoWinXY = gkViewContext.viewOffsetIsoXYZ.convertToWin();

	var transX = ((isoWinXY.x - gkViewContext.lowPanX) / (gkViewContext.highPanX - gkViewContext.lowPanX)) * 200;
	var transY = ((isoWinXY.y - gkViewContext.lowPanY) / (gkViewContext.highPanY - gkViewContext.lowPanY)) * 200;

	panRect.setAttribute("transform","translate(" + transX + "," + transY +")");

	console.log("transX: " + transX + " transY: " + transY);

//	widthHeightText.textContent = gkViewContext.svgWidth + " X " + gkViewContext.svgHeight;
//
//	var widthHeightRect = document.getElementById("widthHeightPad_widthHeightRect");
//	var transX = (gkViewContext.svgWidth - 300) / 11.3;
//	var transY = (gkViewContext.svgHeight - 300) / 11.3;
//	widthHeightRect.setAttribute("transform","translate(" + transX + "," + transY +")");
//	//console.log("transX: " + transX + " transY: " + transY);

}

function gkControlSetBackgroundVolumePadText() {
	var backgroundVolumeText = document.getElementById("backgroundVolumePad_backgroundVolumeText");
	if (backgroundVolumeText != null) {
		var volumeLevel = gkAudioGetVolume(gkAudioContext.backgroundVolumeSelect);
		backgroundVolumeText.textContent = "vol: " + volumeLevel;

		var backgroundVolumeRect = document.getElementById("backgroundVolumePad_backgroundVolumeRect");
		var transX = volumeLevel * 190;
		backgroundVolumeRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	}
}

function gkControlSetEffectsVolumePadText() {
	var effectsVolumeText = document.getElementById("effectsVolumePad_effectsVolumeText");
	if (effectsVolumeText != null) {
		var volumeLevel = gkAudioGetVolume(gkAudioContext.effectsVolumeSelect);
		effectsVolumeText.textContent = "vol: " + volumeLevel;

		var effectsVolumeRect = document.getElementById("effectsVolumePad_effectsVolumeRect");
		var transX = volumeLevel * 190;
		effectsVolumeRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	}
}

function gkControlSetAttributePadText() {
	var terrainAttributeText = document.getElementById("terrainAttributePad_terrainAttributeText");
	if (terrainAttributeText != null) {
		var attributeIndex = gkTerrainEditGetAttributeIndex();
		terrainAttributeText.textContent = "attribute: " + gkTerrainEditGetAttributeText();

		var terrainAttributeRect = document.getElementById("terrainAttributePad_terrainAttributeRect");
		var transX = attributeIndex * 190;
		terrainAttributeRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	}
}

function gkControlClearCurrentMenu() {
	var layer = document.getElementById(gkControlContext.controlLayer);
	while (layer.firstChild) {
		layer.removeChild(layer.firstChild);
	}
}

function gkControlHandleUserPrefRestoreReq(jsonData) {
	var userPrefList = jsonData.userPrefList;
	var i;

	for (i = 0;i < userPrefList.length;i++) {
		var prefName = userPrefList[i].prefName;
		var prefValue = userPrefList[i].prefValue;

		if (prefName == "screenWidth") {
			gkViewContext.svgWidth = parseInt(prefValue);
		}
		if (prefName == "screenHeight") {
			gkViewContext.svgHeight = parseInt(prefValue);
		}
		if (prefName == "screenZoom") {
			gkViewContext.scale = parseFloat(prefValue);
		}
		if (prefName == "effectsVolume") {
			gkAudioVolumeChange(gkAudioContext.effectsVolumeSelect, prefValue);
		}
		if (prefName == "backgroundVolume") {
			gkAudioVolumeChange(gkAudioContext.backgroundVolumeSelect, prefValue);
		}
	}

	gkViewRender();
}

function gkControlHandleTerrainEditSelect() {
	console.log("got gkControlHandleTerrainEditSelect");

	gkTerrainEditSetAddTileOff();
	gkTerrainEditSetRemoveTileOff();
}

function gkControlHandleCloseTerrainEdit() {
	console.log("got gkControlHandleCloseTerrainEdit need to clear add/remove terrain flags");

	gkTerrainEditSetAddTileOff();
	gkTerrainEditSetRemoveTileOff();
}

function gkControlHandleAddTileSelect() {
	console.log("need to set check on add tile");
	gkTerrainEditSetAddTileOn();
	gkTerrainEditSetRemoveTileOff();

	gkControlChangeTileCheckMark();
}

function gkControlHandleRemoveTileSelect() {
	console.log("need to set check on remove tile");
	gkTerrainEditSetAddTileOff();
	gkTerrainEditSetRemoveTileOn();

	gkControlChangeTileCheckMark();
}

function gkControlChangeTileCheckMark() {
	var checkMarkG;
	var checkMarkScale;

	checkMarkScale = "0.01";
	checkMarkG = document.getElementById("addTile_checkMark");
	if (gkTerrainEditIsAddTileOn()) {
		checkMarkScale = "1";
	}
	checkMarkG.setAttribute("transform","scale(" + checkMarkScale + ")");

	checkMarkScale = "0.01";
	checkMarkG = document.getElementById("removeTile_checkMark");
	if (gkTerrainEditIsRemoveTileOn()) {
		checkMarkScale = "1";
	}
	checkMarkG.setAttribute("transform","scale(" + checkMarkScale + ")");
}

function gkControlIsMenuUp() {
//	console.log("testing for menu up stack length: " + gkControlContext.menuStack.length);

	return gkControlContext.menuStack.length > 1
}

function gkControlGetMenuWidth() {
	return gkControlContext.menuWidth;
}

function gkControlGetMenuHeight() {

	var nextLevelControlId = gkControlContext.menuStack[gkControlContext.menuStack.length - 1];


	var menuItemCount = gkControlContext.menuMap[nextLevelControlId].length;

//	console.log("GetMenuHeight mic: " + menuItemCount + " nlci: " + nextLevelControlId);

	return menuItemCount * gkControlContext.menuItemHeight;
}

function gkControlHandleTerrainSaveEdit() {
	console.log("gkControlHandleTerrainSaveEdit");

	var totalTerrain = gkTerrainGetCurrentTerrain();

	jsonMessage = JSON.stringify(totalTerrain);

	gkWsSendMessage("saveTerrainEditReq~" + jsonMessage + "~");

}

