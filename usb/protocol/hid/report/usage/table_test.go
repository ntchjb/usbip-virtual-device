package usage_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/usage"
	"github.com/stretchr/testify/assert"
)

func TestCreateIndexedUsageTable(t *testing.T) {
	table := usage.CreateIndexedUsageTable()

	assert.Equal(t, "Keyboard/Keypad", table.GetUsagePageName(0x07))
	assert.Equal(t, "Keyboard F", table.GetUsageName(0x07, 0x09))

	assert.Equal(t, "Generic Desktop", table.GetUsagePageName(0x01))
	assert.Equal(t, "Mouse", table.GetUsageName(0x01, 0x02))
	assert.Equal(t, "Keyboard", table.GetUsageName(0x01, 0x06))

	assert.Equal(t, "FIDO Alliance", table.GetUsagePageName(0xF1D0))
	assert.Equal(t, "U2F Authenticator Device", table.GetUsageName(0xF1D0, 0x01))

	assert.Equal(t, "Reserved", table.GetUsagePageName(0xF1D1))
	assert.Equal(t, "Reserved", table.GetUsageName(0xF1D1, 0x01))

	assert.Equal(t, "Vendor", table.GetUsagePageName(0xFF00))
	assert.Equal(t, "Vendor", table.GetUsageName(0xFF00, 0x01))

	assert.Equal(t, "Button", table.GetUsagePageName(0x09))
	assert.Equal(t, "Button-65535", table.GetUsageName(0x09, 0xFFFF))

	assert.Equal(t, "Keyboard/Keypad", table.GetUsagePageName(0x07))
	assert.Equal(t, "Reserved", table.GetUsageName(0x07, 0xE8))
}
