package main

import (
	"fmt"
	"log"
	"net/http"
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
	fmt.Printf("\033[0;32mserver running on localhost %s \033[0;37m\n", port)
	log.Fatal(http.ListenAndServe(port, srv))
}
