#!/bin/sh

watcher --inotify "--quiet --recursive --include .go$ -e create -e modify -e delete -e move $GOPATH/src" --exec "echo restart pkgsite now" --keeplive
