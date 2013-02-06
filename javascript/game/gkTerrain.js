
//This is still experimental
//Begin map data
var MapData=new Array("PapayaWhip","IndianRed","LightSalmon","Wheat","Salmon","PaleGoldenRod","LightSalmon","Moccasin","NavajoWhite","SaddleBrown","Peru","Tan","Wheat","Moccasin","IndianRed","SandyBrown","PeachPuff","Bisque","Brown","BlanchedAlmond","Chocolate","Coral","DarkSalmon","NEXT_SECTION","AliceBlue","Aquamarine","Aqua","Blue","CornflowerBlue","CadetBlue","Cyan","DarkSlateBLue","DarkSeaGreen","LightSeaGreen","MediumSeaGreen","MediumSpringGreen","SeaGreen","Teal","NEXT_SECTION","Salmon","Red","Orange","OrangeRed","Tomato","Yellow","DimGrey","NEXT_SECTION","Yellow","YellowGreen","SpringGreen","MediumSeaGreen","MediumSpringGreen","LimeGreen","LightGreen","LawnGreen","Green","GreenYellow","ForestGreen","DarkSeaGreen","DarkGreen","Chartreuse","OliveDrab","NEXT_SECTION","DarkGoldenRod","DarkGray","DarkKhaki","DarkOliveGreen","Olive","OliveDrab","Peru","SaddleBrown","Sienna");
var Rendered=new Array();

function gkRenderMap (mapId,size) {
	var pos;
	var k = 1;
	var map = MapIndex[mapId];
	var l = map.length
	if (mapId=0) {
		k = Math.floor((Math.random()*22)); 
	}
	else if (mapId=1) {
		k = Math.floor((Math.random()*14)+23); 
	}
	else if (mapId=2) {
		k = Math.floor((Math.random()*7)+37); 
	}
	else if (mapId=3) {
		k = Math.floor((Math.random()*15)+44); 
	}
	else if (mapId=4) {
		k = Math.floor((Math.random()*9)+59); 
	}
	for (var i=1; i<=size; i++) {
		for (var j=1; j<=size; j++) {
			pos.isoXYZ.z = 0
			pos.isoXYZ.x = i
			pos.isoXYZ.y = j
			gkIsoCreateSingleDiamond(pos.isoXYZ, MapData[k]);
			Rendered[k] = MapData[k];
			k++;
		}
	}
}

function gkRestorePixel (xv,yv,size) {
	var pos;
	pos.isoXYZ.x = xv;
	pos.isoXYZ.y = yv;
	pos.isoXYZ.z = 0;
	gkIsoCreateSingleDiamond(pos.isoXYZ, map[k]);
}

function gkGetBgPixelColor (xv,yv,size) {
	var i = xv * size + yv;
	var out = MapData[i];
	return out;
}

function gkLoadSpecialMap (location,size) {
	for (var j=1; j<=size; j++) {
		pos.isoXYZ.z = 0
		pos.isoXYZ.x = i
		pos.isoXYZ.y = j
		gkIsoCreateSingleDiamond(pos.isoXYZ, MapData[k]);
		Rendered[k] = MapData[k];
		k++;
	}
}