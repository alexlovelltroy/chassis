# chassis

This is a basic microservice chassis that makes it easy to do the default things and possible to do more complex things

## HTTP Server

Gin is included as the HTTP server.  A standard health/ping route is included by default.  The `AddRoute()` method is the preferred way of adding additional routes.  See the example below.

## Commandline

Cobra is used for the commands.  `migrate` and `serve` are included commands with reasonable flags

## Logging

Logsrus provides semi-structured logging.  Debug output can be triggered with the -D option

## Database

Postgres is supported through the sqlx package with migration support through golang-migrate/migrate/v4

## Example
```
package main

import (
	"github.com/alexlovelltroy/chassis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type ExampleService struct {
	*chassis.Microservice
}

type ExampleServiceConfig struct {
	chassis.MicroserviceConfig
	parallelism int
}

type Example struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

func (s *ExampleService) CreateExample(c *gin.Context) {
	s.DB.Exec("INSERT INTO example (name) VALUES ($1)", "example")
	c.JSON(200, gin.H{
		"message":    "value",
		"statusCode": "statusCode",
	})
}

func main() {
	cfg := &ExampleServiceConfig{
		MicroserviceConfig: chassis.DefaultMicroserviceConfig(),
		parallelism:        1,
	}
	service := ExampleService{
		Microservice: chassis.NewMicroservice(cfg),
	}
	chassis.ServeCmd.Run = func(cmd *cobra.Command, args []string) {
		service.Init() // Establish connection(s) to external services and configure the gin router
		service.AddRoute("POST", "/example", service.CreateExample)
		service.Serve() // Start the gin router
	}
	chassis.Execute()
}
```
