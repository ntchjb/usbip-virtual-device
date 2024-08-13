package descriptor_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/descriptor"
	"github.com/stretchr/testify/assert"
)

func TestDescriptorTypeAndIndex(t *testing.T) {
	descType, index := descriptor.GetDescriptorTypeAndIndex(0x0212)

	assert.Equal(t, descriptor.DESCRIPTOR_TYPE_CONFIGURATION, descType)
	assert.Equal(t, uint8(0x12), index)
}
