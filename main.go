package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TBabs-codes/mySousChef/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed static/*
var staticFiles embed.FS

// Config for application
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// App represents dependencies and server
type App struct {
	Router *chi.Mux
	Server *http.Server
	Config Config
	apiCfg apiConfig
}

// apiConfig contains db connection
type apiConfig struct {
	DB          *database.Queries
	JWTSecret   string
	JWT_Timeout time.Duration
}

func (a *App) Initialize() {
	a.Router = chi.NewRouter()

	a.Server = &http.Server{
		Addr:         ":" + a.Config.Port,
		Handler:      a.Router,
		ReadTimeout:  a.Config.ReadTimeout,
		WriteTimeout: a.Config.WriteTimeout,
	}

	db, err := sql.Open("sqlite3", "./recipe_book.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	a.apiCfg.DB = database.New(db)

	a.apiCfg.JWTSecret = os.Getenv("JWT_SECRET")
	a.apiCfg.JWT_Timeout = 24 * time.Hour
}

func (a *App) setupRoutes() {
	a.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//Serves start page
	a.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Current working directory:", cwd)
		f, err := staticFiles.Open("static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	//Helps server the other pages of the app
	// Add this to your router setup
	a.Router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		// Remove the leading "/static/" from the path
		filePath := strings.TrimPrefix(r.URL.Path, "/static/")

		// Open the file from the embedded filesystem
		f, err := staticFiles.Open("static/" + filePath)
		if err != nil {
			// If file not found, return a 404 error
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		// Determine the content type based on file extension
		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		// Copy the file contents to the response
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)

	if a.apiCfg.DB != nil {
		v1Router.Post("/register", a.apiCfg.handlerCreateUser)
		v1Router.Post("/login", a.apiCfg.handlerLoginUser)
		// v1Router.Get("/notes", a.apiCfg.middlewareAuth(apiCfg.handlerNotesGet))
		// v1Router.Post("/notes", a.apiCfg.middlewareAuth(apiCfg.handlerNotesCreate))
	}

	a.Router.Mount("/v1", v1Router)
}

func (a *App) Run() {
	log.Printf("Serving on port: %s\n", a.Config.Port)
	log.Fatal(a.Server.ListenAndServe())
}

func main() {
	fmt.Printf("Recipe Book server starting...\n")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("unable to read environment variables: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	config := Config{
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	app := App{
		Config: config,
	}

	app.Initialize()
	app.setupRoutes()
	app.Run()
}
