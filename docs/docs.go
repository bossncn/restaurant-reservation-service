// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/bossncn/restaurant-reservation-service/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/public/table/init": {
            "post": {
                "description": "Initializes the total number of tables in the restaurant. This endpoint must be called first and only once.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Table"
                ],
                "summary": "Initialize tables in the restaurant",
                "parameters": [
                    {
                        "description": "Initialize Number of Table Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.InitializeTableRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Total Initialized Tables",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.InitializeTableResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Table Already Initialized",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    }
                }
            }
        },
        "/secure/reservations": {
            "post": {
                "description": "Reserves tables for a group of customers.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reservation"
                ],
                "summary": "Reserve tables",
                "parameters": [
                    {
                        "description": "Number of customers in the group.",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ReservationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Tables reserved successfully.",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.ReservationResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Reservation error.",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    }
                }
            }
        },
        "/secure/reservations/{id}": {
            "delete": {
                "description": "Cancels a reservation and releases the reserved tables.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Reservation"
                ],
                "summary": "Cancel a reservation",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The reservation ID to cancel.",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Reservation canceled successfully.",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.CancelReservationResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Cancellation error.",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.CancelReservationResponse": {
            "type": "object",
            "properties": {
                "freed_tables": {
                    "type": "integer"
                },
                "remaining_tables": {
                    "type": "integer"
                }
            }
        },
        "dto.InitializeTableRequest": {
            "type": "object",
            "properties": {
                "num_tables": {
                    "type": "integer"
                }
            }
        },
        "dto.InitializeTableResponse": {
            "type": "object",
            "properties": {
                "total_tables": {
                    "type": "integer"
                }
            }
        },
        "dto.ReservationRequest": {
            "type": "object",
            "properties": {
                "num_customers": {
                    "type": "integer"
                }
            }
        },
        "dto.ReservationResponse": {
            "type": "object",
            "properties": {
                "booking_id": {
                    "type": "string"
                },
                "remaining_tables": {
                    "type": "integer"
                },
                "tables_reserved": {
                    "type": "integer"
                }
            }
        },
        "model.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "data": {},
                "memory": {
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
	Title:            "Restaurant Reservation Service",
	Description:      "Service for managing table reservations in a restaurant.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
