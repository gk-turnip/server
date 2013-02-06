
//This is still experimental
//Begin map data
var MapData=new Array("PapayaWhip","IndianRed","LightSalmon","Wheat","Salmon","PaleGoldenRod","LightSalmon","Moccasin","NavajoWhite","SaddleBrown","Peru","Tan","Wheat","Moccasin","IndianRed","SandyBrown","PeachPuff","Bisque","Brown","BlanchedAlmond","Chocolate","Coral","DarkSalmon","NEXT_SECTION","AliceBlue","Aquamarine","Aqua","Blue","CornflowerBlue","CadetBlue","Cyan","DarkSlateBLue","DarkSeaGreen","LightSeaGreen","MediumSeaGreen","MediumSpringGreen","SeaGreen","Teal","NEXT_SECTION","Salmon","Red","Orange","OrangeRed","Tomato","Yellow","DimGrey","NEXT_SECTION","Yellow","YellowGreen","SpringGreen","MediumSeaGreen","MediumSpringGreen","LimeGreen","LightGreen","LawnGreen","Green","GreenYellow","ForestGreen","DarkSeaGreen","DarkGreen","Chartreuse","OliveDrab","NEXT_SECTION","DarkGoldenRod","DarkGray","DarkKhaki","DarkOliveGreen","Olive","OliveDrab","Peru","SaddleBrown","Sienna");
var Rendered=new Array();
var TraverseX=new Array();
var TraverseY=new Array();

function gkRenderMap (mapId,size) {
	//MapIds: 0=desert, 1=ocean, 2=fire, 3=grassland, 4=bog
	var a;
	var isoXYZ;
	var k = 0;
	var map = MapIndex[mapId];
	var l = map.length
	if (mapId=0) {
		a = Math.floor((Math.random()*22)); 
	}
	else if (mapId=1) {
		a = Math.floor((Math.random()*14)+24); 
	}
	else if (mapId=2) {
		a = Math.floor((Math.random()*7)+39); 
	}
	else if (mapId=3) {
		a = Math.floor((Math.random()*15)+47); 
	}
	else if (mapId=4) {
		a = Math.floor((Math.random()*9)+63); 
	}
	for (var i=1; i<=size; i++) {
		for (var j=1; j<=size; j++) {
			isoXYZ = new GkIsoXYZDef(i, j, 0);
			diamond = gkIsoCreateSingleDiamond(isoXYZ, MapData[a]);
			field.appendChild(diamond);
			Rendered[k] = MapData[a];
			k++;
		}
	}
}

function gkRestorePixel (xv,yv,size) {
	isoXYZ = new GkIsoXYZDef(xv, yv, 0);
	diamond = gkIsoCreateSingleDiamond(isoXYZ, map[k]);
	field.appendChild(diamond);
}

//This has been delayed due to the me needing to talk to Turnip about some dine details here.
/*
function gkLoadSpecialMap (location,size) {
	var pos;
	var k;
	var mapRequest = makeHttpObject();
	mapRequest.open("GET", "assets/gk/javascript/game/gkAudio.js", false);
	mapRequest.send(null);
	print(request.responseText);
	for (var i=1; i<=size; i++) {
		for (var j=1; j<=size; j++) {
			isoXYZ = new GkIsoXYZDef(i, j, 0);
			diamond = gkIsoCreateSingleDiamond(isoXYZ, SpecialMap[k]);
			field.appendChild(diamond);
			Rendered[k] = SpecialMap[k];
			k++;
		}
	}
}
*/

function gkTraverseAll (size) {
	var a = 0;
	var b = 0;
	var c = 0;
	var x = 0;
	var isoXYZ = new GkIsoXYZDef(a, b, c);
	for (; a<=size; a++) {
		for (var b=1; b<=size; b++) {
			winX, winY = GkIsoXYZDef(a, b, c);
			TraverseX[x] = winX;
			TraverseY[x] = winY;
		}
	}
}

			