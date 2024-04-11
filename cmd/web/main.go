package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.oleksandrholub.com/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

type config struct {
	addr             string
	staticDir        string
	logSource        bool
	connectionString string
}

func main() {
	cfg := setConfig()
	logger := createLogger(cfg)
	db, dbErr := openDbConnection(cfg)
	if dbErr != nil {
		logger.Error(dbErr.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func setConfig() config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.connectionString, "connectionString", "root:golanguserpassword@/snippetbox?parseTime=true", "MySQL connection string")
	flag.BoolVar(&cfg.logSource, "logSource", false, "Include source file and line number in log output")

	flag.Parse()

	return cfg
}

func createLogger(cfg config) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: cfg.logSource,
	}))
}

func openDbConnection(cfg config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
