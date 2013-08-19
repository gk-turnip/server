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

// browser compability check

function gkBrowserCheckCompatibility() {
	var status = document.getElementById("browserCompatibleStatus");
	var ok = true

	status.innerHTML = "browser support:";
	status.style.backgroundColor = "#4fff4f";
	if ("WebSocket" in window) {
		status.innerHTML = status.innerHTML + " WebSocket ok";
	} else {
		status.innerHTML = status.innerHTML + " WebSocket not supported";
		ok = false;
	}

	if (document.implementation.hasFeature("http://www.w3.org/TR/SVG11/feature#Image", "1.1")) {
		status.innerHTML = status.innerHTML + ", SVG ok";
	} else {
		status.innerHTML = status.innerHTML + ", SVG not supported";
		ok = false;
	}

	if (window.XMLHttpRequest) {
		status.innerHTML = status.innerHTML + ", XMLHttpRequest ok";
	} else {
		status.innerHTML = status.innerHTML + ", XMLHttpRequest not supported";
		ok = false;
	}

	var auidoOk = false;
	var tempAudio = document.createElement('audio');
	if (tempAudio.canPlayType('audio/mpeg')) {
		status.innerHTML = status.innerHTML + ", can play mp3";
		audioOk = true;
	}
	if (tempAudio.canPlayType('audio/ogg')) {
		status.innerHTML = status.innerHTML + ", can play ogg";
		audioOk = true;
	}

	if (!ok) {
		status.innerHTML = status.innerHTML + " YOUR BROWSER IS NOT SUPPORTED";
		status.style.backgroundColor = "red";
	} else {
		if (!audioOk) {
			status.innerHTML = status.innerHTML + " no audio support ";
		}
	}
}


