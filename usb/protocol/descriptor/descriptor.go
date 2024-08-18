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

type LangID uint16

const (
	LANGID_AFRIKAANS                  LangID = 0x0436
	LANGID_ALBANIAN                   LangID = 0x041c
	LANGID_ARABIC_SAUDI_ARABIA        LangID = 0x0401
	LANGID_ARABIC_IRAQ                LangID = 0x0801
	LANGID_ARABIC_EGYPT               LangID = 0x0c01
	LANGID_ARABIC_LIBYA               LangID = 0x1001
	LANGID_ARABIC_ALGERIA             LangID = 0x1401
	LANGID_ARABIC_MOROCCO             LangID = 0x1801
	LANGID_ARABIC_TUNISIA             LangID = 0x1c01
	LANGID_ARABIC_OMAN                LangID = 0x2001
	LANGID_ARABIC_YEMEN               LangID = 0x2401
	LANGID_ARABIC_SYRIA               LangID = 0x2801
	LANGID_ARABIC_JORDAN              LangID = 0x2c01
	LANGID_ARABIC_LEBANON             LangID = 0x3001
	LANGID_ARABIC_KUWAIT              LangID = 0x3401
	LANGID_ARABIC_UAE                 LangID = 0x3801
	LANGID_ARABIC_BAHRAIN             LangID = 0x3c01
	LANGID_ARABIC_QATAR               LangID = 0x4001
	LANGID_ARMENIAN                   LangID = 0x042b
	LANGID_ASSAMESE                   LangID = 0x044d
	LANGID_AZERI_LATIN                LangID = 0x042c
	LANGID_AZERI_CYRILLIC             LangID = 0x082c
	LANGID_BASQUE                     LangID = 0x042d
	LANGID_BELARUSSIAN                LangID = 0x0423
	LANGID_BENGALI                    LangID = 0x0445
	LANGID_BULGARIAN                  LangID = 0x0402
	LANGID_BURMESE                    LangID = 0x0455
	LANGID_CATALAN                    LangID = 0x0403
	LANGID_CHINESE_TAIWAN             LangID = 0x0404
	LANGID_CHINESE_PRC                LangID = 0x0804
	LANGID_CHINESE_HONG_KONG_SAR_PRC  LangID = 0x0c04
	LANGID_CHINESE_SINGAPORE          LangID = 0x1004
	LANGID_CHINESE_MACAU_SAR          LangID = 0x1404
	LANGID_CROATIAN                   LangID = 0x041a
	LANGID_CZECH                      LangID = 0x0405
	LANGID_DANISH                     LangID = 0x0406
	LANGID_DUTCH_NETHERLANDS          LangID = 0x0413
	LANGID_DUTCH_BELGIUM              LangID = 0x0813
	LANGID_ENGLISH_UNITED_STATES      LangID = 0x0409
	LANGID_ENGLISH_UNITED_KINGDOM     LangID = 0x0809
	LANGID_ENGLISH_AUSTRALIAN         LangID = 0x0c09
	LANGID_ENGLISH_CANADIAN           LangID = 0x1009
	LANGID_ENGLISH_NEW_ZEALAND        LangID = 0x1409
	LANGID_ENGLISH_IRELAND            LangID = 0x1809
	LANGID_ENGLISH_SOUTH_AFRICA       LangID = 0x1c09
	LANGID_ENGLISH_JAMAICA            LangID = 0x2009
	LANGID_ENGLISH_CARIBBEAN          LangID = 0x2409
	LANGID_ENGLISH_BELIZE             LangID = 0x2809
	LANGID_ENGLISH_TRINIDAD           LangID = 0x2c09
	LANGID_ENGLISH_ZIMBABWE           LangID = 0x3009
	LANGID_ENGLISH_PHILIPPINES        LangID = 0x3409
	LANGID_ESTONIAN                   LangID = 0x0425
	LANGID_FAEROESE                   LangID = 0x0438
	LANGID_FARSI                      LangID = 0x0429
	LANGID_FINNISH                    LangID = 0x040b
	LANGID_FRENCH_STANDARD            LangID = 0x040c
	LANGID_FRENCH_BELGIAN             LangID = 0x080c
	LANGID_FRENCH_CANADIAN            LangID = 0x0c0c
	LANGID_FRENCH_SWITZERLAND         LangID = 0x100c
	LANGID_FRENCH_LUXEMBOURG          LangID = 0x140c
	LANGID_FRENCH_MONACO              LangID = 0x180c
	LANGID_GEORGIAN                   LangID = 0x0437
	LANGID_GERMAN_STANDARD            LangID = 0x0407
	LANGID_GERMAN_SWITZERLAND         LangID = 0x0807
	LANGID_GERMAN_AUSTRIA             LangID = 0x0c07
	LANGID_GERMAN_LUXEMBOURG          LangID = 0x1007
	LANGID_GERMAN_LIECHTENSTEIN       LangID = 0x1407
	LANGID_GREEK                      LangID = 0x0408
	LANGID_GUJARATI                   LangID = 0x0447
	LANGID_HEBREW                     LangID = 0x040d
	LANGID_HINDI                      LangID = 0x0439
	LANGID_HUNGARIAN                  LangID = 0x040e
	LANGID_ICELANDIC                  LangID = 0x040f
	LANGID_INDONESIAN                 LangID = 0x0421
	LANGID_ITALIAN_STANDARD           LangID = 0x0410
	LANGID_ITALIAN_SWITZERLAND        LangID = 0x0810
	LANGID_JAPANESE                   LangID = 0x0411
	LANGID_KANNADA                    LangID = 0x044b
	LANGID_KASHMIRI_INDIA             LangID = 0x0860
	LANGID_KAZAKH                     LangID = 0x043f
	LANGID_KONKANI                    LangID = 0x0457
	LANGID_KOREAN                     LangID = 0x0412
	LANGID_KOREAN_JOHAB               LangID = 0x0812
	LANGID_LATVIAN                    LangID = 0x0426
	LANGID_LITHUANIAN                 LangID = 0x0427
	LANGID_LITHUANIAN_CLASSIC         LangID = 0x0827
	LANGID_MACEDONIAN                 LangID = 0x042f
	LANGID_MALAY_MALAYSIAN            LangID = 0x043e
	LANGID_MALAY_BRUNEI_DARUSSALAM    LangID = 0x083e
	LANGID_MALAYALAM                  LangID = 0x044c
	LANGID_MANIPURI                   LangID = 0x0458
	LANGID_MARATHI                    LangID = 0x044e
	LANGID_NEPALI_INDIA               LangID = 0x0861
	LANGID_NORWEGIAN_BOKMAL           LangID = 0x0414
	LANGID_NORWEGIAN_NYNORSK          LangID = 0x0814
	LANGID_ORIYA                      LangID = 0x0448
	LANGID_POLISH                     LangID = 0x0415
	LANGID_PORTUGUESE_BRAZIL          LangID = 0x0416
	LANGID_PORTUGUESE_STANDARD        LangID = 0x0816
	LANGID_PUNJABI                    LangID = 0x0446
	LANGID_ROMANIAN                   LangID = 0x0418
	LANGID_RUSSIAN                    LangID = 0x0419
	LANGID_SANSKRIT                   LangID = 0x044f
	LANGID_SERBIAN_CYRILLIC           LangID = 0x0c1a
	LANGID_SERBIAN_LATIN              LangID = 0x081a
	LANGID_SINDHI                     LangID = 0x0459
	LANGID_SLOVAK                     LangID = 0x041b
	LANGID_SLOVENIAN                  LangID = 0x0424
	LANGID_SPANISH_TRADITIONAL_SORT   LangID = 0x040a
	LANGID_SPANISH_MEXICAN            LangID = 0x080a
	LANGID_SPANISH_MODERN_SORT        LangID = 0x0c0a
	LANGID_SPANISH_GUATEMALA          LangID = 0x100a
	LANGID_SPANISH_COSTA_RICA         LangID = 0x140a
	LANGID_SPANISH_PANAMA             LangID = 0x180a
	LANGID_SPANISH_DOMINICAN_REPUBLIc LangID = 0x1c0a
	LANGID_SPANISH_VENEZUELA          LangID = 0x200a
	LANGID_SPANISH_COLOMBIA           LangID = 0x240a
	LANGID_SPANISH_PERU               LangID = 0x280a
	LANGID_SPANISH_ARGENTINA          LangID = 0x2c0a
	LANGID_SPANISH_ECUADOR            LangID = 0x300a
	LANGID_SPANISH_CHILE              LangID = 0x340a
	LANGID_SPANISH_URUGUAY            LangID = 0x380a
	LANGID_SPANISH_PARAGUAY           LangID = 0x3c0a
	LANGID_SPANISH_BOLIVIA            LangID = 0x400a
	LANGID_SPANISH_EL_SALVADOR        LangID = 0x440a
	LANGID_SPANISH_HONDURAS           LangID = 0x480a
	LANGID_SPANISH_NICARAGUA          LangID = 0x4c0a
	LANGID_SPANISH_PUERTO_RICO        LangID = 0x500a
	LANGID_SUTU                       LangID = 0x0430
	LANGID_SWAHILI_KENYA              LangID = 0x0441
	LANGID_SWEDISH                    LangID = 0x041d
	LANGID_SWEDISH_FINLAND            LangID = 0x081d
	LANGID_TAMIL                      LangID = 0x0449
	LANGID_TATAR_TATARSTAN            LangID = 0x0444
	LANGID_TELUGU                     LangID = 0x044a
	LANGID_THAI                       LangID = 0x041e
	LANGID_TURKISH                    LangID = 0x041f
	LANGID_UKRAINIAN                  LangID = 0x0422
	LANGID_URDU_PAKISTAN              LangID = 0x0420
	LANGID_URDU_INDIA                 LangID = 0x0820
	LANGID_UZBEK_LATIN                LangID = 0x0443
	LANGID_UZBEK_CYRILLIC             LangID = 0x0843
	LANGID_VIETNAMESE                 LangID = 0x042a
	LANGID_HID_USAGE_DATA_DESCRIPTOR  LangID = 0x04ff
	LANGID_HID_VENDOR_DEFINED_1       LangID = 0xf0ff
	LANGID_HID_VENDOR_DEFINED_2       LangID = 0xf4ff
	LANGID_HID_VENDOR_DEFINED_3       LangID = 0xf8ff
	LANGID_HID_VENDOR_DEFINED_4       LangID = 0xfcff
)

func GetDescriptorTypeAndIndex(wValue uint16) (descriptorType DescriptorType, index uint8) {
	descriptorType = DescriptorType(wValue >> 8)
	index = uint8(wValue & 0x00FF)

	return
}
