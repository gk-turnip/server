
var gkWsContext = new gkWsContextDef();

function gkWsContextDef() {
	this.ws = null;
	this.dispatchFunction = null;
}

function gkWsInit(dispatchFunction) {
	gkWsContext.dispatchFunction = dispatchFunction
	if (gkWsContext.ws != null) {
		console.log("closing extra ws")
		gkWsContext.ws.close();
		gkWsContext.ws = null;
	}

	gkWsContext.ws = new WebSocket("ws://www.gourdianknot.com:14003/ws/svg");
	gkWsContext.ws.onopen = function() { gkWsDoOpen(); };
	gkWsContext.ws.onmessage = function(evt) { gkWsDoMessage(evt); };
	gkWsContext.ws.onclose = function() { gkWsDoOnClose(); };
}

function gkWsDoOpen() {
	console.log("gkWsDoOpen");
}

function gkWsDoMessage(e) {
	console.log("gkWsDoMessage");

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
		console.log("command from gameServer: " + command);
		if ((nlIndex1 + 1) == nlIndex2) {
			console.log("no jsonData");
			jsonData = null;
		} else {
			jsonRawData = e.data.substring(nlIndex1 + 1, nlIndex2);
			console.log("got jsonRawData len: " + jsonRawData.length);
			jsonData = JSON.parse(jsonRawData, null);
		}

		if ((nlIndex2 + 1) == e.data.length) {
			console.log("no ws data");
			data = null;
		} else {
			data = e.data.substring(nlIndex2 + 1);
			console.log("got ws data len: " + data.length);
		}

		gkWsContext.dispatchFunction(command, jsonData, data);

	} else {
		console.log("did not understand input from game server");
	}
}

function gkWsDoOnClose() {
	console.log("gkWsDoOnClose");
}

function gkWsSendMessage(message) {
	gkWsContext.ws.send(message);
}

