
var gkFieldContext = new gkFieldContextDef();

function gkFieldContextDef() {
	this.objectMap = new Object();
	this.avatarId = null;
	this.avatarDestination = null;
}

function gkFieldInit() {
	setInterval(gkFieldMoveObjects,50)
}

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

function gkFieldAddSvg(jsonData, rawSvgData) {
console.log("gkFieldAddSvg");

	var g = gkIsoCreateSvgDiamond(rawSvgData);

	var isoXYZ = new GkIsoXYZDef(parseInt(jsonData.x), parseInt(jsonData.y), parseInt(jsonData.z))
	var originX = parseInt(jsonData.origin_x)
	var originY = parseInt(jsonData.origin_y)
	gkIsoSetSvgPositionWithOffset(g, isoXYZ, originX, originY);

	g.setAttribute("id",jsonData.id)

	var text = document.createElementNS(GK_SVG_NAMESPACE, "text");
	text.setAttribute("stroke","#000000");
	text.setAttribute("stroke-width","0");
	text.setAttribute("x","0");
	text.setAttribute("y","0");
	text.setAttribute("font-size","24");
	text.setAttribute("id",jsonData.id + "_userName");
	var userNameText = document.createTextNode(jsonData.userName);
	text.appendChild(userNameText)

	g.appendChild(text)

	var field;
	field = document.getElementById("gkField");
	field.appendChild(g)

	var fieldObject = new gkFieldObjectDef(jsonData.id, jsonData.userName, g, isoXYZ, isoXYZ, originX, originY)
	gkFieldContext.objectMap[fieldObject.id] = fieldObject

	console.log("got new field object userName: " + jsonData.userName);
}

function gkFieldDelSvg(jsonData) {
console.log("gkFieldDelSvg");
	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	if (fieldObject != undefined) {
		var g = document.getElementById(fieldObject.id);
		g.parentNode.removeChild(g);
		delete gkFieldContext.objectMap[jsonData.id];
	}
}

function gkFieldMoveSvg(jsonData) {
console.log("gkFieldMoveSvg");
console.log(jsonData.id + " " + jsonData.x + " " + jsonData.y)
	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	if (fieldObject != undefined) {
		fieldObject.isoXYZCurrent.x = parseInt(jsonData.x)
		fieldObject.isoXYZCurrent.y = parseInt(jsonData.y)
		fieldObject.isoXYZCurrent.z = parseInt(jsonData.z)
		gkIsoSetSvgPositionWithOffset(fieldObject.g, fieldObject.isoXYZCurrent, fieldObject.originX, fieldObject.originY);
	}
}

function gkFieldLoadNewAvatar(avatarName) {
	gkFieldRemoveExistingAvatar()
	gkWsSendMessage("delAvatarSvgReq~~");
	if (avatarName != "") {
		gkWsSendMessage("getAvatarSvgReq~{\"SvgName\":\"" + avatarName + "\"}~");
	}
}

function gkFieldRemoveExistingAvatar() {
	if (gkFieldContext.avatarId != undefined) {
		var fieldObject = gkFieldContext.objectMap[gkFieldContext.avatarId];
		var g = document.getElementById(fieldObject.id);
		g.parentNode.removeChild(g);
		delete gkFieldContext.objectMap[jsonData.id];
		gkFieldContext.avatarId = null
	}
}

function gkFieldAddAvatar(jsonData, data) {
	gkFieldContext.avatarId = jsonData.id
	gkFieldAddSvg(jsonData, data);
}

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
	gkFieldContext.avatarId = null
}

