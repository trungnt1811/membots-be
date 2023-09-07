# Affiliate System Service

Astra Affiliate System API and SDK repository.

## A typical top-level directory layout

    .
    ├── cmd                    # Main applications for this project.
    ├── conf                   # Configuration file templates or default configs.
        ├── database           # Config for database.
        ├── environment        # Configuration for environment variables.
    ├── docs                   # Swagger API.
    ├── internal               # Private application and library code.
        ├── app
            ├── v1
                ├── dto                # Data transfer object.
                ├── middleware          # Middleware of REST server.
                ├── model              # Model of RDBMS or NoSQL.
                ├── route              # Routing.
                ├── util               # Utility.
            ├── another-app     # If your project contains many application.    
        ├── pkg                 # External code that you want to publish as library.
    ├── scripts                 # Scripts to perform various build, test,...
    └── README.md

## Dependency

Golang Template uses a number of open source projects to work properly:

- [Gorm] - The fantastic ORM library for Golang
- [Gin] - Gin is a web framework written in Go
- [Gin-Swagger] - Middleware to automatically generate RESTful API documentation with Swagger 2.0
- [Go-Redis] - Go redis cache
- [Prometheus-Golang-Client] - Go client library for Prometheus

## Installation

This application requires [Golang](https://golang.org/) v1.19+ to run.

Install the dependencies and devDependencies and start the server.

```sh
export PATH=$PATH:$HOME/go/bin
go install github.com/swaggo/swag/cmd/swag@latest
go mod init github.com/AstraProtocol/affiliate-system
go mod tidy
go get github.com/codegangsta/gin
cp example.env .env
```

## Run

```sh
./script/dev.sh
```

## Build

```sh
go build -o ./main ./cmd/main.go
```

## Test

```sh
./script/test.sh
```

## Docker

This application is very easy to install and deploy in a Docker container.

By default, the Docker will expose port 8080, so change this within the
Dockerfile if necessary. When ready, simply use the Dockerfile to
build the image.

```sh
cd your-project-name
docker build -t <youruser>/your-project-name:${version} .
```

Once done, run the Docker image and map the port to whatever you wish on
your host. In this example, we simply map port 8080 of the host to
port 8080 of the Docker (or whatever port was exposed in the Dockerfile):

```sh
docker run -d -p 8080:8080 --restart=always --name=your-project-name <youruser>/your-project-name:${version}
```

Verify the deployment by navigating to your server address in
your preferred browser.

```sh
127.0.0.1:8080
```

## Swagger init
    swag init -g internal/app.go

## CI/CD at TIKI

Read document here: https://docs.tiki.com.vn/pages/viewpage.action?spaceKey=HAN&title=Tiki+CI

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen)

   [Gin]: <https://github.com/gin-gonic/gin>
   [Gorm]: <https://gorm.io/>
   [Logrus]: <https://github.com/sirupsen/logrus>
   [Go-Redis]: <https://github.com/go-redis/redis>
   [Prometheus-Golang-Client]: <https://github.com/prometheus/client_golang>
   [Gin-Swagger]: <https://github.com/swaggo/gin-swagger>
   [Newrelic]: <https://github.com/newrelic/go-agent>
   [MongoDB]: <https://www.mongodb.com/2>
   [Jaeger]: <https://github.com/jaegertracing/jaeger-client-go>

