package hid

import (
	"encoding/hex"
	"fmt"
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
		return "", fmt.Errorf("empty data")
	}

	for cursor < len(h) {
		prefix := h.GetItemPrefix(cursor)
		tag := HIDReportTag((prefix.BTag << 4) | (uint8(prefix.BType) << 2))
		dataLength := hidReportItemSize[prefix.BSize]
		var detail string

		if tag == HID_REPORT_TAG_LONG_ITEM {
			if cursor+1 >= len(h) {
				return "", fmt.Errorf("data is too short for long item, cannot read data length")
			}
			dataLength = int(h[cursor+1]) + 2
		}
		if cursor+dataLength >= len(h) {
			return "", fmt.Errorf("data is too short, need parsing tag: %x, need to read from index %d with size %d", tag, cursor, hidReportItemSize[prefix.BSize])
		}

		if prefix.BType == HID_REPORT_TYPE_GLOBAL {
			if int(tag) < len(HIDReportGlobalStateUpdaterMap) && HIDReportGlobalStateUpdaterMap[tag] != nil {
				HIDReportGlobalStateUpdaterMap[tag](&globalState, h[cursor+1:cursor+dataLength])
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
			builder.WriteString("Unknown Item: 0x")
			builder.WriteString(hex.EncodeToString([]byte{byte(tag)}))
		}
		// #3: Write data part of that item
		if int(tag) < len(HIDReportDataStringParserMap) && HIDReportDataStringParserMap[tag] != nil {
			detail = HIDReportDataStringParserMap[tag](globalState, h[cursor+1:cursor+dataLength])
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
		} else if tag == HID_REPORT_TAG_END_COLLECTION {
			currTabs--
		}
		cursor += 1 + dataLength
	}

	return builder.String(), nil
}
