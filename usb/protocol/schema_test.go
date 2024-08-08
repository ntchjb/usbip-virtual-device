package protocol_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSetupRequestType(t *testing.T) {
	setup := protocol.SetupRequestType(0b10100001)

	assert.Equal(t, protocol.SETUP_DATA_DIRECTION_OUT, setup.Direction())
	assert.Equal(t, protocol.SETUP_DATA_TYPE_CLASS, setup.Type())
	assert.Equal(t, protocol.SETUP_RECIPIENT_INTERFACE, setup.Recipient())

	setup.SetDirection(protocol.SETUP_DATA_DIRECTION_IN)
	setup.SetType(protocol.SETUP_DATA_TYPE_VENDOR)
	setup.SetRecipient(protocol.SETUP_RECIPIENT_ENDPOINT)

	assert.Equal(t, protocol.SETUP_DATA_DIRECTION_IN, setup.Direction())
	assert.Equal(t, protocol.SETUP_DATA_TYPE_VENDOR, setup.Type())
	assert.Equal(t, protocol.SETUP_RECIPIENT_ENDPOINT, setup.Recipient())
}
