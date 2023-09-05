#!/usr/bin/env bash

# Run install first: "go install github.com/swaggo/swag/cmd/swag@latest"
swag init -g internal/app/affiliate/app.go
gin --appPort 8080 --port 9000 --build cmd/main.go --immediate run cmd/main.go
