
var gkAvatarContext = new gkAvatarContextDef();

function gkAvatarContextDef() {
	this.currentAvatar = null;
	this.isoXYZCurrent = new GkIsoXYZDef(20, 20, 0);
	this.isoXYZDestination = new GkIsoXYZDef(20, 20, 0);
	this.avatarMeta = null;
}

function gkAvatarInit() {
	setInterval(gkAvatarAnimateLoop,50);
}

function gkAvatarAnimateLoop() {
	var field;
	var avatar;
	var needMove = false;

	field = document.getElementById("gkField");
	avatar = field.getElementById("avatar");

	if (gkAvatarContext.isoXYZDestination.x > gkAvatarContext.isoXYZCurrent.x) {
		needMove = true;
		gkAvatarContext.isoXYZCurrent.x += 1;
	}
	if (gkAvatarContext.isoXYZDestination.y > (gkAvatarContext.isoXYZCurrent.y + 0)) {
		needMove = true;
		gkAvatarContext.isoXYZCurrent.y += 1;
	}

	if (gkAvatarContext.isoXYZDestination.x < gkAvatarContext.isoXYZCurrent.x) {
		needMove = true;
		gkAvatarContext.isoXYZCurrent.x -= 1;
	}
	if (gkAvatarContext.isoXYZDestination.y < (gkAvatarContext.isoXYZCurrent.y + 0)) {
		needMove = true;
		gkAvatarContext.isoXYZCurrent.y -= 1;
	}

	if (needMove) {
		gkAvatarSetCurrentTransform();
	}
}

function gkAvatarRemoveCurrentAvatar() {
	var field;
	var avatar;

	field = document.getElementById("gkField");
// id should be soft
	avatar = field.getElementById("avatar");
	if (avatar != null) {
		avatar.parentNode.removeChild(avatar);
	}
}

function gkAvatarAddNewAvatar(avatarMeta, avatarSvg) {
	var r1 = new DOMParser().parseFromString(avatarSvg, 'text/xml');
	var field;
	field = document.getElementById("gkField");
	field.appendChild(document.importNode(r1.documentElement.firstChild,true));
	gkAvatarContext.avatarMeta = avatarMeta
	gkAvatarSetNewDestination(gkAvatarContext.isoXYZCurrent);
	gkAvatarSetCurrentTransform();
}

function gkAvatarSetNewDestination(isoXYZ) {
	field = document.getElementById("gkField");
	avatar = field.getElementById("avatar");

	gkAvatarContext.isoXYZDestination = isoXYZ;

	var diamond;
	diamond = gkIsoCreateSingleDiamond(gkAvatarContext.isoXYZDestination, "#8f8fff");
	field.appendChild(diamond);
}

function gkAvatarSetCurrentTransform() {
	var avatarWinXY;

	avatarWinXY = gkAvatarContext.isoXYZCurrent.convertToWin()
	avatar.setAttribute("transform","translate(" + (avatarWinXY.x - gkAvatarContext.avatarMeta.origin_x) + "," + (avatarWinXY.y - gkAvatarContext.avatarMeta.origin_y) + ")");
}

function gkAvatarLoadNewAvatar(avatarName) {
	gkWsSendMessage("getSvgReq~{\"SvgName\":\"" + avatarName + "\"}~");
}

