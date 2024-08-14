package report

import (
	"strings"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/common"
)

type HIDReportInputData struct {
	IsConstant                bool
	IsVariableOrArray         bool
	IsRelativeOrAbsolute      bool
	IsWrap                    bool
	IsNonLinear               bool
	IsNoPreferredState        bool
	IsNullState               bool
	IsBufferedBytesOrBitField bool
}

func (h HIDReportInputData) String() string {
	var builder strings.Builder

	if h.IsConstant {
		builder.WriteString("Constant")
		return builder.String()
	} else {
		builder.WriteString("Data")
	}
	if h.IsVariableOrArray {
		builder.WriteString(",Variable")
	} else {
		builder.WriteString(",Array")
	}
	if h.IsRelativeOrAbsolute {
		builder.WriteString(",Relative")
	} else {
		builder.WriteString(",Absolute")
	}
	if h.IsWrap {
		builder.WriteString(",Wrap")
	}
	if h.IsNonLinear {
		builder.WriteString(",Non Linear")
	}
	if h.IsNoPreferredState {
		builder.WriteString(",No Preferred")
	}
	if h.IsNullState {
		builder.WriteString(",Null State")
	}
	if h.IsBufferedBytesOrBitField {
		builder.WriteString(",Buffered Bytes")
	}

	return builder.String()
}

type HIDReportOutputData struct {
	HIDReportInputData
	IsVolatile bool
}

func (h HIDReportOutputData) String() string {
	var builder strings.Builder
	builder.WriteString(h.HIDReportInputData.String())
	if h.IsVolatile {
		builder.WriteString(",Volatile")
	}

	return builder.String()
}

type HIDReportFeatureData struct {
	HIDReportOutputData
}

type HIDReportCollectionData uint8

const (
	HID_REPORT_COLLECTION_PHYSICAL       = 0x00
	HID_REPORT_COLLECTION_APPLICATION    = 0x01
	HID_REPORT_COLLECTION_LOGICAL        = 0x02
	HID_REPORT_COLLECTION_REPORT         = 0x03
	HID_REPORT_COLLECTION_NAMED_ARRAY    = 0x04
	HID_REPORT_COLLECTION_USAGE_SWITCH   = 0x05
	HID_REPORT_COLLECTION_USAGE_MODIFIER = 0x06
)

var (
	HIDReportCollectionNames = map[HIDReportCollectionData]string{
		HID_REPORT_COLLECTION_PHYSICAL:       "Physical (group of axes)",
		HID_REPORT_COLLECTION_APPLICATION:    "Application (mouse, keyboard)",
		HID_REPORT_COLLECTION_LOGICAL:        "Logical (interrelated data)",
		HID_REPORT_COLLECTION_REPORT:         "Report",
		HID_REPORT_COLLECTION_NAMED_ARRAY:    "Named Array",
		HID_REPORT_COLLECTION_USAGE_SWITCH:   "Usage Switch",
		HID_REPORT_COLLECTION_USAGE_MODIFIER: "Usage Modifier",
	}
)

func ParseInputReportItem(item []byte) HIDReportInputData {
	var res HIDReportInputData
	if len(item) > 0 {
		res.IsConstant = common.ByteToBool(item[0] & 0b0000_0001)
		res.IsVariableOrArray = common.ByteToBool((item[0] & 0b0000_0010) >> 1)
		res.IsRelativeOrAbsolute = common.ByteToBool((item[0] & 0b0000_0100) >> 2)
		res.IsWrap = common.ByteToBool((item[0] & 0b0000_1000) >> 3)
		res.IsNonLinear = common.ByteToBool((item[0] & 0b0001_0000) >> 4)
		res.IsNoPreferredState = common.ByteToBool((item[0] & 0b0010_0000) >> 5)
		res.IsNullState = common.ByteToBool((item[0] & 0b0100_0000) >> 6)
	}
	if len(item) > 1 {
		res.IsBufferedBytesOrBitField = common.ByteToBool(item[1] & 0b0000_0001)
	}

	return res
}

func ParseOutputReportItem(item []byte) HIDReportOutputData {
	var res HIDReportOutputData
	res.HIDReportInputData = ParseInputReportItem(item)
	if len(item) > 0 {
		res.IsVolatile = common.ByteToBool((item[0] & 0b1000_0000) >> 7)
	}

	return res
}

func ParseFeatureReportItem(item []byte) HIDReportFeatureData {
	return HIDReportFeatureData{
		HIDReportOutputData: ParseOutputReportItem(item),
	}
}

func ParseCollectionReportItem(item byte) HIDReportCollectionData {
	return HIDReportCollectionData(item)
}
