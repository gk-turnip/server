
// handle the game playing field
// objectMap is a list of objects on the field (dandelions, avatars etc.)
// avatarId is the id of the current users avatar
var gkFieldContext = new gkFieldContextDef();

function gkFieldContextDef() {
	this.objectMap = new Object();
//	this.avatarId = null;
	this.avatarDestination = null;
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

	var g = gkIsoCreateSvgDiamond(rawSvgData);

	var isoXYZ = new GkIsoXYZDef(parseInt(jsonData.x), parseInt(jsonData.y), parseInt(jsonData.z))
	var originX = parseInt(jsonData.origin_x)
	var originY = parseInt(jsonData.origin_y)
	gkIsoSetSvgPositionWithOffset(g, isoXYZ, originX, originY);

	g.setAttribute("id",jsonData.id)

	if ((jsonData.userName != undefined) && (jsonData.userName.length > 0)) {
		var text = document.createElementNS(GK_SVG_NAMESPACE, "text");
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
		var g = document.getElementById(fieldObject.id);
		g.parentNode.removeChild(g);
		delete gkFieldContext.objectMap[jsonData.id];
		delete gkFieldContext.avatarId;
	}
}

// add a new avatar to the field for the current user
function gkFieldAddAvatar(jsonData, data) {
	gkFieldContext.avatarId = jsonData.id
	gkFieldAddSvg(jsonData, data);
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
				gkIsoSetSvgPositionWithOffset(fieldObject.g, fieldObject.isoXYZCurrent, fieldObject.originX, fieldObject.originY);
			}
		}
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
			g.parentNode.removeChild(g);
			delete gkFieldContext.objectMap[jsonData.id];
		}
	}
	delete gkFieldContext.avatarId;
}

