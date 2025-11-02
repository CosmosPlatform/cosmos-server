package monitoring

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oasdiff/oasdiff/checker"
	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/load"
	"gopkg.in/yaml.v3"
)

type OpenApiService interface {
	ParseOpenApiSpec(specContent string) (*openapi3.T, error)
	CompareOpenApiSpecs(spec1, spec2 *openapi3.T) (checker.Changes, error)
}

type openApiService struct{}

func NewOpenApiService() OpenApiService {
	return &openApiService{}
}

func (s *openApiService) ParseOpenApiSpec(specContent string) (*openapi3.T, error) {
	version, err := s.detectOpenAPIVersion(specContent)
	if err != nil {
		return nil, fmt.Errorf("failed to detect OpenAPI version: %s", err.Error())
	}

	var doc *openapi3.T

	if version == 3 {
		loader := openapi3.NewLoader()
		doc, err = loader.LoadFromData([]byte(specContent))
		if err != nil {
			return nil, fmt.Errorf("failed to load OpenAPI 3 spec: %s", err.Error())
		}
	} else if version == 2 {
		var swagger2 openapi2.T
		if err := swagger2.UnmarshalJSON([]byte(specContent)); err != nil {
			// If not JSON, we try YAML
			if err := yaml.Unmarshal([]byte(specContent), &swagger2); err != nil {
				return nil, fmt.Errorf("spec is neither valid JSON nor YAML: %s", err.Error())
			}
		}

		doc, err = openapi2conv.ToV3(&swagger2)
		if err != nil {
			return nil, fmt.Errorf("failed to convert OpenAPI 2 spec to 3: %s", err.Error())
		}
	} else {
		return nil, fmt.Errorf("unsupported OpenAPI version: %d", version)
	}

	return doc, nil
}

func (s *openApiService) detectOpenAPIVersion(spec string) (int, error) {
	// We try JSON first
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(spec), &obj); err != nil {
		// If not JSON, we try YAML
		if err := yaml.Unmarshal([]byte(spec), &obj); err != nil {
			return 0, fmt.Errorf("spec is neither valid JSON nor YAML: %w", err)
		}
	}

	if v, ok := obj["swagger"]; ok {
		if s, ok := v.(string); ok && strings.HasPrefix(s, "2.0") {
			return 2, nil
		}
	}
	if v, ok := obj["openapi"]; ok {
		if s, ok := v.(string); ok && strings.HasPrefix(s, "3") {
			return 3, nil
		}
	}
	return 0, fmt.Errorf("could not detect OpenAPI version")
}

func (s *openApiService) CompareOpenApiSpecs(spec1, spec2 *openapi3.T) (checker.Changes, error) {
	loadSpec1 := &load.SpecInfo{Spec: spec1}
	loadSpec2 := &load.SpecInfo{Spec: spec2}

	diffReport, operationsSources, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), loadSpec1, loadSpec2)
	if err != nil {
		return nil, fmt.Errorf("failed to compare OpenAPI specs: %s", err.Error())
	}

	checkerConfig := checker.NewConfig(checker.GetAllChecks())
	changes := checker.CheckBackwardCompatibilityUntilLevel(checkerConfig, diffReport, operationsSources, checker.INFO)

	return changes, nil
}
