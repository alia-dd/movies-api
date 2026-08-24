# Movies API

A Go-based REST API for managing movies, actors, and genres with SQLite persistence. The project includes CRUD operations, filtering, pagination, and optional data seeding from local JSON files.

## Features

- Movie management with create, read, update, and delete endpoints
- Actor management with pagination and search by ID or name
- Genre management with search, filtering, and deletion controls
- SQLite database integration
- Seedable demo data from TMDB-style JSON files
- Middleware-based panic recovery for request safety

## Tech Stack

- Go
- SQLite via `github.com/mattn/go-sqlite3`
- Standard library HTTP server
- JSON-based request/response payloads

## Project Structure

```text
movies-api/
├── cmd/
│   ├── main.go
│   └── router.go
├── internal/
│   ├── database/
│   │   ├── data/
│   │   │   ├── tmdb_actors.json
│   │   │   ├── tmdb_genres.json
│   │   │   └── tmdb_movies.json
│   │   ├── db.go
│   │   ├── FetchSeedData.go
│   │   └── schema.go
│   ├── errors/
│   │   └── errors.go
│   ├── handler/
│   │   ├── actor_handler.go
│   │   ├── genre_handler.go
│   │   └── movie_handler.go
│   ├── models/
│   │   ├── actor.go
│   │   ├── genre.go
│   │   └── movie.go
│   ├── repository/
│   │   ├── actor_repository.go
│   │   ├── genre_repository.go
│   │   └── movie_repository.go
│   └── service/
│       ├── actor_service.go
│       ├── genre_service.go
│       └── movie_service.go
├── go.mod
├── movie_api.postman_collection.json
└── README.md
```

## Getting Started

### Clone the repository


```bash
git clone git@github.com:alia-dd/movies-api.git
cd movies-api
```

### Prerequisites

- Go 1.25+ (the project declares `go 1.25.7`)
- Git

### Install dependencies

```bash
go mod download
```

### Run the server

```bash
go run ./cmd
```

The API runs on:

```text
http://localhost:8800
```

### Seed the database

The server supports a `-s` flag to populate SQLite with seed data from the JSON files in `internal/database/data`.

```bash
go run ./cmd -s
```

> This will initialize the database and insert seeded actors, genres, and movies.

## API Endpoints

### Movies

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/movies` | Create a movie |
| GET | `/api/movies` | Get all movies with optional filters and pagination |
| GET | `/api/movies/{id}` | Get movie by ID or title |
| PATCH | `/api/movies/{id}` | Update a movie |
| DELETE | `/api/movies/{id}` | Delete a movie; use `?force=true` to delete its associations first |
| GET | `/api/movies/{movieId}/actors` | Get actors linked to a movie |
| GET | `/api/movies/{movieId}/genres` | Get genres linked to a movie |

### Actors

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/actors` | Create an actor |
| PATCH | `/api/actors/{id}` | Update an actor |
| GET | `/api/actors` | Get all actors with pagination |
| GET | `/api/actors/{id}` | Get actor by ID |
| GET | `/api/actors/name/{name}` | Get actor by name |
| DELETE | `/api/actors/{id}` | Delete actor by ID |
| DELETE | `/api/actors/name/{name}` | Delete actor by name |

### Genres

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/genres` | Create a genre |
| GET | `/api/genres` | Get all genres with pagination |
| GET | `/api/genres/search?q=...` | Search genres by name |
| GET | `/api/genres/{id}` | Get a genre by ID |
| PATCH | `/api/genres/{id}` | Update a genre |
| DELETE | `/api/genres/{id}` | Delete genre by ID |
| DELETE | `/api/genres/name/{name}` | Delete genre by name |

## Query Parameters

### Movie listing

```text
GET /api/movies?genre=Action&actor=Tom Hanks&year=1994&page=1&size=10
```

### Actor listing

```text
GET /api/actors?page=1&limit=10
```

### Genre search

```text
GET /api/genres/search?q=drama&page=1&limit=10
```

## Example Requests

You can test the API with either the `curl` examples below or the included Postman collection. `curl` is optional if you use Postman. To use Postman, import `movie_api.postman_collection.json`, start the server, and send the requests from the collection.

### Create a movie

```bash
curl -X POST http://localhost:8800/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception",
    "release_date": "2010",
    "duration": 148,
    "overview": "A thief who steals corporate secrets through dream-sharing technology.",
    "original_language": "en",
    "genre_ids": [1, 2],
    "Actor_ids": [10, 11]
  }'
```

### Get all movies

```bash
curl http://localhost:8800/api/movies?page=1&size=10
```

### Create an actor

```bash
curl -X POST http://localhost:8800/api/actors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Leonardo DiCaprio",
    "birthDate": "1974-11-11"
  }'
```

### Search genres

```bash
curl "http://localhost:8800/api/genres/search?q=drama&page=1&limit=10"
```

### Force-delete a movie

If `force` is omitted or set to `false`, the API performs a regular deletion and will not delete a movie that still has actor or genre associations:

```bash
curl -X DELETE "http://localhost:8800/api/movies/123"
curl -X DELETE "http://localhost:8800/api/movies/123?force=false"
```

Use `force=true` to delete a movie together with its actor and genre associations:

```bash
curl -X DELETE "http://localhost:8800/api/movies/123?force=true"
```

## Notes

- The API uses SQLite and creates its tables automatically on startup if they do not exist.
- Seed data is loaded from JSON files under `internal/database/data`.
- A Postman collection is included at `movie_api.postman_collection.json` for testing the endpoints quickly.
