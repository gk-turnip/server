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

// audio control
// several audio streams or "audioSelect"
// which are 1, 2, or 3 (more in the future
// the first is automatically looping for background music
var gkAudioContext = new gkAudioContextDef();

function gkAudioContextDef() {
	this.preferredSuffix = ".none";
	this.preferredType = "audio/none";
	this.canPlayMp3 = false;
	this.canPlayOgg = false;
	this.canPlayWav = false;
	this.sourceDir = "unknown";
	this.backgroundVolumeSelect = 1;
	this.effectsVolumeSelect = 3;
}

// check what audio formats are supported by the browser
function gkAudioInit(sourceDir) {
	gkAudioContext.sourceDir = sourceDir;
	var audio1 = document.getElementById("audio1")
	if (audio1.canPlayType('audio/wav;')) {
		gkAudioContext.canPlayWav = true;
		gkAudioContext.preferredSuffix = ".wav";
		gkAudioContext.preferredType = "audio/wav";
	}
	if (audio1.canPlayType('audio/mpeg;')) {
		gkAudioContext.canPlayMp3 = true;
		gkAudioContext.preferredSuffix = ".mp3";
		gkAudioContext.preferredType = "audio/mpeg";
	}
	if (audio1.canPlayType('audio/ogg;')) {
		gkAudioContext.canPlayOgg = true;
		gkAudioContext.preferredSuffix = ".ogg";
		gkAudioContext.preferredType = "audio/ogg";
	}

	gkAudioVolumeChange("1",0.3);
	gkAudioVolumeChange("2",0.3);
	gkAudioVolumeChange("3",0,3);
	gkAudioVolumeChange("4",0,3);
}

// start a new audio source, loop if specified
function gkAudioStartAudio(audioSelect, sourceFile, loop) {
	var source = document.createElement("source");
	var audio = document.getElementById("audio" + audioSelect);

	for (i = (audio.children.length - 1); i >= 0; i--) {
		audio.removeChild(audio.children[i]);
	}

	source.type = gkAudioContext.preferredType;
	source.src = gkAudioContext.sourceDir + "/assets/gk/audio/" + sourceFile + gkAudioContext.preferredSuffix;
	audio.appendChild(source);
	if (loop) {
		audio.addEventListener('ended', function() {
  			this.currentTime = 0;
			this.play();
			}, false);
	}
	this.currentTime = 0;
	audio.play();
}

// change the volume for one audioSelect
function gkAudioVolumeChange(audioSelect, volumeValue) {
	var audio = document.getElementById("audio" + audioSelect);
	audio.volume = volumeValue;
}

function gkAudioGetVolume(audioSelect) {
	var audio = document.getElementById("audio" + audioSelect);
	return audio.volume;
}


