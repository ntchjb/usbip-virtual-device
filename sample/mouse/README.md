# Virtual Mouse

This is USB/IP Server that implements Mouse device, which emits mouse location events to between (5,5) and (-5,-5), that is, mouse cursor should move from top-left to bottom-right, and bottom-right to top-left, and so on.

This mouse have polling rate of ~10Hz to not consuming too much CPU :P.

## Usage

```sh
cd sample/mouse

# Start server at 127.0.0.1:3240
go run .

# Attach mouse to begin moving mouse cursor.
sudo usbip --debug attach -r 127.0.0.1 -b 1-1

# Detach the mouse device to stop moving mouse cursor
sudo usbip detach --port=0
```

Feel free to play with `processHIDData` function in `sample/mouse/mouse.go` to move mouse to different directions.
