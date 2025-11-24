# OpenClient Specification

As part of the development of this project, we have created the OpenClient specification, which defines a standardized file to describe the dependencies of your applications.

The OpenClient specification is designed to be simple and easy to use, allowing developers to quickly define the dependencies of their applications in a clear and concise manner.

## Overview

The OpenClient file is a JSON document that describes all external service dependencies for an application. It provides:

- **Transparency**: Clear documentation of which services your application depends on
- **Traceability**: Explicit reasons for each dependency at both service and endpoint levels
- **Maintainability**: Easy-to-understand structure that evolves with your application
- **Integration**: Can be used by monitoring and documentation tools to visualize service interactions

## File Structure

The OpenClient file must be a valid JSON document named `openclient.json` placed in the root directory of your application.

### Root Object

The root object contains a single required property:

- `dependencies` (object, required): A map of service dependencies where each key is the service name

### Dependency Object

Each dependency object describes a single external service and contains:

- `reasons` (array of strings, required): High-level reasons why this service is needed
- `endpoints` (object, required): A map of endpoint paths that are consumed from this service

### Endpoint Object

Each endpoint is keyed by its path (e.g., `/users/{id}`) and contains:

- HTTP method objects (e.g., `get`, `post`, `put`, `delete`, `patch`)

### HTTP Method Object

Each HTTP method object contains:

- `reasons` (array of strings, required): Specific reasons why this endpoint method is called

## Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["dependencies"],
  "properties": {
    "dependencies": {
      "type": "object",
      "patternProperties": {
        "^[a-zA-Z0-9-_]+$": {
          "type": "object",
          "required": ["reasons", "endpoints"],
          "properties": {
            "reasons": {
              "type": "array",
              "items": {
                "type": "string"
              },
              "minItems": 1
            },
            "endpoints": {
              "type": "object",
              "patternProperties": {
                "^/.*$": {
                  "type": "object",
                  "properties": {
                    "get": {
                      "$ref": "#/definitions/methodObject"
                    },
                    "post": {
                      "$ref": "#/definitions/methodObject"
                    },
                    "put": {
                      "$ref": "#/definitions/methodObject"
                    },
                    "delete": {
                      "$ref": "#/definitions/methodObject"
                    },
                    "patch": {
                      "$ref": "#/definitions/methodObject"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "methodObject": {
      "type": "object",
      "required": ["reasons"],
      "properties": {
        "reasons": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "minItems": 1
        }
      }
    }
  }
}
```

## Examples

### Example 1: Notification Service

A notification service that depends on user and order services:

```json
{
  "dependencies": {
    "user-service": {
      "reasons": [
        "Get user contact information for sending notifications"
      ],
      "endpoints": {
        "/users/{id}": {
          "get": {
            "reasons": [
              "Fetch email or phone number"
            ]
          }
        }
      }
    },
    "order-service": {
      "reasons": [
        "Send order-related notifications"
      ],
      "endpoints": {
        "/orders/{id}": {
          "get": {
            "reasons": [
              "Retrieve order details for notification content"
            ]
          }
        }
      }
    }
  }
}
```

### Example 2: Order Service

An order service with multiple dependencies for product information and inventory management:

```json
{
  "dependencies": {
    "user-service": {
      "reasons": [
        "Get user information for shipping and billing"
      ],
      "endpoints": {
        "/users/{id}": {
          "get": {
            "reasons": [
              "Retrieve user details for order"
            ]
          }
        }
      }
    },
    "product-service": {
      "reasons": [
        "Fetch product data for ordered items"
      ],
      "endpoints": {
        "/products/{id}": {
          "get": {
            "reasons": [
              "Retrieve product name and price"
            ]
          }
        }
      }
    },
    "inventory-service": {
      "reasons": [
        "Reserve and release stock for orders"
      ],
      "endpoints": {
        "/inventory/reserve": {
          "post": {
            "reasons": [
              "Reserve stock when order is placed"
            ]
          }
        },
        "/inventory/release": {
          "post": {
            "reasons": [
              "Release stock when order is canceled"
            ]
          }
        }
      }
    }
  }
}
```

### Example 3: Review Service

A review service that validates users and products:

```json
{
  "dependencies": {
    "user-service": {
      "reasons": [
        "Verify reviewer identity"
      ],
      "endpoints": {
        "/users/{id}": {
          "get": {
            "reasons": [
              "Ensure user exists before allowing review"
            ]
          }
        }
      }
    },
    "product-service": {
      "reasons": [
        "Validate product being reviewed"
      ],
      "endpoints": {
        "/products/{id}": {
          "get": {
            "reasons": [
              "Verify product exists before accepting review"
            ]
          }
        }
      }
    }
  }
}
```

## Best Practices

### 1. Service Naming

- Use lowercase with hyphens for service names (e.g., `user-service`, `order-service`)
- Keep names consistent with your actual service names
- Use descriptive names that clearly identify the service's purpose

### 2. Reasons Documentation

- Be specific and action-oriented in your reasons
- Service-level reasons should describe the overall purpose
- Endpoint-level reasons should describe the specific use case
- Avoid generic reasons like "for integration" - instead explain what data or functionality is needed

### 3. Endpoint Paths

- Use the exact path as defined in the service's API
- Include path parameters using curly braces (e.g., `{id}`, `{userId}`)
- Keep paths consistent with the service's OpenAPI/Swagger documentation

### 4. HTTP Methods

- Only include methods that your application actually uses
- Don't document theoretical or future usage - only current dependencies
- If you use the same endpoint with multiple methods, document each separately

### 5. Maintenance

- Update the OpenClient file whenever you add or remove a dependency
- Review and update reasons if the purpose of a dependency changes
- Remove entries for endpoints or services that are no longer used
- Keep the file in version control alongside your code

## Validation

To ensure your OpenClient file is valid:

1. Verify it's valid JSON syntax
2. Ensure all required fields are present (`dependencies`, `reasons`, `endpoints`)
3. Check that all reasons arrays contain at least one item
4. Confirm endpoint paths start with a forward slash
5. Validate HTTP method names are lowercase
