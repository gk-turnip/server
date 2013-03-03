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
	case "pingRes":
		gkWsPingRes(jsonData);
	default:
		console.log("did not understand command from game server " + command);
	}
}
