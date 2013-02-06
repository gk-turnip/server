
//This is still experimental
//Begin map data
var MapIndex=new Array("MapDesert","MapOcean","MapFire","MapGrassland","MapBog")
var MapDesert=new Array("PapayaWhip","IndianRed","LightSalmon","Wheat","Salmon","PaleGoldenRod","LightSalmon","Moccasin","NavajoWhite","SaddleBrown","Peru","Tan","Wheat","Moccasin","IndianRed","SandyBrown","PeachPuff","Bisque","Brown","BlanchedAlmond","Chocolate","Coral","DarkSalmon");
var MapOcean=new Array("AliceBlue","Aquamarine","Aqua","Blue","CornflowerBlue","CadetBlue","Cyan","DarkSlateBLue","DarkSeaGreen","LightSeaGreen","MediumSeaGreen","MediumSpringGreen","SeaGreen","Teal");
var MapFire=new Array("Salmon","Red","Orange","OrangeRed","Tomato","Yellow","DimGrey");
var MapGrassland=new Array("Yellow","YellowGreen","SpringGreen","MediumSeaGreen","MediumStringGreen","LimeGreen","LightGreen","LawnGreen","Green","GreenYellow","ForestGreen","DarkSeaGreen","DarkGreen","Chartreuse","OliveDrab");
var MapBog=new Array("DarkGoldenRod","DarkGray","DarkKhaki","DarkOliveGreen","Olive","OliveDrab","Peru","SaddleBrown","Sienna");
var Rendered=new Array();

function gkRenderMap (mapId,size) {
	var pos;
	var k = 1;
	var wdth = Math.sqrt(size);
	var map = MapIndex[mapId];
	var l = map.length
	for (var i=1; i<=wdth; i++) {
		for (var j=1; j<=wdth; j++) {
			pos.isoXYZ.z = 0
			pos.isoXYZ.x = i
			pos.isoXYZ.y = j
			k = Math.floor((Math.random()*l)+1);
			gkIsoCreateSingleDiamond(pos.isoXYZ, map[k]);
			Rendered[k] = map[k];
		}
	}
}

function gkRestorePixel (xv,yv,size) {
	var pos.isoXYZ.x = xv;
	pos.isoXYZ.y = yv;
	pos.isoXYZ.z = 0;
	gkIsoCreateSingleDiamond(pos.isoXYZ, map[k]);
}
