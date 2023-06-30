package chassis

import (
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Microservice struct {
	ExecutableName string
	InstanceName   string
	InstaceID      uuid.UUID
	Version        string
	Router         *gin.Engine
	DB             *sqlx.DB
	Config         Config
}

func (m *Microservice) AddRoute(method string, path string, handler gin.HandlerFunc) {
	m.Router.Group("/").Handle(method, path, handler)
}

// NewMicroservice initializes the bare microservice.
func NewMicroservice(config Config) *Microservice {
	m := &Microservice{}
	m.InstaceID = uuid.New()
	m.InstanceName = namesgenerator.GetRandomName(0)
	m.Config = config
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	m.ExecutableName = filepath.Base(path)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.Info("Starting ", m.ExecutableName, " instance ", m.InstanceName, " (", m.InstaceID, ")")
	return m
}

// Init is a convenience function.  Developers that want to accept all reasonable defaults can use it to call the other initialization functions in one go.
func (m *Microservice) Init() {
	m.InitPostgres(DBConfig) // DBConfig is a global variable defined in config.go
	m.InitRouter()
}

// Sets up the Router
func (m *Microservice) InitRouter() {
	// Create a new Gin object
	gin.SetMode(gin.ReleaseMode)
	m.Router = gin.Default()
	m.Router.Use(otelgin.Middleware(m.ExecutableName))
	hostname, _ := os.Hostname()
	m.AddRoute("GET", "/health/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":    "value",
			"statusCode": "statusCode",
			"instance":   m.InstanceName,
			"instanceID": m.InstaceID,
			"version":    m.Version,
			"Executable": m.ExecutableName,
			"Hostname":   hostname,
		})
	})
}

// InitPostgres initializes the connection to Postgres
func (m *Microservice) InitPostgres(config PostgresConfig) {
	var err error
	m.DB, err = ConnectDB("postgres", config.ConnectionString)
	// Set reasonable connection parameters for our database connection(s)
	// https://www.alexedwards.net/blog/configuring-sqldb
	m.DB.SetConnMaxIdleTime(time.Duration(config.ConnectionMaxIdleTime) * time.Second)
	m.DB.SetConnMaxLifetime(time.Duration(config.ConnectionMaxLifetime) * time.Second)
	m.DB.SetMaxIdleConns(config.MaxIdleConnections)
	m.DB.SetMaxOpenConns(config.MaxOpenConnections)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
}

func (m *Microservice) Serve() {
	log.Info(m.ExecutableName, " gin server is available at http://", m.Config.GetListenAddress())
	log.Debug("Registered routes:")
	for _, r := range m.Router.Routes() {
		log.Debug(r.Method, " ", r.Path)
	}
	m.Router.Run(m.Config.GetListenAddress())

}
