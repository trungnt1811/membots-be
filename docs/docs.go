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
        "/api/v1/launchpad": {
            "get": {
                "description": "Get launchpad history by address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "launchpad"
                ],
                "summary": "Get launchpad history by address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "memeAddress to query, default is ",
                        "name": "memeAddress",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.LaunchpadInfoResp"
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
        "/api/v1/meme": {
            "get": {
                "description": "Get meme by meme address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memeception"
                ],
                "summary": "Get meme by meme address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "memeAddress to query, default is ",
                        "name": "memeAddress",
                        "in": "query"
                    },
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
                            "$ref": "#/definitions/dto.MemeceptionDetailResp"
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
                            "$ref": "#/definitions/dto.MemeceptionsResp"
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
        "/api/v1/memes": {
            "post": {
                "description": "Create meme",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "memeception"
                ],
                "summary": "Create meme",
                "parameters": [
                    {
                        "description": "Request create meme, required",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateMemePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "When success, return {\"success\": true}",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateMemePayload"
                        }
                    },
                    "417": {
                        "description": "Expectation Failed",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    },
                    "424": {
                        "description": "Failed Dependency",
                        "schema": {
                            "$ref": "#/definitions/util.GeneralError"
                        }
                    }
                }
            }
        },
        "/api/v1/quote": {
            "get": {
                "description": "Swap router",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "swap"
                ],
                "summary": "Swap router",
                "parameters": [
                    {
                        "type": "string",
                        "description": "protocols to query, default is v3",
                        "name": "protocols",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tokenInAddress to query",
                        "name": "tokenInAddress",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tokenInChainId to query",
                        "name": "tokenInChainId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tokenOutAddress to query",
                        "name": "tokenOutAddress",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "tokenOutChainId to query",
                        "name": "tokenOutChainId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "amount to query",
                        "name": "amount",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "type to query",
                        "name": "type",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
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
        "/api/v1/stats": {
            "get": {
                "description": "Get stats by memeAddress",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stats"
                ],
                "summary": "Get stats by memeAddress",
                "parameters": [
                    {
                        "type": "string",
                        "description": "memeAddress to query",
                        "name": "memeAddress",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
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
        "/api/v1/swaps": {
            "get": {
                "description": "Get swaps history by memeAddress",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "swap"
                ],
                "summary": "Get swaps history by memeAddress",
                "parameters": [
                    {
                        "type": "string",
                        "description": "memeAddress to query, default is ",
                        "name": "memeAddress",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.SwapHistoryByAddressResp"
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
        "dto.BlockchainPayload": {
            "type": "object",
            "properties": {
                "chainId": {
                    "type": "string"
                },
                "creatorAddress": {
                    "type": "string"
                }
            }
        },
        "dto.CreateMemePayload": {
            "type": "object",
            "properties": {
                "blockchain": {
                    "$ref": "#/definitions/dto.BlockchainPayload"
                },
                "memeInfo": {
                    "$ref": "#/definitions/dto.MemeInfoPayload"
                },
                "memeception": {
                    "$ref": "#/definitions/dto.MemeceptionPayload"
                },
                "socials": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/dto.Social"
                    }
                }
            }
        },
        "dto.LaunchpadInfo": {
            "type": "object",
            "properties": {
                "collectedETH": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "targetETH": {
                    "type": "string"
                },
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Transaction"
                    }
                }
            }
        },
        "dto.LaunchpadInfoResp": {
            "type": "object",
            "properties": {
                "launchpadInfo": {
                    "$ref": "#/definitions/dto.LaunchpadInfo"
                }
            }
        },
        "dto.LaunchpadTx": {
            "type": "object",
            "properties": {
                "amountUSD": {
                    "type": "string"
                },
                "logoUrl": {
                    "type": "string"
                },
                "memeContractAddress": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                },
                "txHash": {
                    "type": "string"
                },
                "txType": {
                    "type": "string"
                },
                "walletAddress": {
                    "type": "string"
                }
            }
        },
        "dto.MemeCommon": {
            "type": "object",
            "properties": {
                "bannerUrl": {
                    "type": "string"
                },
                "contract_address": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "logoUrl": {
                    "type": "string"
                },
                "meta": {
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
        "dto.MemeDetail": {
            "type": "object",
            "properties": {
                "bannerUrl": {
                    "type": "string"
                },
                "contractAddress": {
                    "type": "string"
                },
                "creatorAddress": {
                    "type": "string"
                },
                "decimals": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "live": {
                    "type": "boolean"
                },
                "logoUrl": {
                    "type": "string"
                },
                "memeception": {
                    "$ref": "#/definitions/dto.Memeception"
                },
                "meta": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "networkId": {
                    "type": "integer"
                },
                "socials": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/dto.Social"
                    }
                },
                "swapFeeBps": {
                    "type": "integer"
                },
                "symbol": {
                    "type": "string"
                },
                "totalSupply": {
                    "type": "string"
                },
                "vestingAllocBps": {
                    "type": "integer"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "dto.MemeInfoPayload": {
            "type": "object",
            "properties": {
                "ama": {
                    "type": "boolean"
                },
                "bannerUrl": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "logoUrl": {
                    "type": "string"
                },
                "meta": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "salt": {
                    "type": "string"
                },
                "swapFeeBps": {
                    "type": "string"
                },
                "swapFeePct": {
                    "type": "number"
                },
                "symbol": {
                    "type": "string"
                },
                "targetETH": {
                    "type": "number"
                },
                "telegram": {
                    "type": "string"
                },
                "vestingAllocBps": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "dto.Memeception": {
            "type": "object",
            "properties": {
                "ama": {
                    "type": "boolean"
                },
                "collectedETH": {
                    "type": "string"
                },
                "contractAddress": {
                    "type": "string"
                },
                "enabled": {
                    "type": "boolean"
                },
                "meme": {
                    "$ref": "#/definitions/dto.MemeCommon"
                },
                "memeID": {
                    "type": "integer"
                },
                "startAt": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "targetETH": {
                    "type": "string"
                },
                "updatedAtEpoch": {
                    "type": "integer"
                }
            }
        },
        "dto.MemeceptionByStatus": {
            "type": "object",
            "properties": {
                "live": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Memeception"
                    }
                },
                "past": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Memeception"
                    }
                }
            }
        },
        "dto.MemeceptionDetailResp": {
            "type": "object",
            "properties": {
                "meme": {
                    "$ref": "#/definitions/dto.MemeDetail"
                },
                "price": {
                    "type": "integer"
                }
            }
        },
        "dto.MemeceptionPayload": {
            "type": "object",
            "properties": {
                "ema": {
                    "type": "boolean"
                },
                "startAt": {
                    "type": "string"
                },
                "targetETH": {
                    "type": "string"
                }
            }
        },
        "dto.MemeceptionsResp": {
            "type": "object",
            "properties": {
                "latestCoins": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Memeception"
                    }
                },
                "latestLaunchpadTx": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.LaunchpadTx"
                    }
                },
                "memeceptionsByStatus": {
                    "$ref": "#/definitions/dto.MemeceptionByStatus"
                },
                "price": {
                    "type": "integer"
                }
            }
        },
        "dto.Social": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string"
                },
                "photoUrl": {
                    "type": "string"
                },
                "provider": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "dto.Swap": {
            "type": "object",
            "properties": {
                "amount0": {
                    "type": "string"
                },
                "amount1": {
                    "type": "string"
                },
                "amountUSD": {
                    "type": "string"
                },
                "buy": {
                    "type": "boolean"
                },
                "priceETH": {
                    "type": "string"
                },
                "priceUSD": {
                    "type": "string"
                },
                "swapAt": {
                    "type": "integer"
                },
                "token1IsMeme": {
                    "type": "boolean"
                },
                "txHash": {
                    "type": "string"
                },
                "walletAddress": {
                    "type": "string"
                }
            }
        },
        "dto.SwapHistoryByAddressResp": {
            "type": "object",
            "properties": {
                "swaps": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Swap"
                    }
                }
            }
        },
        "dto.Transaction": {
            "type": "object",
            "properties": {
                "amountETH": {
                    "type": "string"
                },
                "amountMeme": {
                    "type": "string"
                },
                "epoch": {
                    "type": "integer"
                },
                "memeId": {
                    "type": "string"
                },
                "txHash": {
                    "type": "string"
                },
                "txType": {
                    "type": "string"
                },
                "walletAddress": {
                    "type": "string"
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
