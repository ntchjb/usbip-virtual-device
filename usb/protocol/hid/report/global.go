package report

import "github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/common"

type HIDReportGlobalState struct {
	UsagePage       uint16
	LogicalMinimum  int32
	LogicalMaximum  int32
	PhysicalMinimum int32
	PhysicalMaximum int32
	UnitExponent    int8
	Unit            HIDReportUnitExponent
	ReportSize      uint32
	ReportID        uint8
	ReportCount     uint32

	Stack []HIDReportGlobalState
}

type HIDGlobalStateUpdater func(globalState *HIDReportGlobalState, item []byte)

var (
	HIDReportGlobalStateUpdaterMap = []HIDGlobalStateUpdater{
		HID_REPORT_TAG_USAGE_PAGE: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.UsagePage = uint16(common.ParseUint(item))
		},
		HID_REPORT_TAG_LOGICAL_MINIMUM: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.LogicalMinimum = common.ParseInt(item)
		},
		HID_REPORT_TAG_LOGICAL_MAXIMUM: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.LogicalMaximum = common.ParseInt(item)
		},
		HID_REPORT_TAG_PHYSICAL_MINIMUM: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.PhysicalMinimum = common.ParseInt(item)
		},
		HID_REPORT_TAG_PHYSICAL_MAXIMUM: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.PhysicalMaximum = common.ParseInt(item)
		},
		HID_REPORT_TAG_UNIT_EXPONENT: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.UnitExponent = int8(common.ParseInt(item))
		},
		HID_REPORT_TAG_UNIT: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.Unit = ParseUnits(item)
		},
		HID_REPORT_TAG_REPORT_SIZE: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.ReportSize = common.ParseUint(item)
		},
		HID_REPORT_TAG_REPORT_ID: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.ReportID = uint8(common.ParseUint(item))
		},
		HID_REPORT_TAG_REPORT_COUNT: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.ReportCount = common.ParseUint(item)
		},
		HID_REPORT_TAG_PUSH: func(globalState *HIDReportGlobalState, item []byte) {
			globalState.Stack = append(globalState.Stack, *globalState)
		},
		HID_REPORT_TAG_POP: func(globalState *HIDReportGlobalState, item []byte) {
			newLength := len(globalState.Stack) - 1
			state := globalState.Stack[newLength]
			globalState.Stack = globalState.Stack[:newLength]

			globalState.LogicalMaximum = state.LogicalMaximum
			globalState.LogicalMinimum = state.LogicalMinimum
			globalState.PhysicalMaximum = state.PhysicalMaximum
			globalState.PhysicalMinimum = state.PhysicalMinimum
			globalState.ReportCount = state.ReportCount
			globalState.ReportID = state.ReportID
			globalState.ReportSize = state.ReportSize
			globalState.Unit = state.Unit
			globalState.UnitExponent = state.UnitExponent
			globalState.UsagePage = state.UsagePage
		},
	}
)
