package usage

import (
	"encoding/json"
	"strconv"
)

type UsagePageID uint16

type UsageID uint16

const (
	USAGE_TABLE_KIND_DEFINED   = 0
	USAGE_TABLE_KIND_GENERATED = 1

	USAGE_RESERVED_NAME = "Reserved"
	USAGE_VENDOR_NAME   = "Vendor"
)

type Usage struct {
	ID    UsageID
	Name  string
	Kinds []string
}

type UsageGeneratorData struct {
	Prefix  string
	StartID UsageID
	EndID   UsageID
}

type UsageTable struct {
	Kind          uint8
	PageName      string
	Defined       map[UsageID]Usage
	GeneratorData UsageGeneratorData
}

func (u *UsageTable) generate(id UsageID) string {
	if id >= u.GeneratorData.StartID && id <= u.GeneratorData.EndID {
		return u.GeneratorData.Prefix + "-" + strconv.FormatUint(uint64(id), 10)
	}
	return USAGE_RESERVED_NAME
}

func (u *UsageTable) GetUsageName(id UsageID) string {
	if u.Kind == USAGE_TABLE_KIND_DEFINED {
		if usage, ok := u.Defined[id]; ok {
			return usage.Name
		} else {
			return USAGE_RESERVED_NAME
		}
	} else {
		return u.generate(id)
	}
}

type UsagePageTable map[UsagePageID]UsageTable

func (u UsagePageTable) GetUsageName(pageID UsagePageID, id UsageID) string {
	if indexed, ok := u[pageID]; !ok {
		if pageID >= 0xFF00 && pageID <= 0xFFFF {
			return USAGE_VENDOR_NAME
		}
		return USAGE_RESERVED_NAME
	} else {
		return indexed.GetUsageName(id)
	}
}

func (u UsagePageTable) GetUsagePageName(pageID UsagePageID) string {
	if indexed, ok := u[pageID]; !ok {
		if pageID >= 0xFF00 && pageID <= 0xFFFF {
			return USAGE_VENDOR_NAME
		}
		return USAGE_RESERVED_NAME
	} else {
		return indexed.PageName
	}
}

var IndexedUsageTable UsagePageTable

func CreateIndexedUsageTable() UsagePageTable {
	var usageTable OUsageTable

	if err := json.Unmarshal([]byte(USAGE_SPEC_JSON), &usageTable); err != nil {
		panic(err)
	}

	return usageTable.Indexed()
}

func init() {
	IndexedUsageTable = CreateIndexedUsageTable()
}
