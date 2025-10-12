package api

type GetApplicationOpenAPISpecificationResponse struct {
	ApplicationName string `json:"applicationName"`
	OpenAPISpec     string `json:"openAPISpec"`
}
