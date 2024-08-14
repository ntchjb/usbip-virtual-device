package report_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report"
	"github.com/stretchr/testify/assert"
)

func TestHIDReportDescriptor_String(t *testing.T) {
	tests := []struct {
		name string
		desc []byte
		out  string
		err  error
	}{
		{
			name: "Mouse",
			desc: []byte{
				0x05, 0x01,
				0x09, 0x02,
				0xA1, 0x01,
				0x09, 0x01,
				0xA1, 0x00,
				0x05, 0x09,
				0x19, 0x01,
				0x29, 0x03,
				0x15, 0x00,
				0x25, 0x01,
				0x95, 0x03,
				0x75, 0x01,
				0x81, 0x02,
				0x95, 0x01,
				0x75, 0x05,
				0x81, 0x01,
				0x05, 0x01,
				0x09, 0x30,
				0x09, 0x31,
				0x15, 0x81,
				0x25, 0x7F,
				0x75, 0x08,
				0x95, 0x02,
				0x81, 0x06,
				0xC0,
				0xC0,
			},
			out: `Usage Page (Generic Desktop)
Usage (Mouse)
Collection (Application (mouse, keyboard))
	Usage (Pointer)
	Collection (Physical (group of axes))
		Usage Page (Button)
		Usage Minimum (Button-1 (0x09,0x01))
		Usage Maximum (Button-3 (0x09,0x03))
		Logical Minimum (0)
		Logical Maximum (1)
		Report Count (3)
		Report Size (1)
		Input (Data,Variable,Absolute)
		Report Count (1)
		Report Size (5)
		Input (Constant)
		Usage Page (Generic Desktop)
		Usage (X)
		Usage (Y)
		Logical Minimum (-127)
		Logical Maximum (127)
		Report Size (8)
		Report Count (2)
		Input (Data,Variable,Relative)
	End Collection
End Collection
`,
			err: nil,
		},
		{
			name: "Keyboard",
			desc: []byte{
				0x05, 0x01,
				0x09, 0x06,
				0xA1, 0x01,
				0x05, 0x07,
				0x19, 0xE0,
				0x29, 0xE7,
				0x15, 0x00,
				0x25, 0x01,
				0x75, 0x01,
				0x95, 0x08,
				0x81, 0x02,
				0x95, 0x01,
				0x75, 0x08,
				0x81, 0x01,
				0x95, 0x05,
				0x75, 0x01,
				0x05, 0x08,
				0x19, 0x01,
				0x29, 0x05,
				0x91, 0x02,
				0x95, 0x01,
				0x75, 0x03,
				0x91, 0x01,
				0x95, 0x06,
				0x75, 0x08,
				0x15, 0x00,
				0x25, 0x65,
				0x05, 0x07,
				0x19, 0x00,
				0x29, 0x65,
				0x81, 0x00,
				0xC0,
			},
			out: `Usage Page (Generic Desktop)
Usage (Keyboard)
Collection (Application (mouse, keyboard))
	Usage Page (Keyboard/Keypad)
	Usage Minimum (Keyboard LeftControl (0x07,0xE0))
	Usage Maximum (Keyboard Right GUI (0x07,0xE7))
	Logical Minimum (0)
	Logical Maximum (1)
	Report Size (1)
	Report Count (8)
	Input (Data,Variable,Absolute)
	Report Count (1)
	Report Size (8)
	Input (Constant)
	Report Count (5)
	Report Size (1)
	Usage Page (LED)
	Usage Minimum (Num Lock (0x08,0x01))
	Usage Maximum (Kana (0x08,0x05))
	Output (Data,Variable,Absolute)
	Report Count (1)
	Report Size (3)
	Output (Constant)
	Report Count (6)
	Report Size (8)
	Logical Minimum (0)
	Logical Maximum (101)
	Usage Page (Keyboard/Keypad)
	Usage Minimum (Reserved (0x07,0x00))
	Usage Maximum (Keyboard Application (0x07,0x65))
	Input (Data,Array,Absolute)
End Collection
`,
			err: nil,
		},
		{
			name: "No Data",
			desc: []byte{},
			out:  "",
			err:  report.ErrEmptyData,
		},
		{
			name: "Long Item",
			desc: []byte{
				0x05, 0x01,
				0x09, 0x06,
				0xFE, 0x05, 0xA1, 0x01, 0x02, 0x03, 0x04, 0x05,
			},
			out: `Usage Page (Generic Desktop)
Usage (Keyboard)
Long Item (Tag: 0xA1, Data: 0x0504030201)
`,
			err: nil,
		},
		{
			name: "Unknown Item",
			desc: []byte{
				0xF3, 0x01, 0x02, 0x03, 0x0A,
			},
			out: `Unknown Item: 0xF0 (0x0A030201)
`,
			err: nil,
		},
	}

	for _, test := range tests {
		desc := report.HIDReportDescriptor(test.desc)
		out, err := desc.String()

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.out, out)
	}
}
