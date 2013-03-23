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

// handle the game playing field (other than terrain)
// objectMap is a list of objects on the field (dandelions, avatars etc.)
// avatarId is the id of the current users avatar
var gkFieldContext = new gkFieldContextDef();

function gkFieldContextDef() {
	this.objectMap = new Object();
//	this.avatarId = null;
	this.avatarDestination = null;
	this.leftKeyDown = false;
	this.rightKeyDown = false;
	this.upKeyDown = false;
	this.downKeyDown = false;
}

function gkFieldInit() {
	setInterval(gkFieldMoveObjects,50)
}

// the attributes for a single object on the field
function gkFieldObjectDef(id, userName, g, isoXYZCurrent, isoXYZDestination, originX, originY) {
	this.id = id
	this.userName = userName
	this.g = g
	this.isoXYZCurrent = isoXYZCurrent
	this.isoXYZDestination = isoXYZDestination
	this.originX = originX
	this.originY = originY

	gkFieldObjectDef.prototype.setDestination = function (isoXYZ) {
		this.isoXYZDestination = isoXYZ
	}
}

// add an svg image into the field
function gkFieldAddSvg(jsonData, rawSvgData) {
//console.log("gkFieldAddSvg id: " + jsonData.id);

	var g = gkIsoCreateSvgObject(rawSvgData);

	var isoXYZ = new GkIsoXYZDef(parseInt(jsonData.x), parseInt(jsonData.y), parseInt(jsonData.z))
	var originX = parseInt(jsonData.originX)
	var originY = parseInt(jsonData.originY)
	gkIsoSetSvgObjectPositionWithOffset(g, isoXYZ, originX, originY);

	g.setAttribute("id",jsonData.id)
	if (jsonData.clickFunction != undefined) {
		g.setAttribute("onClick", jsonData.clickFunction);
	}
	if ((jsonData.userName != undefined) && (jsonData.userName.length > 0)) {
		var text = document.createElementNS(gkIsoContext.svgNameSpace, "text");
		text.setAttribute("stroke","#000000");
		text.setAttribute("stroke-width","0");
		text.setAttribute("x","0");
		text.setAttribute("y",originY);
		text.setAttribute("font-size","24");
		text.setAttribute("id",jsonData.id + "_userName");
		var userNameText = document.createTextNode(jsonData.userName);
		text.appendChild(userNameText)

		g.appendChild(text)
	}

	var layer;
	layer = document.getElementById(jsonData.layer);
	layer.appendChild(g)

	var destIsoXYZ = new GkIsoXYZDef(isoXYZ.x, isoXYZ.y, isoXYZ.z)
	var fieldObject = new gkFieldObjectDef(jsonData.id, jsonData.userName, g, isoXYZ, destIsoXYZ, originX, originY)
	gkFieldContext.objectMap[fieldObject.id] = fieldObject

	//console.log("got new field object userName: " + jsonData.userName + " id: " + jsonData.id);
}

// delete an svg object from the field
function gkFieldDelSvg(jsonData) {
//console.log("gkFieldDelSvg id: " + jsonData.id);
	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	if (fieldObject != undefined) {
		var g = document.getElementById(fieldObject.id);
		g.parentNode.removeChild(g);
		delete gkFieldContext.objectMap[jsonData.id];
	}
}

// move the svg object in the field
function gkFieldMoveSvg(jsonData) {
//console.log("gkFieldMoveSvg id: " + jsonData.id);
	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	if (fieldObject != undefined) {
		fieldObject.isoXYZDestination.x = parseInt(jsonData.x)
		fieldObject.isoXYZDestination.y = parseInt(jsonData.y)
		fieldObject.isoXYZDestination.z = parseInt(jsonData.z)
	}
}

// request a new avatar svg and jsonData from the server
function gkFieldLoadNewAvatar(avatarName) {
	gkFieldRemoveExistingAvatar()
	gkWsSendMessage("delAvatarSvgReq~~");
	if (avatarName != "") {
		gkWsSendMessage("getAvatarSvgReq~{\"SvgName\":\"" + avatarName + "\"}~");
	}
}

// if there is an avarar already, remove it
function gkFieldRemoveExistingAvatar() {
	if (gkFieldContext.avatarId != undefined) {
		var fieldObject = gkFieldContext.objectMap[gkFieldContext.avatarId];
		if (fieldObject == undefined) {
			console.error("ERROR undefined fieldObject trying to remove avatar");
		} else {
			var g = document.getElementById(fieldObject.id);
			if (g == undefined) {
				console.error("ERROR undefined g trying to remove avatar");
			} else {
				g.parentNode.removeChild(g);
				delete gkFieldContext.objectMap[fieldObject.id];
				delete gkFieldContext.avatarId;
			}
		}
	}
}

// add a new avatar to the field for the current user
function gkFieldAddAvatar(jsonData, data) {
	gkFieldContext.avatarId = jsonData.id
	gkFieldAddSvg(jsonData, data);
	var fieldObject = gkFieldContext.objectMap[gkFieldContext.avatarId]
	gkFieldUpdatePositionDisplay(fieldObject.isoXYZCurrent);
}

// set the current users avatar to a new position
function gkFieldSetNewAvatarDestination(isoXYZ) {
	var fieldObject
	if (gkFieldContext.avatarId != undefined) {
		fieldObject = gkFieldContext.objectMap[gkFieldContext.avatarId]
	}
	if (fieldObject != undefined) {
		fieldObject.setDestination(isoXYZ)
		gkWsSendMessage("moveAvatarSvgReq~{ \"id\":\"" + fieldObject.id + "\", \"x\":\"" + isoXYZ.x + "\", \"y\": \"" + isoXYZ.y + "\", \"z\": \"" + isoXYZ.z + "\" }~");
	}
}

// move all objects closer to their proper destination
function gkFieldMoveObjects() {
	gkFieldContext.objectMap

	for (var prop in gkFieldContext.objectMap) {
		var fieldObject;
		fieldObject = gkFieldContext.objectMap[prop];
		if (fieldObject.id != undefined) {
			var curIsoXYZ = fieldObject.isoXYZCurrent
			var destIsoXYZ = fieldObject.isoXYZDestination
			if ((curIsoXYZ.x != destIsoXYZ.x) ||
				(curIsoXYZ.y != destIsoXYZ.y)) {
				if (destIsoXYZ.x > curIsoXYZ.x) {
					fieldObject.isoXYZCurrent.x += 1;
				}
				if (destIsoXYZ.x < curIsoXYZ.x) {
					fieldObject.isoXYZCurrent.x -= 1;
				}
				if (destIsoXYZ.y > curIsoXYZ.y) {
					fieldObject.isoXYZCurrent.y += 1;
				}
				if (destIsoXYZ.y < curIsoXYZ.y) {
					fieldObject.isoXYZCurrent.y -= 1;
				}
				gkIsoSetSvgObjectPositionWithOffset(fieldObject.g, fieldObject.isoXYZCurrent, fieldObject.originX, fieldObject.originY);
				if (gkFieldContext.avatarId != undefined) {
					if (gkFieldContext.avatarId == fieldObject.id) {
						gkFieldUpdatePositionDisplay(fieldObject.isoXYZCurrent);
					}
				}
			}
		}
		if (gkFieldContext.avatarId != undefined) {
			if (gkFieldContext.avatarId == fieldObject.id) {

				var localIsoXYZ = new GkIsoXYZDef(
					fieldObject.isoXYZCurrent.x,
					fieldObject.isoXYZCurrent.y,
					fieldObject.isoXYZCurrent.z);

				localIsoXYZ.x -= gkViewContext.viewOffsetIsoXYZ.x;
				localIsoXYZ.y -= gkViewContext.viewOffsetIsoXYZ.y;
				localIsoXYZ.z -= gkViewContext.viewOffsetIsoXYZ.z;

				var winXY = localIsoXYZ.convertToWin();
				if (winXY.x < gkViewContext.scrollEdgeX) {
					gkViewContext.viewOffsetIsoXYZ.x -= 1;
					gkViewContext.viewOffsetIsoXYZ.y += 1;
					gkViewRender();
				}
				if ((winXY.x + gkViewContext.scrollEdgeX) > (gkViewContext.svgWidth / gkViewContext.scale)) {
					gkViewContext.viewOffsetIsoXYZ.x += 1;
					gkViewContext.viewOffsetIsoXYZ.y -= 1;
					gkViewRender();
				}
				if (winXY.y < gkViewContext.scrollEdgeY) {
					gkViewContext.viewOffsetIsoXYZ.x -= 1;
					gkViewContext.viewOffsetIsoXYZ.y -= 1;
					gkViewRender();
				}
				if ((winXY.y + gkViewContext.scrollEdgeY) > (gkViewContext.svgHeight / gkViewContext.scale)) {
					gkViewContext.viewOffsetIsoXYZ.x += 1;
					gkViewContext.viewOffsetIsoXYZ.y += 1;
					gkViewRender();
				}
			}
		}
	}

	var moveX, moveY
	moveX = 0;
	moveY = 0;
	if (gkFieldContext.leftKeyDown) {
		moveX -= 1;
		moveY += 1;
	}
	if (gkFieldContext.rightKeyDown) {
		moveX += 1;
		moveY -= 1;
	}
	if (gkFieldContext.upKeyDown) {
		moveX -= 1;
		moveY -= 1;
	}
	if (gkFieldContext.downKeyDown) {
		moveX += 1;
		moveY += 1;
	}
	if ((gkFieldContext.leftKeyDown) || (gkFieldContext.rightKeyDown) || (gkFieldContext.upKeyDown) || (gkFieldContext.downKeyDown)) {
		gkFieldSetArrowKeyDestination(moveX, moveY);
	}
}

// delete all objects from the field
// called if we lose communications from the server
function gkFieldDelAllObjects() {
	for (var prop in gkFieldContext.objectMap) {
		var fieldObject;
		fieldObject = gkFieldContext.objectMap[prop];
		if (fieldObject.id != undefined) {
			var g = document.getElementById(fieldObject.id);
			if (g == undefined) {
				console.error("ERROR did not find g in delete all id: " + fieldObject.id);
			} else {
				g.parentNode.removeChild(g);
			}
			delete gkFieldContext.objectMap[fieldObject.id];
		}
	}
	delete gkFieldContext.avatarId;
}

// fix this, an Iso function in the Field javascript file
function gkFieldUpdatePositionDisplay(isoXYZCurrent) {
	var v;

	v = document.getElementById("posValueX");
	v.innerHTML = isoXYZCurrent.x;
	v = document.getElementById("posValueY");
	v.innerHTML = isoXYZCurrent.y;
	v = document.getElementById("posValueZ");
	v.innerHTML = isoXYZCurrent.z;
}

function gkFieldSetArrowKeyDestination(x,y) {
	if (gkFieldContext.avatarId != undefined) {
		var fieldObject = gkFieldContext.objectMap[gkFieldContext.avatarId];
		if (fieldObject != undefined) {
			fieldObject.isoXYZDestination.x = fieldObject.isoXYZCurrent.x + x;
			fieldObject.isoXYZDestination.y = fieldObject.isoXYZCurrent.y + y;
		}
	}
}

function gkFieldSetRightKeyDown() {
	gkFieldContext.rightKeyDown = true;
}

function gkFieldSetRightKeyUp() {
	gkFieldContext.rightKeyDown = false;
}

function gkFieldSetLeftKeyDown() {
	gkFieldContext.leftKeyDown = true;
}

function gkFieldSetLeftKeyUp() {
	gkFieldContext.leftKeyDown = false;
}

function gkFieldSetUpKeyDown() {
	gkFieldContext.upKeyDown = true;
}

function gkFieldSetUpKeyUp() {
	gkFieldContext.upKeyDown = false;
}

function gkFieldSetDownKeyDown() {
	gkFieldContext.downKeyDown = true;
}

function gkFieldSetDownKeyUp() {
	gkFieldContext.downKeyDown = false;
}

