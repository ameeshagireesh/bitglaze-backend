package main

import (
	"context"
	"net/http"
	"io"
	"github.com/go-chi/chi/v5"
)

type moviesResource struct{

}

//Deckaring Routes
func (rs moviesResource) Routes() chi.Router{
	r := chi.NewRouter()

	r.Post("/", rs.Create) // POST/movies - Create a new movie
	r.Get("/", rs.List)    // GET/movies - Read a list of movies

	r.Route("/{id}", func(r chi.Router){
		r.Use(MovieCtx)
		r.Get("/", rs.Get)   		// GET /movies/{id} - Reads a single movie with given :id
		r.Put("/", rs.Update)      //  PUT /movies/{id} - Updates a single movie with given :id
		r.Delete("/", rs.Delete)  //   DELETE /movies/{id} - Delets a single movie with given :id
	})

	return r
}

//Request Handler - POST/movies (Create a new movie) - C
func (rs moviesResource) Create(w http.ResponseWriter, r *http.Request){
	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Request Handler - GET/movies - Reads a list of movies - R
func (rs moviesResource) List(w http.ResponseWriter, r *http.Request){
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	
	if err!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

//To define routes in which id is passed
func MovieCtx(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))		
	})
}

//Request Handler - GET/movies/{id} - Reads a single movie with given :id - R
func (rs moviesResource) Get (w http.ResponseWriter, r*http.Request){
	id := r.Context().Value("id").(string)

	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/" + id)

	if err!= nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err!= nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Request Handler - PUT/movies/id - Updates a single movie with the given :id - U

func (rs moviesResource) Update ( w http.ResponseWriter, r *http.Request){
	id := r.Context().Value("id").(string)
	client := &http.Client{}

	req, err := http.NewRequest("PUT", "https://jsonplaceholder.typicode.com/posts/" + id, r.Body)
	req.Header.Add("Content-Type", "application/json")

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Request Handler - DELETE/movies/{id} - Deletes a single movie with given :id - D
func (rs moviesResource) Delete(w http.ResponseWriter, r *http.Request){
	id := r.Context().Value("id").(string)
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "https://jsonplaceholder.typicode.com/posts/" + id, nil)

	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)

	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	if _, err := io.Copy(w, resp.Body); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}