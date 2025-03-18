package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nimilgp/URLcommentary/internal/dblayer"
	"github.com/nimilgp/URLcommentary/internal/graphdb"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  *slog.Logger
	queries *dblayer.Queries
	ctx     context.Context
}

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "user=newuser password=2003 dbname=URLC host=localhost port=5432")

	if err != nil {
		log.Fatalf("Could Not Establish DB Connection %v", err)
	}

	defer conn.Close(ctx)
	queries := dblayer.New(conn)

	err = graphdb.InitializeNeo4j()
	if err != nil {
		log.Fatalf("Neo4j Initialization Error: %v", err)
	}
	defer graphdb.CloseNeo4j()

	var cfg config
	flag.IntVar(&cfg.port, "port", 3333, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config:  cfg,
		logger:  logger,
		queries: queries,
		ctx:     ctx,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Reminder", "CORS policy", "set before production")
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
