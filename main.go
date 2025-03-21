package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TBabs-codes/recipe_book/internal/database"
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

//apiConfig contains db connection
type apiConfig struct {
	DB *database.Queries
}

func (a *App) Initialize() {
	a.Router = chi.NewRouter()
	a.setupRoutes()

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

	a.apiCfg.DB = database.New(db)
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

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)

	if a.apiCfg.DB != nil {
		v1Router.Post("/users", a.apiCfg.handlerUsersCreate)
		v1Router.Get("/users", a.apiCfg.middlewareAuth(a.apiCfg.handlerUsersGet))
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
	app.Run()
}
