function gkDispatchWsMessage(command, jsonData, data) {
console.log("got command: " + command);
	switch (command) {
	case "turnOnRainReq":
		gkRainOn();
		break;
	case "turnOffRainReq":
		gkRainOff();
		break;
	case "getAvatarSvgRes":
		gkFieldAddAvatar(jsonData, data)
		break;
	case "addSvgReq":
		gkFieldAddSvg(jsonData, data);
		break;
	case "delSvgReq":
		gkFieldDelSvg(jsonData);
		break;
	case "moveSvgReq":
		gkFieldMoveSvg(jsonData);
		break;
	case "loadTerrainReq":
		gkTerrainLoad(jsonData, data);
		break;
	case "setTerrainReq":
		gkTerrainSetDiamond(jsonData);
		break;
	default:
		console.log("did not understand command from game server " + command);
	}
}
