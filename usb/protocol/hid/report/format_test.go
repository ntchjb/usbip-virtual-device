package report_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report"
	"github.com/stretchr/testify/assert"
)

func TestHIDReportDataStringParseMap(t *testing.T) {
	tests := []struct {
		name        string
		tag         report.HIDReportTag
		globalState report.HIDReportGlobalState
		item        []byte
		expected    string
	}{
		{
			name:        "HID_REPORT_TAG_INPUT Constant",
			tag:         report.HID_REPORT_TAG_INPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x03,
			},
			expected: "Constant",
		},
		{
			name:        "HID_REPORT_TAG_INPUT DataAll",
			tag:         report.HID_REPORT_TAG_INPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFE, 0x01,
			},
			expected: "Data,Variable,Relative,Wrap,Non Linear,No Preferred,Null State,Buffered Bytes",
		},
		{
			name:        "HID_REPORT_TAG_INPUT Data",
			tag:         report.HID_REPORT_TAG_INPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x00,
			},
			expected: "Data,Array,Absolute",
		},
		{
			name:        "HID_REPORT_TAG_OUTPUT Constant",
			tag:         report.HID_REPORT_TAG_OUTPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x03,
			},
			expected: "Constant",
		},
		{
			name:        "HID_REPORT_TAG_OUTPUT DataAll",
			tag:         report.HID_REPORT_TAG_OUTPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFE, 0x01,
			},
			expected: "Data,Variable,Relative,Wrap,Non Linear,No Preferred,Null State,Buffered Bytes,Volatile",
		},
		{
			name:        "HID_REPORT_TAG_OUTPUT Data",
			tag:         report.HID_REPORT_TAG_OUTPUT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x00,
			},
			expected: "Data,Array,Absolute",
		},
		{
			name:        "HID_REPORT_TAG_FEATURE Constant",
			tag:         report.HID_REPORT_TAG_FEATURE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x03,
			},
			expected: "Constant",
		},
		{
			name:        "HID_REPORT_TAG_FEATURE DataAll",
			tag:         report.HID_REPORT_TAG_FEATURE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFE, 0x01,
			},
			expected: "Data,Variable,Relative,Wrap,Non Linear,No Preferred,Null State,Buffered Bytes,Volatile",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Physical",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x00,
			},
			expected: "Physical (group of axes)",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Application",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x01,
			},
			expected: "Application (mouse, keyboard)",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Logical",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x02,
			},
			expected: "Logical (interrelated data)",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Report",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x03,
			},
			expected: "Report",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION NamedArray",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x04,
			},
			expected: "Named Array",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION UsageSwitch",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "Usage Switch",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION UsageModifier",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x06,
			},
			expected: "Usage Modifier",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Reserved",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "Reserved",
		},
		{
			name:        "HID_REPORT_TAG_COLLECTION Vendor",
			tag:         report.HID_REPORT_TAG_COLLECTION,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFE,
			},
			expected: "Vendor-defined",
		},
		{
			name:        "HID_REPORT_TAG_USAGE_PAGE Name",
			tag:         report.HID_REPORT_TAG_USAGE_PAGE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "Keyboard/Keypad",
		},
		{
			name:        "HID_REPORT_TAG_USAGE_PAGE Vendor",
			tag:         report.HID_REPORT_TAG_USAGE_PAGE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFE, 0xFF,
			},
			expected: "Vendor",
		},
		{
			name:        "HID_REPORT_TAG_USAGE_PAGE Reserved",
			tag:         report.HID_REPORT_TAG_USAGE_PAGE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xD3, 0xF1,
			},
			expected: "Reserved",
		},
		{
			name:        "HID_REPORT_TAG_LOGICAL_MINIMUM positive",
			tag:         report.HID_REPORT_TAG_LOGICAL_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0x7F,
			},
			expected: "32767",
		},
		{
			name:        "HID_REPORT_TAG_LOGICAL_MINIMUM negative",
			tag:         report.HID_REPORT_TAG_LOGICAL_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_LOGICAL_MAXIMUM positive",
			tag:         report.HID_REPORT_TAG_LOGICAL_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0x7F,
			},
			expected: "32767",
		},
		{
			name:        "HID_REPORT_TAG_LOGICAL_MAXIMUM negative",
			tag:         report.HID_REPORT_TAG_LOGICAL_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_PHYSICAL_MINIMUM positive",
			tag:         report.HID_REPORT_TAG_PHYSICAL_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0x7F,
			},
			expected: "32767",
		},
		{
			name:        "HID_REPORT_TAG_PHYSICAL_MINIMUM negative",
			tag:         report.HID_REPORT_TAG_PHYSICAL_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_PHYSICAL_MAXIMUM positive",
			tag:         report.HID_REPORT_TAG_PHYSICAL_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0x7F,
			},
			expected: "32767",
		},
		{
			name:        "HID_REPORT_TAG_PHYSICAL_MAXIMUM negative",
			tag:         report.HID_REPORT_TAG_PHYSICAL_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_UNIT_EXPONENT positive",
			tag:         report.HID_REPORT_TAG_UNIT_EXPONENT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0x7F,
			},
			expected: "32767",
		},
		{
			name:        "HID_REPORT_TAG_UNIT_EXPONENT negative",
			tag:         report.HID_REPORT_TAG_UNIT_EXPONENT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF, 0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_UNIT_EXPONENT positive",
			tag:         report.HID_REPORT_TAG_UNIT_EXPONENT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "7",
		},
		{
			name:        "HID_REPORT_TAG_UNIT_EXPONENT negative",
			tag:         report.HID_REPORT_TAG_UNIT_EXPONENT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF,
			},
			expected: "-1",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System1 All",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x31, 0x33, 0x33, 0x03,
			},
			expected: "cm^3 g^3 s^3 K^3 A^3 cd^3",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System1 All Negative",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xF1, 0xFF, 0xFF, 0x0F,
			},
			expected: "cm^-1 g^-1 s^-1 K^-1 A^-1 cd^-1",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System0 All",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x30, 0x33, 0x33, 0x03,
			},
			expected: "None",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System2 All",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x72, 0x77, 0x77, 0x07,
			},
			expected: "rad^7 g^7 s^7 K^7 A^7 cd^7",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System2 All Negative",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x82, 0x88, 0x88, 0x08,
			},
			expected: "rad^-8 g^-8 s^-8 K^-8 A^-8 cd^-8",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System3 All",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x73, 0x77, 0x77, 0x07,
			},
			expected: "in^7 slug^7 s^7 °F^7 A^7 cd^7",
		},
		{
			name:        "HID_REPORT_TAG_UNIT System4 All",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x74, 0x77, 0x77, 0x07,
			},
			expected: "deg^7 slug^7 s^7 °F^7 A^7 cd^7",
		},
		{
			name:        "HID_REPORT_TAG_UNIT SystemUnknown",
			tag:         report.HID_REPORT_TAG_UNIT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x78, 0x77, 0x77, 0x07,
			},
			expected: "Unknown System: 0x08",
		},
		{
			name:        "HID_REPORT_TAG_REPORT_SIZE positive",
			tag:         report.HID_REPORT_TAG_REPORT_SIZE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "7",
		},
		{
			name:        "HID_REPORT_TAG_REPORT_SIZE 16bit",
			tag:         report.HID_REPORT_TAG_REPORT_SIZE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x12, 0xFF,
			},
			expected: "65298",
		},
		{
			name:        "HID_REPORT_TAG_REPORT_ID positive",
			tag:         report.HID_REPORT_TAG_REPORT_ID,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "7",
		},
		{
			name:        "HID_REPORT_TAG_REPORT_ID largest",
			tag:         report.HID_REPORT_TAG_REPORT_ID,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0xFF,
			},
			expected: "255",
		},
		{
			name:        "HID_REPORT_TAG_REOPORT_COUNT positive",
			tag:         report.HID_REPORT_TAG_REOPORT_COUNT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07,
			},
			expected: "7",
		},
		{
			name:        "HID_REPORT_TAG_REOPORT_COUNT 16bit",
			tag:         report.HID_REPORT_TAG_REOPORT_COUNT,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x07, 0x45,
			},
			expected: "17671",
		},
		{
			name: "HID_REPORT_TAG_USAGE withGlobalState",
			tag:  report.HID_REPORT_TAG_USAGE,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0x1D,
			},
			expected: "Keyboard Z",
		},
		{
			name:        "HID_REPORT_TAG_USAGE withoutGlobalState",
			tag:         report.HID_REPORT_TAG_USAGE,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x1D, 0x00, 0x07, 0x00,
			},
			expected: "Keyboard Z",
		},
		{
			name: "HID_REPORT_TAG_USAGE Reserved",
			tag:  report.HID_REPORT_TAG_USAGE,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0xF0, 0xFF,
			},
			expected: "Reserved",
		},
		{
			name: "HID_REPORT_TAG_USAGE_MINIMUM withGlobalState",
			tag:  report.HID_REPORT_TAG_USAGE_MINIMUM,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0x1D,
			},
			expected: "Keyboard Z (0x07,0x1D)",
		},
		{
			name:        "HID_REPORT_TAG_USAGE_MINIMUM withoutGlobalState",
			tag:         report.HID_REPORT_TAG_USAGE_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x1D, 0x00, 0x07, 0x00,
			},
			expected: "Keyboard Z (0x07,0x1D)",
		},
		{
			name: "HID_REPORT_TAG_USAGE_MINIMUM Reserved",
			tag:  report.HID_REPORT_TAG_USAGE_MINIMUM,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0xF0, 0xFF,
			},
			expected: "Reserved (0x07,0xFFF0)",
		},
		{
			name: "HID_REPORT_TAG_USAGE_MAXIMUM withGlobalState",
			tag:  report.HID_REPORT_TAG_USAGE_MAXIMUM,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0x1D,
			},
			expected: "Keyboard Z (0x07,0x1D)",
		},
		{
			name:        "HID_REPORT_TAG_USAGE_MAXIMUM withoutGlobalState",
			tag:         report.HID_REPORT_TAG_USAGE_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x1D, 0x00, 0x07, 0x00,
			},
			expected: "Keyboard Z (0x07,0x1D)",
		},
		{
			name: "HID_REPORT_TAG_USAGE_MAXIMUM Reserved",
			tag:  report.HID_REPORT_TAG_USAGE_MAXIMUM,
			globalState: report.HIDReportGlobalState{
				UsagePage: 0x07,
			},
			item: []byte{
				0xF0, 0xFF,
			},
			expected: "Reserved (0x07,0xFFF0)",
		},
		{
			name:        "HID_REPORT_TAG_DESIGNATOR_INDEX",
			tag:         report.HID_REPORT_TAG_DESIGNATOR_INDEX,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_DESIGNATOR_MINIMUM",
			tag:         report.HID_REPORT_TAG_DESIGNATOR_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_DESIGNATOR_MAXIMUM",
			tag:         report.HID_REPORT_TAG_DESIGNATOR_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_STRING_INDEX",
			tag:         report.HID_REPORT_TAG_STRING_INDEX,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_STRING_MINIMUM",
			tag:         report.HID_REPORT_TAG_STRING_MINIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_STRING_MAXIMUM",
			tag:         report.HID_REPORT_TAG_STRING_MAXIMUM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x05,
			},
			expected: "5",
		},
		{
			name:        "HID_REPORT_TAG_DELIMITER UsageSwitch",
			tag:         report.HID_REPORT_TAG_DELIMITER,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x00,
			},
			expected: "Open Set",
		},
		{
			name:        "HID_REPORT_TAG_DELIMITER UsageSwitch",
			tag:         report.HID_REPORT_TAG_DELIMITER,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x01,
			},
			expected: "Close Set",
		},
		{
			name:        "HID_REPORT_TAG_LONG_ITEM UsageSwitch",
			tag:         report.HID_REPORT_TAG_LONG_ITEM,
			globalState: report.HIDReportGlobalState{},
			item: []byte{
				0x0A,                                                       // bSize
				0x03,                                                       // bTag
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, // data
			},
			expected: "Tag: 0x03, Data: 0x0A090807060504030201",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			out := report.HIDReportDataStringParserMap[test.tag](test.globalState, test.item)

			assert.Equal(t, test.expected, out)
		})
	}
}
