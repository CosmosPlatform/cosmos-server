package model

type SentinelSettings struct {
	Interval int
	Enabled  bool
}

type SentinelSettingsUpdate struct {
	Interval *int
	Enabled  *bool
}
