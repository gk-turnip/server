
var gkDrops = new Array();
var gkRainContext = new gkRainContextDef();

function gkRainContextDef () {
	this.dropsRequired = 0;
	this.dropsWidth = 500;
	this.dropsHeight = 300;
	this.dropsStateCount = 0;
}

function gkRainStart() {
	setInterval(gkRainLoop,100);
}

function gkRainOn() {
	gkRainContext.dropsRequired = 30
}

function gkRainOff() {
	gkRainContext.dropsRequired = 0
}

function gkRainLoop() {
	var field;

//	gkRainContext.dropsStateCount += 1;
//
//	if ((gkRainContext.dropsStateCount > 100) && (gkRainContext.dropsStateCount < 130)){
//		gkRainContext.dropsRequired += 1;
//	}
//
//	if (gkRainContext.dropsStateCount > 500) {
//		gkRainContext.dropsRequired = 0;
//	}
//
//	if (gkRainContext.dropsStateCount > 1000) {
//		gkRainContext.dropsStateCount = 0;
//	}

//				console.log("drops.length: " + drops.length);

	field = document.getElementById("gkField");
	var undefinedIndex = -1;
	var dropsCounted = 0;
	for (i = 0;i < drops.length;i++) {
		if (drops[i] == undefined) {
//						console.log("undefined found at: " + i);
			undefinedIndex = i;
		} else {
			if (drops[i].y > gkRainContext.dropsHeight) {
				drops[i].svgGroup.parentNode.removeChild(drops[i].svgGroup);
				delete drops[i];
			} else {
				drops[i].fallOne();
				dropsCounted += 1;
			}
		}
	}
	if (dropsCounted < gkRainContext.dropsRequired) {
		if (undefinedIndex != -1) {
			drops[undefinedIndex] = new GkDropDef();
			field.appendChild(drops[undefinedIndex].svgGroup);
		} else {
			drops.push(new GkDropDef());
			field.appendChild(drops[drops.length - 1].svgGroup);
		}
	}
}

function GkDropDef() {
	this.x = Math.floor(Math.random() * gkRainContext.dropsWidth) + 4;
	this.y = 0 - Math.floor(Math.random() * 200);
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
		this.svgGroup.setAttribute("transform","translate(" + this.x + "," + this.y + ") scale(" + this.scale + ")");
	}

	GkDropDef.prototype.fallOne = function() {
		this.y += (this.speed / 3) + 6;
		this.setTranslate();
	}

	this.setTranslate();
}

