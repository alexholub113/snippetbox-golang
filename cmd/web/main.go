package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

type config struct {
	addr      string
	staticDir string
	logSource bool
}

func main() {
	cfg := setConfig()
	logger := createLogger(cfg)
	app := &application{logger: logger}

	logger.Info("starting server", slog.String("addr", cfg.addr))

	err := http.ListenAndServe(cfg.addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func setConfig() config {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticDir", "./ui/static", "Path to static assets")
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
