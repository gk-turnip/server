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
	this.objectMap = new Array();
	this.refObjectMap = new Array();
//	this.avatarId = null;
	this.avatarDestination = null;
	this.leftKeyDown = false;
	this.rightKeyDown = false;
	this.upKeyDown = false;
	this.downKeyDown = false;
	this.lastIntervalTime = 0;
	this.duration1 = 0;
	this.duration2 = 0;
	this.duration3 = 0;
	this.inFocus = true;
	this.frameRate = 0;
	this.maxElevationMove = 11;
	this.oldAvatarDestination = undefined;
}

function gkFieldInit() {
	setInterval(gkFieldMoveObjects,50)
}

// the attributes for a single object on the field <g> tag put into <defs>
function gkFieldObjectDef(id, g) {
	this.id = id
	this.g = g
//	this.isoXYZCurrent = isoXYZCurrent
//	this.isoXYZDestination = isoXYZDestination
//	this.originX = originX
//	this.originY = originY

//	gkFieldObjectDef.prototype.setDestination = function (isoXYZ) {
//		this.isoXYZDestination = isoXYZ
//	}
}

// the attributes for a single reference on the field <use> tag
function gkFieldRefObjectDef(id, userName, isoXYZCurrent, isoXYZDestination, originX, originY, originZ) {
	this.id = id;
	this.userName = userName;
	this.isoXYZCurrent = isoXYZCurrent;
	this.isoXYZDestination = isoXYZDestination;
	this.originX = originX;
	this.originY = originY;
	this.originZ = originZ;
	this.pushIsoXYZDestination = null;

	gkFieldRefObjectDef.prototype.setDestination = function (isoXYZ) {
		this.isoXYZDestination = isoXYZ;
	}

	gkFieldRefObjectDef.prototype.setPushDestination = function (isoXYZ) {
		this.pushIsoXYZDestination = isoXYZ;
	}
}
// add an svg image into the field
function gkFieldAddSvg(jsonData, rawSvgData) {
//console.log("gkFieldAddSvg id: " + jsonData.id);

	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	if (fieldObject == undefined) {
		var g = gkIsoCreateSvgObject(rawSvgData);
		g.setAttribute("id",jsonData.id);
		var gkDefs = document.getElementById("gkDefs");
		gkDefs.appendChild(g);
		var fieldObject = new gkFieldObjectDef(jsonData.id, g)
		gkFieldContext.objectMap[fieldObject.id] = fieldObject
	}

	var g2 = document.createElementNS(gkIsoContext.svgNameSpace,"g");

	var ref = document.createElementNS(gkIsoContext.svgNameSpace,"use");
	ref.setAttributeNS(gkIsoContext.xlinkNameSpace,"href","#" + jsonData.id);
	var isoXYZ = new GkIsoXYZDef(parseInt(jsonData.x), parseInt(jsonData.y), parseInt(jsonData.z))
	var originX = parseInt(jsonData.originX)
	var originY = parseInt(jsonData.originY)
	var originZ = parseInt(jsonData.originZ)
	gkIsoSetSvgObjectPositionWithOffset(g2, isoXYZ, originX, originY, originZ);
	g2.setAttribute("id","ref_" + jsonData.id)

	if ((jsonData.userName != undefined) && (jsonData.userName.length > 0)) {
		var text = document.createElementNS(gkIsoContext.svgNameSpace, "text");
		text.setAttribute("stroke","#000000");
		text.setAttribute("stroke-width","0");
		text.setAttribute("x","40");
		text.setAttribute("y",originY + 50);
		text.setAttribute("font-size","24");
		text.setAttribute("style","text-anchor: middle");
		text.setAttribute("id",jsonData.id + "_userName");
		var userNameText = document.createTextNode(jsonData.userName);
		text.appendChild(userNameText)

		g2.appendChild(text)
	}

	g2.appendChild(ref);
	var layer;
	layer = document.getElementById(jsonData.layer);
	layer.appendChild(g2)

	var destIsoXYZ = new GkIsoXYZDef(isoXYZ.x, isoXYZ.y, isoXYZ.z)
	var refObject = new gkFieldRefObjectDef(jsonData.id, jsonData.userName, isoXYZ, destIsoXYZ, originX, originY, originZ);
	gkFieldContext.refObjectMap[refObject.id] = refObject;
//	gkFieldContext.refObjectMap[refObject.id + 1] = "NEXT";
	//console.log("got new field object userName: " + jsonData.userName + " id: " + jsonData.id);
}

// delete an svg object from the field
function gkFieldDelSvg(jsonData) {
//console.log("gkFieldDelSvg id: " + jsonData.id);
	var refObject = gkFieldContext.refObjectMap[jsonData.id];
	if (refObject != undefined) {
		var ref = document.getElementById("ref_" + refObject.id);
		if (ref != undefined) {
			ref.parentNode.removeChild(ref);
		}
		delete gkFieldContext.refObjectMap[jsonData.id];
	}
}

// move the svg object in the field, animated move
function gkFieldMoveSvg(jsonData) {
//console.log("gkFieldMoveSvg id: " + jsonData.id);
	var refObject = gkFieldContext.refObjectMap[jsonData.id];
	if (refObject != undefined) {
		refObject.isoXYZDestination.x = parseInt(jsonData.x)
		refObject.isoXYZDestination.y = parseInt(jsonData.y)
		refObject.isoXYZDestination.z = parseInt(jsonData.z)
	}
}

// set the svg object position in the field directly
function gkFieldSetSvg(jsonData) {
console.log("gkFieldSetSvg id: " + jsonData.id);
	var refObject = gkFieldContext.refObjectMap[jsonData.id];
	if (refObject != undefined) {
		refObject.isoXYZDestination.x = parseInt(jsonData.x)
		refObject.isoXYZDestination.y = parseInt(jsonData.y)
		refObject.isoXYZDestination.z = parseInt(jsonData.z)
		refObject.isoXYZCurrent.x = parseInt(jsonData.x)
		refObject.isoXYZCurrent.y = parseInt(jsonData.y)
		refObject.isoXYZCurrent.z = parseInt(jsonData.z)

		var ref = document.getElementById("ref_" + refObject.id);

		gkIsoSetSvgObjectPositionWithOffset(ref, refObject.isoXYZCurrent, refObject.originX, refObject.originY, refObject.originZ);
	}
}

/*
function gkFieldEnumObjects(scanIsoXYZ, acceptedOffset) {
	console.log("gkFieldEnumObjects");
	var holder = document.getElementById("objectList");
	holder.innerHTML = "";
	var l = gkFieldContext.refObjectMap.push("last") - 1;
	gkFieldContext.refObjectMap.pop();
	for (var i=0;i<=l;i++) {
		var fieldObject = gkFieldContext.refObjectMap[i];
		if (fieldObject) {
			var isoXYZ = fieldObject.isoXYZCurrent;
			if ((Math.abs(isoXYZ.x - scanIsoXYZ.x) <= acceptedOffset) && (Math.abs(isoXYZ.y - scanIsoXYZ.y) <= acceptedOffset) && (Math.abs(isoXYZ.x - scanIsoXYZ.x) <= acceptedOffset)) {
				var listing = document.createTextNode("id: " + i + " x: " + isoXYZ.x + " y: " + isoXYZ.y + " z: " + isoXYZ.z + "\n");
				holder.appendChild(listing);
			}
		}
	}
}
*/

// request a new avatar svg and jsonData from the server
function gkFieldLoadNewAvatar(avatarName) {
	gkFieldContext.oldAvatarDestination = gkFieldRemoveExistingAvatar();
	gkWsSendMessage("delAvatarSvgReq~~");
	if (avatarName != "") {
		gkWsSendMessage("getAvatarSvgReq~{\"SvgName\":\"" + avatarName + "\"}~");
	}
}

// if there is an avarar already, remove it and return position
function gkFieldRemoveExistingAvatar() {
	if (gkFieldContext.avatarId != undefined) {
		var refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId];
		var refPos = refObject.isoXYZCurrent;
		if (refObject == undefined) {
			console.error("ERROR undefined fieldObject trying to remove avatar");
		} else {
			var ref = document.getElementById("ref_" + refObject.id);
			if (ref == undefined) {
				console.error("ERROR undefined g trying to remove avatar");
			} else {
				ref.parentNode.removeChild(ref);
				delete gkFieldContext.refObjectMap[refObject.id];
				delete gkFieldContext.avatarId;
			}
		}
		return refPos;
	}
}

function gkFieldRemoveOtherAvatars() {
	for (var prop in gkFieldContext.refObjectMap) {
		var refObject;
		refObject = gkFieldContext.refObjectMap[prop];

		if (gkFieldContext.avatarId != refObject.id) {
			var ref = document.getElementById("ref_" + refObject.id);
			if (ref == undefined) {
				console.error("ERROR undefined g trying to remove avatar");
			} else {
				ref.parentNode.removeChild(ref);
				delete gkFieldContext.refObjectMap[refObject.id];
			}
		}
	}
}

// add a new avatar to the field for the current user
function gkFieldAddAvatar(jsonData, data) {
console.log(jsonData);
	gkFieldContext.avatarId = jsonData.id
	gkFieldAddSvg(jsonData, data);
	var refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId]
	if (gkFieldContext.oldAvatarDestination != undefined) {
		gkFieldSetNewAvatarDestination(gkFieldContext.oldAvatarDestination);
		gkFieldContext.oldAvatarDestination = undefined;
	}
	gkFieldUpdatePositionDisplay(refObject.isoXYZCurrent);
}

function gkFieldHandleFieldClick(winX, winY) {
	console.log("doClick " + winX + "," + winY);

	var refObject
	if (gkFieldContext.avatarId != undefined) {
		refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId]
	}
	if (refObject != undefined) {

		var isoXYZ = gkViewConvertWinToIso(winX, gkViewContext.marginX, winY, gkViewContext.marginY, 0);

		var g = gkIsoCreateSingleDiamond(isoXYZ, "#00ff00", 0.5);

		isoXYZ = new GkIsoXYZDef(isoXYZ.x + refObject.isoXYZCurrent.z, isoXYZ.y + refObject.isoXYZCurrent.z, 0);

		//gkFieldSetNewAvatarDestination(winXY.convertToIso(0));
		gkFieldSetNewAvatarDestination(isoXYZ);

		gkAudioStartAudio(2, "boing", false);

		gkTerrainClearMoveMarker();
		gkTerrainSetMoveMarker(g);
	}
}

// set the current users avatar to a new position
function gkFieldSetNewAvatarDestination(isoXYZ) {
	var refObject
	if (gkFieldContext.avatarId != undefined) {
		refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId]
	}
	if (refObject != undefined) {
		refObject.setDestination(isoXYZ)
		gkWsSendMessage("moveAvatarSvgReq~{ \"id\":\"" + refObject.id + "\", \"x\":\"" + isoXYZ.x + "\", \"y\": \"" + isoXYZ.y + "\", \"z\": \"" + isoXYZ.z + "\" }~");
console.log("setting new destination: " + isoXYZ.x + "," + isoXYZ.y);
	}
}

// push the current users avatar to a new position, delayed for new pod
function gkFieldPushNewAvatarDestination(isoXYZ) {
	var refObject
	if (gkFieldContext.avatarId != undefined) {
		refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId]
	}
	if (refObject != undefined) {
		refObject.setPushDestination(isoXYZ)
console.log("pushing new destination: " + isoXYZ.x + "," + isoXYZ.y);
	}
}

// move all objects closer to their proper destination
function gkFieldMoveObjects() {
//	console.log("gkFieldMoveObjects");
	//gkFieldContext.objectMap

	var moveFlag = false;

	// check all of the svg objects	in the pod
	for (var prop in gkFieldContext.refObjectMap) {
		var refObject;
		refObject = gkFieldContext.refObjectMap[prop];
		if (refObject.id != undefined) {
			if (refObject.pushIsoXYZDestination != null) {
console.log("handle new push dest: " + refObject.pushIsoXYZDestination.x + "," + refObject.pushIsoXYZDestination.y);
				curIsoXYZ = refObject.pushIsoXYZDestination;
				refObject.isoXYZCurrent.x = refObject.pushIsoXYZDestination.x;
				refObject.isoXYZCurrent.y = refObject.pushIsoXYZDestination.y;
				refObject.isoXYZCurrent.z = refObject.pushIsoXYZDestination.z;
				destIsoXYZ = refObject.pushIsoXYZDestination;
				refObject.isoXYZDestination.x = refObject.pushIsoXYZDestination.x;
				refObject.isoXYZDestination.y = refObject.pushIsoXYZDestination.y;
				refObject.isoXYZDestination.z = refObject.pushIsoXYZDestination.z;
				var ref = document.getElementById("ref_" + refObject.id);
				gkIsoSetSvgObjectPositionWithOffset(ref, refObject.isoXYZCurrent, refObject.originX, refObject.originY, refObject.originZ);
				gkViewContext.viewOffsetIsoXYZ.x = refObject.isoXYZDestination.x - 40;

				gkWsSendMessage("setAvatarSvgReq~{ \"id\":\"" + refObject.id + "\", \"x\":\"" + refObject.isoXYZCurrent.x + "\", \"y\": \"" + refObject.isoXYZCurrent.y + "\", \"z\": \"" + refObject.isoXYZCurrent.z + "\" }~");

				refObject.pushIsoXYZDestination = null;
			}
			// now test to see if it needs to be moved
			var curIsoXYZ = refObject.isoXYZCurrent
			var destIsoXYZ = refObject.isoXYZDestination
			if ((curIsoXYZ.x != destIsoXYZ.x) ||
				(curIsoXYZ.y != destIsoXYZ.y)) {

				var newCurrentX = refObject.isoXYZCurrent.x;
				var newCurrentY = refObject.isoXYZCurrent.y;
	
				if (destIsoXYZ.x > curIsoXYZ.x) {
					newCurrentX += 1;
				}
				if (destIsoXYZ.x < curIsoXYZ.x) {
					newCurrentX -= 1;
				}
				if (destIsoXYZ.y > curIsoXYZ.y) {
					newCurrentY += 1;
				}
				if (destIsoXYZ.y < curIsoXYZ.y) {
					newCurrentY -= 1;
				}

				var testMove = gkTerrainTestMoveElevation(newCurrentX, newCurrentY, refObject.isoXYZCurrent.z, gkFieldContext.maxElevationMove)

				if (testMove.canMove) {
					refObject.isoXYZCurrent.x = newCurrentX;
					refObject.isoXYZCurrent.y = newCurrentY;
					refObject.isoXYZCurrent.z = testMove.z;
				} else {
					refObject.isoXYZDestination.x = refObject.isoXYZCurrent.x;
					refObject.isoXYZDestination.y = refObject.isoXYZCurrent.y;
					refObject.isoXYZDestination.z = refObject.isoXYZCurrent.z;
				}
/*
				var z = gkFieldGetElevation2(newCurrentX, newCurrentY);
				if (Math.abs(z - refObject.isoXYZCurrent.z) <= gkFieldContext.maxElevationMove) {
					refObject.isoXYZCurrent.x = newCurrentX;
					refObject.isoXYZCurrent.y = newCurrentY;
					refObject.isoXYZCurrent.z = z;
				} else {
					z = gkFieldGetElevation1(newCurrentX, newCurrentY);
					if (Math.abs(z - refObject.isoXYZCurrent.z) <= gkFieldContext.maxElevationMove) {
						refObject.isoXYZCurrent.x = newCurrentX;
						refObject.isoXYZCurrent.y = newCurrentY;
						refObject.isoXYZCurrent.z = z;
					} else {
						refObject.isoXYZDestination.x = refObject.isoXYZCurrent.x;
						refObject.isoXYZDestination.y = refObject.isoXYZCurrent.y;
						refObject.isoXYZDestination.z = refObject.isoXYZCurrent.z;
					}
				}
*/

				var ref = document.getElementById("ref_" + refObject.id);

				gkIsoSetSvgObjectPositionWithOffset(ref, refObject.isoXYZCurrent, refObject.originX, refObject.originY, refObject.originZ);
				if (gkFieldContext.avatarId != undefined) {
					if (gkFieldContext.avatarId == refObject.id) {
						gkFieldUpdatePositionDisplay(refObject.isoXYZCurrent);
					}
				}
				moveFlag = true;
			}
		}

		// test if the current users avatar has moved
		if (gkFieldContext.avatarId != undefined) {
			if (gkFieldContext.avatarId == refObject.id) {
//				gkFieldEnumObjects(refObject.isoXYZCurrent, 1);

				// the local users avatar is moving
				// so test if the view needs to be shifed
				// because the avatar is moving off the edge
				var localIsoXYZ = new GkIsoXYZDef(
					refObject.isoXYZCurrent.x,
					refObject.isoXYZCurrent.y,
					refObject.isoXYZCurrent.z);

				localIsoXYZ.x -= gkViewContext.viewOffsetIsoXYZ.x;
				localIsoXYZ.y -= gkViewContext.viewOffsetIsoXYZ.y;
				localIsoXYZ.z -= gkViewContext.viewOffsetIsoXYZ.z;

				var renderFlag = false;

				var winXY = localIsoXYZ.convertToWin();
				if (winXY.x < gkViewContext.scrollEdgeX) {
					gkViewContext.viewOffsetIsoXYZ.x -= 1;
					gkViewContext.viewOffsetIsoXYZ.y += 1;
					renderFlag = true;
				}
				if ((winXY.x + gkViewContext.scrollEdgeX) > (gkViewContext.svgWidth / gkViewContext.scale)) {
					gkViewContext.viewOffsetIsoXYZ.x += 1;
					gkViewContext.viewOffsetIsoXYZ.y -= 1;
					renderFlag = true;
				}
				if (winXY.y < gkViewContext.scrollEdgeY) {
					gkViewContext.viewOffsetIsoXYZ.x -= 1;
					gkViewContext.viewOffsetIsoXYZ.y -= 1;
					renderFlag = true;
				}
				if ((winXY.y + gkViewContext.scrollEdgeY) > (gkViewContext.svgHeight / gkViewContext.scale)) {
					gkViewContext.viewOffsetIsoXYZ.x += 1;
					gkViewContext.viewOffsetIsoXYZ.y += 1;
					renderFlag = true;
				}

				if (renderFlag) {
					moveFlag = true;
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
	if ((gkFieldContext.inFocus) && ((gkFieldContext.leftKeyDown) || (gkFieldContext.rightKeyDown) || (gkFieldContext.upKeyDown) || (gkFieldContext.downKeyDown))) {
		gkFieldSetArrowKeyDestination(moveX, moveY);
		moveFlag = true;
	}

	var endTime = (new Date()).getTime();

	if (moveFlag) {
		for (var evAC=0;evAC<gkTerrainContext.terrainAudioMap.length;evAC++) {
			if (gkFieldContext.avatarId == undefined) {
				console.error("no avatar id");
				break;
			}
			var refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId];
			if (refObject == undefined) {
				console.error("no avatar id");
				break;
			}
			if ((Math.abs(gkTerrainContext.terrainAudioMap[evAC].x - refObject.isoXYZCurrent.x) < 25) && (Math.abs(gkTerrainContext.terrainAudioMap[evAC].y - refObject.isoXYZCurrent.y) < 25) && (Math.abs(gkTerrainContext.terrainAudioMap[evAC].z - refObject.isoXYZCurrent.z) < 50)) {
				gkAudioStartAudio(4, gkTerrainContext.terrainAudioMap[evAC].clip, true)
				break;
			} else {
				eAudio = document.getElementById("audio4")
				eAudio.pause();
			}
		}
		if (gkFieldContext.lastIntervalTime > 0) {
			var duration = endTime - gkFieldContext.lastIntervalTime;
			gkFieldContext.duration3 = gkFieldContext.duration2;
			gkFieldContext.duration2 = gkFieldContext.duration1;
			gkFieldContext.duration1 = duration;
			gkFieldContext.frameRate = (3000.0 / (gkFieldContext.duration1 + gkFieldContext.duration2 + gkFieldContext.duration3)).toFixed(2);
			var frameRate = document.getElementById("frameRate");
			frameRate.innerHTML = "fps: " + gkFieldContext.frameRate;
		}
	}

	gkFieldContext.lastIntervalTime = endTime;
}

/*
function gkFieldGetElevation1(x, y) {
	return gkTerrainGetElevation1(x,y);
}

function gkFieldGetElevation2(x, y) {
	return gkTerrainGetElevation2(x,y);
}
*/

// delete all objects from the field
// called if we lose communications from the server
function gkFieldDelAllObjects() {
	for (var prop in gkFieldContext.refObjectMap) {
		var refObject;
		refObject = gkFieldContext.refObjectMap[prop];
		if (refObject.id != undefined) {
			var ref = document.getElementById("ref_" + refObject.id);
			if (ref == undefined) {
				console.error("ERROR did not find g in delete all id: " + fieldObject.id);
			} else {
				ref.parentNode.removeChild(ref);
			}
			delete gkFieldContext.refObjectMap[refObject.id];
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

function gkFieldNewPodTitleReq(jsonData) {
	var v;

	v = document.getElementById("podTitle");
	v.innerHTML = jsonData.podTitle;
}

function gkFieldSetArrowKeyDestination(x,y) {
	if (gkFieldContext.avatarId != undefined) {
		var refObject = gkFieldContext.refObjectMap[gkFieldContext.avatarId];
		if (refObject != undefined) {
			refObject.isoXYZDestination.x = refObject.isoXYZCurrent.x + x;
			refObject.isoXYZDestination.y = refObject.isoXYZCurrent.y + y;
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

function gkFieldGetFocus() {
	var inFocus = document.getElementById("inFocus");
	inFocus.innerHTML="game in focus";
	inFocus.style.backgroundColor = "green";
	inFocus.style.color = "white";
	gkFieldContext.inFocus = true;
}

function gkFieldLoseFocus() {
	var inFocus = document.getElementById("inFocus");
	inFocus.innerHTML="game not in focus";
	inFocus.style.backgroundColor = "red";
	inFocus.style.color = "white";
	gkFieldContext.inFocus = false;
}

