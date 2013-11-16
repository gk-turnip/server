
var gkTerrainEditContext = new gkTerrainEditContextDef();

function gkTerrainEditContextDef() {
	this.terrainEditAddTile = false;
	this.terrainEditRemoveTile = false;
	this.defaultTerrainName = "fern_2d";
	this.editTileIdPrefix = "et_";
	this.attributeList = ["none","wall","field"];
	this.currentAttributeIndex = 0;
}

function gkTerrainEditSetAddTileOn() {
	gkTerrainEditContext.terrainEditAddTile = true;
}

function gkTerrainEditSetAddTileOff() {
	gkTerrainEditContext.terrainEditAddTile = false;
}

function gkTerrainEditIsAddTileOn() {
	return gkTerrainEditContext.terrainEditAddTile;
}

function gkTerrainEditSetRemoveTileOn() {
	gkTerrainEditContext.terrainEditRemoveTile = true;
}

function gkTerrainEditSetRemoveTileOff() {
	gkTerrainEditContext.terrainEditRemoveTile = false;
}

function gkTerrainEditIsRemoveTileOn() {
	return gkTerrainEditContext.terrainEditRemoveTile;
}

function gkTerrainEditNeedClick() {
	return gkTerrainEditContext.terrainEditAddTile || gkTerrainEditContext.terrainEditRemoveTile || (gkTerrainEditContext.currentAttributeIndex > 0);
}

function gkTerrainEditHandleClick(x, y) {
	console.log("gkTerrainEditHandleClick: " + x + "," + y);

	if (gkTerrainEditIsAddTileOn()) {
		gkTerrainAddTile(x, y, gkTerrainEditContext.defaultTerrainName, gkTerrainEditContext.editTileIdPrefix);
	} else {
		if (gkTerrainEditIsRemoveTileOn()) {
			gkTerrainRemoveTile(x, y, gkTerrainEditContext.editTileIdPrefix);
		} else {
			if (gkTerrainEditContext.currentAttributeIndex > 0) {
				gkTerrainSetAttribute(x, y, gkTerrainEditContext.editTileIdPrefix);
			} else {
				console.error("internal error, gkTerrainEditHandleClick invalid state");
			}
		}
	}
}

function gkTerrainEditSetAttributeIndex(attributeIndex) {
	gkTerrainEditContext.currentAttributeIndex = attributeIndex;
}

function gkTerrainEditGetAttributeIndex() {
	return gkTerrainEditContext.currentAttributeIndex;
}

function gkTerrainEditGetAttributeText() {
	return gkTerrainEditContext.attributeList[gkTerrainEditContext.currentAttributeIndex];
}

function gkTerrainEditGetAttributeCount() {
	return gkTerrainEditContext.attributeList.length;
}

