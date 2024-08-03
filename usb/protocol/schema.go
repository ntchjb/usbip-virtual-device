package protocol

type BusID [32]byte

const (
	HIDSpecVersion      uint16 = 0x0110
	HIDClassSpecVersion uint16 = 0x0101
)

const (
	SpeedUSB1Low uint32 = iota
	SpeedUSB1Full
	SpeedUSB2High
	SpeedUSB3Super
)

const (
	ClassBasedOnInterface    uint8 = 0x00
	ClassAudio               uint8 = 0x01
	ClassCDCControl          uint8 = 0x02
	ClassHID                 uint8 = 0x03
	ClassPhysical            uint8 = 0x05
	ClassImage               uint8 = 0x06
	ClassPrinter             uint8 = 0x07
	ClassMassStorage         uint8 = 0x08
	ClassHub                 uint8 = 0x09
	ClassCDCData             uint8 = 0x0a
	ClassSmartCard           uint8 = 0x0b
	ClassContentSecutiry     uint8 = 0x0d
	ClassVideo               uint8 = 0x0e
	ClassPersonalHealthcare  uint8 = 0x0f
	ClassAudioAndVideo       uint8 = 0x10
	ClassBillboardDevice     uint8 = 0x11
	ClassUSBTypeCBridge      uint8 = 0x12
	ClassBulkDisplayProtocol uint8 = 0x13
	ClassMCTPOverUSB         uint8 = 0x14
	ClassI3CDevice           uint8 = 0x3c
	ClassDiagnostic          uint8 = 0xdc
	ClassWirelessController  uint8 = 0xe0
	ClassMiscellaneous       uint8 = 0xef
	ClassApplicationSpecific uint8 = 0xfe
	ClassVendorSpecific      uint8 = 0xff
)

const (
	EndpointControl   uint32 = 0
	EndpointDevToHost uint32 = 1
	EndpointHostToDev uint32 = 2
)

const (
	SubclassNone uint8 = 0x00

	HIDSubclassBootInterface uint8 = 0x01
)

const (
	ProtocolNone        uint8 = 0x00
	HIDProtocolKeyboard uint8 = 0x01
	HIDProtocolMouse    uint8 = 0x02
)

type SetupRequest uint8

const (
	RequestGetStatus        SetupRequest = 0
	RequestClearFeature     SetupRequest = 1
	RequestSetFeature       SetupRequest = 3
	RequestSetAddress       SetupRequest = 5
	RequestGetDescriptor    SetupRequest = 6
	RequestSetDescriptor    SetupRequest = 7
	RequestGetConfiguration SetupRequest = 8
	RequestSetConfiguration SetupRequest = 9
	RequestGetInterface     SetupRequest = 10
	RequestSetInterface     SetupRequest = 11
	RequestSynchFrame       SetupRequest = 12
	RequestSetSel           SetupRequest = 48
	RequestSetISOCHDelay    SetupRequest = 49
)

const (
	RequestHIDGetReport   SetupRequest = 0x01
	RequestHIDGetIdle     SetupRequest = 0x02
	RequestHIDGetProtocol SetupRequest = 0x03
	RequestHIDSetReport   SetupRequest = 0x09
	RequestHIDSetIdle     SetupRequest = 0x0A
	RequestHIDSetProtocol SetupRequest = 0x0B
)

const (
	SETUP_PACKET_LENGTH = 8
)

type SetupDataDirection byte

const (
	SetupDataDirectionIn  SetupDataDirection = 0
	SetupDataDirectionOut SetupDataDirection = 1
)

type SetupDataType byte

const (
	SetupDataTypeStandard SetupDataType = 0
	SetupDataTypeClass    SetupDataType = 1
	SetupDataTypeVendor   SetupDataType = 2
)

type SetupRecipient byte

const (
	SetupRecipientDevice    SetupRecipient = 0
	SetupRecipientInterface SetupRecipient = 1
	SetupRecipientEndpoint  SetupRecipient = 2
	SetupRecipientOther     SetupRecipient = 3
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
