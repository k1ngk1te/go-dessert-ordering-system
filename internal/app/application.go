package app

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/redisstore"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type Application struct {
	DB        *sql.DB
	Loggers   *ApplicationLoggers
	Models    *ApplicationModels
	Services  *ApplicationServices
	RedisPool *redis.Pool
	Session   *ApplicationSession
	Templates *template.Template
}

// Render Template Helper Function
func (a *Application) RenderHtmlTemplate(w http.ResponseWriter, templateName string, data any, statusCode int) {
	// Set Content-Type header to HTML
	w.Header().Set("Content-Type", "text/html;charset=utf-8")

	w.WriteHeader(statusCode)

	// Execute the specified template, passing the data.
	err := a.Templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error rendering template %s: %v", templateName, err)
	}
}

func NewApplication() *Application {
	// Initiate loggers
	loggers := NewApplicationLoggers()

	err := godotenv.Load()
	if err != nil {
		loggers.Info.Printf("Warning: Could not load .env file: %v. Assuming environment variables are set externally.", err)
	}

	templates, err := NewApplicationTemplates()
	if err != nil {
		loggers.Error.Fatalf("Error parsing templates: %v", err)
	}

	// Initialize Redigo Redis Pool instead of go-redis client
	redisPool, err := openRedisPool()
	if err != nil {
		loggers.Error.Fatalf("could not connect to Redis: %v", err)
	} else {
		loggers.Info.Println("Successfully connected to Redis.")
	}

	// Initialize session manager
	sessionManager := openSession(loggers, redisstore.New(redisPool))
	session := NewApplicationSession(sessionManager)

	// Open a database connection
	// sql.Open doesn't actually connect to the database yet; it just validates the DSN format.
	db, err := openDB(os.Getenv("DSN"))
	if err != nil {
		loggers.Error.Fatalf("Error opening or connecting to database: %v", err)
	}
	loggers.Info.Println("Successfully connected to Database!")

	models := NewApplicationModels(db)
	services := NewApplicationServices(models)

	a := &Application{
		DB:        db,
		Loggers:   loggers,
		Models:    models,
		Services:  services,
		RedisPool: redisPool,
		Session:   session,
		Templates: templates,
	}
	return a
}
