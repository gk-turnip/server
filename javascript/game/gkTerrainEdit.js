
var gkTerrainEditContext = new gkTerrainEditContextDef();

function gkTerrainEditContextDef() {
	this.terrainEditAddTile = false;
	this.terrainEditRemoveTile = false;
	this.defaultTerrainName = "fern_2d";
	this.editTileIdPrefix = "et_";
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
	return gkTerrainEditContext.terrainEditAddTile || gkTerrainEditContext.terrainEditRemoveTile;
}

function gkTerrainEditHandleClick(x, y) {
	console.log("gkTerrainEditHandleClick: " + x + "," + y);

	if (gkTerrainEditIsAddTileOn()) {
		gkTerrainAddTile(x, y, gkTerrainEditContext.defaultTerrainName, gkTerrainEditContext.editTileIdPrefix);
	} else {
		if (gkTerrainEditIsRemoveTileOn()) {
			gkTerrainRemoveTile(x, y, gkTerrainEditContext.editTileIdPrefix);
		} else {
			console.error("one of the terrain edit tile should be on");
		}
	}
}

