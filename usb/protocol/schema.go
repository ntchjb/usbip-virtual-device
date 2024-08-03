package protocol

type BusID [32]byte

const (
	HID_SPEC_VERSION       uint16 = 0x0110
	HID_CLASS_SPEC_VERSION uint16 = 0x0101
)

const (
	SPEED_USB1_LOW uint32 = iota
	SPEED_USB1_FULL
	SPEED_USB2_HIGH
	SPEED_USB2_WIRELESS
	SPEED_USB3_SUPER
	SPEED_USB3_SUPER_PLUS
)

const (
	CLASS_BASEDON_INTERFACE     uint8 = 0x00
	CLASS_AUDIO                 uint8 = 0x01
	CLASS_CDC_CONTROL           uint8 = 0x02
	CLASS_HID                   uint8 = 0x03
	CLASS_PHYSICAL              uint8 = 0x05
	CLASS_IMAGE                 uint8 = 0x06
	CLASS_PRINTER               uint8 = 0x07
	CLASS_MASS_STORAGE          uint8 = 0x08
	CLASS_HUB                   uint8 = 0x09
	CLASS_CDC_DATA              uint8 = 0x0a
	CLASS_SMART_CARD            uint8 = 0x0b
	CLASS_CONTENT_SECURITY      uint8 = 0x0d
	CLASS_VIDEO                 uint8 = 0x0e
	CLASS_PERSONAL_HEALTHCARE   uint8 = 0x0f
	CLASS_AUDIO_AND_VIDEO       uint8 = 0x10
	CLASS_BILLBOARD_DEVICE      uint8 = 0x11
	CLASS_USB_TYPE_C_BRIDGE     uint8 = 0x12
	CLASS_BULK_DISPLAY_PROTOCOL uint8 = 0x13
	CLASS_MCTP_OVER_USB         uint8 = 0x14
	CLASS_I3C_DEVICE            uint8 = 0x3c
	CLASS_DIAGNOSTIC            uint8 = 0xdc
	CLASS_WIRELESS_CONTROLLER   uint8 = 0xe0
	CLASS_MISCELLANEOUS         uint8 = 0xef
	CLASS_APPLICATION_SPECIFIC  uint8 = 0xfe
	CLASS_VENDOR_SPECIFIC       uint8 = 0xff
)

const (
	ENDPOINT_CONTROL     uint32 = 0
	ENDPOINT_DEV_TO_HOST uint32 = 1
	ENDPOINT_HOST_TO_DEV uint32 = 2
)

const (
	SUBCLASS_NONE uint8 = 0x00

	SUBCLASS_HID_BOOT_INTERFACE uint8 = 0x01
)

const (
	PROTOCOL_NONE         uint8 = 0x00
	PROTOCOL_HID_KEYBOARD uint8 = 0x01
	PROTOCOL_HID_MOUSE    uint8 = 0x02
)

type SetupRequest uint8

const (
	REQUEST_GET_STATUS        SetupRequest = 0
	REQUEST_CLEAR_FEATURE     SetupRequest = 1
	REQUEST_SET_FEATURE       SetupRequest = 3
	REQUEST_SET_ADDRESS       SetupRequest = 5
	REQUEST_GET_DESCRIPTOR    SetupRequest = 6
	REQUEST_SET_DESCRIPTOR    SetupRequest = 7
	REQUEST_GET_CONFIGURATION SetupRequest = 8
	REQUEST_SET_CONFIGURATION SetupRequest = 9
	REQUEST_GET_INTERFACE     SetupRequest = 10
	REQUEST_SET_INTERFACE     SetupRequest = 11
	REQUEST_SYNCH_FRAME       SetupRequest = 12
	REQUEST_SET_SEL           SetupRequest = 48
	REQUEST_SET_ISOCH_DELAY   SetupRequest = 49
)

const (
	REQUEST_HID_GET_REPORT   SetupRequest = 0x01
	REQUEST_HID_GET_IDLE     SetupRequest = 0x02
	REQUEST_HID_GET_PROTOCOL SetupRequest = 0x03
	REQUEST_HID_SET_REPORT   SetupRequest = 0x09
	REQUEST_HID_SET_IDLE     SetupRequest = 0x0A
	REQUEST_HID_SET_PROTOCOL SetupRequest = 0x0B
)

const (
	SETUP_PACKET_LENGTH = 8
)

type SetupDataDirection byte

const (
	SETUP_DATA_DIRECTION_IN  SetupDataDirection = 0
	SETUP_DATA_DIRECTION_OUT SetupDataDirection = 1
)

type SetupDataType byte

const (
	SETUP_DATA_TYPE_STANDARD SetupDataType = 0
	SETUP_DATA_TYPE_CLASS    SetupDataType = 1
	SETUP_DATA_TYPE_VENDOR   SetupDataType = 2
)

type SetupRecipient byte

const (
	SETUP_RECIPIENT_DEVICE    SetupRecipient = 0
	SETUP_RECIPIENT_INTERFACE SetupRecipient = 1
	SETUP_RECIPIENT_ENDPOINT  SetupRecipient = 2
	SETUP_RECIPIENT_OTHER     SetupRecipient = 3
)

type SetupRequestType uint8

func (s SetupRequestType) Direction() SetupDataDirection {
	return SetupDataDirection((s >> 7) & 1)
}

func (s SetupRequestType) Type() SetupDataType {
	return SetupDataType((s >> 5) & 0b11)
}

func (s SetupRequestType) Recipient() SetupRecipient {
	return SetupRecipient(s & 0b11111)
}

func (s *SetupRequestType) SetDirection(direction SetupDataDirection) {
	*s = SetupRequestType(uint8(*s) & ^(uint8(1)<<7) | (uint8(direction) << 7))
}

func (s *SetupRequestType) SetType(dataType SetupDataType) {
	*s = SetupRequestType(uint8(*s) & ^(uint8(0b11)<<5) | (uint8(dataType) << 5))
}

func (s *SetupRequestType) SetRecipient(recipient SetupRecipient) {
	*s = SetupRequestType(uint8(*s) & ^(uint8(0b11111)) | uint8(recipient))
}
