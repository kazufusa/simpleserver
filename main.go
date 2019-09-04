// +build !windows

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const port = "8080"

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprint(w, "hello world")
	log.Printf("request %v\n", r)
}

func main() {
	server := http.Server{Addr: ":" + port, Handler: http.HandlerFunc(handler)}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Failed to gracefully shutdown HTTPServer:", err)
	}
	log.Println("HTTPServer shutdown.")
}
