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
        "/api/v1/app/aff-campaign": {
            "get": {
                "description": "Get list of all aff campaign",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "Get list of all aff campaign",
                "parameters": [
                    {
                        "type": "string",
                        "description": "page to query, default is 1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "size to query, default is 10",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.AffCampaignAppDtoResponse"
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
        "/api/v1/app/aff-campaign/{id}": {
            "get": {
                "description": "Get aff campaign by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "app"
                ],
                "summary": "Get aff campaign by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id to query",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.AffCampaignAppDto"
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
        "/api/v1/campaign/link": {
            "post": {
                "description": "Request create a new campaign link or pick the active old one",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "redeem"
                ],
                "summary": "PostGenerateAffLink",
                "parameters": [
                    {
                        "description": "Request create link payload, required",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CreateLinkPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "when success, return the created link for this request campaign",
                        "schema": {
                            "$ref": "#/definitions/dto.CreateLinkResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
        "/api/v1/console/aff-campaign": {
            "get": {
                "description": "Get list aff campaign",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "console"
                ],
                "summary": "Get list aff campaign",
                "parameters": [
                    {
                        "type": "string",
                        "description": "by to query, default is all",
                        "name": "stella_status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "order to query, default is desc",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "page to query, default is 1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "size to query, default is 10",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.AffCampaignDtoResponse"
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
        "/api/v1/console/aff-campaign/{id}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "update campaign info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "console"
                ],
                "summary": "update campaign info",
                "parameters": [
                    {
                        "description": "Campaign info to update, required",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.AffCampaignAppDto"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "id to query",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.AffCampaignAppDto"
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
        "/api/v1/order/postback": {
            "post": {
                "description": "A callback to receive order from AccessTrade",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "redeem"
                ],
                "summary": "PostBackOrderHandle",
                "parameters": [
                    {
                        "description": "Request create link payload, required",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ATPostBackRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "when success, return the modified order for this post back",
                        "schema": {
                            "$ref": "#/definitions/dto.ATPostBackResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
        "/api/v1/redeem/request": {
            "post": {
                "description": "Send cashback to customer wallet by redeem code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "redeem"
                ],
                "summary": "PostRequestRedeem",
                "parameters": [
                    {
                        "description": "Redeem payload, required",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.RedeemRequestPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "when redeem code is available, only valid if not claimed yet",
                        "schema": {
                            "$ref": "#/definitions/types.RedeemRewardResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
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
        "/api/v1/users/{uid}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get a user detail by uid",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "GetUserDetail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User uid to query",
                        "name": "uid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "dto.ATPostBackRequest": {
            "type": "object",
            "properties": {
                "browser": {
                    "description": "Trình duyệt sử dụng",
                    "type": "string"
                },
                "campaign_id": {
                    "description": "ID của campaign trên hệ thống",
                    "type": "string"
                },
                "click_time": {
                    "description": "Thời gian phát sinh click",
                    "type": "string"
                },
                "conversion_platform": {
                    "description": "Platform sử dụng",
                    "type": "string"
                },
                "customer_type": {
                    "description": "Thuộc tính của khách hàng phụ thuộc theo campaigns",
                    "type": "string"
                },
                "ip": {
                    "description": "IP phát sinh đơn hàng",
                    "type": "string"
                },
                "is_confirmed": {
                    "description": "Đơn hàng khóa data và được thanh toán: 0: chưa đối soát, 1: đã đối soát",
                    "type": "integer"
                },
                "order_id": {
                    "description": "Mã đơn hàng hiển thị trên trang pub",
                    "type": "string"
                },
                "product_category": {
                    "description": "Group commission của sản phẩm",
                    "type": "string"
                },
                "product_id": {
                    "description": "Mã sản phẩm",
                    "type": "string"
                },
                "product_price": {
                    "description": "Giá của một sản phẩm",
                    "type": "number"
                },
                "publisher_login_name": {
                    "type": "string"
                },
                "quantity": {
                    "description": "Số lượng sản phẩm",
                    "type": "integer"
                },
                "referrer": {
                    "description": "click_referrer",
                    "type": "string"
                },
                "reward": {
                    "description": "Hoa hồng nhận được",
                    "type": "number"
                },
                "sales_time": {
                    "description": "Thời gian phát sinh của đơn hàng",
                    "type": "string"
                },
                "status": {
                    "description": "Status của đơn hàng gồm 3 giá trị: 0: new, 1: approved, 2: rejected",
                    "type": "integer"
                },
                "transaction_id": {
                    "description": "Mã unique trên hệ thống AccessTrade",
                    "type": "string"
                },
                "utm_campaign": {
                    "description": "Thông tin tùy biến pub truyền vào url trong param utm_campaign",
                    "type": "string"
                },
                "utm_content": {
                    "description": "Thông tin tùy biến pub truyền vào url trong param utm_content",
                    "type": "string"
                },
                "utm_medium": {
                    "description": "Thông tin tùy biến pub truyền vào url trong param utm_medium",
                    "type": "string"
                },
                "utm_source": {
                    "description": "Thông tin tùy biến pub truyền vào url trong param utm_source",
                    "type": "string"
                }
            }
        },
        "dto.ATPostBackResponse": {
            "type": "object",
            "properties": {
                "order_id": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "dto.AffCampaignAppDto": {
            "type": "object",
            "properties": {
                "accesstrade_id": {
                    "type": "string"
                },
                "brand_id": {
                    "type": "integer"
                },
                "category_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "end_time": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                },
                "stella_description": {},
                "stella_max_com": {
                    "type": "number"
                },
                "stella_status": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "dto.AffCampaignAppDtoResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.AffCampaignAppDto"
                    }
                },
                "next_page": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "dto.AffCampaignDto": {
            "type": "object",
            "properties": {
                "accesstrade_id": {
                    "type": "string"
                },
                "active_status": {
                    "type": "integer"
                },
                "approval": {
                    "type": "string"
                },
                "brand_id": {
                    "type": "integer"
                },
                "category": {
                    "type": "string"
                },
                "category_id": {
                    "type": "integer"
                },
                "cookie_duration": {
                    "type": "integer"
                },
                "cookie_policy": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "$ref": "#/definitions/dto.CampaignDescriptionDto"
                },
                "end_time": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "logo": {
                    "type": "string"
                },
                "max_com": {
                    "type": "string"
                },
                "merchant": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "scope": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "stella_description": {},
                "stella_max_com": {
                    "type": "number"
                },
                "stella_status": {
                    "type": "string"
                },
                "sub_category": {
                    "type": "string"
                },
                "thumbnail": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "dto.AffCampaignDtoResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.AffCampaignDto"
                    }
                },
                "next_page": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "dto.CampaignDescriptionDto": {
            "type": "object",
            "properties": {
                "action_point": {
                    "type": "string"
                },
                "campaign_id": {
                    "type": "integer"
                },
                "commission_policy": {
                    "type": "string"
                },
                "cookie_policy": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "introduction": {
                    "type": "string"
                },
                "other_notice": {
                    "type": "string"
                },
                "rejected_reason": {
                    "type": "string"
                },
                "traffic_building_policy": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "dto.CreateLinkPayload": {
            "type": "object",
            "properties": {
                "campaign_id": {
                    "type": "integer"
                },
                "original_url": {
                    "type": "string"
                },
                "shorten_link": {
                    "type": "boolean"
                }
            }
        },
        "dto.CreateLinkResponse": {
            "type": "object",
            "properties": {
                "aff_link": {
                    "type": "string"
                },
                "campaign_id": {
                    "type": "integer"
                },
                "original_url": {
                    "type": "string"
                },
                "short_link": {
                    "type": "string"
                }
            }
        },
        "types.RedeemRequestPayload": {
            "type": "object",
            "required": [
                "redeemCode",
                "walletAddress"
            ],
            "properties": {
                "redeemCode": {
                    "description": "a unique text to redeem reward",
                    "type": "string"
                },
                "walletAddress": {
                    "description": "valid wallet address to send reward to",
                    "type": "string"
                }
            }
        },
        "types.RedeemRewardResponse": {
            "type": "object",
            "properties": {
                "deadline": {
                    "type": "integer"
                },
                "holderAddress": {
                    "type": "string"
                },
                "signature": {
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
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Use for authorization of reward creator",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Affiliate System API",
	Description:      "Use for authorization during server to server calls",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
