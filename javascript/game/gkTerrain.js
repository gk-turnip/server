
//This is still experimental
//Begin map data
var Map=new Array("PapayaWhip","IndianRed","LightSalmon","Wheat","Salmon","PaleGoldenRod","LightSalmon","Moccasin","NavajoWhite","SaddleBrown","Peru","Tan","Wheat","Moccasin","IndianRed","SandyBrown","PeachPuff");

function gkRenderMap (mapId,size) {
	var pos;
	var k = 1;
	for (var i=1; i<=size; i++) {
		for (var j=1; j<= size; j++) {
			pos.isoXYZ.z = 0
			pos.isoXYZ.x = i
			pos.isoXYZ.y = j
			gkIsoCreateSingleDiamond(pos.isoXYZ, Map[k]);
			++k;
			if (k = 16) {
				k = 1;
			}
		}
	}
}