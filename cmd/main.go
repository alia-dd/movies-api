package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8000"
	// database.FetchSeedData()
	srv, err := route()
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal(http.ListenAndServe(port, srv))
	fmt.Println("server running on localhost", port)
}

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func main() {

// 	port := ":8000"
// 	fmt.Println("server running on localhost", port)
// 	srv, err := route()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Fatal(http.ListenAndServe(port, srv))

// }
