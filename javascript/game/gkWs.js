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

// handle the websocket communications

var gkWsContext = new gkWsContextDef();

// common "global" stuff put into a single context
function gkWsContextDef() {
	this.ws = null;
	this.dispatchFunction = null;
	this.websocketAddressPrefix = null;
	this.websocketPath = null;
	this.sessionId = null;
	this.reconnectTries = 0;
	this.pingId = Math.floor(Math.random() * 32767)
	this.userName = "unknown"
}

// the websocket needs to be "initialized" so it can open a connection to the server
function gkWsInit(dispatchFunction, websocketAddressPrefix, websocketPath, sessionId, recon) {
	console.log("gkWsInit");
	gkWsContext.dispatchFunction = dispatchFunction
	gkWsContext.websocketAddressPrefix = websocketAddressPrefix
	gkWsContext.websocketPath = websocketPath
	gkWsContext.sessionId = sessionId
	if (gkWsContext.ws != null) {
		console.log("closing extra ws")
		gkWsContext.ws.close();
		gkWsContext.ws = null;
	}

	if (!recon) {
		gkWsSetStatusTryInitialConnection();
	}

	gkWsContext.ws = new WebSocket(gkWsContext.websocketAddressPrefix + websocketPath + "?ses=" + sessionId);
	gkWsContext.ws.onopen = function() { gkWsDoOpen(); };
	gkWsContext.ws.onmessage = function(evt) { gkWsDoMessage(evt); };
	gkWsContext.ws.onclose = function() { gkWsDoOnClose(); };
	gkWsContext.ws.onerror = function() { gkWsDoOnError(); };
	//scriptUse = document.getElementById("scriptUse");
	//scriptUse.removeChild(loaded);
}

// this is called (by the browser) when the websocket completes a connection
function gkWsDoOpen() {
	console.log("gkWsDoOpen");
	gkWsSetStatusWaitingPing();
	gkWsContext.reconnectTries = 0;

	gkWsDoPing()
	setInterval("gkWsDoPing", 5000);
}

// send a ping
function gkWsDoPing() {
	gkWsSendMessage("pingReq~{ \"pingId\":\"" + gkWsContext.pingId + "\" }~");
}
// this is called (by the browser) when a new message is received from the server
// it is decoded and sent to gkDispatch
function gkWsDoMessage(e) {
//	console.log("gkWsDoMessage");
	var nlIndex1 = -1;
	var nlIndex2 = -1;

	for (i = 0;i < e.data.length;i++) {
		if (e.data[i] == '~') {
			nlIndex1 = i
			break;
		}
	}
	for (i = (nlIndex1 + 1);i < e.data.length;i++) {
		if (e.data[i] == '~') {
			nlIndex2 = i
			break;
		}
	}

	if ((nlIndex1 != -1) && (nlIndex2 != -1)) {
		command = e.data.substring(0,nlIndex1);
		//console.log("command from gameServer: " + command);
		if ((nlIndex1 + 1) == nlIndex2) {
			// no json data
			jsonData = null;
		} else {
			jsonRawData = e.data.substring(nlIndex1 + 1, nlIndex2);
			jsonData = JSON.parse(jsonRawData, null);
		}

		if ((nlIndex2 + 1) == e.data.length) {
			// no data
			data = null;
		} else {
			data = e.data.substring(nlIndex2 + 1);
		}

		gkWsContext.dispatchFunction(command, jsonData, data);

	} else {
		console.log("did not understand input from game server");
	}
}

// this is called (by the browser) if the web socket is closed
function gkWsDoOnClose() {
	console.log("gkWsDoOnClose");
	gkFieldDelAllObjects();
	var errorOut = document.createTextNode("The WebSocket connection was closed.");
	if (mode == "debug") {
		alert("The WebSocket connection was closed.");
	}
	gkWsSetStatusNotConnected();
console.log("reconnectTries: " + gkWsContext.reconnectTries);
	if (gkWsContext.reconnectTries < 3) {
		console.log("Attempting to reconnect in 5 seconds");
		gkWsContext.reconnectTries++;
		setTimeout(gkWsAttemptReconnect, 5000);
	}
	else {
		console.error("Max reconnects exceeded");
	}
}

// this is called (by the browser) when the websockets get an error
function gkWsDoOnError() {
	console.log("gkWsDoOnError");
}

// this is called when we attempt to reconnect to the game server
function gkWsAttemptReconnect() {
	console.log("gkWsAttemptReconnect");
	gkWsSetStatusReconnecting();
	gkWsContext.ws = null;
	gkWsInit(gkWsContext.dispatchFunction, gkWsContext.websocketAddressPrefix, gkWsContext.websocketPath, gkWsContext.sessionId, true)
}

// this is called when the client side wants to send a message to the server
function gkWsSendMessage(message) {
	console.log("gkWsSendMessage bufferedAmount: " + gkWsContext.ws.bufferedAmount);
	gkWsContext.ws.send(message);
}

function gkWsSetStatusTryInitialConnection() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket trying initial connection";
	connectionStatus.style.backgroundColor = "blue";
	connectionStatus.style.color = "white";
}

function gkWsSetStatusConnected() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket connected";
	connectionStatus.style.backgroundColor = "green";
	connectionStatus.style.color = "white";
}

function gkWsSetStatusNotConnected() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket not connected";
	connectionStatus.style.backgroundColor = "red";
	connectionStatus.style.color = "white";
}

function gkWsSetStatusReconnecting() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket reconnecting...";
	connectionStatus.style.backgroundColor = "DarkOrange";
	connectionStatus.style.color = "black";
}

function gkWsSetStatusWaitingPing() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket waiting ping";
	connectionStatus.style.backgroundColor = "yellow";
	connectionStatus.style.color = "black";
}

function gkWsSetStatusPingError() {
	var connectionStatus = document.getElementById("wsConnectionStatus");
	connectionStatus.innerHTML="web socket ping error";
	connectionStatus.style.backgroundColor = "red";
	connectionStatus.style.color = "black";
}

function gkWsPingRes(jsonData) {
	if (jsonData.pingId == gkWsContext.pingId) {
		gkWsSetStatusConnected();
	} else {
		gkWsSetStatusPingError();
	}

	gkWsContext.pingId += 1

	if (gkWsContext.pingId > 32767) {
		gkWsContext.pingId = 1;
	}
}

function gkWsUserNameReq(jsonData) {
	gkWsContext.userName = jsonData.userName;
	console.log("got userName: " + gkWsContext.userName);
}

function gkWsChatReq(jsonData) {
	var userName = jsonData.userName
	var message = jsonData.message

	gkChatMessageFromServer(userName, message);
}

