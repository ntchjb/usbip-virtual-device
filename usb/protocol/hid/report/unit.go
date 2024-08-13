package hid

import (
	"strconv"
	"strings"
)

type HIDReportUnitSystem uint8

const (
	HID_REPORT_UNIT_SYSTEM_NONE             HIDReportUnitSystem = 0
	HID_REPORT_UNIT_SYSTEM_SI_LINEAR        HIDReportUnitSystem = 1
	HID_REPORT_UNIT_SYSTEM_SI_ROTATION      HIDReportUnitSystem = 2
	HID_REPORT_UNIT_SYSTEM_ENGLISH_LINEAR   HIDReportUnitSystem = 3
	HID_REPORT_UNIT_SYSTEM_ENGLISH_ROTATION HIDReportUnitSystem = 4
)

type HIDReportUnitName struct {
	Length            string
	Mass              string
	Time              string
	Temperature       string
	Current           string
	LuminousIntensity string
}

var (
	HIDReportUnitMap = map[HIDReportUnitSystem]HIDReportUnitName{
		HID_REPORT_UNIT_SYSTEM_NONE: {
			Length:            "none",
			Mass:              "none",
			Time:              "none",
			Temperature:       "none",
			Current:           "none",
			LuminousIntensity: "none",
		},
		HID_REPORT_UNIT_SYSTEM_SI_LINEAR: {
			Length:            "cm", // Centimeter
			Mass:              "g",  // Gram
			Time:              "s",  // Second
			Temperature:       "K",  // Kelvin
			Current:           "A",  // Ampere
			LuminousIntensity: "cd", // Candela
		},
		HID_REPORT_UNIT_SYSTEM_SI_ROTATION: {
			Length:            "rad", // Radians
			Mass:              "g",   // Gram
			Time:              "s",   // Second
			Temperature:       "K",   // Kelvin
			Current:           "A",   // Ampere
			LuminousIntensity: "cd",  // Candela
		},
		HID_REPORT_UNIT_SYSTEM_ENGLISH_LINEAR: {
			Length:            "in",   // Inch
			Mass:              "slug", // Slug
			Time:              "s",    // Second
			Temperature:       "°F",   // Fahrenheit
			Current:           "A",    // Ampere
			LuminousIntensity: "cd",   // Candela
		},
		HID_REPORT_UNIT_SYSTEM_ENGLISH_ROTATION: {
			Length:            "°",    // Degrees
			Mass:              "slug", // Slug
			Time:              "s",    // Second
			Temperature:       "°F",   // Fahrenheit
			Current:           "A",    // Ampere
			LuminousIntensity: "cd",   // Candela
		},
	}
)

type HIDReportUnitExponent struct {
	Length            int8
	Mass              int8
	Time              int8
	Temperature       int8
	Current           int8
	LuminousIntensity int8
}

func buildUnitItemString(builder *strings.Builder, name string, exponentValue int8) {
	builder.WriteString(name)
	builder.WriteRune('^')
	builder.WriteString(strconv.FormatInt(int64(exponentValue), 10))
	builder.WriteRune(' ')
}

func ParseUnits(item []byte) HIDReportUnitExponent {
	return HIDReportUnitExponent{
		Length:            int8((item[0] & 0b1111_0000) >> 4),
		Mass:              int8(item[1] & 0b0000_1111),
		Time:              int8((item[1] & 0b1111_0000) >> 4),
		Temperature:       int8(item[2] & 0b0000_1111),
		Current:           int8((item[2] & 0b1111_0000) >> 4),
		LuminousIntensity: int8(item[3] & 0b0000_1111),
	}
}
