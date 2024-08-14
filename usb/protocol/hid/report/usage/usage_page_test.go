package usage_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/usage"
	"github.com/stretchr/testify/assert"
)

func TestParseUsagePageID(t *testing.T) {
	usagePageID, usageID := usage.ParseUsagePageID(0x07, []byte{0x03})

	assert.Equal(t, usage.UsagePageID(0x07), usagePageID)
	assert.Equal(t, usage.UsageID(0x03), usageID)

	usagePageID, usageID = usage.ParseUsagePageID(0x07, []byte{0x02, 0x00, 0x01, 0x00})

	assert.Equal(t, usage.UsagePageID(0x01), usagePageID)
	assert.Equal(t, usage.UsageID(0x02), usageID)
}
