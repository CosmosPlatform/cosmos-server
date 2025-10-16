package api

type UpdateSentinelSettingsRequest struct {
	Interval *int  `json:"interval"`
	Enabled  *bool `json:"enabled"`
}

type GetSentinelSettingsResponse struct {
	Interval int  `json:"interval"`
	Enabled  bool `json:"enabled"`
}
