package descriptor

type DescriptorType uint8

const (
	DESCRIPTOR_TYPE_DEVICE                                          DescriptorType = 1
	DESCRIPTOR_TYPE_CONFIGURATION                                   DescriptorType = 2
	DESCRIPTOR_TYPE_STRING                                          DescriptorType = 3
	DESCRIPTOR_TYPE_INTERFACE                                       DescriptorType = 4
	DESCRIPTOR_TYPE_ENDPOINT                                        DescriptorType = 5
	DESCRIPTOR_TYPE_INTERFACE_POWER                                 DescriptorType = 8
	DESCRIPTOR_TYPE_OTG                                             DescriptorType = 9
	DESCRIPTOR_TYPE_DEBUG                                           DescriptorType = 10
	DESCRIPTOR_TYPE_INTERFACE_ASSOCIATION                           DescriptorType = 11
	DESCRIPTOR_TYPE_BOS                                             DescriptorType = 15
	DESCRIPTOR_TYPE_DEVICE_CAPABILITY                               DescriptorType = 16
	DESCRIPTOR_TYPE_SUPER_SPEED_USB_ENDPOINT_COMPANION              DescriptorType = 48
	DESCRIPTOR_TYPE_SUPER_SPEED_PLUS_ISOCHRONOUS_ENDPOINT_COMPANION DescriptorType = 49
)

const (
	DESCRIPTOR_TYPE_HID                     DescriptorType = 0x21
	DESCRIPTOR_TYPE_HID_REPORT              DescriptorType = 0x22
	DESCRIPTOR_TYPE_HIS_PHYSICAL_DESCRIPTOR DescriptorType = 0x23
)

const (
	STANDARD_DEVICE_DESCRIPTOR_LENGTH        = 18
	STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH = 9
	STANDARD_INTERFACE_DESCRIPTOR_LENGTH     = 9
	STANDARD_ENDPOINT_DESCRIPTOR_LENGTH      = 7
)

const (
	LangIDEnglishUSA uint16 = 0x0409
)

func GetDescriptorTypeAndIndex(wValue uint16) (descriptorType DescriptorType, index uint8) {
	descriptorType = DescriptorType(wValue >> 8)
	index = uint8(wValue & 0x00FF)

	return
}
