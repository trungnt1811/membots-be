// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/memeception": {
            "get": {
                "description": "Get memeception by symbol",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memeception"
                ],
                "summary": "Get memeception by symbol",
                "parameters": [
                    {
                        "type": "string",
                        "description": "symbol to query, default is ",
                        "name": "symbol",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.MemeceptionResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    }
                }
            }
        },
        "/api/v1/memeceptions": {
            "get": {
                "description": "Get memeceptions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memeception"
                ],
                "summary": "Get memeceptions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.MemeceptionResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.MemeCommon": {
            "type": "object",
            "properties": {
                "bannerUrl": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "logoUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                }
            }
        },
        "dto.MemeceptionCommon": {
            "type": "object",
            "properties": {
                "ama": {
                    "type": "boolean"
                },
                "contractAddress": {
                    "type": "string"
                },
                "meme": {
                    "$ref": "#/definitions/dto.MemeCommon"
                },
                "startAt": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "dto.MemeceptionResp": {
            "type": "object",
            "properties": {
                "live": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.MemeceptionCommon"
                    }
                },
                "past": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.MemeceptionCommon"
                    }
                },
                "upcoming": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.MemeceptionCommon"
                    }
                }
            }
        },
        "util.GeneralError": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "HTTP error code, or custom error code",
                    "type": "integer"
                },
                "errors": {
                    "description": "List of error send server 2 server",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "string": {
                    "description": "Friendly error message",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "membots-be API",
	Description:      "This Swagger docs for membots-be.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
