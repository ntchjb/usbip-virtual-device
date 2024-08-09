# USB/IP Virtual Device Library

This is a library for developing server side of USB/IP to emulate a USB device. The library contains features as follows

- Data schema for encoding/decoding via USB/IP protocol
- Data schema for encoding/decoding USB device descriptors (+ HID device descriptors)
- A Server code for running USB/IP server, with request handling.
- A worker pool to help managing URB requests i.e. unlinking URB, process URB in sequences, etc.
- Device registrar to register multiple devices to the server.

User of this library only need to implement `Device` interface located at `/usb/device.go`, and use Server with registrar to run it. See samples in `/sample` folder.

## Why do we need this?

Virtual USB device can be used for...
- For USB driver software development, to be used with automation testing on clouds
- For emulating access of USB devices without using actual hardware

Ideas of use cases
- [ ] Implement Keyboard device as a sample
- [ ] Implement FIDO2 USB device as a sample
- [ ] Implement virtual audio cable with effect (if possible)
- [ ] Implement virtual USB flash drive using a pre-allocated file as storage
- [ ] Support USB 3.0 and beyond (not sure if running them on CPU is possible or not)

Made by [@ntchjb](https://github.com/ntchjb).
