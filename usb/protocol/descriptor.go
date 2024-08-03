package protocol

type DescriptorType uint8

const (
	DescriptorTypeDevice                                     DescriptorType = 1
	DescriptorTypeConfiguration                              DescriptorType = 2
	DescriptorTypeString                                     DescriptorType = 3
	DescriptorTypeInterface                                  DescriptorType = 4
	DescriptorTypeEndpoint                                   DescriptorType = 5
	DescriptorTypeInterfacePower                             DescriptorType = 8
	DescriptorTypeOTG                                        DescriptorType = 9
	DescriptorTypeDebug                                      DescriptorType = 10
	DescriptorTypeInterfaceAssociation                       DescriptorType = 11
	DescriptorTypeBOS                                        DescriptorType = 15
	DescriptorTypeDeviceCapability                           DescriptorType = 16
	DescriptorTypeSuperSpeedUSBEndpointCompanion             DescriptorType = 48
	DescriptorTypeSuperSpeedPlusIsochronousEndpointCompanion DescriptorType = 49
)

const (
	DescriptorTypeHID                   DescriptorType = 0x21
	DescriptorTypeHIDReport             DescriptorType = 0x22
	DescriptorTypeHIDPhysicalDescriptor DescriptorType = 0x23
)

const (
	STANDARD_DEVICE_DESCRIPTOR_LENGTH        = 18
	STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH = 9
	STANDARD_INTERFACE_DESCRIPTOR_LENGTH     = 9
	STANDARD_ENDPOINT_DESCRIPTOR_LENGTH      = 7

	HID_DESCRIPTOR_LENGTH = 9
)

const (
	LangIDEnglishUSA uint16 = 0x0409
)

func GetDescriptorTypeAndIndex(wValue uint16) (descriptorType DescriptorType, index uint8) {
	descriptorType = DescriptorType(wValue >> 8)
	index = uint8(wValue & 0x00FF)

	return
}
