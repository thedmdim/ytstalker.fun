package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ytstalker/backend/api"
	"ytstalker/backend/conf"
	"ytstalker/backend/youtube"

	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	config := conf.ParseConfig("conf.json")

	// prepare db
	db, err := sqlitex.Open(config.DSN, 0, 100)
	if err != nil {
		log.Fatal("cannot open db", err)
	}
	dbScheme, err := os.ReadFile("tables.sql")
	if err != nil {
		log.Fatal("cannot open tables.sql: ", err.Error())
	}
	conn := db.Get(context.Background())
	if err := sqlitex.ExecuteScript(conn, string(dbScheme), nil); err != nil {
		log.Fatal("cannot create db: ", err)
	}
	db.Put(conn)
	log.Println("database ready")

	// init youtube api requester
	ytr := youtube.NewYouTubeRequester(config)

	// make router
	handler := api.NewRouter(db, ytr)
	server := &http.Server{
		Handler: handler,
	}

	// serve
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe: %v", err)
        }
	}()

	// preparation for gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// wait for a stop
	<- stop

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("successfully finished serving")

	err = db.Close()
	if err != nil {
		log.Println("error gracefully closing db:", err.Error())
	}
	log.Println("successfully closed db")
	log.Println("thanks :)")
}
