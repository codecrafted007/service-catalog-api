package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codecrafted007/service-catalog-api/internal/handler"
	"github.com/codecrafted007/service-catalog-api/internal/logger"
	"github.com/codecrafted007/service-catalog-api/internal/middleware"
	"github.com/codecrafted007/service-catalog-api/internal/storage"
	"github.com/codecrafted007/service-catalog-api/internal/storage/sqlite"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()

	defer func() {
		if err := logger.L().Sync(); err != nil {
			os.Stderr.WriteString("Failed to sync logger: " + err.Error() + "\n")
		}
	}()
	logger.L().Info("Starting service catalog API")

	driver := flag.String("db-driver", "sqlite3", "Database driver: sqlite3|mysql|postgres")
	dataSourceName := flag.String("db-dsn", "services.db", "Data source name or file path")
	httpPort := flag.String("port", "8080", "HTTP server port")
	flag.Parse()

	if err := ensureSchemaExists(*driver, *dataSourceName); err != nil {
		log.Fatal("failed to initialize schema ", err)
	}

	// Init sqlite store (for MySQL/Postgres we can add support later)
	store, err := sqlite.New("services.db")
	if err != nil {
		log.Fatal("failed to connect to db", err)
	}
	initDatabase(store, logger.L())
	r := mux.NewRouter()
	r.Use(middleware.APIKeyAuth(store.IsValidAPIKey))

	h := handler.New(store, logger.L())

	r.HandleFunc("/services", h.ListServices).Methods("GET")
	r.HandleFunc("/services/{id}", h.GetServiceByID).Methods("GET")

	r.HandleFunc("/services", h.CreateService).Methods("POST")
	r.HandleFunc("/services/{id}", h.UpdateService).Methods("PUT")
	r.HandleFunc("/services/{id}", h.DeleteService).Methods("DELETE")

	addr := fmt.Sprintf(":%s", *httpPort)
	logger.L().Infof("Listening on %s", addr)
	http.ListenAndServe(addr, r)
}

func initDatabase(store storage.Storage, logger *zap.SugaredLogger) {
	var count int
	err := store.DB().Get(&count, "SELECT COUNT(*) FROM api_keys")
	if err != nil {
		logger.Fatal("failed to check api_keys table:", err)
	}

	if count == 0 {
		apiKey := generateAPIKey()
		_, err := store.DB().Exec("INSERT INTO api_keys (key) VALUES (?)", apiKey)
		if err != nil {
			logger.Fatal("failed to insert default API key:", err)
		}
		logger.Infof("Default API key generated: %s", apiKey)
	}
}

func generateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func ensureSchemaExists(driver, dsn string) error {
	schemaFile := map[string]string{
		"sqlite3":  "db/schema.sqlite.sql",
		"mysql":    "db/schema.mysql.sql",
		"postgres": "db/schema.postgres.sql",
	}[driver]

	if schemaFile == "" {
		return fmt.Errorf("unsupported DB driver: %s", driver)
	}

	// only needed for sqlite
	if driver == "sqlite3" {
		if _, err := os.Stat(dsn); err == nil {
			logger.L().Infof("SQLite DB already exists: %s", dsn)
			return nil
		}
	}

	content, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to apply schema: %w", err)
	}

	logger.L().Info("Schema applied successfully")
	return nil
}
