#!/usr/bin/zsh

sudo usbip --debug attach -r 127.0.0.1 -b 1-1 && sleep 3 && sudo usbip detach --port=0


