package docs

import "github.com/swaggo/swag"

type s struct{}

func (s *s) ReadDoc() string {
    return `{
  "swagger": "2.0",
  "info": {
    "title": "natmap-go API",
    "description": "REST API with JWT auth, refresh tokens, roles and permissions",
    "version": "1.0.0"
  },
  "schemes": ["http"],
  "basePath": "/",
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "JWT token. Use: Bearer <token>"
    }
  },
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "responses": {"200": {"description": "ok"}}
      }
    },
    "/api/v1/auth/register": {
      "post": {
        "summary": "Register new user",
        "consumes": ["application/json"],
        "parameters": [{"in":"body","name":"body","schema":{"type":"object","required":["username","password"],"properties":{"username":{"type":"string"},"password":{"type":"string"},"email":{"type":"string"}}}}],
        "responses": {"201": {"description": "created"}, "400": {"description": "invalid_request"}}
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login",
        "consumes": ["application/json"],
        "parameters": [{"in":"body","name":"body","schema":{"type":"object","required":["username","password"],"properties":{"username":{"type":"string"},"password":{"type":"string"}}}}],
        "responses": {"200": {"description": "token returned"}, "401": {"description": "invalid_credentials"}}
      }
    },
    "/api/v1/auth/refresh": {
      "post": {
        "summary": "Refresh token",
        "consumes": ["application/json"],
        "parameters": [{"in":"body","name":"body","schema":{"type":"object","required":["refresh_token"],"properties":{"refresh_token":{"type":"string"}}}}],
        "responses": {"200": {"description": "new tokens"}, "401": {"description": "invalid_refresh"}}
      }
    },
    "/api/v1/auth/logout": {
      "post": {
        "summary": "Logout by revoking refresh token",
        "consumes": ["application/json"],
        "parameters": [{"in":"body","name":"body","schema":{"type":"object","required":["refresh_token"],"properties":{"refresh_token":{"type":"string"}}}}],
        "responses": {"204": {"description": "no content"}, "400": {"description": "invalid_request"}}
      }
    },
    "/api/v1/auth/logout_all": {
      "post": {
        "summary": "Logout all sessions",
        "security": [{"BearerAuth": []}],
        "responses": {"204": {"description": "no content"}, "401": {"description": "unauthorized"}}
      }
    },
    "/api/v1/me": {
      "get": {
        "summary": "Current user info",
        "security": [{"BearerAuth": []}],
        "responses": {"200": {"description": "user info"}, "401": {"description": "unauthorized"}}
      }
    },
    "/api/v1/admin/users": {
      "get": {
        "summary": "List users",
        "security": [{"BearerAuth": []}],
        "responses": {"200": {"description": "users"}, "403": {"description": "forbidden"}}
      }
    },
    "/api/v1/admin/users/{id}/role": {
      "put": {
        "summary": "Set user role",
        "security": [{"BearerAuth": []}],
        "parameters": [
          {"in":"path","name":"id","required":true,"type":"integer"},
          {"in":"body","name":"body","schema":{"type":"object","required":["role"],"properties":{"role":{"type":"string"}}}}
        ],
        "responses": {"204": {"description": "no content"}, "403": {"description": "forbidden"}}
      }
    },
    "/api/v1/admin/permissions": {
      "get": {
        "summary": "List permissions",
        "security": [{"BearerAuth": []}],
        "responses": {"200": {"description": "permissions"}, "403": {"description": "forbidden"}}
      },
      "post": {
        "summary": "Create permission",
        "security": [{"BearerAuth": []}],
        "consumes": ["application/json"],
        "parameters": [{"in":"body","name":"body","schema":{"type":"object","required":["role","resource","action"],"properties":{"role":{"type":"string"},"resource":{"type":"string"},"action":{"type":"string"},"allowed":{"type":"boolean"}}}}],
        "responses": {"201": {"description": "created"}, "403": {"description": "forbidden"}}
      }
    },
    "/api/v1/admin/permissions/{id}": {
      "delete": {
        "summary": "Delete permission",
        "security": [{"BearerAuth": []}],
        "parameters": [{"in":"path","name":"id","required":true,"type":"integer"}],
        "responses": {"204": {"description": "no content"}, "403": {"description": "forbidden"}}
      }
    },
    "/api/v1/reports": {
      "get": {
        "summary": "Read reports (permission protected)",
        "security": [{"BearerAuth": []}],
        "responses": {"200": {"description": "items"}, "403": {"description": "forbidden"}}
      }
    }
  }
}`
}

func init() { swag.Register(swag.Name, &s{}) }

