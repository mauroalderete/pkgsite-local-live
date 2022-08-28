#!/bin/sh

pkgsite -http "0.0.0.0:8080" $(ls $GOPATH/src/**/go.mod | sed 's/\/go.mod//' | paste -sd ',')
