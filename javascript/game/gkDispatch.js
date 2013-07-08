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

// dispatch the message from the websocket connection
// to the correct function
function gkDispatchWsMessage(command, jsonData, data) {
//console.log("got command: " + command);
	switch (command) {
	case "turnOnRainReq":
		//gkRainOn();
		break;
	case "turnOffRainReq":
		//gkRainOff();
		break;
	case "getAvatarSvgRes":
		console.log("gkDispatch getAvatarSvgRes");
		gkFieldAddAvatar(jsonData, data)
		break;
	case "addSvgReq":
		console.log("gkDispatch addSvgReq");
		gkFieldAddSvg(jsonData, data);
		break;
	case "delSvgReq":
		console.log("gkDispatch delSvgReq");
		gkFieldDelSvg(jsonData);
		break;
	case "moveSvgReq":
		console.log("gkDispatch moveSvgReq");
		gkFieldMoveSvg(jsonData);
		break;
	case "setSvgReq":
		console.log("gkDispatch setSvgReq");
		gkFieldSetSvg(jsonData);
		break;
	case "clearTerrainReq":
		console.log("gkDispatch clearTerrainReq");
		gkTerrainClearTerrain(jsonData);
		break;
	case "setTerrainSvgReq":
		console.log("gkDispatch setTerrainSvgReq");
		gkTerrainSetTerrainSvg(jsonData, data);
		break;
	case "setTerrainMapReq":
		console.log("gkDispatch setTerrainMapReq");
		gkTerrainSetTerrainMap(jsonData);
		break;
	case "pingRes":
		gkWsPingRes(jsonData);
		break;
	case "userNameReq":
		console.log("gkDispatch userNameReq");
		gkWsUserNameReq(jsonData);
		break;
	case "chatReq":
		gkWsChatReq(jsonData);
		break;
	case "sendPastChatReq":
		gkWsChatSendPastChatReq(jsonData);
		break;
	case "newPodTitleReq":
		console.log("gkDispatch newPodTitleReq");
		gkFieldNewPodTitleReq(jsonData);
		break;
	default:
		console.error("did not understand command from game server " + command);
	}
}
