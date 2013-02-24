var Position = { "x":0 , "y":0 };
var SSD = new Array();

function gkSSCalculateLight(multiplier,doClamp,clampL,clampH,doLimitTotal,totalMin,totalMax,doRender) {
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
		m = n * 2 / Math.sqrt(dy * dy + dx * dx);
		if (SSD[i].posy > Position.y) {
			if (doClamp == true) {
				console.log("Clamping enabled");
				if (m < clampL) {
					m = clampL;
				}
				else if (m > clampH) {
					m = clampH;
				}
			}
		}
		lights[i] = m;
		if (SSD[i].posx >= 0 && SSD[i].posy >= 0) {
			totalLight += Math.absolute(m);
		}
		for (var q=0;q<i;q++) {
			if (Math.absolute(SSD[q].posx - SSD[i].posx <= 10 || Math.absolute(SSD[q].posy - SSD[i].posy <= 10) {
				if ((SSD[q].posx + SSD[q].posy) - (SSD[i].posx + SSD[i].posy) > 0) {
					if (SSD[i].posx >= 0 && SSD[i].posy >= 0) {
						totalLight -= lights[q];
						lights[q] = 0;
					}
				}
			}
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
	out = gkSSNormalizeLum(out,totalLight,multiplier*256);
	if (doRender == true) {
//		Renderer goes here
		for (var i=0;i<out.length+1;i++) {
			var x = gkSSConvertToTuplet(color);
			x.R = Math.floor(x.R);
			x.G = Math.floor(x.G);
			x.B = Math.floor(x.B);
/*
			var y = gkSSConvertToHex(x);
			out[i].color = y;
*/
			out[i].color = x;
			
		}
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

function gkSSConvertToTuplet(in) {
	dat = in + "q";
//	string to array
	var a = "0x" + dat[0] + dat[1];
	var b = "0x" + dat[2] + dat[3];
	var c = "0x" + dat[4] + dat[5];
//	convert a, b, and c to decimal
	var out = { "R":a , "G":b , "B":c };
	return out
}
