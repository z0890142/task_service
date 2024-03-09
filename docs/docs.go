// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/task-service/api/v1/tasks": {
            "get": {
                "summary": "list tasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "order",
                        "name": "order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "post": {
                "summary": "create task",
                "parameters": [
                    {
                        "description": "task",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        },
        "/task-service/api/v1/tasks/{taskId}": {
            "get": {
                "summary": "get tasks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "put": {
                "summary": "update task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "task",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "summary": "delete task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "task ID",
                        "name": "taskId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "code.Code": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3,
                4,
                5,
                6,
                7,
                16,
                8,
                9,
                10,
                11,
                12,
                13,
                14,
                15
            ],
            "x-enum-varnames": [
                "Code_OK",
                "Code_CANCELLED",
                "Code_UNKNOWN",
                "Code_INVALID_ARGUMENT",
                "Code_DEADLINE_EXCEEDED",
                "Code_NOT_FOUND",
                "Code_ALREADY_EXISTS",
                "Code_PERMISSION_DENIED",
                "Code_UNAUTHENTICATED",
                "Code_RESOURCE_EXHAUSTED",
                "Code_FAILED_PRECONDITION",
                "Code_ABORTED",
                "Code_OUT_OF_RANGE",
                "Code_UNIMPLEMENTED",
                "Code_INTERNAL",
                "Code_UNAVAILABLE",
                "Code_DATA_LOSS"
            ]
        },
        "models.HttpError": {
            "type": "object",
            "required": [
                "code",
                "message"
            ],
            "properties": {
                "code": {
                    "$ref": "#/definitions/code.Code"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "$ref": "#/definitions/code.Code"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Task"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.Task": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "tag": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "version": {
                    "type": "integer"
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
	Title:            "Task Service",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}