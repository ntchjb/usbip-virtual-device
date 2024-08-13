package usage

import "github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/common"

func ParseUsagePageID(currentUsagePageID uint16, item []byte) (usagePageID UsagePageID, usageID UsageID) {
	usagePageID = UsagePageID(currentUsagePageID)
	usagePageIDAndID := common.ParseUint(item)
	if usagePageIDAndID > 0xFFFF {
		usagePageID = UsagePageID(usagePageIDAndID >> 16)
	}

	return usagePageID, UsageID(usagePageIDAndID & 0xFFFF)
}
