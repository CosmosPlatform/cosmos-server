package api

type ErrorResponse struct {
	StatusCode int      `json:"-"`
	Error      string   `json:"error"`
	Details    []string `json:"details,omitempty"`
}
