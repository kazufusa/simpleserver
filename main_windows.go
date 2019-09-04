package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"
)

const port = "8080"

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprint(w, "hello world")
	log.Println("sleep sitayo")
}

func main() {
	fullpath, _ := filepath.Abs(os.Args[0])
	name := filepath.Base(fullpath)
	err := exec.Command(
		"netsh",
		"advfirewall",
		"firewall",
		"add",
		"rule",
		fmt.Sprintf("name=\"%s\"", name),
		"dir=in",
		"action=allow",
		fmt.Sprintf("program=\"%s\"", fullpath),
		"enable=yes",
		"protocol=TCP",
		fmt.Sprintf("localport=%s", port),
	).Run()
	if err != nil {
		log.Fatal(err)
	}

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

	err = exec.Command(
		"netsh",
		"advfirewall",
		"firewall",
		"delete",
		"rule",
		fmt.Sprintf("name=\"%s\"", name),
	).Run()
	if err != nil {
		log.Fatal(err)
	}
}
