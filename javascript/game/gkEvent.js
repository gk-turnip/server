function gkEventInit(websocketAddressPrefx, websocketPath, audioPathPrefix, sessionId) {
	console.log("calling event init for seesion id: " + sessionId);

	gkChatInit();
	gkControlInit();
	
	var s1
	s1 = document.getElementById("selectAvatar");
	s1.selectedIndex = 0;
	//s1 = document.getElementById("volume1");
	//s1.selectedIndex = 1;
	//s1 = document.getElementById("volume2");
	//s1.selectedIndex = 1;

	gkWsInit(gkDispatchWsMessage, websocketAddressPrefx, websocketPath, sessionId);

	gkAudioInit(audioPathPrefix);

	gkAudioStartAudio(gkAudioContext.backgroundVolumeSelect, "Underwater2", true);

	gkTerrainInit(5);

	gkFieldInit();
}

function gkEventDoClick(e) {
	console.log("gkEventDoClick");

	if ((!gkControlIsMenuUp()) || (e.pageX > gkControlGetMenuWidth()) || (e.pageY > gkControlGetMenuHeight())) {
		if (gkTerrainEditNeedClick()) {
			gkTerrainEditHandleClick(e.pageX, e.pageY);
		} else {
			gkFieldHandleFieldClick(e.pageX, e.pageY);
		}
	}
}

function gkEventChangeAvatar() {
	avatarName = document.getElementById("selectAvatar").value;
	console.log("avatarName: " + avatarName)
	gkFieldLoadNewAvatar(avatarName);
}

//function volumeChange(volumeId) {
//	volumeSelect = document.getElementById("volume" + volumeId);
//	value = volumeSelect.options[volumeSelect.selectedIndex].value
//	gkAudioVolumeChange(volumeId, value);
//}

//function doSubmitControls() {
//	return false;
//}

//function setValue(id, value) {
//	var v = document.getElementById(id)
//	v.innerHTML = value
//}

function gkEventDoKeyDown(event) {
	if (event.keyCode == 37) {
		gkFieldSetLeftKeyDown();
	}
	if (event.keyCode == 39) {
		gkFieldSetRightKeyDown();
	}
	if (event.keyCode == 38) {
		gkFieldSetUpKeyDown();
	}
	if (event.keyCode == 40) {
		gkFieldSetDownKeyDown();
	}
}

function gkEventDoKeyUp(event) {
	if (event.keyCode == 37) {
		gkFieldSetLeftKeyUp();
	}
	if (event.keyCode == 39) {
		gkFieldSetRightKeyUp();
	}
	if (event.keyCode == 38) {
		gkFieldSetUpKeyUp();
	}
	if (event.keyCode == 40) {
		gkFieldSetDownKeyUp();
	}
}

function nothing() {
}
