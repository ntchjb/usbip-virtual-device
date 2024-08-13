package usage

const (
	USAGE_PAGE_KIND_DEFINED   = "Defined"
	USAGE_PAGE_KIND_GENERATED = "Generated"
)

type OUsage struct {
	ID    UsageID  `json:"Id"`
	Name  string   `json:"Name"`
	Kinds []string `json:"Kinds"`
}

type OUsageIDGenerator struct {
	NamePrefix   string   `json:"NamePrefix"`
	StartUsageID uint16   `json:"StartUsageId"`
	EndUsageID   uint16   `json:"EndUsageId"`
	Kinds        []string `json:"Kinds"`
}

type OUsagePage struct {
	Kind             string            `json:"Kind"`
	ID               UsagePageID       `json:"Id"`
	Name             string            `json:"Name"`
	UsageIDs         []OUsage          `json:"UsageIds"`
	UsageIDGenerator OUsageIDGenerator `json:"UsageIdGenerator"`
}

type OUsageTable struct {
	UsageTableVersion             uint         `json:"UsageTableVersion"`
	UsageTableRevision            uint         `json:"UsageTableRevision"`
	UsageTableSubRevisionInternal uint         `json:"UsageTableSubRevisionInternal"`
	LastGenerated                 string       `json:"LastGenerated"`
	UsagePages                    []OUsagePage `json:"UsagePages"`
}

func (u *OUsageTable) Indexed() UsagePageTable {
	res := make(UsagePageTable)

	for _, usagePage := range u.UsagePages {
		usageTable := UsageTable{
			PageName: usagePage.Name,
		}
		if usagePage.Kind == USAGE_PAGE_KIND_DEFINED {
			usageTable.Kind = USAGE_TABLE_KIND_DEFINED
			usageTable.Defined = make(map[UsageID]Usage)
			for _, usageDetail := range usagePage.UsageIDs {
				usageTable.Defined[usageDetail.ID] = Usage(usageDetail)
			}
		} else {
			usageTable.Kind = USAGE_TABLE_KIND_GENERATED
			usageTable.GeneratorData = UsageGeneratorData{
				Prefix:  usagePage.UsageIDGenerator.NamePrefix,
				StartID: UsageID(usagePage.UsageIDGenerator.StartUsageID),
				EndID:   UsageID(usagePage.UsageIDGenerator.EndUsageID),
			}
		}

		res[usagePage.ID] = usageTable
	}

	return res
}
