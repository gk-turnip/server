
var gkFieldContext = new gkFieldContextDef();

function gkFieldContextDef() {
	this.objectMap = new Object();
}

function gkFieldInit() {
}

function gkFieldObjectDef(id, x, y, z) {
	this.id = id
	this.x = x
	this.y = y
	this.z = z
}

function gkFieldAddSvg(jsonData, rawSvgData) {
console.log("gkFieldAddSvg");

	var g = gkIsoCreateSvgDiamond(rawSvgData);

	var isoXYZ = new GkIsoXYZDef(parseInt(jsonData.x), parseInt(jsonData.y), parseInt(jsonData.z))
	gkIsoSetSvgDiamondPosition(g, isoXYZ);

	g.setAttribute("id",jsonData.id)

	var field;
	field = document.getElementById("gkField");
	field.appendChild(g)

	var fieldObject = new gkFieldObjectDef(jsonData.id, jsonData.x, jsonData.y, jsonData.z)
	gkFieldContext.objectMap[fieldObject.id] = fieldObject
}

function gkFieldDelSvg(jsonData) {
console.log("gkFieldDelSvg");
	var fieldObject = gkFieldContext.objectMap[jsonData.id];
	var g = document.getElementById(fieldObject.id);
	g.parentNode.removeChild(g);
	delete gkFieldContext.objectMap[jsonData.id];
}

function gkFieldMoveSvg(jsonData) {
console.log("gkFieldMoveSvg");
}

function gkFieldAnimateLoop() {
	var field;
	var avatar;
	var needMove = false;

	field = document.getElementById("gkField");
	avatar = field.getElementById("avatar");

	if (gkFieldContext.isoXYZDestination.x > gkFieldContext.isoXYZCurrent.x) {
		needMove = true;
		gkFieldContext.isoXYZCurrent.x += 1;
	}
	if (gkFieldContext.isoXYZDestination.y > (gkFieldContext.isoXYZCurrent.y + 0)) {
		needMove = true;
		gkFieldContext.isoXYZCurrent.y += 1;
	}

	if (gkFieldContext.isoXYZDestination.x < gkFieldContext.isoXYZCurrent.x) {
		needMove = true;
		gkFieldContext.isoXYZCurrent.x -= 1;
	}
	if (gkFieldContext.isoXYZDestination.y < (gkFieldContext.isoXYZCurrent.y + 0)) {
		needMove = true;
		gkFieldContext.isoXYZCurrent.y -= 1;
	}

	if (needMove) {
		gkFieldSetCurrentTransform();
	}
}

function gkFieldRemoveCurrentAvatar() {
	var field;
	var avatar;

	field = document.getElementById("gkField");
// id should be soft
	avatar = field.getElementById("avatar");
	if (avatar != null) {
		avatar.parentNode.removeChild(avatar);
	}
}

function gkFieldAddNewAvatar(avatarMeta, avatarSvg) {
	var r1 = new DOMParser().parseFromString(avatarSvg, 'text/xml');
	var field;
	field = document.getElementById("gkField");
	field.appendChild(document.importNode(r1.documentElement.firstChild,true));
	gkFieldContext.avatarMeta = avatarMeta
	gkFieldSetNewDestination(gkFieldContext.isoXYZCurrent);
	gkFieldSetCurrentTransform();
}

function gkFieldSetNewDestination(isoXYZ) {
	field = document.getElementById("gkField");
	avatar = field.getElementById("avatar");

	gkFieldContext.isoXYZDestination = isoXYZ;

	if (gkFieldContext.lastMovePoint != undefined) {
		field.removeChild(gkFieldContext.lastMovePoint);
	}
	gkFieldContext.lastMovePoint = gkIsoCreateSingleDiamond(gkFieldContext.isoXYZDestination, "#8f8fff", 1.0);
	field.appendChild(gkFieldContext.lastMovePoint);
}

function gkFieldSetCurrentTransform() {
	var avatarWinXY;

	avatarWinXY = gkFieldContext.isoXYZCurrent.convertToWin()
	avatar.setAttribute("transform","translate(" + (avatarWinXY.x - gkFieldContext.avatarMeta.origin_x) + "," + (avatarWinXY.y - gkFieldContext.avatarMeta.origin_y) + ")");
}

function gkFieldLoadNewAvatar(avatarName) {
	gkWsSendMessage("getSvgReq~{\"SvgName\":\"" + avatarName + "\"}~");
}

