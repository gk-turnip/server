
var gkAudioContext = new gkAudioContextDef();

function gkAudioContextDef() {
	this.preferredSuffix = ".none";
	this.preferredType = "audio/none";
	this.canPlayMp3 = false;
	this.canPlayOgg = false;
	this.canPlayWav = false;
	this.sourceDir = "unknown";
}

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
}

function gkAudioStartAudio(audioSelect, sourceFile) {
	var source = document.createElement("source");
	var audio = document.getElementById("audio" + audioSelect);

	for (i = (audio.children.length - 1); i >= 0; i--) {
		audio.removeChild(audio.children[i]);
	}

	source.type = gkAudioContext.preferredType;
	source.src = gkAudioContext.sourceDir + "/assets/gk/audio/" + sourceFile + gkAudioContext.preferredSuffix;
	audio.appendChild(source);
	if (audioSelect == 1) {
		audio.addEventListener('ended', function() {
  			this.currentTime = 0;
			this.play();
			}, false);
	}
	audio.play();
}

function gkAudioVolumeChange(audioSelect, volumeValue) {
	var audio = document.getElementById("audio" + audioSelect);
	audio.volume = volumeValue;
}

