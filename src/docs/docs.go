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
        "/v0.1/auth/signin": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 400\nUSER_NOT_EXIST : 유저가 존재하지 않음\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "로그인",
                "parameters": [
                    {
                        "description": "이메일, 비밀번호",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqSignin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ResSignin"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/v0.1/auth/signup": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 400\nUSER_ALREADY_EXISTED : 유저가 이미 존재\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "회원 가입",
                "parameters": [
                    {
                        "description": "이름, 이메일, 비밀번호",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqSignup"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/v0.1/room/create": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "room"
                ],
                "summary": "방 생성",
                "parameters": [
                    {
                        "type": "string",
                        "description": "accessToken",
                        "name": "tkn",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "json body",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqCreate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/v0.1/room/join": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "room"
                ],
                "summary": "방 참여",
                "parameters": [
                    {
                        "type": "string",
                        "description": "accessToken",
                        "name": "tkn",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "json body",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqJoin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/v0.1/room/out": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "room"
                ],
                "summary": "방 나가기",
                "parameters": [
                    {
                        "type": "string",
                        "description": "accessToken",
                        "name": "tkn",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "json body",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqOut"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/v0.1/room/ready": {
            "post": {
                "description": "■ errCode with 400\nPARAM_BAD : 파라미터 오류\n\n■ errCode with 500\nINTERNAL_SERVER : 내부 로직 처리 실패\nINTERNAL_DB : DB 처리 실패",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "room"
                ],
                "summary": "게임 준비 상태 변경",
                "parameters": [
                    {
                        "type": "string",
                        "description": "accessToken",
                        "name": "tkn",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "json body",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReqReady"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "request.ReqCreate": {
            "type": "object",
            "required": [
                "max_count",
                "min_count",
                "name"
            ],
            "properties": {
                "max_count": {
                    "type": "integer"
                },
                "min_count": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.ReqJoin": {
            "type": "object",
            "required": [
                "roomID"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "roomID": {
                    "type": "integer"
                }
            }
        },
        "request.ReqOut": {
            "type": "object",
            "properties": {
                "roomID": {
                    "type": "integer"
                }
            }
        },
        "request.ReqReady": {
            "type": "object",
            "properties": {
                "playerState": {
                    "type": "string"
                },
                "roomID": {
                    "type": "integer"
                }
            }
        },
        "request.ReqSignin": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.ReqSignup": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "response.ResSignin": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
