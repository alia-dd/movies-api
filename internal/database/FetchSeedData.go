package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type TMDBMovie struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	ReleaseDate      string `json:"release_date"`
	Overview         string `json:"overview"`
	OriginalLanguage string `json:"original_language"`
	GenreIDs         []int  `json:"genre_ids"`
	ActorIDs         []int  `json:"Actor_ids"`
}

type TMDBActor struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"biography"`
	Character    string `json:"character"`
	PersonInfo   Person `json:"Person_Info"`
}
type Person struct {
	Bio          string `json:"biography"`
	BirthDate    string `json:"birthday"`
	KnownFor     string `json:"known_for_department"`
	PlaceOfBirth string `json:"place_of_birth"`
}
type TMDGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type TMDBMovieResponse struct {
	Results []TMDBMovie `json:"results"`
}
type TMDBGenreResponse struct {
	Genres []TMDGenre `json:"genres"`
}
type TMDBActorResponse struct {
	Cast []TMDBActor `json:"cast"`
}

func FetchSeedData() {
	Load()
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		log.Println("TMDB_API_KEY not set")
		return
	}

	var allMovies []TMDBMovie
	var allGenres = make(map[int]TMDGenre)
	var allActors = make(map[int]TMDBActor)
	var foundGenres = make(map[int]string)
	nextGenreID := 0
	nextActorID := 0

	//  --url 'https://api.themoviedb.org/3/genre/movie/list?language=en' \
	urlG := fmt.Sprintf("https://api.themoviedb.org/3/genre/movie/list?api_key=%s", apiKey)
	bodyG := ResponsData(urlG)
	var tmdbGenreResp TMDBGenreResponse
	if err := json.Unmarshal(bodyG, &tmdbGenreResp); err != nil {
		log.Printf("Failed to Unmarshal to tmdbGenreResp")
		return
	}

	for _, genre := range tmdbGenreResp.Genres {
		foundGenres[genre.ID] = genre.Name
	}

	for page := 1; page <= 50; page++ {
		//  --url 'https://api.themoviedb.org/3/movie/popular?language=en-US&page=1' \

		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d", apiKey, page)
		body := ResponsData(url)
		var tmdbMovieResp TMDBMovieResponse
		if err := json.Unmarshal(body, &tmdbMovieResp); err != nil {
			log.Printf("Failed to Unmarshal to tmdbMovieResp")
			return
		}
		for _, movie := range tmdbMovieResp.Results {
			movieGenres := []int{}
			movieActors := []int{}

			for _, genreId := range movie.GenreIDs {
				genreData, ok := allGenres[genreId]
				if !ok {
					genreData.ID = nextGenreID
					genrename := foundGenres[genreId]
					allGenres[genreId] = TMDGenre{nextGenreID, genrename}
					nextGenreID++
				}
				movieGenres = append(movieGenres, genreData.ID)
			}

			TMDBActorResp := fetchAllMovieActors(movie.ID, apiKey)
			if len(TMDBActorResp.Cast) == 0 {
				continue
			}
			for _, actor := range TMDBActorResp.Cast {
				actorData, ok := allActors[actor.ID]
				if !ok {
					person := fetchActorData(actor.ID, apiKey)
					if person == (Person{}) {
						continue
					}
					if person.KnownFor != "acting" {
						break
					}
					actorData = TMDBActor{
						ID:           nextActorID,
						Name:         actor.Name,
						OriginalName: actor.OriginalName,
						Character:    actor.Character,
						PersonInfo:   person,
					}
					allActors[actor.ID] = actorData
					nextActorID++
				}
				movieActors = append(movieActors, actorData.ID)
			}

			allMovies = append(allMovies, TMDBMovie{
				ID:               movie.ID,
				Title:            movie.Title,
				ReleaseDate:      movie.ReleaseDate,
				Overview:         movie.Overview,
				OriginalLanguage: movie.OriginalLanguage,
				GenreIDs:         movieGenres,
				ActorIDs:         movieActors,
			})
		}
		fmt.Printf("Page %d fetched, total so far: %d\n", page, len(allMovies))
	}

	fmt.Printf("Total movies fetched: %d\n", len(allMovies))

	movieOut, err := json.MarshalIndent(allMovies, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal to allMovies")
		return
	}

	var genreList []TMDGenre
	for _, g := range allGenres {
		genreList = append(genreList, g)
	}
	genreOut, err := json.MarshalIndent(genreList, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal to genreList")
		return
	}
	var actorList []TMDBActor
	for _, a := range allActors {
		actorList = append(actorList, a)
	}
	actorOut, err := json.MarshalIndent(actorList, "", " ")
	if err != nil {
		log.Printf("Failed to marshal to actorList")
		return
	}
	if movieErr := os.WriteFile("internal/database/data/tmdb_movies.json", movieOut, 0644); movieErr != nil {
		log.Printf("Failed to write movies data to json Error: %v", movieErr)
		return
	}
	if genreErr := os.WriteFile("internal/database/data/tmdb_genres.json", genreOut, 0644); genreErr != nil {
		log.Printf("Failed to write genre data to json Error: %v", genreErr)
		return
	}
	if actorErr := os.WriteFile("internal/database/data/tmdb_actors.json", actorOut, 0644); actorErr != nil {
		log.Printf("Failed to write actor data to json Error: %v", actorErr)
		return
	}
}
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func ResponsData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch from TMDB api: %v", url)
		return []byte{}
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

func fetchAllMovieActors(movieID int, apiKey string) TMDBActorResponse {
	// https: //api.themoviedb.org/3/movie/27205/credits?api_key=add63e21c755467e1763245167c231b4
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/credits?api_key=%s", movieID, apiKey)
	bodyG := ResponsData(url)
	var TMDBActorResp TMDBActorResponse
	if err := json.Unmarshal(bodyG, &TMDBActorResp); err != nil {
		log.Printf("Failed to Unmarshal to TMDBActorResp")
		return TMDBActorResponse{}
	}
	return TMDBActorResp
}

func fetchActorData(actorID int, apiKey string) Person {
	// https://api.themoviedb.org/3/person/6193?api_key=add63e21c755467e1763245167c231b4
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/%d?api_key=%s", actorID, apiKey)
	bodyG := ResponsData(url)
	var person Person
	if err := json.Unmarshal(bodyG, &person); err != nil {
		log.Printf("Failed to Unmarshal to TMDBActorResp")
		return Person{}
	}
	return person
}
