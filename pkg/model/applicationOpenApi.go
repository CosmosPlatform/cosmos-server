package model

import "github.com/getkin/kin-openapi/openapi3"

type ApplicationOpenAPISpecification struct {
	Application *Application
	OpenAPISpec *openapi3.T
}
