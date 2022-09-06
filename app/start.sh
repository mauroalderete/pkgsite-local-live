#!/bin/sh

reloader --origin http://localhost:$PKGSITE_PORT --public http://0.0.0.0:$PROXY_PORT --snippet $APPDIR/websocket.html &
goat -c $APPDIR/goat.yml -i 500
