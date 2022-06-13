package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bitglaze Backend."))
	})

	r.Route("/movies", func(r chi.Router) {
		r.With(paginate).Get("/", ListMovies) // GET /articles
		r.Post("/", CreateMovie) // POST /articles

		r.Route("/{movieID}", func(r chi.Router) {
			r.Use(MovieCtx)            
			r.Get("/", GetMovie)       
			r.Put("/", UpdateMovie)    
			r.Delete("/", DeleteMovie) 
		})

	})

	http.ListenAndServe(":3333", r)
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, NewMovieListResponse(movies)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func MovieCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var movie *Movie
		var err error

		if movieID := chi.URLParam(r, "movieID"); movieID != "" {
			movie, err = dbGetMovie(movieID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "movie", movie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	data := &MovieRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	movie := data.Movie
	dbNewMovie(movie)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewMovieResponse(movie))
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value("movie").(*Movie)

	if err := render.Render(w, r, NewMovieResponse(movie)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	movie := r.Context().Value("movie").(*Movie)

	data := &MovieRequest{Movie: movie}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	movie = data.Movie
	dbUpdateMovie(movie.ID, movie)

	render.Render(w, r, NewMovieResponse(movie))
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	var err error
	movie := r.Context().Value("movie").(*Movie)

	movie, err = dbRemoveMovie(movie.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewMovieResponse(movie))
}

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type MovieRequest struct {
	*Movie
}

func (a *MovieRequest) Bind(r *http.Request) error {
	if a.Movie == nil {
		return errors.New("missing required Movie fields.")
	}
	a.Movie.MovieName = strings.ToLower(a.Movie.MovieName) 
	return nil
}

type MovieResponse struct {
	*Movie
}

func NewMovieResponse(movie *Movie) *MovieResponse {
	resp := &MovieResponse{Movie: movie}
	return resp
}

func (rd *MovieResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewMovieListResponse(movies []*Movie) []render.Renderer {
	list := []render.Renderer{}
	for _, movie := range movies {
		list = append(list, NewMovieResponse(movie))
	}
	return list
}

// Errors for different scenarios
type ErrResponse struct {
	Err            error `json:"-"` 
	HTTPStatusCode int   `json:"-"` 

	StatusText string `json:"status"`          
	AppCode    int64  `json:"code,omitempty"`  
	ErrorText  string `json:"error,omitempty"` 
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

type Movie struct {
	ID     string `json:"id"`
	MovieName string  `json:"movie_name"` 
	Genre  string `json:"genre"`
	Director   string `json:"director"`
}

// Article fixture data
var movies = []*Movie{
	{ID: "1", MovieName: "100", Genre: "Hi", Director: "hi"},
	{ID: "2", MovieName: "200", Genre: "sup", Director: "sup"},
	{ID: "3", MovieName: "300", Genre: "alo", Director: "alo"},
	{ID: "4", MovieName: "400", Genre: "bonjour", Director: "bonjour"},
}

func dbNewMovie(movie *Movie) (string, error) {
	movie.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	movies = append(movies, movie)
	return movie.ID, nil
}

func dbGetMovie(id string) (*Movie, error) {
	for _, a := range movies {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, errors.New("Movie not found.")
}


func dbUpdateMovie(id string, movie *Movie) (*Movie, error) {
	for i, a := range movies {
		if a.ID == id {
			movies[i] = movie
			return movie, nil
		}
	}
	return nil, errors.New("Movie not found.")
}

func dbRemoveMovie(id string) (*Movie, error) {
	for i, a := range movies {
		if a.ID == id {
			movies = append((movies)[:i], (movies)[i+1:]...)
			return a, nil
		}
	}
	return nil, errors.New("Movie not found.")
}