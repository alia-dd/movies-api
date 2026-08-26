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
git clone https://gitea.kood.tech/aliadaahirmohamed/movies-api.git
cd movies-api
```

### Prerequisites

- Go 1.25+
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
### API Key Authentication

All API endpoints are protected by an API key.

The API key is checked by the middleware before the request reaches the handler. Requests without a valid API key will be rejected.

### Generate an API Key

You can generate an API key using the `-g` flag:

```bash
go run ./cmd -g generate
```

The generated API key will be printed in the terminal.

### Using the API Key

Add the API key to your request using the `x-api-key` header:

For example:

```bash
curl -k -H "x-api-key: YOUR_API_KEY" "https://localhost:8800/api/actors"
```

## API Endpoints

### Movies

| Method | Endpoint | Description |
| --- | --- | --- |
| POST | `/api/movies` | Create a movie |
| GET | `/api/movies` | Get all movies with optional filters and pagination |
| GET | `/api/movies/{id}` | Get movie by ID or title |
| GET | `/api/movies/search?search=matrix` | Get movie by title |
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
| GET | `/api/actors/name?name=Leonardo` | Get actor by name |
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

You can test the API with either the `curl` examples below or the included Postman collection. `curl` is optional if you use Postman. To use Postman, import `Movie Database API generated.postman_collection.json`, start the server, and send the requests from the collection.

### Create a movie

```bash
curl -k -X POST https://localhost:8800/api/movies \
  -H "x-api-key: YOUR_API_KEY" \
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
### Update a movie

The API uses `PATCH` to partially update a movie. You only need to include the fields you want to change.

For example, to update the title and duration:

```bash
curl -k -X PATCH "https://localhost:8800/api/movies/1" \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception Updated",
    "duration": 150
  }'
```

You can also update the movie's actors and genres using their IDs:

```bash
curl -k -X PATCH "https://localhost:8800/api/movies/1" \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Inception Updated",
    "addActorIds": [10, 11],
    "removeActorIds": [5],
    "addGenreIds": [2],
    "removeGenreIds": [4]
  }'
```

Because this endpoint uses `PATCH`, you do not need to provide every movie field. Only the fields included in the request are changed.
### Update an Actor

The API uses `PATCH` for updating existing actors.

```bash
curl -k -X PATCH "https://localhost:8800/api/actors/1" \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Leonardo DiCaprio",
    "birthDate": "1974-11-11"
  }'
```

### Get all movies

```bash
curl -k "https://localhost:8800/api/movies?page=1&size=10"
```

### Create an actor

```bash
curl -k -X POST https://localhost:8800/api/actors \
  -H "x-api-key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Leonardo DiCaprio",
    "birthDate": "1974-11-11"
  }'
```

### Search genres

```bash
curl -k -H "x-api-key: YOUR_API_KEY" "https://localhost:8800/api/genres/search?q=drama&page=1&limit=10"
```

### Force-delete a movie

If `force` is omitted or set to `false`, the API performs a regular deletion and will not delete a movie that still has actor or genre associations:

```bash
curl -k -X DELETE -H "x-api-key: YOUR_API_KEY" "https://localhost:8800/api/movies/123"
curl -k -X DELETE -H "x-api-key: YOUR_API_KEY" "https://localhost:8800/api/movies/123?force=false"
```

Use `force=true` to delete a movie together with its actor and genre associations:

```bash
curl -k -X DELETE -H "x-api-key: YOUR_API_KEY" "https://localhost:8800/api/movies/123?force=true"
```

## Notes

- The API uses SQLite and creates its tables automatically on startup if they do not exist.
- Seed data is loaded from JSON files under `internal/database/data`.
- A Postman collection is included at `movie_api.postman_collection.json` for testing the endpoints quickly.