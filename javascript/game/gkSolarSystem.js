var Position = { "x":0 , "y":0 };
var SSD = new Array();


function gkSSCalculateLight(doClamp,clampL,clampH,doLimitTotal,totalMin,totalMax,doRender) {
	var dy;
	var dx;
	var m;
	var n;
	var totalLight;
	var lights = new Array();
	var out = new Array();
	for (var i=0;i<=SSD.length;i++) {
		dy = SSD[i].posy - Position.y;
		dx = SSD[i].posx - Position.x;
		n = dy/dx;
		m = n * Math.sqrt(dy * dy + dx * dx);
		if (SSD[i].posy > Position.y) {
			if (doClamp == true) {
				if (m < clampL) {
					m = clampL;
				}
				else if (m > clampH) {
					m = clampH;
				}
			}
		lights[i] = m;
		if (m > 0) {
			totalLight += m;
		}
	}
	if (doLimitTotal == true) {
		if (totalLight < totalMin) {
			lights = gkSSNormalizeLum(lights,totalLight,totalMin);
		}
		if (totalLight > totalMax) {
			lights = gkSSNormalizeLum(lights,totalLight,totalMax);
		}
	}
	for (var i=0;i<=lights.length;i++) {
		out[i].color = SSD[i].color;
		out[i].lum = lights[i];
	}
	if (doRender == true) {
//		Renderer goes here
	}
	else {
		return out;
	}
}

function gkSSNormalizeLum(in,total,target) {
	var x = target / total;
	for (var i=0;i<=in.length;i++) {
		in[i] *= x;
	}
	return in;
}