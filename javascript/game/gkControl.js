
// this program handles the in-game controls

var gkControlContext = new gkControlContextDef();

function gkControlContextDef() {
	this.controlLayer = "gkControlLayer"
	this.controlUrlPrefix = "/assets/gk/controls/"
	this.loadMap = new Object();
	this.menuMap = null;
	this.menuItemHeight = 50;
	this.menuStack = new Array();
	this.mouseDown = false;
}

function gkControlInit() {
	//gkControlLoad("start", gkControlHandleLoadStart);
	gkControlLoadMenuMap();
}

function gkControlLoadMenuMap() {
	var xmlhttp = new XMLHttpRequest();
	var fullUrl = gkControlContext.controlUrlPrefix + "menuMap.json";

	xmlhttp.onreadystatechange=function() {
		if (xmlhttp.readyState == 4) {
			if (xmlhttp.status == 200) {
				gkControlHandleLoadMenuMap(xmlhttp.responseText);
			} else {
				console.log("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
			}
		}
	}

	xmlhttp.open("GET", fullUrl, true);
	xmlhttp.send();
}

function gkControlHandleLoadMenuMap(menuMapText) {
	gkControlContext.menuMap = JSON.parse(menuMapText, null);

	gkControlContext.menuStack.push("menu");
	for (var i = 0;i < gkControlContext.menuMap.menu.length;i++) {
		gkControlLoad(gkControlContext.menuMap.menu[i].display, i, gkControlHandleLoadMenuItem);
	}
}

function gkControlHandleLoadMenuItem(controlId, index) {
	console.log("got menu item control loaded controlId: " + controlId + " index: " + index);
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
				console.log("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
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
				console.log("error loading XMLHttpRequest " + fullUrl + " " + xmlhttp.status);
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

	if ((controlId == "widthHeightPad") || (controlId == "zoomPad")) {
		g.onmousedown = function() {
			gkControlMenuItemMouseDown(controlId, index);
		};
		g.onmouseup = function() {
			gkControlMenuItemMouseUp(controlId, index);
		};
		g.onmousemove = function(evt) {
			gkControlMenuItemMouseMove(evt, controlId, index);
		};
	}
	var layer = document.getElementById(gkControlContext.controlLayer);
	layer.appendChild(g);
	console.log("added to layer: " + gkControlContext.controlLayer);

	if (controlId == "zoomPad") {
		gkControlSetZoomPadText();
	}
}

function gkControlMenuItemClick(controlId, index) {
	console.log("gkControlMenuItemClick " + controlId + " " + index);

	var nextLevelControlId = controlId;

	if (controlId == "close") {
		gkControlContext.mouseDown = false;

		gkControlContext.menuStack.pop();
		nextLevelControlId = gkControlContext.menuStack[gkControlContext.menuStack.length - 1];
	} else {
		if (gkControlContext.menuMap[nextLevelControlId] != undefined) {
			gkControlContext.menuStack.push(controlId);
		}
	}

	console.log("new menu controlId: " + nextLevelControlId);
	if (gkControlContext.menuMap[nextLevelControlId] != undefined) {
		gkControlClearCurrentMenu();
		for (var i = 0;i < gkControlContext.menuMap[nextLevelControlId].length;i++) {
			gkControlLoad(gkControlContext.menuMap[nextLevelControlId][i].display, i, gkControlHandleLoadMenuItem);
		}
	}
}

function gkControlMenuItemMouseDown(controlId, index) {
	console.log("gkControlMenuItemMouseDown " + controlId + " " + index);
	gkControlContext.mouseDown = true;
}

function gkControlMenuItemMouseUp(controlId, index) {
	console.log("gkControlMenuItemMouseUp " + controlId + " " + index);
	gkControlContext.mouseDown = false;
}

function gkControlMenuItemMouseMove(evt, controlId, index) {
	evt.preventDefault();

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
			gkControlHandleWidthHeightPad(x,y);
		}
		if (controlId == "zoomPad") {
			gkControlHandleZoomPad(x);
		}
	}
}

function gkControlHandleWidthHeightPad(x, y) {
	var width, height

	width = Math.floor(300 + (x * 2700));
	height = Math.floor(300 + (y * 2700));

	console.log("new width: " + width + " height: " + height);
}

function gkControlHandleZoomPad(x) {
	var zoomLevel;

	zoomLevel = 0.5 + (x * 1.5);
	zoomLevel = Math.floor(zoomLevel * 10) / 10;

	gkViewContext.scale = zoomLevel;

	gkControlSetZoomPadText();

	gkViewRender();
	console.log("new zoom level: " + zoomLevel);
}

function gkControlSetZoomPadText() {
	var zoomText = document.getElementById("zoomPad_zoomText");
	zoomText.textContent = "zoom: " + gkViewContext.scale;

	var zoomRect = document.getElementById("zoomPad_zoomRect");
	var transX = (gkViewContext.scale - 0.5) * 133;
	zoomRect.setAttribute("transform","translate(" + transX + "," + 0 +")");
	console.log("transX: " + transX);
}

function gkControlClearCurrentMenu() {
	var layer = document.getElementById(gkControlContext.controlLayer);
	while (layer.firstChild) {
		layer.removeChild(layer.firstChild);
	}
}

