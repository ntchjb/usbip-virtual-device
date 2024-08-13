package report

type HIDReportType uint8

const (
	HID_REPORT_TYPE_MAIN     HIDReportType = 0x00
	HID_REPORT_TYPE_GLOBAL   HIDReportType = 0x01
	HID_REPORT_TYPE_LOCAL    HIDReportType = 0x02
	HID_REPORT_TYPE_RESERVED HIDReportType = 0x03
)

type HIDReportTag uint8

const (
	// Main tags
	HID_REPORT_TAG_INPUT          HIDReportTag = HIDReportTag(uint8(0b1000<<4) | (uint8(HID_REPORT_TYPE_MAIN) << 2))
	HID_REPORT_TAG_OUTPUT         HIDReportTag = HIDReportTag(uint8(0b1001<<4) | (uint8(HID_REPORT_TYPE_MAIN) << 2))
	HID_REPORT_TAG_FEATURE        HIDReportTag = HIDReportTag(uint8(0b1011<<4) | (uint8(HID_REPORT_TYPE_MAIN) << 2))
	HID_REPORT_TAG_COLLECTION     HIDReportTag = HIDReportTag(uint8(0b1010<<4) | (uint8(HID_REPORT_TYPE_MAIN) << 2))
	HID_REPORT_TAG_END_COLLECTION HIDReportTag = HIDReportTag(uint8(0b1100<<4) | (uint8(HID_REPORT_TYPE_MAIN) << 2))

	// Global tags
	HID_REPORT_TAG_USAGE_PAGE       HIDReportTag = HIDReportTag(uint8(0b0000<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_LOGICAL_MINIMUM  HIDReportTag = HIDReportTag(uint8(0b0001<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_LOGICAL_MAXIMUM  HIDReportTag = HIDReportTag(uint8(0b0010<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_PHYSICAL_MINIMUM HIDReportTag = HIDReportTag(uint8(0b0011<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_PHYSICAL_MAXIMUM HIDReportTag = HIDReportTag(uint8(0b0100<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_UNIT_EXPONENT    HIDReportTag = HIDReportTag(uint8(0b0101<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_UNIT             HIDReportTag = HIDReportTag(uint8(0b0110<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_REPORT_SIZE      HIDReportTag = HIDReportTag(uint8(0b0111<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_REPORT_ID        HIDReportTag = HIDReportTag(uint8(0b1000<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_REOPORT_COUNT    HIDReportTag = HIDReportTag(uint8(0b1001<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_PUSH             HIDReportTag = HIDReportTag(uint8(0b1010<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))
	HID_REPORT_TAG_POP              HIDReportTag = HIDReportTag(uint8(0b1011<<4) | (uint8(HID_REPORT_TYPE_GLOBAL) << 2))

	HID_REPORT_TAG_USAGE              HIDReportTag = HIDReportTag(uint8(0b0000<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_USAGE_MINIMUM      HIDReportTag = HIDReportTag(uint8(0b0001<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_USAGE_MAXIMUM      HIDReportTag = HIDReportTag(uint8(0b0010<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_DESIGNATOR_INDEX   HIDReportTag = HIDReportTag(uint8(0b0011<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_DESIGNATOR_MINIMUM HIDReportTag = HIDReportTag(uint8(0b0100<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_DESIGNATOR_MAXIMUM HIDReportTag = HIDReportTag(uint8(0b0101<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_STRING_INDEX       HIDReportTag = HIDReportTag(uint8(0b0111<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_STRING_MINIMUM     HIDReportTag = HIDReportTag(uint8(0b1000<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_STRING_MAXIMUM     HIDReportTag = HIDReportTag(uint8(0b1001<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))
	HID_REPORT_TAG_DELIMITER          HIDReportTag = HIDReportTag(uint8(0b1010<<4) | (uint8(HID_REPORT_TYPE_LOCAL) << 2))

	HID_REPORT_TAG_LONG_ITEM HIDReportTag = 0b11111110
)

var (
	HIDReportTagNames = map[HIDReportTag]string{
		HID_REPORT_TAG_INPUT:          "Input",
		HID_REPORT_TAG_OUTPUT:         "Output",
		HID_REPORT_TAG_FEATURE:        "Feature",
		HID_REPORT_TAG_COLLECTION:     "Collection",
		HID_REPORT_TAG_END_COLLECTION: "End Collection",

		HID_REPORT_TAG_USAGE_PAGE:       "Usage Page",
		HID_REPORT_TAG_LOGICAL_MINIMUM:  "Logical Minimum",
		HID_REPORT_TAG_LOGICAL_MAXIMUM:  "Logical Maximum",
		HID_REPORT_TAG_PHYSICAL_MINIMUM: "Physical Minimum",
		HID_REPORT_TAG_PHYSICAL_MAXIMUM: "Physical Maximum",
		HID_REPORT_TAG_UNIT_EXPONENT:    "Unit Exponent",
		HID_REPORT_TAG_UNIT:             "Unit",
		HID_REPORT_TAG_REPORT_SIZE:      "Report Size",
		HID_REPORT_TAG_REPORT_ID:        "Report ID",
		HID_REPORT_TAG_REOPORT_COUNT:    "Report Count",
		HID_REPORT_TAG_PUSH:             "Push",
		HID_REPORT_TAG_POP:              "Pop",

		HID_REPORT_TAG_USAGE:              "Usage",
		HID_REPORT_TAG_USAGE_MINIMUM:      "Usage Minimum",
		HID_REPORT_TAG_USAGE_MAXIMUM:      "Usage Maximum",
		HID_REPORT_TAG_DESIGNATOR_INDEX:   "Designator Index",
		HID_REPORT_TAG_DESIGNATOR_MINIMUM: "Designator Minimum",
		HID_REPORT_TAG_DESIGNATOR_MAXIMUM: "Designator Maximum",
		HID_REPORT_TAG_STRING_INDEX:       "String Index",
		HID_REPORT_TAG_STRING_MINIMUM:     "String Minimum",
		HID_REPORT_TAG_STRING_MAXIMUM:     "String Maximum",
		HID_REPORT_TAG_DELIMITER:          "Delimiter",

		HID_REPORT_TAG_LONG_ITEM: "Long Item",
	}
)
