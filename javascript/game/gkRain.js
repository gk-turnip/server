
var rainVolumeOriginal;
var rainFadeAmount;
var rainFadeInterval;
var rainFadeTime;
var gkDrops = new Array();
var rainLast;
var gkRainContext = new gkRainContextDef();

function gkRainContextDef () {
	this.dropsRequired = 0;
	this.dropsWidth = 500;
	this.dropsHeight = 300;
	this.dropsStateCount = 0;
}

function gkRainStart() {
	setInterval(gkRainLoop,100);
	var rainFadeAmount = 0.1;
}

function gkRainOn() {
	gkRainContext.dropsRequired = 30
	var rainTBP = document.getElementById("audio3");
	rainTBP.play();
	rainFadeAmount = 0.1;
	rainFadeInterval = setInterval(gkRainVolumeFader,rainFadeTime);
}

function gkRainOff() {
	gkRainContext.dropsRequired = 0
	rainFadeAmount = -0.1;
	rainFadeInterval = setInterval(gkRainVolumeFader,rainFadeTime);
}

function gkRainVolumeFader() {
	var rainTBP = document.getElementById("audio3");
	var rainVolumeOriginal = rainTBP.volume;
	if ((rainTBP.volume + rainFadeAmount <= 1) && (rainTBP.volume + rainFadeAmount >= 0)) {
		rainTBP.volume += rainFadeAmount;
	}
	if ((Math.min(rainTBP.volume,rainVolumeOriginal) == rainVolumeOriginal) && rainTBP.volume != rainVolumeOriginal) {
		rainTBP.volume -= rainFadeAmount;
		clearInterval(rainFadeInterval);
	}
	if (rainTBP.volume == rainLast) {
		clearInterval(rainFadeInterval);
	}
	rainLast = rainTBP.volume;
}
function gkRainLoop() {
	var tileLayer;

	tileLayer = document.getElementById("gkField");
	var undefinedIndex = -1;
	var dropsCounted = 0;
	for (i = 0;i < gkDrops.length;i++) {
		if (gkDrops[i] == undefined) {
			undefinedIndex = i;
		} else {
			if (gkDrops[i].isoXYZ.z < 3) {
						gkDrops[i].fallOne();
				if (gkDrops[i].isoXYZ.z < -20) {
						gkDrops[i].fallOne();
					gkDrops[i].diamond.parentNode.removeChild(gkDrops[i].diamond);
					delete gkDrops[i];
				} else {
					if (gkDrops[i].diamond == undefined) {
						gkDrops[i].isoXYZ.z = 0;
						var diamond;
						diamond = gkIsoCreateSingleDiamond(gkDrops[i].isoXYZ, "#0000ff", 0.5);
						tileLayer.appendChild(diamond);
						gkDrops[i].diamond = diamond;
						gkDrops[i].svgGroup.parentNode.removeChild(gkDrops[i].svgGroup);
					}
					gkDrops[i].fallOne();
					dropsCounted += 1;
				}
			} else {
				gkDrops[i].fallOne();
				dropsCounted += 1;
			}
		}
	}
	if (dropsCounted < gkRainContext.dropsRequired) {
		if (undefinedIndex != -1) {
			gkDrops[undefinedIndex] = new GkDropDef();
			tileLayer.appendChild(gkDrops[undefinedIndex].svgGroup);
		} else {
			gkDrops.push(new GkDropDef());
			tileLayer.appendChild(gkDrops[gkDrops.length - 1].svgGroup);
		}
	}
}

function GkDropDef() {
	var x, y, z;

	x = Math.floor(Math.random() * GK_SVG_WIDTH)
	y = Math.floor(Math.random() * GK_SVG_HEIGHT)
	z = GK_SVG_HEIGHT / 10;

	var tempWinXY;
	tempWinXY = new GkWinXYDef(x,y);
	this.isoXYZ = tempWinXY.convertToIso(z);
	//this.diamond = null;

	this.speed = Math.floor(Math.random() * 4);
	this.scale = 0.1 + (Math.floor(Math.random() * 3) / 20);

	this.svgGroup = document.createElementNS("http://www.w3.org/2000/svg","g");
	this.path = document.createElementNS("http://www.w3.org/2000/svg","path");
	this.path.setAttributeNS(null,"d","m-17.42583,43.18433c6.46532,13.97175 58.40341,-23.50802 56.67932,-27.05641c-1.72409,-3.54838 -63.14463,13.08466 -56.67932,27.05641z");
	this.path.setAttributeNS(null,"stroke","#000000");
	this.path.setAttributeNS(null,"fill-opacity","0.4");
	this.path.setAttributeNS(null,"transform","rotate(-65.2283 10.6965 30.967)");
	this.path.setAttributeNS(null,"stroke-width","0");
	this.path.setAttributeNS(null,"fill","#1e49bf");
	this.svgGroup.appendChild(this.path);

	GkDropDef.prototype.setTranslate = function() {
		var winXY
		winXY = this.isoXYZ.convertToWin();

		this.svgGroup.setAttribute("transform","translate(" + winXY.x + "," + winXY.y + ") scale(" + this.scale + ")");
	}

	GkDropDef.prototype.fallOne = function() {
		this.isoXYZ.z -= 3;
		this.setTranslate();
	}

	this.setTranslate();
}

