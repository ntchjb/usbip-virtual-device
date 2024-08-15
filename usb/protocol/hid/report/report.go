package report

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type HIDReportItemPrefix struct {
	BSize uint8
	BType HIDReportType
	BTag  uint8
}

type HIDReportShortItem struct {
	HIDReportItemPrefix
	Data []uint8
}

type HIDReportLongItem struct {
	BDataSize    uint8
	BLongItemTag uint8
	Data         []uint8
}

type HIDReportLocalState struct {
	Usage             uint32
	UsageMinimum      uint32
	UsageMaximum      uint32
	DesignatorIndex   uint8
	DesignatorMinimum uint8
	DesignatorMaximum uint8
	StringIndex       uint8
	StringMinimum     uint8
	StringMaximum     uint8
	Delimiter         bool
}

var (
	ErrEmptyData = errors.New("empty data")
)

type HIDReportDescriptor []byte

func (h HIDReportDescriptor) GetItemPrefix(idx int) HIDReportItemPrefix {
	return HIDReportItemPrefix{
		BSize: h[idx] & 0b0000_0011,
		BType: HIDReportType((h[idx] & 0b0000_1100) >> 2),
		BTag:  h[idx] >> 4,
	}
}

func (h HIDReportDescriptor) String() (string, error) {
	var builder strings.Builder
	var cursor, currTabs int
	var globalState HIDReportGlobalState
	hidReportItemSize := []int{0, 1, 2, 4}

	if len(h) == 0 {
		return "", ErrEmptyData
	}

	for cursor < len(h) && h[cursor] != 0x00 {
		prefix := h.GetItemPrefix(cursor)
		tag := HIDReportTag((prefix.BTag << 4) | (uint8(prefix.BType) << 2))
		// Handle Long item as special case
		if h[cursor] == byte(HID_REPORT_TAG_LONG_ITEM) {
			tag = HIDReportTag(h[cursor])
		}
		dataLength := hidReportItemSize[prefix.BSize]
		var detail string
		dataStartIdx := cursor + 1

		if tag == HID_REPORT_TAG_LONG_ITEM {
			if dataStartIdx >= len(h) {
				return "", fmt.Errorf("data is too short for long item, cannot read data length")
			}
			dataLength = int(h[dataStartIdx]) + 2
		}
		if dataStartIdx+dataLength > len(h) {
			return "", fmt.Errorf("data is too short, need parsing data for tag: %x, need to read from index %d with size %d", tag, dataStartIdx, dataLength)
		}

		if prefix.BType == HID_REPORT_TYPE_GLOBAL {
			if int(tag) < len(HIDReportGlobalStateUpdaterMap) && HIDReportGlobalStateUpdaterMap[tag] != nil {
				HIDReportGlobalStateUpdaterMap[tag](&globalState, h[dataStartIdx:dataStartIdx+dataLength])
			}
		}

		if tag == HID_REPORT_TAG_END_COLLECTION {
			if currTabs > 0 {
				currTabs--
			}
		}
		// #1: Write indents, if any
		for i := 0; i < currTabs; i++ {
			builder.WriteRune('\t')
		}
		// #2: Write Report Item name
		if itemName, ok := HIDReportTagNames[tag]; ok {
			builder.WriteString(itemName)
		} else {
			builder.WriteString("Unknown Item: ")
			builder.WriteString(fmt.Sprintf("0x%02X", byte(tag)))
		}
		// #3: Write data part of that item
		if int(tag) < len(HIDReportDataStringParserMap) && HIDReportDataStringParserMap[tag] != nil {
			detail = HIDReportDataStringParserMap[tag](globalState, h[dataStartIdx:dataStartIdx+dataLength])
		} else {
			reversed := make([]byte, dataLength)
			copy(reversed, h[dataStartIdx:dataStartIdx+dataLength])
			slices.Reverse(reversed)
			detail = fmt.Sprintf("0x%X", reversed)
		}
		if len(detail) > 0 {
			builder.WriteString(" (")
			builder.WriteString(detail)
			builder.WriteRune(')')
		}
		// #4: Write newline as the end of item
		builder.WriteRune('\n')

		if tag == HID_REPORT_TAG_COLLECTION {
			currTabs++
		}
		cursor = dataStartIdx + dataLength
	}

	return builder.String(), nil
}
