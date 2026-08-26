package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// running this Fetch seed data func takes a long time
	// (approximately an hour for 200 movie its actors and genres)
	// so dont run it if not must
	// database.FetchSeedData()

	port := ":8800"
	srv, err := route()
	if err != nil {
		fmt.Println(err)
		return
	}
	server := &http.Server{
		Addr: port,
		Handler: srv,
	}
	serverError := make(chan error, 1)
	go func() {
		fmt.Printf("\033[0;32mserver running on localhost %s \033[0;37m\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- err
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select{
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
		log.Printf("Server shutdown error: %v", err)
		return
	}
	log.Println("server exited properly")
}
