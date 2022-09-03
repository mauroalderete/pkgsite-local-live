#!/bin/sh

pkgsite -http "0.0.0.0:$PKGSITE_PORT" $(ls $GOPATH/src/**/go.mod | sed 's/\/go.mod//' | paste -sd ',')
