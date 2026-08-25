package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type TMDBMovie struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	ReleaseDate      string `json:"release_date"`
	Duration         uint16 `json:"duration"`
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
	nextGenreID := 1
	nextActorID := 1

	// we fetch all of the genres since there aren't many and save it in a map
	// with the original ganre ID as a key to map it to the movies fetched from the api
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

	// each page contains 20 movie so 10 page comes up to 200 movie to get less or more
	//adjust the page limit
	for page := 1; page <= 10; page++ {
		//  --url 'https://api.themoviedb.org/3/movie/popular?language=en-US&page=1' \

		url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&page=%d", apiKey, page)
		body := ResponsData(url)
		var tmdbMovieResp TMDBMovieResponse
		if err := json.Unmarshal(body, &tmdbMovieResp); err != nil {
			log.Printf("Failed to Unmarshal to tmdbMovieResp")
			return
		}
		for _, movie := range tmdbMovieResp.Results {
			seenActors := make(map[int]bool)
			movieGenres := []int{}
			movieActors := []int{}

			// this maped the genre with the movie using the genre.id
			// while also giving it a new incrementing id if its new to the allGenre which holds
			// all of the officaily registred id
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
					if person.KnownFor != "Acting" {
						continue
					}
					// the person data doesn't come with a birthdate sometime with can couse a unique contrain
					// for the name of the actor so we generate random birthdate with actors age 10 to 80
					if person.BirthDate == "" {
						year := time.Now().Year() - (rand.Intn(80) + 10)
						month := rand.Intn(12) + 1
						day := rand.Intn(28) + 1
						person.BirthDate = fmt.Sprintf("%04d-%02d-%02d", year, month, day)
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
				if seenActors[actorData.ID] {
					continue
				}
				seenActors[actorData.ID] = true
				movieActors = append(movieActors, actorData.ID)
			}
			yearString := strings.Split(movie.ReleaseDate, "-")

			allMovies = append(allMovies, TMDBMovie{
				ID:               movie.ID,
				Title:            movie.Title,
				ReleaseDate:      yearString[0],
				Duration:         uint16(rand.Uint32()),
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

	// we're soring the genre and movie map so that when we add it to the db table its in order we want
	// and we can use hte ids we have
	var genreList []TMDGenre
	for _, g := range allGenres {
		genreList = append(genreList, g)
	}
	sort.Slice(genreList, func(i, j int) bool {
		return genreList[i].ID < genreList[j].ID
	})
	genreOut, err := json.MarshalIndent(genreList, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal to genreList")
		return
	}
	var actorList []TMDBActor
	for _, a := range allActors {
		actorList = append(actorList, a)
	}
	sort.Slice(actorList, func(i, j int) bool {
		return actorList[i].ID < actorList[j].ID
	})
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
	var lastErr error
	for i := 0; i < 4; i++ {
		if i > 0 {
			time.Sleep(time.Duration(1) * 500 * time.Millisecond)
		}
		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			continue
		}

		body, bodyErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if bodyErr != nil {
			lastErr = bodyErr
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("status %d for %s", resp.StatusCode, url)
			continue
		}
		return body
	}
	log.Printf("Failed to fetch from TMDB api %s after retries: %v", url, lastErr)
	return []byte{}
}

func fetchAllMovieActors(movieID int, apiKey string) TMDBActorResponse {
	// https: //api.themoviedb.org/3/movie/27205/credits?api_key=?
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
	// https://api.themoviedb.org/3/person/6193?api_key=?
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/%d?api_key=%s", actorID, apiKey)
	bodyG := ResponsData(url)
	var person Person
	if err := json.Unmarshal(bodyG, &person); err != nil {
		log.Printf("Failed to Unmarshal to TMDBActorResp")
		return Person{}
	}
	return person
}
