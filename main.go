package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"schtodo/config/app"
	"schtodo/config/env"
	"schtodo/config/pg"
	"schtodo/domains/user"
	"schtodo/internal"
	"syscall"
	"time"
)

func main() {

	env.LoadEnv()

	SSLMode := os.Getenv("PG_SSL_MODE")

	if SSLMode == "" {
		SSLMode = "disable"
	}

	db := pg.Connect(&pg.PgConnectCfg{
		DBName:   os.Getenv("PG_DATABASE"),
		Host:     os.Getenv("PG_HOST"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		PORT:     os.Getenv("PG_PORT"),
		SSLMode:  SSLMode,
	})

	// gracefully close db
	defer db.Close()

	appCnt := &internal.AppContainer{DB: db, UserRepo: user.NewUserRepo(db)}

	appInstance := app.NewApp(appCnt)

	addr := fmt.Sprintf("127.0.0.1:%v", os.Getenv("PORT"))

	srv := &http.Server{
		Addr:    addr,
		Handler: appInstance.Handler(),
	}

	go func() {
		// service connections
		log.Printf("listening on: %v\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")

}
