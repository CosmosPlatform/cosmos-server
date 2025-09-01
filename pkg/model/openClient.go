package model

import (
	"fmt"
	"regexp"
	"strings"
)

type OpenClientSpecification struct {
	Dependencies map[string]DependencySpecification `json:"dependencies"`
}

type DependencySpecification struct {
	Reasons   []string                   `json:"reasons"`
	Endpoints map[string]EndpointMethods `json:"endpoints"`
}

type EndpointMethods map[string]EndpointSpecification

type EndpointSpecification struct {
	Reasons []string `json:"reasons"`
}

var (
	validHTTPMethods = map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true,
		"PATCH": true, "HEAD": true, "OPTIONS": true, "TRACE": true,
	}

	// Matches paths like /users, /users/{id}, /api/v1/users/{userId}/orders/{orderId}
	validPathRegex = regexp.MustCompile(`^/[a-zA-Z0-9\-_.~!*'();:@&=+$,/?#\[\]{}|%]*$`)
)

func (spec *OpenClientSpecification) Validate() error {
	for depName, dep := range spec.Dependencies {
		if depName == "" {
			return fmt.Errorf("dependency name cannot be empty")
		}

		for path, methods := range dep.Endpoints {
			if path == "" {
				return fmt.Errorf("endpoint path cannot be empty for dependency %s", depName)
			}

			if !isValidPath(path) {
				return fmt.Errorf("invalid path '%s' for dependency %s: path must start with '/' and contain valid URL characters", path, depName)
			}

			for method := range methods {
				if method == "" {
					return fmt.Errorf("endpoint method cannot be empty for dependency %s", depName)
				}

				if !isValidHTTPMethod(method) {
					return fmt.Errorf("invalid HTTP method '%s' for dependency %s: must be one of GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE", method, depName)
				}
			}
		}
	}
	return nil
}

func isValidPath(path string) bool {
	if !validPathRegex.MatchString(path) {
		return false
	}

	segments := strings.Split(path, "/")

	// Regex to match either plain text or {parameter}
	segmentRegex := regexp.MustCompile(`^([a-zA-Z0-9\-_.~!*'();:@&=+$,?#|%]+|\{[a-zA-Z0-9\-_.~!*'();:@&=+$,?#|%]+\})$`)

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		if !segmentRegex.MatchString(segment) {
			return false
		}
	}

	return true
}

func isValidHTTPMethod(method string) bool {
	return validHTTPMethods[strings.ToUpper(method)]
}
