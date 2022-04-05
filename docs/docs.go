// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/auth/access-token-verify": {
            "post": {
                "description": "jwt access token verify",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "access token verify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "accessToken",
                        "name": "accessToken",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/forgot-password": {
            "post": {
                "description": "forgot password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "forgot password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password that have at least 8 length and contain an alphabet and number ",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "repeatPassword that have at least 8 length and contain an alphabet and number ",
                        "name": "repeatPassword",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "firstName",
                        "name": "firstName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "lastName",
                        "name": "lastName",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "jwt login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "send user operating system + browser name in this param",
                        "name": "deviceName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.LoginResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedLoginResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "jwt logout , atuhentication required",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "logout",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.UnauthenticatedResponse"
                        }
                    }
                }
            }
        },
        "/auth/recover-password": {
            "post": {
                "description": "Let user change it password with forgot token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "recover-password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password that have at least 8 length and contain an alphabet and number ",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "repeatPassword that have at least 8 length and contain an alphabet and number ",
                        "name": "repeatPassword",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.NotFoundResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "jwt register",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "register",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password that have at least 8 length and contain an alphabet and number ",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "repeatPassword that have at least 8 length and contain an alphabet and number ",
                        "name": "repeatPassword",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "firstName",
                        "name": "firstName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "lastName",
                        "name": "lastName",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/renew-access-token": {
            "post": {
                "description": "jwt renew access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "renew access token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "refreshToken",
                        "name": "refreshToken",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessVerifyAccessTokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/resend-verify-email": {
            "post": {
                "description": "resend-verify-email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "resend-verify-email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique email",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "post": {
                "description": "verify-email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "verify-email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "generic"
                ],
                "summary": "ping",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.PingResponse"
                        }
                    }
                }
            }
        },
        "/profile/change-password": {
            "post": {
                "description": "Change Password , authentication required",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "change-password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "password that have at least 8 length and contain an alphabet and number ",
                        "name": "password",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "repeatPassword that have at least 8 length and contain an alphabet and number ",
                        "name": "repeatPassword",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.UnauthenticatedResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        },
        "/profile/devices": {
            "post": {
                "description": "return logged in devices in user's account , authentication required",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "devices",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.DevicesResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.UnauthenticatedResponse"
                        }
                    }
                }
            }
        },
        "/profile/terminate-device": {
            "post": {
                "description": "jwt terminate-device , atuhentication required",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "terminate-device",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/swagger.SuccessResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/swagger.UnauthenticatedResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/swagger.NotFoundResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/swagger.FailedValidationResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AccessTokenRes": {
            "type": "object",
            "required": [
                "accessToken",
                "expAccessToken"
            ],
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expAccessToken": {
                    "type": "string"
                }
            }
        },
        "models.Device": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "deviceName": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.DeviceListResponse": {
            "type": "object",
            "properties": {
                "current": {
                    "type": "string"
                },
                "devices": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Device"
                    }
                }
            }
        },
        "models.LoginResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expAccessToken": {
                    "type": "string"
                },
                "expRefreshToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/models.UserResponse"
                }
            }
        },
        "models.UserResponse": {
            "type": "object",
            "required": [
                "firstName",
                "lastName"
            ],
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "verifiedEmail": {
                    "type": "boolean"
                }
            }
        },
        "swagger.DevicesResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.DeviceListResponse"
                },
                "msg": {
                    "type": "string",
                    "example": ""
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "swagger.EmptyData": {
            "type": "object"
        },
        "swagger.FailedLoginResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.EmptyData"
                },
                "msg": {
                    "type": "string",
                    "example": "No user found with entered credentials"
                },
                "ok": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "swagger.FailedResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.EmptyData"
                },
                "msg": {
                    "type": "string",
                    "example": "Error or warnnig message"
                },
                "ok": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "swagger.FailedValidationResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.validationErrors"
                },
                "msg": {
                    "type": "string",
                    "example": "Please review your entered data"
                },
                "ok": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "swagger.LoginResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.LoginResponse"
                },
                "msg": {
                    "type": "string",
                    "example": "Successful message"
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "swagger.NotFoundResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.EmptyData"
                },
                "msg": {
                    "type": "string",
                    "example": "404 not found!"
                },
                "ok": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "swagger.PingResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "pingpong": "🏓🏓🏓🏓🏓🏓"
                    }
                },
                "msg": {
                    "type": "string",
                    "example": "pong"
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "swagger.SuccessResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.EmptyData"
                },
                "msg": {
                    "type": "string",
                    "example": "Successful message"
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "swagger.SuccessVerifyAccessTokenResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.AccessTokenRes"
                },
                "msg": {
                    "type": "string",
                    "example": "Successful message"
                },
                "ok": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "swagger.UnauthenticatedResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/swagger.EmptyData"
                },
                "msg": {
                    "type": "string",
                    "example": "You must login first!"
                },
                "ok": {
                    "type": "boolean",
                    "example": false
                }
            }
        },
        "swagger.validationErrors": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    },
                    "example": {
                        "field1": "This field is required",
                        "field2": "This field must be numeric"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
