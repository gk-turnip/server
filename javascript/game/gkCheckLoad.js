//This script should be loaded last
var x = 0;
x += gkCheckLoadIso();
x += gkCheckLoadAvatar();
x += gkCheckLoadAudio();
x += gkCheckLoadRain();
x += gkCheckLoadTerrain();
x += gkCheckLoadWs();
if (x != 41) {
	document.write("Oh my knots! It seems that <strong>something is not working right!</strong> You can try reloading your browser, or clearing your cache. If you still get this message, you can contact support.");
}
