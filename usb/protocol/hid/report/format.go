package report

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/common"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/usage"
)

type HIDReportDataStringParser func(globalItem HIDReportGlobalState, item []byte) string

func displayUintData(globalState HIDReportGlobalState, item []byte) string {
	return strconv.FormatUint(uint64(common.ParseUint(item)), 10)
}

func displayUintHexData(num uint32) string {
	if num < 0x0100 {
		return fmt.Sprintf("0x%02X", num)
	}
	if num < 0x0001_0000 {
		return fmt.Sprintf("0x%04X", num)
	}
	return fmt.Sprintf("0x%08X", num)
}

func displayIntData(_ HIDReportGlobalState, item []byte) string {
	return strconv.FormatInt(int64(common.ParseInt(item)), 10)
}

var (
	HIDReportDataStringParserMap = []HIDReportDataStringParser{
		HID_REPORT_TAG_INPUT: func(globalState HIDReportGlobalState, item []byte) string {
			res := ParseInputReportItem(item)
			return res.String()
		},
		HID_REPORT_TAG_OUTPUT: func(globalState HIDReportGlobalState, item []byte) string {
			res := ParseOutputReportItem(item)
			return res.String()
		},
		HID_REPORT_TAG_FEATURE: func(globalState HIDReportGlobalState, item []byte) string {
			res := ParseFeatureReportItem(item)
			return res.String()
		},
		HID_REPORT_TAG_COLLECTION: func(globalState HIDReportGlobalState, item []byte) string {
			if name, ok := HIDReportCollectionNames[HIDReportCollectionData(item[0])]; ok {
				return name
			} else if item[0] >= 0x80 && item[0] <= 0xFF {
				return "Vendor-defined"
			}

			return "Reserved"
		},
		// HID_REPORT_TAG_END_COLLECTION:
		HID_REPORT_TAG_USAGE_PAGE: func(globalState HIDReportGlobalState, item []byte) string {
			return usage.IndexedUsageTable.GetUsagePageName(usage.UsagePageID(common.ParseUint(item)))
		},
		HID_REPORT_TAG_LOGICAL_MINIMUM:  displayIntData,
		HID_REPORT_TAG_LOGICAL_MAXIMUM:  displayIntData,
		HID_REPORT_TAG_PHYSICAL_MINIMUM: displayIntData,
		HID_REPORT_TAG_PHYSICAL_MAXIMUM: displayIntData,
		HID_REPORT_TAG_UNIT_EXPONENT:    displayIntData,
		HID_REPORT_TAG_UNIT: func(globalState HIDReportGlobalState, item []byte) string {
			var builder strings.Builder
			unitSystemID := HIDReportUnitSystem(item[0] & 0b0000_1111)
			if unitSystemID == HID_REPORT_UNIT_SYSTEM_NONE {
				return "None"
			}
			unitNames, ok := HIDReportUnitMap[unitSystemID]
			if !ok {
				return "Unknown System: " + displayUintHexData(uint32(item[0]&0b0000_1111))
			}
			units := ParseUnits(item)
			if units.Length != 0 {
				buildUnitItemString(&builder, unitNames.Length, units.Length)
			}
			if units.Mass != 0 {
				buildUnitItemString(&builder, unitNames.Mass, units.Mass)
			}
			if units.Time != 0 {
				buildUnitItemString(&builder, unitNames.Time, units.Time)
			}
			if units.Temperature != 0 {
				buildUnitItemString(&builder, unitNames.Temperature, units.Temperature)
			}
			if units.Current != 0 {
				buildUnitItemString(&builder, unitNames.Current, units.Current)
			}
			if units.LuminousIntensity != 0 {
				buildUnitItemString(&builder, unitNames.LuminousIntensity, units.LuminousIntensity)
			}
			return strings.TrimSpace(builder.String())
		},
		HID_REPORT_TAG_REPORT_SIZE:   displayUintData,
		HID_REPORT_TAG_REPORT_ID:     displayUintData,
		HID_REPORT_TAG_REOPORT_COUNT: displayUintData,
		// HID_REPORT_TAG_PUSH:
		// HID_REPORT_TAG_POP:
		HID_REPORT_TAG_USAGE: func(globalState HIDReportGlobalState, item []byte) string {
			usagePageID, usageID := usage.ParseUsagePageID(globalState.UsagePage, item)
			return usage.IndexedUsageTable.GetUsageName(usagePageID, usageID)
		},
		HID_REPORT_TAG_USAGE_MINIMUM: func(globalState HIDReportGlobalState, item []byte) string {
			usagePageID, usageID := usage.ParseUsagePageID(globalState.UsagePage, item)
			return usage.IndexedUsageTable.GetUsageName(usagePageID, usageID) + " (" + displayUintHexData(uint32(usagePageID)) + "," + displayUintHexData(uint32(usageID)) + ")"
		},
		HID_REPORT_TAG_USAGE_MAXIMUM: func(globalState HIDReportGlobalState, item []byte) string {
			usagePageID, usageID := usage.ParseUsagePageID(globalState.UsagePage, item)
			return usage.IndexedUsageTable.GetUsageName(usagePageID, usageID) + " (" + displayUintHexData(uint32(usagePageID)) + "," + displayUintHexData(uint32(usageID)) + ")"
		},
		HID_REPORT_TAG_DESIGNATOR_INDEX:   displayUintData,
		HID_REPORT_TAG_DESIGNATOR_MINIMUM: displayUintData,
		HID_REPORT_TAG_DESIGNATOR_MAXIMUM: displayUintData,
		HID_REPORT_TAG_STRING_INDEX:       displayUintData,
		HID_REPORT_TAG_STRING_MINIMUM:     displayUintData,
		HID_REPORT_TAG_STRING_MAXIMUM:     displayUintData,
		HID_REPORT_TAG_DELIMITER: func(globalItem HIDReportGlobalState, item []byte) string {
			delimiter := common.ParseUint(item)
			if delimiter == 0 {
				return "Open Set"
			} else {
				return "Close Set"
			}
		},
		HID_REPORT_TAG_LONG_ITEM: func(globalItem HIDReportGlobalState, item []byte) string {
			bDataSize := item[0]
			bLongItemTag := item[1]
			var builder strings.Builder
			reversed := make([]byte, bDataSize)
			copy(reversed, item[2:2+bDataSize])
			slices.Reverse(reversed)

			builder.WriteString("Tag: ")
			builder.WriteString(displayUintHexData(uint32(bLongItemTag)))
			builder.WriteString(", Data: 0x")
			builder.WriteString(fmt.Sprintf("%X", reversed))

			return builder.String()
		},
	}
)
