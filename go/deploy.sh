sudo /etc/init.d/loginServer stop
sudo /etc/init.d/gameServer stop

cp gk/bin/loginServerMain /home/diver2/loginServer/bin
cp gk/bin/gameServerMain /home/diver2/gameServer/bin

cp -r ../stylesheets /var/www/gourdianknot/assets/gk
cp -r ../javascript /var/www/gourdianknot/assets/gk

sudo /etc/init.d/loginServer start
sudo /etc/init.d/gameServer start
