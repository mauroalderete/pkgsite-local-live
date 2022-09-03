#!/bin/sh

reloader-proxy --origin http://localhost:$PKGSITE_PORT --endpoint http://0.0.0.0:$PROXY_PORT --snippet $APPDIR/websocket.html --reloadEndpoint http://0.0.0.0:$RELOAD_PORT &
goat -c $APPDIR/goat.yml -i 500
