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

	}
	log.Fatal(http.ListenAndServe(port, srv))
}

// this lowerpart still not done
// its getting movie data from TMD website's api
// we might be albe to use this as database seed
// if you're interested you can test it
/// this is the link im using https://developer.themoviedb.org/docs/getting-started
// you can use the apikey ist there if you want

// type TMDBMovie struct {
// 	ID          int    `json:"id"`
// 	Title       string `json:"title"`
// 	ReleaseDate string `json:"release_date"`
// 	GenreIDs    []int  `json:"genre_ids"`
// }
// type TMDGenre struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }
// type Actor struct {
// 	Id        int    `json:"id"`
// 	Name      string `json:"name"`
// 	BirthDate string `json:"birthDate"`
// }
// type TMDBResponse struct {
// 	Results []TMDBMovie `json:"results"`
// }

// func main() {
// 	apiKey := "add63e21c755467e1763245167c231b4"
// 	var allMovies []TMDBMovie
// 	for page := 1; page <= 50; page++ {
// 		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d", apiKey, page)

// 		resp, err := http.Get(url)
// 		if err != nil {
// 			panic(err)
// 		}

// 		body, _ := io.ReadAll(resp.Body)
// 		resp.Body.Close()

// 		var tmdbResp TMDBResponse
// 		if err := json.Unmarshal(body, &tmdbResp); err != nil {
// 			panic(err)
// 		}

// 		allMovies = append(allMovies, tmdbResp.Results...)

// 		fmt.Printf("Page %d fetched, total so far: %d\n", page, len(allMovies))

// 	}

// 	fmt.Printf("Total movies fetched: %d\n", len(allMovies))

// 	out, _ := json.MarshalIndent(allMovies, "", "  ")
// 	fmt.Println(out)
// 	// os.WriteFile("../internal/database/data/tmdb_movies.json", out, 0644)
// }
