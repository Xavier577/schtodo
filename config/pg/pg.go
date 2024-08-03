package pg

import (
	"context"
	"fmt"
	"github.com/Xavier577/schtodo/pkg/objects"
	"log"
	"strings"
	"time"

	"database/sql"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type PgConnectCfg struct {
	Host     string `json:"host"`
	PORT     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

func parseCfg(pgCfg *PgConnectCfg) string {
	cfgMap, err := objects.MarshalStructToMap(pgCfg)

	if err != nil {
		log.Fatal(err)
	}

	var cfgPairs []string

	for key, val := range cfgMap {
		if val != nil && val != "" {
			cfgPairs = append(cfgPairs, fmt.Sprintf("%s=%v", key, val))
		}
	}

	return strings.Join(cfgPairs, " ")
}

func migrate(db *sql.DB) error {
	driver := "postgres"

	goose.SetBaseFS(os.DirFS("."))

	if err := goose.SetDialect(driver); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

var PgClient *sqlx.DB

func Connect(pgCfg *PgConnectCfg) *sqlx.DB {

	dataSource := parseCfg(pgCfg)

	db, err := sqlx.Connect("postgres", dataSource)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Run Goose migrations
	if err := migrate(db.DB); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	return db

}
