package obj

type SentinelSetting struct {
	CosmosObj
	Name     string `gorm:"uniqueIndex"`
	Interval int
	Enabled  bool
}
