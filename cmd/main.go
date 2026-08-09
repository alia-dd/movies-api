package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := ":8000"
	fmt.Println("server running on localhost", port)
	srv, err := route()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(port, srv))

}
